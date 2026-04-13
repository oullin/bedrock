package session_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/session"
	"github.com/bedrock/packages/session/handlers"
)

// --- mock types ---

type mockRequestAwareHandler struct {
	*handlers.ArrayHandler
	request *http.Request
}

type mockExistenceAwareHandler struct {
	*handlers.ArrayHandler
	exists bool
}

func (h *mockRequestAwareHandler) SetRequest(r *http.Request) { h.request = r }

func (h *mockExistenceAwareHandler) SetExists(v bool) { h.exists = v }

// --- helpers ---

func newStore() *session.Store {
	s := session.New("test", handlers.NewArrayHandler())
	_ = s.Start(context.Background())

	return s
}

func newStoreWith(h session.Handler) *session.Store {
	s := session.New("test", h)
	_ = s.Start(context.Background())

	return s
}

// --- lifecycle tests ---

func TestStoreIsLoadedFromHandler(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	// Pre-populate via a first session.
	s1 := session.New("test", h)
	_ = s1.Start(ctx)
	s1.Put("foo", "bar")
	_ = s1.Save(ctx)

	// Load into a new store with the same ID.
	s2 := session.NewWithID("test", h, s1.GetID())
	_ = s2.Start(ctx)

	if got := s2.Get("foo", nil); got != "bar" {
		t.Errorf("expected bar, got %v", got)
	}
}

func TestStoreSessionMigration(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	s := session.New("test", h)
	_ = s.Start(ctx)
	s.Put("k", "v")
	_ = s.Save(ctx)

	oldID := s.GetID()

	// Migrate without destroy.
	_ = s.Migrate(ctx, false)

	if s.GetID() == oldID {
		t.Error("Migrate should produce a new ID")
	}

	if got := s.Get("k", nil); got != "v" {
		t.Error("attributes should be preserved after Migrate")
	}

	// Migrate with destroy.
	_ = s.Migrate(ctx, true)

	if s.Get("k", nil) != "v" {
		t.Error("attributes should be preserved; only the old handler data is destroyed")
	}
}

func TestStoreCantSetInvalidID(t *testing.T) {
	s := newStore()

	if err := s.SetID("invalid"); err == nil {
		t.Error("expected error for invalid ID")
	}

	if err := s.SetID(""); err == nil {
		t.Error("expected error for empty ID")
	}

	valid := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	if err := s.SetID(valid); err != nil {
		t.Errorf("unexpected error for valid ID: %v", err)
	}

	if s.GetID() != valid {
		t.Errorf("expected ID %q, got %q", valid, s.GetID())
	}
}

func TestStoreBrandNewSessionIsProperlySaved(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	s := session.New("test", h)
	_ = s.Start(ctx)
	s.Put("foo", "bar")
	s.Flash("baz", "boom")
	s.Now("qux", "norf")
	_ = s.Save(ctx)

	// Verify handler has the data by loading into a new store.
	s2 := session.NewWithID("test", h, s.GetID())
	_ = s2.Start(ctx)

	if got := s2.Get("foo", nil); got != "bar" {
		t.Errorf("expected bar, got %v", got)
	}

	// Flash data "baz" should be available on next request.
	if got := s2.Get("baz", nil); got != "boom" {
		t.Errorf("expected boom, got %v", got)
	}

	// "qux" was Now(), so it was in old flash on save, aged out on s2.Start.
	if s2.Has("qux") {
		t.Error("Now() data should not survive to next request")
	}
}

func TestStoreSessionIsProperlyUpdated(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	s1 := session.New("test", h)
	_ = s1.Start(ctx)
	s1.Put("foo", "bar")
	_ = s1.Save(ctx)

	s2 := session.NewWithID("test", h, s1.GetID())
	_ = s2.Start(ctx)
	s2.Put("foo", "baz")
	_ = s2.Save(ctx)

	s3 := session.NewWithID("test", h, s1.GetID())
	_ = s3.Start(ctx)

	if got := s3.Get("foo", nil); got != "baz" {
		t.Errorf("expected baz, got %v", got)
	}
}

