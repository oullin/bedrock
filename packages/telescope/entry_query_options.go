package telescope

// EntryQueryOptions carries filter parameters for repository queries. It
type EntryQueryOptions struct {
	BatchID        string
	Tag            string
	FamilyHash     string
	BeforeSequence int64
	UUIDs          []string
	Limit          int
}

const defaultQueryLimit = 50

// DefaultQueryOptions returns an EntryQueryOptions with the default query
// limit (50), matching the upstream default.
func DefaultQueryOptions() EntryQueryOptions {
	return EntryQueryOptions{Limit: defaultQueryLimit}
}

// ForBatchID returns a new EntryQueryOptions scoped to a specific batch,
func (o EntryQueryOptions) ForBatchID(id string) EntryQueryOptions {
	o.BatchID = id

	return o
}

// WithTag returns a new EntryQueryOptions filtered to a specific tag.
func (o EntryQueryOptions) WithTag(tag string) EntryQueryOptions {
	o.Tag = tag

	return o
}

// WithFamilyHash returns a new EntryQueryOptions filtered by family hash.
func (o EntryQueryOptions) WithFamilyHash(hash string) EntryQueryOptions {
	o.FamilyHash = hash

	return o
}

// Before returns a new EntryQueryOptions paginating before the given sequence.
func (o EntryQueryOptions) Before(sequence int64) EntryQueryOptions {
	o.BeforeSequence = sequence

	return o
}

// WithLimit returns a new EntryQueryOptions with the specified result limit.
func (o EntryQueryOptions) WithLimit(n int) EntryQueryOptions {
	o.Limit = n

	return o
}

// WithUUIDs returns a new EntryQueryOptions filtered to specific UUIDs.
func (o EntryQueryOptions) WithUUIDs(uuids []string) EntryQueryOptions {
	o.UUIDs = uuids

	return o
}
