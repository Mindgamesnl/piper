package client

import (
	"github.com/radovskyb/watcher"
	"strconv"
)

var ChangedFiles []ChangedFile

const (
	InfoColor    = "\033[1;34m"
	NoticeColor  = "\033[1;36m"
	WarningColor = "\033[1;33m"
	ErrorColor   = "\033[1;31m"
	DebugColor   = "\033[0;36m"
)


type ChangedFile struct {
	Name string
	Path string
	Operation watcher.Op
}

func AddChangedFile(name string, path string, operation watcher.Op) {
	cf := ChangedFile{
		Name: name,
		Path: path,
		Operation: operation,
	}

	ChangedFiles = append(ChangedFiles, cf)
	
	reRender()
}

func reRender()  {
	FilesView.Clear()
	for i := range ChangedFiles {
		file := ChangedFiles[i]

		if file.Operation == watcher.Remove {
			PrintFiles(ErrorColor + "[-] " + file.Name)
			continue
		}

		if file.Operation == watcher.Create {
			PrintFiles(NoticeColor + "[+] " + file.Name)
			continue
		}

		PrintFiles(DebugColor + "[~] " + file.Name)
	}
}

func PushChanges()  {
	Log("Pushing " + strconv.Itoa(len(ChangedFiles)) + " file updates...")
	ChangedFiles = []ChangedFile{}
	reRender()
}