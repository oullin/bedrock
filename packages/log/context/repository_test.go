package context_test

import (
	"testing"

	logctx "github.com/bedrock/packages/log/context"
)

func TestRepositoryGetSet(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("name", "john")

	if r.Get("name") != "john" {
		t.Fatalf("expected 'john', got %v", r.Get("name"))
	}
}

func TestRepositoryGetWithFallback(t *testing.T) {
	t.Parallel()

	r := logctx.New()

	if r.Get("missing", "default") != "default" {
		t.Fatalf("expected 'default', got %v", r.Get("missing", "default"))
	}
}

func TestRepositoryGetReturnsNilForUnset(t *testing.T) {
	t.Parallel()

	r := logctx.New()

	if r.Get("missing") != nil {
		t.Fatalf("expected nil, got %v", r.Get("missing"))
	}
}

func TestRepositoryHas(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("key", nil)

	if !r.Has("key") {
		t.Fatal("expected Has to return true for nil value")
	}

	if r.Has("other") {
		t.Fatal("expected Has to return false for unset key")
	}
}

func TestRepositoryHasMultipleKeys(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("a", 1).Add("b", 2)

	if !r.Has("a", "b") {
		t.Fatal("expected Has to return true for both keys")
	}

	if r.Has("a", "c") {
		t.Fatal("expected Has to return false when one key is missing")
	}
}

func TestRepositoryMissing(t *testing.T) {
	t.Parallel()

	r := logctx.New()

	if !r.Missing("key") {
		t.Fatal("expected Missing to return true for unset key")
	}

	r.Add("key", "value")

	if r.Missing("key") {
		t.Fatal("expected Missing to return false for set key")
	}
}

func TestRepositoryAdd(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("str", "hello").Add("int", 42).Add("bool", true).Add("nil", nil)

	if r.Get("str") != "hello" {
		t.Fatalf("expected 'hello', got %v", r.Get("str"))
	}

	if r.Get("int") != 42 {
		t.Fatalf("expected 42, got %v", r.Get("int"))
	}

	if r.Get("bool") != true {
		t.Fatalf("expected true, got %v", r.Get("bool"))
	}

	if r.Get("nil") != nil {
		t.Fatalf("expected nil, got %v", r.Get("nil"))
	}
}

func TestRepositoryAddIf(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("key", "original")
	r.AddIf("key", "replacement")

	if r.Get("key") != "original" {
		t.Fatalf("expected 'original', got %v", r.Get("key"))
	}

	r.AddIf("new", "value")

	if r.Get("new") != "value" {
		t.Fatalf("expected 'value', got %v", r.Get("new"))
	}
}

func TestRepositoryForget(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("a", 1).Add("b", 2).Add("c", 3)
	r.Forget("a", "c")

	if r.Has("a") {
		t.Fatal("expected 'a' to be forgotten")
	}

	if !r.Has("b") {
		t.Fatal("expected 'b' to remain")
	}

	if r.Has("c") {
		t.Fatal("expected 'c' to be forgotten")
	}
}

func TestRepositoryPull(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("key", "value")

	v := r.Pull("key")
	if v != "value" {
		t.Fatalf("expected 'value', got %v", v)
	}

	if r.Has("key") {
		t.Fatal("expected key to be removed after pull")
	}
}

func TestRepositoryPullWithFallback(t *testing.T) {
	t.Parallel()

	r := logctx.New()

	v := r.Pull("missing", "default")
	if v != "default" {
		t.Fatalf("expected 'default', got %v", v)
	}
}

func TestRepositoryOnly(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("a", 1).Add("b", 2).Add("c", 3)

	result := r.Only("a", "c")

	if result["a"] != 1 {
		t.Fatalf("expected a = 1, got %v", result["a"])
	}

	if _, ok := result["b"]; ok {
		t.Fatal("expected 'b' to be excluded")
	}

	if result["c"] != 3 {
		t.Fatalf("expected c = 3, got %v", result["c"])
	}
}

func TestRepositoryExcept(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("a", 1).Add("b", 2).Add("c", 3)

	result := r.Except("b")

	if result["a"] != 1 {
		t.Fatalf("expected a = 1, got %v", result["a"])
	}

	if _, ok := result["b"]; ok {
		t.Fatal("expected 'b' to be excluded")
	}

	if result["c"] != 3 {
		t.Fatalf("expected c = 3, got %v", result["c"])
	}
}

