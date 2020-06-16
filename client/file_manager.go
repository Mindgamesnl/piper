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

	cancel := false

	// check if a similar one already exists, dont write to just delete it later
	for i := range ChangedFiles {
		file := ChangedFiles[i]

		if file.Path == path && file.Name == name && file.Operation != watcher.Remove && operation == watcher.Remove {
			cancel = true
			ChangedFiles = remove(ChangedFiles, i)
			Log("Ignoring file " + name + " since it has been created and delete in the same sync")
		} else if file.Path == path && file.Name == name {
			cancel = true
		}
	}

	if !cancel {
		ChangedFiles = append(ChangedFiles, cf)
	}

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

func remove(s []ChangedFile, i int) []ChangedFile {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}