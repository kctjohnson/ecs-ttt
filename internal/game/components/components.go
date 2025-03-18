package components

import "ttt/pkg/ecs"

const (
	Board      ecs.ComponentType = "board"
	Player     ecs.ComponentType = "player"
	MoveIntent ecs.ComponentType = "move_intent"
)

type BoardComponent struct {
	ecs.Component
	Board [][]string
}

type PlayerComponent struct {
	ecs.Component
	Character string
}

type MoveIntentComponent struct {
	ecs.Component
	Row int
	Col int
}

var ComponentTypes = []ecs.ComponentType{
	Board,
	Player,
	MoveIntent,
}
