package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Pool struct {
	Reg       chan *Client
	Unreg     chan *Client
	Clients   map[*Client]bool
	Broadcast chan Message
}

func newPool() *Pool {
	return &Pool{
		Reg:       make(chan *Client),
		Unreg:     make(chan *Client),
		Clients:   make(map[*Client]bool),
		Broadcast: make(chan Message),
	}
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type int
	Body string
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unreg <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Error:", err)
			return
		}
		message := Message{Type: messageType, Body: string(p)}
		c.Pool.Broadcast <- message
		fmt.Println("Received:", message)
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Reg:
			pool.Clients[client] = true
			fmt.Println("Size of connection pool:", len(pool.Clients))
			for client := range pool.Clients {
				fmt.Println("")
				client.Conn.WriteJSON(Message{Type: 1, Body: "New user joined"})
			}
		case client := <-pool.Unreg:
			delete(pool.Clients, client)
			fmt.Println("Size of connection pool:", len(pool.Clients))
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "New user disconnected"})
			}
		case message := <-pool.Broadcast:
			fmt.Println("Received message", message)
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
