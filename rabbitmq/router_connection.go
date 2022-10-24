package rabbitmq

import (
	"context"
	"encoding/json"

	mqc "github.com/Velnbur/mq-connector"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Check on compile time, that RouterConnection implements
// mqc.RouterConnection
var _ mqc.RouterConnection = &RouterConnection{}

type RouterConnection struct {
	id uuid.UUID

	requests  chan<- message
	responses <-chan message
}

func (c *RouterConnection) Send(ctx context.Context, data json.RawMessage) (json.RawMessage, error) {
	c.requests <- message{
		CorrelationID: c.id.String(),
		Data:          data,
		Type:          messageSend,
	}

	select {
	case <-ctx.Done():
		return nil, errors.New("context closed")
	case msg := <-c.responses:
		return msg.Data, nil
	}
}

func (c *RouterConnection) Close() error {
	c.requests <- message{
		CorrelationID: c.id.String(),
		Type:          messageClose,
	}

	return nil
}
