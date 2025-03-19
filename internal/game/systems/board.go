package systems

import (
	"ttt/internal/game/components"
	"ttt/internal/game/events"
	"ttt/pkg/ecs"
)

// BoardSystem checks the state of the board, and sends out winner / tie events
type BoardSystem struct {
	ComponentAccess *components.ComponentAccess
}

func (b *BoardSystem) Update(world *ecs.World) {
	// Get the board
	boardEnts := world.ComponentManager.GetAllEntitiesWithComponent(components.Board)
	if len(boardEnts) == 0 {
		return
	}

	board, hasBoardComp := b.ComponentAccess.GetBoardComponent(boardEnts[0])
	if !hasBoardComp {
		return
	}

	// Get the players
	playerEnts := world.ComponentManager.GetAllEntitiesWithComponent(components.Player)
	if len(playerEnts) != 2 {
		return
	}

	for _, playerEnt := range playerEnts {
		player, hasPlayerComp := b.ComponentAccess.GetPlayerComponent(playerEnt)
		if !hasPlayerComp {
			continue
		}

		// Check for winner
		if b.checkIfWin(board, player.Character) {
			world.QueueEvent(events.PlayerWon, playerEnt, nil)
			return
		}
	}

	// Check for a draw (No more spaces to move and no winner)
	if b.checkIfDraw(board) {
		world.QueueEvent(events.Tie, -1, nil)
		return
	}
}

func (b BoardSystem) checkIfWin(board *components.BoardComponent, playerChar string) bool {
	// Check rows
	if board.Board[0][0] == playerChar && board.Board[0][1] == playerChar &&
		board.Board[0][2] == playerChar {
		return true
	}
	if board.Board[1][0] == playerChar && board.Board[1][1] == playerChar &&
		board.Board[1][2] == playerChar {
		return true
	}
	if board.Board[2][0] == playerChar && board.Board[2][1] == playerChar &&
		board.Board[2][2] == playerChar {
		return true
	}

	// Check columns
	if board.Board[0][0] == playerChar && board.Board[1][0] == playerChar &&
		board.Board[2][0] == playerChar {
		return true
	}
	if board.Board[0][1] == playerChar && board.Board[1][1] == playerChar &&
		board.Board[2][1] == playerChar {
		return true
	}
	if board.Board[0][2] == playerChar && board.Board[1][2] == playerChar &&
		board.Board[2][2] == playerChar {
		return true
	}

	// Check diagonals
	if board.Board[0][0] == playerChar && board.Board[1][1] == playerChar &&
		board.Board[2][2] == playerChar {
		return true
	}
	if board.Board[0][2] == playerChar && board.Board[1][1] == playerChar &&
		board.Board[2][0] == playerChar {
		return true
	}

	return false
}

func (b BoardSystem) checkIfDraw(board *components.BoardComponent) bool {
	for i := range 3 {
		for j := range 3 {
			if board.Board[i][j] == "" {
				return false
			}
		}
	}
	return true
}
