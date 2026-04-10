package drivers

import (
	"context"
	"sync"
	"time"

	"github.com/bedrock/packages/queue"
)

// DeferredEntry holds a deferred job payload.
type DeferredEntry struct {
	Queue   string
	Payload []byte
	After   time.Time
}

// DeferredDriver buffers jobs in memory and dispatches them after-response.
// Call Flush() to process all buffered jobs (e.g. in an after-response hook).
type DeferredDriver struct {
	mu         sync.Mutex
	deferred   []DeferredEntry
	connection string
	dispatcher func(ctx context.Context, entry DeferredEntry) error
}

// NewDeferredDriver creates a DeferredDriver.
// dispatcher is called for each entry during Flush.
func NewDeferredDriver(connection string, dispatcher func(ctx context.Context, entry DeferredEntry) error) *DeferredDriver {
	return &DeferredDriver{connection: connection, dispatcher: dispatcher}
}

func (d *DeferredDriver) Push(_ context.Context, queueName string, payload []byte) (string, error) {
	d.mu.Lock()

	defer d.mu.Unlock()

	d.deferred = append(d.deferred, DeferredEntry{Queue: queueName, Payload: payload, After: time.Now()})

	return "", nil
}

func (d *DeferredDriver) PushDelayed(_ context.Context, queueName string, payload []byte, delay time.Duration) (string, error) {
	d.mu.Lock()

	defer d.mu.Unlock()

	d.deferred = append(d.deferred, DeferredEntry{Queue: queueName, Payload: payload, After: time.Now().Add(delay)})

	return "", nil
}

func (d *DeferredDriver) PushMultiple(_ context.Context, queueName string, payloads [][]byte) ([]string, error) {
	d.mu.Lock()

	defer d.mu.Unlock()

	for _, p := range payloads {
		d.deferred = append(d.deferred, DeferredEntry{Queue: queueName, Payload: p, After: time.Now()})
	}

	return make([]string, len(payloads)), nil
}

func (d *DeferredDriver) Pop(_ context.Context, _ string) (queue.Job, error) {
	return nil, queue.ErrNoJob
}

func (d *DeferredDriver) Size(_ context.Context, _ string) (int64, error) {
	d.mu.Lock()

	defer d.mu.Unlock()

	return int64(len(d.deferred)), nil
}

func (d *DeferredDriver) PendingSize(ctx context.Context, q string) (int64, error) {
	return d.Size(ctx, q)
}

func (d *DeferredDriver) DelayedSize(_ context.Context, _ string) (int64, error)  { return 0, nil }
func (d *DeferredDriver) ReservedSize(_ context.Context, _ string) (int64, error) { return 0, nil }
func (d *DeferredDriver) ConnectionName() string                                  { return d.connection }

// Flush dispatches all buffered deferred jobs via the dispatcher.
func (d *DeferredDriver) Flush(ctx context.Context) error {
	d.mu.Lock()
	entries := d.deferred
	d.deferred = nil
	d.mu.Unlock()

	for _, entry := range entries {
		if d.dispatcher != nil {
			if err := d.dispatcher(ctx, entry); err != nil {
				return err
			}
		}
	}

	return nil
}
