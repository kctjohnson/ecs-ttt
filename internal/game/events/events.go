package events

import "ttt/pkg/ecs"

const (
	PlayerMoved ecs.EventType = "player_moved"
	PlayerWon   ecs.EventType = "player_won"
	Tie         ecs.EventType = "tie"
)
