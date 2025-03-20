package systems

import (
	"encoding/json"

	"ttt/internal/game/components"
	"ttt/internal/game/events"
	"ttt/internal/network"
	"ttt/pkg/ecs"
)

// ClientNetworkSystem handles communication with the server on the client side
type ClientNetworkSystem struct {
	ComponentAccess *components.ComponentAccess
	Client          *network.GameClient
	gameID          string
	localPlayerID   ecs.Entity
	isPlayerTurn    bool
	mark            string // "X" or "O"
}

func (s *ClientNetworkSystem) Update(world *ecs.World) {
	// Process any pending network messages
	messages := s.Client.GetPendingMessages()
	for _, msg := range messages {
		s.handleMessage(msg, world)
	}

	// Check for move intent components and send to server
	s.processMoveIntents(world)
}

func (s *ClientNetworkSystem) handleMessage(msg *network.Message, world *ecs.World) {
	switch msg.Type {
	case network.MsgGameJoined:
		s.handleGameJoined(msg, world)
	case network.MsgBoardUpdate:
		s.handleBoardUpdate(msg, world)
	case network.MsgGameOver:
		s.handleGameOver(msg, world)
	}
}

func (s *ClientNetworkSystem) handleGameJoined(msg *network.Message, world *ecs.World) {
	var payload network.GameJoinedPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return
	}

	s.gameID = payload.GameID
	s.mark = payload.YourMark
	s.isPlayerTurn = payload.IsYourTurn

	// Find the local player entity
	playerEnts := world.ComponentManager.GetAllEntitiesWithComponent(components.Player)
	for _, ent := range playerEnts {
		player, _ := s.ComponentAccess.GetPlayerComponent(ent)
		if player.Character == s.mark {
			s.localPlayerID = ent
			break
		}
	}

	// Update game state to reflect whose turn it is
	gameStateEnts := world.ComponentManager.GetAllEntitiesWithComponent(components.GameState)
	if len(gameStateEnts) > 0 {
		gameState, _ := s.ComponentAccess.GetGameStateComponent(gameStateEnts[0])
		if s.isPlayerTurn {
			gameState.PlayerTurn = s.localPlayerID
		} else {
			// Set to the other player
			for _, ent := range playerEnts {
				if ent != s.localPlayerID {
					gameState.PlayerTurn = ent
				}
			}
		}
	}
}

func (s *ClientNetworkSystem) handleBoardUpdate(msg *network.Message, world *ecs.World) {
	var payload network.BoardUpdatePayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return
	}

	// Update the board state
	boardEnts := world.ComponentManager.GetAllEntitiesWithComponent(components.Board)
	if len(boardEnts) == 0 {
		return
	}

	board, _ := s.ComponentAccess.GetBoardComponent(boardEnts[0])

	// Convert string representation to CellState
	for i := range 3 {
		for j := range 3 {
			switch payload.Board[i][j] {
			case "X":
				board.Board[i][j] = components.Player1
			case "O":
				board.Board[i][j] = components.Player2
			case "":
				board.Board[i][j] = components.Empty
			}
		}
	}

	// Update whose turn it is
	s.isPlayerTurn = (payload.Turn == s.mark)

	// Update game state
	gameStateEnts := world.ComponentManager.GetAllEntitiesWithComponent(components.GameState)
	if len(gameStateEnts) > 0 {
		gameState, _ := s.ComponentAccess.GetGameStateComponent(gameStateEnts[0])

		// Find all player entities
		playerEnts := world.ComponentManager.GetAllEntitiesWithComponent(components.Player)

		// Set turn to correct player
		for _, ent := range playerEnts {
			player, _ := s.ComponentAccess.GetPlayerComponent(ent)
			if (payload.Turn == "X" && player.Character == "X") ||
				(payload.Turn == "O" && player.Character == "O") {
				gameState.PlayerTurn = ent
				break
			}
		}
	}
}

func (s *ClientNetworkSystem) handleGameOver(msg *network.Message, world *ecs.World) {
	var payload network.GameOverPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return
	}

	// Update game state to game over
	gameStateEnts := world.ComponentManager.GetAllEntitiesWithComponent(components.GameState)
	if len(gameStateEnts) > 0 {
		gameState, _ := s.ComponentAccess.GetGameStateComponent(gameStateEnts[0])
		gameState.GameOver = true
	}

	// Generate appropriate local event
	switch payload.Result {
	case "X_won":
		if s.mark == "X" {
			world.QueueEvent(events.PlayerWonEvent{Ent: s.localPlayerID})
		} else {
			// Find the other player
			playerEnts := world.ComponentManager.GetAllEntitiesWithComponent(components.Player)
			for _, ent := range playerEnts {
				if ent != s.localPlayerID {
					world.QueueEvent(events.PlayerWonEvent{Ent: ent})
					break
				}
			}
		}
	case "O_won":
		if s.mark == "O" {
			world.QueueEvent(events.PlayerWonEvent{Ent: s.localPlayerID})
		} else {
			// Find the other player
			playerEnts := world.ComponentManager.GetAllEntitiesWithComponent(components.Player)
			for _, ent := range playerEnts {
				if ent != s.localPlayerID {
					world.QueueEvent(events.PlayerWonEvent{Ent: ent})
					break
				}
			}
		}
	case "tie":
		world.QueueEvent(events.TieEvent{Ent: -1})
	}
}
