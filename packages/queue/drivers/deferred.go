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

// QueueNames returns the unique queue names currently buffered for
// deferred dispatch.
func (d *DeferredDriver) QueueNames(_ context.Context) ([]string, error) {
	d.mu.Lock()

	defer d.mu.Unlock()

	seen := make(map[string]struct{}, len(d.deferred))

	var out []string

	for _, e := range d.deferred {
		if _, ok := seen[e.Queue]; ok {
			continue
		}

		seen[e.Queue] = struct{}{}

		out = append(out, e.Queue)
	}

	return out, nil
}

// PendingJobs returns snapshots of the entries due immediately (After
// is in the past).
func (d *DeferredDriver) PendingJobs(_ context.Context, queueName string) ([]queue.InspectedJob, error) {
	return d.snapshots(queueName, func(e DeferredEntry) bool { return !e.After.After(time.Now()) }), nil
}

// DelayedJobs returns snapshots of entries whose After is still in
// the future.
func (d *DeferredDriver) DelayedJobs(_ context.Context, queueName string) ([]queue.InspectedJob, error) {
	return d.snapshots(queueName, func(e DeferredEntry) bool { return e.After.After(time.Now()) }), nil
}

// ReservedJobs returns an empty slice. The deferred driver never
// reserves jobs — every entry is dispatched in Flush.
func (d *DeferredDriver) ReservedJobs(_ context.Context, _ string) ([]queue.InspectedJob, error) {
	return nil, nil
}

func (d *DeferredDriver) snapshots(queueName string, include func(DeferredEntry) bool) []queue.InspectedJob {
	d.mu.Lock()

	defer d.mu.Unlock()

	var out []queue.InspectedJob

	for _, e := range d.deferred {
		if e.Queue != queueName || !include(e) {
			continue
		}

		out = append(out, queue.InspectedJob{
			Queue:       e.Queue,
			Connection:  d.connection,
			Payload:     e.Payload,
			AvailableAt: e.After,
		})
	}

	return out
}

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
