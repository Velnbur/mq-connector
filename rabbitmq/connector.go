package rabbitmq

import (
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type connector struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

func newConnector(conn *amqp.Connection, queueName string) (*connector, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get channel from conection")
	}

	queue, err := channel.QueueDeclare(queueName, false, true, false, false, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to declare queue")
	}

	return &connector{
		channel: channel,
		queue:   queue,
	}, nil
}
