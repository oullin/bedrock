package drivers

import (
	"context"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/events"
)

// FailoverDriver tries each driver in order; reads from the first that succeeds,
// writes to the first that accepts the push. On a fallthrough the driver emits
// a QueueFailedOver event through the configured EventEmitter (if any) so
// operators can alert on backend degradation.
type FailoverDriver struct {
	drivers    []queue.Queue
	connection string
	emitter    queue.EventEmitter
}

// NewFailoverDriver creates a FailoverDriver. drivers are tried in order.
func NewFailoverDriver(connection string, drivers ...queue.Queue) *FailoverDriver {
	return &FailoverDriver{drivers: drivers, connection: connection}
}

// SetEmitter installs an EventEmitter the driver will use to dispatch
// QueueFailedOver events. Passing nil disables emission. Mirrors the
// constructor-injected Events\Dispatcher that Laravel's FailoverQueue
// receives.
func (d *FailoverDriver) SetEmitter(e queue.EventEmitter) *FailoverDriver {
	d.emitter = e

	return d
}

// emitFailover dispatches a QueueFailedOver event for a fallthrough
// from one driver to the next. It is a no-op when the emitter is nil.
func (d *FailoverDriver) emitFailover(from, to queue.Queue, err error) {
	if d.emitter == nil {
		return
	}

	var fromName, toName string

	if from != nil {
		fromName = from.ConnectionName()
	}

	if to != nil {
		toName = to.ConnectionName()
	}

	d.emitter.Emit(events.QueueFailedOver{From: fromName, To: toName, Err: err})
}

func (d *FailoverDriver) Push(ctx context.Context, queueName string, payload []byte) (string, error) {
	var lastErr error

	for i, drv := range d.drivers {
		id, err := drv.Push(ctx, queueName, payload)

		if err == nil {
			return id, nil
		}

		lastErr = err

		if i+1 < len(d.drivers) {
			d.emitFailover(drv, d.drivers[i+1], err)
		}
	}

	return "", lastErr
}

func (d *FailoverDriver) PushDelayed(ctx context.Context, queueName string, payload []byte, delay time.Duration) (string, error) {
	var lastErr error

	for i, drv := range d.drivers {
		id, err := drv.PushDelayed(ctx, queueName, payload, delay)

		if err == nil {
			return id, nil
		}

		lastErr = err

		if i+1 < len(d.drivers) {
			d.emitFailover(drv, d.drivers[i+1], err)
		}
	}

	return "", lastErr
}

func (d *FailoverDriver) PushMultiple(ctx context.Context, queueName string, payloads [][]byte) ([]string, error) {
	var lastErr error

	for i, drv := range d.drivers {
		ids, err := drv.PushMultiple(ctx, queueName, payloads)

		if err == nil {
			return ids, nil
		}

		lastErr = err

		if i+1 < len(d.drivers) {
			d.emitFailover(drv, d.drivers[i+1], err)
		}
	}

	return nil, lastErr
}

func (d *FailoverDriver) Pop(ctx context.Context, queueName string) (queue.Job, error) {
	for _, drv := range d.drivers {
		job, err := drv.Pop(ctx, queueName)

		if err == nil && job != nil {
			return job, nil
		}
	}

	return nil, queue.ErrNoJob
}

func (d *FailoverDriver) Size(ctx context.Context, queueName string) (int64, error) {
	for _, drv := range d.drivers {
		n, err := drv.Size(ctx, queueName)

		if err == nil {
			return n, nil
		}
	}

	return 0, nil
}

func (d *FailoverDriver) PendingSize(ctx context.Context, queueName string) (int64, error) {
	for _, drv := range d.drivers {
		n, err := drv.PendingSize(ctx, queueName)

		if err == nil {
			return n, nil
		}
	}

	return 0, nil
}

func (d *FailoverDriver) DelayedSize(ctx context.Context, queueName string) (int64, error) {
	for _, drv := range d.drivers {
		n, err := drv.DelayedSize(ctx, queueName)

		if err == nil {
			return n, nil
		}
	}

	return 0, nil
}

func (d *FailoverDriver) ReservedSize(ctx context.Context, queueName string) (int64, error) {
	for _, drv := range d.drivers {
		n, err := drv.ReservedSize(ctx, queueName)

		if err == nil {
			return n, nil
		}
	}

	return 0, nil
}

func (d *FailoverDriver) ConnectionName() string { return d.connection }
