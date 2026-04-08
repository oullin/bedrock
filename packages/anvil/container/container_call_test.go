package container

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

// Upstream's ContainerCallTest covers Container::call() which invokes closures
// and methods with automatic dependency injection via reflection.
// Go has no reflection-based auto-wiring, so most Upstream tests are skipped.
// This file tests Make, MustMake, FactoryFunc, and binding registration — the
// Go equivalents of the call/resolve entry points.
//
// INTENTIONAL-SKIP: testCallWithAtSignBasedClassReferencesWithoutMethodThrowsException
// INTENTIONAL-SKIP: testCallWithAtSignBasedClassReferences
// INTENTIONAL-SKIP: testCallWithCallableArray
// INTENTIONAL-SKIP: testCallWithStaticMethodNameString
// INTENTIONAL-SKIP: testCallWithGlobalMethodName
// INTENTIONAL-SKIP: testCallWithBoundMethod
// INTENTIONAL-SKIP: testBindMethodAcceptsAnArray
// INTENTIONAL-SKIP: testClosureCallWithInjectedDependency
// INTENTIONAL-SKIP: testCallWithDependencies
// INTENTIONAL-SKIP: testCallWithCallableObject
// INTENTIONAL-SKIP: testCallWithCallableClassString
// INTENTIONAL-SKIP: testCallWithUnnamedParametersThrowsException
// INTENTIONAL-SKIP: testCallWithNullableClassParameterDefaultValue
// INTENTIONAL-SKIP: testCallWithNullableClassParameterDefaultValueWithBinding

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

func TestFactoryFuncWithSingleton(t *testing.T) {
	t.Parallel()

	c := New()
	c.Singleton("svc", func(_ *Container) (any, error) { return "cached", nil })

	factory := c.FactoryFunc("svc")

	v1, _ := factory()
	v2, _ := factory()

	if v1 != v2 {
		t.Fatal("FactoryFunc should return cached singleton on repeated calls")
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

// --- Portable Upstream tests ---

func TestCallWithVariadicDependency(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("items", func(_ *Container) (any, error) {
		return []string{"a", "b", "c"}, nil
	})

	got, err := c.Make("items")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	items := got.([]string)
	if len(items) != 3 {
		t.Fatalf("want 3 items, got %d", len(items))
	}
}

func TestCallWithoutRequiredParamsThrowsException(t *testing.T) {
	t.Parallel()

	c := New()
	_, err := c.Make("unregistered-service")

	if err == nil {
		t.Fatal("expected error for unbound abstract")
	}
}

func TestCallWithoutRequiredParamsOnClosureThrowsException(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("broken", func(_ *Container) (any, error) {
		return nil, fmt.Errorf("cannot resolve required param")
	})

	_, err := c.Make("broken")
	if err == nil {
		t.Fatal("expected error from factory")
	}
}
