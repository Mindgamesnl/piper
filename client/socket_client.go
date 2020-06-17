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
	if err == websocket.ErrBadHandshake {
		log.Printf("handshake failed with status %d", resp.StatusCode)
	}
	if err != nil {
		Log(ErrorColor + "Connection refused. Is the hostname and password correct, or is the server down?")
		return
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

	callback()
}