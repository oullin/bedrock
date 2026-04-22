package bus

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// Batch tracks the state of a group of dispatched jobs.
type Batch struct {
	mu sync.RWMutex

	ID           string
	Name         string
	TotalJobs    int
	PendingJobs  int
	FailedJobs   int
	FailedJobIDs []string
	Options      map[string]any
	CreatedAt    time.Time
	CancelledAt  *time.Time
	FinishedAt   *time.Time

	// Callbacks — invoked by the dispatcher at batch lifecycle events.
	ProgressCallbacks []func(ctx context.Context, batch *Batch)
	ThenCallbacks     []func(ctx context.Context, batch *Batch)
	CatchCallbacks    []FailureCallback
	FinallyCallbacks  []func(ctx context.Context, batch *Batch)

	repo       BatchRepository
	dispatcher QueueingDispatcher
	eventFunc  EventFunc
	started    bool
}

// NewBatchWithRepo creates a Batch with the given ID and repository.

// Finished reports whether all jobs have run (pending == 0).

// Cancelled reports whether the batch was cancelled.

// HasFailures reports whether any jobs have failed.

// Cancel marks the batch as cancelled.

// RecordSuccessfulJob decrements PendingJobs. Returns updated counts.

// RecordFailedJob increments FailedJobs and decrements PendingJobs.

// Fresh retrieves the latest batch state from the repository.

// Progress returns the batch completion percentage (0–100).

// ProcessedJobs returns the number of jobs that have been processed.

// Delete removes the batch from the repository.

// AllowsFailures reports whether the batch is configured to tolerate job failures.

// Add adds jobs to a dispatched batch and dispatches them.

// MarshalJSON serialises the batch to JSON.

// UpdatedBatchJobCounts is a DTO for updated batch job counts.
type UpdatedBatchJobCounts struct {
	PendingJobs int
	FailedJobs  int
}

func NewBatchWithRepo(id string, repo BatchRepository) *Batch {
	return &Batch{
		ID:      id,
		Options: make(map[string]any),
		repo:    repo,
	}
}

// SetEventFunc sets the event callback for batch lifecycle events.
func (b *Batch) SetEventFunc(fn EventFunc) {
	b.eventFunc = fn
}

// SetDispatcher sets the dispatcher used when adding jobs to an existing batch.
func (b *Batch) SetDispatcher(dispatcher QueueingDispatcher) {
	b.dispatcher = dispatcher
}

func (b *Batch) Finished() bool {
	b.mu.RLock()

	defer b.mu.RUnlock()

	return b.PendingJobs == 0
}

func (b *Batch) Cancelled() bool {
	b.mu.RLock()

	defer b.mu.RUnlock()

	return b.CancelledAt != nil
}

func (b *Batch) HasFailures() bool {
	b.mu.RLock()

	defer b.mu.RUnlock()

	return b.FailedJobs > 0
}

func (b *Batch) Cancel(ctx context.Context) error {
	now := time.Now()

	if b.repo != nil {
		if err := b.repo.Cancel(ctx, b.ID); err != nil {
			return err
		}
	}

	b.mu.Lock()
	b.CancelledAt = &now
	b.mu.Unlock()

	if b.eventFunc != nil {
		b.eventFunc(BatchCanceled{Batch: b})
	}

	return nil
}

func (b *Batch) RecordSuccessfulJob(ctx context.Context) (*UpdatedBatchJobCounts, error) {
	var counts *UpdatedBatchJobCounts
	started := b.shouldDispatchStarted()

	if b.repo != nil {
		var err error

		counts, err = b.repo.DecrementPendingJobs(ctx, b.ID)

		if err != nil {
			return nil, err
		}

		b.mu.Lock()
		b.PendingJobs = counts.PendingJobs
		b.FailedJobs = counts.FailedJobs
		b.mu.Unlock()
	} else {
		b.mu.Lock()

		if b.PendingJobs > 0 {
			b.PendingJobs--
		}

		counts = &UpdatedBatchJobCounts{PendingJobs: b.PendingJobs, FailedJobs: b.FailedJobs}
		b.mu.Unlock()
	}

	if started && b.markStarted() {
		b.dispatchStarted()
	}

	b.InvokeProgressCallbacks(ctx)

	if counts.AllJobsRanExactlyOnce() {
		b.InvokeThenCallbacks(ctx)

		if b.repo != nil {
			_ = b.repo.MarkAsFinished(ctx, b.ID)
		}

		if b.eventFunc != nil {
			b.eventFunc(BatchFinished{Batch: b})
		}
	}

	if counts.PendingJobs == 0 {
		b.InvokeFinallyCallbacks(ctx)
	}

	return counts, nil
}

