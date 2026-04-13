package cache_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func newFileStore(t *testing.T) *cache.FileStore {
	t.Helper()

	dir := t.TempDir()

	return cache.NewFileStore(dir)
}

func TestFileStoreBasicGetPut(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)
	ctx := context.Background()

	_ = s.Put(ctx, "key", "value", time.Minute)

	got, err := s.Get(ctx, "key")

	if err != nil {
		t.Fatal(err)
	}

	if got != "value" {
		t.Fatalf("expected 'value', got %v", got)
	}
}

func TestFileStoreMissingKey(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)

	_, err := s.Get(context.Background(), "missing")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestFileStoreExpiry(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	clk := &fakeClock{now: time.Now()}
	s := cache.NewFileStoreWithOptions(dir, "", 0o644, clk)
	ctx := context.Background()

	_ = s.Put(ctx, "key", "val", 10*time.Second)
	clk.Advance(11 * time.Second)

	_, err := s.Get(ctx, "key")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after expiry, got %v", err)
	}
}

func TestFileStoreNoExpiry(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	clk := &fakeClock{now: time.Now()}
	s := cache.NewFileStoreWithOptions(dir, "", 0o644, clk)
	ctx := context.Background()

	_ = s.Put(ctx, "key", "val", 0)
	clk.Advance(365 * 24 * time.Hour)

	_, err := s.Get(ctx, "key")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFileStoreGetMany(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)
	ctx := context.Background()

	_ = s.Put(ctx, "a", 1, 0)
	_ = s.Put(ctx, "b", 2, 0)

	got, err := s.GetMany(ctx, []string{"a", "b", "c"})

	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFileStorePutMany(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)
	ctx := context.Background()

	_ = s.PutMany(ctx, map[string]any{"x": 10, "y": 20}, time.Minute)

	v, _ := s.Get(ctx, "x")

	if v != 10 {
		t.Fatalf("expected 10, got %v", v)
	}
}

func TestFileStoreAdd(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)
	ctx := context.Background()

	ok, _ := s.Add(ctx, "k", "v1", time.Minute)

	if !ok {
		t.Fatal("expected Add to succeed")
	}

	ok, _ = s.Add(ctx, "k", "v2", time.Minute)

	if ok {
		t.Fatal("expected Add to fail for existing key")
	}
}

func TestFileStoreForever(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)
	ctx := context.Background()

	_ = s.Forever(ctx, "k", "eternal")

	v, err := s.Get(ctx, "k")

	if err != nil {
		t.Fatal(err)
	}

	if v != "eternal" {
		t.Fatalf("expected 'eternal', got %v", v)
	}
}

func TestFileStoreIncrement(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)
	ctx := context.Background()

	v, _ := s.Increment(ctx, "counter", 3)

	if v != 3 {
		t.Fatalf("expected 3, got %d", v)
	}

	v, _ = s.Increment(ctx, "counter", 7)

	if v != 10 {
		t.Fatalf("expected 10, got %d", v)
	}
}

func TestFileStoreDecrement(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)
	ctx := context.Background()

	_, _ = s.Increment(ctx, "k", 10)
	v, _ := s.Decrement(ctx, "k", 3)

	if v != 7 {
		t.Fatalf("expected 7, got %d", v)
	}
}

func TestFileStoreTouch(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	clk := &fakeClock{now: time.Now()}
	s := cache.NewFileStoreWithOptions(dir, "", 0o644, clk)
	ctx := context.Background()

	_ = s.Put(ctx, "k", "v", 10*time.Second)

	ok, _ := s.Touch(ctx, "k", 30*time.Second)

	if !ok {
		t.Fatal("expected Touch to succeed")
	}

	clk.Advance(15 * time.Second)

	_, err := s.Get(ctx, "k")

	if err != nil {
		t.Fatal("expected key to survive after Touch")
	}
}

