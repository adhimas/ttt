package tictactoe

import (
	"errors"
	"slices"
)

var winningTriples = [][]int{
	{0, 1, 2}, // rows
	{3, 4, 5},
	{6, 7, 8},

	{0, 3, 6}, // cols
	{1, 4, 7},
	{2, 5, 8},

	{0, 4, 8}, // diags
	{2, 4, 6},
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
	for _, triple := range winningTriples {
		a, b, c := triple[0], triple[1], triple[2]
		if g.board[a] == X && g.board[b] == X && g.board[c] == X {
			return true, false
		}
		if g.board[a] == O && g.board[b] == O && g.board[c] == O {
			return false, true
		}
	}

	return false, false
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
