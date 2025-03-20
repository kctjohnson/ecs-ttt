package network

import (
	"encoding/json"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

// GameClient handles WebSocket communication with the server
type GameClient struct {
	serverURL    string
	conn         *websocket.Conn
	clientID     string
	messageQueue messageQueue
	connected    bool
	connMutex    sync.Mutex
	logger       *log.Logger
}

// NewGameClient creates a new game client
func NewGameClient(serverURL string, logger *log.Logger) *GameClient {
	return &GameClient{
		serverURL:    serverURL,
		connected:    false,
		messageQueue: messageQueue{messages: []*Message{}},
		logger:       logger,
	}
}

// Connect establishes a WebSocket connection with the server
func (gc *GameClient) Connect() error {
	gc.connMutex.Lock()
	defer gc.connMutex.Unlock()

	if gc.connected {
		return nil
	}

	// Parse the WebSocket URL
	u, err := url.Parse(gc.serverURL)
	if err != nil {
		return err
	}

	// Create the WebSocket connection
	gc.logger.Printf("Connecting to %s", gc.serverURL)
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	gc.conn = conn
	gc.connected = true

	// Start a goroutine to handle incoming messages
	go gc.readMessages()

	// Send a connect message
	err = gc.SendMessage(MsgConnect, &ConnectPayload{
		Username: "Player", // Could be configurable
	})

	return err
}

// Disconnect closes the WebSocket connection
func (gc *GameClient) Disconnect() {
	gc.connMutex.Lock()
	defer gc.connMutex.Unlock()

	if !gc.connected || gc.conn == nil {
		return
	}

	gc.conn.Close()
	gc.connected = false
	gc.clientID = ""
}

// SendMessage sends a message to the server
func (gc *GameClient) SendMessage(msgType MessageType, payload any) error {
	gc.connMutex.Lock()
	defer gc.connMutex.Unlock()

	if !gc.connected || gc.conn == nil {
		return ErrNotConnected
	}

	msg, err := NewMessage(msgType, gc.clientID, payload)
	if err != nil {
		return err
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return gc.conn.WriteMessage(websocket.TextMessage, jsonMsg)
}

// GetPendingMessages returns all pending messages and clears the queue
func (gc *GameClient) GetPendingMessages() []*Message {
	gc.messageQueue.lock.Lock()
	defer gc.messageQueue.lock.Unlock()

	messages := gc.messageQueue.messages
	gc.messageQueue.messages = []*Message{}
	return messages
}

// IsConnected returns whether the client is connected to the server
func (gc *GameClient) IsConnected() bool {
	gc.connMutex.Lock()
	defer gc.connMutex.Unlock()
	return gc.connected
}

// readMessages processes incoming messages from the server
func (gc *GameClient) readMessages() {
	defer gc.Disconnect()

	for {
		_, rawMessage, err := gc.conn.ReadMessage()
		if err != nil {
			gc.logger.Printf("Error reading message: %v", err)
			break
		}

		// Parse the message
		var msg Message
		if err := json.Unmarshal(rawMessage, &msg); err != nil {
			gc.logger.Printf("Invalid message format: %v", err)
			continue
		}

		// Store client ID if this is the first message
		if gc.clientID == "" && msg.ClientID != "" {
			gc.clientID = msg.ClientID
		}

		// Queue the message for processing
		gc.messageQueue.lock.Lock()
		gc.messageQueue.messages = append(gc.messageQueue.messages, &msg)
		gc.messageQueue.lock.Unlock()

		gc.logger.Printf("Received message: %s", msg.Type)
	}
}

// Define error types
var (
	ErrNotConnected = &NetworkError{
		Code:    "note_connected",
		Message: "Client is not connected to the server",
	}
)

// NetworkError represents a network-related error
type NetworkError struct {
	Code    string
	Message string
}

func (e *NetworkError) Error() string {
	return e.Message
}
