package main

import (
	"log"
	"os"
)

const (
	minArgsCount = 1
)

func main() {
	fileName := os.Getenv("GOFILE")

	if len(fileName) == 0 {
		if len(os.Args) != minArgsCount {
			log.Fatalf("Argument amount error: minimum %d arguments are required", minArgsCount)
		}

		fileName = os.Args[1]
	}
	err := Generate(fileName)
	if err != nil {
		log.Fatal(err)
	}
}
