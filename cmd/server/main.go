package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"example.com/ttt/internal/messaging"
	"example.com/ttt/internal/tictactoe"
)

const maxClients = 100

// TODO: scope variables properly, e.g., pass to handler as a dependency
var upgrader = websocket.Upgrader{} // TODO: configure
var clients = make(chan struct{}, maxClients)
var lobby chan tictactoe.Subscription

func handler(w http.ResponseWriter, r *http.Request) {
	// TODO: count active clients, measure duration
	// TODO: collect metrics in general
	select {
	case clients <- struct{}{}:
	default:
		// TODO: improve logging
		log.Println("capacity reached")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer func() {
		<-clients
	}()

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ugprade:", err)
		return
	}
	defer c.Close()

	// TODO: add messaging for game wait

	// TODO: cancel/interrupt
	// TODO: measure duration
	game := <-lobby

	// TODO: improve (debug) logging
	for {
		update, ok := <-game
		if !ok {
			return
		}

		msg := messaging.Message{
			Code: messaging.StatusUpdate,
			Payload: messaging.Payload{
				StatusUpdate: &update,
			},
		}
		err = c.WriteJSON(msg)
		if err != nil {
			log.Println("write:", err)
			break
		}

		if update.Winner != nil {
			break
		}

		msg = messaging.Message{
			Code: messaging.Prompt,
			Payload: messaging.Payload{
				Prompt: update.Piece,
			},
		}
		err = c.WriteJSON(msg)
		if err != nil {
			log.Println("write:", err)
			break
		}

		var respond messaging.Message
		// TODO: handle closed client
		err := c.ReadJSON(&respond)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				break
			}
			log.Println("read:", err)
			break
		}
		if respond.Code != messaging.Move {
			log.Println("msg code:", err)
			break
		}

		update.NextMove <- respond.Payload.Move
		close(update.NextMove)
	}

	// TODO: check for closed client
	err = c.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(5*time.Second), // TODO: configure
	)
	if err != nil {
		log.Println("close:", err)
	}
}

func main() {
	lobby = make(chan tictactoe.Subscription)
	go func() {
		for {
			p1 := make(chan tictactoe.GameUpdate, 1)
			p2 := make(chan tictactoe.GameUpdate, 1)
			lobby <- p1
			lobby <- p2
			tictactoe.Start(context.TODO(), p1, p2)
		}
		// TODO: exit
	}()

	http.HandleFunc("/", handler)

	// TODO: configure server
	log.Fatal(http.ListenAndServe(":8080", nil))

	// TODO: handle graceful shutdown
}
