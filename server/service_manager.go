package server

import (
	"bufio"
	"os/exec"
	"strings"
	"syscall"
)

var Command *exec.Cmd

func StartChildProcess(command string)  {
	parts := strings.Split(command, " ")
	var arguments = ""
	for i := range parts {
		if i == 0 {
			continue
		}
		arguments += parts[i] + " "
	}

	Command = exec.Command(parts[0], arguments)
	Command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmdReader, err := Command.StdoutPipe()
	if err != nil {
		BroadcastCommandError(command, err.Error())
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			BroadcastCommandOutput(scanner.Text())
		}
	}()

	if err := Command.Start(); err != nil {
		BroadcastCommandError(command, err.Error())
	}

	if err := Command.Wait(); err != nil {
		BroadcastCommandError(command, err.Error())
	}
}

func KillChildProcess()  {
	if Command == nil {
		return
	}
	syscall.Kill(-Command.Process.Pid, syscall.SIGKILL)
	Command = nil
}