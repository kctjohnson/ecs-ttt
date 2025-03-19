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
		&components.PlayerComponent{
			Character: "X",
			CellState: components.Player1,
		},
	)

	// Make the player 2 entity
	player2 := g.world.EntityManager.CreateEntity()
	g.world.ComponentManager.AddComponent(
		player2,
		components.Player,
		&components.PlayerComponent{
			Character: "O",
			CellState: components.Player2,
		},
	)

	// Make the board entity
	boardTiles := make([][]components.CellState, 3)
	for i := range boardTiles {
		boardTiles[i] = make([]components.CellState, 3)
		for j := range boardTiles[i] {
			boardTiles[i][j] = components.Empty
		}
	}

	board := g.world.EntityManager.CreateEntity()
	g.world.ComponentManager.AddComponent(
		board,
		components.Board,
		&components.BoardComponent{
			Board: boardTiles,
		},
	)

	// Make the game state entity
	gameState := g.world.EntityManager.CreateEntity()
	g.world.ComponentManager.AddComponent(
		gameState,
		components.GameState,
		&components.GameStateComponent{
			PlayerTurn: player1,
			GameOver:   false,
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
	for {
		// Get the game state entity
		gameState := g.getGameState()
		if gameState == nil || gameState.GameOver {
			break
		}

		// Display the board
		g.displayBoard()

		// Get the player component
		playerEnt := gameState.PlayerTurn
		player, _ := g.componentAccess.GetPlayerComponent(gameState.PlayerTurn)

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

	// Get the display characters from player components
	var p1Char, p2Char string
	playerEnts := g.world.ComponentManager.GetAllEntitiesWithComponent(components.Player)
	for _, ent := range playerEnts {
		player, _ := g.componentAccess.GetPlayerComponent(ent)
		if player.CellState == components.Player1 {
			p1Char = player.Character
		} else {
			p2Char = player.Character
		}
	}

	// Translate the board to a string representation
	displayBoard := make([][]string, 3)
	for y := range board.Board {
		displayBoard[y] = make([]string, 3)
		for x := range board.Board[y] {
			switch board.Board[y][x] {
			case components.Player1:
				displayBoard[y][x] = p1Char
			case components.Player2:
				displayBoard[y][x] = p2Char
			case components.Empty:
				displayBoard[y][x] = ""
			}
		}
	}

	g.displayManager.ShowBoard(displayBoard)
}

func (g *Game) getGameState() *components.GameStateComponent {
	gameStateEnts := g.world.ComponentManager.GetAllEntitiesWithComponent(components.GameState)
	if len(gameStateEnts) == 0 {
		return nil
	}

	gameState, hasComp := g.componentAccess.GetGameStateComponent(gameStateEnts[0])
	if !hasComp {
		return nil
	}

	return gameState
}
