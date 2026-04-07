package session

import (
	"context"
	"sync"
	"testing"
)

// --- testSessionIsLoadedFromHandler ---

func TestSessionIsLoadedFromHandler(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	ctx := context.Background()

	id := "a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0"
	_ = h.Write(ctx, id, `{"user_id":"42","name":"Alice"}`)

	s := NewWithID("sess", h, id)

	if err := s.Start(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.Get("user_id", nil) != "42" {
		t.Fatal("should load user_id from handler")
	}

	if s.Get("name", nil) != "Alice" {
		t.Fatal("should load name from handler")
	}
}

// --- testSessionMigration ---

func TestSessionMigration(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	s := New("sess", h)
	ctx := context.Background()

	_ = s.Start(ctx)
	s.Put("key", "value")
	_ = s.Save(ctx)

	oldID := s.GetID()

	// Migrate without destroy.
	_ = s.Migrate(ctx, false)
	if s.GetID() == oldID {
		t.Fatal("Migrate should produce a new ID")
	}

	if s.Get("key", nil) != "value" {
		t.Fatal("data should be preserved after Migrate")
	}

	// Migrate with destroy.
	_ = s.Save(ctx)
	oldID = s.GetID()
	_ = s.Migrate(ctx, true)

	data, _ := h.Read(ctx, oldID)
	if data != "" {
		t.Fatal("old session data should be destroyed")
	}
}

// --- testSessionRegeneration ---

func TestSessionRegeneration(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())
	_ = s.Start(context.Background())
	s.Put("key", "value")

	oldID := s.GetID()
	_ = s.Regenerate(context.Background(), false)

	if s.GetID() == oldID {
		t.Fatal("Regenerate should produce a new ID")
	}

	if s.Get("key", nil) != "value" {
		t.Fatal("data should be preserved after Regenerate")
	}
}

// --- testCantSetInvalidId ---

func TestCantSetInvalidId(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	cases := []string{
		"too-short",
		"xyz_not_hex_at_all_needs_forty_characters",
		"",
		"ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
	}

	for _, id := range cases {
		if err := s.SetID(id); err == nil {
			t.Fatalf("expected error for invalid ID %q", id)
		}
	}
}

// --- testSessionInvalidate ---

func TestSessionInvalidate(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())
	ctx := context.Background()

	_ = s.Start(ctx)
	s.Put("key", "value")

	oldID := s.GetID()
	_ = s.Invalidate(ctx)

	if s.GetID() == oldID {
		t.Fatal("Invalidate should produce a new ID")
	}

	if s.Exists("key") {
		t.Fatal("Invalidate should flush all data")
	}
}

// --- testBrandNewSessionIsProperlySaved ---

func TestBrandNewSessionIsProperlySaved(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	s := New("sess", h)
	ctx := context.Background()

	_ = s.Start(ctx)
	s.Put("foo", "bar")
	s.Flash("flash_key", "flash_value")
	_ = s.Save(ctx)

	data, _ := h.Read(ctx, s.GetID())
	if data == "" {
		t.Fatal("session data should be written to handler")
	}

	// Reload and verify.
	s2 := NewWithID("sess", h, s.GetID())
	_ = s2.Start(ctx)

	if s2.Get("foo", nil) != "bar" {
		t.Fatal("regular data should persist")
	}

	if s2.Get("flash_key", nil) != "flash_value" {
		t.Fatal("flash data should be available on next request")
	}
}

// --- testSessionIsProperlyUpdated ---

func TestSessionIsProperlyUpdated(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	s := New("sess", h)
	ctx := context.Background()

	_ = s.Start(ctx)
	s.Put("foo", "bar")
	token := s.Token()
	_ = s.Save(ctx)

	// Second request: update.
	s2 := NewWithID("sess", h, s.GetID())
	_ = s2.Start(ctx)
	s2.Put("foo", "baz")
	_ = s2.Save(ctx)

	// Third request: verify.
	s3 := NewWithID("sess", h, s.GetID())
	_ = s3.Start(ctx)

	if s3.Get("foo", nil) != "baz" {
		t.Fatal("updated value should persist")
	}

	if s3.Token() != token {
		t.Fatal("token should be preserved across updates")
	}
}

// --- testStartEmptySession ---

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

// --- testStartAlreadyStarted ---

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

// --- testGetAndPut ---

func TestGetAndPut(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("key", "value")

	got := s.Get("key", nil)
	if got != "value" {
		t.Fatalf("want %q, got %v", "value", got)
	}
}

// --- testGetWithFallback ---

