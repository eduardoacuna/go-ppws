package engine

// Player models a user playing the game
type Player struct {
	send chan *State

	Position  int       `json:"position"`
	Direction Direction `json:"direction"`
	Score     int       `json:"score"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
}

// Direction models the orientation a player is heading to
type Direction string

// Directions constant values
const (
	DirectionNorth Direction = "north"
	DirectionSouth Direction = "south"
	DirectionEast  Direction = "east"
	DirectionWest  Direction = "west"
)

// Action models a players intent in the game
type Action struct {
	player *Player

	Command Command `json:"command"`
}

// Command models the type of actions a player can execute
type Command string

// Commands constant values
const (
	CommandMoveForward  Command = "move-forward"
	CommandMoveBackward Command = "move-backward"
	CommandTurnRight    Command = "turn-right"
	CommandTurnLeft     Command = "turn-left"
	CommandAttack       Command = "attack"
)

// NewPlayer creates a new player
func NewPlayer(name, color string) *Player {
	return &Player{
		send:  make(chan *State),
		Name:  name,
		Color: color,
	}
}
