package pennant_test

import (
	"context"
	"sync"
	"testing"

	"github.com/bedrock/packages/featureflags"
)

// ---------------------------------------------------------------------------
// Interface assertions
// ---------------------------------------------------------------------------

func TestDecorator_InterfaceAssertions(t *testing.T) {
	t.Parallel()

	var _ featureflags.Driver       = (*featureflags.Decorator)(nil)
	var _ featureflags.CacheFlusher = (*featureflags.Decorator)(nil)
}

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

func TestDecorator_Get_CacheMiss_PopulatesCache(t *testing.T) {
	t.Parallel()

	calls := 0
	drv := featureflags.NewArrayDriver()
	ctx := context.Background()

	drv.Define("flag", func(_ context.Context, _ any) (any, error) {
		calls++
		return true, nil
	})

	dec := featureflags.NewDecorator(drv)

	val, err := dec.Get(ctx, "flag", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != true {
		t.Fatalf("expected true, got %v", val)
	}
	if calls != 1 {
		t.Fatalf("expected 1 driver call, got %d", calls)
	}
}

func TestDecorator_Get_CacheHit(t *testing.T) {
	t.Parallel()

	calls := 0
	drv := featureflags.NewArrayDriver()
	ctx := context.Background()

	drv.Define("flag", func(_ context.Context, _ any) (any, error) {
		calls++
		return "cached-value", nil
	})

	dec := featureflags.NewDecorator(drv)

	// First call populates cache.
	dec.Get(ctx, "flag", nil) //nolint:errcheck

	// Second call must be served from cache — driver should not be called again.
	val, err := dec.Get(ctx, "flag", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "cached-value" {
		t.Fatalf("expected cached-value, got %v", val)
	}
	if calls != 1 {
		t.Fatalf("expected exactly 1 driver call (cache hit on 2nd), got %d", calls)
	}
}

func TestDecorator_Get_DispatchesFeatureResolved(t *testing.T) {
	t.Parallel()

	dispatcher := &testDispatcher{}
	drv := featureflags.NewArrayDriver()
	ctx := context.Background()

	drv.Define("flag", func(_ context.Context, _ any) (any, error) {
		return true, nil
	})

	dec := featureflags.NewDecoratorWithDispatcher(drv, dispatcher)

	dec.Get(ctx, "flag", nil) //nolint:errcheck

	if dispatcher.count("FeatureResolved") != 1 {
		t.Fatalf("expected 1 FeatureResolved event, got %d", dispatcher.count("FeatureResolved"))
	}

	ev, ok := dispatcher.last().(featureflags.FeatureResolved)
	if !ok {
		t.Fatalf("last event is not FeatureResolved: %T", dispatcher.last())
	}
	if ev.Feature != "flag" {
		t.Fatalf("expected feature=flag, got %q", ev.Feature)
	}
	if ev.Value != true {
		t.Fatalf("expected value=true, got %v", ev.Value)
	}
}

func TestDecorator_Get_ErrorNotCached(t *testing.T) {
	t.Parallel()

	drv := featureflags.NewArrayDriver()
	ctx := context.Background()

	// No resolver defined → ErrFeatureNotDefined.
	dec := featureflags.NewDecorator(drv)

	_, err := dec.Get(ctx, "undefined-flag", nil)
	if err == nil {
		t.Fatal("expected error for undefined feature, got nil")
	}

	// Second call must also go to the driver (no error stored in cache).
	_, err2 := dec.Get(ctx, "undefined-flag", nil)
	if err2 == nil {
		t.Fatal("expected error on second call too")
	}
}

// ---------------------------------------------------------------------------
// Set
// ---------------------------------------------------------------------------

func TestDecorator_Set_UpdatesCache_DispatchesEvent(t *testing.T) {
	t.Parallel()

	calls := 0
	dispatcher := &testDispatcher{}
	drv := featureflags.NewArrayDriver()
	ctx := context.Background()

	drv.Define("flag", func(_ context.Context, _ any) (any, error) {
		calls++
		return true, nil
	})

	dec := featureflags.NewDecoratorWithDispatcher(drv, dispatcher)

	if err := dec.Set(ctx, "flag", nil, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Value should come from the cache (set to false), resolver not called.
	val, err := dec.Get(ctx, "flag", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != false {
		t.Fatalf("expected false (set value), got %v", val)
	}
	if calls != 0 {
		t.Fatalf("expected 0 resolver calls after Set, got %d", calls)
	}

	if dispatcher.count("FeatureUpdated") != 1 {
		t.Fatalf("expected 1 FeatureUpdated event, got %d", dispatcher.count("FeatureUpdated"))
	}
}

// ---------------------------------------------------------------------------
// SetForAllScopes
// ---------------------------------------------------------------------------

func TestDecorator_SetForAllScopes_ClearsCacheForFeature_DispatchesEvent(t *testing.T) {
	t.Parallel()

	dispatcher := &testDispatcher{}
	drv := featureflags.NewArrayDriver()
	ctx := context.Background()

	drv.Define("flag", func(_ context.Context, _ any) (any, error) {
		return true, nil
	})

	dec := featureflags.NewDecoratorWithDispatcher(drv, dispatcher)

	// Populate cache for two scopes.
	dec.Get(ctx, "flag", "user:1") //nolint:errcheck
	dec.Get(ctx, "flag", "user:2") //nolint:errcheck

	if err := dec.SetForAllScopes(ctx, "flag", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dispatcher.count("FeatureUpdatedForAllScopes") != 1 {
		t.Fatalf("expected 1 FeatureUpdatedForAllScopes, got %d", dispatcher.count("FeatureUpdatedForAllScopes"))
	}

	// After SetForAllScopes the decorator cache is wiped, so the next Get goes
	// to the driver which now stores false for both scopes.
	val, err := dec.Get(ctx, "flag", "user:1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != false {
		t.Fatalf("expected false after SetForAllScopes, got %v", val)
	}
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestDecorator_Delete_RemovesCacheEntry_DispatchesEvent(t *testing.T) {
	t.Parallel()

	calls := 0
	dispatcher := &testDispatcher{}
	drv := featureflags.NewArrayDriver()
	ctx := context.Background()

	drv.Define("flag", func(_ context.Context, _ any) (any, error) {
		calls++
		return true, nil
	})

	dec := featureflags.NewDecoratorWithDispatcher(drv, dispatcher)

	// Resolve once to populate cache.
	dec.Get(ctx, "flag", nil) //nolint:errcheck
	if calls != 1 {
		t.Fatalf("setup: expected 1 call, got %d", calls)
	}

	if err := dec.Delete(ctx, "flag", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dispatcher.count("FeatureDeleted") != 1 {
		t.Fatalf("expected 1 FeatureDeleted, got %d", dispatcher.count("FeatureDeleted"))
	}

	// Next Get must go back to the driver because cache entry was removed.
	dec.Get(ctx, "flag", nil) //nolint:errcheck
	if calls != 2 {
		t.Fatalf("expected 2 driver calls after delete, got %d", calls)
	}
}

// ---------------------------------------------------------------------------
// Purge
// ---------------------------------------------------------------------------

func TestDecorator_Purge_Nil_ClearsAllCache_DispatchesAllFeaturesPurged(t *testing.T) {
	t.Parallel()

	dispatcher := &testDispatcher{}
	drv := featureflags.NewArrayDriver()
	ctx := context.Background()

	drv.Define("flag-a", func(_ context.Context, _ any) (any, error) { return true, nil })
	drv.Define("flag-b", func(_ context.Context, _ any) (any, error) { return true, nil })

	dec := featureflags.NewDecoratorWithDispatcher(drv, dispatcher)

	dec.Get(ctx, "flag-a", nil) //nolint:errcheck
	dec.Get(ctx, "flag-b", nil) //nolint:errcheck

	if err := dec.Purge(ctx, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dispatcher.count("AllFeaturesPurged") != 1 {
		t.Fatalf("expected 1 AllFeaturesPurged, got %d", dispatcher.count("AllFeaturesPurged"))
	}
	if dispatcher.count("FeaturesPurged") != 0 {
		t.Fatalf("expected 0 FeaturesPurged, got %d", dispatcher.count("FeaturesPurged"))
	}
}

func TestDecorator_Purge_List_ClearsNamedFeatures_DispatchesFeaturesPurged(t *testing.T) {
	t.Parallel()

	dispatcher := &testDispatcher{}
	drv := featureflags.NewArrayDriver()
	ctx := context.Background()

	drv.Define("flag-a", func(_ context.Context, _ any) (any, error) { return true, nil })
	drv.Define("flag-b", func(_ context.Context, _ any) (any, error) { return "v1", nil })

	dec := featureflags.NewDecoratorWithDispatcher(drv, dispatcher)

	dec.Get(ctx, "flag-a", nil) //nolint:errcheck
	dec.Get(ctx, "flag-b", nil) //nolint:errcheck

	if err := dec.Purge(ctx, []string{"flag-a"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dispatcher.count("FeaturesPurged") != 1 {
		t.Fatalf("expected 1 FeaturesPurged, got %d", dispatcher.count("FeaturesPurged"))
	}

	ev, ok := dispatcher.last().(featureflags.FeaturesPurged)
	if !ok {
		t.Fatalf("last event is not FeaturesPurged: %T", dispatcher.last())
	}
	if len(ev.Features) != 1 || ev.Features[0] != "flag-a" {
		t.Fatalf("expected Features=[flag-a], got %v", ev.Features)
	}
}

func TestDecorator_Purge_EmptySlice_IsNoOp_NoEvents(t *testing.T) {
	t.Parallel()

	dispatcher := &testDispatcher{}
	drv := featureflags.NewArrayDriver()
	ctx := context.Background()

	drv.Define("flag", func(_ context.Context, _ any) (any, error) { return true, nil })

	dec := featureflags.NewDecoratorWithDispatcher(drv, dispatcher)

	dec.Get(ctx, "flag", nil) //nolint:errcheck

	if err := dec.Purge(ctx, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	totalEvents := dispatcher.count("FeaturesPurged") + dispatcher.count("AllFeaturesPurged")
	if totalEvents != 0 {
		t.Fatalf("expected no purge events for empty slice, got %d", totalEvents)
	}
}

// ---------------------------------------------------------------------------
// FlushCache
// ---------------------------------------------------------------------------

func TestDecorator_FlushCache_ClearsWithoutEvents(t *testing.T) {
	t.Parallel()

	dispatcher := &testDispatcher{}
	drv := featureflags.NewArrayDriver()
	ctx := context.Background()

	drv.Define("flag", func(_ context.Context, _ any) (any, error) {
		return true, nil
	})

	dec := featureflags.NewDecoratorWithDispatcher(drv, dispatcher)

	// First Get: cache miss → FeatureResolved dispatched.
	dec.Get(ctx, "flag", nil) //nolint:errcheck

	resolvedBefore := dispatcher.count("FeatureResolved")
	if resolvedBefore != 1 {
		t.Fatalf("setup: expected 1 FeatureResolved, got %d", resolvedBefore)
	}

	// Snapshot total event count before flush — flush must add zero.
	purgeEventsBefore := dispatcher.count("AllFeaturesPurged") + dispatcher.count("FeaturesPurged")

	dec.FlushCache()

	purgeEventsAfter := dispatcher.count("AllFeaturesPurged") + dispatcher.count("FeaturesPurged")
	if purgeEventsAfter != purgeEventsBefore {
		t.Fatalf("FlushCache must not dispatch purge events (before=%d, after=%d)", purgeEventsBefore, purgeEventsAfter)
	}

	// After flush the decorator cache is empty, so the next Get is a cache miss
	// and FeatureResolved is dispatched again (even though the driver still has
	// the value stored).
	dec.Get(ctx, "flag", nil) //nolint:errcheck

	resolvedAfter := dispatcher.count("FeatureResolved")
	if resolvedAfter != 2 {
		t.Fatalf("expected 2 FeatureResolved events after FlushCache + Get, got %d", resolvedAfter)
	}
}

// ---------------------------------------------------------------------------
// Nil dispatcher
// ---------------------------------------------------------------------------

func TestDecorator_NilDispatcher_NoPanic(t *testing.T) {
	t.Parallel()

	drv := featureflags.NewArrayDriver()
	ctx := context.Background()

	drv.Define("flag", func(_ context.Context, _ any) (any, error) { return true, nil })

	// NewDecorator has a nil dispatcher — operations must not panic.
	dec := featureflags.NewDecorator(drv)

	if err := dec.Set(ctx, "flag", nil, true); err != nil {
		t.Fatalf("Set panicked or returned error: %v", err)
	}

	dec.Get(ctx, "flag", nil)             //nolint:errcheck
	dec.SetForAllScopes(ctx, "flag", false) //nolint:errcheck
	dec.Delete(ctx, "flag", nil)          //nolint:errcheck
	dec.Purge(ctx, nil)                   //nolint:errcheck
	dec.Purge(ctx, []string{"flag"})      //nolint:errcheck
	dec.FlushCache()
}

// ---------------------------------------------------------------------------
// GetAll
// ---------------------------------------------------------------------------

func TestDecorator_GetAll_MixedCacheHitsMisses(t *testing.T) {
	t.Parallel()

	calls := map[string]int{}
	var mu sync.Mutex

	drv := featureflags.NewArrayDriver()
	ctx := context.Background()

	drv.Define("flag-a", func(_ context.Context, _ any) (any, error) {
		mu.Lock()
		calls["flag-a"]++
		mu.Unlock()
		return "a", nil
	})

	drv.Define("flag-b", func(_ context.Context, _ any) (any, error) {
		mu.Lock()
		calls["flag-b"]++
		mu.Unlock()
		return "b", nil
	})

	dec := featureflags.NewDecorator(drv)

	// Prime the cache for flag-a / nil scope.
	dec.Get(ctx, "flag-a", nil) //nolint:errcheck

	mu.Lock()
	callsBeforeGetAll := calls["flag-a"]
	mu.Unlock()
	if callsBeforeGetAll != 1 {
		t.Fatalf("setup: expected 1 call for flag-a, got %d", callsBeforeGetAll)
	}

	result, err := dec.GetAll(ctx, map[string][]any{
		"flag-a": {nil},        // cache hit — driver should NOT be called
		"flag-b": {"user:1"},   // cache miss — driver must be called
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["flag-a"][0] != "a" {
		t.Fatalf("expected a for flag-a, got %v", result["flag-a"][0])
	}
	if result["flag-b"][0] != "b" {
		t.Fatalf("expected b for flag-b, got %v", result["flag-b"][0])
	}

	mu.Lock()
	callsAfter := calls["flag-a"]
	mu.Unlock()

	// flag-a was already cached — driver must not be called again.
	if callsAfter != 1 {
		t.Fatalf("expected flag-a driver calls to remain 1 (cache hit), got %d", callsAfter)
	}
}

// ---------------------------------------------------------------------------
// Define / Defined
// ---------------------------------------------------------------------------

func TestDecorator_Define_PassedToDriver(t *testing.T) {
	t.Parallel()

	drv := featureflags.NewArrayDriver()
	ctx := context.Background()

	dec := featureflags.NewDecorator(drv)

	dec.Define("flag", func(_ context.Context, _ any) (any, error) {
		return "from-decorator-define", nil
	})

	val, err := dec.Get(ctx, "flag", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "from-decorator-define" {
		t.Fatalf("expected from-decorator-define, got %v", val)
	}
}

func TestDecorator_Defined_DelegatestoDriver(t *testing.T) {
	t.Parallel()

	drv := featureflags.NewArrayDriver()

	drv.Define("flag-x", func(_ context.Context, _ any) (any, error) { return true, nil })
	drv.Define("flag-y", func(_ context.Context, _ any) (any, error) { return true, nil })

	dec := featureflags.NewDecorator(drv)

	names := dec.Defined()

	if len(names) != 2 {
		t.Fatalf("expected 2 defined features, got %d: %v", len(names), names)
	}
}
