package tictactoe

import (
	"testing"
	"time"
)

func TestCheckWinner(t *testing.T) {
	for _, triple := range winningTriples {
		game := newGame()
		for _, cell := range triple {
			game.board[cell] = X
		}
		x, _ := game.checkWinner()
		if !x {
			t.Errorf("x is %t; want true", x)
		}
	}
}

func TestStart(t *testing.T) {
	p1 := make(chan GameUpdate)
	p2 := make(chan GameUpdate)
	Start(t.Context(), p1, p2)

	p1Moves := []int{0, 1, 2}
	p2Moves := []int{8, 7, 6}

	for p1 != nil || p2 != nil {
		select {
		case update, ok := <-p1:
			if !ok { // game has closed the channel
				p1 = nil
				continue
			}
			if update.Winner != nil { // final update
				if *update.Winner != X {
					t.Errorf("winner is %s; want X", update.Winner)
				}
				continue
			}
			move := p1Moves[0]
			p1Moves = p1Moves[1:]
			update.NextMove <- move

		case update, ok := <-p2:
			if !ok {
				p2 = nil
				continue
			}
			if update.Winner != nil {
				if *update.Winner != X {
					t.Errorf("winner is %s; want X", update.Winner)
				}
				continue
			}
			move := p2Moves[0]
			p2Moves = p2Moves[1:]
			update.NextMove <- move

		case <-time.After(time.Second):
			t.Errorf("game didn't complete")
			return
		}
	}
}
