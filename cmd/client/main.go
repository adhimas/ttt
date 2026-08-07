package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"example.com/ttt/internal/messaging"
	"example.com/ttt/internal/tictactoe"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// connect to server
	u := url.URL{Scheme: "ws", Host: ":8080", Path: "/"} // TODO: support config

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// listen to messages from server
	done := make(chan struct{})
	srvMessages := make(chan messaging.Message)
	go func() {
		defer close(done)
		for {
			var message messaging.Message
			err := c.ReadJSON(&message)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					return
				}
				log.Println("read:", err)
				return
			}
			srvMessages <- message
		}
	}()

	for {
		select {
		case m := <-srvMessages:
			if m.Code == messaging.StatusUpdate && m.Payload.StatusUpdate != nil {
				printBoard(m.Payload.StatusUpdate.Board)
				if m.Payload.StatusUpdate.Winner != nil {
					switch *m.Payload.StatusUpdate.Winner {
					case tictactoe.Empty:
						log.Printf("tied!")
					case m.Payload.StatusUpdate.Piece:
						log.Printf("you won!")
					default:
						log.Printf("game over")
					}
				}
				continue
			}
			if m.Code != messaging.Prompt {
				continue
			}
			move := prompt(m.Payload.Prompt)
			msg := messaging.Message{
				Code: messaging.Move,
				Payload: messaging.Payload{
					Move: move,
				},
			}
			err := c.WriteJSON(msg)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func printBoard(board []tictactoe.Cell) {
	for i, cell := range board {
		if cell == tictactoe.Empty {
			fmt.Printf(" %d", i)
		} else {
			fmt.Printf(" %s", cell.String())
		}

		if i%3 == 2 {
			fmt.Println()
		}
	}
}

func prompt(player tictactoe.Cell) int {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("enter number for %s: ", player)

		if more := scanner.Scan(); !more {
			continue
		}
		// TODO: handle interrupt
		if err := scanner.Err(); err != nil {
			log.Println("scanner:", err)
			continue
		}

		line := scanner.Text()
		cellNum, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			log.Println("prompt input:", err)
			continue
		}
		if cellNum < 0 || cellNum > 8 { // 0-8
			log.Println("prompt input: invalid move")
			continue
		}

		fmt.Println()
		return cellNum
	}
}
