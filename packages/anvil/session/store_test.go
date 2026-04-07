package session

import (
	"context"
	"sync"
	"testing"
)

func TestNewStore(t *testing.T) {
	t.Parallel()

	s := New("app_session", NewArrayHandler())

	if s.GetName() != "app_session" {
		t.Fatalf("want name %q, got %q", "app_session", s.GetName())
	}

	id := s.GetID()
	if !isValidID(id) {
		t.Fatalf("generated ID is not valid: %q", id)
	}
}

func TestStartLoadsFromHandler(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	ctx := context.Background()

	_ = h.Write(ctx, "a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0", `{"user_id":"42"}`)

	s := NewWithID("sess", h, "a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0")

	if err := s.Start(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := s.Get("user_id", nil)
	if got != "42" {
		t.Fatalf("want %q, got %v", "42", got)
	}
}

func TestStartEmptySession(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	if err := s.Start(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !s.IsStarted() {
		t.Fatal("session should be started")
	}
}

func TestStartAlreadyStarted(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())
	ctx := context.Background()

	_ = s.Start(ctx)

	err := s.Start(ctx)
	if err != ErrAlreadyStarted {
		t.Fatalf("want ErrAlreadyStarted, got %v", err)
	}
}

func TestSave(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	s := New("sess", h)
	ctx := context.Background()

	_ = s.Start(ctx)
	s.Put("key", "value")

	if err := s.Save(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := h.Read(ctx, s.GetID())
	if data == "" {
		t.Fatal("handler should have data after Save")
	}
}

func TestGetAndPut(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("key", "value")

	got := s.Get("key", nil)
	if got != "value" {
		t.Fatalf("want %q, got %v", "value", got)
	}
}

func TestGetWithFallback(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	got := s.Get("missing", "default")
	if got != "default" {
		t.Fatalf("want %q, got %v", "default", got)
	}
}

func TestHasExistsMissing(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("key", "value")
	s.Put("nil_key", nil)

	if !s.Has("key") {
		t.Fatal("Has should return true for non-nil value")
	}

	if s.Has("nil_key") {
		t.Fatal("Has should return false for nil value")
	}

	if !s.Exists("nil_key") {
		t.Fatal("Exists should return true for nil value")
	}

	if !s.Missing("unknown") {
		t.Fatal("Missing should return true for non-existent key")
	}
}

func TestPull(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("key", "value")

	got := s.Pull("key", nil)
	if got != "value" {
		t.Fatalf("want %q, got %v", "value", got)
	}

	if s.Exists("key") {
		t.Fatal("key should be removed after Pull")
	}
}

func TestPullFallback(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	got := s.Pull("missing", "default")
	if got != "default" {
		t.Fatalf("want %q, got %v", "default", got)
	}
}

func TestPush(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Push("items", "a")
	s.Push("items", "b")

	got := s.Get("items", nil)
	slice, ok := got.([]any)
	if !ok {
		t.Fatalf("want []any, got %T", got)
	}

	if len(slice) != 2 || slice[0] != "a" || slice[1] != "b" {
		t.Fatalf("want [a, b], got %v", slice)
	}
}

func TestPushCreatesSlice(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Push("items", "first")

	got := s.Get("items", nil)
	slice, ok := got.([]any)
	if !ok || len(slice) != 1 {
		t.Fatalf("want slice of length 1, got %v", got)
	}
}

func TestAll(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("a", 1)
	s.Put("b", 2)

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("want 2 attributes, got %d", len(all))
	}

	// Verify it's a copy.
	all["c"] = 3
	if s.Exists("c") {
		t.Fatal("All should return a copy, not a reference")
	}
}

func TestForgetKeys(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	s.Forget("a", "c")

	if s.Exists("a") || s.Exists("c") {
		t.Fatal("forgotten keys should not exist")
	}

	if !s.Exists("b") {
		t.Fatal("non-forgotten key should still exist")
	}
}

func TestFlush(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("a", 1)
	s.Put("b", 2)
	s.Flush()

	all := s.All()
	if len(all) != 0 {
		t.Fatalf("want empty map after Flush, got %d items", len(all))
	}
}

func TestFlashAndAging(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	ctx := context.Background()

	// Request 1: flash a value.
	s1 := New("sess", h)
	_ = s1.Start(ctx)
	s1.Flash("message", "hello")
	_ = s1.Save(ctx)

	// Request 2: flash data is available, then aged.
	s2 := NewWithID("sess", h, s1.GetID())
	_ = s2.Start(ctx)

	got := s2.Get("message", nil)
	if got != "hello" {
		t.Fatalf("want %q, got %v", "hello", got)
	}

	_ = s2.Save(ctx)

	// Request 3: flash data should be gone.
	s3 := NewWithID("sess", h, s2.GetID())
	_ = s3.Start(ctx)

	got = s3.Get("message", nil)
	if got != nil {
		t.Fatalf("flash data should be gone after aging, got %v", got)
	}
}

func TestReflash(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	ctx := context.Background()

	s1 := New("sess", h)
	_ = s1.Start(ctx)
	s1.Flash("msg", "kept")
	_ = s1.Save(ctx)

	// Request 2: reflash to keep it.
	s2 := NewWithID("sess", h, s1.GetID())
	_ = s2.Start(ctx)
	s2.Reflash()
	_ = s2.Save(ctx)

	// Request 3: should still be available.
	s3 := NewWithID("sess", h, s2.GetID())
	_ = s3.Start(ctx)

	got := s3.Get("msg", nil)
	if got != "kept" {
		t.Fatalf("reflashed data should still be available, got %v", got)
	}
}

func TestKeep(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	ctx := context.Background()

	s1 := New("sess", h)
	_ = s1.Start(ctx)
	s1.Flash("a", "1")
	s1.Flash("b", "2")
	_ = s1.Save(ctx)

	// Request 2: keep only "a".
	s2 := NewWithID("sess", h, s1.GetID())
	_ = s2.Start(ctx)
	s2.Keep("a")
	_ = s2.Save(ctx)

	// Request 3: "a" should survive, "b" should be gone.
	s3 := NewWithID("sess", h, s2.GetID())
	_ = s3.Start(ctx)

	if s3.Get("a", nil) != "1" {
		t.Fatal("kept flash key 'a' should still exist")
	}

	if s3.Get("b", nil) != nil {
		t.Fatal("non-kept flash key 'b' should be aged out")
	}
}

func TestToken(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	token := s.Token()
	if token == "" {
		t.Fatal("Token should generate a non-empty token")
	}

	if s.Token() != token {
		t.Fatal("Token should return the same token on repeated calls")
	}
}

func TestRegenerateToken(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	old := s.Token()
	s.RegenerateToken()

	if s.Token() == old {
		t.Fatal("RegenerateToken should produce a different token")
	}
}

func TestGetID(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	id := s.GetID()
	if len(id) != 40 {
		t.Fatalf("want 40-char ID, got %d chars", len(id))
	}
}

func TestSetID(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	validID := "abcdef0123456789abcdef0123456789abcdef01"

	if err := s.SetID(validID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.GetID() != validID {
		t.Fatalf("want %q, got %q", validID, s.GetID())
	}
}

func TestSetIDInvalid(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	cases := []string{
		"too-short",
		"xyz_not_hex_at_all_needs_forty_characters",
		"",
	}

	for _, id := range cases {
		if err := s.SetID(id); err == nil {
			t.Fatalf("expected error for invalid ID %q", id)
		}
	}
}

func TestGetSetName(t *testing.T) {
	t.Parallel()

	s := New("original", NewArrayHandler())

	s.SetName("renamed")

	if s.GetName() != "renamed" {
		t.Fatalf("want %q, got %q", "renamed", s.GetName())
	}
}

func TestRegenerate(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())
	_ = s.Start(context.Background())
	s.Put("key", "value")

	oldID := s.GetID()

	if err := s.Regenerate(context.Background(), false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.GetID() == oldID {
		t.Fatal("Regenerate should produce a new ID")
	}

	if s.Get("key", nil) != "value" {
		t.Fatal("data should be preserved after Regenerate without destroy")
	}
}

func TestRegenerateWithDestroy(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	s := New("sess", h)
	ctx := context.Background()

	_ = s.Start(ctx)
	s.Put("key", "value")
	_ = s.Save(ctx)

	oldID := s.GetID()

	if err := s.Regenerate(ctx, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := h.Read(ctx, oldID)
	if data != "" {
		t.Fatal("old session data should be destroyed")
	}
}

func TestInvalidate(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())
	ctx := context.Background()

	_ = s.Start(ctx)
	s.Put("key", "value")

	oldID := s.GetID()

	if err := s.Invalidate(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.GetID() == oldID {
		t.Fatal("Invalidate should produce a new ID")
	}

	if s.Exists("key") {
		t.Fatal("Invalidate should flush all data")
	}
}

func TestIsStarted(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	if s.IsStarted() {
		t.Fatal("should not be started initially")
	}

	_ = s.Start(context.Background())

	if !s.IsStarted() {
		t.Fatal("should be started after Start")
	}
}

func TestPreviousURL(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	if s.PreviousURL() != "" {
		t.Fatal("PreviousURL should be empty initially")
	}

	s.SetPreviousURL("/dashboard")

	if s.PreviousURL() != "/dashboard" {
		t.Fatalf("want %q, got %q", "/dashboard", s.PreviousURL())
	}
}

func TestConcurrentGetPut(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			s.Put("key", n)
		}(i)
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = s.Get("key", nil)
		}()
	}

	wg.Wait()
}
