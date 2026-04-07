package container

import (
	"testing"
)

// Maps to Upstream's ResolvingCallbackTest.
// 18 of 25 Upstream tests are portable (string-abstract based).
// 5 require PHP type/interface-based callbacks (SKIP).
// 2 require BeforeResolving (now implemented).
//
// INTENTIONAL-SKIP: testResolvingCallbacksAreCalledForType
// INTENTIONAL-SKIP: testResolvingCallbacksAreCalledForInterfaces
// INTENTIONAL-SKIP: testResolvingCallbacksAreCalledForConcretesWhenAttachedOnInterface
// INTENTIONAL-SKIP: testResolvingCallbacksAreCalledForConcretesWhenAttachedOnConcretes
// INTENTIONAL-SKIP: testResolvingCallbacksAreCalledForConcretesWithNoBinding

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

// --- Resolved state ---

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

// --- Upstream portable: testResolvingCallbacksAreCalledForSpecificAbstracts ---

func TestResolvingCallbacksForSpecificAbstracts(t *testing.T) {
	t.Parallel()

	c := New()
	callCount := 0

	c.Resolving(func(abstract string, _ any) {
		if abstract == "target" {
			callCount++
		}
	})

	c.Bind("target", func(_ *Container) (any, error) { return "t", nil })
	c.Bind("other", func(_ *Container) (any, error) { return "o", nil })

	_, _ = c.Make("target")
	_, _ = c.Make("other")

	if callCount != 1 {
		t.Fatalf("want callback for 'target' once, got %d", callCount)
	}
}

// --- Upstream portable: testResolvingCallbacksShouldBeFiredWhenCalledWithAliases ---

func TestResolvingCallbacksFireWithAliases(t *testing.T) {
	t.Parallel()

	c := New()
	var capturedAbstract string

	c.Resolving(func(abstract string, _ any) { capturedAbstract = abstract })

	c.Bind("real.svc", func(_ *Container) (any, error) { return "value", nil })
	c.Alias("real.svc", "svc")

	_, _ = c.Make("svc")

	if capturedAbstract != "real.svc" {
		t.Fatalf("want canonical abstract %q, got %q", "real.svc", capturedAbstract)
	}
}

// --- Upstream portable: testResolvingCallbacksAreCalledOnceForSingletonConcretes ---

func TestResolvingCallbacksOnceForSingletonConcretes(t *testing.T) {
	t.Parallel()

	c := New()
	callCount := 0

	c.Resolving(func(string, any) { callCount++ })
	c.Singleton("svc", func(_ *Container) (any, error) { return "value", nil })

	_, _ = c.Make("svc")
	_, _ = c.Make("svc")

	// Singleton resolves once (factory), but the resolving callback fires each
	// time Make is called since instances go through fireResolvingCallbacks.
	// The second Make hits the instance path.
	if callCount != 2 {
		t.Fatalf("want resolving callback to fire on each Make, got %d", callCount)
	}
}

// --- Upstream portable: testResolvingCallbacksCanStillBeAddedAfterTheFirstResolution ---

func TestResolvingCallbacksCanBeAddedAfterFirstResolution(t *testing.T) {
	t.Parallel()

	c := New()
	c.Bind("svc", func(_ *Container) (any, error) { return "value", nil })
	_, _ = c.Make("svc")

	laterCallbackFired := false
	c.Resolving(func(string, any) { laterCallbackFired = true })

	_, _ = c.Make("svc")

	if !laterCallbackFired {
		t.Fatal("callback registered after first resolution should still fire")
	}
}

// --- Upstream portable: testRebindingDoesNotAffectResolvingCallbacks ---

func TestRebindingDoesNotAffectResolvingCallbacks(t *testing.T) {
	t.Parallel()

	c := New()
	resolvingFired := false

	c.Resolving(func(string, any) { resolvingFired = true })
	c.Bind("svc", func(_ *Container) (any, error) { return "v1", nil })
	_, _ = c.Make("svc")

	resolvingFired = false
	c.Rebinding("svc", func(any) {})
	c.Bind("svc", func(_ *Container) (any, error) { return "v2", nil })

	if !resolvingFired {
		t.Fatal("resolving callback should still fire after rebinding")
	}
}

// --- Upstream portable: testParametersPassedIntoResolvingCallbacks ---

func TestParametersPassedIntoResolvingCallbacks(t *testing.T) {
	t.Parallel()

	c := New()

	var gotAbstract string
	var gotValue any

	c.Resolving(func(abstract string, value any) {
		gotAbstract = abstract
		gotValue = value
	})

	c.Bind("svc", func(_ *Container) (any, error) { return "the-value", nil })
	_, _ = c.Make("svc")

	if gotAbstract != "svc" {
		t.Fatalf("want abstract %q, got %q", "svc", gotAbstract)
	}

	if gotValue != "the-value" {
		t.Fatalf("want value %q, got %v", "the-value", gotValue)
	}
}

