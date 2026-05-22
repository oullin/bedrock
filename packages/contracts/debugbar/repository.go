package debugbar

import "time"

// IncomingEntry is the minimal entry shape the repository contract operates on.
// The concrete type lives in the debugbar package; this alias keeps the
// contracts package free of circular dependencies.
type IncomingEntry = struct {
	UUID       string
	BatchID    string
	Type       string
	FamilyHash string
	Content    map[string]any
	Tags       []string
	RecordedAt time.Time
	Hostname   string
}

// EntryUpdate describes a mutation to apply to a stored entry.
type EntryUpdate = struct {
	UUID    string
	Type    string
	Changes map[string]any
	Tags    struct {
		Add    []string
		Remove []string
	}
}

// EntryResult is the serialised form of a stored entry returned by queries.
type EntryResult = struct {
	ID         string
	Sequence   int64
	BatchID    string
	Type       string
	FamilyHash string
	Content    map[string]any
	CreatedAt  time.Time
	Tags       []string
	Avatar     string
}

// EntryQueryOptions carries filter parameters for repository queries.
type EntryQueryOptions = struct {
	BatchID        string
	Tag            string
	FamilyHash     string
	BeforeSequence int64
	UUIDs          []string
	Limit          int
}

// EntriesRepository defines the storage contract for DebugBar entries.
type EntriesRepository interface {
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
}

// ClearableRepository extends EntriesRepository with a clear operation.
type ClearableRepository interface {
	EntriesRepository
	// Clear removes all stored entries.
	Clear() error
}

// PrunableRepository extends EntriesRepository with a prune operation.
type PrunableRepository interface {
	EntriesRepository
	// Prune removes entries older than before. When keepExceptions is true
	// handled exceptions are retained.
	Prune(before time.Time, keepExceptions bool) (int64, error)
}
