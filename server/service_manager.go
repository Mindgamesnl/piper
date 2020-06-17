package server

import (
	"bufio"
	"github.com/sirupsen/logrus"
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

	arguments = strings.TrimSuffix(arguments, " ")

	Command = exec.Command(parts[0], arguments)
	Command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmdReader, err := Command.StdoutPipe()
	cmdErrr, _ := Command.StderrPipe()
	if err != nil {
		logrus.Error(err)
		BroadcastCommandError(command, err.Error())
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			BroadcastServiceOutput("SERVICE >> " + scanner.Text())
		}
	}()

	errorScanner := bufio.NewScanner(cmdErrr)
	go func() {
		for errorScanner.Scan() {
			BroadcastCommandOutput(errorScanner.Text())
		}
	}()

	if err := Command.Start(); err != nil {
		logrus.Error(err)
		BroadcastCommandError(command, err.Error())
	}

	if err := Command.Wait(); err != nil {
		logrus.Error(err)
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