func (b *Batch) RecordFailedJob(ctx context.Context, failedJobID string, err error) (*UpdatedBatchJobCounts, error) {
	var counts *UpdatedBatchJobCounts
	started := b.shouldDispatchStarted()
	alreadyFailed := b.HasFailures()
	allowsFailures := b.AllowsFailures()

	if b.repo != nil {
		var repoErr error

		counts, repoErr = b.repo.IncrementFailedJobs(ctx, b.ID, failedJobID)

		if repoErr != nil {
			return nil, repoErr
		}

		b.mu.Lock()
		b.PendingJobs = counts.PendingJobs
		b.FailedJobs = counts.FailedJobs
		b.FailedJobIDs = append(b.FailedJobIDs, failedJobID)
		b.mu.Unlock()
	} else {
		b.mu.Lock()

		b.FailedJobs++
		b.FailedJobIDs = append(b.FailedJobIDs, failedJobID)

		if allowsFailures && b.PendingJobs > 0 {
			b.PendingJobs--
		} else {
			b.PendingJobs = 0
			now := time.Now()
			b.CancelledAt = &now
			b.FinishedAt = &now
		}

		counts = &UpdatedBatchJobCounts{PendingJobs: b.PendingJobs, FailedJobs: b.FailedJobs}
		b.mu.Unlock()
	}

	if started && b.markStarted() {
		b.dispatchStarted()
	}

	if !alreadyFailed {
		b.InvokeCatchCallbacks(ctx, err)
	}

	if allowsFailures {
		b.InvokeProgressCallbacks(ctx)

		return counts, nil
	}

	if b.repo != nil {
		if cancelErr := b.repo.Cancel(ctx, b.ID); cancelErr != nil {
			return nil, cancelErr
		}

		_ = b.repo.MarkAsFinished(ctx, b.ID)
		counts.PendingJobs = 0

		b.mu.Lock()
		now := time.Now()
		b.PendingJobs = 0
		b.CancelledAt = &now
		b.FinishedAt = &now
		b.mu.Unlock()
	}

	if b.eventFunc != nil {
		b.eventFunc(BatchCanceled{Batch: b})
	}

	if counts.PendingJobs == 0 {
		b.InvokeFinallyCallbacks(ctx)
	}

	return counts, nil
}

func (b *Batch) Fresh(ctx context.Context) (*Batch, error) {
	if b.repo == nil {
		return b, nil
	}

	return b.repo.Get(ctx, b.ID)
}

func (b *Batch) Progress() float64 {
	b.mu.RLock()

	defer b.mu.RUnlock()

	if b.TotalJobs == 0 {
		return 100
	}

	return float64(b.TotalJobs-b.PendingJobs) / float64(b.TotalJobs) * 100
}

func (b *Batch) ProcessedJobs() int {
	b.mu.RLock()

	defer b.mu.RUnlock()

	return b.TotalJobs - b.PendingJobs
}

func (b *Batch) Delete(ctx context.Context) error {
	if b.repo != nil {
		return b.repo.Delete(ctx, b.ID)
	}

	return nil
}

func (b *Batch) AllowsFailures() bool {
	b.mu.RLock()

	defer b.mu.RUnlock()

	v, ok := b.Options["allowFailures"]

	if !ok {
		return false
	}

	allowed, _ := v.(bool)

	return allowed
}

// Canceled is an alias for Cancelled (American spelling).
func (b *Batch) Canceled() bool {
	return b.Cancelled()
}

// HasProgressCallbacks reports whether the batch has progress callbacks registered.
func (b *Batch) HasProgressCallbacks() bool {
	b.mu.RLock()

	defer b.mu.RUnlock()

	return len(b.ProgressCallbacks) > 0
}

// HasThenCallbacks reports whether the batch has then callbacks registered.
func (b *Batch) HasThenCallbacks() bool {
	b.mu.RLock()

	defer b.mu.RUnlock()

	return len(b.ThenCallbacks) > 0
}

// HasCatchCallbacks reports whether the batch has catch callbacks registered.
func (b *Batch) HasCatchCallbacks() bool {
	b.mu.RLock()

	defer b.mu.RUnlock()

	return len(b.CatchCallbacks) > 0
}

