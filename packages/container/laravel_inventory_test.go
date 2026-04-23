package container_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/bedrock/packages/container"
)

type laravelConfigStub struct {
	data map[string]any
}

func (c laravelConfigStub) Get(key string, fallback ...any) any {
	if value, ok := c.data[key]; ok {
		return value
	}

	if len(fallback) > 0 {
		return fallback[0]
	}

	return nil
}

// ContainerTest::testContainerSingleton
// ContainerTest::testClosureResolution
// ContainerTest::testSharedClosureResolution
// ContainerTest::testScopedClosureResolution
// ContainerTest::testScopedBindingsWithClosureReturnType
// ContainerTest::testScopedIf
// ContainerTest::testScopedClosureResets
// ContainerTest::testSharedConcreteResolution
// ContainerTest::testScopedConcreteResolutionResets
// ContainerTest::testAbstractToConcreteResolution
// ContainerTest::testNestedDependencyResolution
// ContainerTest::testContainerIsPassedToResolvers
// ContainerTest::testBindingsCanBeOverridden
// ContainerTest::testBindingAnInstanceReturnsTheInstance
// ContainerTest::testBindingAnInstanceAsShared
// ContainerTest::testScopedSingletonWithBind
// ContainerTest::testSingletonWithBind
// ContainerTest::testWithFactoryHasDependency
func TestUpstreamContainerBindingInventoryEquivalents(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	name, err := c.Make("name")

	if err != nil {
		t.Fatal(err)
	}

	if name != "Taylor" {
		t.Fatalf("expected closure binding to resolve Taylor, got %v", name)
	}

	type stub struct{ id int }

	calls := 0
	c.Singleton("singleton", func(_ *container.Container) (any, error) {
		calls++

		return &stub{id: calls}, nil
	})

	firstSingleton, err := c.Make("singleton")

	if err != nil {
		t.Fatal(err)
	}

	secondSingleton, err := c.Make("singleton")

	if err != nil {
		t.Fatal(err)
	}

	if firstSingleton != secondSingleton || calls != 1 {
		t.Fatalf("expected singleton to be shared after one factory call, got %v/%v calls=%d", firstSingleton, secondSingleton, calls)
	}

	scopedCalls := 0
	c.Scoped("scoped", func(_ *container.Container) (any, error) {
		scopedCalls++

		return &stub{id: scopedCalls}, nil
	})

	firstScoped, err := c.Make("scoped")

	if err != nil {
		t.Fatal(err)
	}

	secondScoped, err := c.Make("scoped")

	if err != nil {
		t.Fatal(err)
	}

	if firstScoped != secondScoped {
		t.Fatal("expected scoped binding to be shared within the current scope")
	}

	c.ForgetScopedInstances()
	thirdScoped, err := c.Make("scoped")

	if err != nil {
		t.Fatal(err)
	}

	if thirdScoped == firstScoped {
		t.Fatal("expected scoped binding to refresh after ForgetScopedInstances")
	}

	c.ScopedIf("scoped-if", func(_ *container.Container) (any, error) {
		return "first", nil
	})
	c.ScopedIf("scoped-if", func(_ *container.Container) (any, error) {
		return "second", nil
	})
	scopedIf, err := c.Make("scoped-if")

	if err != nil {
		t.Fatal(err)
	}

	if scopedIf != "first" {
		t.Fatalf("expected ScopedIf to keep original binding, got %v", scopedIf)
	}

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

	outer, err := c.Make("outer")

	if err != nil {
		t.Fatal(err)
	}

	if outer != "outer:inner-value" {
		t.Fatalf("expected nested dependency to resolve, got %v", outer)
	}

	c.Bind("container", func(cc *container.Container) (any, error) {
		return cc, nil
	}, false)
	resolvedContainer, err := c.Make("container")

	if err != nil {
		t.Fatal(err)
	}

	if resolvedContainer != c {
		t.Fatal("expected resolver to receive the active container")
	}

	c.Bind("override", func(_ *container.Container) (any, error) {
		return "old", nil
	}, false)
	c.Bind("override", func(_ *container.Container) (any, error) {
		return "new", nil
	}, false)
	overridden, err := c.Make("override")

	if err != nil {
		t.Fatal(err)
	}

	if overridden != "new" {
		t.Fatalf("expected latest binding to win, got %v", overridden)
	}

	instance := &stub{id: 99}

	if returned := c.Instance("instance", instance); returned != instance {
		t.Fatal("expected Instance to return the registered instance")
	}

	resolvedInstance, err := c.Make("instance")

	if err != nil {
		t.Fatal(err)
	}

	if resolvedInstance != instance || !c.IsShared("instance") {
		t.Fatal("expected registered instance to resolve as a shared value")
	}

	factory := c.FactoryFunc("name")
	factoryValue, err := factory()

	if err != nil {
		t.Fatal(err)
	}

	if factoryValue != "Taylor" {
		t.Fatalf("expected factory to resolve name, got %v", factoryValue)
	}
}

