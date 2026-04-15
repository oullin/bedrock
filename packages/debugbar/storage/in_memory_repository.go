// Package storage provides DebugBar repository implementations.
package storage

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/bedrock/packages/debugbar"
)

// ErrNotFound is returned when no entry matches the requested UUID.
var ErrNotFound = errors.New("debugbar: entry not found")

// storedEntry pairs an IncomingEntry with its auto-assigned sequence number.
type storedEntry struct {
	entry    *debugbar.IncomingEntry
	sequence int64
}

// InMemoryRepository is a thread-safe, in-memory DebugBar repository intended
// for use in tests. It mirrors the behaviour of DatabaseEntriesRepository
// without requiring a database connection.
type InMemoryRepository struct {
	mu             sync.RWMutex
	entries        []*storedEntry
	monitoredTags  map[string]struct{}
	nextSequence   int64
}

// NewInMemoryRepository creates an empty InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		monitoredTags: make(map[string]struct{}),
		nextSequence:  1,
	}
}

// compile-time interface satisfaction check.
var _ debugbar.Repository = (*InMemoryRepository)(nil)

// Find retrieves a single entry by UUID.
func (r *InMemoryRepository) Find(id string) (*debugbar.EntryResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, se := range r.entries {
		if se.entry.UUID == id {
			return r.toResult(se), nil
		}
	}

	return nil, ErrNotFound
}

// Get retrieves entries of the given type filtered by EntryQueryOptions.
func (r *InMemoryRepository) Get(entryType string, opts debugbar.EntryQueryOptions) ([]*debugbar.EntryResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	limit := opts.Limit
	if limit <= 0 {
		limit = 50
	}

	var results []*debugbar.EntryResult

	// Iterate in reverse (newest first).
	for i := len(r.entries) - 1; i >= 0; i-- {
		se := r.entries[i]
		e := se.entry

		if entryType != "" && e.Type != entryType {
			continue
		}

		if opts.BatchID != "" && e.BatchID != opts.BatchID {
			continue
		}

		if opts.FamilyHash != "" && e.FamilyHash != opts.FamilyHash {
			continue
		}

		if opts.BeforeSequence > 0 && se.sequence >= opts.BeforeSequence {
			continue
		}

		if opts.Tag != "" && !hasTag(e.Tags, opts.Tag) {
			continue
		}

		if len(opts.UUIDs) > 0 && !containsStr(opts.UUIDs, e.UUID) {
			continue
		}

		results = append(results, r.toResult(se))

		if len(results) >= limit {
			break
		}
	}

	return results, nil
}

// Store persists a batch of incoming entries.
func (r *InMemoryRepository) Store(entries []*debugbar.IncomingEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, e := range entries {
		r.entries = append(r.entries, &storedEntry{
			entry:    e,
			sequence: r.nextSequence,
		})
		r.nextSequence++
	}

	return nil
}

// Update applies field mutations to stored entries.
func (r *InMemoryRepository) Update(updates []*debugbar.EntryUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	index := make(map[string]*storedEntry, len(r.entries))

	for _, se := range r.entries {
		index[se.entry.UUID] = se
	}

	for _, u := range updates {
		se, ok := index[u.UUID]
		if !ok {
			continue
		}

		for k, v := range u.Changes {
			se.entry.Content[k] = v
		}

		// Add tags.
		existing := make(map[string]struct{}, len(se.entry.Tags))

		for _, t := range se.entry.Tags {
			existing[t] = struct{}{}
		}

		for _, t := range u.Tags.Add {
			if _, ok := existing[t]; !ok {
				se.entry.Tags = append(se.entry.Tags, t)
				existing[t] = struct{}{}
			}
		}

		// Remove tags.
		remove := make(map[string]struct{}, len(u.Tags.Remove))

		for _, t := range u.Tags.Remove {
			remove[t] = struct{}{}
		}

		filtered := se.entry.Tags[:0]

		for _, t := range se.entry.Tags {
			if _, ok := remove[t]; !ok {
				filtered = append(filtered, t)
			}
		}

		se.entry.Tags = filtered
	}

	return nil
}

// LoadMonitoredTags is a no-op for the in-memory repository; monitored tags
// are managed entirely in memory.
func (r *InMemoryRepository) LoadMonitoredTags() error { return nil }

// IsMonitoring reports whether any of the given tags are monitored.
func (r *InMemoryRepository) IsMonitoring(tags []string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, t := range tags {
		if _, ok := r.monitoredTags[t]; ok {
			return true
		}
	}

	return false
}

// Monitoring returns the list of currently monitored tags.
func (r *InMemoryRepository) Monitoring() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]string, 0, len(r.monitoredTags))

	for t := range r.monitoredTags {
		out = append(out, t)
	}

	sort.Strings(out)

	return out
}

// Monitor activates monitoring for the given tags.
func (r *InMemoryRepository) Monitor(tags []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, t := range tags {
		r.monitoredTags[t] = struct{}{}
	}

	return nil
}

// StopMonitoring deactivates monitoring for the given tags.
func (r *InMemoryRepository) StopMonitoring(tags []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, t := range tags {
		delete(r.monitoredTags, t)
	}

	return nil
}

// Clear removes all stored entries.
func (r *InMemoryRepository) Clear() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.entries = nil
	r.nextSequence = 1

	return nil
}

// Prune removes entries older than before. When keepExceptions is true,
// exception entries are retained regardless of age.
func (r *InMemoryRepository) Prune(before time.Time, keepExceptions bool) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var kept []*storedEntry
	var pruned int64

	for _, se := range r.entries {
		if se.entry.RecordedAt.Before(before) {
			if keepExceptions && se.entry.Type == debugbar.EntryTypeException {
				kept = append(kept, se)
				continue
			}

			pruned++

			continue
		}

		kept = append(kept, se)
	}

	r.entries = kept

	return pruned, nil
}

// Entries returns a snapshot of all stored entries (for test assertions).
func (r *InMemoryRepository) Entries() []*debugbar.IncomingEntry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]*debugbar.IncomingEntry, len(r.entries))

	for i, se := range r.entries {
		out[i] = se.entry
	}

	return out
}

// Count returns the total number of stored entries.
func (r *InMemoryRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.entries)
}

// toResult converts a storedEntry to an EntryResult.
func (r *InMemoryRepository) toResult(se *storedEntry) *debugbar.EntryResult {
	e := se.entry

	tags := make([]string, len(e.Tags))
	copy(tags, e.Tags)

	content := make(map[string]any, len(e.Content))

	for k, v := range e.Content {
		content[k] = v
	}

	return &debugbar.EntryResult{
		ID:         e.UUID,
		Sequence:   se.sequence,
		BatchID:    e.BatchID,
		Type:       e.Type,
		FamilyHash: e.FamilyHash,
		Content:    content,
		CreatedAt:  e.RecordedAt,
		Tags:       tags,
	}
}

// hasTag reports whether the tag slice contains the target tag.
func hasTag(tags []string, target string) bool {
	for _, t := range tags {
		if t == target {
			return true
		}
	}

	return false
}

// containsStr reports whether the slice contains s.
func containsStr(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}

	return false
}
