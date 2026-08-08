package tictactoe

type GameUpdate struct {
	Board    []Cell
	NextMove chan<- int `json:"-"`
	Piece    Cell       // next player
	Winner   *Cell
}

// TODO: support cancellation
// TODO: improve comments in general
type Subscription chan GameUpdate
