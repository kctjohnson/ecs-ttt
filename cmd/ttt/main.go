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
		server := network.NewGameServer(logger)
		logger.Fatal(server.Start(*address))
	} else {
		g := game.NewGame()
		g.Initialize()
		g.Run()
	}
}
