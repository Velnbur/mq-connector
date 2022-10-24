package mqconnector

import (
	"context"
	"encoding/json"
)

type Handler func(ctx context.Context, data json.RawMessage) error

type Consumer interface {
	Runner

	SetHandler(handler Handler) Consumer
}
