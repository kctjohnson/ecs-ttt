package console

import "fmt"

type ConsoleDisplayManager struct{}

func NewConsoleDisplayManager() ConsoleDisplayManager {
	return ConsoleDisplayManager{}
}

func (c ConsoleDisplayManager) ShowBoard(board [][]string) {
	fmt.Println()
	for y := range 3 {
		for x := range 3 {
			if board[y][x] == "" {
				fmt.Print(".")
			} else {
				fmt.Print(board[y][x])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (c ConsoleDisplayManager) ShowTurnPrompt(player string) {
	fmt.Printf("%s, enter column and row (0-2) separated by a space:\n", player)
}

func (c ConsoleDisplayManager) ShowGameResult(result string) {
	fmt.Println(result)
}
