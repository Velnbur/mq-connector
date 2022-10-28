package mqconnector

import (
	"context"
	"encoding/json"
)

type Handler func(ctx context.Context, data json.RawMessage) error

type Consumer interface {
	Runner
	Closable

	SetContexters(ctxers ...Contexter) Consumer
	SetHandler(handler Handler) Consumer
}
