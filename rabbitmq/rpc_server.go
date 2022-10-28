package rabbitmq

import (
	"context"
	"log"
	"sync"

	mqc "github.com/Velnbur/mq-connector"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Check on compile time, that RabbitRouter implements
// mqc.Router
var _ mqc.RpcServer = &RabbitRpcServer{}

type RabbitRpcServer struct {
	channel       *amqp.Channel
	queueRequest  amqp.Queue
	queueResponse amqp.Queue

	requests chan message

	receivers *responseReceivers
}

func (r *RabbitRpcServer) Close() error {
	if err := r.channel.Close(); err != nil {
		return errors.Wrap(err, "failed to close channel")
	}
	return nil
}

const requestsMaxLen = 1024

func NewRouter(conn *amqp.Connection, reqQueue, respQueue string) (*RabbitRpcServer, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create channel")
	}

	queueReq, err := channel.QueueDeclare(
		reqQueue, // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create queue")
	}

	queueResp, err := channel.QueueDeclare(
		respQueue, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create queue")
	}

	return &RabbitRpcServer{
		channel:       channel,
		queueRequest:  queueReq,
		queueResponse: queueResp,
		requests:      make(chan message, requestsMaxLen),
		receivers:     newReceivers(),
	}, nil
}

const responseReceiverMaxLen = 2

func (r *RabbitRpcServer) Client() mqc.RpcClient {
	uuid := uuid.New()

	responseReceiver := make(chan message, responseReceiverMaxLen)

	r.receivers.Add(uuid, responseReceiver)

	return &RabbitRpcClient{
		id:        uuid,
		requests:  r.requests,
		responses: responseReceiver,
	}
}

func (r *RabbitRpcServer) Run(ctx context.Context) {
	deliveries, err := r.channel.Consume(
		r.queueResponse.Name, "", true, false, false, false, nil)
	if err != nil {
		// TODO: may be use antither router or handle it in another way
		log.Fatalf("failed to consume msg: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("context closed, finishing...")
			return
		case msg := <-r.requests:
			err = r.handleMessage(ctx, msg)
		case delivery := <-deliveries:
			err = r.handleResponse(delivery)
		}
		if err != nil {
			// TODO: may be use antither router or handle it in another way
			log.Fatalf("failed to handle delivery or msg: %s", err)
		}
	}
}

func (r *RabbitRpcServer) handleMessage(ctx context.Context, msg message) error {
	var err error

	switch msg.Type {
	case messageClose:
		id := uuid.MustParse(msg.CorrelationID)
		r.receivers.Close(id)
	case messageSend:
		err = r.channel.PublishWithContext(
			ctx,
			"",
			r.queueRequest.Name,
			false,
			false,
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: msg.CorrelationID,
				ReplyTo:       r.queueResponse.Name,
				Body:          msg.Data,
			})

		if err != nil {
			return errors.Wrap(err, "failed to publish delivery")
		}
	}

	return err
}

func (r *RabbitRpcServer) handleResponse(delivery amqp.Delivery) error {
	id := uuid.MustParse(delivery.CorrelationId)

	recv := r.receivers.Get(id)

	recv <- message{
		Data: delivery.Body,
	}

	return nil
}

// responseReceivers - is just a thread read/write save wrapper around map
// of Receiver's
type responseReceivers struct {
	mapMux *sync.RWMutex

	routes map[uuid.UUID]chan<- message
}

func newReceivers() *responseReceivers {
	return &responseReceivers{
		mapMux: new(sync.RWMutex),
		routes: make(map[uuid.UUID]chan<- message),
	}
}

func (r *responseReceivers) Add(id uuid.UUID, recv chan<- message) {
	r.mapMux.Lock()
	r.routes[id] = recv
	r.mapMux.Unlock()
}

func (r *responseReceivers) Get(id uuid.UUID) chan<- message {
	r.mapMux.RLock()
	route := r.routes[id]
	r.mapMux.RUnlock()

	return route
}

func (r *responseReceivers) Close(id uuid.UUID) {
	r.mapMux.Lock()
	close(r.routes[id])
	delete(r.routes, id)
	r.mapMux.Unlock()
}
