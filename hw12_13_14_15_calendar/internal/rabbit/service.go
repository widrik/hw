package rabbit

import (
	"context"
	"net"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/streadway/amqp"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/entities"
	"go.uber.org/zap"
)

type ServiceRabbit struct {
	address      string
	connection   *amqp.Connection
	connectionCh *amqp.Channel
	queueName    string
	exchangeName string
	doneCh       chan error
}

func InitService(configuration config.RabbitConfiguration) (entities.EvenetsSeviceInterface, error) {
	service := &ServiceRabbit{
		doneCh:       make(chan error),
		queueName:    configuration.QueueName,
		exchangeName: configuration.ExchangeName,
	}

	if configuration.Host == "" || configuration.Port == "" {
		return service, ErrNotValidAddressData
	}
	service.address = "amqp://" + net.JoinHostPort(configuration.Host, configuration.Port)

	return service, nil
}

func (service *ServiceRabbit) Connect() error {
	connection, err := amqp.Dial(service.address)
	if err != nil {
		return err
	}
	service.connection = connection

	connectionCh, err := service.connection.Channel()
	if err != nil {
		return err
	}
	service.connectionCh = connectionCh

	notifies := make(chan *amqp.Error)

	zap.L().Warn("connecting to channel")
	go func() {
		zap.L().Warn("closing channel")
		zap.L().Error("error: ", zap.Error(<-service.connection.NotifyClose(notifies)))
		service.doneCh <- ErrChannelClosed
	}()
	zap.L().Info("connected was successful")



	if err != nil {
		return err
	}

	err = service.connectionCh.ExchangeDeclare(service.exchangeName, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = service.connectionCh.QueueBind(service.queueName, "", service.exchangeName, false, nil)
	if err != nil {
		return err
	}

	return nil
}

func (service *ServiceRabbit) reconnect(context context.Context) (<-chan amqp.Delivery, error) {
	expBackOff := backoff.NewExponentialBackOff()
	expBackOff.MaxElapsedTime = MaxElapsedTime
	expBackOff.MaxInterval = MaxInterval
	expBackOff.InitialInterval = InitialInterval
	expBackOff.Multiplier = Multiplier

	backOffContext := backoff.WithContext(expBackOff, context)
	for {
		durationToWait := backOffContext.NextBackOff()
		if durationToWait == backoff.Stop {
			return nil, ErrReconnect
		}

		select {
		case <-context.Done():
			return nil, nil
		case <-time.After(durationToWait):
			err := service.Connect()
			if err != nil {
				zap.L().Error("error: ", zap.Error(err))

				continue
			}
			deliveredMessages, err := service.connectionCh.Consume(service.queueName, "", false, false, false, false, nil)
			if err != nil {
				zap.L().Error("error: ", zap.Error(err))

				continue
			}

			return deliveredMessages, nil
		}
	}
}

func (service *ServiceRabbit) Stop() error {
	err := service.connectionCh.Close()
	if err != nil {
		return err
	}

	err = service.connection.Close()
	if err != nil {
		return err
	}

	return nil
}

func (service *ServiceRabbit) Listen(context context.Context, msgsHandler func(<-chan amqp.Delivery)) error {
	deliveredMessages, err := service.reconnect(context)
	if err != nil {
		return err
	}

	for {
		go msgsHandler(deliveredMessages)

		if <-service.doneCh != nil {
			deliveredMessages, err = service.reconnect(context)
			if err != nil {
				return err
			}
		}
	}
}
