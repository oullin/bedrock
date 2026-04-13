package container_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/container"
)

func newContainer() *container.Container {
	return container.New()
}

// ---------- Binding & Resolution ----------

func TestClosureResolution(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	v, err := c.Make("name")

	if err != nil {
		t.Fatal(err)
	}

	if v != "Taylor" {
		t.Fatalf("expected Taylor, got %v", v)
	}
}

func TestBindReturnsDifferentInstances(t *testing.T) {
	t.Parallel()

	c := newContainer()

	type stub struct{ id int }

	n := 0

	c.Bind("stub", func(_ *container.Container) (any, error) {
		n++

		return &stub{id: n}, nil
	}, false)

	a, _ := c.Make("stub")
	b, _ := c.Make("stub")

	if a == b {
		t.Fatal("expected different instances for non-shared binding")
	}
}

func TestBindIfDoesNotRegisterIfAlreadyBound(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.BindIf("name", func(_ *container.Container) (any, error) {
		return "Dayle", nil
	}, false)

	v, _ := c.Make("name")

	if v != "Taylor" {
		t.Fatalf("expected Taylor, got %v", v)
	}
}

func TestBindIfDoesRegisterIfNotBound(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.BindIf("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	v, _ := c.Make("name")

	if v != "Taylor" {
		t.Fatalf("expected Taylor, got %v", v)
	}
}

func TestSingletonReturnsSameInstance(t *testing.T) {
	t.Parallel()

	c := newContainer()

	type stub struct{ value string }

	c.Singleton("stub", func(_ *container.Container) (any, error) {
		return &stub{value: "hello"}, nil
	})

	a, _ := c.Make("stub")
	b, _ := c.Make("stub")

	if a != b {
		t.Fatal("expected same instance for singleton")
	}
}

func TestSingletonIfDoesNotRegisterIfAlreadyBound(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Singleton("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	})

	c.SingletonIf("name", func(_ *container.Container) (any, error) {
		return "Dayle", nil
	})

	v, _ := c.Make("name")

	if v != "Taylor" {
		t.Fatalf("expected Taylor, got %v", v)
	}
}

func TestSingletonIfDoesRegisterIfNotBound(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.SingletonIf("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	})

	v, _ := c.Make("name")

	if v != "Taylor" {
		t.Fatalf("expected Taylor, got %v", v)
	}
}

func TestScopedReturnsSameInstanceThenResets(t *testing.T) {
	t.Parallel()

	c := newContainer()

	type stub struct{ id int }

	n := 0

	c.Scoped("stub", func(_ *container.Container) (any, error) {
		n++

		return &stub{id: n}, nil
	})

	a, _ := c.Make("stub")
	b, _ := c.Make("stub")

	if a != b {
		t.Fatal("expected same instance for scoped binding")
	}

	c.ForgetScopedInstances()

	d, _ := c.Make("stub")

	if a == d {
		t.Fatal("expected new instance after ForgetScopedInstances")
	}
}

func TestScopedIfDoesNotRegisterIfAlreadyBound(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Scoped("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	})

	c.ScopedIf("name", func(_ *container.Container) (any, error) {
		return "Dayle", nil
	})

	v, _ := c.Make("name")

	if v != "Taylor" {
		t.Fatalf("expected Taylor, got %v", v)
	}
}

func TestInstanceReturnsRegisteredValue(t *testing.T) {
	t.Parallel()

	c := newContainer()

	type stub struct{ Value string }

	obj := &stub{Value: "hello"}
	c.Instance("stub", obj)

	v, _ := c.Make("stub")

	if v != obj {
		t.Fatal("expected exact same instance")
	}
}

func TestInstanceReturnsTheInstance(t *testing.T) {
	t.Parallel()

	c := newContainer()

	type stub struct{}

	obj := &stub{}
	ret := c.Instance("stub", obj)

	if ret != obj {
		t.Fatal("Instance should return the instance")
	}
}

func TestBindingAnInstanceAsShared(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Instance("stub", "hello")

	if !c.IsShared("stub") {
		t.Fatal("instance should be shared")
	}
}

func TestContainerIsPassedToResolvers(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("something", func(cc *container.Container) (any, error) {
		return cc, nil
	}, false)

	v, _ := c.Make("something")

	if v != c {
		t.Fatal("expected container to be passed to factory")
	}
}

func TestNestedDependencyResolution(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("inner", func(_ *container.Container) (any, error) {
		return "inner-value", nil
	}, false)

	c.Bind("outer", func(cc *container.Container) (any, error) {
		inner, err := cc.Make("inner")

		if err != nil {
			return nil, err
		}

		return "outer:" + inner.(string), nil
	}, false)

	v, err := c.Make("outer")

	if err != nil {
		t.Fatal(err)
	}

	if v != "outer:inner-value" {
		t.Fatalf("expected outer:inner-value, got %v", v)
	}
}

func TestMakeWithParameters(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("greeting", func(cc *container.Container) (any, error) {
		params := cc.Parameters()

		name, _ := params["name"]

		return "Hello, " + name.(string), nil
	}, false)

	v, err := c.MakeWith("greeting", map[string]any{"name": "Taylor"})

	if err != nil {
		t.Fatal(err)
	}

	if v != "Hello, Taylor" {
		t.Fatalf("expected Hello, Taylor, got %v", v)
	}
}

func TestSingletonNotCachedWithMakeParameters(t *testing.T) {
	t.Parallel()

	c := newContainer()

	count := 0

	c.Singleton("counter", func(_ *container.Container) (any, error) {
		count++

		return count, nil
	})

	a, _ := c.MakeWith("counter", map[string]any{"force": true})
	b, _ := c.MakeWith("counter", map[string]any{"force": true})

	if a == b {
		t.Fatal("singleton with parameters should not be cached")
	}
}

func TestNestedParametersAreReset(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("inner", func(cc *container.Container) (any, error) {
		params := cc.Parameters()

		if params != nil {
			return "has-params", nil
		}

		return "no-params", nil
	}, false)

	c.Bind("outer", func(cc *container.Container) (any, error) {
		return cc.Make("inner")
	}, false)

	v, _ := c.MakeWith("outer", map[string]any{"key": "value"})

	if v != "no-params" {
		t.Fatalf("expected no-params for nested make, got %v", v)
	}
}

func TestBindingsCanBeOverridden(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Dayle", nil
	}, false)

	v, _ := c.Make("name")

	if v != "Dayle" {
		t.Fatalf("expected Dayle, got %v", v)
	}
}

func TestContainerCanBindAnyWord(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("Taylor", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	if v, _ := c.Make("Taylor"); v != "Taylor" {
		t.Fatalf("expected Taylor, got %v", v)
	}
}

func TestBuildExecutesFactory(t *testing.T) {
	t.Parallel()

	c := newContainer()

	v, err := c.Build(func(_ *container.Container) (any, error) {
		return "built", nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if v != "built" {
		t.Fatalf("expected built, got %v", v)
	}
}

func TestGetReturnsErrNotBoundForUnknown(t *testing.T) {
	t.Parallel()

	c := newContainer()

	_, err := c.Get("unknown")

	if !errors.Is(err, container.ErrNotBound) {
		t.Fatalf("expected ErrNotBound, got %v", err)
	}
}

func TestGetResolvesBinding(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	v, err := c.Get("name")

	if err != nil {
		t.Fatal(err)
	}

	if v != "Taylor" {
		t.Fatalf("expected Taylor, got %v", v)
	}
}

func TestFactoryFuncReturnsNewInstanceEachTime(t *testing.T) {
	t.Parallel()

	c := newContainer()

	type stub struct{ id int }

	n := 0

	c.Bind("stub", func(_ *container.Container) (any, error) {
		n++

		return &stub{id: n}, nil
	}, false)

	factory := c.FactoryFunc("stub")

	a, _ := factory()
	b, _ := factory()

	if a == b {
		t.Fatal("factory should return new instance each time")
	}
}

func TestMakeWithReturnsErrNotBound(t *testing.T) {
	t.Parallel()

	c := newContainer()

	_, err := c.Make("nonexistent")

	if !errors.Is(err, container.ErrNotBound) {
		t.Fatalf("expected ErrNotBound, got %v", err)
	}
}

func TestCircularDependencyDetection(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("a", func(cc *container.Container) (any, error) {
		return cc.Make("b")
	}, false)

	c.Bind("b", func(cc *container.Container) (any, error) {
		return cc.Make("c")
	}, false)

	c.Bind("c", func(cc *container.Container) (any, error) {
		return cc.Make("a")
	}, false)

	_, err := c.Make("a")

	if !errors.Is(err, container.ErrCircularDependency) {
		t.Fatalf("expected ErrCircularDependency, got %v", err)
	}
}

// ---------- Aliases ----------

func TestAliasResolvesToAbstract(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.Alias("name", "shortName")

	v, _ := c.Make("shortName")

	if v != "Taylor" {
		t.Fatalf("expected Taylor, got %v", v)
	}
}

func TestAliasChainResolution(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.Alias("name", "alias1")
	c.Alias("alias1", "alias2")

	v, _ := c.Make("alias2")

	if v != "Taylor" {
		t.Fatalf("expected Taylor, got %v", v)
	}
}

func TestGetAliasResolvesChain(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Alias("name", "alias1")
	c.Alias("alias1", "alias2")

	if c.GetAlias("alias2") != "name" {
		t.Fatalf("expected name, got %s", c.GetAlias("alias2"))
	}
}

func TestGetAliasReturnsAbstractIfNotAlias(t *testing.T) {
	t.Parallel()

	c := newContainer()

	if c.GetAlias("notanalias") != "notanalias" {
		t.Fatal("expected same name back")
	}
}

func TestIsAlias(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Alias("name", "shortName")

	if !c.IsAlias("shortName") {
		t.Fatal("expected shortName to be an alias")
	}

	if c.IsAlias("name") {
		t.Fatal("name is not an alias")
	}
}

func TestSelfAliasPanics(t *testing.T) {
	t.Parallel()

	c := newContainer()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for self-alias")
		}
	}()

	c.Alias("name", "name")
}

func TestBoundWithAlias(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.Alias("name", "shortName")

	if !c.Bound("shortName") {
		t.Fatal("alias should report as bound")
	}
}

func TestResolvedResolvesAlias(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.Alias("name", "shortName")
	c.Make("name") //nolint:errcheck

	if !c.Resolved("shortName") {
		t.Fatal("alias should show as resolved")
	}
}

// ---------- Queries ----------

func TestBound(t *testing.T) {
	t.Parallel()

	c := newContainer()

	if c.Bound("something") {
		t.Fatal("should not be bound")
	}

	c.Bind("something", func(_ *container.Container) (any, error) {
		return "value", nil
	}, false)

	if !c.Bound("something") {
		t.Fatal("should be bound")
	}
}

func TestBoundWithInstance(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Instance("name", "Taylor")

	if !c.Bound("name") {
		t.Fatal("instance should be bound")
	}
}

func TestHasDelegatesToBound(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	if !c.Has("name") {
		t.Fatal("Has should delegate to Bound")
	}

	if c.Has("missing") {
		t.Fatal("Has should return false for unbound")
	}
}

func TestResolved(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	if c.Resolved("name") {
		t.Fatal("should not be resolved before Make")
	}

	c.Make("name") //nolint:errcheck

	if !c.Resolved("name") {
		t.Fatal("should be resolved after Make")
	}
}

func TestResolvedWithInstance(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Instance("name", "Taylor")

	if !c.Resolved("name") {
		t.Fatal("instance should report as resolved")
	}
}

func TestIsShared(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("transient", func(_ *container.Container) (any, error) {
		return "value", nil
	}, false)

	c.Singleton("shared", func(_ *container.Container) (any, error) {
		return "value", nil
	})

	if c.IsShared("transient") {
		t.Fatal("transient should not be shared")
	}

	if !c.IsShared("shared") {
		t.Fatal("singleton should be shared")
	}
}

func TestIsSharedWithScoped(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Scoped("scoped", func(_ *container.Container) (any, error) {
		return "value", nil
	})

	if !c.IsShared("scoped") {
		t.Fatal("scoped should be shared")
	}
}

func TestIsSharedWithInstance(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Instance("name", "Taylor")

	if !c.IsShared("name") {
		t.Fatal("instance should be shared")
	}
}

func TestCurrentlyResolving(t *testing.T) {
	t.Parallel()

	c := newContainer()

	var resolving string

	c.Bind("outer", func(cc *container.Container) (any, error) {
		resolving = cc.CurrentlyResolving()

		return "outer", nil
	}, false)

	c.Make("outer") //nolint:errcheck

	if resolving != "outer" {
		t.Fatalf("expected outer, got %s", resolving)
	}
}

func TestCurrentlyResolvingEmpty(t *testing.T) {
	t.Parallel()

	c := newContainer()

	if c.CurrentlyResolving() != "" {
		t.Fatal("expected empty string when not resolving")
	}
}

// ---------- Instance Management ----------

func TestForgetInstanceForgetsInstance(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Instance("name", "Taylor")

	c.ForgetInstance("name")

	if c.Bound("name") {
		t.Fatal("instance should be forgotten (not in bindings)")
	}
}

func TestForgetInstancesForgetsAll(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Instance("a", "1")
	c.Instance("b", "2")

	c.ForgetInstances()

	if c.Bound("a") || c.Bound("b") {
		t.Fatal("all instances should be forgotten")
	}
}

func TestForgetScopedInstancesOnlyClearsScoped(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Singleton("shared", func(_ *container.Container) (any, error) {
		return "shared", nil
	})

	c.Scoped("scoped", func(_ *container.Container) (any, error) {
		return "scoped", nil
	})

	c.Make("shared") //nolint:errcheck
	c.Make("scoped") //nolint:errcheck

	c.ForgetScopedInstances()

	// Shared should still return cached instance.
	a, _ := c.Make("shared")

	if a != "shared" {
		t.Fatal("shared instance should survive ForgetScopedInstances")
	}
}

func TestFlushResetsEverything(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.Alias("name", "shortName")
	c.Instance("inst", "value")
	c.Make("name") //nolint:errcheck

	c.Flush()

	if c.Bound("name") || c.Bound("shortName") || c.Bound("inst") {
		t.Fatal("flush should clear all state")
	}

	if c.Resolved("name") {
		t.Fatal("flush should clear resolved")
	}
}

func TestGetBindings(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("a", func(_ *container.Container) (any, error) { return "a", nil }, false)
	c.Singleton("b", func(_ *container.Container) (any, error) { return "b", nil })

	bindings := c.GetBindings()

	if len(bindings) != 2 {
		t.Fatalf("expected 2 bindings, got %d", len(bindings))
	}
}

func TestUnsetRemovesBoundInstances(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Instance("stub", "value")

	if !c.Bound("stub") {
		t.Fatal("expected bound")
	}

	c.ForgetInstance("stub")

	if c.Bound("stub") {
		t.Fatal("expected unbound after forget")
	}
}

// ---------- Static Instance ----------

func TestGetSetInstance(t *testing.T) {
	t.Parallel()

	c := container.New()
	container.SetInstance(c)

	if container.GetInstance() != c {
		t.Fatal("expected same container")
	}

	container.SetInstance(nil)

	got := container.GetInstance()

	if got == nil {
		t.Fatal("GetInstance should create a new container if nil")
	}

	if got == c {
		t.Fatal("should be a new container")
	}

	// Cleanup.
	container.SetInstance(nil)
}

// ---------- Concurrent Access ----------

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()

	c := newContainer()

	done := make(chan bool, 40)

	for i := 0; i < 10; i++ {
		go func() {
			c.Bind("name", func(_ *container.Container) (any, error) {
				return "Taylor", nil
			}, false)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			c.Make("name") //nolint:errcheck
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			c.Bound("name")
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			c.Has("name")
			done <- true
		}()
	}

	for i := 0; i < 40; i++ {
		<-done
	}
}

func TestMakeWithIsAliasForMake(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	v, err := c.MakeWith("name", nil)

	if err != nil {
		t.Fatal(err)
	}

	if v != "Taylor" {
		t.Fatalf("expected Taylor, got %v", v)
	}
}

// ---------- Rebinding ----------

func TestReboundListeners(t *testing.T) {
	t.Parallel()

	c := newContainer()

	var rebound bool

	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.Make("name") //nolint:errcheck

	c.Rebinding("name", func(_ any, _ *container.Container) { //nolint:errcheck
		rebound = true
	})

	// Rebinding triggers the callback because "name" is already bound.
	if !rebound {
		t.Fatal("expected rebound callback to fire")
	}
}

func TestReboundListenersOnInstances(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Instance("name", "Taylor")

	var called bool

	c.Rebinding("name", func(_ any, _ *container.Container) { //nolint:errcheck
		called = true
	})

	if !called {
		t.Fatal("expected callback on already-bound instance")
	}
}

func TestReboundListenersNotCalledWhenNotBound(t *testing.T) {
	t.Parallel()

	c := newContainer()

	var called bool

	c.Rebinding("name", func(_ any, _ *container.Container) { //nolint:errcheck
		called = true
	})

	if called {
		t.Fatal("should not call callback when not bound")
	}
}

func TestReboundListenersFiredOnRebind(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	c.Make("name") //nolint:errcheck

	var newVal any

	c.Rebinding("name", func(instance any, _ *container.Container) { //nolint:errcheck
		newVal = instance
	})

	// Rebind.
	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Dayle", nil
	}, false)

	if newVal != "Dayle" {
		t.Fatalf("expected Dayle, got %v", newVal)
	}
}

func TestRefreshUpdatesTarget(t *testing.T) {
	t.Parallel()

	c := newContainer()
	c.Instance("name", "Taylor")

	var updated any

	c.Refresh("name", func(v any) {
		updated = v
	})

	if updated != "Taylor" {
		t.Fatalf("expected Taylor, got %v", updated)
	}

	c.Instance("name", "Dayle")

	if updated != "Dayle" {
		t.Fatalf("expected Dayle after refresh, got %v", updated)
	}
}
