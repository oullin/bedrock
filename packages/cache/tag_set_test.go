package cache_test

import (
	"context"
	"strings"
	"testing"

	"github.com/bedrock/packages/cache"
)

func TestTagSetNamespace(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"users", "posts"})
	ctx := context.Background()

	ns, err := ts.Namespace(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if ns == "" {
		t.Fatal("expected non-empty namespace")
	}

	parts := strings.Split(ns, "|")

	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
}

func TestTagSetNamespaceConsistency(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"users"})
	ctx := context.Background()

	ns1, _ := ts.Namespace(ctx)
	ns2, _ := ts.Namespace(ctx)

	if ns1 != ns2 {
		t.Fatal("expected consistent namespace between calls")
	}
}

func TestTagSetReset(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"users"})
	ctx := context.Background()

	before, _ := ts.Namespace(ctx)
	_ = ts.Reset(ctx)
	after, _ := ts.Namespace(ctx)

	if before == after {
		t.Fatal("expected namespace to change after Reset")
	}
}

func TestTagSetResetTag(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"users", "posts"})
	ctx := context.Background()

	before, _ := ts.TagID(ctx, "users")
	_, _ = ts.ResetTag(ctx, "users")
	after, _ := ts.TagID(ctx, "users")

	if before == after {
		t.Fatal("expected tag ID to change after ResetTag")
	}
}

func TestTagSetGetNames(t *testing.T) {
	t.Parallel()

	ts := cache.NewTagSet(cache.NewArrayStore(), []string{"a", "b", "c"})

	names := ts.GetNames()

	if len(names) != 3 {
		t.Fatalf("expected 3, got %d", len(names))
	}

	if names[0] != "a" || names[1] != "b" || names[2] != "c" {
		t.Fatal("unexpected names")
	}
}

func TestTagSetTagIDs(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"x", "y"})
	ctx := context.Background()

	ids, err := ts.TagIDs(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if len(ids) != 2 {
		t.Fatalf("expected 2 IDs, got %d", len(ids))
	}

	if ids[0] == "" || ids[1] == "" {
		t.Fatal("expected non-empty IDs")
	}
}

func TestTagSetTagIDPersistence(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"cached"})
	ctx := context.Background()

	id1, _ := ts.TagID(ctx, "cached")
	id2, _ := ts.TagID(ctx, "cached")

	if id1 != id2 {
		t.Fatal("expected same tag ID on repeated calls")
	}
}

func TestTagSetFlush(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"a", "b"})
	ctx := context.Background()

	// Generate tag IDs (stored in cache).
	_, _ = ts.TagID(ctx, "a")
	_, _ = ts.TagID(ctx, "b")

	// Verify tag keys exist.
	_, err := store.Get(ctx, ts.TagKey("a"))

	if err != nil {
		t.Fatal("expected tag key 'a' to exist")
	}

	// Flush deletes tag keys entirely.
	if err := ts.Flush(ctx); err != nil {
		t.Fatal(err)
	}

	_, err = store.Get(ctx, ts.TagKey("a"))

	if err == nil {
		t.Fatal("expected tag key 'a' to be deleted after flush")
	}

	_, err = store.Get(ctx, ts.TagKey("b"))

	if err == nil {
		t.Fatal("expected tag key 'b' to be deleted after flush")
	}
}

func TestTagSetFlushTag(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"a", "b"})
	ctx := context.Background()

	_, _ = ts.TagID(ctx, "a")
	_, _ = ts.TagID(ctx, "b")

	if err := ts.FlushTag(ctx, "a"); err != nil {
		t.Fatal(err)
	}

	_, err := store.Get(ctx, ts.TagKey("a"))

	if err == nil {
		t.Fatal("expected tag 'a' to be deleted")
	}

	// Tag 'b' should still exist.
	_, err = store.Get(ctx, ts.TagKey("b"))

	if err != nil {
		t.Fatal("expected tag 'b' to still exist")
	}
}

func TestTagSetTagKey(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"foo"})

	key := ts.TagKey("foo")

	if key != "tag:foo:key" {
		t.Fatalf("expected 'tag:foo:key', got %q", key)
	}
}
