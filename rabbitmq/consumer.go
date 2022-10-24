package rabbitmq

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	mqconnector "github.com/Velnbur/mq-connector"
)

type Consumer struct {
	connector
	handler mqconnector.Handler
}

func NewConsumer(conn *amqp.Connection, queueName string) (mqconnector.Consumer, error) {
	connector, err := newConnector(conn, queueName)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		connector: *connector,
	}, nil
}

func (c *Consumer) SetHandler(handler mqconnector.Handler) mqconnector.Consumer {
	c.handler = handler
	return c
}

func (c *Consumer) Run(ctx context.Context) {
	deliveries, err := c.channel.Consume(
		c.queue.Name,
		"", // consumer name
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		// FIXME: replace with another behavior
		log.Fatal("[ERROR]: failed to create amqp consumer:", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("[INFO]: context canceled, finishing...")
			return
		case delivery := <-deliveries:
			err = c.handler(ctx, delivery.Body)

			if err != nil {
				log.Println("[ERROR]: failed to handle delivery:", err)
				delivery.Reject(true)
				continue
			}

			delivery.Ack(false)
		}
	}
}
