package rabbitmq

import (
	mqconnector "github.com/Velnbur/mq-connector"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

const (
	rabbitYamlKey = "rabbit"
)

type RabbitConnectioner interface {
	RabbitConnection() *amqp.Connection
}

func NewRabbitConfiger(getter kv.Getter) *RabbitConfiger {
	return &RabbitConfiger{
		getter: getter,
	}
}

type RabbitConfiger struct {
	getter kv.Getter
	once   comfig.Once
}

type rabbitConnectionConfig struct {
	URL string `fig:"url,required"`
}

func (r *RabbitConfiger) RabbitConnection() *amqp.Connection {
	return r.once.Do(func() interface{} {
		var cfg rabbitConnectionConfig

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(r.getter, rabbitYamlKey)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to parse rabbit config"))
		}

		conn, err := amqp.Dial(cfg.URL)
		if err != nil {
			panic(errors.Wrap(err, "failed to create rabbit connection"))
		}

		return conn
	}).(*amqp.Connection)
}

const rabbitRouterYamlKey = "rabbit_router"

type RabbitRouterer interface {
	RabbitRouter() *RabbitRouter
}

type rabbitRouterConfig struct {
	ResponsesQueue string `fig:"responses_queue,required"`
	RequestsQueue  string `fig:"requests_queue,required"`
}

func (r *RabbitConfiger) RabbitRouter() *RabbitRouter {
	return r.once.Do(func() interface{} {
		var cfg rabbitRouterConfig

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(r.getter, rabbitRouterYamlKey)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to parse rabbit router config"))
		}

		conn, err := NewRouter(r.RabbitConnection(), cfg.RequestsQueue, cfg.ResponsesQueue)
		if err != nil {
			panic(errors.Wrap(err, "failed to create rabbit router"))
		}

		return conn
	}).(*RabbitRouter)
}

const rabbitConsumerYamlKey = "rabbit_consumer"

type RabbitConsumerer interface {
	RabbitConsumer() mqconnector.Consumer
}

type rabbitQueueConfig struct {
	Queue string `fig:"queue,required"`
}

func (r *RabbitConfiger) RabbitConsumer() mqconnector.Consumer {
	return r.once.Do(func() interface{} {
		var cfg rabbitQueueConfig

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(r.getter, rabbitConsumerYamlKey)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to parse rabbit consumer config"))
		}

		cons, err := NewConsumer(r.RabbitConnection(), cfg.Queue)
		if err != nil {
			panic(errors.Wrap(err, "failed to create rabbit consumer"))
		}

		return cons
	}).(mqconnector.Consumer)
}

const rabbitProducerYamlKey = "rabbit_producer"

type RabbitProducer interface {
	RabbitConsumer() mqconnector.Producer
}

func (r *RabbitConfiger) RabbitProducer() mqconnector.Producer {
	return r.once.Do(func() interface{} {
		var cfg rabbitQueueConfig

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(r.getter, rabbitProducerYamlKey)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to parse rabbit producer config"))
		}

		cons, err := NewProducer(r.RabbitConnection(), cfg.Queue)
		if err != nil {
			panic(errors.Wrap(err, "failed to create rabbit producer"))
		}

		return cons
	}).(mqconnector.Producer)
}
