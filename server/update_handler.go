package server

import (
	"github.com/Mindgamesnl/piper/common"
	"github.com/radovskyb/watcher"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func HandleFileUpdate(update common.FileUpdate) error {
	logrus.Info("Request to write file " + update.Name)

	if update.Operation == watcher.Create || update.Operation == watcher.Write {
		localPath := update.RelativePath
		localPath = strings.Replace(localPath, update.Name, "", -1)
		os.MkdirAll("." + localPath, os.ModePerm)

		file, err := os.OpenFile("." + update.RelativePath, os.O_RDWR | os.O_CREATE, 0666)

		if err != nil {
			return err
		}

		defer file.Close()

		_, err = file.Write(update.Content)
		if err != nil {
			return err
		}

		ConnectionPool.Broadcast <- "Wrote file " + update.Name
		return nil
	}

	if update.Operation == watcher.Remove {
		ConnectionPool.Broadcast <- "Deleted file " + update.Name
		return os.Remove("." + update.RelativePath)
		return nil
	}

	ConnectionPool.Broadcast <- "Unknown action for " + update.Name
	return nil
}
