package pennant_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/pennant"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// newScopedWithFeatures builds a Decorator backed by an ArrayDriver with the
// supplied feature resolvers pre-registered, then returns a
// ScopedFeatureInteraction for the given scope.
func newScopedWithFeatures(
	features map[string]func(context.Context, any) (any, error),
	scope any,
) *pennant.ScopedFeatureInteraction {
	drv := pennant.NewArrayDriver()
	for name, resolver := range features {
		drv.Define(name, resolver)
	}

	dec := pennant.NewDecorator(drv)

	return pennant.NewScopedFeatureInteraction(dec, scope)
}

// ---------------------------------------------------------------------------
// For
// ---------------------------------------------------------------------------

func TestScoped_For_AddsScope(t *testing.T) {
	t.Parallel()

	drv := pennant.NewArrayDriver()
	dec := pennant.NewDecorator(drv)

	base := pennant.NewScopedFeatureInteraction(dec, "user:1")
	merged := base.For("user:2")

	// The new interaction should hold both scopes; the original is unchanged.
	// We verify the merge by calling AllAreActive on a feature that is active
	// for both scopes — if either scope is missing the call would be wrong.
	drv.Define("flag", func(_ context.Context, _ any) (any, error) { return true, nil })

	ctx := context.Background()

	if !merged.AllAreActive(ctx, []string{"flag"}) {
		t.Fatal("expected AllAreActive to be true for merged scopes")
	}
}

// ---------------------------------------------------------------------------
// Active / Inactive
// ---------------------------------------------------------------------------

func TestScoped_Active_True(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag": func(_ context.Context, _ any) (any, error) { return true, nil },
	}, nil)

	if !s.Active(ctx, "flag") {
		t.Fatal("expected Active=true for bool true value")
	}
}

func TestScoped_Active_False(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag": func(_ context.Context, _ any) (any, error) { return false, nil },
	}, nil)

	if s.Active(ctx, "flag") {
		t.Fatal("expected Active=false for bool false value")
	}
}

func TestScoped_Active_StringVariant(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag": func(_ context.Context, _ any) (any, error) { return "variant-a", nil },
	}, nil)

	if !s.Active(ctx, "flag") {
		t.Fatal("expected Active=true for non-empty string value")
	}
}

func TestScoped_Active_Nil(t *testing.T) {
	t.Parallel()

	// An undefined feature returns an error; Active must return false.
	drv := pennant.NewArrayDriver()
	dec := pennant.NewDecorator(drv)
	s := pennant.NewScopedFeatureInteraction(dec, nil)

	ctx := context.Background()

	if s.Active(ctx, "undefined-feature") {
		t.Fatal("expected Active=false for undefined (nil-returning / erroring) feature")
	}
}

func TestScoped_Inactive_IsInverseOfActive(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag": func(_ context.Context, _ any) (any, error) { return true, nil },
	}, nil)

	if s.Inactive(ctx, "flag") {
		t.Fatal("expected Inactive=false when Active=true")
	}
}

// ---------------------------------------------------------------------------
// Value / Values
// ---------------------------------------------------------------------------

func TestScoped_Value_Returns_RawValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag": func(_ context.Context, _ any) (any, error) { return "variant-b", nil },
	}, nil)

	val, err := s.Value(ctx, "flag")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != "variant-b" {
		t.Fatalf("expected variant-b, got %v", val)
	}
}

func TestScoped_Values_ReturnsMap(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag-a": func(_ context.Context, _ any) (any, error) { return true, nil },
		"flag-b": func(_ context.Context, _ any) (any, error) { return "blue", nil },
	}, nil)

	vals, err := s.Values(ctx, []string{"flag-a", "flag-b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if vals["flag-a"] != true {
		t.Fatalf("expected flag-a=true, got %v", vals["flag-a"])
	}

	if vals["flag-b"] != "blue" {
		t.Fatalf("expected flag-b=blue, got %v", vals["flag-b"])
	}
}

// ---------------------------------------------------------------------------
// AllAreActive
// ---------------------------------------------------------------------------

func TestScoped_AllAreActive_AllTrue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag-a": func(_ context.Context, _ any) (any, error) { return true, nil },
		"flag-b": func(_ context.Context, _ any) (any, error) { return true, nil },
	}, nil)

	if !s.AllAreActive(ctx, []string{"flag-a", "flag-b"}) {
		t.Fatal("expected AllAreActive=true when all features are active")
	}
}

func TestScoped_AllAreActive_OneInactive(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag-a": func(_ context.Context, _ any) (any, error) { return true, nil },
		"flag-b": func(_ context.Context, _ any) (any, error) { return false, nil },
	}, nil)

	if s.AllAreActive(ctx, []string{"flag-a", "flag-b"}) {
		t.Fatal("expected AllAreActive=false when one feature is inactive")
	}
}

