package systems

import (
	"encoding/json"

	"ttt/internal/game/components"
	"ttt/internal/network"
	"ttt/pkg/ecs"
)

// NetworkSystem processes network events and synchronizes game state
type NetworkSystem struct {
	ComponentAccess *components.ComponentAccess
	Server          *network.GameServer
}

func (s *NetworkSystem) Update(world *ecs.World) {
	// Process messages received from the network
	messages := s.Server.GetPendingMessages()
	for _, msg := range messages {
		s.handleMessage(msg, world)
	}

	// Sync game state to connected clients
	s.syncGameState(world)
}

func (s *NetworkSystem) handleMessage(msg *network.Message, world *ecs.World) {
	switch msg.Type {
	case network.MsgJoinGame:
		s.handleJoinGame(msg, world)
	case network.MsgPlayerMove:
		s.handlePlayerMove(msg, world)
	}
}

func (s *NetworkSystem) handleJoinGame(msg *network.Message, world *ecs.World) {
	var payload network.JoinGamePayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return
	}

	// Find or create a game session
	session := s.findAvailableSession(world)
	if session == nil {
		// Create a new game session
		session = s.createGameSession(world)
	}

	// Add the player to the session
	s.addPlayerToSession(msg.ClientID, session, world)
}

func (s *NetworkSystem) handlePlayerMove(msg *network.Message, world *ecs.World) {
	var payload network.PlayerMovePayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return
	}

	// Find the game session
	sessionEntity := s.findSessionByID(payload.GameID, world)
	if sessionEntity == ecs.Entity(-1) {
		return
	}

	// Find the player entity for this client
	playerEntity := s.findPlayerEntityByClientID(msg.ClientID, world)
	if playerEntity == ecs.Entity(-1) {
		return
	}

	// Create a move intent
	moveIntent := &components.MoveIntentComponent{
		Row: payload.Row,
		Col: payload.Col,
	}

	// Add the component to the player entity
	world.ComponentManager.AddComponent(
		playerEntity,
		components.MoveIntent,
		moveIntent,
	)
}

func (s *NetworkSystem) syncGameState(world *ecs.World) {
	// Find all active game sessions
	sessionEntities := world.ComponentManager.GetAllEntitiesWithComponent(components.GameSession)

	for _, entity := range sessionEntities {
		session, found := s.ComponentAccess.GetGameSessionComponent(entity)
		if !found {
			continue
		}

		if !session.IsActive {
			continue
		}

		// Get the board component
		board, found := s.ComponentAccess.GetBoardComponent(session.BoardEntity)
		if !found {
			continue
		}

		// Create a board update payload
		boardUpdate := s.createBoardUpdatePayload(session, board, world)

		// Send to both players
		s.sendToPlayer(session.Player1, network.MsgBoardUpdate, boardUpdate, world)
		s.sendToPlayer(session.Player2, network.MsgBoardUpdate, boardUpdate, world)
	}
}

// Helper methods
func (s *NetworkSystem) findAvailableSession(world *ecs.World) *components.GameSessionComponent {
	// Find a session with only one player
	sessionEntities := world.ComponentManager.GetAllEntitiesWithComponent(components.GameSession)

	for _, entity := range sessionEntities {
		session, found := s.ComponentAccess.GetGameSessionComponent(entity)
		if !found {
			continue
		}

		if !session.IsActive {
			continue
		}

		if session.Player2 == ecs.Entity(-1) {
			return session
		}
	}

	return nil
}

// TODO: Implement all of these functions

func (s *NetworkSystem) createGameSession(world *ecs.World) *components.GameSessionComponent {
	// Implementation details omitted for brevity
	// Create a new game session with board and game state
	return nil
}

func (s *NetworkSystem) addPlayerToSession(
	clientID string,
	session *components.GameSessionComponent,
	world *ecs.World,
) {
	// Implementation details omitted for brevity
	// Add a player to the session and notify them
}

func (s *NetworkSystem) findSessionByID(gameID string, world *ecs.World) ecs.Entity {
	// Implementation details omitted for brevity
	// Find a session by its ID
	return -1
}

func (s *NetworkSystem) findPlayerEntityByClientID(clientID string, world *ecs.World) ecs.Entity {
	// Implementation details omitted for brevity
	// Find a player entity by client ID
	return -1
}

func (s *NetworkSystem) createBoardUpdatePayload(
	session *components.GameSessionComponent,
	board *components.BoardComponent,
	world *ecs.World,
) *network.BoardUpdatePayload {
	// Implementation details omitted for brevity
	// Create a board update payload from the current state
	return nil
}

func (s *NetworkSystem) sendToPlayer(
	playerEntity ecs.Entity,
	msgType network.MessageType,
	payload interface{},
	world *ecs.World,
) {
	// Implementation details omitted for brevity
	// Send a message to a specific player
}
