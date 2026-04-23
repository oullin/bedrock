package cache_test

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/bedrock/packages/cache"
)

// fakeClock is a controllable clock for deterministic tests.
type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

// spyStore wraps an ArrayStore and records all method calls.
type spyStore struct {
	inner *cache.ArrayStore
	mu    sync.Mutex
	calls []string
}

type spyTaggableStore struct {
	*spyStore
	tagged *spyTaggedCache
}

type spyTaggedCache struct {
	*spyStore
}

// mockEventDispatcher records dispatched events.
type mockEventDispatcher struct {
	mu     sync.Mutex
	events []cache.Event
}

// mockRedisClient implements cache.RedisClient for testing.
type mockRedisClient struct {
	mu   sync.Mutex
	data map[string]string
}

// Simple implementation: parse existing value or start from 0.

// Try to parse as number from JSON.

// Simplified

// mockDBConnection implements cache.DBConnection for testing.
type mockDBConnection struct {
	mu   sync.Mutex
	data map[string]mockDBEntry
}

type mockDBEntry struct {
	value      string
	expiration int64
}

// Flush.

// Delete operation.

// Insert/update.

// mockDBRow implements cache.DBRow for testing.
type mockDBRow struct {
	value      string
	expiration int64
	err        error
}

// mockDynamoClient implements cache.DynamoClient for testing.
type mockDynamoClient struct {
	mu   sync.Mutex
	data map[string]map[string]any
}

// mockSession implements cache.Session for testing.
type mockSession struct {
	mu   sync.Mutex
	data map[string]any
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()

	defer c.mu.Unlock()

	return c.now
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.now = c.now.Add(d)
}

func newSpyStore() *spyStore {
	return &spyStore{inner: cache.NewArrayStore()}
}

func newSpyTaggableStore() *spyTaggableStore {
	return &spyTaggableStore{
		spyStore: newSpyStore(),
		tagged:   &spyTaggedCache{spyStore: newSpyStore()},
	}
}

func (s *spyStore) record(name string) {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.calls = append(s.calls, name)
}

func (s *spyStore) callCount(name string) int {
	s.mu.Lock()

	defer s.mu.Unlock()

	count := 0

	for _, c := range s.calls {
		if c == name {
			count++
		}
	}

	return count
}

func (s *spyStore) Get(ctx context.Context, key string) (any, error) {
	s.record("Get")

	return s.inner.Get(ctx, key)
}

func (s *spyStore) GetMany(ctx context.Context, keys []string) (map[string]any, error) {
	s.record("GetMany")

	return s.inner.GetMany(ctx, keys)
}

func (s *spyStore) Put(ctx context.Context, key string, value any, ttl time.Duration) error {
	s.record("Put")

	return s.inner.Put(ctx, key, value, ttl)
}

func (s *spyStore) PutMany(ctx context.Context, values map[string]any, ttl time.Duration) error {
	s.record("PutMany")

	return s.inner.PutMany(ctx, values, ttl)
}

func (s *spyStore) Add(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	s.record("Add")

	return s.inner.Add(ctx, key, value, ttl)
}

func (s *spyStore) Forever(ctx context.Context, key string, value any) error {
	s.record("Forever")

	return s.inner.Forever(ctx, key, value)
}

func (s *spyStore) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	s.record("Increment")

	return s.inner.Increment(ctx, key, delta)
}

func (s *spyStore) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	s.record("Decrement")

	return s.inner.Decrement(ctx, key, delta)
}

func (s *spyStore) Touch(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	s.record("Touch")

	return s.inner.Touch(ctx, key, ttl)
}

func (s *spyStore) Forget(ctx context.Context, key string) error {
	s.record("Forget")

	return s.inner.Forget(ctx, key)
}

func (s *spyStore) Flush(ctx context.Context) error {
	s.record("Flush")

	return s.inner.Flush(ctx)
}

func (s *spyStore) GetPrefix() string { return s.inner.GetPrefix() }

func (s *spyTaggableStore) Tags(_ ...string) cache.TaggedCache {
	s.record("Tags")

	return s.tagged
}

func (tc *spyTaggedCache) FlushTagged(context.Context) error {
	tc.record("FlushTagged")

	return nil
}

func (tc *spyTaggedCache) TaggedItemKey(_ context.Context, key string) (string, error) {
	return key, nil
}

func (tc *spyTaggedCache) GetTags() *cache.TagSet { return nil }

func (d *mockEventDispatcher) Dispatch(_ context.Context, event cache.Event) {
	d.mu.Lock()

	defer d.mu.Unlock()

	d.events = append(d.events, event)
}

func (d *mockEventDispatcher) Events() []cache.Event {
	d.mu.Lock()

	defer d.mu.Unlock()

	cp := make([]cache.Event, len(d.events))
	copy(cp, d.events)

	return cp
}

func (d *mockEventDispatcher) count(eventType string) int {
	d.mu.Lock()

	defer d.mu.Unlock()

	n := 0

	for _, e := range d.events {
		switch e.(type) {
		case cache.CacheHit:
			if eventType == "CacheHit" {
				n++
			}
		case cache.CacheMissed:
			if eventType == "CacheMissed" {
				n++
			}
		case cache.RetrievingKey:
			if eventType == "RetrievingKey" {
				n++
			}
		case cache.RetrievingManyKeys:
			if eventType == "RetrievingManyKeys" {
				n++
			}
		case cache.WritingKey:
			if eventType == "WritingKey" {
				n++
			}
		case cache.WritingManyKeys:
			if eventType == "WritingManyKeys" {
				n++
			}
		case cache.KeyWritten:
			if eventType == "KeyWritten" {
				n++
			}
		case cache.ForgettingKey:
			if eventType == "ForgettingKey" {
				n++
			}
		case cache.KeyForgotten:
			if eventType == "KeyForgotten" {
				n++
			}
		case cache.KeyForgetFailed:
			if eventType == "KeyForgetFailed" {
				n++
			}
		case cache.CacheFlushing:
			if eventType == "CacheFlushing" {
				n++
			}
		case cache.CacheFlushed:
			if eventType == "CacheFlushed" {
				n++
			}
		case cache.CacheFlushFailed:
			if eventType == "CacheFlushFailed" {
				n++
			}
		case cache.CacheLocksFlushing:
			if eventType == "CacheLocksFlushing" {
				n++
			}
		case cache.CacheLocksFlushed:
			if eventType == "CacheLocksFlushed" {
				n++
			}
		case cache.CacheLocksFlushFailed:
			if eventType == "CacheLocksFlushFailed" {
				n++
			}
		}
	}

	return n
}

