package main

import (
	"log"
	"os"
)

const (
	argsLenWithoutDirAndEnv = 2
	argsLenWithoutEnv       = 2
)

func main() {
	args := os.Args

	if len(args) < argsLenWithoutDirAndEnv {
		log.Fatal("Few arguments: dir to environment and execute command are missed")
	} else if len(args) < argsLenWithoutEnv {
		log.Fatal("Few arguments: execute command is missed")
	}

	environmentPath := args[1]
	command := args[2:]

	environment, err := ReadDir(environmentPath)

	if err != nil {
		log.Fatal("Environment error: $d", err)
	}

	os.Exit(RunCmd(command, environment))
}
