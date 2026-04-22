package pennant_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bedrock/packages/featureflags"
)

// ---------------------------------------------------------------------------
// Store / Driver
// ---------------------------------------------------------------------------

func TestManager_StoreReturnsDefault(t *testing.T) {
	t.Parallel()

	m := featureflags.NewManager("array")

	dec, err := m.Store()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dec == nil {
		t.Fatal("expected non-nil Decorator")
	}
}

func TestManager_StoreNamed(t *testing.T) {
	t.Parallel()

	m := featureflags.NewManager("array")

	dec, err := m.Store("array")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dec == nil {
		t.Fatal("expected non-nil Decorator")
	}
}

func TestManager_StoreIsSingleton(t *testing.T) {
	t.Parallel()

	m := featureflags.NewManager("array")

	first, err := m.Store()

	if err != nil {
		t.Fatalf("first Store: %v", err)
	}

	second, err := m.Store()

	if err != nil {
		t.Fatalf("second Store: %v", err)
	}

	if first != second {
		t.Fatal("expected the same Decorator pointer on both calls (singleton)")
	}
}

func TestManager_Driver_IsAliasForStore(t *testing.T) {
	t.Parallel()

	m := featureflags.NewManager("array")

	fromStore, err := m.Store("array")

	if err != nil {
		t.Fatalf("Store: %v", err)
	}

	fromDriver, err := m.Driver("array")

	if err != nil {
		t.Fatalf("Driver: %v", err)
	}

	if fromStore != fromDriver {
		t.Fatal("Driver() must return the same pointer as Store()")
	}
}

// ---------------------------------------------------------------------------
// Extend
// ---------------------------------------------------------------------------

func TestManager_Extend_CustomFactory(t *testing.T) {
	t.Parallel()

	m := featureflags.NewManager("array")

	invoked := false
	m.Extend("custom", func(_ map[string]any) (featureflags.Driver, error) {
		invoked = true

		return featureflags.NewArrayDriver(), nil
	})

	dec, err := m.Store("custom")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dec == nil {
		t.Fatal("expected non-nil Decorator")
	}

	if !invoked {
		t.Fatal("expected the custom factory to have been invoked")
	}
}

func TestManager_Extend_CustomFactory_UsedByStore(t *testing.T) {
	t.Parallel()

	m := featureflags.NewManager("my-driver")

	called := 0
	m.Extend("my-driver", func(_ map[string]any) (featureflags.Driver, error) {
		called++

		return featureflags.NewArrayDriver(), nil
	})

	// First call creates the driver via the factory.
	if _, err := m.Store("my-driver"); err != nil {
		t.Fatalf("first Store: %v", err)
	}

	// Second call returns the cached instance; factory must NOT be called again.
	if _, err := m.Store("my-driver"); err != nil {
		t.Fatalf("second Store: %v", err)
	}

	if called != 1 {
		t.Fatalf("expected factory to be called exactly once, got %d", called)
	}
}

// ---------------------------------------------------------------------------
// SetDefaultDriver / GetDefaultDriver
// ---------------------------------------------------------------------------

func TestManager_SetDefaultDriver_ChangesDefault(t *testing.T) {
	t.Parallel()

	m := featureflags.NewManager("array")

	m.Extend("other", func(_ map[string]any) (featureflags.Driver, error) {
		return featureflags.NewArrayDriver(), nil
	})

	m.SetDefaultDriver("other")

	dec, err := m.Store()

	if err != nil {
		t.Fatalf("unexpected error after SetDefaultDriver: %v", err)
	}

	if dec == nil {
		t.Fatal("expected non-nil Decorator for new default driver")
	}

	if m.GetDefaultDriver() != "other" {
		t.Fatalf("expected default driver to be %q, got %q", "other", m.GetDefaultDriver())
	}
}

func TestManager_GetDefaultDriver(t *testing.T) {
	t.Parallel()

	m := featureflags.NewManager("array")

	if got := m.GetDefaultDriver(); got != "array" {
		t.Fatalf("expected %q, got %q", "array", got)
	}
}

// ---------------------------------------------------------------------------
// FlushCache
// ---------------------------------------------------------------------------