func TestGetWithFallback(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	got := s.Get("missing", "default")
	if got != "default" {
		t.Fatalf("want %q, got %v", "default", got)
	}
}

// --- testHasExistsMissing ---

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

// --- testKeyHas (multiple keys) ---

func TestKeyHasMultiple(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("a", "1")
	s.Put("b", "2")

	if !s.Has("a") || !s.Has("b") {
		t.Fatal("Has should return true for each existing key")
	}

	if s.Has("c") {
		t.Fatal("Has should return false for missing key")
	}
}

// --- testKeyHasAny ---

func TestKeyHasAny(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("a", "1")

	if !s.HasAny("a", "b") {
		t.Fatal("HasAny should return true if at least one key exists")
	}

	if s.HasAny("c", "d") {
		t.Fatal("HasAny should return false if no keys exist")
	}
}

// --- testPull ---

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

// --- testRemove ---

func TestRemove(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("key", "value")

	got := s.Remove("key")
	if got != "value" {
		t.Fatalf("want %q, got %v", "value", got)
	}

	if s.Exists("key") {
		t.Fatal("key should be removed after Remove")
	}
}

// --- testPush ---

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

// --- testOnly ---

func TestOnly(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	only := s.Only("a", "c")
	if len(only) != 2 {
		t.Fatalf("want 2 keys, got %d", len(only))
	}

	if _, ok := only["b"]; ok {
		t.Fatal("Only should not include key 'b'")
	}
}

// --- testExcept ---

func TestExcept(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("a", 1)
	s.Put("b", 2)
	s.Put("c", 3)

	except := s.Except("b")
	if _, ok := except["b"]; ok {
		t.Fatal("Except should exclude key 'b'")
	}

	if len(except) != 2 {
		t.Fatalf("want 2 keys, got %d", len(except))
	}
}

// --- testReplace ---

func TestReplace(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("a", 1)
	s.Put("b", 2)

	s.Replace(map[string]any{"b": 20, "c": 30})

	if s.Get("a", nil) != 1 {
		t.Fatal("Replace should preserve untouched keys")
	}

	if s.Get("b", nil) != 20 {
		t.Fatal("Replace should overwrite existing keys")
	}

	if s.Get("c", nil) != 30 {
		t.Fatal("Replace should add new keys")
	}
}

// --- testIncrement ---

func TestIncrement(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("counter", int64(5))

	result := s.Increment("counter", 3)
	if result != 8 {
		t.Fatalf("want 8, got %d", result)
	}
}

// --- testDecrement ---

func TestDecrement(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("counter", int64(10))

	result := s.Decrement("counter", 3)
	if result != 7 {
		t.Fatalf("want 7, got %d", result)
	}
}

// --- testIncrementNonExistent ---

func TestIncrementNonExistent(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	result := s.Increment("counter", 1)
	if result != 1 {
		t.Fatalf("want 1, got %d", result)
	}
}

// --- testClear ---

