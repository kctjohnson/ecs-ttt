package network

import (
	"encoding/json"
	"time"
)

// MessageType identifies the kind of message
type MessageType string

// Message represents a network message
type Message struct {
	Type      MessageType     `json:"type"`
	ClientID  string          `json:"client_id"`
	Timestamp int64           `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

// NewMessage creates a new message of the specified type
func NewMessage(msgType MessageType, clientID string, payload any) (*Message, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type:      msgType,
		ClientID:  clientID,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payloadBytes,
	}, nil
}
