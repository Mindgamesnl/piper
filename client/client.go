package client

import "github.com/sirupsen/logrus"

func StartClient()  {
	logrus.Info("Starting client")
	LoadConfiguration()
	StartFileWatcher()
}