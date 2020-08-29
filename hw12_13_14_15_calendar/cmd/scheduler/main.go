package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/widrik/hw/hw12_13_14_15_calendar/api/spec"
	schedulerApp "github.com/widrik/hw/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/logging"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/rabbit"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const ConfigFlag = "config"

var configFile string

func init() {
	flag.StringVar(&configFile, ConfigFlag, "", "Path to config file")
}

func main() {
	flag.Parse()

	// Config
	configuration, err := config.InitSchedulerConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Logger
	initLogging(configuration)

	publisher, err := rabbit.InitPublisher(configuration)
	if err != nil {
		zap.L().Error("fatal error", zap.Error(err))
	}

	grpcConnection, err := grpc.Dial(net.JoinHostPort(configuration.GRPCServer.Host, configuration.GRPCServer.Port), grpc.WithInsecure())
	if err != nil {
		zap.L().Error("connection error", zap.Error(err))
	}
	defer func() {
		err := grpcConnection.Close()
		if err != nil {
			zap.L().Error("fatal error", zap.Error(err))
		}
	}()

	client := spec.NewCalendarServiceClient(grpcConnection)
	scheduler := schedulerApp.InitScheduler(publisher, client)
	defer func() {
		err := scheduler.Stop()
		if err != nil {
			zap.L().Error("fatal error", zap.Error(err))
		}
	}()

	notifyCh := make(chan os.Signal, 1)
	signal.Notify(notifyCh, syscall.SIGINT, syscall.SIGTERM)

	notifyTick := time.NewTicker(5 * time.Second)
	defer notifyTick.Stop()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for {
			select {
			case <-notifyCh:
				signal.Stop(notifyCh)

				return
			case <-notifyTick.C:
				scheduler.Process()
			case <-ctx.Done():
				return
			}
		}
	}()
	<-notifyCh
}

func initLogging(configuration config.SchedulerConfiguration) {
	err := logging.Init(configuration.Logging.Level, configuration.Logging.File)
	if err != nil {
		log.Fatal(err)
	}
}
