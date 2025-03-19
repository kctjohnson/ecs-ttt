package game

import (
	"ttt/pkg/ecs"
)

func (g *Game) playerMovedEventHandler(event ecs.EventInterface) {
	g.world.Logger.Println("Player moved event received, toggling turn")
	g.isPlayer1Turn = !g.isPlayer1Turn
}

func (g *Game) playerWonEventHandler(event ecs.EventInterface) {
	player, _ := g.componentAccess.GetPlayerComponent(event.Entity())
	g.displayManager.ShowGameResult(player.Character + " won!")
	g.gameOver = true
}

func (g *Game) tieEventHandler(event ecs.EventInterface) {
	g.displayManager.ShowGameResult("It's a tie!")
	g.gameOver = true
}
