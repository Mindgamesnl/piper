package server

import (
	"io"
	"os/exec"
	"strings"
)

func ExecuteTask(command string)  {
	// split between commands and arguments
	parts := strings.Split(command, " ")
	var arguments []string
	for i := range parts {
		if i == 0 {
			continue
		}
		arguments = append(arguments, parts[i])
	}

	cmd := exec.Command(parts[0], arguments...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		BroadcastCommandError(command, err.Error())
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		BroadcastCommandError(command, err.Error())
	}

	feedback := string(out)
	feedback = strings.TrimSuffix(feedback, "\n")

	BroadcastCommandOutput(feedback)
}
