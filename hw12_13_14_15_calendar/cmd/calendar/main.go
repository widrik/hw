package main

import (
	"flag"
	"log"

	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/logging"
)

const (
	ConfigFlag string = "config"
	SqlStorage = "sql"
)

func main() {
	flag.Parse()

	var configFile string
	flag.StringVar(&configFile, ConfigFlag, "", "Path to config file")

	// Config
	configuration, err := config.Init(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Logger
	err = logging.Init(configuration.Logging.Level, configuration.Logging.File)
	if err != nil {
		log.Fatal(err)
	}

	// Storage
	if configuration.Storage.Type == SqlStorage {

	} else {

	}

	// Server


}
