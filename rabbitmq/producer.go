package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"

	mqconnector "github.com/Velnbur/mq-connector"
)

type Producer struct {
	connector
}

func (p *Producer) Publish(ctx context.Context, data json.RawMessage) error {
	err := p.channel.PublishWithContext(
		ctx,
		"", // exchange
		p.queue.Name,
		false, // mandatory
		false, // immidiate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
	if err != nil {
		return errors.Wrap(err, "failed to publish delivery")
	}
	return nil
}

func NewProducer(conn *amqp.Connection, queueName string) (mqconnector.Producer, error) {
	connector, err := newConnector(conn, queueName)
	if err != nil {
		return nil, err
	}

	return &Producer{
		connector: *connector,
	}, nil
}
