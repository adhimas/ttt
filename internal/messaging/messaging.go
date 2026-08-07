package messaging

import "example.com/ttt/internal/tictactoe"

type Code int

const (
	Other Code = iota
	Prompt
	Move
	StatusUpdate
)

type Payload struct { // TODO: generalize
	StatusUpdate *tictactoe.GameUpdate
	Move         int
	Prompt       tictactoe.Cell
}

type Message struct {
	Payload Payload
	Text    string
	Code    Code
}
