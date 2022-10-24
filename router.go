package mqconnector

import (
	"context"
	"encoding/json"
)

// A genereal interface for realization router for request-response
// model from message broker.
type Router interface {
	Runner

	// Returns an interface with which "client side" could
	// communicate with the router safely
	Connection() RouterConnection

	Close() error
}

// An interface with which a thread/goroutine can safely communicate
// with router to send requests into other service.
type RouterConnection interface {
	// Send delivery through router to another service and block current execution
	// until the returning og the response
	Send(ctx context.Context, request json.RawMessage) (response json.RawMessage, err error)

	Close() error
}
