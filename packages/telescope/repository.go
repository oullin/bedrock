package telescope

import "time"

// Repository is the local storage contract used by the Telescope orchestrator.
// It is a superset of the contracts/telescope.EntriesRepository interface,
// adding clear and prune capabilities.
type Repository interface {
	// Find retrieves a single entry by its UUID.
	Find(id string) (*EntryResult, error)
	// Get retrieves entries of the given type using the provided query options.
	Get(entryType string, options EntryQueryOptions) ([]*EntryResult, error)
	// Store persists a batch of incoming entries.
	Store(entries []*IncomingEntry) error
	// Update applies a set of mutations to stored entries.
	Update(updates []*EntryUpdate) error
	// LoadMonitoredTags loads the set of monitored tags from storage.
	LoadMonitoredTags() error
	// IsMonitoring reports whether any of the given tags are being monitored.
	IsMonitoring(tags []string) bool
	// Monitoring returns the list of currently monitored tags.
	Monitoring() []string
	// Monitor activates monitoring for the given tags.
	Monitor(tags []string) error
	// StopMonitoring deactivates monitoring for the given tags.
	StopMonitoring(tags []string) error
	// Clear removes all stored entries.
	Clear() error
	// Prune removes entries older than before.
	Prune(before time.Time, keepExceptions bool) (int64, error)
}