func TestStoreSessionIsReSavedWhenNothingChanged(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	s := session.New("test", h)
	_ = s.Start(ctx)
	_ = s.Save(ctx)

	// The handler should have data (at least flash metadata).
	s2 := session.NewWithID("test", h, s.GetID())
	_ = s2.Start(ctx)

	// No panic, no error — session saved even with no explicit puts.
	if !s2.IsStarted() {
		t.Error("session should be started")
	}
}

func TestStoreSessionIsReSavedAfterMigration(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	s := session.New("test", h)
	_ = s.Start(ctx)
	s.Put("foo", "bar")
	_ = s.Regenerate(ctx, false)
	_ = s.Save(ctx)

	s2 := session.NewWithID("test", h, s.GetID())
	_ = s2.Start(ctx)

	if got := s2.Get("foo", nil); got != "bar" {
		t.Errorf("expected bar, got %v", got)
	}
}

// --- flash data tests ---

func TestStoreOldInputFlashing(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	s1 := session.New("test", h)
	_ = s1.Start(ctx)
	s1.FlashInput(map[string]any{"name": "Taylor", "email": "taylor@example.com"})
	_ = s1.Save(ctx)

	s2 := session.NewWithID("test", h, s1.GetID())
	_ = s2.Start(ctx)

	if !s2.HasOldInput("name") {
		t.Error("expected old input for name")
	}

	if got := s2.GetOldInput("name", nil); got != "Taylor" {
		t.Errorf("expected Taylor, got %v", got)
	}

	if got := s2.GetOldInput("email", nil); got != "taylor@example.com" {
		t.Errorf("expected taylor@example.com, got %v", got)
	}

	if got := s2.GetOldInput("missing", "default"); got != "default" {
		t.Errorf("expected default, got %v", got)
	}
}

func TestStoreDataFlashingNow(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	s1 := session.New("test", h)
	_ = s1.Start(ctx)
	s1.Now("foo", "bar")

	// Available in current request.
	if got := s1.Get("foo", nil); got != "bar" {
		t.Errorf("expected bar, got %v", got)
	}

	_ = s1.Save(ctx)

	// Not available on next request.
	s2 := session.NewWithID("test", h, s1.GetID())
	_ = s2.Start(ctx)

	if s2.Has("foo") {
		t.Error("Now() data should not survive to next request")
	}
}

func TestStoreDataMergeNewFlashes(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	s1 := session.New("test", h)
	_ = s1.Start(ctx)
	s1.Flash("foo", "bar")
	_ = s1.Save(ctx)

	s2 := session.NewWithID("test", h, s1.GetID())
	_ = s2.Start(ctx)

	// foo is now in "old". Keep it for another request.
	s2.Keep("foo")
	_ = s2.Save(ctx)

	s3 := session.NewWithID("test", h, s1.GetID())
	_ = s3.Start(ctx)

	if got := s3.Get("foo", nil); got != "bar" {
		t.Errorf("Keep should preserve flash data, got %v", got)
	}
}

func TestStoreReflashWithNow(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	s1 := session.New("test", h)
	_ = s1.Start(ctx)
	s1.Now("foo", "bar")
	s1.Reflash()
	_ = s1.Save(ctx)

	// Reflash moves old→new, so "foo" should survive to next request.
	s2 := session.NewWithID("test", h, s1.GetID())
	_ = s2.Start(ctx)

	if got := s2.Get("foo", nil); got != "bar" {
		t.Errorf("Reflash should keep Now() data, got %v", got)
	}
}

