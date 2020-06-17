package server

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Pool struct {
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	Clients    map[*websocket.Conn]bool
	Broadcast  chan string
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan string),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			break
		case message := <-pool.Broadcast:
			for client, _ := range pool.Clients {
				if err := client.WriteMessage(1, []byte(message)); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}