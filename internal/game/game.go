package game

import (
	"fmt"
	"log"

	"ttt/internal/game/components"
	"ttt/internal/game/events"
	"ttt/internal/game/systems"
	"ttt/pkg/ecs"
)

type Game struct {
	world         *ecs.World
	gameOver      bool
	isPlayer1Turn bool
}

func NewGame() *Game {
	world := ecs.NewWorld(log.Default())

	// Register core ECS systems
	world.AddSystem(&systems.MoveSystem{})
	world.AddSystem(&systems.BoardSystem{})

	return &Game{
		world:         world,
		gameOver:      false,
		isPlayer1Turn: true,
	}
}

func (g *Game) Initialize() {
	// Register component types
	g.registerComponentTypes()

	// Register event handlers
	g.world.RegisterEventHandler(events.PlayerMoved, g.playerMovedEventHandler)
	g.world.RegisterEventHandler(events.PlayerWon, g.playerWonEventHandler)
	g.world.RegisterEventHandler(events.Tie, g.tieEventHandler)

	// Make the player 1 entity
	player1 := g.world.EntityManager.CreateEntity()
	g.world.ComponentManager.AddComponent(
		player1,
		components.Player,
		&components.PlayerComponent{Character: "X"},
	)

	// Make the player 2 entity
	player2 := g.world.EntityManager.CreateEntity()
	g.world.ComponentManager.AddComponent(
		player2,
		components.Player,
		&components.PlayerComponent{Character: "O"},
	)

	// Make the board entity
	boardTiles := make([][]string, 3)
	for i := range boardTiles {
		boardTiles[i] = make([]string, 3)
	}

	board := g.world.EntityManager.CreateEntity()
	g.world.ComponentManager.AddComponent(
		board,
		components.Board,
		&components.BoardComponent{
			Board: boardTiles,
		},
	)
}

func (g *Game) registerComponentTypes() {
	// Register all component types with the component manager
	for _, componentType := range components.ComponentTypes {
		g.world.ComponentManager.RegisterComponentType(componentType)
	}
}

func (g *Game) Run() {
	g.world.Logger.Println("Starting game...")

	// Main game loop
	for !g.gameOver {
		// Display the board
		g.displayBoard()

		// Get the player entity
		playerEnts := g.world.ComponentManager.GetAllEntitiesWithComponent(components.Player)
		if len(playerEnts) != 2 {
			g.world.Logger.Println("Error: Missing player components")
			continue
		}
		var playerEnt ecs.Entity
		if g.isPlayer1Turn {
			playerEnt = playerEnts[0]
		} else {
			playerEnt = playerEnts[1]
		}

		// Get the player component
		playerComp, _ := g.world.ComponentManager.GetComponent(
			playerEnt,
			components.Player,
		)
		player := playerComp.(*components.PlayerComponent)

		g.world.Logger.Printf(
			"%s, please enter column and row (0-2) separated by a space:",
			player.Character,
		)
		var row, col int
		_, err := fmt.Scanf("%d %d", &col, &row)
		if err != nil {
			g.world.Logger.Println("Invalid input. Please try again.")
			continue
		}

		if row < 0 || row > 2 || col < 0 || col > 2 {
			g.world.Logger.Println("Invalid input. Please try again.")
			continue
		}

		// Create a move intent component
		moveIntent := &components.MoveIntentComponent{
			Row: row,
			Col: col,
		}

		// Add the move intent component to the player entity
		g.world.ComponentManager.AddComponent(
			playerEnt,
			components.MoveIntent,
			moveIntent,
		)

		g.world.Update()
	}
}

func (g Game) displayBoard() {
	boardEnts := g.world.ComponentManager.GetAllEntitiesWithComponent(components.Board)
	if len(boardEnts) == 0 {
		return
	}
	boardComp, hasBoardComp := g.world.ComponentManager.GetComponent(boardEnts[0], components.Board)
	if !hasBoardComp {
		return
	}
	board := boardComp.(*components.BoardComponent)

	fmt.Println()
	for y := range 3 {
		for x := range 3 {
			if board.Board[y][x] == "" {
				fmt.Print(".")
			} else {
				fmt.Print(board.Board[y][x])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}