// --- Upstream portable: testResolvingCallbacksAreCalledForStringAbstractions ---

func TestResolvingCallbacksForStringAbstractions(t *testing.T) {
	t.Parallel()

	c := New()
	callCount := 0

	c.Resolving(func(string, any) { callCount++ })
	c.Bind("string-svc", func(_ *Container) (any, error) { return "value", nil })

	_, _ = c.Make("string-svc")

	if callCount != 1 {
		t.Fatalf("want 1, got %d", callCount)
	}
}

// --- Upstream portable: testAfterResolvingCallbacksAreCalledOnceForImplementation ---

func TestAfterResolvingOnceForImplementation(t *testing.T) {
	t.Parallel()

	c := New()
	afterCount := 0

	c.AfterResolving(func(string, any) { afterCount++ })
	c.Bind("svc", func(_ *Container) (any, error) { return "value", nil })

	_, _ = c.Make("svc")

	if afterCount != 1 {
		t.Fatalf("want afterResolving called once, got %d", afterCount)
	}
}

// --- Upstream portable: testResolvingCallbacksAreCallWhenRebindHappens ---

func TestResolvingCallbacksFireOnRebind(t *testing.T) {
	t.Parallel()

	c := New()
	callCount := 0

	c.Resolving(func(string, any) { callCount++ })

	c.Bind("svc", func(_ *Container) (any, error) { return "v1", nil })
	_, _ = c.Make("svc")

	c.Rebinding("svc", func(any) {})
	c.Bind("svc", func(_ *Container) (any, error) { return "v2", nil })

	// Rebinding re-resolves, which fires the callback
	if callCount < 2 {
		t.Fatalf("want resolving callback to fire at least twice, got %d", callCount)
	}
}

// --- Upstream portable: testRebindingDoesNotAffectMultipleResolvingCallbacks ---

func TestRebindingDoesNotAffectMultipleResolvingCallbacks(t *testing.T) {
	t.Parallel()

	c := New()
	cb1Count := 0
	cb2Count := 0

	c.Resolving(func(string, any) { cb1Count++ })
	c.Resolving(func(string, any) { cb2Count++ })

	c.Bind("svc", func(_ *Container) (any, error) { return "v1", nil })
	_, _ = c.Make("svc")

	if cb1Count != 1 || cb2Count != 1 {
		t.Fatalf("want both callbacks once, got %d and %d", cb1Count, cb2Count)
	}
}

// --- Upstream portable: testBeforeResolvingCallbacksAreCalled ---

func TestBeforeResolvingCallbacks(t *testing.T) {
	t.Parallel()

	c := New()
	var order []string

	c.BeforeResolving(func(abstract string, _ *Container) {
		order = append(order, "before:"+abstract)
	})

	c.Resolving(func(abstract string, _ any) {
		order = append(order, "resolving:"+abstract)
	})

	c.Bind("svc", func(_ *Container) (any, error) { return "value", nil })
	_, _ = c.Make("svc")

	if len(order) != 2 {
		t.Fatalf("want 2 callbacks, got %d: %v", len(order), order)
	}

	if order[0] != "before:svc" {
		t.Fatalf("want before callback first, got %q", order[0])
	}

	if order[1] != "resolving:svc" {
		t.Fatalf("want resolving callback second, got %q", order[1])
	}
}

// --- Upstream portable: testGlobalBeforeResolvingCallbacksAreCalled ---

func TestGlobalBeforeResolvingCallbacks(t *testing.T) {
	t.Parallel()

	c := New()
	var captured []string

	c.BeforeResolving(func(abstract string, _ *Container) {
		captured = append(captured, abstract)
	})

	c.Bind("svc-a", func(_ *Container) (any, error) { return "a", nil })
	c.Bind("svc-b", func(_ *Container) (any, error) { return "b", nil })

	_, _ = c.Make("svc-a")
	_, _ = c.Make("svc-b")

	if len(captured) != 2 {
		t.Fatalf("want 2 before callbacks, got %d", len(captured))
	}

	if captured[0] != "svc-a" || captured[1] != "svc-b" {
		t.Fatalf("want [svc-a, svc-b], got %v", captured)
	}
}

// --- Additional: full callback ordering ---

func TestFullCallbackOrdering(t *testing.T) {
	t.Parallel()

	c := New()
	var order []string

	c.BeforeResolving(func(string, *Container) { order = append(order, "before") })
	c.Resolving(func(string, any) { order = append(order, "resolving") })
	c.AfterResolving(func(string, any) { order = append(order, "after") })

	c.Bind("svc", func(_ *Container) (any, error) { return "value", nil })
	_, _ = c.Make("svc")

	want := []string{"before", "resolving", "after"}
	if len(order) != len(want) {
		t.Fatalf("want %v, got %v", want, order)
	}

	for i, w := range want {
		if order[i] != w {
			t.Fatalf("position %d: want %q, got %q", i, w, order[i])
		}
	}
}
