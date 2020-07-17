package main

import (
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

type Environment map[string]string

var regexpSplit = regexp.MustCompile(`[\r\n]+`)

func splitTextToParts(inputText string) []string {
	return regexpSplit.Split(inputText, -1)
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envValues := make(Environment)

	for _, file := range files {
		fileName := file.Name()
		filePath := path.Join(dir, fileName)

		fileData, err := ioutil.ReadFile(filePath)

		if err == nil {
			fileString := splitTextToParts(string(fileData))[0]
			fileString = strings.ReplaceAll(fileString, string(0), "\n")
			fileString = strings.TrimRight(fileString, "\t")
			fileString = strings.TrimRight(fileString, " ")

			envValues[fileName] = fileString
		}
	}

	return envValues, nil
}
