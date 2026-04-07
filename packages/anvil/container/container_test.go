package container

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

// --- Bind / Make round-trip ---

func TestClosureResolution(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("greeting", func(_ *Container) (any, error) {
		return "hello", nil
	})

	got, err := c.Make("greeting")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "hello" {
		t.Fatalf("want %q, got %q", "hello", got)
	}
}

func TestBindIfDoesntRegisterIfServiceAlreadyRegistered(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "first", nil })
	c.BindIf("svc", func(_ *Container) (any, error) { return "second", nil })

	got, _ := c.Make("svc")
	if got != "first" {
		t.Fatalf("want %q, got %q", "first", got)
	}
}

func TestBindIfDoesRegisterIfServiceNotRegisteredYet(t *testing.T) {
	t.Parallel()

	c := New()
	c.BindIf("svc", func(_ *Container) (any, error) { return "value", nil })

	got, err := c.Make("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "value" {
		t.Fatalf("want %q, got %q", "value", got)
	}
}

func TestSingletonIfDoesntRegisterIfBindingAlreadyRegistered(t *testing.T) {
	t.Parallel()

	c := New()
	c.Singleton("svc", func(_ *Container) (any, error) { return "first", nil })
	c.SingletonIf("svc", func(_ *Container) (any, error) { return "second", nil })

	got, _ := c.Make("svc")
	if got != "first" {
		t.Fatalf("want %q, got %q", "first", got)
	}
}

func TestSingletonIfDoesRegisterIfBindingNotRegisteredYet(t *testing.T) {
	t.Parallel()

	c := New()
	c.SingletonIf("svc", func(_ *Container) (any, error) { return "value", nil })

	got, err := c.Make("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "value" {
		t.Fatalf("want %q, got %q", "value", got)
	}
}

// --- Singleton ---

func TestSharedClosureResolution(t *testing.T) {
	t.Parallel()

	type svc struct{ id int }

	calls := 0
	c := New()
	c.Singleton("svc", func(_ *Container) (any, error) {
		calls++
		return &svc{id: calls}, nil
	})

	first, err := c.Make("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second, err := c.Make("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != second {
		t.Fatal("singleton should return same pointer")
	}

	if calls != 1 {
		t.Fatalf("factory should be called once, got %d", calls)
	}
}

func TestSharedClosureResolutionConcurrent(t *testing.T) {
	t.Parallel()

	var calls atomic.Int64

	c := New()
	c.Singleton("counter", func(_ *Container) (any, error) {
		calls.Add(1)
		return &struct{ v int }{1}, nil
	})

	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	results := make([]any, goroutines)
	errs := make([]error, goroutines)

	for i := range goroutines {
		go func(idx int) {
			defer wg.Done()
			results[idx], errs[idx] = c.Make("counter")
		}(i)
	}

	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Fatalf("goroutine %d: unexpected error: %v", i, err)
		}
	}

	first := results[0]
	for i := 1; i < goroutines; i++ {
		if results[i] != first {
			t.Fatalf("goroutine %d returned different pointer", i)
		}
	}

	if got := calls.Load(); got != 1 {
		t.Fatalf("factory should be called once, got %d", got)
	}
}

// --- Scoped ---

func TestScopedClosureResolution(t *testing.T) {
	t.Parallel()

	type svc struct{ id int }

	calls := 0
	c := New()
	c.Scoped("svc", func(_ *Container) (any, error) {
		calls++
		return &svc{id: calls}, nil
	})

	first, _ := c.Make("svc")
	second, _ := c.Make("svc")

	if first != second {
		t.Fatal("scoped should return same instance within scope")
	}
}

