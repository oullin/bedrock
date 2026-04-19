package telescope_test

import (
	"testing"

	"github.com/bedrock/packages/telescope"
	"github.com/bedrock/packages/telescope/watchers"
)

// TestLogWatcherRecordsDebugLevel mirrors RequestWatchersTest for log entries.
func TestLogWatcherRecordsDebugLevel(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewLogWatcher(scope, nil)

	w.Record("debug", "debug message", nil)

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeLog, 1)

	e := repo.Entries()[0]

	if e.Content["level"] != "debug" {
		t.Fatalf("expected level=debug, got %v", e.Content["level"])
	}

	if e.Content["message"] != "debug message" {
		t.Fatalf("expected message=%q, got %v", "debug message", e.Content["message"])
	}
}

func TestLogWatcherRecordsAllPSR3Levels(t *testing.T) {
	t.Parallel()

	levels := []string{"debug", "info", "notice", "warning", "error", "critical", "alert", "emergency"}

	for _, level := range levels {
		level := level

		t.Run(level, func(t *testing.T) {
			t.Parallel()

			scope, repo := newTestScope(t)
			w := watchers.NewLogWatcher(scope, nil)

			w.Record(level, level+" message", nil)

			storeAndAssertCount(t, scope, repo, telescope.EntryTypeLog, 1)

			e := repo.Entries()[0]

			if e.Content["level"] != level {
				t.Fatalf("expected level=%q, got %v", level, e.Content["level"])
			}
		})
	}
}

func TestLogWatcherFiltersLevelsBelowMinimum(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewLogWatcher(scope, map[string]any{"level": "error"})

	w.Record("debug", "too low", nil)
	w.Record("info", "also too low", nil)
	w.Record("warning", "still too low", nil)
	w.Record("error", "meets threshold", nil)
	w.Record("critical", "above threshold", nil)

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeLog, 2)
}

func TestLogWatcherInterpolatesPlaceholders(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewLogWatcher(scope, nil)

	w.Record("info", "Hello {name}, you are {age} years old", map[string]any{
		"name": "Alice",
		"age":  "30",
	})

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeLog, 1)

	msg := repo.Entries()[0].Content["message"]

	if msg != "Hello Alice, you are 30 years old" {
		t.Fatalf("unexpected interpolated message: %v", msg)
	}
}

func TestLogWatcherSkipsExceptionLogs(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewLogWatcher(scope, nil)

	w.Record("error", "an error occurred", map[string]any{
		"exception": "some exception",
	})

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeLog, 0)
}

func TestLogWatcherExtractsTelescopeTagsFromContext(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewLogWatcher(scope, nil)

	w.Record("info", "tagged log", map[string]any{
		"telescope": []string{"billing", "user"},
	})

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeLog, 1)

	entry := repo.Entries()[0]
	assertHasTag(t, entry.Tags, "billing")
	assertHasTag(t, entry.Tags, "user")
}

func TestLogWatcherStripsTelescopeKeyFromStoredContext(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewLogWatcher(scope, nil)

	w.Record("info", "tagged log", map[string]any{
		"telescope": []string{"tag1"},
		"user_id":   42,
	})

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeLog, 1)

	ctx, _ := repo.Entries()[0].Content["context"].(map[string]any)

	if _, hasTelescope := ctx["telescope"]; hasTelescope {
		t.Fatal("expected telescope key to be stripped from context")
	}

	if _, hasUserID := ctx["user_id"]; !hasUserID {
		t.Fatal("expected user_id to be preserved in context")
	}
}

// ─── Shared test helpers ──────────────────────────────────────────────────────

func storeAndAssertCount(t *testing.T, scope *telescope.Telescope, repo interface {
	Count() int
	Entries() []*telescope.IncomingEntry
}, entryType string, want int) {
	t.Helper()

	if err := scope.Store(testContext()); err != nil {
		t.Fatalf("store: %v", err)
	}

	got := 0

	for _, e := range repo.Entries() {
		if e.Type == entryType {
			got++
		}
	}

	if got != want {
		t.Fatalf("expected %d %q entries, got %d", want, entryType, got)
	}
}

func assertHasTag(t *testing.T, tags []string, want string) {
	t.Helper()

	for _, tag := range tags {
		if tag == want {
			return
		}
	}

	t.Fatalf("expected tag %q in %v", want, tags)
}
