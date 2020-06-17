package server

import (
	"github.com/Mindgamesnl/piper/common"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var ConnectionPool = NewPool()

func StartSocket()  {
	http.HandleFunc("/piper", func(w http.ResponseWriter, r *http.Request) {
		if !ValidatePassword(r) {
			logrus.Error("Warning! blocked incoming connection without valid authentication")
			return
		}

		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		ConnectionPool.Register <- conn
		conn.WriteMessage(1, []byte("Welcome! Waiting for output"))

		defer func() {
			ConnectionPool.Unregister <- conn
			conn.Close()
		}()

		for {
			// Read message from browser
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			err = HandleFileUpdate(common.FromJson(msg))
			if err != nil {
				logrus.Error(err)
				conn.WriteMessage(1, []byte("Error while writing file: " + err.Error()))
			}
		}
	})

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, "websockets.html")
	// })

	go ConnectionPool.Start()

	http.ListenAndServe(":" + strconv.Itoa(LoadedInstance.Port), nil)
}

func ValidatePassword(r *http.Request) bool {
	return r.URL.Query().Get("password") == LoadedInstance.Password
}