func TestScopedClosureResets(t *testing.T) {
	t.Parallel()

	type svc struct{ id int }

	calls := 0
	c := New()
	c.Scoped("svc", func(_ *Container) (any, error) {
		calls++
		return &svc{id: calls}, nil
	})

	first, _ := c.Make("svc")
	c.ForgetScopedInstances()
	second, _ := c.Make("svc")

	if first == second {
		t.Fatal("scoped should return new instance after reset")
	}

	if first.(*svc).id == second.(*svc).id {
		t.Fatal("scoped should have called factory again")
	}
}

func TestScopedIf(t *testing.T) {
	t.Parallel()

	c := New()
	c.Scoped("svc", func(_ *Container) (any, error) { return "first", nil })
	c.ScopedIf("svc", func(_ *Container) (any, error) { return "second", nil })

	got, _ := c.Make("svc")
	if got != "first" {
		t.Fatalf("want %q, got %q", "first", got)
	}
}

// --- Instance ---

func TestBindingAnInstanceReturnsTheInstance(t *testing.T) {
	t.Parallel()

	c := New()
	val := &struct{ Name string }{"test"}
	c.Instance("obj", val)

	got, err := c.Make("obj")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != val {
		t.Fatal("instance should return exact same value")
	}
}

func TestBindingAnInstanceAsShared(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("key", "first")
	c.Instance("key", "second")

	got, err := c.Make("key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "second" {
		t.Fatalf("want %q, got %q", "second", got)
	}
}

// --- Container is passed to resolvers ---

func TestContainerIsPassedToResolvers(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("db.dsn", "postgres://localhost/test")
	c.Bind("db", func(c *Container) (any, error) {
		dsn, err := c.Make("db.dsn")
		if err != nil {
			return nil, err
		}
		return "connected:" + dsn.(string), nil
	})

	got, err := c.Make("db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "connected:postgres://localhost/test" {
		t.Fatalf("want connected DSN, got %v", got)
	}
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

// --- Binding override ---

func TestBindingsCanBeOverridden(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("val", func(_ *Container) (any, error) { return "first", nil })
	c.Bind("val", func(_ *Container) (any, error) { return "second", nil })

	got, err := c.Make("val")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "second" {
		t.Fatalf("want %q, got %q", "second", got)
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

func TestContainerCanBindAnyWord(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("Taylor", "Taylor Otwell")

	got, err := c.Make("Taylor")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "Taylor Otwell" {
		t.Fatalf("want %q, got %q", "Taylor Otwell", got)
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

// --- Resolved ---

func TestResolvedResolvesAliasToBindingNameBeforeChecking(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "ok", nil })
	c.Alias("svc", "short")

	_, _ = c.Make("svc")

	if !c.Resolved("short") {
		t.Fatal("expected Resolved to follow alias")
	}
}

func TestResolved(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "ok", nil })

	if c.Resolved("svc") {
		t.Fatal("expected not resolved before Make")
	}

	_, _ = c.Make("svc")

	if !c.Resolved("svc") {
		t.Fatal("expected resolved after Make")
	}
}

// --- Rebound listeners ---

func TestReboundListeners(t *testing.T) {
	t.Parallel()

	c := New()

	var captured any
	c.Bind("svc", func(_ *Container) (any, error) { return "v1", nil })
	_, _ = c.Make("svc")

	c.Rebinding("svc", func(val any) { captured = val })

	c.Bind("svc", func(_ *Container) (any, error) { return "v2", nil })

	if captured != "v2" {
		t.Fatalf("want rebound callback with %q, got %v", "v2", captured)
	}
}

func TestReboundListenersOnInstances(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("svc", "v1")
	_, _ = c.Make("svc")

	var captured any
	c.Rebinding("svc", func(val any) { captured = val })

	c.Instance("svc", "v2")

	if captured != "v2" {
		t.Fatalf("want rebound callback with %q, got %v", "v2", captured)
	}
}

func TestReboundListenersOnInstancesOnlyFiresIfWasAlreadyBound(t *testing.T) {
	t.Parallel()

	c := New()

	var fired bool
	c.Rebinding("svc", func(any) { fired = true })

	c.Instance("svc", "v1")

	if fired {
		t.Fatal("rebound should not fire for first Instance")
	}
}

// --- Error handling ---

func TestUnknownEntryThrowsException(t *testing.T) {
	t.Parallel()

	c := New()
	_, err := c.Make("missing")

	if err == nil {
		t.Fatal("expected error for unbound abstract")
	}

	if !errors.Is(err, ErrNotBound) {
		t.Fatalf("expected ErrNotBound, got: %v", err)
	}
}

func TestBindingResolutionExceptionMessage(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("broken", func(_ *Container) (any, error) {
		return nil, errors.New("factory failed")
	})

	_, err := c.Make("broken")
	if err == nil {
		t.Fatal("expected error from factory")
	}

	if !errors.Is(err, ErrResolve) {
		t.Fatalf("expected ErrResolve, got: %v", err)
	}
}

func TestSingletonFactoryErrorCached(t *testing.T) {
	t.Parallel()

	c := New()
	c.Singleton("broken", func(_ *Container) (any, error) {
		return nil, errors.New("init failed")
	})

	_, err1 := c.Make("broken")
	_, err2 := c.Make("broken")

	if !errors.Is(err1, ErrResolve) {
		t.Fatalf("first call: expected ErrResolve, got %v", err1)
	}

	if !errors.Is(err2, ErrResolve) {
		t.Fatalf("second call: expected same error, got %v", err2)
	}
}

// --- MustMake ---

func TestMustMakePanics(t *testing.T) {
	t.Parallel()

	c := New()

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic from MustMake")
		}
	}()

	c.MustMake("missing")
}

