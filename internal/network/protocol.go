package network

// MessageType identifiers
const (
	MsgConnect     MessageType = "connect"      // Client connects to server
	MsgDisconnect  MessageType = "disconnect"   // Client disconnects
	MsgJoinGame    MessageType = "join_game"    // Client requests to join a game
	MsgGameJoined  MessageType = "game_joined"  // Client successfully joined a game
	MsgPlayerMove  MessageType = "player_move"  // Player makes a move
	MsgBoardUpdate MessageType = "board_update" // Board state changed
	MsgGameOver    MessageType = "game_over"    // Game has ended
	MsgError       MessageType = "error"        // Error occurred
)

// ConnectPayload contains data for a connect message
type ConnectPayload struct {
	Username string `json:"username"`
}

// JoinGamePayload contains data for joining a game
type JoinGamePayload struct {
	GameID string `json:"game_id,omitempty"` // Optional - if empty, join any game
}

// GameJoinedPayload contains data about the game joined
type GameJoinedPayload struct {
	GameID     string `json:"game_id"`
	YourMark   string `json:"your_mark"` // X or O
	IsYourTurn bool   `json:"is_your_turn"`
}

// PlayerMovePayload contains data for a move
type PlayerMovePayload struct {
	GameID string `json:"game_id"`
	Row    int    `json:"row"`
	Col    int    `json:"col"`
}

// BoardUpdatePayload contains the current board state
type BoardUpdatePayload struct {
	GameID string     `json:"game_id"`
	Board  [][]string `json:"board"`
	Turn   string     `json:"turn"` // X or O
}

// GameOverPayload contains information about the game result
type GameOverPayload struct {
	GameID string `json:"game_id"`
	Result string `json:"result"` // "X_won", "O_won", or "tie"
}

// ErrorPayload contains error information
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
