package bus

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"
)

// FailureCallback is invoked when a failure-tolerant batch records a failed job.
type FailureCallback = func(ctx context.Context, batch *Batch, err error)

// PendingBatch is a fluent builder for creating and dispatching a Batch.
type PendingBatch struct {
	name              string
	jobs              []any
	connection        string
	queue             string
	allowFailures     bool
	options           map[string]any
	progressCallbacks []func(ctx context.Context, batch *Batch)
	thenCallbacks     []func(ctx context.Context, batch *Batch)
	catchCallbacks    []FailureCallback
	finallyCallbacks  []func(ctx context.Context, batch *Batch)
	beforeCallbacks   []func(ctx context.Context, batch *Batch)
	dispatcher        QueueingDispatcher
	batchRepo         BatchRepository
	eventFunc         EventFunc
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
func (p *PendingBatch) Catch(fn FailureCallback) *PendingBatch {
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

// DisallowFailures disables failure tolerance for the batch.
func (p *PendingBatch) DisallowFailures() *PendingBatch {
	p.allowFailures = false

	return p
}

// OnFailure enables failure tolerance and registers typed failure callbacks.
func (p *PendingBatch) OnFailure(callbacks ...FailureCallback) *PendingBatch {
	p.allowFailures = true
	p.catchCallbacks = append(p.catchCallbacks, callbacks...)

	return p
}

// WithOption sets a custom option on the batch.
func (p *PendingBatch) WithOption(key string, value any) *PendingBatch {
	if p.options == nil {
		p.options = make(map[string]any)
	}

	p.options[key] = value

	return p
}

// Connection returns the configured connection.
func (p *PendingBatch) Connection() string { return p.connection }

// Queue returns the configured queue.
func (p *PendingBatch) Queue() string { return p.queue }

// GetName returns the batch name.
func (p *PendingBatch) GetName() string { return p.name }

// Jobs returns the configured jobs.
func (p *PendingBatch) Jobs() []any { return p.jobs }

// AllowsFailures reports whether the batch allows failures.
func (p *PendingBatch) AllowsFailures() bool { return p.allowFailures }

// CatchCallbacks returns the catch callbacks.
func (p *PendingBatch) CatchCallbacks() []FailureCallback {
	return p.catchCallbacks
}

// FailureCallbacks returns callbacks registered through OnFailure or Catch.
func (p *PendingBatch) FailureCallbacks() []FailureCallback {
	return p.catchCallbacks
}

// Options returns the batch options map.
func (p *PendingBatch) Options() map[string]any {
	opts := make(map[string]any)

	for k, v := range p.options {
		opts[k] = v
	}

	if p.allowFailures {
		opts["allowFailures"] = true
	}

	return opts
}

// Dispatch creates and persists the batch, then dispatches all jobs.
func (p *PendingBatch) Dispatch(ctx context.Context) (*Batch, error) {
	id, err := generateBatchID()

	if err != nil {
		return nil, err
	}

	batch := p.newBatch(id)

	// Persist the batch if a repository is available.
	if p.batchRepo != nil {
		if err = p.batchRepo.Store(ctx, batch); err != nil {
			return nil, err
		}
	}

	// Set batch ID on batchable jobs.
	p.prepareBatchedJobs(batch.ID)

	// Invoke before callbacks.
	for _, fn := range p.beforeCallbacks {
		fn(ctx, batch)
	}

	// Dispatch each job.
	for _, job := range p.jobs {
		if err = p.dispatcher.DispatchToQueue(ctx, job); err != nil {
			if p.batchRepo != nil {
				_ = p.batchRepo.Delete(ctx, batch.ID)
			}

			return nil, err
		}
	}

	if p.eventFunc != nil {
		p.eventFunc(BatchDispatched{Batch: batch})
	}

	return batch, nil
}

// prepareBatchedJobs sets the batch ID on all jobs that implement WithBatchID.
func (p *PendingBatch) prepareBatchedJobs(batchID string) {
	for _, job := range p.jobs {
		if b, ok := job.(interface{ WithBatchID(string) }); ok {
			b.WithBatchID(batchID)
		}
	}
}

// DispatchIf dispatches the batch only if the condition is true.
func (p *PendingBatch) DispatchIf(ctx context.Context, condition bool) (*Batch, error) {
	if !condition {
		return nil, nil
	}

	return p.Dispatch(ctx)
}

// DispatchUnless dispatches the batch only if the condition is false.
func (p *PendingBatch) DispatchUnless(ctx context.Context, condition bool) (*Batch, error) {
	return p.DispatchIf(ctx, !condition)
}

// DispatchAfterResponse creates and persists the batch, then defers job dispatch.
func (p *PendingBatch) DispatchAfterResponse(ctx context.Context) (*Batch, error) {
	id, err := generateBatchID()

	if err != nil {
		return nil, err
	}

	batch := p.newBatch(id)

	if p.batchRepo != nil {
		if err = p.batchRepo.Store(ctx, batch); err != nil {
			return nil, err
		}
	}

	p.prepareBatchedJobs(batch.ID)

	for _, fn := range p.beforeCallbacks {
		fn(ctx, batch)
	}

	for _, job := range p.jobs {
		if err = p.dispatcher.DispatchAfterResponse(ctx, job); err != nil {
			if p.batchRepo != nil {
				_ = p.batchRepo.Delete(ctx, batch.ID)
			}

			return nil, err
		}
	}

	if p.eventFunc != nil {
		p.eventFunc(BatchDispatched{Batch: batch})
	}

	return batch, nil
}

func (p *PendingBatch) newBatch(id string) *Batch {
	opts := make(map[string]any)

	for k, v := range p.options {
		opts[k] = v
	}

	if p.allowFailures {
		opts["allowFailures"] = true
	}

	return &Batch{
		ID:                id,
		Name:              p.name,
		TotalJobs:         len(p.jobs),
		PendingJobs:       len(p.jobs),
		Options:           opts,
		CreatedAt:         time.Now(),
		ProgressCallbacks: p.progressCallbacks,
		ThenCallbacks:     p.thenCallbacks,
		CatchCallbacks:    p.catchCallbacks,
		FinallyCallbacks:  p.finallyCallbacks,
		repo:              p.batchRepo,
		dispatcher:        p.dispatcher,
		eventFunc:         p.eventFunc,
	}
}

func generateBatchID() (string, error) {
	b := make([]byte, 16)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