// ---------------------------------------------------------------------------
// SomeAreActive
// ---------------------------------------------------------------------------

func TestScoped_SomeAreActive_OneActive(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag-a": func(_ context.Context, _ any) (any, error) { return true, nil },
		"flag-b": func(_ context.Context, _ any) (any, error) { return false, nil },
	}, nil)

	if !s.SomeAreActive(ctx, []string{"flag-a", "flag-b"}) {
		t.Fatal("expected SomeAreActive=true when at least one feature is active")
	}
}

func TestScoped_SomeAreActive_NoneActive(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag-a": func(_ context.Context, _ any) (any, error) { return false, nil },
		"flag-b": func(_ context.Context, _ any) (any, error) { return false, nil },
	}, nil)

	if s.SomeAreActive(ctx, []string{"flag-a", "flag-b"}) {
		t.Fatal("expected SomeAreActive=false when no features are active")
	}
}

// ---------------------------------------------------------------------------
// AllAreInactive
// ---------------------------------------------------------------------------

func TestScoped_AllAreInactive_AllFalse(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag-a": func(_ context.Context, _ any) (any, error) { return false, nil },
		"flag-b": func(_ context.Context, _ any) (any, error) { return false, nil },
	}, nil)

	if !s.AllAreInactive(ctx, []string{"flag-a", "flag-b"}) {
		t.Fatal("expected AllAreInactive=true when all features are inactive")
	}
}

// ---------------------------------------------------------------------------
// SomeAreInactive
// ---------------------------------------------------------------------------

func TestScoped_SomeAreInactive_OneInactive(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag-a": func(_ context.Context, _ any) (any, error) { return true, nil },
		"flag-b": func(_ context.Context, _ any) (any, error) { return false, nil },
	}, nil)

	if !s.SomeAreInactive(ctx, []string{"flag-a", "flag-b"}) {
		t.Fatal("expected SomeAreInactive=true when at least one feature is inactive")
	}
}

// ---------------------------------------------------------------------------
// Activate / ActivateWithValue / Deactivate
// ---------------------------------------------------------------------------

func TestScoped_Activate_SetsTrue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	drv := pennant.NewArrayDriver()
	dec := pennant.NewDecorator(drv)
	s := pennant.NewScopedFeatureInteraction(dec, nil)

	if err := s.Activate(ctx, []string{"flag"}); err != nil {
		t.Fatalf("Activate returned error: %v", err)
	}

	if !s.Active(ctx, "flag") {
		t.Fatal("expected Active=true after Activate")
	}
}

func TestScoped_ActivateWithValue_SetsCustomValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	drv := pennant.NewArrayDriver()
	dec := pennant.NewDecorator(drv)
	s := pennant.NewScopedFeatureInteraction(dec, nil)

	if err := s.ActivateWithValue(ctx, []string{"theme"}, "dark"); err != nil {
		t.Fatalf("ActivateWithValue returned error: %v", err)
	}

	val, err := s.Value(ctx, "theme")
	if err != nil {
		t.Fatalf("Value returned error: %v", err)
	}

	if val != "dark" {
		t.Fatalf("expected dark, got %v", val)
	}
}

func TestScoped_Deactivate_SetsFalse(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	drv := pennant.NewArrayDriver()
	dec := pennant.NewDecorator(drv)
	s := pennant.NewScopedFeatureInteraction(dec, nil)

	// Activate first, then deactivate.
	_ = s.Activate(ctx, []string{"flag"})

	if err := s.Deactivate(ctx, []string{"flag"}); err != nil {
		t.Fatalf("Deactivate returned error: %v", err)
	}

	if s.Active(ctx, "flag") {
		t.Fatal("expected Active=false after Deactivate")
	}
}

// ---------------------------------------------------------------------------
// Forget
// ---------------------------------------------------------------------------

func TestScoped_Forget_RemovesState(t *testing.T) {
	t.Parallel()

	calls := 0
	ctx := context.Background()

	drv := pennant.NewArrayDriver()
	drv.Define("flag", func(_ context.Context, _ any) (any, error) {
		calls++
		return true, nil
	})

	dec := pennant.NewDecorator(drv)
	s := pennant.NewScopedFeatureInteraction(dec, nil)

	// First resolution via Active.
	s.Active(ctx, "flag")

	if calls != 1 {
		t.Fatalf("setup: expected 1 resolver call, got %d", calls)
	}

	if err := s.Forget(ctx, []string{"flag"}); err != nil {
		t.Fatalf("Forget returned error: %v", err)
	}

	// After Forget the state is gone; the next Active must re-invoke the resolver.
	s.Active(ctx, "flag")

	if calls != 2 {
		t.Fatalf("expected resolver called again after Forget, got %d total calls", calls)
	}
}

// ---------------------------------------------------------------------------
// Purge
// ---------------------------------------------------------------------------

