package components

import "ttt/pkg/ecs"

const (
	GameState  ecs.ComponentType = "game_state"
	Board      ecs.ComponentType = "board"
	Player     ecs.ComponentType = "player"
	MoveIntent ecs.ComponentType = "move_intent"
)

type GameStateComponent struct {
	ecs.Component
	PlayerTurn ecs.Entity
	GameOver   bool
}

func (c GameStateComponent) IsComponent() {}
func (c GameStateComponent) GetType() ecs.ComponentType {
	return GameState
}

type CellState int

const (
	Empty CellState = iota
	Player1
	Player2
)

type BoardComponent struct {
	Board [][]CellState
}

func (c BoardComponent) IsComponent() {}
func (c BoardComponent) GetType() ecs.ComponentType {
	return Board
}

type PlayerComponent struct {
	ecs.Component
	Character string
	CellState CellState
}

func (c PlayerComponent) IsComponent() {}
func (c PlayerComponent) GetType() ecs.ComponentType {
	return Player
}

type MoveIntentComponent struct {
	ecs.Component
	Row int
	Col int
}

func (c MoveIntentComponent) IsComponent() {}
func (c MoveIntentComponent) GetType() ecs.ComponentType {
	return MoveIntent
}

var ComponentTypes = []ecs.ComponentType{
	Board,
	Player,
	MoveIntent,
}
