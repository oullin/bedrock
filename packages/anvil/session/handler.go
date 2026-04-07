package session

import "context"

// Handler abstracts the session storage backend. It mirrors PHP's
// SessionHandlerInterface.
type Handler interface {
	Open(ctx context.Context, path string, name string) error
	Close(ctx context.Context) error
	Read(ctx context.Context, id string) (string, error)
	Write(ctx context.Context, id string, data string) error
	Destroy(ctx context.Context, id string) error
	GC(ctx context.Context, maxLifetime int) error
}
