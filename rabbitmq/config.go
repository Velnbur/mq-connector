package rabbitmq

import (
	mqconnector "github.com/Velnbur/mq-connector"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

func NewRabbitConfiger(getter kv.Getter) *RabbitConfiger {
	return &RabbitConfiger{
		getter: getter,
	}
}

type RabbitConfiger struct {
	getter         kv.Getter
	onceConnection comfig.Once
	onceProducer   comfig.Once
	onceRouter     comfig.Once
	onceConsumer   comfig.Once
}

const (
	rabbitYamlKey = "rabbit"
)

var _ RabbitConnectioner = &RabbitConfiger{}

type RabbitConnectioner interface {
	RabbitConnection() *amqp.Connection
}

type rabbitConnectionConfig struct {
	URL string `fig:"url,required"`
}

func (r *RabbitConfiger) RabbitConnection() *amqp.Connection {
	return r.onceConnection.Do(func() interface{} {
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

var _ RabbitRpcServerer = &RabbitConfiger{}

type RabbitRpcServerer interface {
	RabbitRpcServer(yamlKey string) *RabbitRpcServer
}

type rabbitRouterConfig struct {
	ResponsesQueue string `fig:"responses_queue,required"`
	RequestsQueue  string `fig:"requests_queue,required"`
}

func (r *RabbitConfiger) RabbitRpcServer(yamlKey string) *RabbitRpcServer {
	return r.onceRouter.Do(func() interface{} {
		var cfg rabbitRouterConfig

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(r.getter, yamlKey)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to parse rabbit router config"))
		}

		conn, err := NewRouter(r.RabbitConnection(), cfg.RequestsQueue, cfg.ResponsesQueue)
		if err != nil {
			panic(errors.Wrap(err, "failed to create rabbit router"))
		}

		return conn
	}).(*RabbitRpcServer)
}

var _ RabbitConsumerer = &RabbitConfiger{}

type RabbitConsumerer interface {
	RabbitConsumer(yamlKey string) mqconnector.Consumer
}

type rabbitQueueConfig struct {
	Queue string `fig:"queue,required"`
}

func (r *RabbitConfiger) RabbitConsumer(yamlKey string) mqconnector.Consumer {
	return r.onceConsumer.Do(func() interface{} {
		var cfg rabbitQueueConfig

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(r.getter, yamlKey)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to parse rabbit consumer config"))
		}

		cons, err := NewRabbitConsumer(r.RabbitConnection(), cfg.Queue)
		if err != nil {
			panic(errors.Wrap(err, "failed to create rabbit consumer"))
		}

		return cons
	}).(mqconnector.Consumer)
}

var _ RabbitProducerer = &RabbitConfiger{}

type RabbitProducerer interface {
	RabbitProducer(yamlKey string) mqconnector.Producer
}

func (r *RabbitConfiger) RabbitProducer(yamlKey string) mqconnector.Producer {
	return r.onceProducer.Do(func() interface{} {
		var cfg rabbitQueueConfig

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(r.getter, yamlKey)).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to parse rabbit producer config"))
		}

		cons, err := NewRabbitProducer(r.RabbitConnection(), cfg.Queue)
		if err != nil {
			panic(errors.Wrap(err, "failed to create rabbit producer"))
		}

		return cons
	}).(mqconnector.Producer)
}