func TestStoreHasOldInputWithoutKey(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	s1 := session.New("test", h)
	_ = s1.Start(ctx)
	s1.FlashInput(map[string]any{"name": "Taylor"})
	_ = s1.Save(ctx)

	s2 := session.NewWithID("test", h, s1.GetID())
	_ = s2.Start(ctx)

	// Empty key checks if any old input exists.
	if !s2.HasOldInput("") {
		t.Error("expected HasOldInput('') to return true when old input exists")
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

// --- attribute operation tests ---

func TestStoreOnly(t *testing.T) {
	s := newStore()

	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	result := s.Only("a", "c")

	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}

	if result["a"] != 1 || result["c"] != 3 {
		t.Errorf("unexpected values: %v", result)
	}

	if _, ok := result["b"]; ok {
		t.Error("b should not be in Only result")
	}
}

func TestStoreExcept(t *testing.T) {
	s := newStore()

	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	result := s.Except("b")

	if _, ok := result["b"]; ok {
		t.Error("b should not be in Except result")
	}

	if result["a"] != 1 || result["c"] != 3 {
		t.Errorf("unexpected values: %v", result)
	}
}

func TestStoreReplace(t *testing.T) {
	s := newStore()

	s.Put("a", 1)
	s.Put("b", 2)

	s.Replace(map[string]any{"b": 99, "c": 3})

	if got := s.Get("a", nil); got != 1 {
		t.Errorf("expected 1, got %v", got)
	}

	if got := s.Get("b", nil); got != 99 {
		t.Errorf("expected 99, got %v", got)
	}

	if got := s.Get("c", nil); got != 3 {
		t.Errorf("expected 3, got %v", got)
	}
}

func TestStoreRemove(t *testing.T) {
	s := newStore()

	s.Put("k", "v")

	got := s.Remove("k")

	if got != "v" {
		t.Errorf("expected v, got %v", got)
	}

	if s.Exists("k") {
		t.Error("key should be removed after Remove")
	}
}

func TestStoreClear(t *testing.T) {
	s := newStore()

	s.Put("a", 1)
	s.Put("b", 2)
	s.Flush()

	if len(s.All()) != 0 {
		t.Error("Flush should clear all attributes")
	}
}

func TestStoreForgetMultiple(t *testing.T) {
	s := newStore()

	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	s.Forget("a", "c")

	if s.Exists("a") {
		t.Error("a should be forgotten")
	}

	if !s.Exists("b") {
		t.Error("b should remain")
	}

	if s.Exists("c") {
		t.Error("c should be forgotten")
	}
}

func TestStoreKeyHas(t *testing.T) {
	s := newStore()

	s.Put("a", "val")
	s.Put("b", nil)

	if !s.Has("a") {
		t.Error("Has should return true for non-nil value")
	}

	if s.Has("b") {
		t.Error("Has should return false for nil value")
	}

	if s.Has("c") {
		t.Error("Has should return false for missing key")
	}
}

func TestStoreKeyHasAny(t *testing.T) {
	s := newStore()

	s.Put("a", "val")

	if !s.HasAny("a", "b", "c") {
		t.Error("HasAny should return true when at least one key exists")
	}

	if s.HasAny("b", "c") {
		t.Error("HasAny should return false when no keys exist")
	}
}

func TestStoreKeyExists(t *testing.T) {
	s := newStore()

	s.Put("a", "val")
	s.Put("b", nil)

	if !s.Exists("a") {
		t.Error("Exists should return true for non-nil value")
	}

	if !s.Exists("b") {
		t.Error("Exists should return true for nil value")
	}

	if s.Exists("c") {
		t.Error("Exists should return false for missing key")
	}
}

func TestStoreKeyMissing(t *testing.T) {
	s := newStore()

	s.Put("a", "val")

	if s.Missing("a") {
		t.Error("Missing should return false for existing key")
	}

	if !s.Missing("b") {
		t.Error("Missing should return true for non-existent key")
	}
}

func TestStoreRememberNew(t *testing.T) {
	s := newStore()

	called := false

	got := s.Remember("key", func() any {
		called = true

		return "computed"
	})

	if !called {
		t.Error("callback should be called for new key")
	}

	if got != "computed" {
		t.Errorf("expected computed, got %v", got)
	}

	if s.Get("key", nil) != "computed" {
		t.Error("Remember should store the result")
	}
}

func TestStoreRememberExisting(t *testing.T) {
	s := newStore()

	s.Put("key", "existing")

	called := false

	got := s.Remember("key", func() any {
		called = true

		return "new"
	})

	if called {
		t.Error("callback should not be called when key exists")
	}

	if got != "existing" {
		t.Errorf("expected existing, got %v", got)
	}
}

// --- existing tests preserved ---

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

func TestStoreIncrementDefault(t *testing.T) {
	s := newStore()

	got := s.Increment("counter", 1)

	if got != 1 {
		t.Errorf("expected 1 for new key, got %d", got)
	}
}

func TestStoreDecrementDefault(t *testing.T) {
	s := newStore()

	got := s.Decrement("counter", 1)

	if got != -1 {
		t.Errorf("expected -1 for new key, got %d", got)
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

// --- new method tests ---

func TestStoreGetSetHandler(t *testing.T) {
	h1 := handlers.NewArrayHandler()
	s := newStoreWith(h1)

	if s.GetHandler() != h1 {
		t.Error("GetHandler should return the original handler")
	}

	h2 := handlers.NewArrayHandler()
	old := s.SetHandler(h2)

	if old != h1 {
		t.Error("SetHandler should return the old handler")
	}

	if s.GetHandler() != h2 {
		t.Error("GetHandler should return the new handler")
	}
}

func TestStoreHandlerNeedsRequest(t *testing.T) {
	// ArrayHandler does not implement RequestAware.
	s := newStore()

	if s.HandlerNeedsRequest() {
		t.Error("ArrayHandler should not need request")
	}

	// mockRequestAwareHandler does implement RequestAware.
	mh := &mockRequestAwareHandler{ArrayHandler: handlers.NewArrayHandler()}
	s2 := newStoreWith(mh)

	if !s2.HandlerNeedsRequest() {
		t.Error("mockRequestAwareHandler should need request")
	}
}

func TestStoreSetRequestOnHandler(t *testing.T) {
	mh := &mockRequestAwareHandler{ArrayHandler: handlers.NewArrayHandler()}
	s := newStoreWith(mh)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	s.SetRequestOnHandler(req)

	if mh.request != req {
		t.Error("SetRequestOnHandler should forward request to handler")
	}
}

func TestStoreSetRequestOnHandlerNoOp(t *testing.T) {
	// ArrayHandler does not implement RequestAware — should not panic.
	s := newStore()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	s.SetRequestOnHandler(req) // no panic
}

func TestStoreSetExists(t *testing.T) {
	mh := &mockExistenceAwareHandler{ArrayHandler: handlers.NewArrayHandler()}
	s := newStoreWith(mh)

	s.SetExists(true)

	if !mh.exists {
		t.Error("SetExists should forward to ExistenceAware handler")
	}

	s.SetExists(false)

	if mh.exists {
		t.Error("SetExists(false) should update handler")
	}
}

func TestStoreSetExistsNoOp(t *testing.T) {
	// ArrayHandler does not implement ExistenceAware — should not panic.
	s := newStore()
	s.SetExists(true) // no panic
}

func TestStoreIsValidID(t *testing.T) {
	if !session.IsValidID("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa") {
		t.Error("expected valid")
	}

	if session.IsValidID("short") {
		t.Error("expected invalid for short string")
	}

	if session.IsValidID("") {
		t.Error("expected invalid for empty string")
	}

	if session.IsValidID("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz") {
		t.Error("expected invalid for non-hex characters")
	}
}

func TestStoreIDAlias(t *testing.T) {
	s := newStore()

	if s.ID() != s.GetID() {
		t.Error("ID() should return the same as GetID()")
	}
}

func TestStoreHasPreviousURL(t *testing.T) {
	s := newStore()

	if s.HasPreviousURL() {
		t.Error("should not have previous URL initially")
	}

	s.SetPreviousURL("https://example.com")

	if !s.HasPreviousURL() {
		t.Error("should have previous URL after setting it")
	}
}

func TestStoreSetPreviousURL(t *testing.T) {
	s := newStore()

	s.SetPreviousURL("https://example.com")

	if got := s.PreviousURL(); got != "https://example.com" {
		t.Errorf("expected https://example.com, got %v", got)
	}
}

func TestStorePreviousRoute(t *testing.T) {
	s := newStore()

	if s.PreviousRoute() != "" {
		t.Error("expected empty route initially")
	}

	s.SetPreviousRoute("home")

	if got := s.PreviousRoute(); got != "home" {
		t.Errorf("expected home, got %v", got)
	}
}

func TestStoreGetSetName(t *testing.T) {
	s := newStore()

	if s.GetName() != "test" {
		t.Errorf("expected test, got %v", s.GetName())
	}

	s.SetName("changed")

	if s.GetName() != "changed" {
		t.Errorf("expected changed, got %v", s.GetName())
	}
}