func TestRepositoryRemember(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	calls := 0

	v1 := r.Remember("key", func() any {
		calls++
		return "computed"
	})

	v2 := r.Remember("key", func() any {
		calls++
		return "recomputed"
	})

	if v1 != "computed" {
		t.Fatalf("expected 'computed', got %v", v1)
	}

	if v2 != "computed" {
		t.Fatalf("expected cached 'computed', got %v", v2)
	}

	if calls != 1 {
		t.Fatalf("expected fn called once, called %d times", calls)
	}
}

func TestRepositoryRememberExistingKey(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("key", "existing")

	v := r.Remember("key", func() any { return "new" })

	if v != "existing" {
		t.Fatalf("expected 'existing', got %v", v)
	}
}

func TestRepositoryAll(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("a", 1).Add("b", 2)

	all := r.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}

	if all["a"] != 1 || all["b"] != 2 {
		t.Fatalf("expected {a:1, b:2}, got %v", all)
	}
}

func TestRepositoryFlush(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("a", 1).AddHidden("b", 2)
	r.Flush()

	if !r.IsEmpty() {
		t.Fatal("expected repository to be empty after flush")
	}
}

// Hidden data tests

func TestRepositoryGetHidden(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.AddHidden("secret", "password")

	if r.GetHidden("secret") != "password" {
		t.Fatalf("expected 'password', got %v", r.GetHidden("secret"))
	}
}

func TestRepositoryGetHiddenWithFallback(t *testing.T) {
	t.Parallel()

	r := logctx.New()

	if r.GetHidden("missing", "default") != "default" {
		t.Fatalf("expected 'default', got %v", r.GetHidden("missing", "default"))
	}
}

func TestRepositoryHasHidden(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.AddHidden("key", "value")

	if !r.HasHidden("key") {
		t.Fatal("expected HasHidden to return true")
	}

	if r.HasHidden("other") {
		t.Fatal("expected HasHidden to return false for unset key")
	}
}

func TestRepositoryAddHidden(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.AddHidden("token", "abc123")

	if r.GetHidden("token") != "abc123" {
		t.Fatalf("expected 'abc123', got %v", r.GetHidden("token"))
	}

	all := r.All()
	if _, ok := all["token"]; ok {
		t.Fatal("expected hidden data to not appear in All()")
	}
}

func TestRepositoryAddHiddenIf(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.AddHidden("key", "original")
	r.AddHiddenIf("key", "replacement")

	if r.GetHidden("key") != "original" {
		t.Fatalf("expected 'original', got %v", r.GetHidden("key"))
	}
}

func TestRepositoryForgetHidden(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.AddHidden("a", 1).AddHidden("b", 2)
	r.ForgetHidden("a")

	if r.HasHidden("a") {
		t.Fatal("expected 'a' to be forgotten")
	}

	if !r.HasHidden("b") {
		t.Fatal("expected 'b' to remain")
	}
}

func TestRepositoryPullHidden(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.AddHidden("key", "value")

	v := r.PullHidden("key")
	if v != "value" {
		t.Fatalf("expected 'value', got %v", v)
	}

	if r.HasHidden("key") {
		t.Fatal("expected key to be removed after pull")
	}
}

func TestRepositoryOnlyHidden(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.AddHidden("a", 1).AddHidden("b", 2).AddHidden("c", 3)

	result := r.OnlyHidden("a", "c")
	if result["a"] != 1 || result["c"] != 3 {
		t.Fatalf("expected {a:1, c:3}, got %v", result)
	}

	if _, ok := result["b"]; ok {
		t.Fatal("expected 'b' to be excluded")
	}
}

func TestRepositoryExceptHidden(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.AddHidden("a", 1).AddHidden("b", 2).AddHidden("c", 3)

	result := r.ExceptHidden("b")
	if _, ok := result["b"]; ok {
		t.Fatal("expected 'b' to be excluded")
	}

	if result["a"] != 1 || result["c"] != 3 {
		t.Fatalf("expected {a:1, c:3}, got %v", result)
	}
}

func TestRepositoryRememberHidden(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	calls := 0

	v := r.RememberHidden("key", func() any {
		calls++
		return "secret"
	})

	if v != "secret" {
		t.Fatalf("expected 'secret', got %v", v)
	}

	v2 := r.RememberHidden("key", func() any {
		calls++
		return "new"
	})

	if v2 != "secret" {
		t.Fatalf("expected cached 'secret', got %v", v2)
	}

	if calls != 1 {
		t.Fatalf("expected fn called once, called %d times", calls)
	}
}

