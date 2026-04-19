package watchers

import (
	"github.com/bedrock/packages/telescope"
)

// JobStatus constants mirror the three lifecycle states recorded by Laravel's
// JobWatcher.

// JobWatcher monitors queued job lifecycle events (pending, processed, failed)
// and records them as Telescope entries. It mirrors Laravel's JobWatcher class.
//
// Options:
//   - "ignore" ([]string): fully-qualified job type names to skip.
type JobWatcher struct {
	telescope.BaseWatcher
}

// NewJobWatcher creates a JobWatcher with the given options.

// Register is a no-op for JobWatcher; callers drive it via Pending, Processed,
// and Failed.

// ShouldIgnore reports whether the job type should be skipped.

// JobMeta carries metadata about a queued job dispatched into the queue.
type JobMeta struct {
	// UUID is the job's unique identifier.
	UUID string
	// Type is the fully-qualified job struct or interface name.
	Type string
	// BatchID links the job to the Telescope batch that dispatched it.
	BatchID string
	// Connection is the queue connection name (e.g. "redis", "database").
	Connection string
	// Queue is the queue name (e.g. "default", "mail").
	Queue string
	// MaxTries is the maximum number of attempts.
	MaxTries int
	// Timeout is the job timeout in seconds.
	Timeout int
	// Tags are additional searchable tags.
	Tags []string
}

const (
	JobStatusPending   = "pending"
	JobStatusProcessed = "processed"
	JobStatusFailed    = "failed"
)

func NewJobWatcher(t *telescope.Telescope, options map[string]any) *JobWatcher {
	w := &JobWatcher{}
	w.SetTelescope(t)
	w.Options = options

	return w
}

func (w *JobWatcher) Register(_ any) error { return nil }

func (w *JobWatcher) ShouldIgnore(jobType string) bool {
	for _, name := range w.StringsOption("ignore") {
		if name == jobType {
			return true
		}
	}

	return false
}

// Pending records a pending (dispatched) job entry.
func (w *JobWatcher) Pending(meta JobMeta) {
	if w.ShouldIgnore(meta.Type) {
		return
	}

	content := map[string]any{
		"status":     JobStatusPending,
		"name":       meta.Type,
		"connection": meta.Connection,
		"queue":      meta.Queue,
		"max_tries":  meta.MaxTries,
		"timeout":    meta.Timeout,
	}

	entry := telescope.NewEntry(telescope.EntryTypeJob, content)

	if meta.UUID != "" {
		entry.UUID = meta.UUID
	}

	if meta.BatchID != "" {
		entry.WithBatchID(meta.BatchID)
	}

	entry.AddTags(meta.Tags...)

	w.Scope().RecordJob(entry)
}

// Processed records a successfully processed job by updating the existing
// entry status. If a matching pending entry exists in the repository it is
// updated; otherwise a new entry is created.
func (w *JobWatcher) Processed(jobUUID string, jobType string) {
	if w.ShouldIgnore(jobType) {
		return
	}

	update := telescope.NewEntryUpdate(jobUUID, telescope.EntryTypeJob, map[string]any{
		"status": JobStatusProcessed,
	})

	w.Scope().Update([]*telescope.EntryUpdate{update})
}

// Failed records a failed job entry, capturing the error details.
func (w *JobWatcher) Failed(jobUUID string, jobType string, err error) {
	if w.ShouldIgnore(jobType) {
		return
	}

	changes := map[string]any{
		"status": JobStatusFailed,
	}

	if err != nil {
		changes["exception"] = err.Error()
	}

	update := telescope.NewEntryUpdate(jobUUID, telescope.EntryTypeJob, changes)

	w.Scope().Update([]*telescope.EntryUpdate{update})
}
