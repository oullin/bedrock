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

type eventTaggedCache struct {
	TaggedCache
	storeName string
	tags      []string
	events    EventDispatcher
}

var _ TaggedCache = (*taggedCache)(nil)
var _ TaggedCache = (*eventTaggedCache)(nil)

// NewTaggedCache creates a TaggedCache wrapping the given store and tags.
func NewTaggedCache(store Store, tags *TagSet) TaggedCache {
	return &taggedCache{store: store, tags: tags}
}

func newEventTaggedCache(inner TaggedCache, storeName string, tags []string, events EventDispatcher) TaggedCache {
	return &eventTaggedCache{TaggedCache: inner, storeName: storeName, tags: append([]string(nil), tags...), events: events}
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

func (tc *eventTaggedCache) dispatch(ctx context.Context, event Event) {
	if tc.events != nil {
		tc.events.Dispatch(ctx, event)
	}
}

func (tc *eventTaggedCache) Get(ctx context.Context, key string) (any, error) {
	tc.dispatch(ctx, RetrievingKey{StoreName: tc.storeName, Key: key, Tags: tc.tags})

	v, err := tc.TaggedCache.Get(ctx, key)

	if err != nil {
		tc.dispatch(ctx, CacheMissed{StoreName: tc.storeName, Key: key, Tags: tc.tags})

		return nil, err
	}

	tc.dispatch(ctx, CacheHit{StoreName: tc.storeName, Key: key, Value: v, Tags: tc.tags})

	return v, nil
}

func (tc *eventTaggedCache) GetMany(ctx context.Context, keys []string) (map[string]any, error) {
	tc.dispatch(ctx, RetrievingManyKeys{StoreName: tc.storeName, Keys: keys, Tags: tc.tags})

	result, err := tc.TaggedCache.GetMany(ctx, keys)

	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		if v, ok := result[key]; ok {
			tc.dispatch(ctx, CacheHit{StoreName: tc.storeName, Key: key, Value: v, Tags: tc.tags})
		} else {
			tc.dispatch(ctx, CacheMissed{StoreName: tc.storeName, Key: key, Tags: tc.tags})
		}
	}

	return result, nil
}

func (tc *eventTaggedCache) Put(ctx context.Context, key string, value any, ttl time.Duration) error {
	tc.dispatch(ctx, WritingKey{StoreName: tc.storeName, Key: key, Value: value, TTL: ttl, Tags: tc.tags})

	if err := tc.TaggedCache.Put(ctx, key, value, ttl); err != nil {
		return err
	}

	tc.dispatch(ctx, KeyWritten{StoreName: tc.storeName, Key: key, Value: value, TTL: ttl, Tags: tc.tags})

	return nil
}

func (tc *eventTaggedCache) PutMany(ctx context.Context, values map[string]any, ttl time.Duration) error {
	keys := make([]string, 0, len(values))

	for key := range values {
		keys = append(keys, key)
	}

	tc.dispatch(ctx, WritingManyKeys{StoreName: tc.storeName, Keys: keys, Values: values, TTL: ttl, Tags: tc.tags})

	if err := tc.TaggedCache.PutMany(ctx, values, ttl); err != nil {
		return err
	}

	for key, value := range values {
		tc.dispatch(ctx, KeyWritten{StoreName: tc.storeName, Key: key, Value: value, TTL: ttl, Tags: tc.tags})
	}

	return nil
}

func (tc *eventTaggedCache) Add(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	tc.dispatch(ctx, RetrievingKey{StoreName: tc.storeName, Key: key, Tags: tc.tags})

	if existing, err := tc.TaggedCache.Get(ctx, key); err == nil {
		tc.dispatch(ctx, CacheHit{StoreName: tc.storeName, Key: key, Value: existing, Tags: tc.tags})

		return false, nil
	}

	tc.dispatch(ctx, CacheMissed{StoreName: tc.storeName, Key: key, Tags: tc.tags})
	tc.dispatch(ctx, WritingKey{StoreName: tc.storeName, Key: key, Value: value, TTL: ttl, Tags: tc.tags})

	ok, err := tc.TaggedCache.Add(ctx, key, value, ttl)

	if err != nil {
		return ok, err
	}

	if ok {
		tc.dispatch(ctx, KeyWritten{StoreName: tc.storeName, Key: key, Value: value, TTL: ttl, Tags: tc.tags})
	}

	return ok, nil
}

func (tc *eventTaggedCache) Forever(ctx context.Context, key string, value any) error {
	tc.dispatch(ctx, WritingKey{StoreName: tc.storeName, Key: key, Value: value, Tags: tc.tags})

	if err := tc.TaggedCache.Forever(ctx, key, value); err != nil {
		return err
	}

	tc.dispatch(ctx, KeyWritten{StoreName: tc.storeName, Key: key, Value: value, Tags: tc.tags})

	return nil
}

func (tc *eventTaggedCache) Forget(ctx context.Context, key string) error {
	tc.dispatch(ctx, ForgettingKey{StoreName: tc.storeName, Key: key, Tags: tc.tags})

	if err := tc.TaggedCache.Forget(ctx, key); err != nil {
		tc.dispatch(ctx, KeyForgetFailed{StoreName: tc.storeName, Key: key, Err: err, Tags: tc.tags})

		return err
	}

	tc.dispatch(ctx, KeyForgotten{StoreName: tc.storeName, Key: key, Tags: tc.tags})

	return nil
}

func (tc *eventTaggedCache) Flush(ctx context.Context) error {
	tc.dispatch(ctx, CacheFlushing{StoreName: tc.storeName, Tags: tc.tags})

	if err := tc.TaggedCache.Flush(ctx); err != nil {
		tc.dispatch(ctx, CacheFlushFailed{StoreName: tc.storeName, Err: err, Tags: tc.tags})

		return err
	}

	tc.dispatch(ctx, CacheFlushed{StoreName: tc.storeName, Tags: tc.tags})

	return nil
}
