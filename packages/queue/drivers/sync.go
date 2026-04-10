package drivers

import (
	"context"
	"time"

	"github.com/bedrock/packages/queue"
)

// SyncDriver executes jobs immediately in the same goroutine. Useful for testing.
type SyncDriver struct {
	connection string
	handler    queue.Handler
}

// NewSyncDriver creates a SyncDriver. handler is called synchronously for every Push.

type syncJob struct{ BaseJob }

func NewSyncDriver(connection string, handler queue.Handler) *SyncDriver {
	return &SyncDriver{connection: connection, handler: handler}
}

func (d *SyncDriver) Push(ctx context.Context, queueName string, payload []byte) (string, error) {
	job := &syncJob{BaseJob: BaseJob{payload: payload, queue: queueName, connection: d.connection}}
	job.fireFunc = func(ctx context.Context) error {
		return d.handler.Handle(ctx, job)
	}

	return "", job.Fire(ctx)
}

func (d *SyncDriver) PushDelayed(ctx context.Context, queueName string, payload []byte, delay time.Duration) (string, error) {
	time.Sleep(delay)

	return d.Push(ctx, queueName, payload)
}

func (d *SyncDriver) PushMultiple(ctx context.Context, queueName string, payloads [][]byte) ([]string, error) {
	ids := make([]string, 0, len(payloads))

	for _, p := range payloads {
		id, err := d.Push(ctx, queueName, p)

		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (d *SyncDriver) Pop(_ context.Context, _ string) (queue.Job, error) { return nil, queue.ErrNoJob }

func (d *SyncDriver) Size(_ context.Context, _ string) (int64, error)         { return 0, nil }
func (d *SyncDriver) PendingSize(_ context.Context, _ string) (int64, error)  { return 0, nil }
func (d *SyncDriver) DelayedSize(_ context.Context, _ string) (int64, error)  { return 0, nil }
func (d *SyncDriver) ReservedSize(_ context.Context, _ string) (int64, error) { return 0, nil }
func (d *SyncDriver) ConnectionName() string                                  { return d.connection }
