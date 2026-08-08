package tictactoe

import (
	"context"
	"slices"
	"time"
)

func Start(ctx context.Context, p1, p2 Subscription) {
	go start(ctx, p1, p2)
}

// start creates a new game and runs it.
// the game state is kept private and incoming/outgoing communications are handled via channels.
//   - updates are sent to each player via their subscription channels.
//   - player moves are sent via a channel included in the updates
//
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

		// TODO: handle blocked update, add timeout
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
