package server

import (
	"fmt"
	"github.com/Mindgamesnl/piper/client"
	"github.com/gorilla/websocket"
	"sync"
)

type Pool struct {
	Register   chan Player
	Unregister chan Player
	Clients    map[Player]bool
	Broadcast  chan string
}

type Player struct {
	Socket *websocket.Conn // websocket connection of the player
	mu      sync.Mutex
}

func (p *Player) send(v []byte) error {
	p.mu.Lock()
	result := p.Socket.WriteMessage(1, v)
	defer p.mu.Unlock()
	return result
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan Player),
		Unregister: make(chan Player),
		Clients:    make(map[Player]bool),
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
 				if err := client.send([]byte(message)); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

func BroadcastMessage(message string)  {
	ConnectionPool.Broadcast <- client.InfoColor + "PIPER: " + message + "\033[0m"
}

func BroadcastCommandError(command string, message string)  {
	ConnectionPool.Broadcast <- client.WarningColor + "Error while executing '" + command + "', " + message + "\033[0m"
}

func BroadcastServiceError(command string)  {
	ConnectionPool.Broadcast <- client.WarningColor + command + "\033[0m"
}

func BroadcastCommandOutput(message string) {
	ConnectionPool.Broadcast <- client.NoticeColor + message + "\033[0m"
}

func BroadcastServiceOutput(message string) {
	ConnectionPool.Broadcast <- client.DebugColor + message + "\033[0m"
}
