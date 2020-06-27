package main

import (
	"log"
	"os"
	"os/exec"
)

const (
	cmdExistsLenCheck = 1
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < cmdExistsLenCheck {
		returnCode = -1
	} else {
		command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
		command.Env = os.Environ()
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		for envName, envValue := range env {
			command.Env = append(command.Env, envName+"="+envValue)
		}

		err := command.Run()
		if err != nil {
			log.Fatal(err)
		}
		returnCode = command.ProcessState.ExitCode()
	}

	return returnCode
}
