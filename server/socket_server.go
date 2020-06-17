package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartSocket()  {
	http.HandleFunc("/piper", func(w http.ResponseWriter, r *http.Request) {
		if !ValidatePassword(r) {
			logrus.Error("Warning! blocked incoming connection without valid authentication")
			return
		}

		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, "websockets.html")
	// })

	http.ListenAndServe(":" + strconv.Itoa(LoadedInstance.Port), nil)
}

func ValidatePassword(r *http.Request) bool {
	return r.URL.Query().Get("password") == LoadedInstance.Password
}
