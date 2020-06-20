package client

import (
	"github.com/radovskyb/watcher"
	"os"
	"strings"
)

func SyncAllAndClose()  {
	// add all files as new
	appPath, _ := os.Getwd()
	for path, f := range W.WatchedFiles() {
		// create new file
		localPath := path
		localPath = strings.Replace(localPath, appPath, "", -1)
		AddChangedFile(f.Name(), localPath, watcher.Create)
	}

	PushChanges()
}
