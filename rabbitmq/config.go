package rabbitmq

import (
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

const (
	rabbitYamlKey = "rabbit"
)

type RabbiteConfiger interface {
	RabbitConnection() *amqp.Connection
}

func NewRabbitConfiger(getter kv.Getter) RabbiteConfiger {
	return &rabbitConfiger{
		getter: getter,
	}
}

type rabbitConfiger struct {
	getter kv.Getter
	once   comfig.Once
}

type rabbitConfig struct {
	URL string `fig:"url,required"`
}

func (r *rabbitConfiger) RabbitConnection() *amqp.Connection {
	return r.once.Do(func() interface{} {
		var cfg rabbitConfig

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
