package container

import (
	"errors"
	"testing"
)

// Maps to Upstream's UtilTest. Covers generic Make[T]/MustMake[T], alias
// resolution, Bound/Has, ForgetInstance(s), Flush, and IsAlias.
//
// INTENTIONAL-SKIP: testGetParameterClassName (PHP reflection)

// --- Shared test types ---

type testService struct {
	Name string
}

type testLogger interface {
	Log(msg string)
}

type testConsoleLogger struct{}

func (testConsoleLogger) Log(string) {}

// --- Generic Make[T] / MustMake[T] ---

func TestMakeGeneric(t *testing.T) {
	t.Parallel()

	c := New()
	svc := &testService{Name: "auth"}
	c.Instance("svc", svc)

	got, err := Make[*testService](c, "svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Name != "auth" {
		t.Fatalf("want Name=%q, got %q", "auth", got.Name)
	}

	if got != svc {
		t.Fatal("expected same pointer")
	}
}

func TestMakeGenericInterface(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("logger", testConsoleLogger{})

	got, err := Make[testLogger](c, "logger")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestMakeGenericMismatch(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("val", "a string")

	_, err := Make[int](c, "val")
	if err == nil {
		t.Fatal("expected type mismatch error")
	}

	if !errors.Is(err, ErrTypeMismatch) {
		t.Fatalf("expected ErrTypeMismatch, got: %v", err)
	}
}

func TestMustMakeGenericPanic(t *testing.T) {
	t.Parallel()

	c := New()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic from MustMake")
		}
	}()

	MustMake[string](c, "missing")
}

// --- Aliases ---

func TestAliases(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("app.config", func(_ *Container) (any, error) { return "cfg", nil })
	c.Alias("app.config", "config")

	got, err := c.Make("config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "cfg" {
		t.Fatalf("want %q, got %q", "cfg", got)
	}
}

func TestAliasesWithChain(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("real", func(_ *Container) (any, error) { return "value", nil })
	c.Alias("real", "middle")
	c.Alias("middle", "short")

	got, err := c.Make("short")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "value" {
		t.Fatalf("want %q, got %q", "value", got)
	}
}

func TestGetAlias(t *testing.T) {
	t.Parallel()

	c := New()
	c.Alias("target", "alias")

	if got := c.GetAlias("alias"); got != "target" {
		t.Fatalf("want %q, got %q", "target", got)
	}
}

func TestGetAliasRecursive(t *testing.T) {
	t.Parallel()

	c := New()
	c.Alias("real", "middle")
	c.Alias("middle", "short")

	if got := c.GetAlias("short"); got != "real" {
		t.Fatalf("want %q, got %q", "real", got)
	}
}

func TestGetAliasNonAlias(t *testing.T) {
	t.Parallel()

	c := New()

	if got := c.GetAlias("not-an-alias"); got != "not-an-alias" {
		t.Fatalf("want %q returned unchanged, got %q", "not-an-alias", got)
	}
}

func TestItThrowsExceptionWhenAbstractIsSameAsAlias(t *testing.T) {
	t.Parallel()

	c := New()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for self-alias")
		}
	}()

	c.Alias("same", "same")
}

func TestIsAlias(t *testing.T) {
	t.Parallel()

	c := New()
	c.Alias("target", "alias")

	if !c.IsAlias("alias") {
		t.Fatal("expected IsAlias to return true")
	}

	if c.IsAlias("target") {
		t.Fatal("expected IsAlias to return false for non-alias")
	}
}

func TestIsAliasAfterFlush(t *testing.T) {
	t.Parallel()

	c := New()
	c.Alias("target", "alias")
	c.Flush()

	if c.IsAlias("alias") {
		t.Fatal("expected IsAlias to return false after Flush")
	}
}

// --- Bound / Has ---

func TestBound(t *testing.T) {
	t.Parallel()

	c := New()

	c.Bind("bound-bind", func(_ *Container) (any, error) { return nil, nil })
	if !c.Bound("bound-bind") {
		t.Fatal("expected Bound to return true for Bind")
	}

	c.Singleton("bound-singleton", func(_ *Container) (any, error) { return nil, nil })
	if !c.Bound("bound-singleton") {
		t.Fatal("expected Bound to return true for Singleton")
	}

	c.Instance("bound-instance", "val")
	if !c.Bound("bound-instance") {
		t.Fatal("expected Bound to return true for Instance")
	}
}

func TestBoundWithAlias(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("full.name", func(_ *Container) (any, error) { return "ok", nil })
	c.Alias("full.name", "short")

	if !c.Bound("short") {
		t.Fatal("expected Bound to return true via alias")
	}
}

func TestContainerKnowsEntry(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "ok", nil })

	if !c.Has("svc") {
		t.Fatal("expected Has to return true")
	}

	if c.Has("missing") {
		t.Fatal("expected Has to return false for unregistered")
	}
}

// --- Unset / Forget ---

func TestUnsetRemoveBoundInstances(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("key", "val")
	c.ForgetInstance("key")

	if c.Bound("key") {
		t.Fatal("expected instance removed after ForgetInstance")
	}
}

func TestForgetInstanceForgetsInstance(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("key", "val")
	c.ForgetInstance("key")

	if c.Bound("key") {
		t.Fatal("expected instance removed")
	}
}

func TestForgetInstancesForgetsAllInstances(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("a", 1)
	c.Instance("b", 2)
	c.ForgetInstances()

	if c.Bound("a") || c.Bound("b") {
		t.Fatal("expected all instances removed")
	}
}

// --- Flush ---

func TestContainerFlushFlushesAllBindingsAliasesAndResolvedInstances(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("a", func(_ *Container) (any, error) { return 1, nil })
	c.Instance("b", 2)
	c.Alias("a", "aa")
	c.Tag([]string{"a"}, "nums")

	_, _ = c.Make("a")

	c.Flush()

	if c.Bound("a") {
		t.Fatal("expected bindings cleared after Flush")
	}

	if c.Bound("b") {
		t.Fatal("expected instances cleared after Flush")
	}

	if c.IsAlias("aa") {
		t.Fatal("expected aliases cleared after Flush")
	}

	if c.Resolved("a") {
		t.Fatal("expected resolved cleared after Flush")
	}
}

// --- Upstream portable: testUnwrapIfClosure ---

func TestUnwrapIfClosure(t *testing.T) {
	t.Parallel()

	c := New()

	c.When("consumer").Needs("val").Give("plain-string")
	got, _ := c.MakeFor("consumer", "val")

	if got != "plain-string" {
		t.Fatalf("plain value: want %q, got %v", "plain-string", got)
	}

	c.When("consumer").Needs("dynamic").Give(Factory(func(_ *Container) (any, error) {
		return "from-closure", nil
	}))

	got, _ = c.MakeFor("consumer", "dynamic")
	if got != "from-closure" {
		t.Fatalf("closure value: want %q, got %v", "from-closure", got)
	}
}
