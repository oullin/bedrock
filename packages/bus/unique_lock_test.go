package bus_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/bus"
)

func TestUniqueLockAcquireSuccess(t *testing.T) {
	cache := newMockCacheStore()
	lock := bus.NewUniqueLock(cache)

	if !lock.Acquire(context.Background(), "my-job", 60*time.Second) {
		t.Error("expected Acquire to succeed on empty cache")
	}

	// Verify cache.Put was called with correct key and TTL.
	found := false

	for _, call := range cache.calls {
		if call.Method == "Put" && call.Key == "unique_lock:my-job" && call.TTL == 60 {
			found = true

			break
		}
	}

	if !found {
		t.Error("expected cache.Put to be called with key 'unique_lock:my-job' and TTL 60")
	}
}

func TestUniqueLockAcquireAlreadyHeld(t *testing.T) {
	cache := newMockCacheStore()
	cache.data["unique_lock:my-job"] = "1"

	lock := bus.NewUniqueLock(cache)

	if lock.Acquire(context.Background(), "my-job", 60*time.Second) {
		t.Error("expected Acquire to fail when lock is already held")
	}
}

func TestUniqueLockAcquireDefaultTTL(t *testing.T) {
	cache := newMockCacheStore()
	lock := bus.NewUniqueLock(cache)

	lock.Acquire(context.Background(), "key", 0)

	for _, call := range cache.calls {
		if call.Method == "Put" && call.TTL != 3600 {
			t.Errorf("expected default TTL 3600, got %d", call.TTL)
		}
	}
}

func TestUniqueLockRelease(t *testing.T) {
	cache := newMockCacheStore()
	cache.data["unique_lock:my-job"] = "1"

	lock := bus.NewUniqueLock(cache)

	if err := lock.Release(context.Background(), "my-job"); err != nil {
		t.Fatal(err)
	}

	found := false

	for _, call := range cache.calls {
		if call.Method == "Forget" && call.Key == "unique_lock:my-job" {
			found = true

			break
		}
	}

	if !found {
		t.Error("expected cache.Forget to be called with key 'unique_lock:my-job'")
	}
}

func TestUniqueLockKeyPrefix(t *testing.T) {
	cache := newMockCacheStore()
	lock := bus.NewUniqueLock(cache)

	lock.Acquire(context.Background(), "foo", 30*time.Second)

	for _, call := range cache.calls {
		if call.Method == "Get" && call.Key != "unique_lock:foo" {
			t.Errorf("expected key 'unique_lock:foo', got %q", call.Key)
		}

		if call.Method == "Put" && call.Key != "unique_lock:foo" {
			t.Errorf("expected key 'unique_lock:foo', got %q", call.Key)
		}
	}
}
