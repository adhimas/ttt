package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"example.com/ttt/internal/messaging"
)

var upgrader = websocket.Upgrader{} // use default options

var clients = make(chan struct{}, 100)

var lobby chan struct{}

func handler(w http.ResponseWriter, r *http.Request) {
	select {
	case clients <- struct{}{}:
	default:
		log.Println("capacity reached")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer func() {
		<-clients
	}()

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	// ctx := r.Context()
	// TODO
	_ = <-lobby // select, timeout/interrupt

	for {
		err = c.WriteJSON(messaging.Message{Code: messaging.Other, Text: "test"})
		if err != nil {
			log.Println("write:", err)
			break
		}
		time.Sleep(time.Second)
	}
}

func main() {
	lobby = make(chan struct{})
	go func() {
		for {
			lobby <- struct{}{}
			lobby <- struct{}{}
			// TODO
		}
		// TODO: exit
	}()

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