// HasFinallyCallbacks reports whether the batch has finally callbacks registered.
func (b *Batch) HasFinallyCallbacks() bool {
	b.mu.RLock()

	defer b.mu.RUnlock()

	return len(b.FinallyCallbacks) > 0
}

// InvokeProgressCallbacks calls all progress callbacks.
func (b *Batch) InvokeProgressCallbacks(ctx context.Context) {
	b.mu.RLock()
	callbacks := b.ProgressCallbacks
	b.mu.RUnlock()

	for _, fn := range callbacks {
		fn(ctx, b)
	}
}

// InvokeThenCallbacks calls all then callbacks.
func (b *Batch) InvokeThenCallbacks(ctx context.Context) {
	b.mu.RLock()
	callbacks := b.ThenCallbacks
	b.mu.RUnlock()

	for _, fn := range callbacks {
		fn(ctx, b)
	}
}

// InvokeCatchCallbacks calls all catch callbacks.
func (b *Batch) InvokeCatchCallbacks(ctx context.Context, err error) {
	b.mu.RLock()
	callbacks := b.CatchCallbacks
	b.mu.RUnlock()

	for _, fn := range callbacks {
		fn(ctx, b, err)
	}
}

// InvokeFinallyCallbacks calls all finally callbacks.
func (b *Batch) InvokeFinallyCallbacks(ctx context.Context) {
	b.mu.RLock()
	callbacks := b.FinallyCallbacks
	b.mu.RUnlock()

	for _, fn := range callbacks {
		fn(ctx, b)
	}
}

func (b *Batch) Add(ctx context.Context, jobs []any) error {
	if b.repo != nil {
		if err := b.repo.IncrementTotalJobs(ctx, b.ID, len(jobs)); err != nil {
			return err
		}
	}

	for _, job := range jobs {
		if batchable, ok := job.(interface{ WithBatchID(string) }); ok {
			batchable.WithBatchID(b.ID)
		}
	}

	b.mu.Lock()
	b.TotalJobs += len(jobs)
	b.PendingJobs += len(jobs)
	b.mu.Unlock()

	if b.dispatcher != nil {
		for _, job := range jobs {
			if err := b.dispatcher.DispatchToQueue(ctx, job); err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Batch) MarshalJSON() ([]byte, error) {
	b.mu.RLock()

	defer b.mu.RUnlock()

	return json.Marshal(struct {
		ID           string         `json:"id"`
		Name         string         `json:"name"`
		TotalJobs    int            `json:"total_jobs"`
		PendingJobs  int            `json:"pending_jobs"`
		FailedJobs   int            `json:"failed_jobs"`
		FailedJobIDs []string       `json:"failed_job_ids"`
		Progress     float64        `json:"progress"`
		Options      map[string]any `json:"options"`
		CreatedAt    time.Time      `json:"created_at"`
		CancelledAt  *time.Time     `json:"cancelled_at"`
		FinishedAt   *time.Time     `json:"finished_at"`
	}{
		ID:           b.ID,
		Name:         b.Name,
		TotalJobs:    b.TotalJobs,
		PendingJobs:  b.PendingJobs,
		FailedJobs:   b.FailedJobs,
		FailedJobIDs: b.FailedJobIDs,
		Progress:     b.progressLocked(),
		Options:      b.Options,
		CreatedAt:    b.CreatedAt,
		CancelledAt:  b.CancelledAt,
		FinishedAt:   b.FinishedAt,
	})
}

func (b *Batch) progressLocked() float64 {
	if b.TotalJobs == 0 {
		return 100
	}

	return float64(b.TotalJobs-b.PendingJobs) / float64(b.TotalJobs) * 100
}

func (b *Batch) shouldDispatchStarted() bool {
	b.mu.RLock()

	defer b.mu.RUnlock()

	return !b.started && b.TotalJobs > 0 && b.PendingJobs == b.TotalJobs && b.FailedJobs == 0
}

func (b *Batch) markStarted() bool {
	b.mu.Lock()

	defer b.mu.Unlock()

	if b.started {
		return false
	}

	b.started = true

	return true
}

func (b *Batch) dispatchStarted() {
	if b.eventFunc != nil {
		b.eventFunc(BatchStarted{Batch: b})
	}
}

// AllJobsRanExactlyOnce reports whether all jobs completed without failure.
func (c UpdatedBatchJobCounts) AllJobsRanExactlyOnce() bool {
	return c.PendingJobs == 0 && c.FailedJobs == 0
}
