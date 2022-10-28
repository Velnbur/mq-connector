package mqconnector

import (
	"context"
	"encoding/json"
)

//go:generate mockery --name=Producer
type Producer interface {
	Publish(ctx context.Context, data json.RawMessage) error
}
