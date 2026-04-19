package telescope_test

import (
	"testing"

	"github.com/bedrock/packages/telescope"
	"github.com/bedrock/packages/telescope/watchers"
)

func TestCacheWatcherRecordsHit(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewCacheWatcher(scope, nil)

	w.Hit("users:1", map[string]any{"id": 1})

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeCache, 1)

	e := repo.Entries()[0]

	if e.Content["type"] != "hit" {
		t.Fatalf("expected type=hit, got %v", e.Content["type"])
	}

	if e.Content["key"] != "users:1" {
		t.Fatalf("expected key=users:1, got %v", e.Content["key"])
	}
}

func TestCacheWatcherRecordsMiss(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewCacheWatcher(scope, nil)

	w.Missed("users:999")

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeCache, 1)

	e := repo.Entries()[0]

	if e.Content["type"] != "missed" {
		t.Fatalf("expected type=missed, got %v", e.Content["type"])
	}
}

func TestCacheWatcherRecordsWrite(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewCacheWatcher(scope, nil)

	w.Written("session:abc", "data", 3600)

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeCache, 1)

	e := repo.Entries()[0]

	if e.Content["type"] != "set" {
		t.Fatalf("expected type=set, got %v", e.Content["type"])
	}

	if e.Content["key"] != "session:abc" {
		t.Fatalf("expected key=session:abc, got %v", e.Content["key"])
	}

	if e.Content["expires"] != int64(3600) {
		t.Fatalf("expected expires=3600, got %v", e.Content["expires"])
	}
}

func TestCacheWatcherRecordsForget(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewCacheWatcher(scope, nil)

	w.Forgotten("users:1")

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeCache, 1)

	if repo.Entries()[0].Content["type"] != "forget" {
		t.Fatalf("expected type=forget, got %v", repo.Entries()[0].Content["type"])
	}
}

func TestCacheWatcherMasksHiddenKeys(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewCacheWatcher(scope, map[string]any{
		"hidden": []string{"secret_token"},
	})

	w.Hit("secret_token", "super-secret-value")

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeCache, 1)

	e := repo.Entries()[0]

	if e.Content["value"] != "********" {
		t.Fatalf("expected masked value, got %v", e.Content["value"])
	}
}

func TestCacheWatcherIgnoresFrameworkKeys(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewCacheWatcher(scope, nil)

	w.Hit("illuminate:queue:restart", "value")
	w.Hit("telescope:recording", "value")
	w.Written("framework/schedule/lock", "val", 60)

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeCache, 0)
}
