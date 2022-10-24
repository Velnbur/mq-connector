package rabbitmq

import (
	"context"
	"encoding/json"
	mqconnector "github.com/Velnbur/mq-connector"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel *amqp.Channel

	subscription string
}

func (p *Publisher) Publish(ctx context.Context, data json.RawMessage) error {
	err := p.channel.PublishWithContext(ctx,
		p.subscription, // exchange
		"",             // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
	if err != nil {
		return errors.Wrap(err, "failed publish message")
	}
	return nil
}

func NewPublisher(conn *amqp.Connection, subscription string) (mqconnector.Producer, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get channel from connection")
	}

	err = channel.ExchangeDeclare(
		subscription, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to declare exchange")
	}

	return &Publisher{
		channel:      channel,
		subscription: subscription,
	}, nil
}
