package queue_test

import (
	"context"
	"sort"
	"testing"

	"github.com/bedrock/packages/queue"
)

type registerableJob struct {
	_    struct{} `queue:"tries=7,timeout=90s,queue=high"`
	Body string
}

func noopHandler() queue.Handler {
	return queue.HandlerFunc(func(ctx context.Context, j queue.Job) error { return nil })
}

func TestHandlerRegistryRegisterAndResolve(t *testing.T) {
	t.Parallel()

	r := queue.NewHandlerRegistry()
	h := noopHandler()

	if err := r.Register("App\\Jobs\\Foo", h, queue.JobOptions{MaxTries: 3}); err != nil {
		t.Fatalf("Register: %v", err)
	}

	entry, ok := r.Resolve("App\\Jobs\\Foo")

	if !ok {
		t.Fatal("Resolve: not found")
	}

	if entry.Name != "App\\Jobs\\Foo" {
		t.Errorf("Name: got %q", entry.Name)
	}

	if entry.Options.MaxTries != 3 {
		t.Errorf("MaxTries: got %d", entry.Options.MaxTries)
	}

	if entry.Handler == nil {
		t.Error("Handler: nil")
	}
}

func TestHandlerRegistryRejectsDuplicateRegister(t *testing.T) {
	t.Parallel()

	r := queue.NewHandlerRegistry()
	h := noopHandler()

	if err := r.Register("dup", h, queue.JobOptions{}); err != nil {
		t.Fatal(err)
	}

	if err := r.Register("dup", h, queue.JobOptions{}); err == nil {
		t.Fatal("expected duplicate register to fail")
	}
}

func TestHandlerRegistryRejectsEmptyName(t *testing.T) {
	t.Parallel()

	r := queue.NewHandlerRegistry()

	if err := r.Register("", noopHandler(), queue.JobOptions{}); err == nil {
		t.Fatal("expected empty-name register to fail")
	}
}

func TestHandlerRegistryRejectsNilHandler(t *testing.T) {
	t.Parallel()

	r := queue.NewHandlerRegistry()

	if err := r.Register("x", nil, queue.JobOptions{}); err == nil {
		t.Fatal("expected nil-handler register to fail")
	}
}

func TestHandlerRegistryReplaceOverrides(t *testing.T) {
	t.Parallel()

	r := queue.NewHandlerRegistry()

	_ = r.Register("x", noopHandler(), queue.JobOptions{MaxTries: 1})
	r.Replace("x", noopHandler(), queue.JobOptions{MaxTries: 9})

	entry, _ := r.Resolve("x")

	if entry.Options.MaxTries != 9 {
		t.Errorf("MaxTries after Replace: got %d, want 9", entry.Options.MaxTries)
	}
}

func TestHandlerRegistryForget(t *testing.T) {
	t.Parallel()

	r := queue.NewHandlerRegistry()

	_ = r.Register("x", noopHandler(), queue.JobOptions{})

	r.Forget("x")

	if _, ok := r.Resolve("x"); ok {
		t.Error("expected x to be forgotten")
	}
}

func TestHandlerRegistryNamesReturnsAll(t *testing.T) {
	t.Parallel()

	r := queue.NewHandlerRegistry()

	_ = r.Register("a", noopHandler(), queue.JobOptions{})
	_ = r.Register("b", noopHandler(), queue.JobOptions{})
	_ = r.Register("c", noopHandler(), queue.JobOptions{})

	names := r.Names()

	sort.Strings(names)

	if len(names) != 3 || names[0] != "a" || names[1] != "b" || names[2] != "c" {
		t.Errorf("Names: got %v", names)
	}
}

func TestHandlerRegistryRegisterForParsesTags(t *testing.T) {
	t.Parallel()

	r := queue.NewHandlerRegistry()

	if err := r.RegisterFor(registerableJob{}, noopHandler()); err != nil {
		t.Fatalf("RegisterFor: %v", err)
	}

	entry, ok := r.Resolve(queue.DisplayName(registerableJob{}))

	if !ok {
		t.Fatal("Resolve: not found")
	}

	if entry.Options.MaxTries != 7 {
		t.Errorf("MaxTries: got %d, want 7", entry.Options.MaxTries)
	}

	if entry.Options.Queue != "high" {
		t.Errorf("Queue: got %q, want high", entry.Options.Queue)
	}
}
