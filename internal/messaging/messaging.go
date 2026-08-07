package messaging

type Code int

const (
	Other Code = iota
	Prompt
	Move
	StatusUpdate
)

type Message struct {
	Payload any
	Text    string
	Code    Code
}
