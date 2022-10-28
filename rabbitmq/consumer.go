package rabbitmq

import (
	"context"
	"log"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"

	mqc "github.com/Velnbur/mq-connector"
)

var _ mqc.Consumer = &RabbitConsumer{}

type RabbitConsumer struct {
	connector
	contexters []mqc.ContextFunc
	handler    mqc.Handler
}

func NewRabbitConsumer(conn *amqp.Connection, queueName string) (mqc.Consumer, error) {
	connector, err := newConnector(conn, queueName)
	if err != nil {
		return nil, err
	}

	return &RabbitConsumer{
		connector: *connector,
	}, nil
}

func (c *RabbitConsumer) SetContexters(ctxers ...mqc.ContextFunc) mqc.Consumer {
	c.contexters = ctxers
	return c
}

func (c *RabbitConsumer) SetHandler(handler mqc.Handler) mqc.Consumer {
	c.handler = handler
	return c
}

func (c *RabbitConsumer) Close() error {
	return errors.Wrap(c.channel.Close(), "failed to close channel")
}

func (c *RabbitConsumer) Run(ctx context.Context) {
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
			ctx := c.applyContexts(ctx)
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

// apply contexters to context
func (c *RabbitConsumer) applyContexts(ctx context.Context) context.Context {
	for _, contexter := range c.contexters {
		ctx = contexter(ctx)
	}
	return ctx
}
