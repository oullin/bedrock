package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/bedrock/packages/session"
)

type complianceCacheStore struct {
	mu   sync.Mutex
	data map[string]string
	ttls map[string]int
}

func newComplianceCacheStore() *complianceCacheStore {
	return &complianceCacheStore{
		data: make(map[string]string),
		ttls: make(map[string]int),
	}
}

func (c *complianceCacheStore) Get(_ context.Context, key string) (string, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	v, ok := c.data[key]

	if !ok {
		return "", fmt.Errorf("missing key: %s", key)
	}

	return v, nil
}

func (c *complianceCacheStore) Put(_ context.Context, key, value string, ttlSeconds int) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.data[key] = value
	c.ttls[key] = ttlSeconds

	return nil
}

func (c *complianceCacheStore) Forget(_ context.Context, key string) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	delete(c.data, key)
	delete(c.ttls, key)

	return nil
}

// Ports of Upstream framework session handler tests from tests/Session.
func TestFrameworkSessionUpstreamInventoryCoverage(t *testing.T) {
	t.Run("ArraySessionHandlerTest::test_it_implements_the_session_handler_interface", func(t *testing.T) {
		if _, ok := any(NewArrayHandler()).(session.Handler); !ok {
			t.Fatal("expected ArrayHandler to satisfy session.Handler")
		}
	})

	t.Run("ArraySessionHandlerTest::test_it_initializes_the_session", func(t *testing.T) {
		var h ArrayHandler

		if err := h.Open(context.Background(), "", "session"); err != nil {
			t.Fatalf("Open: %v", err)
		}

		if h.sessions == nil {
			t.Fatal("expected sessions map to be initialized")
		}
	})

	t.Run("ArraySessionHandlerTest::test_it_closes_the_session", func(t *testing.T) {
		if err := NewArrayHandler().Close(context.Background()); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})

	t.Run("ArraySessionHandlerTest::test_it_reads_data_from_the_session", func(t *testing.T) {
		h := NewArrayHandler()
		ctx := context.Background()

		if err := h.Write(ctx, "id", "payload"); err != nil {
			t.Fatalf("Write: %v", err)
		}

		got, err := h.Read(ctx, "id")

		if err != nil {
			t.Fatalf("Read: %v", err)
		}

		if got != "payload" {
			t.Fatalf("got %q, want payload", got)
		}
	})

	t.Run("ArraySessionHandlerTest::test_it_reads_data_from_an_almost_expired_session", func(t *testing.T) {
		h := NewArrayHandler()
		h.maxLifetime = 1
		h.sessions["id"] = sessionRecord{data: "payload", writtenAt: time.Now().Add(-500 * time.Millisecond)}

		got, err := h.Read(context.Background(), "id")

		if err != nil {
			t.Fatalf("Read: %v", err)
		}

		if got != "payload" {
			t.Fatalf("got %q, want payload", got)
		}
	})

	t.Run("ArraySessionHandlerTest::test_it_reads_data_from_an_expired_session", func(t *testing.T) {
		h := NewArrayHandler()
		h.maxLifetime = 1
		h.sessions["id"] = sessionRecord{data: "payload", writtenAt: time.Now().Add(-2 * time.Second)}

		got, err := h.Read(context.Background(), "id")

		if err != nil {
			t.Fatalf("Read: %v", err)
		}

		if got != "" {
			t.Fatalf("got %q, want empty string", got)
		}
	})

	t.Run("ArraySessionHandlerTest::test_it_reads_data_from_a_non_existing_session", func(t *testing.T) {
		got, err := NewArrayHandler().Read(context.Background(), "missing")

		if err != nil {
			t.Fatalf("Read: %v", err)
		}

		if got != "" {
			t.Fatalf("got %q, want empty string", got)
		}
	})

	t.Run("ArraySessionHandlerTest::test_it_writes_session_data", func(t *testing.T) {
		h := NewArrayHandler()

		if err := h.Write(context.Background(), "id", "payload"); err != nil {
			t.Fatalf("Write: %v", err)
		}

		got, err := h.Read(context.Background(), "id")

		if err != nil {
			t.Fatalf("Read: %v", err)
		}

		if got != "payload" {
			t.Fatalf("got %q, want payload", got)
		}
	})

	t.Run("ArraySessionHandlerTest::test_it_destroys_a_session", func(t *testing.T) {
		h := NewArrayHandler()
		ctx := context.Background()

		if err := h.Write(ctx, "id", "payload"); err != nil {
			t.Fatalf("Write: %v", err)
		}

		if err := h.Destroy(ctx, "id"); err != nil {
			t.Fatalf("Destroy: %v", err)
		}

		got, err := h.Read(ctx, "id")

		if err != nil {
			t.Fatalf("Read: %v", err)
		}

		if got != "" {
			t.Fatalf("got %q, want empty string", got)
		}
	})

	t.Run("ArraySessionHandlerTest::test_it_cleans_up_old_sessions", func(t *testing.T) {
		h := NewArrayHandler()
		ctx := context.Background()
		h.maxLifetime = 1
		h.sessions["old"] = sessionRecord{data: "old", writtenAt: time.Now().Add(-2 * time.Second)}
		h.sessions["new"] = sessionRecord{data: "new", writtenAt: time.Now()}

		if err := h.GC(ctx, 1); err != nil {
			t.Fatalf("GC: %v", err)
		}

		if got, err := h.Read(ctx, "old"); err != nil {
			t.Fatalf("Read old: %v", err)
		} else if got != "" {
			t.Fatalf("got %q for old, want empty string", got)
		}

		if got, err := h.Read(ctx, "new"); err != nil {
			t.Fatalf("Read new: %v", err)
		} else if got != "new" {
			t.Fatalf("got %q for new, want new", got)
		}
	})

	t.Run("CacheBasedSessionHandlerTest::test_open", func(t *testing.T) {
		if err := NewCacheBasedHandler(newComplianceCacheStore(), 10).Open(context.Background(), "", "session"); err != nil {
			t.Fatalf("Open: %v", err)
		}
	})

	t.Run("CacheBasedSessionHandlerTest::test_close", func(t *testing.T) {
		if err := NewCacheBasedHandler(newComplianceCacheStore(), 10).Close(context.Background()); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})

	t.Run("CacheBasedSessionHandlerTest::test_read_returns_data_from_cache", func(t *testing.T) {
		cache := newComplianceCacheStore()
		cache.data["session:id"] = "payload"
		h := NewCacheBasedHandler(cache, 10)

		got, err := h.Read(context.Background(), "id")

		if err != nil {
			t.Fatalf("Read: %v", err)
		}

		if got != "payload" {
			t.Fatalf("got %q, want payload", got)
		}
	})

	t.Run("CacheBasedSessionHandlerTest::test_read_returns_empty_string_if_no_data", func(t *testing.T) {
		got, err := NewCacheBasedHandler(newComplianceCacheStore(), 10).Read(context.Background(), "missing")

		if err != nil {
			t.Fatalf("Read: %v", err)
		}

		if got != "" {
			t.Fatalf("got %q, want empty string", got)
		}
	})

	t.Run("CacheBasedSessionHandlerTest::test_write_stores_data_in_cache", func(t *testing.T) {
		cache := newComplianceCacheStore()
		h := NewCacheBasedHandler(cache, 10)

		if err := h.Write(context.Background(), "id", "payload"); err != nil {
			t.Fatalf("Write: %v", err)
		}

		if got := cache.data["session:id"]; got != "payload" {
			t.Fatalf("got %q, want payload", got)
		}

		if ttl := cache.ttls["session:id"]; ttl != 600 {
			t.Fatalf("got ttl %d, want 600", ttl)
		}
	})

	t.Run("CacheBasedSessionHandlerTest::test_destroy_removes_data_from_cache", func(t *testing.T) {
		cache := newComplianceCacheStore()
		cache.data["session:id"] = "payload"
		h := NewCacheBasedHandler(cache, 10)

		if err := h.Destroy(context.Background(), "id"); err != nil {
			t.Fatalf("Destroy: %v", err)
		}

		if _, ok := cache.data["session:id"]; ok {
			t.Fatal("expected cache entry to be removed")
		}
	})

	t.Run("CacheBasedSessionHandlerTest::test_gc_returns_zero", func(t *testing.T) {
		if err := NewCacheBasedHandler(newComplianceCacheStore(), 10).GC(context.Background(), 3600); err != nil {
			t.Fatalf("GC: %v", err)
		}
	})

	t.Run("CacheBasedSessionHandlerTest::test_get_cache_returns_cache_instance", func(t *testing.T) {
		cache := newComplianceCacheStore()
		h := NewCacheBasedHandler(cache, 10)

		if got := h.GetCache(); got != cache {
			t.Fatalf("got %p, want %p", got, cache)
		}
	})

	t.Run("FileSessionHandlerTest::test_open", func(t *testing.T) {
		dir := t.TempDir()
		sub := filepath.Join(dir, "sessions")
		h := NewFileHandler(sub)

		if err := h.Open(context.Background(), "", "session"); err != nil {
			t.Fatalf("Open: %v", err)
		}

		if info, err := os.Stat(sub); err != nil {
			t.Fatalf("Stat: %v", err)
		} else if !info.IsDir() {
			t.Fatal("expected sessions path to be a directory")
		}
	})

	t.Run("FileSessionHandlerTest::test_close", func(t *testing.T) {
		if err := NewFileHandler(t.TempDir()).Close(context.Background()); err != nil {
			t.Fatalf("Close: %v", err)
		}
	})

	t.Run("FileSessionHandlerTest::test_read_returns_data_when_file_exists_and_is_valid", func(t *testing.T) {
		dir := t.TempDir()
		h := NewFileHandler(dir)
		h.maxLifetime = 1
		ctx := context.Background()

		if err := h.Open(ctx, "", "session"); err != nil {
			t.Fatalf("Open: %v", err)
		}

		if err := h.Write(ctx, "id", "payload"); err != nil {
			t.Fatalf("Write: %v", err)
		}

		got, err := h.Read(ctx, "id")

		if err != nil {
			t.Fatalf("Read: %v", err)
		}

		if got != "payload" {
			t.Fatalf("got %q, want payload", got)
		}
	})

	t.Run("FileSessionHandlerTest::test_read_returns_data_when_file_exists_but_expired", func(t *testing.T) {
		dir := t.TempDir()
		h := NewFileHandler(dir)
		h.maxLifetime = 1
		ctx := context.Background()

		if err := h.Open(ctx, "", "session"); err != nil {
			t.Fatalf("Open: %v", err)
		}

		if err := h.Write(ctx, "id", "payload"); err != nil {
			t.Fatalf("Write: %v", err)
		}

		past := time.Now().Add(-2 * time.Second)

		if err := os.Chtimes(filepath.Join(dir, "id"), past, past); err != nil {
			t.Fatalf("Chtimes: %v", err)
		}

		got, err := h.Read(ctx, "id")

		if err != nil {
			t.Fatalf("Read: %v", err)
		}

		if got != "" {
			t.Fatalf("got %q, want empty string", got)
		}
	})

	t.Run("FileSessionHandlerTest::test_read_returns_empty_string_when_file_does_not_exist", func(t *testing.T) {
		h := NewFileHandler(t.TempDir())
		ctx := context.Background()

		if err := h.Open(ctx, "", "session"); err != nil {
			t.Fatalf("Open: %v", err)
		}

		got, err := h.Read(ctx, "missing")

		if err != nil {
			t.Fatalf("Read: %v", err)
		}

		if got != "" {
			t.Fatalf("got %q, want empty string", got)
		}
	})

	t.Run("FileSessionHandlerTest::test_write_stores_data", func(t *testing.T) {
		h := NewFileHandler(t.TempDir())
		ctx := context.Background()

		if err := h.Open(ctx, "", "session"); err != nil {
			t.Fatalf("Open: %v", err)
		}

		if err := h.Write(ctx, "id", "payload"); err != nil {
			t.Fatalf("Write: %v", err)
		}

		got, err := h.Read(ctx, "id")

		if err != nil {
			t.Fatalf("Read: %v", err)
		}

		if got != "payload" {
			t.Fatalf("got %q, want payload", got)
		}
	})

	t.Run("FileSessionHandlerTest::test_destroy_deletes_session_file", func(t *testing.T) {
		dir := t.TempDir()
		h := NewFileHandler(dir)
		ctx := context.Background()

		if err := h.Open(ctx, "", "session"); err != nil {
			t.Fatalf("Open: %v", err)
		}

		if err := h.Write(ctx, "id", "payload"); err != nil {
			t.Fatalf("Write: %v", err)
		}

		if err := h.Destroy(ctx, "id"); err != nil {
			t.Fatalf("Destroy: %v", err)
		}

		got, err := h.Read(ctx, "id")

		if err != nil {
			t.Fatalf("Read: %v", err)
		}

		if got != "" {
			t.Fatalf("got %q, want empty string", got)
		}
	})

	t.Run("FileSessionHandlerTest::test_gc_deletes_old_session_files", func(t *testing.T) {
		dir := t.TempDir()
		h := NewFileHandler(dir)
		ctx := context.Background()
		h.maxLifetime = 1

		if err := h.Open(ctx, "", "session"); err != nil {
			t.Fatalf("Open: %v", err)
		}

		if err := h.Write(ctx, "old", "old"); err != nil {
			t.Fatalf("Write old: %v", err)
		}

		if err := h.Write(ctx, "new", "new"); err != nil {
			t.Fatalf("Write new: %v", err)
		}

		past := time.Now().Add(-2 * time.Second)

		if err := os.Chtimes(filepath.Join(dir, "old"), past, past); err != nil {
			t.Fatalf("Chtimes old: %v", err)
		}

		if err := h.GC(ctx, 1); err != nil {
			t.Fatalf("GC: %v", err)
		}

		if got, err := h.Read(ctx, "old"); err != nil {
			t.Fatalf("Read old: %v", err)
		} else if got != "" {
			t.Fatalf("got %q for old, want empty string", got)
		}

		if got, err := h.Read(ctx, "new"); err != nil {
			t.Fatalf("Read new: %v", err)
		} else if got != "new" {
			t.Fatalf("got %q for new, want new", got)
		}
	})
}
