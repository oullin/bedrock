package telescope_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/telescope"
	"github.com/bedrock/packages/telescope/watchers"
)

func TestExceptionWatcherRecordsException(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewExceptionWatcher(scope, nil)

	err := errors.New("something went wrong")
	w.Record(err, 0)

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeException, 1)

	e := repo.Entries()[0]

	if e.Content["message"] != "something went wrong" {
		t.Fatalf("expected message=%q, got %v", "something went wrong", e.Content["message"])
	}

	class, _ := e.Content["class"].(string)

	if class == "" {
		t.Fatal("expected non-empty class in exception entry")
	}

	if e.Content["file"] == "" {
		t.Fatal("expected non-empty file in exception entry")
	}
}

func TestExceptionWatcherCapturesStackTrace(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewExceptionWatcher(scope, nil)

	w.Record(errors.New("traced"), 0)

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeException, 1)

	trace, _ := repo.Entries()[0].Content["trace"].([]map[string]any)

	if len(trace) == 0 {
		t.Fatal("expected non-empty stack trace")
	}
}

func TestExceptionWatcherIgnoresConfiguredExceptionTypes(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewExceptionWatcher(scope, map[string]any{
		"ignore": []string{"*errors.errorString"},
	})

	w.Record(errors.New("ignored error"), 0)

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeException, 0)
}

func TestExceptionWatcherTagsEntryWithExceptionType(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewExceptionWatcher(scope, nil)

	w.Record(errors.New("tagged error"), 0)

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeException, 1)

	entry := repo.Entries()[0]
	class, _ := entry.Content["class"].(string)

	assertHasTag(t, entry.Tags, class)
}

func TestExceptionWatcherRecordRawStoresMetadata(t *testing.T) {
	t.Parallel()

	scope, repo := newTestScope(t)
	w := watchers.NewExceptionWatcher(scope, nil)

	w.RecordRaw("MyApp\\CustomException", "/app/handler.go", 42, "custom error", nil)

	storeAndAssertCount(t, scope, repo, telescope.EntryTypeException, 1)

	e := repo.Entries()[0]

	if e.Content["class"] != "MyApp\\CustomException" {
		t.Fatalf("unexpected class: %v", e.Content["class"])
	}

	if e.Content["file"] != "/app/handler.go" {
		t.Fatalf("unexpected file: %v", e.Content["file"])
	}

	if e.Content["line"] != 42 {
		t.Fatalf("unexpected line: %v", e.Content["line"])
	}

	if e.Content["message"] != "custom error" {
		t.Fatalf("unexpected message: %v", e.Content["message"])
	}
}
