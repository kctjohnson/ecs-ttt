package systems

import (
	"ttt/internal/game/components"
	"ttt/internal/game/events"
	"ttt/pkg/ecs"
)

// MoveSystem is responsible for evaluating and executing player moves
type MoveSystem struct {
	ComponentAccess *components.ComponentAccess
}

func (m *MoveSystem) Update(world *ecs.World) {
	// Get all entities with a move intent component
	moveIntentEnts := world.ComponentManager.GetAllEntitiesWithComponent(components.MoveIntent)
	if len(moveIntentEnts) == 0 {
		return
	}

	// Get the board entity
	boardEnts := world.ComponentManager.GetAllEntitiesWithComponent(components.Board)
	if len(boardEnts) == 0 {
		return
	}

	// Get the board component
	board, hasBoardComp := m.ComponentAccess.GetBoardComponent(boardEnts[0])
	if !hasBoardComp {
		return
	}

	for _, entity := range moveIntentEnts {
		// Get the move intent component
		moveIntent, _ := m.ComponentAccess.GetMoveIntentComponent(entity)

		// Check if the move is valid
		if board.Board[moveIntent.Row][moveIntent.Col] != components.Empty {
			// Invalid move, remove the move intent component
			world.ComponentManager.RemoveComponent(entity, components.MoveIntent)
			continue
		} else {
			// Get the player component
			player, hasPlayerComp := m.ComponentAccess.GetPlayerComponent(entity)
			if !hasPlayerComp {
				continue
			}

			// Update the board
			board.Board[moveIntent.Row][moveIntent.Col] = player.CellState

			// Remove the move intent component
			world.ComponentManager.RemoveComponent(entity, components.MoveIntent)

			// Send out events
			world.QueueEvent(events.PlayerMovedEvent{
				Ent: entity,
				Row: moveIntent.Row,
				Col: moveIntent.Col,
			})
		}
	}

}
