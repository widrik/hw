package sender

import (
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/entities"
)

type Sender struct {
	service entities.EvenetsSeviceInterface
}

func NewSender(handler entities.EvenetsSeviceInterface) Sender {
	return Sender{
		service: handler,
	}
}

func (sender *Sender) Listen(context context.Context) error {
	return sender.service.Listen(
		context,
		func(deliveredMessages <-chan amqp.Delivery) {
			log.Printf("get events")
			for message := range deliveredMessages {
				messageText := message.Body
				event := &entities.Event{}
				err := json.Unmarshal(messageText, event)
				if err != nil {
					log.Printf("parse error: %s", err)
				} else {
					log.Printf("new event %s", messageText)
				}
			}
		},
	)
}

func (sender *Sender) Stop() error {
	err := sender.service.Stop()
	if err != nil {
		return err
	}

	return nil
}
