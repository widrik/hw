package rabbit

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/streadway/amqp"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/entities"
	"go.uber.org/zap"
)

type Publisher struct {
	address      string
	connection   *amqp.Connection
	connectionCh *amqp.Channel
	queueName    string
	exchangeName string
	doneCh       chan error
}

func InitPublisher(configuration config.SchedulerConfiguration) (entities.PublisherInterface, error) {
	publisher := &Publisher{
		queueName:    configuration.QueueName,
		exchangeName: configuration.ExchangeName,
		doneCh:       make(chan error),
	}

	if configuration.Host == "" || configuration.Port == "" {
		return publisher, ErrNotValidAddressData
	}
	publisher.address = "amqp://" + net.JoinHostPort(configuration.Host, configuration.Port)

	return publisher, nil
}

func (publisher *Publisher) Connect() error {
	var err error

	connection, err := amqp.Dial(publisher.address)
	if err != nil {
		return err
	}
	publisher.connection = connection
	zap.L().Warn("got connection")

	publisher.connectionCh, err = publisher.connection.Channel()
	if err != nil {
		return err
	}
	zap.L().Warn("got channel")

	notifies := make(chan *amqp.Error)

	go func() {
		zap.L().Warn("closing channel")
		zap.L().Error("error: ", zap.Error(<-publisher.connection.NotifyClose(notifies)))
		publisher.doneCh <- ErrChannelClosed
	}()

	err = publisher.connectionCh.ExchangeDeclare(publisher.exchangeName, amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		return err
	}

	_, err = publisher.connectionCh.QueueDeclare(publisher.queueName, false, false, false, false, nil)
	if err != nil {
		return ErrQueue
	}

	return nil
}

func (publisher *Publisher) Close() error {
	err := publisher.connectionCh.Close()
	if err != nil {
		return err
	}

	err = publisher.connection.Close()
	if err != nil {
		return err
	}

	return nil
}

func (publisher *Publisher) Send(event *entities.Event) error {
	expBackOff := backoff.NewExponentialBackOff()
	expBackOff.MaxElapsedTime = MaxElapsedTime
	expBackOff.MaxInterval = MaxInterval
	expBackOff.InitialInterval = InitialInterval
	expBackOff.Multiplier = Multiplier

	b := backoff.WithContext(expBackOff, context.Background())
	for {
		durationToWait := b.NextBackOff()
		if durationToWait == backoff.Stop {
			return ErrReconnect
		}

		//nolint:gosimple
		select {
		case <-time.After(durationToWait):
			err := publisher.Connect()
			if err != nil {
				zap.L().Error("error: ", zap.Error(err))

				continue
			}

			eventData, err := json.Marshal(event)
			if err != nil {
				zap.L().Error("event data error: ", zap.Error(err))

				continue
			}

			msg := amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         eventData,
			}

			err = publisher.connectionCh.Publish(publisher.exchangeName, publisher.queueName, false, false, msg)
			if err != nil {
				zap.L().Error("error: ", zap.Error(err))

				continue
			}
			zap.L().Info("message was sent")

			return nil
		}
	}
}
