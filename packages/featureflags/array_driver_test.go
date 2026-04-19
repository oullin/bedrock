package pennant_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/bedrock/packages/featureflags"
)

//nolint:errcheck

//nolint:errcheck
//nolint:errcheck

//nolint:errcheck
//nolint:errcheck

// Resolve for two different scopes.
//nolint:errcheck
//nolint:errcheck

// Deactivate for all scopes.

//nolint:errcheck

// Next Get should re-invoke resolver.
//nolint:errcheck

//nolint:errcheck
//nolint:errcheck

//nolint:errcheck
//nolint:errcheck

//nolint:errcheck

//nolint:errcheck

//nolint:errcheck

// testDispatcher is a simple in-memory EventDispatcher for tests.
type testDispatcher struct {
	mu     sync.Mutex
	events []featureflags.Event
}

func TestArrayDriver_InterfaceAssertions(t *testing.T) {
	t.Parallel()

	var _ featureflags.Driver = (*featureflags.ArrayDriver)(nil)

	var _ featureflags.StoredFeaturesLister = (*featureflags.ArrayDriver)(nil)

	var _ featureflags.BulkFeatureSetter = (*featureflags.ArrayDriver)(nil)
}

func TestArrayDriver_Define_Get(t *testing.T) {
	t.Parallel()

	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	d.Define("dark-mode", func(_ context.Context, _ any) (any, error) {
		return true, nil
	})

	val, err := d.Get(ctx, "dark-mode", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != true {
		t.Fatalf("expected true, got %v", val)
	}
}

func TestArrayDriver_Defined(t *testing.T) {
	t.Parallel()

	d := featureflags.NewArrayDriver()

	d.Define("flag-a", func(_ context.Context, _ any) (any, error) { return true, nil })
	d.Define("flag-b", func(_ context.Context, _ any) (any, error) { return false, nil })

	names := d.Defined()

	if len(names) != 2 {
		t.Fatalf("expected 2 defined features, got %d", len(names))
	}
}

func TestArrayDriver_Get_UndefinedFeature_ReturnsError(t *testing.T) {
	t.Parallel()

	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	_, err := d.Get(ctx, "unknown", nil)

	if !errors.Is(err, featureflags.ErrFeatureNotDefined) {
		t.Fatalf("expected ErrFeatureNotDefined, got %v", err)
	}
}

func TestArrayDriver_Get_UndefinedFeature_DispatchesEvent(t *testing.T) {
	t.Parallel()

	dispatcher := &testDispatcher{}
	d := featureflags.NewArrayDriverWithDispatcher(dispatcher)
	ctx := context.Background()

	d.Get(ctx, "unknown", nil)

	if dispatcher.count("UnknownFeatureResolved") != 1 {
		t.Fatalf("expected 1 UnknownFeatureResolved event, got %d", dispatcher.count("UnknownFeatureResolved"))
	}
}

func TestArrayDriver_Get_CachesResult(t *testing.T) {
	t.Parallel()

	calls := 0
	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	d.Define("flag", func(_ context.Context, _ any) (any, error) {
		calls++

		return "variant-a", nil
	})

	d.Get(ctx, "flag", nil)
	d.Get(ctx, "flag", nil)

	if calls != 1 {
		t.Fatalf("expected resolver to be called once, got %d calls", calls)
	}
}

func TestArrayDriver_Get_DifferentScopes_CallsResolverPerScope(t *testing.T) {
	t.Parallel()

	calls := 0
	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	d.Define("flag", func(_ context.Context, _ any) (any, error) {
		calls++

		return true, nil
	})

	d.Get(ctx, "flag", "user:1")
	d.Get(ctx, "flag", "user:2")

	if calls != 2 {
		t.Fatalf("expected resolver called twice, got %d", calls)
	}
}

func TestArrayDriver_Set_BypassesResolver(t *testing.T) {
	t.Parallel()

	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	d.Define("flag", func(_ context.Context, _ any) (any, error) {
		return true, nil
	})

	_ = d.Set(ctx, "flag", nil, false)

	val, err := d.Get(ctx, "flag", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != false {
		t.Fatalf("expected false (stored), got %v", val)
	}
}

func TestArrayDriver_SetAll(t *testing.T) {
	t.Parallel()

	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	entries := []featureflags.FeatureEntry{
		{Feature: "flag-a", Scope: "user:1", Value: true},
		{Feature: "flag-b", Scope: "user:1", Value: "variant"},
	}

	if err := d.SetAll(ctx, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, err := d.Get(ctx, "flag-a", "user:1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != true {
		t.Fatalf("expected true, got %v", val)
	}

	val, err = d.Get(ctx, "flag-b", "user:1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != "variant" {
		t.Fatalf("expected %q, got %v", "variant", val)
	}
}

func TestArrayDriver_SetForAllScopes(t *testing.T) {
	t.Parallel()

	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	d.Define("flag", func(_ context.Context, _ any) (any, error) { return true, nil })

	d.Get(ctx, "flag", "user:1")
	d.Get(ctx, "flag", "user:2")

	if err := d.SetForAllScopes(ctx, "flag", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, scope := range []any{"user:1", "user:2"} {
		val, err := d.Get(ctx, "flag", scope)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if val != false {
			t.Fatalf("scope %v: expected false, got %v", scope, val)
		}
	}
}

func TestArrayDriver_Delete(t *testing.T) {
	t.Parallel()

	calls := 0
	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	d.Define("flag", func(_ context.Context, _ any) (any, error) {
		calls++

		return true, nil
	})

	d.Get(ctx, "flag", nil)

	if err := d.Delete(ctx, "flag", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	d.Get(ctx, "flag", nil)

	if calls != 2 {
		t.Fatalf("expected resolver called twice (after delete), got %d", calls)
	}
}

func TestArrayDriver_Purge_All(t *testing.T) {
	t.Parallel()

	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	d.Define("flag-a", func(_ context.Context, _ any) (any, error) { return true, nil })
	d.Define("flag-b", func(_ context.Context, _ any) (any, error) { return true, nil })

	d.Get(ctx, "flag-a", nil)
	d.Get(ctx, "flag-b", nil)

	if err := d.Purge(ctx, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stored, _ := d.Stored(ctx)

	if len(stored) != 0 {
		t.Fatalf("expected no stored features after purge all, got %v", stored)
	}
}

func TestArrayDriver_Purge_Specific(t *testing.T) {
	t.Parallel()

	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	d.Define("flag-a", func(_ context.Context, _ any) (any, error) { return true, nil })
	d.Define("flag-b", func(_ context.Context, _ any) (any, error) { return true, nil })

	d.Get(ctx, "flag-a", nil)
	d.Get(ctx, "flag-b", nil)

	if err := d.Purge(ctx, []string{"flag-a"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stored, _ := d.Stored(ctx)

	if len(stored) != 1 || stored[0] != "flag-b" {
		t.Fatalf("expected only flag-b stored, got %v", stored)
	}
}

func TestArrayDriver_Purge_EmptySlice_IsNoOp(t *testing.T) {
	t.Parallel()

	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	d.Define("flag", func(_ context.Context, _ any) (any, error) { return true, nil })

	d.Get(ctx, "flag", nil)

	if err := d.Purge(ctx, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stored, _ := d.Stored(ctx)

	if len(stored) != 1 {
		t.Fatalf("expected 1 stored feature after empty-slice purge, got %v", stored)
	}
}

func TestArrayDriver_Stored(t *testing.T) {
	t.Parallel()

	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	d.Define("flag-a", func(_ context.Context, _ any) (any, error) { return true, nil })
	d.Define("flag-b", func(_ context.Context, _ any) (any, error) { return true, nil })

	d.Get(ctx, "flag-a", nil)

	stored, err := d.Stored(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(stored) != 1 || stored[0] != "flag-a" {
		t.Fatalf("expected [flag-a], got %v", stored)
	}
}

func TestArrayDriver_GetAll(t *testing.T) {
	t.Parallel()

	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	d.Define("flag-a", func(_ context.Context, _ any) (any, error) { return true, nil })
	d.Define("flag-b", func(_ context.Context, _ any) (any, error) { return "v1", nil })

	result, err := d.GetAll(ctx, map[string][]any{
		"flag-a": {nil, "user:1"},
		"flag-b": {"user:2"},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result["flag-a"]) != 2 {
		t.Fatalf("expected 2 results for flag-a, got %d", len(result["flag-a"]))
	}

	if result["flag-b"][0] != "v1" {
		t.Fatalf("expected v1, got %v", result["flag-b"][0])
	}
}

func TestArrayDriver_ConcurrentGet_ResolverCalledOnce(t *testing.T) {
	t.Parallel()

	var calls atomic.Int64
	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	d.Define("flag", func(_ context.Context, _ any) (any, error) {
		calls.Add(1)

		return true, nil
	})

	const goroutines = 100

	var wg sync.WaitGroup

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()

			d.Get(ctx, "flag", nil)
		}()
	}

	wg.Wait()

	if calls.Load() != 1 {
		t.Fatalf("resolver called %d times, expected exactly 1", calls.Load())
	}
}

func TestArrayDriver_RichValue(t *testing.T) {
	t.Parallel()

	d := featureflags.NewArrayDriver()
	ctx := context.Background()

	expected := map[string]any{"color": "blue", "size": "lg"}

	d.Define("button", func(_ context.Context, _ any) (any, error) {
		return expected, nil
	})

	val, err := d.Get(ctx, "button", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m, ok := val.(map[string]any)

	if !ok {
		t.Fatalf("expected map, got %T", val)
	}

	if m["color"] != "blue" {
		t.Fatalf("expected color=blue, got %v", m["color"])
	}
}

func (d *testDispatcher) Dispatch(_ context.Context, event featureflags.Event) {
	d.mu.Lock()

	defer d.mu.Unlock()

	d.events = append(d.events, event)
}

func (d *testDispatcher) count(typeName string) int {
	d.mu.Lock()

	defer d.mu.Unlock()

	n := 0

	for _, e := range d.events {
		switch typeName {
		case "FeatureResolved":
			if _, ok := e.(featureflags.FeatureResolved); ok {
				n++
			}
		case "UnknownFeatureResolved":
			if _, ok := e.(featureflags.UnknownFeatureResolved); ok {
				n++
			}
		case "FeatureUpdated":
			if _, ok := e.(featureflags.FeatureUpdated); ok {
				n++
			}
		case "FeatureDeleted":
			if _, ok := e.(featureflags.FeatureDeleted); ok {
				n++
			}
		case "FeatureUpdatedForAllScopes":
			if _, ok := e.(featureflags.FeatureUpdatedForAllScopes); ok {
				n++
			}
		case "FeaturesPurged":
			if _, ok := e.(featureflags.FeaturesPurged); ok {
				n++
			}
		case "AllFeaturesPurged":
			if _, ok := e.(featureflags.AllFeaturesPurged); ok {
				n++
			}
		}
	}

	return n
}

func (d *testDispatcher) last() featureflags.Event {
	d.mu.Lock()

	defer d.mu.Unlock()

	if len(d.events) == 0 {
		return nil
	}

	return d.events[len(d.events)-1]
}
