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
	CatchCallbacks    []func(ctx context.Context, batch *Batch, err error)
	FinallyCallbacks  []func(ctx context.Context, batch *Batch)

	repo       BatchRepository
	dispatcher QueueingDispatcher
}

// NewBatchWithRepo creates a Batch with the given ID and repository.
func NewBatchWithRepo(id string, repo BatchRepository) *Batch {
	return &Batch{
		ID:      id,
		Options: make(map[string]any),
		repo:    repo,
	}
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

// Fresh retrieves the latest batch state from the repository.
func (b *Batch) Fresh(ctx context.Context) (*Batch, error) {
	if b.repo == nil {
		return b, nil
	}

	return b.repo.Get(ctx, b.ID)
}

// Progress returns the batch completion percentage (0–100).
func (b *Batch) Progress() float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.TotalJobs == 0 {
		return 100
	}

	return float64(b.TotalJobs-b.PendingJobs) / float64(b.TotalJobs) * 100
}

// ProcessedJobs returns the number of jobs that have been processed.
func (b *Batch) ProcessedJobs() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.TotalJobs - b.PendingJobs
}

// Delete removes the batch from the repository.
func (b *Batch) Delete(ctx context.Context) error {
	if b.repo != nil {
		return b.repo.Delete(ctx, b.ID)
	}

	return nil
}

// AllowsFailures reports whether the batch is configured to tolerate job failures.
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

// Add adds jobs to a dispatched batch and dispatches them.
func (b *Batch) Add(ctx context.Context, jobs []any) error {
	if b.repo != nil {
		if err := b.repo.IncrementTotalJobs(ctx, b.ID, len(jobs)); err != nil {
			return err
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

// MarshalJSON serialises the batch to JSON.
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

// UpdatedBatchJobCounts is a DTO for updated batch job counts.
type UpdatedBatchJobCounts struct {
	PendingJobs int
	FailedJobs  int
}

// AllJobsRanExactlyOnce reports whether all jobs completed without failure.
func (c UpdatedBatchJobCounts) AllJobsRanExactlyOnce() bool {
	return c.PendingJobs == 0 && c.FailedJobs == 0
}
