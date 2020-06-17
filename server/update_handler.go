package server

import (
	"github.com/Mindgamesnl/piper/common"
	"github.com/radovskyb/watcher"
	"os"
	"strings"
)

func HandleFileUpdate(update common.FileUpdate) error {
	if update.PiperOpcode == common.ExecuteCommands {
		for i := range update.ExecutableCommands {
			command := update.ExecutableCommands[i]
			ExecuteTask(command)
		}
		return nil
	}

	if update.PiperOpcode == common.StopService {
		KillChildProcess()
		return nil
	}

	if update.PiperOpcode == common.StartService {
		go StartChildProcess(update.ExecutableCommands[0])
		return nil
	}

	if update.Operation == watcher.Create || update.Operation == watcher.Write {
		localPath := update.RelativePath
		localPath = strings.Replace(localPath, update.Name, "", -1)

		if update.RelativePath == "" {
			return nil
		}

		if localPath != "" {
			os.MkdirAll("." + localPath, os.ModePerm)
		}

		file, err := os.OpenFile("." + update.RelativePath, os.O_RDWR | os.O_CREATE, 0666)
		if err != nil {
			return err
		}

		defer file.Close()

		_, err = file.Write(update.Content)
		if err != nil {
			return err
		}

		BroadcastMessage("Wrote file " + update.Name)
		return nil
	}

	if update.Operation == watcher.Remove {
		BroadcastMessage("Deleted file " + update.Name)
		return os.Remove("." + update.RelativePath)
		return nil
	}

	BroadcastMessage("Unknown action for " + update.Name)
	return nil
}
