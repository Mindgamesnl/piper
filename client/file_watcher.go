package client

import (
	"github.com/radovskyb/watcher"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var W *watcher.Watcher

func StartFileWatcher(callback func())  {
	W = watcher.New()

	expression := "^.*\\.("
	for i := range LoadedInstance.WatchedExtensions {
		ex := LoadedInstance.WatchedExtensions[i]
		expression += strings.ToLower(ex) + "|"
		expression += strings.ToUpper(ex) + "|"
	}
	expression = trimSuffix(expression, "|")
	expression += ")$"
	r := regexp.MustCompile(expression)
	W.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-W.Event:
				path, _ := os.Getwd()
				if event.Path == path && event.FileInfo.Name() == "piper" {

				} else {

					if event.Op == watcher.Move {
						// fake move into two seperate events
						// deletion
						localPath := event.OldPath
						localPath = strings.Replace(localPath, path, "", -1)
						AddChangedFile(event.FileInfo.Name(), localPath, watcher.Remove)

						// create new file
						localPath = event.Path
						localPath = strings.Replace(localPath, path, "", -1)
						AddChangedFile(event.FileInfo.Name(), localPath, watcher.Create)
					} else {
						localPath := event.Path
						localPath = strings.Replace(localPath, path, "", -1)
						AddChangedFile(event.FileInfo.Name(), localPath, event.Op)
					}
				}
			case err := <-W.Error:
				log.Fatalln(err)
			case <-W.Closed:
				return
			}
		}
	}()

	path, _ := os.Getwd()
	path += "/"

	if err := W.Add(path); err != nil {
		log.Fatalln(err)
	}

	if err := W.AddRecursive("."); err != nil {
		log.Fatalln(err)
	}

	for i := range LoadedInstance.IgnoredDirectories {
		dir := LoadedInstance.IgnoredDirectories[i]
		path, _ := os.Getwd()
		path += "/"
		Log("Ignoring directory: " + path + dir)
		W.RemoveRecursive(path + dir)
	}

	callback()

	// Start the watching process - it'll check for changes every 100ms.
	if err := W.Start(time.Millisecond * 500); err != nil {
		log.Fatalln(err)
	}
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
