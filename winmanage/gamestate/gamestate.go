package gamestate

// GameState specifies whether the process has last been paused or is currently unpaused
type GameState int

const (
	// Paused indicates that the game has been paused last, and not resumed since
	Paused GameState = 0

	// Playing indicates that the game hasn't been paused, or has been resumed since
	Playing GameState = 1
)

// Check checks the state of the GameState for validity
func (state *GameState) Check() {
	if state == nil {
		panic("Given GameState has nil pointer")
	}
	if (*state != Paused) && (*state != Playing) {
		panic("GameState is invalid")
	}
}
