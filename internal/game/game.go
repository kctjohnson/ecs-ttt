package game

import (
	"log"
	"os"

	"ttt/internal/game/components"
	"ttt/internal/game/events"
	"ttt/internal/game/systems"
	"ttt/internal/game/ui/console"
	"ttt/pkg/ecs"
)

type Game struct {
	world           *ecs.World
	inputManager    console.ConsoleInputManager
	displayManager  console.ConsoleDisplayManager
	componentAccess *components.ComponentAccess
	player1         ecs.Entity
	player2         ecs.Entity
	gameOver        bool
	isPlayer1Turn   bool
}

func NewGame() *Game {
	logger := log.New(os.Stdout, "TicTacToe: ", log.LstdFlags)

	world := ecs.NewWorld(logger)

	// Create the component access manager
	componentAccess := components.NewComponentAccess(world)

	// Register core ECS systems
	world.AddSystem(&systems.MoveSystem{
		ComponentAccess: componentAccess,
	})
	world.AddSystem(&systems.BoardSystem{
		ComponentAccess: componentAccess,
	})

	return &Game{
		world:           world,
		inputManager:    console.NewConsoleInputManager(),
		displayManager:  console.NewConsoleDisplayManager(),
		componentAccess: componentAccess,
		gameOver:        false,
		isPlayer1Turn:   true,
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

	g.player1 = player1

	// Make the player 2 entity
	player2 := g.world.EntityManager.CreateEntity()
	g.world.ComponentManager.AddComponent(
		player2,
		components.Player,
		&components.PlayerComponent{Character: "O"},
	)

	g.player2 = player2

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
		var playerEnt ecs.Entity
		if g.isPlayer1Turn {
			playerEnt = g.player1
		} else {
			playerEnt = g.player2
		}

		// Get the player component
		player, _ := g.componentAccess.GetPlayerComponent(playerEnt)

		g.displayManager.ShowTurnPrompt(player.Character)
		row, col, valid := g.inputManager.GetPlayerMove()

		if !valid {
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

	board, hasBoardComp := g.componentAccess.GetBoardComponent(boardEnts[0])
	if !hasBoardComp {
		return
	}

	g.displayManager.ShowBoard(board.Board)
}