func TestRepositoryAllHidden(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.AddHidden("a", 1).AddHidden("b", 2)

	all := r.AllHidden()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestRepositoryAllHiddenEmpty(t *testing.T) {
	t.Parallel()

	r := logctx.New()

	all := r.AllHidden()
	if len(all) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(all))
	}
}

// Stack operations

func TestRepositoryPush(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Push("items", "a", "b")
	r.Push("items", "c")

	items := r.Get("items").([]any)
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}

	if items[0] != "a" || items[1] != "b" || items[2] != "c" {
		t.Fatalf("expected [a, b, c], got %v", items)
	}
}

func TestRepositoryPop(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Push("items", "a", "b", "c")

	v := r.Pop("items")
	if v != "c" {
		t.Fatalf("expected 'c', got %v", v)
	}

	items := r.Get("items").([]any)
	if len(items) != 2 {
		t.Fatalf("expected 2 items after pop, got %d", len(items))
	}
}

func TestRepositoryPopEmpty(t *testing.T) {
	t.Parallel()

	r := logctx.New()

	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatal("expected panic on empty pop")
		}
	}()

	r.Pop("missing")
}

func TestRepositoryPushToNonSlice(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("key", "string_value")

	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatal("expected panic when pushing to non-slice")
		}
	}()

	r.Push("key", "item")
}

func TestRepositoryStackContains(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Push("items", "a", "b", "c")

	if !r.StackContains("items", "b") {
		t.Fatal("expected stack to contain 'b'")
	}

	if r.StackContains("items", "d") {
		t.Fatal("expected stack to not contain 'd'")
	}
}

func TestRepositoryStackContainsFunc(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Push("nums", 1, 2, 3)

	found := r.StackContainsFunc("nums", func(v any) bool {
		n, ok := v.(int)
		return ok && n > 2
	})

	if !found {
		t.Fatal("expected stack to contain value > 2")
	}
}

func TestRepositoryHiddenStackContains(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.PushHidden("tokens", "abc", "def")

	if !r.HiddenStackContains("tokens", "abc") {
		t.Fatal("expected hidden stack to contain 'abc'")
	}

	if r.StackContains("tokens", "abc") {
		t.Fatal("expected hidden values to not appear in non-hidden stack check")
	}
}

func TestRepositoryHiddenStackContainsFunc(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.PushHidden("nums", 10, 20)

	found := r.HiddenStackContainsFunc("nums", func(v any) bool {
		n, ok := v.(int)
		return ok && n == 20
	})

	if !found {
		t.Fatal("expected hidden stack to contain 20")
	}
}

func TestRepositoryPopHidden(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.PushHidden("items", "x", "y")

	v := r.PopHidden("items")
	if v != "y" {
		t.Fatalf("expected 'y', got %v", v)
	}
}

func TestRepositoryPopHiddenEmpty(t *testing.T) {
	t.Parallel()

	r := logctx.New()

	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatal("expected panic on empty hidden pop")
		}
	}()

	r.PopHidden("missing")
}

// Counter operations

func TestRepositoryIncrement(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Increment("count")

	if r.Get("count") != 1 {
		t.Fatalf("expected 1, got %v", r.Get("count"))
	}

	r.Increment("count")

	if r.Get("count") != 2 {
		t.Fatalf("expected 2, got %v", r.Get("count"))
	}
}

func TestRepositoryIncrementBy(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Increment("count", 5)

	if r.Get("count") != 5 {
		t.Fatalf("expected 5, got %v", r.Get("count"))
	}

	r.Increment("count", 3)

	if r.Get("count") != 8 {
		t.Fatalf("expected 8, got %v", r.Get("count"))
	}
}

func TestRepositoryDecrement(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("count", 10)
	r.Decrement("count")

	if r.Get("count") != 9 {
		t.Fatalf("expected 9, got %v", r.Get("count"))
	}
}

func TestRepositoryDecrementBy(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("count", 10)
	r.Decrement("count", 3)

	if r.Get("count") != 7 {
		t.Fatalf("expected 7, got %v", r.Get("count"))
	}
}

// Scoped execution

func TestRepositoryScope(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("a", 1)

	r.Scope(func(scoped *logctx.Repository) {
		scoped.Add("b", 2)

		if scoped.Get("a") != 1 {
			t.Fatal("expected scoped to inherit parent data")
		}

		if scoped.Get("b") != 2 {
			t.Fatal("expected scoped to have new data")
		}
	})

	if r.Has("b") {
		t.Fatal("expected scope changes to not leak to parent")
	}
}

