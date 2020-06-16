package client

import (
	"github.com/radovskyb/watcher"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

func StartFileWatcher()  {
	w := watcher.New()

	expression := "^.*\\.("
	for i := range LoadedInstance.WatchedExtensions {
		ex := LoadedInstance.WatchedExtensions[i]
		expression += strings.ToLower(ex) + "|"
		expression += strings.ToUpper(ex) + "|"
	}
	expression = trimSuffix(expression, "|")
	expression += ")$"
	r := regexp.MustCompile(expression)
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				path, _ := os.Getwd()
				if event.Path == path && event.FileInfo.Name() == "piper" {

				} else {
					AddChangedFile(event.FileInfo.Name(), event.Path, event.Op)
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	path, _ := os.Getwd()
	path += "/"

	if err := w.Add(path); err != nil {
		log.Fatalln(err)
	}

	if err := w.AddRecursive("."); err != nil {
		log.Fatalln(err)
	}

	for i := range LoadedInstance.IgnoredDirectories {
		dir := LoadedInstance.IgnoredDirectories[i]
		path, _ := os.Getwd()
		path += "/"
		Log("Ignoring directory: " + path + dir)
		w.RemoveRecursive(path + dir)
	}

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}