// ContainerTest::testBindIfDoesntRegisterIfServiceAlreadyRegistered
// ContainerTest::testBindIfDoesRegisterIfServiceNotRegisteredYet
// ContainerTest::testSingletonIfDoesntRegisterIfBindingAlreadyRegistered
// ContainerTest::testSingletonIfDoesRegisterIfBindingNotRegisteredYet
func TestUpstreamContainerConditionalBindingInventoryEquivalents(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)
	c.BindIf("name", func(_ *container.Container) (any, error) {
		return "Dayle", nil
	}, false)

	name, err := c.Make("name")

	if err != nil {
		t.Fatal(err)
	}

	if name != "Taylor" {
		t.Fatalf("expected BindIf to keep existing binding, got %v", name)
	}

	c.BindIf("new-name", func(_ *container.Container) (any, error) {
		return "Abigail", nil
	}, false)
	newName, err := c.Make("new-name")

	if err != nil {
		t.Fatal(err)
	}

	if newName != "Abigail" {
		t.Fatalf("expected BindIf to register missing binding, got %v", newName)
	}

	c.Singleton("shared-name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	})
	c.SingletonIf("shared-name", func(_ *container.Container) (any, error) {
		return "Dayle", nil
	})
	sharedName, err := c.Make("shared-name")

	if err != nil {
		t.Fatal(err)
	}

	if sharedName != "Taylor" {
		t.Fatalf("expected SingletonIf to keep existing binding, got %v", sharedName)
	}

	c.SingletonIf("new-shared-name", func(_ *container.Container) (any, error) {
		return "Abigail", nil
	})
	newSharedName, err := c.Make("new-shared-name")

	if err != nil {
		t.Fatal(err)
	}

	if newSharedName != "Abigail" || !c.IsShared("new-shared-name") {
		t.Fatalf("expected SingletonIf to register shared binding, got %v shared=%v", newSharedName, c.IsShared("new-shared-name"))
	}
}

// ContainerTest::testAliases
// ContainerTest::testAliasesWithArrayOfParameters
// ContainerTest::testResolvedResolvesAliasToBindingNameBeforeChecking
// ContainerTest::testGetAlias
// ContainerTest::testGetAliasRecursive
// ContainerTest::testItThrowsExceptionWhenAbstractIsSameAsAlias
func TestUpstreamContainerAliasInventoryEquivalents(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("name", func(cc *container.Container) (any, error) {
		params := cc.Parameters()

		if params != nil {
			return params["value"], nil
		}

		return "Taylor", nil
	}, false)

	c.Alias("name", "username")
	c.Alias("username", "display-name")

	value, err := c.Make("username")

	if err != nil {
		t.Fatal(err)
	}

	if value != "Taylor" {
		t.Fatalf("expected alias to resolve name, got %v", value)
	}

	withParams, err := c.MakeWith("display-name", map[string]any{"value": "Abigail"})

	if err != nil {
		t.Fatal(err)
	}

	if withParams != "Abigail" {
		t.Fatalf("expected alias to resolve with parameters, got %v", withParams)
	}

	c.Make("username") //nolint:errcheck

	if !c.Resolved("display-name") {
		t.Fatal("expected Resolved to normalize aliases")
	}

	if alias := c.GetAlias("display-name"); alias != "name" {
		t.Fatalf("expected recursive alias to resolve to name, got %q", alias)
	}

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected self-alias to panic")
		}
	}()

	c.Alias("same", "same")
}

