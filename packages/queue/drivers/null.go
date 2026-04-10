package drivers

import (
	"context"
	"time"

	"github.com/bedrock/packages/queue"
)

// NullDriver discards all jobs silently.
type NullDriver struct {
	connection string
}

// NewNullDriver creates a NullDriver.
func NewNullDriver(connection string) *NullDriver {
	return &NullDriver{connection: connection}
}

func (d *NullDriver) Push(_ context.Context, _ string, _ []byte) (string, error)  { return "", nil }
func (d *NullDriver) PushDelayed(_ context.Context, _ string, _ []byte, _ time.Duration) (string, error) {
	return "", nil
}
func (d *NullDriver) PushMultiple(_ context.Context, _ string, payloads [][]byte) ([]string, error) {
	return make([]string, len(payloads)), nil
}
func (d *NullDriver) Pop(_ context.Context, _ string) (queue.Job, error)            { return nil, queue.ErrNoJob }
func (d *NullDriver) Size(_ context.Context, _ string) (int64, error)               { return 0, nil }
func (d *NullDriver) PendingSize(_ context.Context, _ string) (int64, error)        { return 0, nil }
func (d *NullDriver) DelayedSize(_ context.Context, _ string) (int64, error)        { return 0, nil }
func (d *NullDriver) ReservedSize(_ context.Context, _ string) (int64, error)       { return 0, nil }
func (d *NullDriver) ConnectionName() string                                        { return d.connection }
