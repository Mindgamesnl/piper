package server

import (
	"fmt"
	"github.com/go-cmd/cmd"
	"os"
	"strings"
)

var Command *cmd.Cmd


func StartChildProcess(command string)  {
	parts := strings.Split(command, " ")
	var arguments = ""
	for i := range parts {
		if i == 0 {
			continue
		}
		arguments += parts[i] + " "
	}

	arguments = strings.TrimSuffix(arguments, " ")

	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}

	Command = cmd.NewCmdOptions(cmdOptions, parts[0], arguments)

	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		// Done when both channels have been closed
		// https://dave.cheney.net/2013/04/30/curious-channels
		for Command.Stdout != nil || Command.Stderr != nil {
			select {
			case line, open := <-Command.Stdout:
				if !open {
					Command.Stdout = nil
					continue
				}
				BroadcastServiceOutput("SERVICE >> " + line)
			case line, open := <-Command.Stderr:
				if !open {
					Command.Stderr = nil
					continue
				}
				fmt.Fprintln(os.Stderr, line)
				BroadcasServiceError(line)
			}
		}
	}()
	<-Command.Start()
	<-doneChan
}

func KillChildProcess()  {
	if Command == nil {
		return
	}
	Command.Stop()
	Command = nil
}