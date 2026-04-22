package bus

import (
	"context"
	"time"
)

// Queueable embeds queue routing options into a command/job struct.
type Queueable struct {
	Connection          string
	Queue               string
	Delay               time.Duration
	ChainJobs           []any
	Middleware          []Pipe
	ChainConnection     string
	ChainQueue          string
	ChainCatchCallbacks []func(ctx context.Context, err error)
	AfterCommit         *bool
}

// OnConnection sets the queue connection.

// OnQueue sets the queue name.

// WithDelay sets the dispatch delay.

// WithoutDelay removes any dispatch delay.

// Chain sets a sequence of jobs to run after this one succeeds.

// AppendToChain appends jobs to the existing chain.

// PrependToChain inserts jobs at the beginning of the chain.

// Through sets the middleware pipeline for this job.

// GetQueue returns the queue name.

// GetConnection returns the connection name.

// GetDelay returns the queue delay.

// AllOnConnection sets the connection on this job and all chain jobs
// that support the OnConnection method.

// AllOnQueue sets the queue on this job and all chain jobs
// that support the OnQueue method.

// Batchable embeds batch membership information into a job struct.
type Batchable struct {
	BatchID   string
	batchInst *Batch
}

func (q *Queueable) OnConnection(connection string) *Queueable {
	q.Connection = connection

	return q
}

func (q *Queueable) OnQueue(queue string) *Queueable {
	q.Queue = queue

	return q
}

func (q *Queueable) WithDelay(d time.Duration) *Queueable {
	q.Delay = d

	return q
}

func (q *Queueable) WithoutDelay() *Queueable {
	q.Delay = 0

	return q
}

func (q *Queueable) Chain(jobs ...any) *Queueable {
	q.ChainJobs = jobs

	return q
}

func (q *Queueable) AppendToChain(jobs ...any) *Queueable {
	q.ChainJobs = append(q.ChainJobs, jobs...)

	return q
}

func (q *Queueable) PrependToChain(jobs ...any) *Queueable {
	q.ChainJobs = append(jobs, q.ChainJobs...)

	return q
}

func (q *Queueable) Through(pipes ...Pipe) *Queueable {
	q.Middleware = pipes

	return q
}

func (q *Queueable) GetQueue() string { return q.Queue }

func (q *Queueable) GetConnection() string { return q.Connection }

func (q *Queueable) GetDelay() time.Duration { return q.Delay }

func (q *Queueable) AllOnConnection(connection string) *Queueable {
	q.Connection = connection
	q.ChainConnection = connection

	for _, job := range q.ChainJobs {
		if c, ok := job.(interface{ OnConnection(string) *Queueable }); ok {
			c.OnConnection(connection)
		}
	}

	return q
}

func (q *Queueable) AllOnQueue(queue string) *Queueable {
	q.Queue = queue
	q.ChainQueue = queue

	for _, job := range q.ChainJobs {
		if c, ok := job.(interface{ OnQueue(string) *Queueable }); ok {
			c.OnQueue(queue)
		}
	}

	return q
}

// SetAfterCommit marks the job to be dispatched after the database transaction commits.
func (q *Queueable) SetAfterCommit() *Queueable {
	v := true
	q.AfterCommit = &v

	return q
}

// SetBeforeCommit marks the job to be dispatched before the database transaction commits.
func (q *Queueable) SetBeforeCommit() *Queueable {
	v := false
	q.AfterCommit = &v

	return q
}

// OnChainCatch registers a callback invoked when a chained job fails.
func (q *Queueable) OnChainCatch(fn func(ctx context.Context, err error)) *Queueable {
	q.ChainCatchCallbacks = append(q.ChainCatchCallbacks, fn)

	return q
}

// DispatchNextJobInChain dispatches the next job in the chain, if any.
func (q *Queueable) DispatchNextJobInChain(ctx context.Context, dispatcher Dispatcher) error {
	if len(q.ChainJobs) == 0 {
		return nil
	}

	next := q.ChainJobs[0]
	remaining := q.ChainJobs[1:]

	// Pass remaining chain to the next job.
	if chainable, ok := next.(interface{ Chain(jobs ...any) *Queueable }); ok {
		chainable.Chain(remaining...)
	}

	// Pass chain connection/queue to the next job.
	if q.ChainConnection != "" {
		if c, ok := next.(interface{ OnConnection(string) *Queueable }); ok {
			c.OnConnection(q.ChainConnection)
		}
	}

	if q.ChainQueue != "" {
		if c, ok := next.(interface{ OnQueue(string) *Queueable }); ok {
			c.OnQueue(q.ChainQueue)
		}
	}

	// Pass chain catch callbacks to the next job.
	if len(q.ChainCatchCallbacks) > 0 {
		if c, ok := next.(interface {
			OnChainCatch(func(context.Context, error)) *Queueable
		}); ok {
			for _, fn := range q.ChainCatchCallbacks {
				c.OnChainCatch(fn)
			}
		}
	}

	_, err := dispatcher.Dispatch(ctx, next)

	return err
}

// InvokeChainCatchCallbacks invokes all chain catch callbacks.
func (q *Queueable) InvokeChainCatchCallbacks(ctx context.Context, err error) {
	for _, fn := range q.ChainCatchCallbacks {
		fn(ctx, err)
	}
}

// Batching reports whether the job is part of an active (non-cancelled) batch.
func (b *Batchable) Batching() bool {
	if b.BatchID == "" {
		return false
	}

	if b.batchInst != nil && b.batchInst.Cancelled() {
		return false
	}

	return true
}

// WithBatchID sets the batch ID.
func (b *Batchable) WithBatchID(id string) { b.BatchID = id }

// Batch returns the parent Batch instance, if set.
func (b *Batchable) Batch() *Batch { return b.batchInst }

// SetBatch sets the parent Batch instance.
func (b *Batchable) SetBatch(batch *Batch) { b.batchInst = batch }

// WithFakeBatch creates a fake in-memory Batch for testing.
// It sets both the BatchID and the local batch instance.
func (b *Batchable) WithFakeBatch(id, name string, totalJobs int) *Batch {
	batch := &Batch{
		ID:        id,
		Name:      name,
		TotalJobs: totalJobs,
		Options:   make(map[string]any),
	}

	b.BatchID = id
	b.batchInst = batch

	return batch
}

// BatchFromRepo retrieves the Batch from the repository if the local instance is nil.
func (b *Batchable) BatchFromRepo(ctx context.Context, repo BatchRepository) (*Batch, error) {
	if b.batchInst != nil {
		return b.batchInst, nil
	}

	if b.BatchID == "" || repo == nil {
		return nil, nil
	}

	batch, err := repo.Get(ctx, b.BatchID)

	if err != nil {
		return nil, err
	}

	b.batchInst = batch

	return batch, nil
}
