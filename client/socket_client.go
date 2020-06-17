package client

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var Connection *websocket.Conn

func ConnectSocket(callback func())  {
	port := LoadedInstance.Port
	host := LoadedInstance.ServerHost
	hostAndPort := host + ":" + strconv.Itoa(port)
	password := LoadedInstance.Password

	Log(InfoColor + "Connecting to server at " + hostAndPort)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: hostAndPort, Path: "/piper", RawQuery: "password=" + password}

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	Connection = c
	if err == websocket.ErrBadHandshake {
		log.Printf("handshake failed with status %d", resp.StatusCode)
	}
	if err != nil {
		Log(ErrorColor + "Connection refused. Is the hostname and password correct, or is the server down?")
		return
	}
	defer Connection.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := Connection.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			PrintRemote(string(message))
		}
	}()

	go callback()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := Connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
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

func WriteSocket(whatToWrite string)  {
	Connection.WriteMessage(websocket.TextMessage, []byte(whatToWrite))
}