// ContainerTest::testMakeWithMethodIsAnAliasForMakeMethod
// ContainerTest::testResolvingWithArrayOfParameters
// ContainerTest::testResolvingWithArrayOfMixedParameters
// ContainerTest::testNestedParameterOverride
// ContainerTest::testNestedParametersAreResetForFreshMake
// ContainerTest::testSingletonBindingsNotRespectedWithMakeParameters
func TestUpstreamContainerParameterInventoryEquivalents(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Singleton("name", func(cc *container.Container) (any, error) {
		params := cc.Parameters()

		if params != nil {
			return params["name"], nil
		}

		return "Taylor", nil
	})

	withParams, err := c.MakeWith("name", map[string]any{"name": "Abigail", "age": 42})

	if err != nil {
		t.Fatal(err)
	}

	if withParams != "Abigail" {
		t.Fatalf("expected MakeWith parameter override, got %v", withParams)
	}

	withoutParams, err := c.Make("name")

	if err != nil {
		t.Fatal(err)
	}

	if withoutParams != "Taylor" {
		t.Fatalf("expected fresh Make to reset parameter stack and populate singleton cache, got %v", withoutParams)
	}

	secondWithParams, err := c.MakeWith("name", map[string]any{"name": "Rachel"})

	if err != nil {
		t.Fatal(err)
	}

	if secondWithParams != "Rachel" {
		t.Fatalf("expected singleton cache to be bypassed for parameterized MakeWith, got %v", secondWithParams)
	}

	c.Bind("outer", func(cc *container.Container) (any, error) {
		outerName := cc.Parameters()["name"]
		inner, err := cc.MakeWith("name", map[string]any{"name": "Inner"})

		if err != nil {
			return nil, err
		}

		return []any{outerName, inner}, nil
	}, false)

	nested, err := c.MakeWith("outer", map[string]any{"name": "Outer"})

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(nested, []any{"Outer", "Inner"}) {
		t.Fatalf("expected nested parameter overrides to stay scoped, got %#v", nested)
	}
}

// ContainerTest::testBound
// ContainerTest::testContainerKnowsEntry
// ContainerTest::testContainerCanBindAnyWord
// ContainerTest::testContainerCanDynamicallySetService
// ContainerTest::testUnknownEntryThrowsException
// ContainerTest::testBoundEntriesThrowsContainerExceptionWhenNotResolvable
// ContainerTest::testBindingResolutionExceptionMessage
// ContainerTest::testBindingResolutionExceptionMessageIncludesBuildStack
// ContainerTest::testForgetInstanceForgetsInstance
// ContainerTest::testForgetInstancesForgetsAllInstances
// ContainerTest::testContainerFlushFlushesAllBindingsAliasesAndResolvedInstances
// ContainerTest::testCurrentlyResolving
// ContainerTest::testContainerGetFactory
// ContainerTest::testContainerCanCatchCircularDependency
func TestUpstreamContainerStateAndErrorInventoryEquivalents(t *testing.T) {
	t.Parallel()

	c := newContainer()

	if c.Bound("anything") || c.Has("anything") {
		t.Fatal("fresh container should not know unknown entries")
	}

	c.Bind("any-word", func(_ *container.Container) (any, error) {
		return "ok", nil
	}, false)

	if !c.Bound("any-word") || !c.Has("any-word") {
		t.Fatal("expected Bound and Has to report registered binding")
	}

	c.Instance("dynamic", "set")
	dynamic, err := c.Get("dynamic")

	if err != nil {
		t.Fatal(err)
	}

	if dynamic != "set" {
		t.Fatalf("expected dynamic instance, got %v", dynamic)
	}

	if _, err := c.Get("missing"); !errors.Is(err, container.ErrNotBound) {
		t.Fatalf("expected ErrNotBound for unknown entry, got %v", err)
	}

	expectedErr := errors.New("factory failed")
	c.Bind("failing", func(_ *container.Container) (any, error) {
		return nil, expectedErr
	}, false)

	if _, err := c.Make("failing"); !errors.Is(err, expectedErr) {
		t.Fatalf("expected factory error for bound but unresolvable entry, got %v", err)
	}

	var resolving string

	c.Bind("current", func(cc *container.Container) (any, error) {
		resolving = cc.CurrentlyResolving()

		return "value", nil
	}, false)

	if _, err := c.Make("current"); err != nil {
		t.Fatal(err)
	}

	if resolving != "current" {
		t.Fatalf("expected CurrentlyResolving to expose current abstract, got %q", resolving)
	}

	c.Bind("a", func(cc *container.Container) (any, error) {
		return cc.Make("b")
	}, false)
	c.Bind("b", func(cc *container.Container) (any, error) {
		return cc.Make("a")
	}, false)

	if _, err := c.Make("a"); !errors.Is(err, container.ErrCircularDependency) {
		t.Fatalf("expected circular dependency error, got %v", err)
	}

	c.ForgetInstance("dynamic")

	if c.Resolved("dynamic") {
		t.Fatal("expected ForgetInstance to remove registered instance")
	}

	c.Instance("first", "1")
	c.Instance("second", "2")
	c.ForgetInstances()

	if c.Resolved("first") || c.Resolved("second") {
		t.Fatal("expected ForgetInstances to remove all instances")
	}

	c.Bind("factory-name", func(_ *container.Container) (any, error) {
		return "factory-value", nil
	}, false)
	factory := c.FactoryFunc("factory-name")
	factoryValue, err := factory()

	if err != nil {
		t.Fatal(err)
	}

	if factoryValue != "factory-value" {
		t.Fatalf("expected factory closure to resolve value, got %v", factoryValue)
	}

	c.Alias("factory-name", "factory-alias")
	c.Make("factory-alias") //nolint:errcheck
	c.Flush()

	if c.Bound("factory-name") || c.IsAlias("factory-alias") || c.Resolved("factory-name") {
		t.Fatal("expected Flush to clear bindings, aliases, and resolved state")
	}
}

