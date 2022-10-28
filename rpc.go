package mqconnector

import (
	"context"
	"encoding/json"
)

// A genereal interface for realization router for request-response
// model from message broker.
type RpcServer interface {
	Runner
	Closable

	// Returns an interface with which "client side" could
	// communicate with the router safely
	Client() RpcClient
}

// An interface with which a thread/goroutine can safely communicate
// with router to send requests into other service.
type RpcClient interface {
	Closable

	// Send delivery through router to another service and block current execution
	// until the returning og the response
	Send(ctx context.Context, request json.RawMessage) (response json.RawMessage, err error)
}
