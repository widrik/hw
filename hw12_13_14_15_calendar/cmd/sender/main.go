package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	appSender "github.com/widrik/hw/hw12_13_14_15_calendar/internal/app/sender"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/logging"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/rabbit"
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
	configuration, err := config.InitSenderConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Logger
	initLogging(configuration)

	// Rabbit hndler
	rabbitHandler, err := rabbit.InitService(configuration)
	if err != nil {
		zap.L().Error("fatal error", zap.Error(err))
	}

	// Sender
	sender := appSender.NewSender(rabbitHandler)
	defer stopSender(&sender)

	errsCh := make(chan error, 1)
	notifiesCh := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())

	signal.Notify(notifiesCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		errsCh <- sender.Listen(ctx)
	}()

	select {
	case <-notifiesCh:
		signal.Stop(notifiesCh)
		cancel()

		return
	case err = <-errsCh:
		if err != nil {
			zap.L().Error("fatal error", zap.Error(err))
		}

		cancel()
	}
}

func initLogging(configuration config.RabbitConfiguration) {
	err := logging.Init(configuration.Logging.Level, configuration.Logging.File)
	if err != nil {
		log.Fatal(err)
	}
}

func stopSender(sender *appSender.Sender) {
	err := sender.Stop()
	if err != nil {
		err := sender.Stop()
		if err != nil {
			zap.L().Error("stop sender err", zap.Error(err))
		}
	}
}
