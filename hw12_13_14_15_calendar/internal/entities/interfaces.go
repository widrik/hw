package entities

import (
	"context"

	"github.com/streadway/amqp"
)

type EvenetsSeviceInterface interface {
	Connect() error
	Reconnect(context.Context) (<-chan amqp.Delivery, error)
	Stop() error
	Listen(context.Context, func(<-chan amqp.Delivery)) error
}

type PublisherInterface interface {
	Connect() error
	Send(event *Event) error
	Close() error
}

type SchedulerInterface interface {
	Process()
	Publish(event *Event) error
	Stop() error
}
