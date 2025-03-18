package game

import (
	"ttt/internal/game/components"
	"ttt/pkg/ecs"
)

func (g *Game) playerMovedEventHandler(event ecs.Event) {
	g.isPlayer1Turn = !g.isPlayer1Turn
}

func (g *Game) playerWonEventHandler(event ecs.Event) {
	playerComp, _ := g.world.ComponentManager.GetComponent(event.Entity, components.Player)
	player := playerComp.(*components.PlayerComponent)
	g.world.Logger.Printf("%s won!", player.Character)
	g.gameOver = true
}

func (g *Game) tieEventHandler(event ecs.Event) {
	g.world.Logger.Printf("It's a tie!")
	g.gameOver = true
}