func newMockRedisClient() *mockRedisClient {
	return &mockRedisClient{data: make(map[string]string)}
}

func (c *mockRedisClient) Get(_ context.Context, key string) (string, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	v, ok := c.data[key]

	if !ok {
		return "", errors.New("redis: nil")
	}

	return v, nil
}

func (c *mockRedisClient) MGet(_ context.Context, keys ...string) ([]any, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	result := make([]any, len(keys))

	for i, k := range keys {
		if v, ok := c.data[k]; ok {
			result[i] = v
		}
	}

	return result, nil
}

func (c *mockRedisClient) Set(_ context.Context, key string, value any, _ time.Duration) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.data[key] = value.(string)

	return nil
}

func (c *mockRedisClient) SetNX(_ context.Context, key string, value any, _ time.Duration) (bool, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	if _, ok := c.data[key]; ok {
		return false, nil
	}

	c.data[key] = value.(string)

	return true, nil
}

func (c *mockRedisClient) Del(_ context.Context, keys ...string) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	for _, k := range keys {
		delete(c.data, k)
	}

	return nil
}

func (c *mockRedisClient) FlushDB(_ context.Context) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.data = make(map[string]string)

	return nil
}

func (c *mockRedisClient) Incr(_ context.Context, key string) (int64, error) {
	return c.IncrBy(nil, key, 1)
}

func (c *mockRedisClient) IncrBy(_ context.Context, key string, value int64) (int64, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	var current int64

	if v, ok := c.data[key]; ok {

		for _, ch := range v {
			if ch >= '0' && ch <= '9' {
				current = current*10 + int64(ch-'0')
			}
		}
	}

	result := current + value
	c.data[key] = string(rune('0' + result))

	return result, nil
}

func (c *mockRedisClient) Expire(_ context.Context, _ string, _ time.Duration) (bool, error) {
	return true, nil
}

func (c *mockRedisClient) Eval(_ context.Context, _ string, _ []string, _ ...any) (any, error) {
	return int64(1), nil
}

func newMockDBConnection() *mockDBConnection {
	return &mockDBConnection{data: make(map[string]mockDBEntry)}
}

func (c *mockDBConnection) QueryRow(_ context.Context, _ string, args ...any) cache.DBRow {
	c.mu.Lock()

	defer c.mu.Unlock()

	key := args[0].(string)
	entry, ok := c.data[key]

	if !ok {
		return &mockDBRow{err: errors.New("not found")}
	}

	return &mockDBRow{value: entry.value, expiration: entry.expiration}
}

func (c *mockDBConnection) Exec(_ context.Context, query string, args ...any) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	if len(args) == 0 {

		c.data = make(map[string]mockDBEntry)

		return nil
	}

	key := args[0].(string)

	if len(args) == 1 {
		delete(c.data, key)

		return nil
	}

	if len(args) >= 3 {
		value := args[1].(string)

		var exp int64

		switch e := args[2].(type) {
		case int64:
			exp = e
		case int:
			exp = int64(e)
		}

		c.data[key] = mockDBEntry{value: value, expiration: exp}
	}

	return nil
}

func (r *mockDBRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}

	if len(dest) >= 1 {
		if p, ok := dest[0].(*string); ok {
			*p = r.value
		}
	}

	if len(dest) >= 2 {
		if p, ok := dest[1].(*int64); ok {
			*p = r.expiration
		}
	}

	return nil
}

func newMockDynamoClient() *mockDynamoClient {
	return &mockDynamoClient{data: make(map[string]map[string]any)}
}

func (c *mockDynamoClient) GetItem(_ context.Context, _ string, keyAttr, key string) (map[string]any, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	item, ok := c.data[key]

	if !ok {
		return nil, errors.New("not found")
	}

	return item, nil
}

func (c *mockDynamoClient) PutItem(_ context.Context, _ string, item map[string]any) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	key, ok := item["key"].(string)

	if !ok {
		return errors.New("missing key")
	}

	c.data[key] = item

	return nil
}

func (c *mockDynamoClient) DeleteItem(_ context.Context, _ string, _, key string) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	delete(c.data, key)

	return nil
}

func (c *mockDynamoClient) Scan(_ context.Context, _ string) ([]map[string]any, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	var items []map[string]any

	for _, item := range c.data {
		items = append(items, item)
	}

	return items, nil
}

func newMockSession() *mockSession {
	return &mockSession{data: make(map[string]any)}
}

func (s *mockSession) Get(key string) (any, bool) {
	s.mu.Lock()

	defer s.mu.Unlock()

	v, ok := s.data[key]

	return v, ok
}

func (s *mockSession) Put(key string, value any) {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.data[key] = value
}

func (s *mockSession) Forget(key string) {
	s.mu.Lock()

	defer s.mu.Unlock()

	delete(s.data, key)
}

func (s *mockSession) Flush() {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.data = make(map[string]any)
}
