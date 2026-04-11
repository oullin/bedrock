package client

import "sync"

// Batch groups multiple HTTP requests to be executed together. Unlike Pool,
// Batch allows incremental addition of requests before execution.
type Batch struct {
	factory   *Factory
	mu        sync.Mutex
	callbacks []PoolCallback
	started   bool
}

// NewBatch creates a Batch backed by the given factory.
func NewBatch(factory *Factory) *Batch {
	return &Batch{factory: factory}
}

// Add appends request callbacks to the batch.
func (b *Batch) Add(callbacks ...PoolCallback) *Batch {
	b.mu.Lock()

	defer b.mu.Unlock()

	if b.started {
		return b
	}

	b.callbacks = append(b.callbacks, callbacks...)

	return b
}

// Execute runs all batched requests concurrently and returns the results.
func (b *Batch) Execute() ([]*PoolResult, error) {
	b.mu.Lock()

	if b.started {
		b.mu.Unlock()

		return nil, ErrBatchInProgress
	}

	b.started = true
	callbacks := make([]PoolCallback, len(b.callbacks))
	copy(callbacks, b.callbacks)
	b.mu.Unlock()

	pool := NewPool(b.factory)

	return pool.Concurrent(callbacks), nil
}

// Pending returns the number of requests not yet executed.
func (b *Batch) Pending() int {
	b.mu.Lock()

	defer b.mu.Unlock()

	return len(b.callbacks)
}
