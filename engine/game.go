package engine

import (
	"fmt"
	"log"
	"math/rand"
)

// Cell models the most basic unit of the grid
type Cell string

// Cells constant values
const (
	CellWall    Cell = "#"
	CellFloor   Cell = "_"
	CellUnknown Cell = "?"
)

// Game models the game rules and the state neccessary to enforce them
type Game struct {
	cells            []Cell
	rows             int
	cols             int
	wallRoots        int
	wallBuildingProb float64
	povRadius        int
	players          map[*Player]bool
}

// NewGame creates a new game
func NewGame() *Game {
	return &Game{
		players: make(map[*Player]bool),
	}
}

// Config models the game parameters
type Config struct {
	Rows             int
	Cols             int
	WallRoots        int
	WallBuildingProb float64
	POVRadius        int
}

// Config the game with sane defaults
func (game *Game) Config(config *Config) {
	game.rows = config.Rows
	log.Printf("[config] rows: %d\n", game.rows)

	game.cols = config.Cols
	log.Printf("[config] cols: %d\n", game.cols)

	game.wallRoots = config.WallRoots
	log.Printf("[config] wall roots: %d\n", game.wallRoots)

	game.wallBuildingProb = config.WallBuildingProb
	log.Printf("[config] wall building probability: %f\n", game.wallBuildingProb)

	game.povRadius = config.POVRadius
	log.Printf("[config] point of view radius: %d\n", game.povRadius)

	game.initializeCells()
	game.initializePlayers()
}

// AddPlayer inserts a player into the game
func (game *Game) AddPlayer(player *Player) {
	game.players[player] = true
}

// KillPlayer ignores the given player in the game
func (game *Game) KillPlayer(player *Player) {
	delete(game.players, player)
}

// Evaluate an action
func (game *Game) Evaluate(action *Action) {
	alive, ok := game.players[action.player]
	if !ok || !alive {
		return
	}
	switch action.Command {
	case CommandAttack:
		game.handleAttack(action.player)
	case CommandMoveBackward:
		game.handleMoveBackward(action.player)
	case CommandMoveForward:
		game.handleMoveForward(action.player)
	case CommandTurnLeft:
		game.handleTurnLeft(action.player)
	case CommandTurnRight:
		game.handleTurnRight(action.player)
	}
}

func (game *Game) handleAttack(player *Player) {
	row, col := game.fromIndex(player.Position)

	switch player.Direction {
	case DirectionEast:
		col++
	case DirectionNorth:
		row--
	case DirectionSouth:
		row++
	case DirectionWest:
		col--
	}

	if !game.inBounds(row, col) {
		return
	}

	target := game.toIndex(row, col)

	if game.cells[target] == CellWall {
		return
	}

	for enemy := range game.players {
		if alive := game.players[enemy]; !alive {
			continue
		}

		if enemy.Position == target {
			player.Score++
		}
	}
}

func (game *Game) handleMoveBackward(player *Player) {
	row, col := game.fromIndex(player.Position)

	switch player.Direction {
	case DirectionEast:
		col--
	case DirectionNorth:
		row++
	case DirectionSouth:
		row--
	case DirectionWest:
		col++
	}

	if !game.inBounds(row, col) {
		return
	}

	position := game.toIndex(row, col)

	if game.cells[position] != CellFloor {
		return
	}

	for enemy := range game.players {
		if enemy.Position == position {
			return
		}
	}

	player.Position = position
}

func (game *Game) handleMoveForward(player *Player) {
	row, col := game.fromIndex(player.Position)

	switch player.Direction {
	case DirectionEast:
		col++
	case DirectionNorth:
		row--
	case DirectionSouth:
		row++
	case DirectionWest:
		col--
	}

	if !game.inBounds(row, col) {
		return
	}

	position := game.toIndex(row, col)

	if game.cells[position] != CellFloor {
		return
	}

	for enemy := range game.players {
		if enemy.Position == position {
			return
		}
	}

	player.Position = position
}

func (game *Game) handleTurnLeft(player *Player) {
	switch player.Direction {
	case DirectionEast:
		player.Direction = DirectionNorth
	case DirectionNorth:
		player.Direction = DirectionWest
	case DirectionSouth:
		player.Direction = DirectionEast
	case DirectionWest:
		player.Direction = DirectionSouth
	}
}

func (game *Game) handleTurnRight(player *Player) {
	switch player.Direction {
	case DirectionEast:
		player.Direction = DirectionSouth
	case DirectionNorth:
		player.Direction = DirectionEast
	case DirectionSouth:
		player.Direction = DirectionWest
	case DirectionWest:
		player.Direction = DirectionNorth
	}
}

