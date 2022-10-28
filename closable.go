package mqconnector

// Closable is an interface for something that can be closed.
type Closable interface {
	Close() error
}
