package session_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/session"
	"github.com/bedrock/packages/session/handlers"
)

func newStore() *session.Store {
	s := session.New("test", handlers.NewArrayHandler())
	_ = s.Start(context.Background())

	return s
}

func TestStorePutGet(t *testing.T) {
	s := newStore()

	s.Put("key", "value")

	if got := s.Get("key", nil); got != "value" {
		t.Errorf("got %v, want value", got)
	}
}

func TestStoreHasExists(t *testing.T) {
	s := newStore()

	s.Put("k", "v")

	if !s.Has("k") {
		t.Error("Has should return true for existing key")
	}

	s.Put("nil_key", nil)

	if s.Has("nil_key") {
		t.Error("Has should return false for nil value")
	}

	if !s.Exists("nil_key") {
		t.Error("Exists should return true for nil value")
	}
}

func TestStorePull(t *testing.T) {
	s := newStore()

	s.Put("k", "v")

	got := s.Pull("k", nil)

	if got != "v" {
		t.Errorf("got %v, want v", got)
	}

	if s.Exists("k") {
		t.Error("key should be removed after Pull")
	}
}

func TestStorePush(t *testing.T) {
	s := newStore()

	s.Push("list", "a")
	s.Push("list", "b")

	v := s.Get("list", nil)
	sl, ok := v.([]any)

	if !ok || len(sl) != 2 {
		t.Errorf("expected slice of length 2, got %v", v)
	}
}

func TestStoreFlash(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	s1 := session.New("test", h)
	_ = s1.Start(ctx)
	s1.Flash("msg", "hello")
	_ = s1.Save(ctx)

	// Second request: flash data should be visible.
	s2 := session.NewWithID("test", h, s1.GetID())
	_ = s2.Start(ctx)

	if got := s2.Get("msg", nil); got != "hello" {
		t.Errorf("flash not visible on next request: got %v", got)
	}

	_ = s2.Save(ctx)

	// Third request: flash data should be gone.
	s3 := session.NewWithID("test", h, s1.GetID())
	_ = s3.Start(ctx)

	if s3.Has("msg") {
		t.Error("flash data should be removed after second request")
	}
}

func TestStoreToken(t *testing.T) {
	s := newStore()

	tok := s.Token()

	if len(tok) == 0 {
		t.Error("expected non-empty CSRF token")
	}

	if s.Token() != tok {
		t.Error("Token should return the same value on repeated calls")
	}

	s.RegenerateToken()

	if s.Token() == tok {
		t.Error("RegenerateToken should produce a different token")
	}
}

func TestStoreRegenerate(t *testing.T) {
	s := newStore()

	oldID := s.GetID()
	_ = s.Regenerate(context.Background(), false)

	if s.GetID() == oldID {
		t.Error("Regenerate should produce a new ID")
	}
}

func TestStoreInvalidate(t *testing.T) {
	s := newStore()

	s.Put("k", "v")
	_ = s.Invalidate(context.Background())

	if s.Has("k") {
		t.Error("Invalidate should flush all data")
	}

	if s.IsStarted() {
		t.Error("Invalidate should mark session as not started")
	}
}

func TestStoreIncrement(t *testing.T) {
	s := newStore()

	s.Put("n", int64(10))

	got := s.Increment("n", 5)

	if got != 15 {
		t.Errorf("got %d, want 15", got)
	}

	got = s.Decrement("n", 3)

	if got != 12 {
		t.Errorf("got %d, want 12", got)
	}
}

func TestStoreReflash(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	s1 := session.New("test", h)
	_ = s1.Start(ctx)
	s1.Flash("x", 1)
	_ = s1.Save(ctx)

	s2 := session.NewWithID("test", h, s1.GetID())
	_ = s2.Start(ctx)
	s2.Reflash()
	_ = s2.Save(ctx)

	// Data should still exist on the third request due to Reflash.
	s3 := session.NewWithID("test", h, s1.GetID())
	_ = s3.Start(ctx)

	if !s3.Has("x") {
		t.Error("Reflash should keep flash data for one more request")
	}
}

func TestStoreAlreadyStarted(t *testing.T) {
	s := newStore()

	if err := s.Start(context.Background()); err != session.ErrAlreadyStarted {
		t.Errorf("expected ErrAlreadyStarted, got %v", err)
	}
}

func TestStorePasswordConfirmed(t *testing.T) {
	s := newStore()

	before := s.PasswordConfirmedAt()

	if before != 0 {
		t.Error("expected 0 before confirmation")
	}

	s.PasswordConfirmed()

	after := s.PasswordConfirmedAt()

	if after == 0 {
		t.Error("expected non-zero timestamp after confirmation")
	}
}
