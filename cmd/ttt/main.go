package main

import "ttt/internal/game"

func main() {
	g := game.NewGame()
	g.Initialize()
	g.Run()
}