// --- IsAlias ---

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

// --- Tags ---

func TestContainerTags(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("report.csv", func(_ *Container) (any, error) { return "csv", nil })
	c.Tag([]string{"report.csv"}, "reports")

	results, err := c.Tagged("reports")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 || results[0] != "csv" {
		t.Fatalf("want [csv], got %v", results)
	}
}

func TestContainerTagsMultiple(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("a", func(_ *Container) (any, error) { return 1, nil })
	c.Bind("b", func(_ *Container) (any, error) { return 2, nil })
	c.Tag([]string{"a", "b"}, "numbers")

	results, err := c.Tagged("numbers")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("want 2 results, got %d", len(results))
	}
}

func TestContainerTagsEmpty(t *testing.T) {
	t.Parallel()

	c := New()
	results, err := c.Tagged("unknown")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("want empty slice, got %v", results)
	}
}

// --- Resolving callbacks ---

func TestResolvingCallback(t *testing.T) {
	t.Parallel()

	c := New()

	var capturedAbstract string
	var capturedValue any

	c.Resolving(func(abstract string, value any) {
		capturedAbstract = abstract
		capturedValue = value
	})

	c.Bind("svc", func(_ *Container) (any, error) { return 42, nil })
	_, _ = c.Make("svc")

	if capturedAbstract != "svc" {
		t.Fatalf("want abstract %q, got %q", "svc", capturedAbstract)
	}

	if capturedValue != 42 {
		t.Fatalf("want value 42, got %v", capturedValue)
	}
}

func TestAfterResolvingCallback(t *testing.T) {
	t.Parallel()

	c := New()

	var order []string

	c.Resolving(func(string, any) { order = append(order, "resolving") })
	c.AfterResolving(func(string, any) { order = append(order, "after") })

	c.Bind("svc", func(_ *Container) (any, error) { return nil, nil })
	_, _ = c.Make("svc")

	if len(order) != 2 || order[0] != "resolving" || order[1] != "after" {
		t.Fatalf("want [resolving, after], got %v", order)
	}
}

// --- Transient vs Singleton behavior ---