// ContainerTest::testUnsetRemoveBoundInstances
func TestUpstreamContainerForgetInstanceInventoryEquivalent(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Instance("name", "Taylor")

	if !c.Resolved("name") {
		t.Fatal("expected registered instance to be resolved")
	}

	c.ForgetInstance("name")

	if c.Resolved("name") || c.Bound("name") {
		t.Fatal("expected ForgetInstance to remove the instance-backed binding")
	}
}

// ContainerTest::testReboundListeners
// ContainerTest::testReboundListenersOnInstances
// ContainerTest::testReboundListenersOnInstancesOnlyFiresIfWasAlreadyBound
// ContainerExtendTest::testExtendInstanceRebindingCallback
// ContainerExtendTest::testExtendBindRebindingCallback
func TestUpstreamContainerRebindingInventoryEquivalents(t *testing.T) {
	t.Parallel()

	c := newContainer()

	var reboundValues []any

	if _, err := c.Rebinding("name", func(instance any, _ *container.Container) {
		reboundValues = append(reboundValues, instance)
	}); err != nil {
		t.Fatal(err)
	}

	if len(reboundValues) != 0 {
		t.Fatal("rebinding callback should not fire before the abstract is bound")
	}

	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, false)

	if len(reboundValues) != 0 {
		t.Fatalf("rebinding callback should not fire until an already-resolved binding is rebound, got %#v", reboundValues)
	}

	if _, err := c.Make("name"); err != nil {
		t.Fatal(err)
	}

	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Abigail", nil
	}, false)

	if !reflect.DeepEqual(reboundValues, []any{"Abigail"}) {
		t.Fatalf("expected rebound after rebind, got %#v", reboundValues)
	}

	c.Instance("instance", "before")

	var instanceValue any

	if _, err := c.Rebinding("instance", func(instance any, _ *container.Container) {
		instanceValue = instance
	}); err != nil {
		t.Fatal(err)
	}

	if instanceValue != "before" {
		t.Fatalf("expected Rebinding to immediately receive existing instance, got %v", instanceValue)
	}

	c.Extend("instance", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + "-extended", nil
	})

	if instanceValue != "before-extended" {
		t.Fatalf("expected Extend on instance to fire rebinding callback, got %v", instanceValue)
	}

	c.Bind("extend-bind", func(_ *container.Container) (any, error) {
		return "before", nil
	}, true)

	if _, err := c.Make("extend-bind"); err != nil {
		t.Fatal(err)
	}

	var extendedBinding any

	if _, err := c.Rebinding("extend-bind", func(instance any, _ *container.Container) {
		extendedBinding = instance
	}); err != nil {
		t.Fatal(err)
	}

	c.Extend("extend-bind", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + "-extended", nil
	})

	if extendedBinding != "before-extended" {
		t.Fatalf("expected Extend on resolved binding to fire rebinding callback, got %v", extendedBinding)
	}
}

