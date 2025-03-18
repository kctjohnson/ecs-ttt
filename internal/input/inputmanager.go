package input

type InputManager interface {
	GetPlayerMove() (row, col int, valid bool)
}
