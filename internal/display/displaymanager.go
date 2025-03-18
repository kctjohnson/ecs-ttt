package display

type DisplayManager interface {
	ShowBoard(board [][]string)
	ShowTurnPrompt(player string)
	ShowGameResult(result string)
}
