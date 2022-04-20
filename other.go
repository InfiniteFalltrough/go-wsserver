package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func run() {
	pool := newPool()
	go pool.Start()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handler(pool, w, r)
	})
}

func handler(pool *Pool, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	client := &Client{
		Conn: ws,
		Pool: pool,
	}
	log.Printf("Client connected\n")
	err = ws.WriteMessage(1, []byte("200 - OK"))
	if err != nil {
		log.Println(err)
	}
	pool.Reg <- client
	client.Read()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
