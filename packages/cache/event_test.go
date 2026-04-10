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

func TestCacheFlushingImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.CacheFlushing{StoreName: "array"}

	e.CacheEvent()
}

func TestCacheFlushedImplementsEvent(t *testing.T) {
	t.Parallel()

	var e cache.Event = cache.CacheFlushed{StoreName: "array"}

	e.CacheEvent()
}
