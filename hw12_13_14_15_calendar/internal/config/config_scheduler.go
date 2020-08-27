package config

import (
	"github.com/spf13/viper"
)

type SchedulerConfiguration struct {
	Host         string
	Port         string
	QueueName    string
	ExchangeName string
	Logging      Logging
	GRPCServer   GRPCServer
}

func InitSchedulerConfig(path string) (SchedulerConfiguration, error) {
	var configuration SchedulerConfiguration

	if path == "" {
		return configuration, ErrFilePathEmpty
	}

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return configuration, ErrReadFile
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		return configuration, ErrReadFile
	}

	return configuration, nil
}
