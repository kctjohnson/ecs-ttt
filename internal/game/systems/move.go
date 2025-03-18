package systems

import (
	"ttt/internal/game/components"
	"ttt/internal/game/events"
	"ttt/pkg/ecs"
)

// MoveSystem is responsible for evaluating and executing player moves
type MoveSystem struct{}

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
	boardEnt := boardEnts[0]
	boardComp, hasBoardComp := world.ComponentManager.GetComponent(boardEnt, components.Board)
	if !hasBoardComp {
		return
	}

	for _, entity := range moveIntentEnts {
		// Get the move intent component
		moveIntentComp, _ := world.ComponentManager.GetComponent(
			entity,
			components.MoveIntent,
		)
		moveIntent := moveIntentComp.(*components.MoveIntentComponent)

		// Check if the move is valid
		if boardComp.(*components.BoardComponent).Board[moveIntent.Row][moveIntent.Col] != "" {
			// Invalid move, remove the move intent component
			world.ComponentManager.RemoveComponent(entity, components.MoveIntent)
			continue
		} else {
			// Get the player component
			playerComp, hasPlayerComp := world.ComponentManager.GetComponent(entity, components.Player)
			if !hasPlayerComp {
				continue
			}
			player := playerComp.(*components.PlayerComponent)

			// Update the board
			boardComp.(*components.BoardComponent).Board[moveIntent.Row][moveIntent.Col] = player.Character

			// Remove the move intent component
			world.ComponentManager.RemoveComponent(entity, components.MoveIntent)

			// Send out events
			world.QueueEvent(events.PlayerMoved, entity, map[string]any{
				"row": moveIntent.Row,
				"col": moveIntent.Col,
			})
		}
	}

}
