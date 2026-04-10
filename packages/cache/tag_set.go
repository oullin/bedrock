package cache

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
)

// TagSet manages a set of cache tags. Each tag is identified by a random ID
// stored in the cache. The namespace is the combination of all tag IDs, which
// is used to prefix cache keys. Resetting a tag invalidates all keys under
// that tag's namespace by generating a new random ID.
type TagSet struct {
	store Store
	names []string
}

// NewTagSet creates a TagSet with the given tag names.
func NewTagSet(store Store, names []string) *TagSet {
	return &TagSet{store: store, names: names}
}

// GetNames returns the tag names.
func (ts *TagSet) GetNames() []string { return ts.names }

// Namespace returns the combined namespace string for all tags. The namespace
// is a pipe-delimited concatenation of all tag IDs.
func (ts *TagSet) Namespace(ctx context.Context) (string, error) {
	ids, err := ts.TagIDs(ctx)

	if err != nil {
		return "", err
	}

	return strings.Join(ids, "|"), nil
}

// TagIDs returns the IDs for all tags.
func (ts *TagSet) TagIDs(ctx context.Context) ([]string, error) {
	ids := make([]string, len(ts.names))

	for i, name := range ts.names {
		id, err := ts.TagID(ctx, name)

		if err != nil {
			return nil, err
		}

		ids[i] = id
	}

	return ids, nil
}

// TagID returns the current ID for a tag. If the tag has no ID yet, a new
// random one is generated and stored.
func (ts *TagSet) TagID(ctx context.Context, name string) (string, error) {
	key := ts.tagKey(name)

	v, err := ts.store.Get(ctx, key)

	if err == nil {
		if id, ok := v.(string); ok {
			return id, nil
		}
	}

	return ts.ResetTag(ctx, name)
}

// Reset regenerates random IDs for all tags, effectively invalidating every
// key stored under this tag set's namespace.
func (ts *TagSet) Reset(ctx context.Context) error {
	for _, name := range ts.names {
		if _, err := ts.ResetTag(ctx, name); err != nil {
			return err
		}
	}

	return nil
}

// ResetTag generates a new random ID for a single tag and stores it.
func (ts *TagSet) ResetTag(ctx context.Context, name string) (string, error) {
	id := randomID()

	if err := ts.store.Forever(ctx, ts.tagKey(name), id); err != nil {
		return "", err
	}

	return id, nil
}

func (ts *TagSet) tagKey(name string) string {
	return "tag:" + name + ":key"
}

// randomID generates a short random hex string.
func randomID() string {
	b := make([]byte, 10)
	_, _ = rand.Read(b)

	return hex.EncodeToString(b)
}
