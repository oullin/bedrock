package telescope_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/debugbar"
	"github.com/bedrock/packages/debugbar/watchers"
)

func testContext() context.Context {
	return context.Background()
}

func TestQueryWatcherRecordsBasicQuery(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewQueryWatcher(scope, nil)

	w.Record("SELECT * FROM users WHERE id = ?", []any{1}, 50*time.Millisecond, "mysql", "", 0)

	storeAndAssertCount(t, scope, repo, debugbar.EntryTypeQuery, 1)

	e := repo.Entries()[0]

	if e.Content["connection"] != "mysql" {
		t.Fatalf("expected connection=mysql, got %v", e.Content["connection"])
	}
}

func TestQueryWatcherTagsSlowQueries(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewQueryWatcher(scope, map[string]any{"slow": float64(100)})

	// Fast query — not slow.
	w.Record("SELECT 1", nil, 50*time.Millisecond, "sqlite", "", 0)

	// Slow query.
	w.Record("SELECT * FROM big_table", nil, 200*time.Millisecond, "sqlite", "", 0)

	if err := scope.Store(testContext()); err != nil {
		t.Fatal(err)
	}

	entries := repo.Entries()

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	// The slow one (stored last, retrieved in reverse order from memory repo — but
	// our InMemoryRepository stores in insert order; so index 0 is the fast one).
	fast := entries[0]
	slow := entries[1]

	if fast.Content["slow"] != false {
		t.Fatalf("expected fast query to have slow=false, got %v", fast.Content["slow"])
	}

	if slow.Content["slow"] != true {
		t.Fatalf("expected slow query to have slow=true, got %v", slow.Content["slow"])
	}

	assertHasTag(t, slow.Tags, "slow")
}

func TestQueryWatcherReplacesPositionalBindings(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewQueryWatcher(scope, nil)

	w.Record(
		"SELECT * FROM users WHERE id = ? AND name = ? AND active = ?",
		[]any{1, "Alice", true},
		10*time.Millisecond,
		"mysql",
		"",
		0,
	)

	storeAndAssertCount(t, scope, repo, debugbar.EntryTypeQuery, 1)

	sql, _ := repo.Entries()[0].Content["sql"].(string)

	expected := "SELECT * FROM users WHERE id = 1 AND name = 'Alice' AND active = 1"

	if sql != expected {
		t.Fatalf("binding replacement failed\nwant: %q\ngot:  %q", expected, sql)
	}
}

func TestQueryWatcherHandlesNullBinding(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewQueryWatcher(scope, nil)

	w.Record("INSERT INTO logs (message, user_id) VALUES (?, ?)", []any{"test", nil}, time.Millisecond, "sqlite", "", 0)

	storeAndAssertCount(t, scope, repo, debugbar.EntryTypeQuery, 1)

	sql, _ := repo.Entries()[0].Content["sql"].(string)

	expected := "INSERT INTO logs (message, user_id) VALUES ('test', NULL)"

	if sql != expected {
		t.Fatalf("null binding failed\nwant: %q\ngot:  %q", expected, sql)
	}
}

func TestQueryWatcherReplaceNamedBindings(t *testing.T) {
	t.Parallel()

	// Direct test of the helper function.
	sql := "SELECT * FROM orders WHERE user_id = :user_id AND status = :status"
	bindings := map[string]any{
		"user_id": 42,
		"status":  "active",
	}

	result := watchers.ReplaceNamedBindings(sql, bindings)
	expected := "SELECT * FROM orders WHERE user_id = 42 AND status = 'active'"

	if result != expected {
		t.Fatalf("named binding replacement failed\nwant: %q\ngot:  %q", expected, result)
	}
}

func TestQueryWatcherGeneratesFamilyHash(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewQueryWatcher(scope, nil)

	// Same query, different bindings → same family hash.
	w.Record("SELECT * FROM users WHERE id = ?", []any{1}, time.Millisecond, "mysql", "", 0)
	w.Record("SELECT * FROM users WHERE id = ?", []any{2}, time.Millisecond, "mysql", "", 0)

	// Different query → different family hash.
	w.Record("SELECT * FROM products", nil, time.Millisecond, "mysql", "", 0)

	if err := scope.Store(testContext()); err != nil {
		t.Fatal(err)
	}

	entries := repo.Entries()

	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	hash0 := entries[0].FamilyHash
	hash1 := entries[1].FamilyHash
	hash2 := entries[2].FamilyHash

	if hash0 == "" || hash1 == "" || hash2 == "" {
		t.Fatal("expected non-empty family hash")
	}

	if hash0 != hash1 {
		t.Fatalf("expected identical queries to share family hash (%q != %q)", hash0, hash1)
	}

	if hash0 == hash2 {
		t.Fatal("expected different queries to have different family hashes")
	}
}
