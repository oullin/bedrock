package drivers

import (
	"context"
	"os/exec"
	"time"

	"github.com/bedrock/packages/queue"
)

// BackgroundDriver dispatches jobs by spawning OS subprocesses.
// The cli command path must be configured; each job is executed as:
//
//	php cli queue:work --once --connection=<conn> --queue=<queue>
//
// In Go projects, replace the command with your own job runner binary.
type BackgroundDriver struct {
	command    string // binary to run for each job
	args       []string
	connection string
	inner      queue.Queue // underlying queue (e.g. database) for storage
}

// NewBackgroundDriver creates a BackgroundDriver backed by inner for storage.
// command is the binary invoked as a subprocess for each job.
func NewBackgroundDriver(command string, args []string, inner queue.Queue, connection string) *BackgroundDriver {
	return &BackgroundDriver{
		command:    command,
		args:       args,
		connection: connection,
		inner:      inner,
	}
}

func (d *BackgroundDriver) Push(ctx context.Context, queueName string, payload []byte) (string, error) {
	id, err := d.inner.Push(ctx, queueName, payload)

	if err != nil {
		return "", err
	}

	go d.spawn()

	return id, nil
}

func (d *BackgroundDriver) PushDelayed(ctx context.Context, queueName string, payload []byte, delay time.Duration) (string, error) {
	id, err := d.inner.PushDelayed(ctx, queueName, payload, delay)

	if err != nil {
		return "", err
	}

	time.AfterFunc(delay, d.spawn)

	return id, nil
}

func (d *BackgroundDriver) PushMultiple(ctx context.Context, queueName string, payloads [][]byte) ([]string, error) {
	return d.inner.PushMultiple(ctx, queueName, payloads)
}

func (d *BackgroundDriver) Pop(ctx context.Context, queueName string) (queue.Job, error) {
	return d.inner.Pop(ctx, queueName)
}

func (d *BackgroundDriver) Size(ctx context.Context, queueName string) (int64, error) {
	return d.inner.Size(ctx, queueName)
}

func (d *BackgroundDriver) PendingSize(ctx context.Context, queueName string) (int64, error) {
	return d.inner.PendingSize(ctx, queueName)
}

func (d *BackgroundDriver) DelayedSize(ctx context.Context, queueName string) (int64, error) {
	return d.inner.DelayedSize(ctx, queueName)
}

func (d *BackgroundDriver) ReservedSize(ctx context.Context, queueName string) (int64, error) {
	return d.inner.ReservedSize(ctx, queueName)
}

func (d *BackgroundDriver) ConnectionName() string { return d.connection }

// QueueNames delegates to the wrapped storage queue when it implements
// QueueNamer; otherwise it surfaces ErrNotSupported so the caller can
// distinguish "no queues" from "this backend can't enumerate".
func (d *BackgroundDriver) QueueNames(ctx context.Context) ([]string, error) {
	if namer, ok := d.inner.(queue.QueueNamer); ok {
		return namer.QueueNames(ctx)
	}

	return nil, queue.ErrNotSupported
}

// PendingJobs delegates to the wrapped storage queue.
func (d *BackgroundDriver) PendingJobs(ctx context.Context, queueName string) ([]queue.InspectedJob, error) {
	if insp, ok := d.inner.(queue.JobInspector); ok {
		return insp.PendingJobs(ctx, queueName)
	}

	return nil, queue.ErrNotSupported
}

// DelayedJobs delegates to the wrapped storage queue.
func (d *BackgroundDriver) DelayedJobs(ctx context.Context, queueName string) ([]queue.InspectedJob, error) {
	if insp, ok := d.inner.(queue.JobInspector); ok {
		return insp.DelayedJobs(ctx, queueName)
	}

	return nil, queue.ErrNotSupported
}

// ReservedJobs delegates to the wrapped storage queue.
func (d *BackgroundDriver) ReservedJobs(ctx context.Context, queueName string) ([]queue.InspectedJob, error) {
	if insp, ok := d.inner.(queue.JobInspector); ok {
		return insp.ReservedJobs(ctx, queueName)
	}

	return nil, queue.ErrNotSupported
}

func (d *BackgroundDriver) spawn() {
	cmd := exec.Command(d.command, d.args...) //nolint:gosec
	_ = cmd.Start()
}
