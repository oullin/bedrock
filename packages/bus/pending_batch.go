package bus

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"
)

// PendingBatch is a fluent builder for creating and dispatching a Batch.
type PendingBatch struct {
	name              string
	jobs              []any
	connection        string
	queue             string
	allowFailures     bool
	progressCallbacks []func(ctx context.Context, batch *Batch)
	thenCallbacks     []func(ctx context.Context, batch *Batch)
	catchCallbacks    []func(ctx context.Context, batch *Batch, err error)
	finallyCallbacks  []func(ctx context.Context, batch *Batch)
	beforeCallbacks   []func(ctx context.Context, batch *Batch)
	dispatcher        QueueingDispatcher
}

// NewPendingBatch creates a PendingBatch.
func NewPendingBatch(dispatcher QueueingDispatcher, jobs []any) *PendingBatch {
	return &PendingBatch{dispatcher: dispatcher, jobs: jobs}
}

// Name sets the batch name.
func (p *PendingBatch) Name(name string) *PendingBatch {
	p.name = name

	return p
}

// Add appends more jobs to the batch.
func (p *PendingBatch) Add(jobs ...any) *PendingBatch {
	p.jobs = append(p.jobs, jobs...)

	return p
}

// OnConnection sets the queue connection.
func (p *PendingBatch) OnConnection(connection string) *PendingBatch {
	p.connection = connection

	return p
}

// OnQueue sets the queue name.
func (p *PendingBatch) OnQueue(queue string) *PendingBatch {
	p.queue = queue

	return p
}

// Before registers a callback invoked before any job is dispatched.
func (p *PendingBatch) Before(fn func(ctx context.Context, batch *Batch)) *PendingBatch {
	p.beforeCallbacks = append(p.beforeCallbacks, fn)

	return p
}

// Progress registers a callback invoked after each job completes.
func (p *PendingBatch) Progress(fn func(ctx context.Context, batch *Batch)) *PendingBatch {
	p.progressCallbacks = append(p.progressCallbacks, fn)

	return p
}

// Then registers a callback invoked when all jobs succeed.
func (p *PendingBatch) Then(fn func(ctx context.Context, batch *Batch)) *PendingBatch {
	p.thenCallbacks = append(p.thenCallbacks, fn)

	return p
}

// Catch registers a callback invoked when any job fails.
func (p *PendingBatch) Catch(fn func(ctx context.Context, batch *Batch, err error)) *PendingBatch {
	p.catchCallbacks = append(p.catchCallbacks, fn)

	return p
}

// Finally registers a callback invoked when all jobs finish (success or failure).
func (p *PendingBatch) Finally(fn func(ctx context.Context, batch *Batch)) *PendingBatch {
	p.finallyCallbacks = append(p.finallyCallbacks, fn)

	return p
}

// AllowFailures prevents the batch from being considered failed when a job fails.
func (p *PendingBatch) AllowFailures() *PendingBatch {
	p.allowFailures = true

	return p
}

// Dispatch creates and persists the batch, then dispatches all jobs.
func (p *PendingBatch) Dispatch(ctx context.Context) (*Batch, error) {
	id, err := generateBatchID()
	if err != nil {
		return nil, err
	}

	batch := &Batch{
		ID:                id,
		Name:              p.name,
		TotalJobs:         len(p.jobs),
		PendingJobs:       len(p.jobs),
		Options:           make(map[string]any),
		CreatedAt:         time.Now(),
		ProgressCallbacks: p.progressCallbacks,
		ThenCallbacks:     p.thenCallbacks,
		CatchCallbacks:    p.catchCallbacks,
		FinallyCallbacks:  p.finallyCallbacks,
	}

	if p.allowFailures {
		batch.Options["allowFailures"] = true
	}

	// Invoke before callbacks.
	for _, fn := range p.beforeCallbacks {
		fn(ctx, batch)
	}

	// Dispatch each job.
	for _, job := range p.jobs {
		if err = p.dispatcher.DispatchToQueue(ctx, job); err != nil {
			return batch, err
		}
	}

	return batch, nil
}

func generateBatchID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
