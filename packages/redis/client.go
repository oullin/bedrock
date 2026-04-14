package redis

import (
	"context"
	"time"
)

// Client is the minimal Redis command surface the Connection depends on.
// Both the go-redis adapter and the in-memory fake used in unit tests
// implement it. Upstream's PhpRedisConnection / PredisConnection methods
// all ultimately map onto these primitives.
//
// Do is the generic command dispatcher (parity with Framework's
// Connection::command). Typed helpers delegate to Do so one implementation
// is enough to satisfy the interface, and specialized backends can
// optionally override hot paths.
type Client interface {
	// Do dispatches a raw Redis command. The first argument is the command
	// name, remaining arguments are positional parameters. The return value
	// is the decoded reply: string, int64, []any, or nil.
	Do(ctx context.Context, args ...any) (any, error)

	// Pipeline returns a new Pipeliner that batches commands.
	Pipeline() Pipeliner

	// TxPipeline returns a new Pipeliner wrapped in MULTI/EXEC.
	TxPipeline() Pipeliner

	// Subscribe starts a subscription on the given channels.
	Subscribe(ctx context.Context, channels ...string) Subscription

	// PSubscribe starts a pattern subscription.
	PSubscribe(ctx context.Context, patterns ...string) Subscription

	// Close releases resources held by the client.
	Close() error
}

// Pipeliner queues commands for batch execution.
type Pipeliner interface {
	Do(ctx context.Context, args ...any) Cmder
	Exec(ctx context.Context) ([]Cmder, error)
	Discard()
	Len() int
}

// Cmder is the result of a single pipelined command.
type Cmder interface {
	Name() string
	Args() []any
	Result() (any, error)
	Err() error
}

// Subscription is a pub/sub handle.
type Subscription interface {
	// Channel returns a receive channel for incoming messages. The
	// caller is responsible for terminating by calling Close.
	Channel() <-chan Message
	Close() error
}

// Message is a published payload.
type Message struct {
	Channel string
	Pattern string // empty for non-pattern subscriptions
	Payload string
}

// TimeoutError reports a context/command timeout.
type TimeoutError struct {
	Op      string
	Elapsed time.Duration
}

func (e *TimeoutError) Error() string {
	return "redis: " + e.Op + " timed out after " + e.Elapsed.String()
}
