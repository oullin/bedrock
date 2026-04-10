package drivers

import (
	"context"
	"time"

	"github.com/bedrock/packages/queue"
)

// BeanstalkdClient is the interface for a Beanstalkd client.
type BeanstalkdClient interface {
	// Put inserts a job with priority, delay, and TTR (time-to-run). Returns job ID.
	Put(ctx context.Context, tube string, body []byte, priority uint32, delay, ttr time.Duration) (uint64, error)
	// ReserveWithTimeout reserves the next job from the given tube within the timeout.
	ReserveWithTimeout(ctx context.Context, tube string, timeout time.Duration) (id uint64, body []byte, err error)
	// Delete removes a job by ID.
	Delete(ctx context.Context, id uint64) error
	// Release puts a job back with a new delay.
	Release(ctx context.Context, id uint64, priority uint32, delay time.Duration) error
	// Bury buries a job (marks as failed) by ID.
	Bury(ctx context.Context, id uint64, priority uint32) error
	// StatsTube returns tube statistics as a map.
	StatsTube(ctx context.Context, tube string) (map[string]string, error)
}

// BeanstalkdDriver enqueues jobs via a Beanstalkd client.
type BeanstalkdDriver struct {
	client     BeanstalkdClient
	connection string
	ttr        time.Duration // time-to-run per job
}

// NewBeanstalkdDriver creates a BeanstalkdDriver. ttr is the job TTR.

type bsJob struct{ BaseJob }

func NewBeanstalkdDriver(client BeanstalkdClient, connection string, ttr time.Duration) *BeanstalkdDriver {
	if ttr == 0 {
		ttr = 60 * time.Second
	}

	return &BeanstalkdDriver{client: client, connection: connection, ttr: ttr}
}

func (d *BeanstalkdDriver) Push(ctx context.Context, queueName string, payload []byte) (string, error) {
	id, err := d.client.Put(ctx, queueName, payload, 1024, 0, d.ttr)

	return toString(id), err
}

func (d *BeanstalkdDriver) PushDelayed(ctx context.Context, queueName string, payload []byte, delay time.Duration) (string, error) {
	id, err := d.client.Put(ctx, queueName, payload, 1024, delay, d.ttr)

	return toString(id), err
}

func (d *BeanstalkdDriver) PushMultiple(ctx context.Context, queueName string, payloads [][]byte) ([]string, error) {
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

func (d *BeanstalkdDriver) Pop(ctx context.Context, queueName string) (queue.Job, error) {
	id, body, err := d.client.ReserveWithTimeout(ctx, queueName, 5*time.Second)

	if err != nil {
		return nil, queue.ErrNoJob
	}

	job := &bsJob{
		BaseJob: BaseJob{
			id:         toString(id),
			payload:    body,
			queue:      queueName,
			connection: d.connection,
		},
	}
	job.deleteFunc = func() error {
		return d.client.Delete(ctx, id)
	}

	job.releaseFunc = func(delay time.Duration) error {
		return d.client.Release(ctx, id, 1024, delay)
	}

	job.failFunc = func(_ error) error {
		return d.client.Bury(ctx, id, 1024)
	}

	return job, nil
}

func (d *BeanstalkdDriver) Size(ctx context.Context, queueName string) (int64, error) {
	stats, err := d.client.StatsTube(ctx, queueName)

	if err != nil {
		return 0, err
	}

	return parseStatInt(stats, "current-jobs-ready"), nil
}

func (d *BeanstalkdDriver) PendingSize(ctx context.Context, queueName string) (int64, error) {
	return d.Size(ctx, queueName)
}

func (d *BeanstalkdDriver) DelayedSize(ctx context.Context, queueName string) (int64, error) {
	stats, err := d.client.StatsTube(ctx, queueName)

	if err != nil {
		return 0, err
	}

	return parseStatInt(stats, "current-jobs-delayed"), nil
}

func (d *BeanstalkdDriver) ReservedSize(ctx context.Context, queueName string) (int64, error) {
	stats, err := d.client.StatsTube(ctx, queueName)

	if err != nil {
		return 0, err
	}

	return parseStatInt(stats, "current-jobs-reserved"), nil
}

func (d *BeanstalkdDriver) ConnectionName() string { return d.connection }
