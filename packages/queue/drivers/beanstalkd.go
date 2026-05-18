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

// BeanstalkdTubeLister is the optional capability that lets the
// driver report which tubes exist on the server (the "list-tubes"
// protocol command). Without it QueueNames returns ErrNotSupported.
type BeanstalkdTubeLister interface {
	ListTubes(ctx context.Context) ([]string, error)
}

// BeanstalkdPeeker is the optional capability that lets the driver
// snapshot jobs without reserving them. PeekReady mirrors
// "peek-ready"; PeekDelayed mirrors "peek-delayed". Beanstalkd does
// not expose a way to enumerate every reserved job, so ReservedJobs
// surfaces ErrNotSupported regardless of whether this interface is
// satisfied.
type BeanstalkdPeeker interface {
	PeekReady(ctx context.Context, tube string) (uint64, []byte, error)
	PeekDelayed(ctx context.Context, tube string) (uint64, []byte, error)
}

// BeanstalkdDriver enqueues jobs via a Beanstalkd client. It is the
// Go port of Upstream's Framework\Queue\BeanstalkdQueue.
//
// Two knobs are tunable after construction:
//
//   - SetDefaultTube configures the fallback tube used when a caller
//     passes an empty queue name (Upstream's $default constructor arg).
//   - SetBlockFor configures the reserve-with-timeout value passed to
//     the client on Pop (Upstream's $blockFor constructor arg). Default
//     is zero for non-blocking reserve.
type BeanstalkdDriver struct {
	client      BeanstalkdClient
	connection  string
	ttr         time.Duration // time-to-run per job
	blockFor    time.Duration // reserve-with-timeout on Pop
	defaultTube string
}

// NewBeanstalkdDriver creates a BeanstalkdDriver. ttr is the job TTR.

type bsJob struct{ BaseJob }

func NewBeanstalkdDriver(client BeanstalkdClient, connection string, ttr time.Duration) *BeanstalkdDriver {
	if ttr == 0 {
		ttr = 60 * time.Second
	}

	return &BeanstalkdDriver{client: client, connection: connection, ttr: ttr}
}

// SetBlockFor sets the reserve-with-timeout value used on Pop. A zero
// value means non-blocking reserve.
func (d *BeanstalkdDriver) SetBlockFor(blockFor time.Duration) *BeanstalkdDriver {
	d.blockFor = blockFor

	return d
}

// SetDefaultTube configures the tube used when a caller passes an empty
// queue name to Push, PushDelayed, or Pop. Mirrors Upstream's $default
// BeanstalkdQueue constructor argument.
func (d *BeanstalkdDriver) SetDefaultTube(tube string) *BeanstalkdDriver {
	d.defaultTube = tube

	return d
}

// resolveTube falls back to defaultTube when queueName is empty,
// matching Upstream's getQueue($queue ?: $this->default) behaviour.
func (d *BeanstalkdDriver) resolveTube(queueName string) string {
	if queueName == "" {
		return d.defaultTube
	}

	return queueName
}

func (d *BeanstalkdDriver) Push(ctx context.Context, queueName string, payload []byte) (string, error) {
	id, err := d.client.Put(ctx, d.resolveTube(queueName), payload, 1024, 0, d.ttr)

	return toString(id), err
}

func (d *BeanstalkdDriver) PushDelayed(ctx context.Context, queueName string, payload []byte, delay time.Duration) (string, error) {
	id, err := d.client.Put(ctx, d.resolveTube(queueName), payload, 1024, delay, d.ttr)

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
	tube := d.resolveTube(queueName)

	id, body, err := d.client.ReserveWithTimeout(ctx, tube, d.blockFor)

	if err != nil {
		return nil, queue.ErrNoJob
	}

	job := &bsJob{
		BaseJob: BaseJob{
			id:         toString(id),
			payload:    body,
			queue:      tube,
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

// QueueNames reports the tubes currently known to the Beanstalkd
// server when the client implements BeanstalkdTubeLister; otherwise
// it returns ErrNotSupported.
func (d *BeanstalkdDriver) QueueNames(ctx context.Context) ([]string, error) {
	lister, ok := d.client.(BeanstalkdTubeLister)

	if !ok {
		return nil, queue.ErrNotSupported
	}

	return lister.ListTubes(ctx)
}

// PendingJobs returns the next ready job in the tube (peek-ready)
// when the client implements BeanstalkdPeeker. Beanstalkd only
// exposes the head of the ready queue, so the slice carries at most
// one entry — operators that want a fuller view should reach for a
// dedicated admin tool.
func (d *BeanstalkdDriver) PendingJobs(ctx context.Context, queueName string) ([]queue.InspectedJob, error) {
	peeker, ok := d.client.(BeanstalkdPeeker)

	if !ok {
		return nil, queue.ErrNotSupported
	}

	id, body, err := peeker.PeekReady(ctx, d.resolveTube(queueName))

	if err != nil || id == 0 {
		return nil, err
	}

	return []queue.InspectedJob{{
		ID:         int64(id),
		Queue:      d.resolveTube(queueName),
		Connection: d.connection,
		Payload:    body,
	}}, nil
}

// DelayedJobs returns the next delayed job in the tube (peek-delayed).
func (d *BeanstalkdDriver) DelayedJobs(ctx context.Context, queueName string) ([]queue.InspectedJob, error) {
	peeker, ok := d.client.(BeanstalkdPeeker)

	if !ok {
		return nil, queue.ErrNotSupported
	}

	id, body, err := peeker.PeekDelayed(ctx, d.resolveTube(queueName))

	if err != nil || id == 0 {
		return nil, err
	}

	return []queue.InspectedJob{{
		ID:         int64(id),
		Queue:      d.resolveTube(queueName),
		Connection: d.connection,
		Payload:    body,
	}}, nil
}

// ReservedJobs returns ErrNotSupported. Beanstalkd does not expose a
// "peek-reserved" command; only the worker holding a reservation can
// see the job. Operators that need this view should rely on
// Beanstalkd's stats-tube (current-jobs-reserved) for counts instead.
func (d *BeanstalkdDriver) ReservedJobs(_ context.Context, _ string) ([]queue.InspectedJob, error) {
	return nil, queue.ErrNotSupported
}
