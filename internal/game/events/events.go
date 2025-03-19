package events

import "ttt/pkg/ecs"

const (
	PlayerMoved ecs.EventType = "player_moved"
	PlayerWon   ecs.EventType = "player_won"
	Tie         ecs.EventType = "tie"
)

type PlayerMovedEvent struct {
	Ent      ecs.Entity
	Row, Col int
}

func (e PlayerMovedEvent) Type() ecs.EventType {
	return PlayerMoved
}

func (e PlayerMovedEvent) Entity() ecs.Entity {
	return e.Ent
}

func (e PlayerMovedEvent) Data() any {
	return map[string]int{"row": e.Row, "col": e.Col}
}

type PlayerWonEvent struct {
	Ent ecs.Entity
}

func (e PlayerWonEvent) Type() ecs.EventType {
	return PlayerWon
}

func (e PlayerWonEvent) Entity() ecs.Entity {
	return e.Ent
}

func (e PlayerWonEvent) Data() any {
	return nil
}

type TieEvent struct {
	Ent ecs.Entity
}

func (e TieEvent) Type() ecs.EventType {
	return Tie
}

func (e TieEvent) Entity() ecs.Entity {
	return e.Ent
}

func (e TieEvent) Data() any {
	return nil
}
