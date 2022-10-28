package mqconnector

import (
	"context"
	"encoding/json"
)

type Handler func(ctx context.Context, data json.RawMessage) error

//go:generate mockery --name=Consumer
type Consumer interface {
	Runner
	Closable

	SetContexters(ctxers ...ContextFunc) Consumer
	SetHandler(handler Handler) Consumer
}
