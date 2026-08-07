package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"

	"example.com/ttt/internal/client/cli"
	"example.com/ttt/internal/messaging"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

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
				cli.PrintGame(m.Payload.StatusUpdate)
				continue
			}
			if m.Code != messaging.Prompt {
				continue
			}
			move := cli.Prompt(m.Payload.Prompt)
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
			err = c.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				time.Now().Add(5*time.Second), // TODO: configure
			)
			if err != nil {
				log.Println("write close:", err)
			}
			return
		}
	}
}
