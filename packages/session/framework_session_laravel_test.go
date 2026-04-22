package session_test

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/bedrock/packages/session"
	"github.com/bedrock/packages/session/handlers"
)

// Ports of Laravel framework session store tests from tests/Session.
func TestFrameworkSessionLaravelInventoryCoverage(t *testing.T) {
	ctx := context.Background()

	t.Run("EncryptedSessionStoreTest::testSessionIsProperlyEncrypted", func(t *testing.T) {
		h := handlers.NewArrayHandler()
		enc := &fakeEncrypter{}
		es := session.NewEncryptedStore("test", h, enc)

		if err := es.Start(ctx); err != nil {
			t.Fatalf("Start: %v", err)
		}

		es.Put("secret", "hello")

		if got := es.Get("secret", nil); got != "hello" {
			t.Fatalf("got %v, want hello", got)
		}

		raw := es.Store.Get("secret", nil)
		want := base64.StdEncoding.EncodeToString([]byte("hello"))

		if raw != want {
			t.Fatalf("got %v, want %v", raw, want)
		}
	})

	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "SessionStoreTest::testSessionIsLoadedFromHandler",
			run: func(t *testing.T) {
				h := handlers.NewArrayHandler()
				s1 := session.New("test", h)

				if err := s1.Start(ctx); err != nil {
					t.Fatalf("Start s1: %v", err)
				}

				s1.Put("foo", "bar")

				if err := s1.Save(ctx); err != nil {
					t.Fatalf("Save s1: %v", err)
				}

				s2 := session.NewWithID("test", h, s1.GetID())

				if err := s2.Start(ctx); err != nil {
					t.Fatalf("Start s2: %v", err)
				}

				if got := s2.Get("foo", nil); got != "bar" {
					t.Fatalf("got %v, want bar", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testSessionMigration",
			run: func(t *testing.T) {
				h := handlers.NewArrayHandler()
				s := session.New("test", h)

				if err := s.Start(ctx); err != nil {
					t.Fatalf("Start: %v", err)
				}

				s.Put("k", "v")

				if err := s.Save(ctx); err != nil {
					t.Fatalf("Save: %v", err)
				}

				oldID := s.GetID()

				if err := s.Migrate(ctx, false); err != nil {
					t.Fatalf("Migrate: %v", err)
				}

				if s.GetID() == oldID {
					t.Fatal("expected Migrate to produce a new ID")
				}

				if got := s.Get("k", nil); got != "v" {
					t.Fatalf("got %v, want v", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testSessionRegeneration",
			run: func(t *testing.T) {
				h := handlers.NewArrayHandler()
				s := session.New("test", h)

				if err := s.Start(ctx); err != nil {
					t.Fatalf("Start: %v", err)
				}

				oldID := s.GetID()

				if err := s.Regenerate(ctx, false); err != nil {
					t.Fatalf("Regenerate: %v", err)
				}

				if s.GetID() == oldID {
					t.Fatal("expected Regenerate to produce a new ID")
				}
			},
		},
		{
			name: "SessionStoreTest::testCantSetInvalidId",
			run: func(t *testing.T) {
				s := newStore()

				if err := s.SetID("invalid"); err == nil {
					t.Fatal("expected error for invalid ID")
				}

				valid := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

				if err := s.SetID(valid); err != nil {
					t.Fatalf("SetID valid: %v", err)
				}

				if got := s.GetID(); got != valid {
					t.Fatalf("got %q, want %q", got, valid)
				}
			},
		},
		{
			name: "SessionStoreTest::testSessionInvalidate",
			run: func(t *testing.T) {
				s := newStore()
				s.Put("k", "v")

				if err := s.Invalidate(ctx); err != nil {
					t.Fatalf("Invalidate: %v", err)
				}

				if s.Has("k") {
					t.Fatal("expected session data to be flushed")
				}

				if s.IsStarted() {
					t.Fatal("expected session to be marked not started")
				}
			},
		},
		{
			name: "SessionStoreTest::testBrandNewSessionIsProperlySaved",
			run: func(t *testing.T) {
				h := handlers.NewArrayHandler()
				s1 := session.New("test", h)

				if err := s1.Start(ctx); err != nil {
					t.Fatalf("Start s1: %v", err)
				}

				s1.Put("foo", "bar")
				s1.Flash("baz", "boom")
				s1.Now("qux", "norf")

				if err := s1.Save(ctx); err != nil {
					t.Fatalf("Save s1: %v", err)
				}

				s2 := session.NewWithID("test", h, s1.GetID())

				if err := s2.Start(ctx); err != nil {
					t.Fatalf("Start s2: %v", err)
				}

				if got := s2.Get("foo", nil); got != "bar" {
					t.Fatalf("got %v, want bar", got)
				}

				if got := s2.Get("baz", nil); got != "boom" {
					t.Fatalf("got %v, want boom", got)
				}

				if s2.Has("qux") {
					t.Fatal("Now data should not survive to next request")
				}
			},
		},
		{
			name: "SessionStoreTest::testSessionIsProperlyUpdated",
			run: func(t *testing.T) {
				h := handlers.NewArrayHandler()
				s1 := session.New("test", h)

				if err := s1.Start(ctx); err != nil {
					t.Fatalf("Start s1: %v", err)
				}

				s1.Put("foo", "bar")

				if err := s1.Save(ctx); err != nil {
					t.Fatalf("Save s1: %v", err)
				}

				s2 := session.NewWithID("test", h, s1.GetID())

				if err := s2.Start(ctx); err != nil {
					t.Fatalf("Start s2: %v", err)
				}

				s2.Put("foo", "baz")

				if err := s2.Save(ctx); err != nil {
					t.Fatalf("Save s2: %v", err)
				}

				s3 := session.NewWithID("test", h, s1.GetID())

				if err := s3.Start(ctx); err != nil {
					t.Fatalf("Start s3: %v", err)
				}

				if got := s3.Get("foo", nil); got != "baz" {
					t.Fatalf("got %v, want baz", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testSessionIsReSavedWhenNothingHasChanged",
			run: func(t *testing.T) {
				h := handlers.NewArrayHandler()
				s := session.New("test", h)

				if err := s.Start(ctx); err != nil {
					t.Fatalf("Start: %v", err)
				}

				if err := s.Save(ctx); err != nil {
					t.Fatalf("Save: %v", err)
				}

				s2 := session.NewWithID("test", h, s.GetID())

				if err := s2.Start(ctx); err != nil {
					t.Fatalf("Start s2: %v", err)
				}

				if !s2.IsStarted() {
					t.Fatal("expected session to be started")
				}
			},
		},
		{
			name: "SessionStoreTest::testSessionIsReSavedWhenNothingHasChangedExceptSessionId",
			run: func(t *testing.T) {
				h := handlers.NewArrayHandler()
				s := session.New("test", h)

				if err := s.Start(ctx); err != nil {
					t.Fatalf("Start: %v", err)
				}

				s.Put("foo", "bar")

				if err := s.Regenerate(ctx, false); err != nil {
					t.Fatalf("Regenerate: %v", err)
				}

				if err := s.Save(ctx); err != nil {
					t.Fatalf("Save: %v", err)
				}

				s2 := session.NewWithID("test", h, s.GetID())

				if err := s2.Start(ctx); err != nil {
					t.Fatalf("Start s2: %v", err)
				}

				if got := s2.Get("foo", nil); got != "bar" {
					t.Fatalf("got %v, want bar", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testOldInputFlashing",
			run: func(t *testing.T) {
				h := handlers.NewArrayHandler()
				s1 := session.New("test", h)

				if err := s1.Start(ctx); err != nil {
					t.Fatalf("Start s1: %v", err)
				}

				s1.FlashInput(map[string]any{"name": "Taylor", "email": "taylor@example.com"})

				if err := s1.Save(ctx); err != nil {
					t.Fatalf("Save s1: %v", err)
				}

				s2 := session.NewWithID("test", h, s1.GetID())

				if err := s2.Start(ctx); err != nil {
					t.Fatalf("Start s2: %v", err)
				}

				if !s2.HasOldInput("name") {
					t.Fatal("expected old input for name")
				}

				if got := s2.GetOldInput("name", nil); got != "Taylor" {
					t.Fatalf("got %v, want Taylor", got)
				}

				if got := s2.GetOldInput("email", nil); got != "taylor@example.com" {
					t.Fatalf("got %v, want taylor@example.com", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testDataFlashing",
			run: func(t *testing.T) {
				h := handlers.NewArrayHandler()
				s1 := session.New("test", h)

				if err := s1.Start(ctx); err != nil {
					t.Fatalf("Start s1: %v", err)
				}

				s1.Flash("msg", "hello")

				if err := s1.Save(ctx); err != nil {
					t.Fatalf("Save s1: %v", err)
				}

				s2 := session.NewWithID("test", h, s1.GetID())

				if err := s2.Start(ctx); err != nil {
					t.Fatalf("Start s2: %v", err)
				}

				if got := s2.Get("msg", nil); got != "hello" {
					t.Fatalf("got %v, want hello", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testDataFlashingNow",
			run: func(t *testing.T) {
				s := newStore()
				s.Now("foo", "bar")

				if got := s.Get("foo", nil); got != "bar" {
					t.Fatalf("got %v, want bar", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testDataMergeNewFlashes",
			run: func(t *testing.T) {
				h := handlers.NewArrayHandler()
				s1 := session.New("test", h)

				if err := s1.Start(ctx); err != nil {
					t.Fatalf("Start s1: %v", err)
				}

				s1.Flash("foo", "bar")

				if err := s1.Save(ctx); err != nil {
					t.Fatalf("Save s1: %v", err)
				}

				s2 := session.NewWithID("test", h, s1.GetID())

				if err := s2.Start(ctx); err != nil {
					t.Fatalf("Start s2: %v", err)
				}

				s2.Keep("foo")

				if err := s2.Save(ctx); err != nil {
					t.Fatalf("Save s2: %v", err)
				}

				s3 := session.NewWithID("test", h, s1.GetID())

				if err := s3.Start(ctx); err != nil {
					t.Fatalf("Start s3: %v", err)
				}

				if got := s3.Get("foo", nil); got != "bar" {
					t.Fatalf("got %v, want bar", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testReflash",
			run: func(t *testing.T) {
				h := handlers.NewArrayHandler()
				s1 := session.New("test", h)

				if err := s1.Start(ctx); err != nil {
					t.Fatalf("Start s1: %v", err)
				}

				s1.Flash("x", 1)

				if err := s1.Save(ctx); err != nil {
					t.Fatalf("Save s1: %v", err)
				}

				s2 := session.NewWithID("test", h, s1.GetID())

				if err := s2.Start(ctx); err != nil {
					t.Fatalf("Start s2: %v", err)
				}

				s2.Reflash()

				if err := s2.Save(ctx); err != nil {
					t.Fatalf("Save s2: %v", err)
				}

				s3 := session.NewWithID("test", h, s1.GetID())

				if err := s3.Start(ctx); err != nil {
					t.Fatalf("Start s3: %v", err)
				}

				if !s3.Has("x") {
					t.Fatal("expected reflash data to survive one more request")
				}
			},
		},
		{
			name: "SessionStoreTest::testReflashWithNow",
			run: func(t *testing.T) {
				h := handlers.NewArrayHandler()
				s1 := session.New("test", h)

				if err := s1.Start(ctx); err != nil {
					t.Fatalf("Start s1: %v", err)
				}

				s1.Now("foo", "bar")
				s1.Reflash()

				if err := s1.Save(ctx); err != nil {
					t.Fatalf("Save s1: %v", err)
				}

				s2 := session.NewWithID("test", h, s1.GetID())

				if err := s2.Start(ctx); err != nil {
					t.Fatalf("Start s2: %v", err)
				}

				if got := s2.Get("foo", nil); got != "bar" {
					t.Fatalf("got %v, want bar", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testOnly",
			run: func(t *testing.T) {
				s := newStore()
				s.Put("a", 1)
				s.Put("b", 2)
				s.Put("c", 3)

				got := s.Only("a", "c")

				if len(got) != 2 || got["a"] != 1 || got["c"] != 3 {
					t.Fatalf("unexpected Only result: %#v", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testExcept",
			run: func(t *testing.T) {
				s := newStore()
				s.Put("a", 1)
				s.Put("b", 2)
				s.Put("c", 3)

				got := s.Except("b")

				if _, ok := got["b"]; ok {
					t.Fatal("expected b to be excluded")
				}
			},
		},
		{
			name: "SessionStoreTest::testReplace",
			run: func(t *testing.T) {
				s := newStore()
				s.Put("a", 1)
				s.Put("b", 2)
				s.Replace(map[string]any{"b": 99, "c": 3})

				if got := s.Get("b", nil); got != 99 {
					t.Fatalf("got %v, want 99", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testRemove",
			run: func(t *testing.T) {
				s := newStore()
				s.Put("k", "v")

				if got := s.Remove("k"); got != "v" {
					t.Fatalf("got %v, want v", got)
				}

				if s.Exists("k") {
					t.Fatal("expected key to be removed")
				}
			},
		},
		{
			name: "SessionStoreTest::testClear",
			run: func(t *testing.T) {
				s := newStore()
				s.Put("a", 1)
				s.Put("b", 2)
				s.Flush()

				if len(s.All()) != 0 {
					t.Fatal("expected all attributes to be cleared")
				}
			},
		},
		{
			name: "SessionStoreTest::testIncrement",
			run: func(t *testing.T) {
				s := newStore()
				s.Put("n", int64(10))

				if got := s.Increment("n", 5); got != 15 {
					t.Fatalf("got %d, want 15", got)
				}

				if got := s.Decrement("n", 3); got != 12 {
					t.Fatalf("got %d, want 12", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testDecrement",
			run: func(t *testing.T) {
				s := newStore()

				if got := s.Decrement("counter", 1); got != -1 {
					t.Fatalf("got %d, want -1", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testHasOldInputWithoutKey",
			run: func(t *testing.T) {
				h := handlers.NewArrayHandler()
				s1 := session.New("test", h)

				if err := s1.Start(ctx); err != nil {
					t.Fatalf("Start s1: %v", err)
				}

				s1.FlashInput(map[string]any{"name": "Taylor"})

				if err := s1.Save(ctx); err != nil {
					t.Fatalf("Save s1: %v", err)
				}

				s2 := session.NewWithID("test", h, s1.GetID())

				if err := s2.Start(ctx); err != nil {
					t.Fatalf("Start s2: %v", err)
				}

				if !s2.HasOldInput("") {
					t.Fatal("expected HasOldInput(\"\") to report true")
				}
			},
		},
		{
			name: "SessionStoreTest::testHandlerNeedsRequest",
			run: func(t *testing.T) {
				s := newStore()

				if s.HandlerNeedsRequest() {
					t.Fatal("ArrayHandler should not need a request")
				}

				mh := &mockRequestAwareHandler{ArrayHandler: handlers.NewArrayHandler()}
				s2 := newStoreWith(mh)

				if !s2.HandlerNeedsRequest() {
					t.Fatal("expected request-aware handler to be detected")
				}
			},
		},
		{
			name: "SessionStoreTest::testToken",
			run: func(t *testing.T) {
				s := newStore()
				tok := s.Token()

				if tok == "" {
					t.Fatal("expected a token")
				}

				if s.Token() != tok {
					t.Fatal("expected token to remain stable")
				}
			},
		},
		{
			name: "SessionStoreTest::testRegenerateToken",
			run: func(t *testing.T) {
				s := newStore()
				tok := s.Token()
				s.RegenerateToken()

				if s.Token() == tok {
					t.Fatal("expected token to change")
				}
			},
		},
		{
			name: "SessionStoreTest::testName",
			run: func(t *testing.T) {
				s := newStore()

				if got := s.GetName(); got != "test" {
					t.Fatalf("got %q, want test", got)
				}

				s.SetName("changed")

				if got := s.GetName(); got != "changed" {
					t.Fatalf("got %q, want changed", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testForget",
			run: func(t *testing.T) {
				s := newStore()
				s.Put("a", 1)
				s.Put("b", 2)
				s.Put("c", 3)
				s.Forget("a", "c")

				if s.Exists("a") || s.Exists("c") || !s.Exists("b") {
					t.Fatal("unexpected forget result")
				}
			},
		},
		{
			name: "SessionStoreTest::testSetPreviousUrl",
			run: func(t *testing.T) {
				s := newStore()
				s.SetPreviousURL("https://example.com")

				if got := s.PreviousURL(); got != "https://example.com" {
					t.Fatalf("got %q, want https://example.com", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testPasswordConfirmed",
			run: func(t *testing.T) {
				s := newStore()

				if before := s.PasswordConfirmedAt(); before != 0 {
					t.Fatalf("got %d, want 0", before)
				}

				s.PasswordConfirmed()

				if after := s.PasswordConfirmedAt(); after == 0 {
					t.Fatal("expected password confirmation timestamp")
				}
			},
		},
		{
			name: "SessionStoreTest::testKeyPush",
			run: func(t *testing.T) {
				s := newStore()
				s.Push("list", "a")
				s.Push("list", "b")
				got := s.Get("list", nil)
				values, ok := got.([]any)

				if !ok || len(values) != 2 || values[0] != "a" || values[1] != "b" {
					t.Fatalf("unexpected pushed value: %#v", got)
				}
			},
		},
		{
			name: "SessionStoreTest::testKeyPull",
			run: func(t *testing.T) {
				s := newStore()
				s.Put("k", "v")

				if got := s.Pull("k", nil); got != "v" {
					t.Fatalf("got %v, want v", got)
				}

				if s.Exists("k") {
					t.Fatal("expected key to be removed")
				}
			},
		},
		{
			name: "SessionStoreTest::testKeyHas",
			run: func(t *testing.T) {
				s := newStore()
				s.Put("a", "val")
				s.Put("b", nil)

				if !s.Has("a") || s.Has("b") || s.Has("c") {
					t.Fatal("unexpected Has result")
				}
			},
		},
		{
			name: "SessionStoreTest::testKeyHasAny",
			run: func(t *testing.T) {
				s := newStore()
				s.Put("a", "val")

				if !s.HasAny("a", "b") {
					t.Fatal("expected HasAny to be true")
				}

				if s.HasAny("b", "c") {
					t.Fatal("expected HasAny to be false")
				}
			},
		},
		{
			name: "SessionStoreTest::testKeyExists",
			run: func(t *testing.T) {
				s := newStore()
				s.Put("a", "val")
				s.Put("b", nil)

				if !s.Exists("a") || !s.Exists("b") || s.Exists("c") {
					t.Fatal("unexpected Exists result")
				}
			},
		},
		{
			name: "SessionStoreTest::testKeyMissing",
			run: func(t *testing.T) {
				s := newStore()
				s.Put("a", "val")

				if s.Missing("a") || !s.Missing("b") {
					t.Fatal("unexpected Missing result")
				}
			},
		},
		{
			name: "SessionStoreTest::testRememberMethodCallsPutAndReturnsDefault",
			run: func(t *testing.T) {
				s := newStore()
				called := false
				got := s.Remember("key", func() any {
					called = true

					return "computed"
				})

				if !called || got != "computed" || s.Get("key", nil) != "computed" {
					t.Fatal("unexpected Remember result")
				}
			},
		},
		{
			name: "SessionStoreTest::testRememberMethodReturnsPreviousValueIfItAlreadySets",
			run: func(t *testing.T) {
				s := newStore()
				s.Put("key", "existing")
				called := false
				got := s.Remember("key", func() any {
					called = true

					return "new"
				})

				if called || got != "existing" {
					t.Fatal("unexpected Remember reuse result")
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, tc.run)
	}
}
