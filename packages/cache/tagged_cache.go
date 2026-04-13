package cache

import (
	"context"
	"fmt"
	"time"
)

// taggedCache wraps a Store with a TagSet, prefixing all keys with the tag
// namespace. Flushing a taggedCache resets the tag IDs, which effectively
// orphans all previously tagged keys without physically deleting them.
type taggedCache struct {
	store Store
	tags  *TagSet
}

var _ TaggedCache = (*taggedCache)(nil)

// NewTaggedCache creates a TaggedCache wrapping the given store and tags.
func NewTaggedCache(store Store, tags *TagSet) TaggedCache {
	return &taggedCache{store: store, tags: tags}
}

func (tc *taggedCache) prefixed(ctx context.Context, key string) (string, error) {
	ns, err := tc.tags.Namespace(ctx)

	if err != nil {
		return "", err
	}

	return ns + ":" + key, nil
}

func (tc *taggedCache) Get(ctx context.Context, key string) (any, error) {
	pk, err := tc.prefixed(ctx, key)

	if err != nil {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
	}

	return tc.store.Get(ctx, pk)
}

func (tc *taggedCache) GetMany(ctx context.Context, keys []string) (map[string]any, error) {
	prefixed := make([]string, len(keys))

	for i, k := range keys {
		pk, err := tc.prefixed(ctx, k)

		if err != nil {
			return nil, err
		}

		prefixed[i] = pk
	}

	result, err := tc.store.GetMany(ctx, prefixed)

	if err != nil {
		return nil, err
	}

	ns, _ := tc.tags.Namespace(ctx)
	out := make(map[string]any, len(result))

	for pk, v := range result {
		original := pk[len(ns)+1:]
		out[original] = v
	}

	return out, nil
}

func (tc *taggedCache) Put(ctx context.Context, key string, value any, ttl time.Duration) error {
	pk, err := tc.prefixed(ctx, key)

	if err != nil {
		return err
	}

	return tc.store.Put(ctx, pk, value, ttl)
}

func (tc *taggedCache) PutMany(ctx context.Context, values map[string]any, ttl time.Duration) error {
	prefixed := make(map[string]any, len(values))

	for k, v := range values {
		pk, err := tc.prefixed(ctx, k)

		if err != nil {
			return err
		}

		prefixed[pk] = v
	}

	return tc.store.PutMany(ctx, prefixed, ttl)
}

func (tc *taggedCache) Add(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	pk, err := tc.prefixed(ctx, key)

	if err != nil {
		return false, err
	}

	return tc.store.Add(ctx, pk, value, ttl)
}

func (tc *taggedCache) Forever(ctx context.Context, key string, value any) error {
	pk, err := tc.prefixed(ctx, key)

	if err != nil {
		return err
	}

	return tc.store.Forever(ctx, pk, value)
}

func (tc *taggedCache) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	pk, err := tc.prefixed(ctx, key)

	if err != nil {
		return 0, err
	}

	return tc.store.Increment(ctx, pk, delta)
}

func (tc *taggedCache) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	pk, err := tc.prefixed(ctx, key)

	if err != nil {
		return 0, err
	}

	return tc.store.Decrement(ctx, pk, delta)
}

func (tc *taggedCache) Touch(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	pk, err := tc.prefixed(ctx, key)

	if err != nil {
		return false, err
	}

	return tc.store.Touch(ctx, pk, ttl)
}

func (tc *taggedCache) Forget(ctx context.Context, key string) error {
	pk, err := tc.prefixed(ctx, key)

	if err != nil {
		return err
	}

	return tc.store.Forget(ctx, pk)
}

func (tc *taggedCache) Flush(ctx context.Context) error {
	return tc.FlushTagged(ctx)
}

func (tc *taggedCache) FlushTagged(ctx context.Context) error {
	return tc.tags.Reset(ctx)
}

func (tc *taggedCache) GetPrefix() string {
	return tc.store.GetPrefix()
}

// TaggedItemKey returns the fully qualified cache key including the tag
// namespace prefix.
func (tc *taggedCache) TaggedItemKey(ctx context.Context, key string) (string, error) {
	return tc.prefixed(ctx, key)
}

// GetTags returns the underlying TagSet.
func (tc *taggedCache) GetTags() *TagSet {
	return tc.tags
}
