package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var command *exec.Cmd
	strCommand := cmd[0]

	if len(cmd) > 1 {
		command = exec.Command(strCommand, cmd[1:]...)
	} else {
		command = exec.Command(strCommand)
	}

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	for key, val := range env {
		if val.NeedRemove {
			_ = os.Unsetenv(key)
		} else {
			_ = os.Setenv(key, val.Value)
		}
	}

	_ = command.Run()
	return command.ProcessState.ExitCode()
}
