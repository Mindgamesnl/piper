package client

import (
	"github.com/Mindgamesnl/piper/common"
	"github.com/radovskyb/watcher"
	"io/ioutil"
	"strconv"
	"time"
)

var (
	ChangedFiles []ChangedFile
	IsSyncing    = false
)

const (
	InfoColor    = "\033[1;34m"
	NoticeColor  = "\033[1;36m"
	WarningColor = "\033[1;33m"
	ErrorColor   = "\033[1;31m"
	DebugColor   = "\033[0;36m"
)

type ChangedFile struct {
	Name      string
	Path      string
	Operation watcher.Op
}

func AddChangedFile(name string, path string, operation watcher.Op) {
	cf := ChangedFile{
		Name:      name,
		Path:      path,
		Operation: operation,
	}

	cancel := false
	changedFiles := ChangedFiles
	// check if a similar one already exists, dont write to just delete it later
	for i := range changedFiles {
		file := changedFiles[i]

		if file.Path == path && file.Name == name && file.Operation != watcher.Remove && operation == watcher.Remove {
			cancel = true
			changedFiles = remove(changedFiles, i)
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

func reRender() {
	FilesView.Clear()
	for i := range ChangedFiles {
		if i > len(ChangedFiles) {
			Log(ErrorColor + "Canceled rendering since the file update index is no longer relevant")
			return
		}
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

func InitManager() {
	if !LoadedInstance.AutoSyncEnabled {
		Log(ErrorColor + "Auto reloading is disabled. Hit <ENTER> to synchronize changed files.")
		return
	}

	Log(DebugColor + "Auto reloading is enabled. Hit <ENTER> to synchronize changed files or wait for the " + strconv.Itoa(LoadedInstance.AutoSyncTimeout) + " second interval")
	ticker := time.NewTicker(time.Duration(LoadedInstance.AutoSyncTimeout) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				PushChanges()
				reRender()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func PushChanges() {
	if IsSyncing {
		Log(ErrorColor + "Cancelled sync because another task is still running")
		return
	}
	changeCount := len(ChangedFiles)
	if changeCount == 0 {
		return
	}

	IsSyncing = true

	// pre commands
	var preCommands = common.FileUpdate{
		PiperOpcode:        common.ExecuteCommands,
		ExecutableCommands: LoadedInstance.PreUpdateCommands,
	}
	WriteSocket(preCommands.ToJson())

	var killChild = common.FileUpdate{
		PiperOpcode: common.StopService,
	}
	WriteSocket(killChild.ToJson())

	handled := 1
	for i := range ChangedFiles {
		file := ChangedFiles[i]
		Log("Handling file " + strconv.Itoa(handled) + " of " + strconv.Itoa(len(ChangedFiles)))
		if file.Operation == watcher.Write || file.Operation == watcher.Create {
			// get contents
			content, _ := ioutil.ReadFile("." + file.Path)

			var update = common.FileUpdate{
				Name:         file.Name,
				RelativePath: file.Path,
				Operation:    file.Operation,
				Content:      content,
			}

			WriteSocket(update.ToJson())
			handled++
		}

		if file.Operation == watcher.Remove {
			var update = common.FileUpdate{
				Name:         file.Name,
				RelativePath: file.Path,
				Operation:    file.Operation,
			}

			WriteSocket(update.ToJson())
			handled++
		}
	}

	var postCommands = common.FileUpdate{
		PiperOpcode:        common.ExecuteCommands,
		ExecutableCommands: LoadedInstance.PostUpdateCommands,
	}
	WriteSocket(postCommands.ToJson())

	var startChild = common.FileUpdate{
		PiperOpcode: common.StartService,
		ExecutableCommands: []string{
			LoadedInstance.ServiceCommand,
		},
	}
	WriteSocket(startChild.ToJson())

	ChangedFiles = []ChangedFile{}
	IsSyncing = false
	reRender()
}

func remove(s []ChangedFile, i int) []ChangedFile {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
