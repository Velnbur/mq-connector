package rabbitmq

import (
	"context"
	"sync"

	mqconnector "github.com/Velnbur/mq-connector"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

var _ mqconnector.Router = &RabbitRouter{}

type rabbitRoute struct {
	handler  mqconnector.Handler
	contexts []mqconnector.ContextFunc
}

type RabbitRouter struct {
	connection *amqp.Connection
	wg         sync.WaitGroup
	handlers   map[mqconnector.Queue]rabbitRoute

	consumers []mqconnector.Consumer
}

func NewRabbitRouter(conn *amqp.Connection) *RabbitRouter {
	return &RabbitRouter{
		handlers:   make(map[mqconnector.Queue]rabbitRoute),
		connection: conn,
	}
}

// Launch all handlers with consumers running in goroutines
func (r *RabbitRouter) Run(ctx context.Context) {
	defer r.Close()

	for queue, route := range r.handlers {
		consumer, err := r.newConsumer(queue)
		if err != nil {
			panic(err)
		}
		r.consumers = append(r.consumers, consumer)

		r.wg.Add(1)
		go r.runConsumer(ctx, consumer, route.handler, route.contexts)
	}

	r.wg.Wait()
}

// create new consumer depending on queue type
func (r *RabbitRouter) newConsumer(queue mqconnector.Queue) (mqconnector.Consumer, error) {
	if queue.Name != "" {
		return NewRabbitConsumer(r.connection, queue.Name)
	}
	if queue.Subscribtion != "" {
		return NewRabbitConsumer(r.connection, queue.Subscribtion)
	}
	return nil, errors.New("queue type is not supported")
}

// launch consumer with handler and contexters
func (r *RabbitRouter) runConsumer(
	ctx context.Context,
	consumer mqconnector.Consumer,
	handler mqconnector.Handler,
	contexts []mqconnector.ContextFunc,
) {
	defer r.wg.Done()
	consumer.
		SetHandler(handler).
		SetContexters(contexts...).
		Run(ctx)
}

// close all consumers
func (r *RabbitRouter) Close() error {
	for _, consumer := range r.consumers {
		if err := consumer.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (r *RabbitRouter) Route(
	queue mqconnector.Queue,
	handler mqconnector.Handler,
	contexters ...mqconnector.ContextFunc,
) {

	r.handlers[queue] = rabbitRoute{
		handler:  handler,
		contexts: contexters,
	}
}
