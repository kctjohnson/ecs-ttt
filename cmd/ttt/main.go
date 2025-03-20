package main

import (
	"flag"
	"log"
	"os"

	"ttt/internal/game"
	"ttt/internal/network"
)

func main() {
	// Parse command-line flags
	serverMode := flag.Bool("server", false, "Run in server mode")
	address := flag.String("address", ":8080", "Server address")
	flag.Parse()

	logger := log.New(os.Stdout, "TicTacToe: ", log.LstdFlags)

	if *serverMode {
		// Run in server mdoe
		server := network.NewGameServer(logger)

		// Create a networked game instance connected to the server
		game := game.NewGame(true, server)
		game.Initialize()

		// Run the game in a separate goroutine
		go game.Run()

		// Start the server
		logger.Fatal(server.Start(*address))
	} else {
		// Run in local game mode
		g := game.NewGame(false, nil)
		g.Initialize()
		g.Run()
	}
}
