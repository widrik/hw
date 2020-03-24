package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	exactTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("current time: %v\n", time.Now())
	fmt.Printf("exact time: %v\n", exactTime)
}
