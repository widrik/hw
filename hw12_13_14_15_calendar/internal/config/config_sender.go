package config

import (
	"github.com/spf13/viper"
)

type RabbitConfiguration struct {
	Host         string
	Port         string
	QueueName    string
	ExchangeName string
	Logging      Logging
}

func InitSenderConfig(path string) (RabbitConfiguration, error) {
	var configuration RabbitConfiguration

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