func TestRepositoryScopeDoesNotLeakToParent(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("key", "original")

	r.Scope(func(scoped *logctx.Repository) {
		scoped.Add("key", "modified")
	})

	if r.Get("key") != "original" {
		t.Fatalf("expected 'original', got %v", r.Get("key"))
	}
}

func TestRepositoryScopeWithExtraData(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("a", 1)

	r.Scope(func(scoped *logctx.Repository) {
		if scoped.Get("extra") != "data" {
			t.Fatal("expected scope to receive extra data")
		}
	}, map[string]any{"extra": "data"})
}

// Serialization

func TestRepositoryDehydrate(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("a", 1).Add("b", "two")

	data := r.Dehydrate()
	if data["a"] != 1 || data["b"] != "two" {
		t.Fatalf("expected {a:1, b:two}, got %v", data)
	}
}

func TestRepositoryDehydrateEmpty(t *testing.T) {
	t.Parallel()

	r := logctx.New()

	data := r.Dehydrate()
	if data != nil {
		t.Fatalf("expected nil for empty dehydrate, got %v", data)
	}
}

func TestRepositoryHydrate(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Hydrate(map[string]any{"x": 10, "y": 20})

	if r.Get("x") != 10 {
		t.Fatalf("expected x = 10, got %v", r.Get("x"))
	}

	if r.Get("y") != 20 {
		t.Fatalf("expected y = 20, got %v", r.Get("y"))
	}
}

func TestRepositoryDehydratingCallback(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("a", 1)

	called := false
	r.Dehydrating(func(_ *logctx.Repository) {
		called = true
	})

	r.Dehydrate()

	if !called {
		t.Fatal("expected dehydrating callback to be called")
	}
}

func TestRepositoryHydratedCallback(t *testing.T) {
	t.Parallel()

	r := logctx.New()

	called := false
	r.Hydrated(func(_ *logctx.Repository) {
		called = true
	})

	r.Hydrate(map[string]any{"a": 1})

	if !called {
		t.Fatal("expected hydrated callback to be called")
	}
}

func TestRepositoryRoundTripDehydrateHydrate(t *testing.T) {
	t.Parallel()

	r1 := logctx.New()
	r1.Add("user_id", 42).Add("req_id", "abc")

	data := r1.Dehydrate()

	r2 := logctx.New()
	r2.Hydrate(data)

	if r2.Get("user_id") != 42 {
		t.Fatalf("expected user_id = 42, got %v", r2.Get("user_id"))
	}

	if r2.Get("req_id") != "abc" {
		t.Fatalf("expected req_id = abc, got %v", r2.Get("req_id"))
	}
}

func TestRepositoryHydrateNil(t *testing.T) {
	t.Parallel()

	r := logctx.New()

	called := false
	r.Hydrated(func(_ *logctx.Repository) {
		called = true
	})

	r.Hydrate(nil)

	if !called {
		t.Fatal("expected hydrated callback to fire even with nil data")
	}
}

// Flush clears hidden and stacks

func TestRepositoryFlushClearsHiddenAndStacks(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Add("a", 1).AddHidden("b", 2).Push("items", "x")
	r.Flush()

	if !r.IsEmpty() {
		t.Fatal("expected empty after flush")
	}

	if len(r.AllHidden()) != 0 {
		t.Fatal("expected hidden to be cleared")
	}
}

func TestRepositoryIsEmpty(t *testing.T) {
	t.Parallel()

	r := logctx.New()

	if !r.IsEmpty() {
		t.Fatal("expected new repository to be empty")
	}

	r.Add("key", "value")

	if r.IsEmpty() {
		t.Fatal("expected repository with data to not be empty")
	}
}

func TestRepositoryMissingHidden(t *testing.T) {
	t.Parallel()

	r := logctx.New()

	if !r.MissingHidden("key") {
		t.Fatal("expected MissingHidden to return true for unset key")
	}

	r.AddHidden("key", "value")

	if r.MissingHidden("key") {
		t.Fatal("expected MissingHidden to return false for set key")
	}
}

func TestRepositoryDecrementFromZero(t *testing.T) {
	t.Parallel()

	r := logctx.New()
	r.Decrement("counter")

	if r.Get("counter") != -1 {
		t.Fatalf("expected -1, got %v", r.Get("counter"))
	}
}
