package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/logging"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/server"
)

const (
	ConfigFlag string = "config"
	// SqlStorage        = "sql"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, ConfigFlag, "", "Path to config file")
}

func main() {
	flag.Parse()

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
	/*
		if configuration.Storage.Type == SqlStorage {
		//	connectionToDb := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", configuration.Database.Host, configuration.Database.Port, configuration.Database.User, configuration.Database.Password, configuration.Database.Name)
			repo, err := sql.NewDbConnection(connectionToDb)
			if err != nil {
				log.Fatal(err)
			}
		} else {
		//	repo := new(inmemory.Repo)
		} */

	// Server
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
	webServer := server.NewWebServer(mux, configuration.HTTPServer.Host+":"+configuration.HTTPServer.Port)
	if err := webServer.Start(); err != nil {
		log.Fatal(err)
	}
}
