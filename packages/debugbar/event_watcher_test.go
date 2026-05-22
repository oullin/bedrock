package debugbar_test

import (
	"testing"

	"github.com/bedrock/packages/debugbar"
	"github.com/bedrock/packages/debugbar/watchers"
)

type UserRegistered struct {
	UserID int
	Email  string
}

func TestEventWatcherRecordsEvent(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewEventWatcher(scope, nil)

	w.Record("user.registered", UserRegistered{UserID: 1, Email: "alice@example.com"}, nil)

	storeAndAssertCount(t, scope, repo, debugbar.EntryTypeEvent, 1)

	e := repo.Entries()[0]

	if e.Content["name"] != "user.registered" {
		t.Fatalf("expected name=user.registered, got %v", e.Content["name"])
	}
}

func TestEventWatcherCapturesPayload(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewEventWatcher(scope, nil)

	w.Record("order.placed", UserRegistered{UserID: 42, Email: "bob@example.com"}, nil)

	storeAndAssertCount(t, scope, repo, debugbar.EntryTypeEvent, 1)

	payload, _ := repo.Entries()[0].Content["payload"].(map[string]any)

	if payload == nil {
		t.Fatal("expected non-nil payload")
	}
}

func TestEventWatcherCapturesListeners(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewEventWatcher(scope, nil)

	listeners := []string{"SendWelcomeEmail", "CreateUserProfile"}
	w.Record("user.created", nil, listeners)

	storeAndAssertCount(t, scope, repo, debugbar.EntryTypeEvent, 1)

	stored, _ := repo.Entries()[0].Content["listeners"].([]string)

	if len(stored) != 2 {
		t.Fatalf("expected 2 listeners, got %v", stored)
	}
}

func TestEventWatcherIgnoresFrameworkEvents(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewEventWatcher(scope, nil)

	// These should all be filtered.
	w.Record("@bedrock\\Auth\\Events\\Login", nil, nil)
	w.Record("Octane\\Events\\RequestReceived", nil, nil)
	w.Record("github.com/bedrock/packages/debugbar.InternalEvent", nil, nil)

	storeAndAssertCount(t, scope, repo, debugbar.EntryTypeEvent, 0)
}

func TestEventWatcherIgnoresConfiguredEvents(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewEventWatcher(scope, map[string]any{
		"ignore": []string{"payment.processed"},
	})

	w.Record("payment.processed", nil, nil)
	w.Record("user.registered", nil, nil)

	storeAndAssertCount(t, scope, repo, debugbar.EntryTypeEvent, 1)
}
