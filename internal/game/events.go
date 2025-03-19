package game

import (
	"ttt/internal/game/components"
	"ttt/pkg/ecs"
)

func (g *Game) playerMovedEventHandler(event ecs.EventInterface) {
	gameState := g.getGameState()
	if gameState == nil {
		return
	}

	// Toggle the turn
	playerEnts := g.world.ComponentManager.GetAllEntitiesWithComponent(components.Player)
	if len(playerEnts) != 2 {
		return
	}

	if gameState.PlayerTurn == playerEnts[0] {
		gameState.PlayerTurn = playerEnts[1]
	} else {
		gameState.PlayerTurn = playerEnts[0]
	}
}

func (g *Game) playerWonEventHandler(event ecs.EventInterface) {
	gameState := g.getGameState()
	if gameState != nil {
		gameState.GameOver = true
	}

	player, _ := g.componentAccess.GetPlayerComponent(event.Entity())
	g.displayManager.ShowGameResult(player.Character + " won!")
}

func (g *Game) tieEventHandler(event ecs.EventInterface) {
	gameState := g.getGameState()
	if gameState != nil {
		gameState.GameOver = true
	}

	g.displayManager.ShowGameResult("It's a tie!")
}