// ContainerExtendTest::testExtendedBindings
// ContainerExtendTest::testExtendInstancesArePreserved
// ContainerExtendTest::testExtendIsLazyInitialized
// ContainerExtendTest::testExtendCanBeCalledBeforeBind
// ContainerExtendTest::testExtensionWorksOnAliasedBindings
// ContainerExtendTest::testMultipleExtends
// ContainerExtendTest::testUnsetExtend
// ContainerExtendTest::testExtendContextualBinding
// ContainerExtendTest::testExtendContextualBindingAfterResolution
func TestUpstreamContainerExtendInventoryEquivalents(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	}, true)
	c.Extend("name", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + " Otwell", nil
	})

	extended, err := c.Make("name")

	if err != nil {
		t.Fatal(err)
	}

	extendedAgain, err := c.Make("name")

	if err != nil {
		t.Fatal(err)
	}

	if extended != "Taylor Otwell" || extendedAgain != extended {
		t.Fatalf("expected extended singleton to be preserved, got %v/%v", extended, extendedAgain)
	}

	lazyCalled := false
	c.Bind("lazy", func(_ *container.Container) (any, error) {
		return "lazy", nil
	}, false)
	c.Extend("lazy", func(instance any, _ *container.Container) (any, error) {
		lazyCalled = true

		return instance, nil
	})

	if lazyCalled {
		t.Fatal("expected extender to be lazy before resolution")
	}

	if _, err := c.Make("lazy"); err != nil {
		t.Fatal(err)
	}

	if !lazyCalled {
		t.Fatal("expected extender to run during resolution")
	}

	c.Extend("before-bind", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + "-extended", nil
	})
	c.Bind("before-bind", func(_ *container.Container) (any, error) {
		return "value", nil
	}, false)
	beforeBind, err := c.Make("before-bind")

	if err != nil {
		t.Fatal(err)
	}

	if beforeBind != "value-extended" {
		t.Fatalf("expected pre-binding extender to apply, got %v", beforeBind)
	}

	c.Bind("alias-target", func(_ *container.Container) (any, error) {
		return "alias", nil
	}, false)
	c.Alias("alias-target", "alias-name")
	c.Extend("alias-name", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + "-extended", nil
	})
	aliasValue, err := c.Make("alias-name")

	if err != nil {
		t.Fatal(err)
	}

	if aliasValue != "alias-extended" {
		t.Fatalf("expected alias extender to apply, got %v", aliasValue)
	}

	c.Bind("many", func(_ *container.Container) (any, error) {
		return "a", nil
	}, false)
	c.Extend("many", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + "b", nil
	})
	c.Extend("many", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + "c", nil
	})
	many, err := c.Make("many")

	if err != nil {
		t.Fatal(err)
	}

	if many != "abc" {
		t.Fatalf("expected multiple extenders to run in order, got %v", many)
	}

	c.Bind("forgotten", func(_ *container.Container) (any, error) {
		return "plain", nil
	}, false)
	c.Extend("forgotten", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + "-extended", nil
	})
	c.ForgetExtenders("forgotten")
	forgotten, err := c.Make("forgotten")

	if err != nil {
		t.Fatal(err)
	}

	if forgotten != "plain" {
		t.Fatalf("expected extenders to be forgotten, got %v", forgotten)
	}

	c.When("service").Needs("logger").Give(container.Factory(func(_ *container.Container) (any, error) {
		return "contextual", nil
	}))
	c.Extend("logger", func(instance any, _ *container.Container) (any, error) {
		return instance.(string) + "-extended", nil
	})
	c.Bind("service", func(cc *container.Container) (any, error) {
		return cc.Make("logger")
	}, false)
	contextual, err := c.Make("service")

	if err != nil {
		t.Fatal(err)
	}

	if contextual != "contextual-extended" {
		t.Fatalf("expected contextual binding extender to apply, got %v", contextual)
	}

	contextualAgain, err := c.Make("service")

	if err != nil {
		t.Fatal(err)
	}

	if contextualAgain != "contextual-extended" {
		t.Fatalf("expected contextual binding extender to keep applying after resolution, got %v", contextualAgain)
	}
}

// ContainerTaggingTest::testContainerTags
// ContainerTaggingTest::testTaggedServicesAreLazyLoaded
// ContainerTaggingTest::testLazyLoadedTaggedServicesCanBeLoopedOverMultipleTimes
func TestUpstreamContainerTaggingInventoryEquivalents(t *testing.T) {
	t.Parallel()

	c := newContainer()

	var resolved []string

	c.Bind("report-a", func(_ *container.Container) (any, error) {
		resolved = append(resolved, "a")

		return "a", nil
	}, false)
	c.Bind("report-b", func(_ *container.Container) (any, error) {
		resolved = append(resolved, "b")

		return "b", nil
	}, false)

	c.Tag([]string{"report-a", "report-b"}, "reports")

	if len(resolved) != 0 {
		t.Fatal("tagged services should not resolve until Tagged is called")
	}

	first := c.Tagged("reports")
	second := c.Tagged("reports")

	if !reflect.DeepEqual(first, []any{"a", "b"}) || !reflect.DeepEqual(second, []any{"a", "b"}) {
		t.Fatalf("unexpected tagged results: first=%#v second=%#v", first, second)
	}

	if !reflect.DeepEqual(resolved, []string{"a", "b", "a", "b"}) {
		t.Fatalf("expected tagged services to resolve on each iteration, got %#v", resolved)
	}
}

