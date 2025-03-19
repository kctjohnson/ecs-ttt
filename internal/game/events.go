package game

import (
	"ttt/pkg/ecs"
)

func (g *Game) playerMovedEventHandler(event ecs.Event) {
	g.isPlayer1Turn = !g.isPlayer1Turn
}

func (g *Game) playerWonEventHandler(event ecs.Event) {
	player, _ := g.componentAccess.GetPlayerComponent(event.Entity)
	g.displayManager.ShowGameResult(player.Character + " won!")
	g.gameOver = true
}

func (g *Game) tieEventHandler(event ecs.Event) {
	g.displayManager.ShowGameResult("It's a tie!")
	g.gameOver = true
}
