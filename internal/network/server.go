package network

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ClientConnection represents a connected client
type ClientConnection struct {
	ID   string
	Conn *websocket.Conn
}

// GameServer manages WebSocket connections and game instances
type GameServer struct {
	clients     map[string]*ClientConnection
	clientsLock sync.Mutex
	upgrader    websocket.Upgrader
	logger      *log.Logger
}

func NewGameServer(logger *log.Logger) *GameServer {
	return &GameServer{
		clients: make(map[string]*ClientConnection),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all connections for now
			},
		},
		logger: logger,
	}
}

func (gs *GameServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := gs.upgrader.Upgrade(w, r, nil)
	if err != nil {
		gs.logger.Printf("Error upgrading connection: %v", err)
		return
	}

	// Generate a unique client ID
	clientID := generateClientID()

	// Create a new client connection
	client := &ClientConnection{
		ID:   clientID,
		Conn: conn,
	}

	// Register the client
	gs.registerClient(client)

	// Handle client messages
	go gs.handleClientMessages(client)
}

func (gs *GameServer) Start(address string) error {
	http.HandleFunc("/ws", gs.HandleConnections)
	gs.logger.Printf("WebSocket server starting on %s", address)
	return http.ListenAndServe(address, nil)
}

// registerClient adds a client to the clients map
func (gs *GameServer) registerClient(client *ClientConnection) {
	gs.clientsLock.Lock()
	defer gs.clientsLock.Unlock()

	gs.clients[client.ID] = client
	gs.logger.Printf("Client %s connected", client.ID)
}

// unregisterClient removes a client from the clients map
func (gs *GameServer) unregisterClient(clientID string) {
	gs.clientsLock.Lock()
	defer gs.clientsLock.Unlock()

	if client, ok := gs.clients[clientID]; ok {
		client.Conn.Close()
		delete(gs.clients, clientID)
		gs.logger.Printf("Client disconnected: %s", clientID)
	}
}

// handleClientMessages processes messages from a client
func (gs *GameServer) handleClientMessages(client *ClientConnection) {
	defer gs.unregisterClient(client.ID)

	for {
		// Read message
		messageType, message, err := client.Conn.ReadMessage()
		if err != nil {
			gs.logger.Printf("Error reading message from client %s: %v", client.ID, err)
			break
		}

		// For now, just log the message
		gs.logger.Printf("Received message from client %s: %s", client.ID, message)

		// Echo the message back (for testing)
		if err := client.Conn.WriteMessage(messageType, message); err != nil {
			gs.logger.Printf("Error writing message to client %s: %v", client.ID, err)
			break
		}
	}
}

// generateClientID generates a unique ID for a client
func generateClientID() string {
	return "client-" + randomString(8)
}

// randomString generates a random string of a given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randInt(len(charset))]
	}
	return string(b)
}

// randInt generates a random int up to max
func randInt(max int) int {
	// Simple random number generation - in production, use crypto/rand
	return int(time.Now().UnixNano() % int64(max))
}
