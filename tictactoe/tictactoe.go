package tictactoe

import (
	"context"
	"slices"
	"time"
)

type Cell int

const (
	Empty Cell = iota
	X
	O
)

var cellNames = map[Cell]string{
	Empty: "_",
	X:     "X",
	O:     "O",
}

func (c Cell) String() string {
	return cellNames[c]
}

type GameUpdate struct {
	Board    []Cell
	NextMove chan<- int `json:"-"`
	Piece    Cell       // nextplayer
	Winner   *Cell
}

// TODO: support cancellation
// TODO: improve comments in general
type Subscription chan GameUpdate

// TODO: handle context
// TODO: improve (debug) logging
func start(_ context.Context, p1, p2 Subscription) {
	game := newGame()
	for game.canContinue() {
		board := slices.Clone(game.board)
		nextMove := make(chan int)
		update := GameUpdate{
			NextMove: nextMove,
			Board:    board,
		}

		var ch Subscription
		var piece Cell
		if game.turn == playerX {
			ch = p1
			piece = X
		} else {
			ch = p2
			piece = O
		}

		update.Piece = piece

		// TODO: check potential blocked update
		ch <- update

		var abort bool
		select {
		case move := <-nextMove:
			if err := game.handleMove(move); err != nil {
				continue
			}
		case <-time.After(30 * time.Second): // TODO: support config
			abort = true
		}

		if abort {
			break
		}

		game.endTurn()
	}

	game.cleanUp(p1, p2)
}

func Start(ctx context.Context, p1, p2 Subscription) {
	go start(ctx, p1, p2)
}