func TestManager_FlushCache_PropagatesAll(t *testing.T) {
	t.Parallel()

	dispatcher := &testDispatcher{}
	m := featureflags.NewManagerWithDispatcher("array", dispatcher)
	ctx := context.Background()

	dec, err := m.Store()

	if err != nil {
		t.Fatalf("Store: %v", err)
	}

	dec.Define("flag", func(_ context.Context, _ any) (any, error) {
		return true, nil
	})

	// First Get: Decorator cache miss → FeatureResolved dispatched (count = 1).
	if _, err := dec.Get(ctx, "flag", nil); err != nil {
		t.Fatalf("Get (1st): %v", err)
	}

	if dispatcher.count("FeatureResolved") != 1 {
		t.Fatalf("setup: expected 1 FeatureResolved event, got %d", dispatcher.count("FeatureResolved"))
	}

	// FlushCache on the Manager clears the Decorator's in-process cache.
	m.FlushCache()

	// After the flush the Decorator cache is empty, so the next Get is a cache
	// miss and FeatureResolved is dispatched again (value served from driver
	// storage, resolver is not necessarily re-called).
	if _, err := dec.Get(ctx, "flag", nil); err != nil {
		t.Fatalf("Get (2nd): %v", err)
	}

	if dispatcher.count("FeatureResolved") != 2 {
		t.Fatalf("expected 2 FeatureResolved events after FlushCache, got %d", dispatcher.count("FeatureResolved"))
	}
}

// ---------------------------------------------------------------------------
// SerializeScope
// ---------------------------------------------------------------------------

func TestManager_SerializeScope_Delegates(t *testing.T) {
	t.Parallel()

	m := featureflags.NewManager("array")

	got, err := m.SerializeScope("user:42")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "user:42" {
		t.Fatalf("expected %q, got %q", "user:42", got)
	}

	got, err = m.SerializeScope(nil)

	if err != nil {
		t.Fatalf("unexpected error for nil scope: %v", err)
	}

	if got != featureflags.NullScope {
		t.Fatalf("expected NullScope %q, got %q", featureflags.NullScope, got)
	}
}

// ---------------------------------------------------------------------------
// ResolveScopeUsing
// ---------------------------------------------------------------------------

func TestManager_ResolveScopeUsing(t *testing.T) {
	t.Parallel()

	m := featureflags.NewManager("array")

	called := 0
	m.ResolveScopeUsing(func(_ context.Context) (any, error) {
		called++

		return "team:1", nil
	})

	dec, err := m.DefaultDecorator()
	if err != nil {
		t.Fatalf("DefaultDecorator: %v", err)
	}

	dec.Define("flag", func(_ context.Context, scope any) (any, error) {
		return scope, nil
	})

	// FeatureManagerTest::test_the_authenticated_user_is_the_default_scope
	scoped, err := m.For()
	if err != nil {
		t.Fatalf("For: %v", err)
	}

	val, err := scoped.Value(context.Background(), "flag")
	if err != nil {
		t.Fatalf("Value: %v", err)
	}

	if val != "team:1" {
		t.Fatalf("expected default scope to resolve to team:1, got %v", val)
	}

	if called != 1 {
		t.Fatalf("expected scope resolver to run once, got %d", called)
	}
}

// ---------------------------------------------------------------------------
// Unknown driver
// ---------------------------------------------------------------------------

func TestManager_StoreUnknownDriver_ReturnsError(t *testing.T) {
	t.Parallel()

	m := featureflags.NewManager("array")

	_, err := m.Store("does-not-exist")

	if err == nil {
		t.Fatal("expected an error for an unknown driver, got nil")
	}

	if !errors.Is(err, featureflags.ErrDriverNotFound) {
		t.Fatalf("expected ErrDriverNotFound, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// NewManagerWithDispatcher
// ---------------------------------------------------------------------------

func TestManager_WithDispatcher_DispatchesEvents(t *testing.T) {
	t.Parallel()

	dispatcher := &testDispatcher{}
	m := featureflags.NewManagerWithDispatcher("array", dispatcher)
	ctx := context.Background()

	dec, err := m.Store()

	if err != nil {
		t.Fatalf("Store: %v", err)
	}

	dec.Define("flag", func(_ context.Context, _ any) (any, error) {
		return true, nil
	})

	if _, err := dec.Get(ctx, "flag", nil); err != nil {
		t.Fatalf("Get: %v", err)
	}

	if dispatcher.count("FeatureResolved") < 1 {
		t.Fatalf("expected at least 1 FeatureResolved event, got %d", dispatcher.count("FeatureResolved"))
	}
}

// ---------------------------------------------------------------------------
// DefaultDecorator
// ---------------------------------------------------------------------------

func TestManager_DefaultDecorator_ReturnsSameAsStore(t *testing.T) {
	t.Parallel()

	m := featureflags.NewManager("array")

	fromStore, err := m.Store()

	if err != nil {
		t.Fatalf("Store: %v", err)
	}

	fromDefault, err := m.DefaultDecorator()

	if err != nil {
		t.Fatalf("DefaultDecorator: %v", err)
	}

	if fromStore != fromDefault {
		t.Fatal("DefaultDecorator() must return the same pointer as Store()")
	}
}
