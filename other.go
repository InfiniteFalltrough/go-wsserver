package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type user_id = string

var connections map[user_id]*websocket.Conn = make(map[user_id]*websocket.Conn)

func run() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	err = ws.WriteMessage(1, []byte("OK"))
	if err != nil {
		log.Println(err)
	}

	// a random way to generate unique id for session
	h := sha256.New()
	h.Write([]byte(r.RemoteAddr))
	var id user_id = hex.EncodeToString(h.Sum(nil))
	connections[id] = ws
	// ws.WriteMessage(websocket.BinaryMessage, []byte(fmt.Sprintf("{'id': '%s'}", id))) //Dummy way, just as example
	payl := map[string]interface{}{"id": id}
	bts, _ := json.Marshal(payl)
	ws.WriteMessage(websocket.BinaryMessage, []byte(bts))
	reader(ws)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Payload struct {
	To   string `json:"to"`
	What string `json:"what"`
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))

		var payload *Payload = &Payload{}
		errUnmarshall := json.Unmarshal(p, payload)
		if errUnmarshall != nil {
			panic(errUnmarshall)
		}

		if sock, exist := connections[payload.To]; exist {
			if err := sock.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				return
			}
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}
