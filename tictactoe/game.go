package tictactoe

import (
	"errors"
	"slices"
)

const magicConstant = 15

// every column/row/main diagonal adds up equally
var magicSquare = []int{
	2, 7, 6,
	9, 5, 1,
	4, 3, 8,
}

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

type turn int

const (
	playerX turn = iota
	playerO
)

type game struct {
	board []Cell
	turn  turn
}

func (g *game) handleMove(i int) error {
	if g.board[i] != Empty {
		return errors.New("invalid move")
	}

	if g.turn == playerX {
		g.board[i] = X
	} else {
		g.board[i] = O
	}

	return nil
}

func (g *game) endTurn() {
	if g.turn == playerX {
		g.turn = playerO
	} else {
		g.turn = playerX
	}
}

func (g *game) checkWinner() (bool, bool) {
	xScore, oScore := 0, 0
	for i, cell := range g.board {
		switch cell {
		case X:
			xScore += magicSquare[i]
		case O:
			oScore += magicSquare[i]
		}
	}
	return xScore == magicConstant, oScore == magicConstant
}

func (g *game) canContinue() bool {
	// TODO: check for early draw

	if xWins, oWins := g.checkWinner(); xWins || oWins {
		return false
	}

	return slices.Contains(g.board, Empty)
}

func newGame() *game {
	board := make([]Cell, 9)
	for i := range board {
		board[i] = Empty
	}

	return &game{
		board: board,
		turn:  playerX,
	}
}

// cleanUp sends the last update then closes the subscriptions.
func (g *game) cleanUp(p1, p2 Subscription) {
	board := slices.Clone(g.board)
	update := GameUpdate{Board: board}

	xWins, oWins := g.checkWinner()
	winner := Empty
	if xWins {
		winner = X
	} else if oWins {
		winner = O
	}
	update.Winner = &winner

	updateP1 := update
	updateP2 := update

	updateP1.Piece = X
	updateP2.Piece = O

	p1 <- updateP1
	p2 <- updateP2

	close(p1)
	close(p2)
}
