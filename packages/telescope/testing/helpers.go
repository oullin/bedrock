// Package testing provides test helpers for the Telescope package.
// Import it in your *_test.go files to assert on captured entries.
package testing

import (
	"context"
	"testing"

	"github.com/bedrock/packages/telescope"
	"github.com/bedrock/packages/telescope/storage"
)

// TelescopeTestCase is an embeddable struct that provides test helpers
// mirroring the FeatureTestCase from the upstream Telescope test suite.
type TelescopeTestCase struct {
	T          *testing.T
	Telescope  *telescope.Telescope
	Repository *storage.InMemoryRepository
}

// NewTestCase creates a TelescopeTestCase with recording already started. The
// in-memory repository is empty and ready for assertions.
func NewTestCase(t *testing.T) *TelescopeTestCase {
	t.Helper()

	repo := storage.NewInMemoryRepository()
	scope := telescope.New(telescope.WithRepository(repo))
	scope.StartRecording()

	return &TelescopeTestCase{
		T:          t,
		Telescope:  scope,
		Repository: repo,
	}
}

// Store flushes queued entries to the repository, mirroring
// FeatureTestCase::terminateTelescope().
func (tc *TelescopeTestCase) Store() {
	tc.T.Helper()

	if err := tc.Telescope.Store(context.Background()); err != nil {
		tc.T.Fatalf("telescope: store failed: %v", err)
	}
}

// Flush discards queued entries without persisting them.
func (tc *TelescopeTestCase) Flush() {
	tc.Telescope.Flush()
}

// StartRecording enables recording.
func (tc *TelescopeTestCase) StartRecording() {
	tc.Telescope.StartRecording()
}

// StopRecording disables recording.
func (tc *TelescopeTestCase) StopRecording() {
	tc.Telescope.StopRecording()
}

// LoadEntries stores queued entries and returns all persisted results,
// mirroring FeatureTestCase::loadTelescopeEntries().
func (tc *TelescopeTestCase) LoadEntries() []*telescope.EntryResult {
	tc.T.Helper()
	tc.Store()

	results, err := tc.Repository.Get("", telescope.DefaultQueryOptions().WithLimit(500))

	if err != nil {
		tc.T.Fatalf("telescope: load entries failed: %v", err)
	}

	return results
}

// LoadEntriesOfType stores queued entries and returns persisted results of the
// given entry type.
func (tc *TelescopeTestCase) LoadEntriesOfType(entryType string) []*telescope.EntryResult {
	tc.T.Helper()
	tc.Store()

	results, err := tc.Repository.Get(entryType, telescope.DefaultQueryOptions().WithLimit(500))

	if err != nil {
		tc.T.Fatalf("telescope: load entries of type %q failed: %v", entryType, err)
	}

	return results
}

// AssertRecorded asserts that at least one entry of the given type was
// recorded.
func (tc *TelescopeTestCase) AssertRecorded(entryType string) []*telescope.EntryResult {
	tc.T.Helper()

	entries := tc.LoadEntriesOfType(entryType)

	if len(entries) == 0 {
		tc.T.Errorf("expected at least one %q entry to be recorded, but none were found", entryType)
	}

	return entries
}

// AssertNotRecorded asserts that no entry of the given type was recorded.
func (tc *TelescopeTestCase) AssertNotRecorded(entryType string) {
	tc.T.Helper()

	entries := tc.LoadEntriesOfType(entryType)

	if len(entries) > 0 {
		tc.T.Errorf("expected no %q entries, but %d were found", entryType, len(entries))
	}
}

// AssertEntryCount asserts the number of entries of the given type.
func (tc *TelescopeTestCase) AssertEntryCount(entryType string, want int) {
	tc.T.Helper()

	entries := tc.LoadEntriesOfType(entryType)

	if len(entries) != want {
		tc.T.Errorf("expected %d %q entries, got %d", want, entryType, len(entries))
	}
}

// Reset clears all stored entries and restarts recording with a new batch.
func (tc *TelescopeTestCase) Reset() {
	tc.T.Helper()

	if err := tc.Repository.Clear(); err != nil {
		tc.T.Fatalf("telescope: clear repository failed: %v", err)
	}

	tc.Telescope.Flush()
	tc.Telescope.StartRecording()
}
