package components

import "ttt/pkg/ecs"

const (
	NetworkIdentity ecs.ComponentType = "network_identity"
	GameSession     ecs.ComponentType = "game_session"
)

// NetworkIdentityComponent associates an entity with a client
type NetworkIdentityComponent struct {
	ecs.Component
	ClientID string
}

func (c NetworkIdentityComponent) GetType() ecs.ComponentType {
	return NetworkIdentity
}

// GameSessionComponent tracks an active game session
type GameSessionComponent struct {
	ecs.Component
	GameID      string
	Player1     ecs.Entity // Entity with Player component
	Player2     ecs.Entity // Entity with Player component
	BoardEntity ecs.Entity // Entity with Board component
	IsActive    bool
}

func (c GameSessionComponent) GetType() ecs.ComponentType {
	return GameSession
}