func TestFileStoreTouchMissing(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)

	ok, _ := s.Touch(context.Background(), "missing", time.Minute)

	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestFileStoreForget(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)
	ctx := context.Background()

	_ = s.Put(ctx, "k", "v", 0)
	_ = s.Forget(ctx, "k")

	_, err := s.Get(ctx, "k")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatal("expected ErrNotFound after Forget")
	}
}

func TestFileStoreForgetMissing(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)

	err := s.Forget(context.Background(), "missing")

	if err != nil {
		t.Fatalf("expected no error for missing key, got %v", err)
	}
}

func TestFileStoreFlush(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)
	ctx := context.Background()

	_ = s.Put(ctx, "a", 1, 0)
	_ = s.Put(ctx, "b", 2, 0)
	_ = s.Flush(ctx)

	_, err := s.Get(ctx, "a")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatal("expected ErrNotFound after Flush")
	}
}

func TestFileStoreGetPrefix(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	s := cache.NewFileStoreWithOptions(dir, "myapp", 0o644, nil)

	if s.GetPrefix() != "myapp" {
		t.Fatalf("expected 'myapp', got %q", s.GetPrefix())
	}
}

func TestFileStoreIncrementNonNumeric(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)
	ctx := context.Background()

	_ = s.Put(ctx, "k", "hello", time.Minute)

	_, err := s.Increment(ctx, "k", 1)

	if !errors.Is(err, cache.ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue, got %v", err)
	}
}

func TestFileStoreLock(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	s := cache.NewFileStore(dir)
	ctx := context.Background()

	l1 := s.Lock("res", "owner-1", time.Minute)
	ok, _ := l1.Acquire(ctx)

	if !ok {
		t.Fatal("expected acquire to succeed")
	}

	l2 := s.Lock("res", "owner-2", time.Minute)
	ok, _ = l2.Acquire(ctx)

	if ok {
		t.Fatal("expected acquire to fail for different owner")
	}

	l1.Release(ctx) //nolint:errcheck

	ok, _ = l2.Acquire(ctx)

	if !ok {
		t.Fatal("expected acquire to succeed after release")
	}

	// Cleanup lock files.
	os.RemoveAll(dir)
}

func TestFileStoreTags(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)
	ctx := context.Background()

	tagged := s.Tags("products")
	_ = tagged.Put(ctx, "item", "widget", time.Minute)

	v, err := tagged.Get(ctx, "item")

	if err != nil {
		t.Fatal(err)
	}

	if v != "widget" {
		t.Fatalf("expected 'widget', got %v", v)
	}
}

func TestFileStoreComplexValues(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)
	ctx := context.Background()

	data := map[string]any{"name": "Alice", "age": 30}
	_ = s.Put(ctx, "user", data, time.Minute)

	v, err := s.Get(ctx, "user")

	if err != nil {
		t.Fatal(err)
	}

	m, ok := v.(map[string]any)

	if !ok {
		t.Fatalf("expected map, got %T", v)
	}

	if m["name"] != "Alice" {
		t.Fatalf("expected 'Alice', got %v", m["name"])
	}
}

func TestFileStorePathPublic(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)

	p := s.Path("mykey")

	if p == "" {
		t.Fatal("expected non-empty path")
	}

	if p == "mykey" {
		t.Fatal("expected hashed path, not raw key")
	}
}

func TestFileStoreGetSetDirectory(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)
	original := s.GetDirectory()

	if original == "" {
		t.Fatal("expected non-empty directory")
	}

	s.SetDirectory("/tmp/new-cache")

	if s.GetDirectory() != "/tmp/new-cache" {
		t.Fatalf("expected '/tmp/new-cache', got %q", s.GetDirectory())
	}
}

func TestFileStoreSetLockDirectory(t *testing.T) {
	t.Parallel()

	s := newFileStore(t)
	s.SetLockDirectory("/tmp/locks")

	ctx := context.Background()
	l := s.Lock("test", "owner", time.Minute)
	ok, _ := l.Acquire(ctx)

	if !ok {
		t.Fatal("expected lock acquire with custom lock directory")
	}

	l.Release(ctx) //nolint:errcheck
}
