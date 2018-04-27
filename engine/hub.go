package engine

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	sendStateDuration = 2 * time.Second
)

// Hub models the mechanism to coordinate the game server connections.
type Hub struct {
	// Registered players
	players map[*Player]bool

	// Player actions
	actions chan *Action

	// Request for registering a player
	register chan *Player

	// Request for unregistering a player
	unregister chan *Player

	// Request for starting (true) or stopping (false) the game
	Playing chan *Config

	// Wether the game is being played or not
	isPlaying bool

	// Game being played
	game *Game
}

// NewHub creates a new connection hub
func NewHub() *Hub {
	return &Hub{
		players:    make(map[*Player]bool),
		actions:    make(chan *Action),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		Playing:    make(chan *Config),
		isPlaying:  false,
		game:       NewGame(),
	}
}

// Run executes the main game loop
func (hub *Hub) Run() {
	log.Printf("executing connection hub (%p) main loop\n", hub)
	for {
		select {
		case player := <-hub.register:
			log.Printf("player (%p) request registration\n", player)
			if hub.isPlaying {
				close(player.send)
				log.Printf("rejecting player (%p) registration request\n", player)
			} else {
				hub.players[player] = true
				hub.game.AddPlayer(player)
				log.Printf("player (%p) registered\n", player)
			}
		case player := <-hub.unregister:
			log.Printf("player (%p) request unregistration\n", player)
			if _, ok := hub.players[player]; ok {
				close(player.send)
				delete(hub.players, player)
				hub.game.KillPlayer(player)
				log.Printf("player (%p) unregistered\n", player)
				if len(hub.players) == 0 {
					hub.isPlaying = false
				}
			}
		case config := <-hub.Playing:
			log.Printf("handling playing=%v\n", config)
			switch {
			case config != nil && !hub.isPlaying:
				hub.isPlaying = true
				log.Printf("configuring game\n")
				hub.game.Config(config)
				log.Printf("sending game state to players\n")
				for player := range hub.players {
					log.Printf("sending game state to player (%p)\n", player)
					select {
					case player.send <- hub.game.State(player):
					case <-time.After(sendStateDuration):
						log.Printf("could not send game state to player (%p)\n", player)
						close(player.send)
						delete(hub.players, player)
						hub.game.KillPlayer(player)
					}
				}
			case config == nil && hub.isPlaying:
				log.Printf("shutting down game\n")
				hub.isPlaying = false
				for player := range hub.players {
					log.Printf("shutting down player (%p)\n", player)
					close(player.send)
					delete(hub.players, player)
					hub.game.KillPlayer(player)
				}
			}
		case action := <-hub.actions:
			log.Printf("action (%p) received\n", action)
			if !hub.isPlaying {
				log.Printf("ignoring action (%p)\n", action)
			} else {
				log.Printf("evaluating action (%p)\n", action)
				hub.game.Evaluate(action)
				log.Printf("sending game state to players\n")
				for player := range hub.players {
					log.Printf("sending game state to player (%p)\n", player)
					select {
					case player.send <- hub.game.State(player):
					case <-time.After(sendStateDuration):
						log.Printf("could not send game state to player (%p)\n", player)
						close(player.send)
						delete(hub.players, player)
						hub.game.KillPlayer(player)
					}
				}
			}
		}
	}
}

// ConnectPlayer returns an http handler function for
func (hub *Hub) ConnectPlayer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		color := r.FormValue("color")
		if name == "" || color == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(nil)
			return
		}

		log.Printf("upgrading from HTTP to WebSocket protocol\n")
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("[error] upgrading to websocket:", err)
			return
		}

		log.Printf("creating new player: %s (%s)\n", name, color)
		player := NewPlayer(name, color)
		hub.register <- player

		// pump the game state into the websocket
		go hub.pumpIntoPlayer(player, ws)

		// pump the player actions out of the websocket
		go hub.pumpOutOfPlayer(player, ws)
	}
}

func (hub *Hub) pumpIntoPlayer(player *Player, ws *websocket.Conn) {
	ticker := time.NewTicker(pingTickPeriod)

	defer func() {
		log.Printf("shutting down player (%p) input communication\n", player)
		ticker.Stop()
		_ = ws.Close()
	}()

	log.Printf("executing player (%p) input connection loop\n", player)
	for {
		select {
		case state, ok := <-player.send:
			_ = ws.SetWriteDeadline(time.Now().Add(writeWaitPeriod))

			if !ok {
				log.Printf("player (%p) input communication was closed\n", player)
				_ = ws.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			log.Printf("sending player (%p) a game state\n", player)
			if err := ws.WriteJSON(state); err != nil {
				log.Printf("[error] problem writing JSON to player (%p): %v\n", player, err)
				return
			}
		case <-ticker.C:
			log.Printf("testing player (%p) connection\n", player)
			_ = ws.SetWriteDeadline(time.Now().Add(writeWaitPeriod))
			if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[error] problem testing player (%p) connection: %v\n", player, err)
				return
			}
		}
	}
}

func (hub *Hub) pumpOutOfPlayer(player *Player, ws *websocket.Conn) {
	var action *Action

	defer func() {
		log.Printf("shutting down player (%p) output communication\n", player)
		hub.unregister <- player
		_ = ws.Close()
	}()

	_ = ws.SetReadDeadline(time.Now().Add(readWaitPeriod))
	ws.SetPongHandler(func(string) error {
		_ = ws.SetReadDeadline(time.Now().Add(readWaitPeriod))
		log.Printf("pong received")
		return nil
	})

	log.Printf("executing player (%p) output connection loop\n", player)
	for {
		if err := ws.ReadJSON(&action); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[error] player (%p) output communication was unexpectedly close: %v\n", player, err)
			}
			log.Printf("player (%p) output communication was closed\n", player)
			return
		}
		log.Printf("player (%p) action was read\n", player)

		action.player = player
		hub.actions <- action
	}
}