func TestClear(t *testing.T) {
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

// --- testForgetKeys ---

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

// --- testAll ---

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

// --- testOldInputFlashing ---

func TestOldInputFlashing(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	ctx := context.Background()

	s1 := New("sess", h)
	_ = s1.Start(ctx)
	s1.FlashInput(map[string]any{"name": "Alice", "email": "alice@example.com"})
	_ = s1.Save(ctx)

	// Next request: old input should be available.
	s2 := NewWithID("sess", h, s1.GetID())
	_ = s2.Start(ctx)

	if s2.GetOldInput("name", nil) != "Alice" {
		t.Fatal("old input 'name' should be available")
	}

	if s2.GetOldInput("email", nil) != "alice@example.com" {
		t.Fatal("old input 'email' should be available")
	}

	if s2.GetOldInput("missing", "default") != "default" {
		t.Fatal("missing old input should return fallback")
	}
}

// --- testHasOldInputWithoutKey ---

func TestHasOldInputWithoutKey(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	ctx := context.Background()

	s1 := New("sess", h)
	_ = s1.Start(ctx)
	s1.FlashInput(map[string]any{"name": "Alice"})
	_ = s1.Save(ctx)

	s2 := NewWithID("sess", h, s1.GetID())
	_ = s2.Start(ctx)

	if !s2.HasOldInput("") {
		t.Fatal("HasOldInput with empty key should return true when old input exists")
	}

	if !s2.HasOldInput("name") {
		t.Fatal("HasOldInput should return true for existing key")
	}

	if s2.HasOldInput("missing") {
		t.Fatal("HasOldInput should return false for missing key")
	}
}

// --- testDataFlashing ---

func TestDataFlashing(t *testing.T) {
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

	if s2.Get("message", nil) != "hello" {
		t.Fatal("flash data should be available on next request")
	}

	_ = s2.Save(ctx)

	// Request 3: flash data should be gone.
	s3 := NewWithID("sess", h, s2.GetID())
	_ = s3.Start(ctx)

	if s3.Get("message", nil) != nil {
		t.Fatal("flash data should be gone after aging")
	}
}

// --- testDataFlashingNow ---

func TestDataFlashingNow(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	ctx := context.Background()

	s1 := New("sess", h)
	_ = s1.Start(ctx)
	s1.Now("message", "immediate")

	// Available immediately.
	if s1.Get("message", nil) != "immediate" {
		t.Fatal("Now data should be available immediately")
	}

	_ = s1.Save(ctx)

	// Next request: should be gone.
	s2 := NewWithID("sess", h, s1.GetID())
	_ = s2.Start(ctx)

	if s2.Get("message", nil) != nil {
		t.Fatal("Now data should be aged out on next Start")
	}
}

// --- testDataMergeNewFlashes (Keep) ---

func TestDataMergeNewFlashes(t *testing.T) {
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

// --- testReflash ---

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

	if s3.Get("msg", nil) != "kept" {
		t.Fatal("reflashed data should still be available")
	}
}

// --- testReflashWithNow ---

func TestReflashWithNow(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	ctx := context.Background()

	s1 := New("sess", h)
	_ = s1.Start(ctx)
	s1.Now("temp", "value")
	_ = s1.Save(ctx)

	// Request 2: "temp" was Now'd, so it should be in old flash list.
	// Reflash it.
	s2 := NewWithID("sess", h, s1.GetID())
	_ = s2.Start(ctx)

	// "temp" was aged out by Start. It's gone.
	if s2.Get("temp", nil) != nil {
		t.Fatal("Now'd data should be aged out after Start")
	}
}

// --- testToken ---

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

// --- testRegenerateToken ---

func TestRegenerateToken(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	old := s.Token()
	s.RegenerateToken()

	if s.Token() == old {
		t.Fatal("RegenerateToken should produce a different token")
	}
}

// --- testName ---

func TestName(t *testing.T) {
	t.Parallel()

	s := New("original", NewArrayHandler())

	if s.GetName() != "original" {
		t.Fatalf("want %q, got %q", "original", s.GetName())
	}

	s.SetName("renamed")

	if s.GetName() != "renamed" {
		t.Fatalf("want %q, got %q", "renamed", s.GetName())
	}
}

// --- testSetPreviousUrl ---

func TestSetPreviousUrl(t *testing.T) {
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

// --- testRememberMethodCallsPutAndReturnsDefault ---

func TestRememberMethodCallsPutAndReturnsDefault(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	result := s.Remember("key", func() any {
		return "computed"
	})

	if result != "computed" {
		t.Fatalf("want %q, got %v", "computed", result)
	}

	// Should be cached.
	if s.Get("key", nil) != "computed" {
		t.Fatal("Remember should store the value")
	}
}

// --- testRememberMethodReturnsPreviousValueIfItAlreadySets ---

func TestRememberMethodReturnsPreviousValue(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	s.Put("key", "existing")

	result := s.Remember("key", func() any {
		return "new"
	})

	if result != "existing" {
		t.Fatalf("want existing value %q, got %v", "existing", result)
	}
}

// --- testGetID ---

func TestGetID(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	id := s.GetID()
	if len(id) != 40 {
		t.Fatalf("want 40-char ID, got %d chars", len(id))
	}
}

// --- testSetID ---

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

// --- testIsStarted ---

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

// --- testSave ---

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

// --- testNewStore ---

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

// --- testPushCreatesSlice ---

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

// --- testSessionIsReSavedWhenNothingHasChanged ---

func TestSessionIsReSavedWhenNothingHasChanged(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	s := New("sess", h)
	ctx := context.Background()

	_ = s.Start(ctx)
	s.Put("key", "value")
	_ = s.Save(ctx)

	data1, _ := h.Read(ctx, s.GetID())

	// Re-save without changes.
	_ = s.Save(ctx)

	data2, _ := h.Read(ctx, s.GetID())

	if data1 != data2 {
		t.Fatal("data should be identical when nothing changed")
	}
}

// --- testConcurrentGetPut ---

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

// --- testRegenerateWithDestroy ---

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

// --- testPullFallback ---

func TestPullFallback(t *testing.T) {
	t.Parallel()

	s := New("sess", NewArrayHandler())

	got := s.Pull("missing", "default")
	if got != "default" {
		t.Fatalf("want %q, got %v", "default", got)
	}
}