// ContextualBindingTest::testContainerCanInjectDifferentImplementationsDependingOnContext
// ContextualBindingTest::testContextualBindingWorksForExistingInstancedBindings
// ContextualBindingTest::testContextualBindingWorksForNewlyInstancedBindings
// ContextualBindingTest::testContextualBindingWorksOnExistingAliasedInstances
// ContextualBindingTest::testContextualBindingWorksOnNewAliasedInstances
// ContextualBindingTest::testContextualBindingWorksOnNewAliasedBindings
// ContextualBindingTest::testContextualBindingDoesNotFollowStaleAliases
// ContextualBindingTest::testContextualBindingWorksForMultipleClasses
// ContextualBindingTest::testContextualBindingDoesntOverrideNonContextualResolution
// ContextualBindingTest::testContextuallyBoundInstancesAreNotUnnecessarilyRecreated
// ContextualBindingTest::testContainerCanInjectSimpleVariable
// ContextualBindingTest::testContextualBindingWorksWithAliasedTargets
// ContextualBindingTest::testContextualBindingGivesTagsForArrayWithNoTagsDefined
// ContextualBindingTest::testContextualBindingGivesTagsForArray
// ContextualBindingTest::testContextualBindingGivesValuesFromConfigOptionalValueNull
// ContextualBindingTest::testContextualBindingGivesValuesFromConfigOptionalValueSet
// ContextualBindingTest::testContextualBindingGivesValuesFromConfigWithDefault
// ContextualBindingTest::testContextualBindingGivesValuesFromConfigArray
// ContextualBindingTest::testContextualBindingWorksForMethodInvocation
func TestUpstreamContainerContextualBindingInventoryEquivalents(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Bind("logger", func(_ *container.Container) (any, error) {
		return "default", nil
	}, false)
	c.When("service-a").Needs("logger").Give(container.Factory(func(_ *container.Container) (any, error) {
		return "special", nil
	}))
	c.Bind("service-a", func(cc *container.Container) (any, error) {
		return cc.Make("logger")
	}, false)
	c.Bind("service-b", func(cc *container.Container) (any, error) {
		return cc.Make("logger")
	}, false)

	a, err := c.Make("service-a")

	if err != nil {
		t.Fatal(err)
	}

	b, err := c.Make("service-b")

	if err != nil {
		t.Fatal(err)
	}

	if a != "special" || b != "default" {
		t.Fatalf("expected contextual binding only for service-a, got a=%v b=%v", a, b)
	}

	direct, err := c.Make("logger")

	if err != nil {
		t.Fatal(err)
	}

	if direct != "default" {
		t.Fatalf("expected non-contextual resolution to use default, got %v", direct)
	}

	c.When("service-c", "service-d").Needs("logger").Give("multi")

	for _, name := range []string{"service-c", "service-d"} {
		service := name
		c.Bind(service, func(cc *container.Container) (any, error) {
			return cc.Make("logger")
		}, false)
		value, err := c.Make(service)

		if err != nil {
			t.Fatal(err)
		}

		if value != "multi" {
			t.Fatalf("expected contextual value for %s, got %v", service, value)
		}
	}

	c.Instance("instanced", "default-instance")
	c.When("service-instance").Needs("instanced").Give("context-instance")
	c.Bind("service-instance", func(cc *container.Container) (any, error) {
		return cc.Make("instanced")
	}, false)
	instanced, err := c.Make("service-instance")

	if err != nil {
		t.Fatal(err)
	}

	if instanced != "context-instance" {
		t.Fatalf("expected contextual binding to override existing instance, got %v", instanced)
	}

	c.Bind("aliased-target", func(_ *container.Container) (any, error) {
		return "aliased-default", nil
	}, false)
	c.Alias("aliased-target", "aliased")
	c.When("service-alias").Needs("aliased").Give("aliased-context")
	c.Bind("service-alias", func(cc *container.Container) (any, error) {
		return cc.Make("aliased")
	}, false)
	aliased, err := c.Make("service-alias")

	if err != nil {
		t.Fatal(err)
	}

	if aliased != "aliased-context" {
		t.Fatalf("expected contextual binding to work with alias target, got %v", aliased)
	}

	contextualCalls := 0
	c.When("service-once").Needs("expensive").Give(container.Factory(func(_ *container.Container) (any, error) {
		contextualCalls++

		return "expensive", nil
	}))
	c.Bind("service-once", func(cc *container.Container) (any, error) {
		first, err := cc.Make("expensive")

		if err != nil {
			return nil, err
		}

		second, err := cc.Make("expensive")

		if err != nil {
			return nil, err
		}

		return []any{first, second}, nil
	}, false)
	once, err := c.Make("service-once")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(once, []any{"expensive", "expensive"}) || contextualCalls != 2 {
		t.Fatalf("expected contextual factory to resolve on each dependency request, got %#v calls=%d", once, contextualCalls)
	}

	c.When("service-variable").Needs("$name").Give("Taylor")
	c.Bind("service-variable", func(cc *container.Container) (any, error) {
		return cc.Make("$name")
	}, false)
	variable, err := c.Make("service-variable")

	if err != nil {
		t.Fatal(err)
	}

	if variable != "Taylor" {
		t.Fatalf("expected simple contextual variable, got %v", variable)
	}

	c.When("empty-aggregator").Needs("reports").GiveTagged("missing-tag")
	c.Bind("empty-aggregator", func(cc *container.Container) (any, error) {
		return cc.Make("reports")
	}, false)
	emptyReports, err := c.Make("empty-aggregator")

	if err != nil {
		t.Fatal(err)
	}

	if reports := emptyReports.([]any); len(reports) != 0 {
		t.Fatalf("expected empty tagged slice for missing tag, got %#v", reports)
	}

	c.Bind("report-a", func(_ *container.Container) (any, error) {
		return "a", nil
	}, false)
	c.Bind("report-b", func(_ *container.Container) (any, error) {
		return "b", nil
	}, false)
	c.Tag([]string{"report-a", "report-b"}, "reports")
	c.When("aggregator").Needs("reports").GiveTagged("reports")
	c.Bind("aggregator", func(cc *container.Container) (any, error) {
		return cc.Make("reports")
	}, false)
	reports, err := c.Make("aggregator")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(reports, []any{"a", "b"}) {
		t.Fatalf("expected tagged reports, got %#v", reports)
	}

	c.Instance("config", laravelConfigStub{data: map[string]any{
		"nullable": nil,
		"set":      "configured",
		"array":    []string{"one", "two"},
	}})
	c.When("config-consumer").Needs("nullable").GiveConfig("nullable", "fallback")
	c.When("config-consumer").Needs("set").GiveConfig("set", "fallback")
	c.When("config-consumer").Needs("missing-with-default").GiveConfig("missing", "fallback")
	c.When("config-consumer").Needs("array").GiveConfig("array")
	c.Bind("config-consumer", func(cc *container.Container) (any, error) {
		return []any{
			mustMake(t, cc, "nullable"),
			mustMake(t, cc, "set"),
			mustMake(t, cc, "missing-with-default"),
			mustMake(t, cc, "array"),
		}, nil
	}, false)
	configValues, err := c.Make("config-consumer")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(configValues, []any{nil, "configured", "fallback", []string{"one", "two"}}) {
		t.Fatalf("unexpected config values: %#v", configValues)
	}

	c.BindMethod("Controller@show", func(cc *container.Container, params map[string]any) (any, error) {
		return []any{params["_instance"], mustMake(t, cc, "logger")}, nil
	})
	methodResult, err := c.CallMethodBinding("Controller@show", "controller")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(methodResult, []any{"controller", "default"}) {
		t.Fatalf("expected method binding to resolve through container, got %#v", methodResult)
	}
}