func TestTransientReturnsNewInstances(t *testing.T) {
	t.Parallel()

	type obj struct{ id int }

	calls := 0
	c := New()
	c.Bind("obj", func(_ *Container) (any, error) {
		calls++
		return &obj{id: calls}, nil
	})

	first, _ := c.Make("obj")
	second, _ := c.Make("obj")

	if first.(*obj).id == second.(*obj).id {
		t.Fatal("transient binding should return new instances")
	}
}

func TestSingletonBindingsNotRespectedWithNewBind(t *testing.T) {
	t.Parallel()

	c := New()
	c.Singleton("svc", func(_ *Container) (any, error) { return "singleton-v1", nil })
	_, _ = c.Make("svc")

	c.Bind("svc", func(_ *Container) (any, error) { return "transient", nil })

	got, _ := c.Make("svc")
	if got != "transient" {
		t.Fatalf("want %q, got %q", "transient", got)
	}
}

// --- Instance / Bind precedence ---

func TestInstanceClearsBinding(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "from-factory", nil })
	c.Instance("svc", "from-instance")

	got, _ := c.Make("svc")
	if got != "from-instance" {
		t.Fatalf("want %q, got %q", "from-instance", got)
	}
}

func TestBindClearsInstance(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("svc", "from-instance")
	c.Bind("svc", func(_ *Container) (any, error) { return "from-factory", nil })

	got, _ := c.Make("svc")
	if got != "from-factory" {
		t.Fatalf("want %q, got %q", "from-factory", got)
	}
}

// --- FactoryFunc ---

func TestContainerGetFactory(t *testing.T) {
	t.Parallel()

	c := New()

	calls := 0
	c.Bind("svc", func(_ *Container) (any, error) {
		calls++
		return calls, nil
	})

	factory := c.FactoryFunc("svc")

	v1, err := factory()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v2, err := factory()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v1 == v2 {
		t.Fatal("FactoryFunc should resolve fresh each call for transient bindings")
	}
}

// --- MakeWith (Make is an alias) ---

func TestMakeWithMethodIsAnAliasForMakeMethod(t *testing.T) {
	t.Parallel()

	c := New()
	c.Instance("svc", "value")

	got, err := c.Make("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "value" {
		t.Fatalf("want %q, got %q", "value", got)
	}
}

// --- Scoped with Singleton behavior ---

func TestScopedConcreteResolutionResets(t *testing.T) {
	t.Parallel()

	type svc struct{ id int }

	calls := 0
	c := New()
	c.Scoped("svc", func(_ *Container) (any, error) {
		calls++
		return &svc{id: calls}, nil
	})

	first, _ := c.Make("svc")
	same, _ := c.Make("svc")

	if first != same {
		t.Fatal("scoped should return same instance within scope")
	}

	c.ForgetScopedInstances()

	after, _ := c.Make("svc")
	if first == after {
		t.Fatal("scoped should return new instance after ForgetScopedInstances")
	}
}

// --- ScopedSingleton with Bind override ---

func TestScopedSingletonWithBind(t *testing.T) {
	t.Parallel()

	c := New()
	c.Scoped("svc", func(_ *Container) (any, error) { return "scoped", nil })

	got, _ := c.Make("svc")
	if got != "scoped" {
		t.Fatalf("want %q, got %q", "scoped", got)
	}

	c.Bind("svc", func(_ *Container) (any, error) { return "transient", nil })

	got, _ = c.Make("svc")
	if got != "transient" {
		t.Fatalf("want %q, got %q", "transient", got)
	}
}

func TestSingletonWithBind(t *testing.T) {
	t.Parallel()

	c := New()
	c.Singleton("svc", func(_ *Container) (any, error) { return "singleton", nil })

	got, _ := c.Make("svc")
	if got != "singleton" {
		t.Fatalf("want %q, got %q", "singleton", got)
	}

	c.Bind("svc", func(_ *Container) (any, error) { return "transient", nil })

	got, _ = c.Make("svc")
	if got != "transient" {
		t.Fatalf("want %q, got %q", "transient", got)
	}
}
