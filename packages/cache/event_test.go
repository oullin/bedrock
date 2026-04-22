package cache_test

import (
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func TestCacheHitImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.CacheHit{
		StoreName: "array",
		Key:       "users",
		Value:     "data",
		Tags:      []string{"tag1"},
	}

	e.CacheEvent()
}

func TestCacheMissedImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.CacheMissed{
		StoreName: "array",
		Key:       "users",
	}

	e.CacheEvent()
}

func TestRetrievingKeyImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.RetrievingKey{StoreName: "array", Key: "users"}

	e.CacheEvent()
}

func TestRetrievingManyKeysImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.RetrievingManyKeys{StoreName: "array", Keys: []string{"users", "teams"}}

	e.CacheEvent()
}

func TestWritingKeyImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.WritingKey{
		StoreName: "array",
		Key:       "users",
		Value:     "data",
		TTL:       time.Minute,
	}

	e.CacheEvent()
}

func TestWritingManyKeysImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.WritingManyKeys{
		StoreName: "array",
		Keys:      []string{"users"},
		Values:    map[string]any{"users": "data"},
		TTL:       time.Minute,
	}

	e.CacheEvent()
}

func TestKeyWrittenImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.KeyWritten{
		StoreName: "array",
		Key:       "users",
		Value:     "data",
		TTL:       time.Minute,
	}

	e.CacheEvent()
}

func TestForgettingKeyImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.ForgettingKey{StoreName: "array", Key: "users"}

	e.CacheEvent()
}

func TestKeyForgottenImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.KeyForgotten{StoreName: "array", Key: "users"}

	e.CacheEvent()
}

func TestKeyForgetFailedImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.KeyForgetFailed{StoreName: "array", Key: "users"}

	e.CacheEvent()
}

func TestCacheFlushingImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.CacheFlushing{StoreName: "array"}

	e.CacheEvent()
}

func TestCacheFlushFailedImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.CacheFlushFailed{StoreName: "array"}

	e.CacheEvent()
}

func TestCacheLocksFlushingImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.CacheLocksFlushing{StoreName: "array"}

	e.CacheEvent()
}

func TestCacheLocksFlushedImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.CacheLocksFlushed{StoreName: "array"}

	e.CacheEvent()
}

func TestCacheLocksFlushFailedImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.CacheLocksFlushFailed{StoreName: "array"}

	e.CacheEvent()
}

func TestCacheFlushedImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.CacheFlushed{StoreName: "array"}

	e.CacheEvent()
}
