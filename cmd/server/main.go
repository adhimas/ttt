package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(websocket.TextMessage, []byte("server reply"))
		if err != nil {
			log.Println("write:", err)
			break
		}
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