// State of the game from a players POV
func (game *Game) State(player *Player) *State {
	state := &State{
		Grid: &Grid{
			Rows:  game.rows,
			Cols:  game.cols,
			Cells: game.cellsWithPOV(player),
		},
		Player:  player,
		Enemies: game.enemiesWithPOV(player),
	}

	return state
}

func (game *Game) cellsWithPOV(player *Player) []Cell {
	cells := make([]Cell, len(game.cells))

	for i := range game.cells {
		if game.inPOV(player, i) {
			cells[i] = game.cells[i]
		} else {
			cells[i] = CellUnknown
		}
	}

	return cells
}

func (game *Game) enemiesWithPOV(player *Player) []*Player {
	enemies := []*Player{}

	for enemy := range game.players {
		if enemy == player {
			continue
		}

		if game.inPOV(player, enemy.Position) {
			enemies = append(enemies, enemy)
		}
	}

	return enemies
}

// State of the game at a particular moment from a players POV
type State struct {
	Grid    *Grid     `json:"grid"`
	Player  *Player   `json:"player"`
	Enemies []*Player `json:"enemies"`
}

// Grid models the game board
type Grid struct {
	Rows  int    `json:"rows"`
	Cols  int    `json:"cols"`
	Cells []Cell `json:"cells"`
}

func (game *Game) initializeCells() {
	log.Printf("[config] initializing cells\n")
	var row, col int
	size := game.rows * game.cols
	cells := make([]Cell, size)

	// let's build the map in layers

	// 1. every cell is a floor
	for i := 0; i < size; i++ {
		cells[i] = CellFloor
	}

	// 2. build walls in the map borders
	for i := 0; i < size; i++ {
		row, col = game.fromIndex(i)

		if col == 0 || col == game.cols-1 || row == 0 || row == game.rows-1 {
			cells[i] = CellWall
		}
	}

	// 3. select random points to start building more walls
	for w := 0; w < game.wallRoots; w++ {
		row = rand.Intn(game.rows)
		col = rand.Intn(game.cols)

		for {
			cells[game.toIndex(row, col)] = CellWall
			row += rand.Intn(3) - 1
			col += rand.Intn(3) - 1

			if !game.inBounds(row, col) {
				break
			}

			if rand.Float64() > game.wallBuildingProb {
				break
			}
		}
	}

	game.cells = cells
}

func (game *Game) initializePlayers() {
	log.Printf("[config] initializing players\n")
	occupiedPositions := []int{}
	for player := range game.players {
		log.Printf("[config] initializing player (%p)\n", player)
		var position int
		for {
			position = game.randomCellIndex(CellFloor)
			if !intSliceMember(position, occupiedPositions) {
				break
			}
			log.Printf("[config] failed player (%p) positioning attempt: %d already occupied\n",
				player, position)
		}
		player.Position = position
		log.Printf("[config] player (%p) position: %d\n", player, player.Position)

		player.Direction = game.randomDirection()
		log.Printf("[config] player (%p) direction: %s\n", player, player.Direction)

		player.Score = 0
		log.Printf("[config] player (%p) score: %d\n", player, player.Score)
	}
}

func (game *Game) randomDirection() Direction {
	directions := [...]Direction{DirectionEast, DirectionNorth, DirectionSouth, DirectionWest}
	return directions[rand.Intn(len(directions))]
}

func (game *Game) randomCellIndex(cell Cell) int {
	var probe int
	for {
		probe = rand.Intn(len(game.cells))
		fmt.Printf("[config] random cell position probe: %d\n", probe)
		if game.cells[probe] == cell {
			return probe
		}
	}
}

func (game *Game) inBounds(row, col int) bool {
	return row >= 0 && row < game.rows && col >= 0 && col < game.cols
}

func (game *Game) indexInBounds(i int) bool {
	return i >= 0 && i < (game.rows*game.cols)
}

func (game *Game) fromIndex(i int) (row, col int) {
	return i / game.cols, i % game.cols
}

func (game *Game) toIndex(row, col int) (i int) {
	return row*game.cols + col
}

func (game *Game) inPOV(player *Player, position int) bool {
	pRow, pCol := game.fromIndex(player.Position)
	row, col := game.fromIndex(position)
	rowDiff := row - pRow
	colDiff := col - pCol
	radius := game.povRadius

	return (rowDiff*rowDiff + colDiff*colDiff) <= (radius * radius)
}

func intSliceMember(elm int, elms []int) bool {
	for i := range elms {
		if elm == elms[i] {
			return true
		}
	}
	return false
}
