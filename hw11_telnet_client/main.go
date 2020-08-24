package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const argsLen = 4

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 0, "connection timeout")
}

func main() {
	log.Print("Start")
	flag.Parse()

	log.Print("Args check")
	args := flag.Args()
	if len(args) != argsLen {
		log.Fatal("few arguments: should be " + strconv.Itoa(argsLen))
	}

	log.Print("Args check was ok")
	host := args[0]
	port := args[1]

	client := NewTelnetClient(
		net.JoinHostPort(host, port),
		timeout,
		os.Stdin,
		os.Stdout,
	)
	errorsCh := make(chan error)
	signlsCh := make(chan os.Signal, 1)

	signal.Notify(signlsCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		errorsCh <- client.Send()
	}()
	go func() {
		errorsCh <- client.Receive()
	}()

	select {
	case <-signlsCh:
		signal.Stop(signlsCh)

		return
	case err = <-errorsCh:
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stderr, "...EOF\n")
		
		return
	}
}
