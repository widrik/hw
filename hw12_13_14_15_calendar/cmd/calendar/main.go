package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/baserepo"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/inmemory"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/sql"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/logging"
	grpcserver "github.com/widrik/hw/hw12_13_14_15_calendar/internal/server/grpc"
	httpserver "github.com/widrik/hw/hw12_13_14_15_calendar/internal/server/http"
	"go.uber.org/zap"
)

const ConfigFlag = "config"

var configFile string

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
	var repo baserepo.EventsRepo
	if configuration.Storage.Type == "SqlStorage" {
		connectionToDb := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", configuration.Database.Host, configuration.Database.Port, configuration.Database.User, configuration.Database.Password, configuration.Database.Name)
		repo, err = sql.NewDBConnection(connectionToDb)
		if err != nil {
			zap.L().Error("init repo error", zap.Error(err))
		}
	} else {
		repo = new(inmemory.Repo)
	}

	// App
	calenderApp := app.Calendar{
		Repository: repo,
	}

	serversErrorsCh := make(chan error)

	// Http server
	httpServer := httpserver.NewServer(&calenderApp, net.JoinHostPort(configuration.HTTPServer.Host, configuration.HTTPServer.Port))
	go func() {
		if err := httpServer.Start(); err != nil {
			serversErrorsCh <- err
		}
	}()
	defer func() {
		if err := httpServer.Stop(); err != nil {
			zap.L().Error("stopping of http server error", zap.Error(err))
			return
		}
	}()

	// Grpc server
	grpcServer := grpcserver.NewServer(&calenderApp, net.JoinHostPort(configuration.GRPCServer.Host, configuration.GRPCServer.Port))
	go func() {
		if err := grpcServer.Start(); err != nil {
			serversErrorsCh <- err
		}
	}()
	defer grpcServer.Stop()

	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signalsCh:
		signal.Stop(signalsCh)
		return
	case err = <-serversErrorsCh:
		if err != nil {
			zap.L().Error("server error", zap.Error(err))
		}
		return
	}
}
