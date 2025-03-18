package console

import (
	"fmt"
)

type ConsoleInputManager struct{}

func NewConsoleInputManager() ConsoleInputManager {
	return ConsoleInputManager{}
}

func (c ConsoleInputManager) GetPlayerMove() (row, col int, valid bool) {
	_, err := fmt.Scanf("%d %d", &col, &row)
	if err != nil || row < 0 || row > 2 || col < 0 || col > 2 {
		return 0, 0, false
	}
	return row, col, true
}