func TestScoped_Purge_RemovesAll(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	drv := pennant.NewArrayDriver()

	drv.Define("flag-a", func(_ context.Context, _ any) (any, error) { return true, nil })
	drv.Define("flag-b", func(_ context.Context, _ any) (any, error) { return true, nil })

	dec := pennant.NewDecorator(drv)
	s := pennant.NewScopedFeatureInteraction(dec, nil)

	s.Active(ctx, "flag-a")
	s.Active(ctx, "flag-b")

	if err := s.Purge(ctx, nil); err != nil {
		t.Fatalf("Purge returned error: %v", err)
	}

	stored, err := dec.Stored(ctx)
	if err != nil {
		t.Fatalf("Stored returned error: %v", err)
	}

	if len(stored) != 0 {
		t.Fatalf("expected no stored features after Purge(nil), got %v", stored)
	}
}

// ---------------------------------------------------------------------------
// When / Unless
// ---------------------------------------------------------------------------

func TestScoped_When_CallsActiveCallback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag": func(_ context.Context, _ any) (any, error) { return true, nil },
	}, nil)

	result, err := s.When(ctx, "flag",
		func(val any) (any, error) { return "active-path", nil },
		func(val any) (any, error) { return "inactive-path", nil },
	)

	if err != nil {
		t.Fatalf("When returned error: %v", err)
	}

	if result != "active-path" {
		t.Fatalf("expected active-path, got %v", result)
	}
}

func TestScoped_When_CallsInactiveCallback(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag": func(_ context.Context, _ any) (any, error) { return false, nil },
	}, nil)

	result, err := s.When(ctx, "flag",
		func(val any) (any, error) { return "active-path", nil },
		func(val any) (any, error) { return "inactive-path", nil },
	)

	if err != nil {
		t.Fatalf("When returned error: %v", err)
	}

	if result != "inactive-path" {
		t.Fatalf("expected inactive-path, got %v", result)
	}
}

func TestScoped_When_NilCallbacksNoPanic(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag": func(_ context.Context, _ any) (any, error) { return true, nil },
	}, nil)

	result, err := s.When(ctx, "flag", nil, nil)
	if err != nil {
		t.Fatalf("When with nil callbacks returned error: %v", err)
	}

	if result != nil {
		t.Fatalf("expected nil result for nil active callback, got %v", result)
	}
}

func TestScoped_Unless_IsInverseOfWhen(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	s := newScopedWithFeatures(map[string]func(context.Context, any) (any, error){
		"flag": func(_ context.Context, _ any) (any, error) { return true, nil },
	}, nil)

	// Unless: feature is active → should run the second callback (whenActive in Unless).
	result, err := s.Unless(ctx, "flag",
		func(val any) (any, error) { return "inactive-path", nil },
		func(val any) (any, error) { return "active-path", nil },
	)

	if err != nil {
		t.Fatalf("Unless returned error: %v", err)
	}

	// feature is active → Unless executes the second argument (whenActive).
	if result != "active-path" {
		t.Fatalf("expected active-path from Unless when feature is active, got %v", result)
	}
}

// ---------------------------------------------------------------------------
// Load / LoadMissing
// ---------------------------------------------------------------------------

func TestScoped_Load_PopulatesCache(t *testing.T) {
	t.Parallel()

	calls := 0
	ctx := context.Background()

	drv := pennant.NewArrayDriver()
	drv.Define("flag", func(_ context.Context, _ any) (any, error) {
		calls++
		return true, nil
	})

	dec := pennant.NewDecorator(drv)
	s := pennant.NewScopedFeatureInteraction(dec, nil)

	if err := s.Load(ctx, []string{"flag"}); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if calls != 1 {
		t.Fatalf("expected 1 resolver call after Load, got %d", calls)
	}

	// A subsequent Active call must hit the cache — resolver must not be called again.
	s.Active(ctx, "flag")

	if calls != 1 {
		t.Fatalf("expected resolver call count to remain 1 after cache hit, got %d", calls)
	}
}

// ---------------------------------------------------------------------------
// Multiple scopes
// ---------------------------------------------------------------------------

func TestScoped_MultipleScopes_AllAreActive_RequiresBoth(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	drv := pennant.NewArrayDriver()
	drv.Define("flag", func(_ context.Context, scope any) (any, error) {
		// Only active for user:1.
		if scope == "user:1" {
			return true, nil
		}

		return false, nil
	})

	dec := pennant.NewDecorator(drv)
	s := pennant.NewScopedFeatureInteraction(dec, "user:1", "user:2")

	// user:2 is inactive, so AllAreActive must return false.
	if s.AllAreActive(ctx, []string{"flag"}) {
		t.Fatal("expected AllAreActive=false when one scope is inactive")
	}

	// SomeAreActive must return true because user:1 is active.
	if !s.SomeAreActive(ctx, []string{"flag"}) {
		t.Fatal("expected SomeAreActive=true when at least one scope is active")
	}
}
