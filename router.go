package mqconnector

import (
	"context"
)

// A struct that describes a queue params for consumer/produer
type Queue struct {
	// A name of the queue for manual common consumers.
	Name string
	// A name of the exhange for Subscribers.
	Subscribtion string
}

// ContextFunc - a type for functions that will include something
// in handler's context before the handler will be called.
type ContextFunc func(parent context.Context) (child context.Context)

// Router - interface for consumer router that will manage and create consumer
// for each queue.
//
//go:generate mockery --name=Router
type Router interface {
	Runner

	// Add a handler for a specific queue. `contexters` - a list of
	// functions that will be called on each delivery from the queue
	// to include something in the context before the handler.
	Route(queue Queue, handler Handler, contexters ...ContextFunc)
}
