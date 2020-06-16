package main

import (
	"github.com/Mindgamesnl/piper/client"
	"github.com/Mindgamesnl/piper/server"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		logrus.Error("invalid arguments! use piper <client/server> <arguments>")
		os.Exit(1)
	}

	if args[0] == "client" {
		client.StartClient()
		return
	}

	if args[0] == "server" {
		server.StartServer()
		return
	}

	logrus.Error("invalid arguments! use piper <client/server> <arguments>")
	os.Exit(1)
}
