package telescope_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/debugbar"
	"github.com/bedrock/packages/debugbar/storage"
)

// ─── Helpers ─────────────────────────────────────────────────────────────────

func newTestScope(t *testing.T) (*debugbar.DebugBar, *storage.InMemoryRepository) {
	t.Helper()

	repo := storage.NewInMemoryRepository()
	scope := debugbar.New(debugbar.WithRepository(repo))
	scope.StartRecording()

	return scope, repo
}

// ─── Core recording ──────────────────────────────────────────────────────────

func TestDebugBarRecordsEntryWhenRecordingEnabled(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)

	scope.RecordLog(debugbar.NewEntry(debugbar.EntryTypeLog, map[string]any{
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

func TestDebugBarDoesNotRecordWhenRecordingDisabled(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	scope.StopRecording()

	scope.RecordLog(debugbar.NewEntry(debugbar.EntryTypeLog, map[string]any{
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

func TestDebugBarDoesNotRecordWhenPaused(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	scope.PauseRecording()

	scope.RecordLog(debugbar.NewEntry(debugbar.EntryTypeLog, map[string]any{
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

func TestDebugBarWithoutRecordingExecutesFnAndResumes(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)

	scope.WithoutRecording(func() {
		scope.RecordLog(debugbar.NewEntry(debugbar.EntryTypeLog, map[string]any{
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

func TestDebugBarNewBatchGeneratesUUID(t *testing.T) {
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

func TestDebugBarEntriesAreTaggedWithBatchID(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)

	batchID := scope.CurrentBatchID()

	scope.RecordLog(debugbar.NewEntry(debugbar.EntryTypeLog, map[string]any{
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

func TestDebugBarFilterPreventsEntryFromBeingRecorded(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)

	scope.Filter(func(e *debugbar.IncomingEntry) bool {
		return e.Type != debugbar.EntryTypeLog
	})

	scope.RecordLog(debugbar.NewEntry(debugbar.EntryTypeLog, map[string]any{
		"level":   "info",
		"message": "filtered out",
	}))

	scope.RecordCache(debugbar.NewEntry(debugbar.EntryTypeCache, map[string]any{
		"type": "hit",
		"key":  "foo",
	}))

	if err := scope.Store(context.Background()); err != nil {
		t.Fatal(err)
	}

	if repo.Count() != 1 {
		t.Fatalf("expected 1 entry after filter, got %d", repo.Count())
	}

	if repo.Entries()[0].Type != debugbar.EntryTypeCache {
		t.Fatalf("expected cache entry to survive filter, got %q", repo.Entries()[0].Type)
	}
}

func TestDebugBarFilterBatchPreventsAllEntriesFromBeingStored(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)

	scope.FilterBatch(func(_ []*debugbar.IncomingEntry) bool {
		return false // drop everything
	})

	scope.RecordLog(debugbar.NewEntry(debugbar.EntryTypeLog, map[string]any{
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

func TestDebugBarTagCallbackAddsTagsToEntry(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)

	scope.Tag(func(e *debugbar.IncomingEntry) []string {
		return []string{"custom-tag"}
	})

	scope.RecordLog(debugbar.NewEntry(debugbar.EntryTypeLog, map[string]any{
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

func TestDebugBarAfterStoringHookIsCalledWithBatchAndEntries(t *testing.T) {
	t.Parallel()

	scope, _ := newTestScope(t)

	var hookBatchID string

	var hookEntryCount int

	scope.AfterStoring(func(batchID string, entries []*debugbar.IncomingEntry) {
		hookBatchID = batchID
		hookEntryCount = len(entries)
	})

	batchID := scope.CurrentBatchID()

	scope.RecordLog(debugbar.NewEntry(debugbar.EntryTypeLog, map[string]any{
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

func TestDebugBarAfterRecordingHookIsCalledForEachEntry(t *testing.T) {
	t.Parallel()

	scope, _ := newTestScope(t)

	var recorded []*debugbar.IncomingEntry

	scope.AfterRecording(func(e *debugbar.IncomingEntry) {
		recorded = append(recorded, e)
	})

	scope.RecordLog(debugbar.NewEntry(debugbar.EntryTypeLog, map[string]any{
		"level":   "info",
		"message": "first",
	}))

	scope.RecordLog(debugbar.NewEntry(debugbar.EntryTypeLog, map[string]any{
		"level":   "debug",
		"message": "second",
	}))

	if len(recorded) != 2 {
		t.Fatalf("expected 2 after-recording callbacks, got %d", len(recorded))
	}
}

// ─── Flush ────────────────────────────────────────────────────────────────────

func TestDebugBarFlushDiscardsQueuedEntries(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)

	scope.RecordLog(debugbar.NewEntry(debugbar.EntryTypeLog, map[string]any{
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
