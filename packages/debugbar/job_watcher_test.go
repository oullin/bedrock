package telescope_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/debugbar"
	"github.com/bedrock/packages/debugbar/watchers"
)

func TestJobWatcherRecordsPendingJob(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewJobWatcher(scope, nil)

	w.Pending(watchers.JobMeta{
		UUID:       "job-uuid-1",
		Type:       "SendEmailJob",
		Connection: "redis",
		Queue:      "mail",
	})

	storeAndAssertCount(t, scope, repo, debugbar.EntryTypeJob, 1)

	e := repo.Entries()[0]

	if e.Content["status"] != watchers.JobStatusPending {
		t.Fatalf("expected status=pending, got %v", e.Content["status"])
	}

	if e.Content["name"] != "SendEmailJob" {
		t.Fatalf("expected name=SendEmailJob, got %v", e.Content["name"])
	}

	if e.Content["connection"] != "redis" {
		t.Fatalf("expected connection=redis, got %v", e.Content["connection"])
	}
}

func TestJobWatcherUpdatesStatusOnProcessed(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewJobWatcher(scope, nil)

	w.Pending(watchers.JobMeta{UUID: "job-abc", Type: "ProcessOrder"})

	// Store the pending entry first.
	if err := scope.Store(testContext()); err != nil {
		t.Fatal(err)
	}

	// Now mark as processed.
	w.Processed("job-abc", "ProcessOrder")

	if err := scope.Store(testContext()); err != nil {
		t.Fatal(err)
	}

	entries, err := repo.Get(debugbar.EntryTypeJob, debugbar.DefaultQueryOptions())
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	if entries[0].Content["status"] != watchers.JobStatusProcessed {
		t.Fatalf("expected status=processed, got %v", entries[0].Content["status"])
	}
}

func TestJobWatcherUpdatesStatusOnFailed(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewJobWatcher(scope, nil)

	w.Pending(watchers.JobMeta{UUID: "job-fail-1", Type: "ImportJob"})

	if err := scope.Store(testContext()); err != nil {
		t.Fatal(err)
	}

	w.Failed("job-fail-1", "ImportJob", errors.New("import error: file not found"))

	if err := scope.Store(testContext()); err != nil {
		t.Fatal(err)
	}

	entries, err := repo.Get(debugbar.EntryTypeJob, debugbar.DefaultQueryOptions())
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	entry := entries[0]

	if entry.Content["status"] != watchers.JobStatusFailed {
		t.Fatalf("expected status=failed, got %v", entry.Content["status"])
	}

	if entry.Content["exception"] == nil {
		t.Fatal("expected exception details in failed job entry")
	}
}

func TestJobWatcherIgnoresConfiguredJobTypes(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewJobWatcher(scope, map[string]any{
		"ignore": []string{"InternalCleanupJob"},
	})

	w.Pending(watchers.JobMeta{UUID: "j1", Type: "InternalCleanupJob"})
	w.Pending(watchers.JobMeta{UUID: "j2", Type: "SendEmailJob"})

	storeAndAssertCount(t, scope, repo, debugbar.EntryTypeJob, 1)
}

func TestJobWatcherBatchIDIsAttachedToEntry(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewJobWatcher(scope, nil)

	batchID := scope.CurrentBatchID()

	w.Pending(watchers.JobMeta{UUID: "j-batch", Type: "BatchedJob", BatchID: batchID})

	storeAndAssertCount(t, scope, repo, debugbar.EntryTypeJob, 1)

	// The batch ID from the debugbar should be stamped on the entry.
	entry := repo.Entries()[0]

	if entry.BatchID == "" {
		t.Fatal("expected batch ID to be set on the entry")
	}
}
