package components

import "ttt/pkg/ecs"

const (
	Board      ecs.ComponentType = "board"
	Player     ecs.ComponentType = "player"
	MoveIntent ecs.ComponentType = "move_intent"
)

type BoardComponent struct {
	Board [][]string
}

func (c BoardComponent) IsComponent() {}
func (c BoardComponent) GetType() ecs.ComponentType {
	return Board
}

type PlayerComponent struct {
	ecs.Component
	Character string
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