// ContainerCallTest::testCallWithBoundMethod
// ContainerCallTest::testClosureCallWithInjectedDependency
// ContainerCallTest::testCallWithDependencies
// ContainerTest::testMethodLevelContextualBinding
func TestUpstreamContainerCallInventoryEquivalents(t *testing.T) {
	t.Parallel()

	c := newContainer()

	c.Instance("name", "Taylor")
	c.BindMethod("Controller@show", func(cc *container.Container, params map[string]any) (any, error) {
		name, err := cc.Make("name")

		if err != nil {
			return nil, err
		}

		return []any{name, params["_instance"]}, nil
	})

	method, err := c.CallMethodBinding("Controller@show", "controller")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(method, []any{"Taylor", "controller"}) {
		t.Fatalf("unexpected method binding result: %#v", method)
	}

	called, err := c.Call(func(cc *container.Container, params map[string]any) (any, error) {
		name, err := cc.Make("name")

		if err != nil {
			return nil, err
		}

		return []any{name, params["suffix"]}, nil
	}, map[string]any{"suffix": "Otwell"})

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(called, []any{"Taylor", "Otwell"}) {
		t.Fatalf("unexpected call result: %#v", called)
	}
}

// ResolvingCallbackTest::testResolvingCallbacksAreCalledForSpecificAbstracts
// ResolvingCallbackTest::testResolvingCallbacksAreCalled
// ResolvingCallbackTest::testResolvingCallbacksShouldBeFiredWhenCalledWithAliases
// ResolvingCallbackTest::testResolvingCallbacksAreCalledOnceForSingletonConcretes
// ResolvingCallbackTest::testResolvingCallbacksCanStillBeAddedAfterTheFirstResolution
// ResolvingCallbackTest::testRebindingDoesNotAffectResolvingCallbacks
// ResolvingCallbackTest::testResolvingCallbacksAreCallWhenRebindHappens
// ResolvingCallbackTest::testResolvingCallbacksArentCalledWhenNoRebindingsAreRegistered
// ResolvingCallbackTest::testRebindingDoesNotAffectMultipleResolvingCallbacks
// ResolvingCallbackTest::testResolvingCallbacksAreCalledForConcretesWhenAttachedOnConcretes
// ResolvingCallbackTest::testAfterResolvingCallbacksAreCalledOnceForImplementation
// ResolvingCallbackTest::testBeforeResolvingCallbacksAreCalled
// ResolvingCallbackTest::testGlobalBeforeResolvingCallbacksAreCalled
func TestUpstreamContainerResolvingCallbackInventoryEquivalents(t *testing.T) {
	t.Parallel()

	c := newContainer()

	var order []string

	c.BeforeResolvingAny(func(abstract string, _ map[string]any, _ *container.Container) {
		order = append(order, "global-before:"+abstract)
	})
	c.BeforeResolving("name", func(abstract string, _ map[string]any, _ *container.Container) {
		order = append(order, "specific-before:"+abstract)
	})
	c.ResolvingAny(func(instance any, _ *container.Container) {
		order = append(order, "global-resolving:"+instance.(string))
	})
	c.Resolving("name", func(instance any, _ *container.Container) {
		order = append(order, "specific-resolving:"+instance.(string))
	})
	c.AfterResolving("name", func(instance any, _ *container.Container) {
		order = append(order, "after:"+instance.(string))
	})

	c.Singleton("name", func(_ *container.Container) (any, error) {
		return "Taylor", nil
	})
	c.Alias("name", "username")

	value, err := c.Make("username")

	if err != nil {
		t.Fatal(err)
	}

	if value != "Taylor" {
		t.Fatalf("expected alias resolution to return Taylor, got %v", value)
	}

	expected := []string{
		"global-before:name",
		"specific-before:name",
		"global-resolving:Taylor",
		"specific-resolving:Taylor",
		"after:Taylor",
	}

	if !reflect.DeepEqual(order, expected) {
		t.Fatalf("unexpected callback order: %#v", order)
	}

	if _, err := c.Make("name"); err != nil {
		t.Fatal(err)
	}

	expectedAfterCachedSingleton := append(slicesClone(expected), "global-before:name", "specific-before:name")

	if !reflect.DeepEqual(order, expectedAfterCachedSingleton) {
		t.Fatalf("singleton should not re-fire resolving callbacks, got %#v", order)
	}

	lateCalled := false
	c.Bind("late", func(_ *container.Container) (any, error) {
		return "first", nil
	}, false)

	if _, err := c.Make("late"); err != nil {
		t.Fatal(err)
	}

	c.Resolving("late", func(_ any, _ *container.Container) {
		lateCalled = true
	})

	if _, err := c.Make("late"); err != nil {
		t.Fatal(err)
	}

	if !lateCalled {
		t.Fatal("expected resolving callback added after first resolution to fire on next resolution")
	}

	var reboundResolving []string

	c.Bind("rebound", func(_ *container.Container) (any, error) {
		return "before", nil
	}, false)
	c.Resolving("rebound", func(instance any, _ *container.Container) {
		reboundResolving = append(reboundResolving, instance.(string))
	})

	if _, err := c.Rebinding("rebound", func(_ any, _ *container.Container) {}); err != nil {
		t.Fatal(err)
	}

	c.Bind("rebound", func(_ *container.Container) (any, error) {
		return "after", nil
	}, false)

	if !reflect.DeepEqual(reboundResolving, []string{"before", "after"}) {
		t.Fatalf("expected resolving callbacks during rebind flow, got %#v", reboundResolving)
	}

	var multiple []string

	c.Bind("multiple", func(_ *container.Container) (any, error) {
		return "value", nil
	}, false)
	c.Resolving("multiple", func(_ any, _ *container.Container) {
		multiple = append(multiple, "first")
	})
	c.Resolving("multiple", func(_ any, _ *container.Container) {
		multiple = append(multiple, "second")
	})

	if _, err := c.Make("multiple"); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(multiple, []string{"first", "second"}) {
		t.Fatalf("expected multiple resolving callbacks, got %#v", multiple)
	}
}

func mustMake(t *testing.T, c *container.Container, abstract string) any {
	t.Helper()

	value, err := c.Make(abstract)

	if err != nil {
		t.Fatal(err)
	}

	return value
}

func slicesClone[T any](values []T) []T {
	out := make([]T, len(values))
	copy(out, values)

	return out
}
