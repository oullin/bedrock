package featureflags

import (
	"context"
	"fmt"
	"sync"
)

// ArrayDriver is a thread-safe in-memory feature-flag driver. Resolved values
// are lost when the process exits. It is the canonical driver for testing and
// short-lived scenarios.
type ArrayDriver struct {
	mu             sync.RWMutex
	resolvers      map[string]func(ctx context.Context, scope any) (any, error)
	resolvedStates map[string]map[string]any // feature → serializedScope → value
	dispatcher     EventDispatcher            // nil = no events
	inflight       sync.Map                  // inflightKey → *resolveOnce
}

// resolveOnce holds the result of a single resolver invocation. It guarantees
// the resolver runs exactly once per (feature, scope) pair under concurrent
// access, acting as a lightweight singleflight without external dependencies.
type resolveOnce struct {
	once sync.Once
	val  any
	err  error
}

var _ Driver               = (*ArrayDriver)(nil)
var _ StoredFeaturesLister = (*ArrayDriver)(nil)
var _ BulkFeatureSetter    = (*ArrayDriver)(nil)

// NewArrayDriver creates an ArrayDriver with no event dispatcher.
func NewArrayDriver() *ArrayDriver {
	return &ArrayDriver{
		resolvers:      make(map[string]func(ctx context.Context, scope any) (any, error)),
		resolvedStates: make(map[string]map[string]any),
	}
}

// NewArrayDriverWithDispatcher creates an ArrayDriver that dispatches events
// via d.
func NewArrayDriverWithDispatcher(d EventDispatcher) *ArrayDriver {
	drv := NewArrayDriver()
	drv.dispatcher = d

	return drv
}

// Define registers a resolver for a named feature.
func (a *ArrayDriver) Define(name string, resolver func(ctx context.Context, scope any) (any, error)) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.resolvers[name] = resolver
}

// Defined returns the names of all features with registered resolvers.
func (a *ArrayDriver) Defined() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	names := make([]string, 0, len(a.resolvers))

	for name := range a.resolvers {
		names = append(names, name)
	}

	return names
}

// Get resolves a single feature for a single scope. The resolver is called
// exactly once per (feature, scope) pair; concurrent callers share the result
// via a per-key sync.Once stored in the inflight map.
func (a *ArrayDriver) Get(ctx context.Context, feature string, scope any) (any, error) {
	key, err := SerializeScope(scope)
	if err != nil {
		return nil, err
	}

	// Fast path: already resolved.
	a.mu.RLock()
	if scopes, ok := a.resolvedStates[feature]; ok {
		if val, ok := scopes[key]; ok {
			a.mu.RUnlock()

			return val, nil
		}
	}

	resolver := a.resolvers[feature]
	a.mu.RUnlock()

	if resolver == nil {
		a.dispatch(ctx, UnknownFeatureResolved{Feature: feature, Scope: scope})

		return nil, fmt.Errorf("%w: %q", ErrFeatureNotDefined, feature)
	}

	// Use a per-(feature, scope) sync.Once so the resolver is called exactly
	// once regardless of how many goroutines race past the fast path above.
	inflightKey := feature + "\x00" + key
	ro := &resolveOnce{}
	actual, _ := a.inflight.LoadOrStore(inflightKey, ro)
	ro = actual.(*resolveOnce)

	ro.once.Do(func() {
		ro.val, ro.err = resolver(ctx, scope)

		if ro.err == nil {
			a.mu.Lock()

			if _, ok := a.resolvedStates[feature]; !ok {
				a.resolvedStates[feature] = make(map[string]any)
			}

			a.resolvedStates[feature][key] = ro.val
			a.mu.Unlock()
		} else {
			// Remove the inflight entry on error so the next caller can retry.
			a.inflight.Delete(inflightKey)
		}
	})

	return ro.val, ro.err
}

// GetAll resolves multiple features for multiple scopes. The returned map is
// parallel-indexed: result[feature][i] corresponds to features[feature][i].
func (a *ArrayDriver) GetAll(ctx context.Context, features map[string][]any) (map[string][]any, error) {
	result := make(map[string][]any, len(features))

	for feature, scopes := range features {
		values := make([]any, len(scopes))

		for i, scope := range scopes {
			val, err := a.Get(ctx, feature, scope)
			if err != nil {
				return nil, err
			}

			values[i] = val
		}

		result[feature] = values
	}

	return result, nil
}

// Set stores a resolved value for the given feature and scope, bypassing the
// resolver.
func (a *ArrayDriver) Set(_ context.Context, feature string, scope any, value any) error {
	key, err := SerializeScope(scope)
	if err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.resolvedStates[feature]; !ok {
		a.resolvedStates[feature] = make(map[string]any)
	}

	a.resolvedStates[feature][key] = value

	return nil
}

// SetAll stores multiple (feature, scope, value) entries atomically.
func (a *ArrayDriver) SetAll(_ context.Context, entries []FeatureEntry) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, e := range entries {
		key, err := SerializeScope(e.Scope)
		if err != nil {
			return err
		}

		if _, ok := a.resolvedStates[e.Feature]; !ok {
			a.resolvedStates[e.Feature] = make(map[string]any)
		}

		a.resolvedStates[e.Feature][key] = e.Value
	}

	return nil
}

// SetForAllScopes updates the resolved value for every stored scope of the
// given feature.
func (a *ArrayDriver) SetForAllScopes(_ context.Context, feature string, value any) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	scopes, ok := a.resolvedStates[feature]
	if !ok {
		return nil
	}

	for key := range scopes {
		scopes[key] = value
	}

	return nil
}

// Delete removes the stored resolved value for the given feature and scope.
func (a *ArrayDriver) Delete(_ context.Context, feature string, scope any) error {
	key, err := SerializeScope(scope)
	if err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if scopes, ok := a.resolvedStates[feature]; ok {
		delete(scopes, key)
	}

	return nil
}

// Purge removes stored state. nil purges all features; a non-nil empty slice
// is a no-op; a non-empty slice purges only the named features.
func (a *ArrayDriver) Purge(_ context.Context, features []string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if features == nil {
		a.resolvedStates = make(map[string]map[string]any)
		return nil
	}

	for _, name := range features {
		delete(a.resolvedStates, name)
	}

	return nil
}

// Stored returns the names of all features that have at least one stored
// scope entry.
func (a *ArrayDriver) Stored(_ context.Context) ([]string, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	names := make([]string, 0, len(a.resolvedStates))

	for name, scopes := range a.resolvedStates {
		if len(scopes) > 0 {
			names = append(names, name)
		}
	}

	return names, nil
}

func (a *ArrayDriver) dispatch(ctx context.Context, event Event) {
	if a.dispatcher != nil {
		a.dispatcher.Dispatch(ctx, event)
	}
}
