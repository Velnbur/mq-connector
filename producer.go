package mqconnector

import (
	"context"
	"encoding/json"
)

type Producer interface {
	Publish(ctx context.Context, data json.RawMessage) error
}
