package bus

import (
	"context"
	"sync"
	"time"
)

// Batch tracks the state of a group of dispatched jobs.
type Batch struct {
	mu sync.RWMutex

	ID          string
	Name        string
	TotalJobs   int
	PendingJobs int
	FailedJobs  int
	FailedJobIDs []string
	Options     map[string]any
	CreatedAt   time.Time
	CancelledAt *time.Time
	FinishedAt  *time.Time

	// Callbacks — invoked by the dispatcher at batch lifecycle events.
	ProgressCallbacks []func(ctx context.Context, batch *Batch)
	ThenCallbacks     []func(ctx context.Context, batch *Batch)
	CatchCallbacks    []func(ctx context.Context, batch *Batch, err error)
	FinallyCallbacks  []func(ctx context.Context, batch *Batch)

	repo BatchRepository
}

// Finished reports whether all jobs have run (pending == 0).
func (b *Batch) Finished() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.PendingJobs == 0
}

// Cancelled reports whether the batch was cancelled.
func (b *Batch) Cancelled() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.CancelledAt != nil
}

// HasFailures reports whether any jobs have failed.
func (b *Batch) HasFailures() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.FailedJobs > 0
}

// Cancel marks the batch as cancelled.
func (b *Batch) Cancel(ctx context.Context) error {
	if b.repo != nil {
		return b.repo.Cancel(ctx, b.ID)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	b.CancelledAt = &now

	return nil
}

// RecordSuccessfulJob decrements PendingJobs. Returns updated counts.
func (b *Batch) RecordSuccessfulJob(ctx context.Context) (*UpdatedBatchJobCounts, error) {
	if b.repo != nil {
		return b.repo.DecrementPendingJobs(ctx, b.ID)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.PendingJobs > 0 {
		b.PendingJobs--
	}

	return &UpdatedBatchJobCounts{PendingJobs: b.PendingJobs, FailedJobs: b.FailedJobs}, nil
}

// RecordFailedJob increments FailedJobs and decrements PendingJobs.
func (b *Batch) RecordFailedJob(ctx context.Context, failedJobID string) (*UpdatedBatchJobCounts, error) {
	if b.repo != nil {
		return b.repo.IncrementFailedJobs(ctx, b.ID, failedJobID)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.FailedJobs++
	b.FailedJobIDs = append(b.FailedJobIDs, failedJobID)
	if b.PendingJobs > 0 {
		b.PendingJobs--
	}

	return &UpdatedBatchJobCounts{PendingJobs: b.PendingJobs, FailedJobs: b.FailedJobs}, nil
}

// UpdatedBatchJobCounts is a DTO for updated batch job counts.
type UpdatedBatchJobCounts struct {
	PendingJobs                  int
	FailedJobs                   int
	AllJobsHaveRanExactlyOnce    bool
}
