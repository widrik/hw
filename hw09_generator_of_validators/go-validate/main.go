package main

import (
	"flag"
	"log"
	"os"
)

const (
	minArgsCount = 1
)

func main() {
	flag.Parse()

	if len(os.Args) != minArgsCount {
		log.Fatalf("Argument amount error: minimum %d arguments are required", minArgsCount)
	}

	fileName := os.Args[1]
	err := Generate(fileName)
	if err != nil {
		log.Fatal(err)
	}
}
