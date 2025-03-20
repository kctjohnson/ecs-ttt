package components

import (
	"ttt/pkg/ecs"
)

type ComponentAccess struct {
	world *ecs.World
}

func NewComponentAccess(world *ecs.World) *ComponentAccess {
	return &ComponentAccess{world: world}
}

func (ca *ComponentAccess) GetGameStateComponent(entity ecs.Entity) (*GameStateComponent, bool) {
	component, found := ca.world.ComponentManager.GetComponent(entity, GameState)
	if !found {
		return nil, false
	}
	return component.(*GameStateComponent), true
}

func (ca *ComponentAccess) GetBoardComponent(entity ecs.Entity) (*BoardComponent, bool) {
	component, found := ca.world.ComponentManager.GetComponent(entity, Board)
	if !found {
		return nil, false
	}
	return component.(*BoardComponent), true
}

func (ca *ComponentAccess) GetPlayerComponent(
	entity ecs.Entity,
) (*PlayerComponent, bool) {
	component, found := ca.world.ComponentManager.GetComponent(entity, Player)
	if !found {
		return nil, false
	}
	return component.(*PlayerComponent), true
}

func (ca *ComponentAccess) GetMoveIntentComponent(
	entity ecs.Entity,
) (*MoveIntentComponent, bool) {
	component, found := ca.world.ComponentManager.GetComponent(entity, MoveIntent)
	if !found {
		return nil, false
	}
	return component.(*MoveIntentComponent), true
}

func (ca *ComponentAccess) GetNetworkIdentityComponent(
	entity ecs.Entity,
) (*NetworkIdentityComponent, bool) {
	component, found := ca.world.ComponentManager.GetComponent(entity, NetworkIdentity)
	if !found {
		return nil, false
	}
	return component.(*NetworkIdentityComponent), true
}

func (ca *ComponentAccess) GetGameSessionComponent(
	entity ecs.Entity,
) (*GameSessionComponent, bool) {
	component, found := ca.world.ComponentManager.GetComponent(entity, GameSession)
	if !found {
		return nil, false
	}
	return component.(*GameSessionComponent), true
}
