package telescope_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/telescope"
	"github.com/bedrock/packages/telescope/storage"
)

// ─── Helpers ─────────────────────────────────────────────────────────────────

func newTestScope(t *testing.T) (*telescope.Telescope, *storage.InMemoryRepository) {
	t.Helper()

	repo := storage.NewInMemoryRepository()
	scope := telescope.New(telescope.WithRepository(repo))
	scope.StartRecording()

	return scope, repo
}

// ─── Core recording ──────────────────────────────────────────────────────────

func TestTelescopeRecordsEntryWhenRecordingEnabled(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)

	scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{
		"level":   "info",
		"message": "hello",
	}))

	if err := scope.Store(context.Background()); err != nil {
		t.Fatal(err)
	}

	if repo.Count() != 1 {
		t.Fatalf("expected 1 entry, got %d", repo.Count())
	}
}

func TestTelescopeDoesNotRecordWhenRecordingDisabled(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	scope.StopRecording()

	scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{
		"level":   "info",
		"message": "should not be recorded",
	}))

	if err := scope.Store(context.Background()); err != nil {
		t.Fatal(err)
	}

	if repo.Count() != 0 {
		t.Fatalf("expected 0 entries, got %d", repo.Count())
	}
}

func TestTelescopeDoesNotRecordWhenPaused(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	scope.PauseRecording()

	scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{
		"level":   "info",
		"message": "paused",
	}))

	scope.ResumeRecording()

	if err := scope.Store(context.Background()); err != nil {
		t.Fatal(err)
	}

	if repo.Count() != 0 {
		t.Fatalf("expected 0 entries while paused, got %d", repo.Count())
	}
}

func TestTelescopeWithoutRecordingExecutesFnAndResumes(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)

	scope.WithoutRecording(func() {
		scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{
			"level":   "info",
			"message": "inside without recording",
		}))
	})

	// Recording should be active again.
	if !scope.IsRecording() {
		t.Fatal("expected recording to resume after WithoutRecording")
	}

	if err := scope.Store(context.Background()); err != nil {
		t.Fatal(err)
	}

	if repo.Count() != 0 {
		t.Fatalf("expected 0 entries inside WithoutRecording, got %d", repo.Count())
	}
}

// ─── Batch management ────────────────────────────────────────────────────────

func TestTelescopeNewBatchGeneratesUUID(t *testing.T) {
	t.Parallel()

	scope, _ := newTestScope(t)

	b1 := scope.NewBatch()
	b2 := scope.NewBatch()

	if b1 == "" || b2 == "" {
		t.Fatal("NewBatch returned empty string")
	}

	if b1 == b2 {
		t.Fatal("NewBatch should return a different UUID each time")
	}
}

func TestTelescopeEntriesAreTaggedWithBatchID(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)

	batchID := scope.CurrentBatchID()

	scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{
		"level":   "info",
		"message": "test",
	}))

	if err := scope.Store(context.Background()); err != nil {
		t.Fatal(err)
	}

	entries := repo.Entries()

	if len(entries) == 0 {
		t.Fatal("no entries stored")
	}

	if entries[0].BatchID != batchID {
		t.Fatalf("expected batch ID %q, got %q", batchID, entries[0].BatchID)
	}
}

// ─── Filters ─────────────────────────────────────────────────────────────────

func TestTelescopeFilterPreventsEntryFromBeingRecorded(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)

	scope.Filter(func(e *telescope.IncomingEntry) bool {
		return e.Type != telescope.EntryTypeLog
	})

	scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{
		"level":   "info",
		"message": "filtered out",
	}))

	scope.RecordCache(telescope.NewEntry(telescope.EntryTypeCache, map[string]any{
		"type": "hit",
		"key":  "foo",
	}))

	if err := scope.Store(context.Background()); err != nil {
		t.Fatal(err)
	}

	if repo.Count() != 1 {
		t.Fatalf("expected 1 entry after filter, got %d", repo.Count())
	}

	if repo.Entries()[0].Type != telescope.EntryTypeCache {
		t.Fatalf("expected cache entry to survive filter, got %q", repo.Entries()[0].Type)
	}
}

func TestTelescopeFilterBatchPreventsAllEntriesFromBeingStored(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)

	scope.FilterBatch(func(_ []*telescope.IncomingEntry) bool {
		return false // drop everything
	})

	scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{
		"level":   "info",
		"message": "batch filtered",
	}))

	if err := scope.Store(context.Background()); err != nil {
		t.Fatal(err)
	}

	if repo.Count() != 0 {
		t.Fatalf("expected 0 entries after batch filter, got %d", repo.Count())
	}
}

// ─── Tag callbacks ────────────────────────────────────────────────────────────

func TestTelescopeTagCallbackAddsTagsToEntry(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)

	scope.Tag(func(e *telescope.IncomingEntry) []string {
		return []string{"custom-tag"}
	})

	scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{
		"level":   "info",
		"message": "tagged",
	}))

	if err := scope.Store(context.Background()); err != nil {
		t.Fatal(err)
	}

	entries := repo.Entries()

	if len(entries) == 0 {
		t.Fatal("no entries stored")
	}

	found := false

	for _, tag := range entries[0].Tags {
		if tag == "custom-tag" {
			found = true

			break
		}
	}

	if !found {
		t.Fatalf("expected \"custom-tag\" in tags, got %v", entries[0].Tags)
	}
}

// ─── After-storing hook ───────────────────────────────────────────────────────

func TestTelescopeAfterStoringHookIsCalledWithBatchAndEntries(t *testing.T) {
	t.Parallel()

	scope, _ := newTestScope(t)

	var hookBatchID string

	var hookEntryCount int

	scope.AfterStoring(func(batchID string, entries []*telescope.IncomingEntry) {
		hookBatchID = batchID
		hookEntryCount = len(entries)
	})

	batchID := scope.CurrentBatchID()

	scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{
		"level":   "info",
		"message": "after storing",
	}))

	if err := scope.Store(context.Background()); err != nil {
		t.Fatal(err)
	}

	if hookBatchID != batchID {
		t.Fatalf("hook received batch ID %q, expected %q", hookBatchID, batchID)
	}

	if hookEntryCount != 1 {
		t.Fatalf("hook received %d entries, expected 1", hookEntryCount)
	}
}

// ─── After-recording hook ─────────────────────────────────────────────────────

func TestTelescopeAfterRecordingHookIsCalledForEachEntry(t *testing.T) {
	t.Parallel()

	scope, _ := newTestScope(t)

	var recorded []*telescope.IncomingEntry

	scope.AfterRecording(func(e *telescope.IncomingEntry) {
		recorded = append(recorded, e)
	})

	scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{
		"level":   "info",
		"message": "first",
	}))

	scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{
		"level":   "debug",
		"message": "second",
	}))

	if len(recorded) != 2 {
		t.Fatalf("expected 2 after-recording callbacks, got %d", len(recorded))
	}
}

// ─── Flush ────────────────────────────────────────────────────────────────────

func TestTelescopeFlushDiscardsQueuedEntries(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)

	scope.RecordLog(telescope.NewEntry(telescope.EntryTypeLog, map[string]any{
		"level":   "info",
		"message": "will be flushed",
	}))

	scope.Flush()

	if err := scope.Store(context.Background()); err != nil {
		t.Fatal(err)
	}

	if repo.Count() != 0 {
		t.Fatalf("expected 0 entries after flush, got %d", repo.Count())
	}
}
