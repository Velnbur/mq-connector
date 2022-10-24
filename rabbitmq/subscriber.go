package rabbitmq

import (
	mqconnector "github.com/Velnbur/mq-connector"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Subscriber struct {
	Consumer
}

func NewSubscriber(conn *amqp.Connection, subscription string) (mqconnector.Consumer, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get channel from conection")
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

	queue, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to declare queue")
	}

	err = channel.QueueBind(
		queue.Name,   // queue name
		"",           // routing key
		subscription, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to bind queue")
	}

	return &Subscriber{
		Consumer: Consumer{
			connector: connector{
				channel: channel,
				queue:   queue,
			},
		},
	}, nil
}
