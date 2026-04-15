// Package testing provides test helpers for the DebugBar package.
// Import it in your *_test.go files to assert on captured entries.
package testing

import (
	"context"
	"testing"

	"github.com/bedrock/packages/debugbar"
	"github.com/bedrock/packages/debugbar/storage"
)

// DebugBarTestCase is an embeddable struct that provides test helpers
// mirroring the FeatureTestCase from Upstream's DebugBar test suite.
type DebugBarTestCase struct {
	T          *testing.T
	DebugBar  *debugbar.DebugBar
	Repository *storage.InMemoryRepository
}

// NewTestCase creates a DebugBarTestCase with recording already started. The
// in-memory repository is empty and ready for assertions.
func NewTestCase(t *testing.T) *DebugBarTestCase {
	t.Helper()

	repo := storage.NewInMemoryRepository()
	scope := debugbar.New(debugbar.WithRepository(repo))
	scope.StartRecording()

	return &DebugBarTestCase{
		T:          t,
		DebugBar:  scope,
		Repository: repo,
	}
}

// Store flushes queued entries to the repository, mirroring
// FeatureTestCase::terminateDebugBar().
func (tc *DebugBarTestCase) Store() {
	tc.T.Helper()

	if err := tc.DebugBar.Store(context.Background()); err != nil {
		tc.T.Fatalf("debugbar: store failed: %v", err)
	}
}

// Flush discards queued entries without persisting them.
func (tc *DebugBarTestCase) Flush() {
	tc.DebugBar.Flush()
}

// StartRecording enables recording.
func (tc *DebugBarTestCase) StartRecording() {
	tc.DebugBar.StartRecording()
}

// StopRecording disables recording.
func (tc *DebugBarTestCase) StopRecording() {
	tc.DebugBar.StopRecording()
}

// LoadEntries stores queued entries and returns all persisted results,
// mirroring FeatureTestCase::loadDebugBarEntries().
func (tc *DebugBarTestCase) LoadEntries() []*debugbar.EntryResult {
	tc.T.Helper()
	tc.Store()

	results, err := tc.Repository.Get("", debugbar.DefaultQueryOptions().WithLimit(500))
	if err != nil {
		tc.T.Fatalf("debugbar: load entries failed: %v", err)
	}

	return results
}

// LoadEntriesOfType stores queued entries and returns persisted results of the
// given entry type.
func (tc *DebugBarTestCase) LoadEntriesOfType(entryType string) []*debugbar.EntryResult {
	tc.T.Helper()
	tc.Store()

	results, err := tc.Repository.Get(entryType, debugbar.DefaultQueryOptions().WithLimit(500))
	if err != nil {
		tc.T.Fatalf("debugbar: load entries of type %q failed: %v", entryType, err)
	}

	return results
}

// AssertRecorded asserts that at least one entry of the given type was
// recorded.
func (tc *DebugBarTestCase) AssertRecorded(entryType string) []*debugbar.EntryResult {
	tc.T.Helper()

	entries := tc.LoadEntriesOfType(entryType)

	if len(entries) == 0 {
		tc.T.Errorf("expected at least one %q entry to be recorded, but none were found", entryType)
	}

	return entries
}

// AssertNotRecorded asserts that no entry of the given type was recorded.
func (tc *DebugBarTestCase) AssertNotRecorded(entryType string) {
	tc.T.Helper()

	entries := tc.LoadEntriesOfType(entryType)

	if len(entries) > 0 {
		tc.T.Errorf("expected no %q entries, but %d were found", entryType, len(entries))
	}
}

// AssertEntryCount asserts the number of entries of the given type.
func (tc *DebugBarTestCase) AssertEntryCount(entryType string, want int) {
	tc.T.Helper()

	entries := tc.LoadEntriesOfType(entryType)

	if len(entries) != want {
		tc.T.Errorf("expected %d %q entries, got %d", want, entryType, len(entries))
	}
}

// Reset clears all stored entries and restarts recording with a new batch.
func (tc *DebugBarTestCase) Reset() {
	tc.T.Helper()

	if err := tc.Repository.Clear(); err != nil {
		tc.T.Fatalf("debugbar: clear repository failed: %v", err)
	}

	tc.DebugBar.Flush()
	tc.DebugBar.StartRecording()
}
