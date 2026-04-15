package bus

import (
	"errors"
	"fmt"

	"github.com/bedrock/packages/container"
	"github.com/bedrock/packages/queue"
)

// BusServiceProvider registers the command bus dispatcher into the container.
// It mirrors Illuminate\Bus\BusServiceProvider.
//
// At resolve time the provider tries to wire a queue backend by looking up
// the "queue.connection" abstract in the container — this is intentionally
// distinct from "queue" (the QueueManager). If "queue.connection" is
// unbound, the bus is constructed without a queue backend (DispatchToQueue
// will then panic). If it IS bound but resolves to a value that is not a
// queue.Queue, that is a wiring error and the factory returns an error
// rather than silently ignoring it.
//
// Wiring a default queue connection is the application's responsibility:
//
//	app.Singleton("queue.connection", func(c *container.Container) (any, error) {
//	    mgr, _ := c.Make("queue")
//	    return mgr.(*queue.Manager).Connection("sync")
//	})
type BusServiceProvider struct {
	app *container.Container
}

// NewBusServiceProvider constructs the provider.
func NewBusServiceProvider(app *container.Container) *BusServiceProvider {
	return &BusServiceProvider{app: app}
}

// Register binds the bus dispatcher as a singleton under "bus".
func (p *BusServiceProvider) Register() {
	p.app.Singleton("bus", func(c *container.Container) (any, error) {
		queueBackend, err := optionalQueue(c)

		if err != nil {
			return nil, err
		}

		return NewDispatcher(queueBackend, nil), nil
	})
}

// optionalQueue resolves "queue.connection" if bound. Returns (nil, nil)
// when unbound, (q, nil) when bound to a queue.Queue, and a non-nil error
// when bound to the wrong type.
func optionalQueue(c *container.Container) (queue.Queue, error) {
	raw, err := c.Make("queue.connection")

	if err != nil {
		if errors.Is(err, container.ErrNotBound) {
			return nil, nil
		}

		return nil, fmt.Errorf("bus: resolving queue.connection: %w", err)
	}

	q, ok := raw.(queue.Queue)

	if !ok {
		return nil, fmt.Errorf("bus: \"queue.connection\" binding has wrong type %T (want queue.Queue)", raw)
	}

	return q, nil
}

// Provides returns the abstract keys registered by this provider.
func (p *BusServiceProvider) Provides() []string {
	return []string{"bus"}
}
