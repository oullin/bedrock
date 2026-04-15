package pennant

import (
	"context"
	"sync"
)

// Decorator wraps a Driver and adds in-process caching and event dispatch.
// Every public API operation goes through the Decorator. Cache hits avoid
// a round-trip to the underlying driver entirely.
//
// Thread-safety: read operations use RLock; all cache mutations use Lock.
type Decorator struct {
	mu         sync.RWMutex
	driver     Driver
	cache      map[string]map[string]any // feature → serializedScope → value
	dispatcher EventDispatcher           // nil = no-op
	serialize  func(any) (string, error) // defaults to SerializeScope
}

var _ Driver       = (*Decorator)(nil)
var _ CacheFlusher = (*Decorator)(nil)

// NewDecorator wraps driver with in-process caching but no event dispatch.
func NewDecorator(driver Driver) *Decorator {
	return &Decorator{
		driver:    driver,
		cache:     make(map[string]map[string]any),
		serialize: SerializeScope,
	}
}

// NewDecoratorWithDispatcher wraps driver with in-process caching and
// dispatches lifecycle events via d.
func NewDecoratorWithDispatcher(driver Driver, d EventDispatcher) *Decorator {
	dec := NewDecorator(driver)
	dec.dispatcher = d

	return dec
}

// Define registers a resolver for a named feature. Delegates to the
// underlying driver; no caching is involved.
func (d *Decorator) Define(name string, resolver func(ctx context.Context, scope any) (any, error)) {
	d.driver.Define(name, resolver)
}

// Defined returns the names of all features with registered resolvers.
// Delegates to the underlying driver.
func (d *Decorator) Defined() []string {
	return d.driver.Defined()
}

// Get resolves a single feature for a single scope. The in-process cache is
// checked first; on a miss the underlying driver is called and the result is
// stored in the cache. FeatureResolved is dispatched for every cache miss
// resolution.
func (d *Decorator) Get(ctx context.Context, feature string, scope any) (any, error) {
	key, err := d.serialize(scope)
	if err != nil {
		return nil, err
	}

	// Fast path: cache hit.
	d.mu.RLock()
	if scopes, ok := d.cache[feature]; ok {
		if val, ok := scopes[key]; ok {
			d.mu.RUnlock()

			return val, nil
		}
	}
	d.mu.RUnlock()

	// Call driver outside any lock — it may be slow.
	value, err := d.driver.Get(ctx, feature, scope)
	if err != nil {
		return nil, err
	}

	// Store result in cache.
	d.mu.Lock()
	if _, ok := d.cache[feature]; !ok {
		d.cache[feature] = make(map[string]any)
	}
	d.cache[feature][key] = value
	d.mu.Unlock()

	d.dispatch(ctx, FeatureResolved{Feature: feature, Scope: scope, Value: value})

	return value, nil
}

// GetAll resolves multiple features for multiple scopes. Cache hits are
// returned directly; only misses are forwarded to the underlying driver.
// FeatureResolved is dispatched for each freshly resolved entry.
func (d *Decorator) GetAll(ctx context.Context, features map[string][]any) (map[string][]any, error) {
	result := make(map[string][]any, len(features))

	// misses[feature] = list of (index, scope) pairs that were not cached.
	type miss struct {
		idx   int
		scope any
	}
	misses := make(map[string][]miss)

	// Allocate result slices and collect cache hits.
	for feature, scopes := range features {
		result[feature] = make([]any, len(scopes))

		for i, scope := range scopes {
			key, err := d.serialize(scope)
			if err != nil {
				return nil, err
			}

			d.mu.RLock()
			var cached any
			var hit bool
			if scopeMap, ok := d.cache[feature]; ok {
				cached, hit = scopeMap[key]
			}
			d.mu.RUnlock()

			if hit {
				result[feature][i] = cached
			} else {
				misses[feature] = append(misses[feature], miss{idx: i, scope: scope})
			}
		}
	}

	if len(misses) == 0 {
		return result, nil
	}

	// Build the input map for the driver containing only misses.
	driverInput := make(map[string][]any, len(misses))
	for feature, ms := range misses {
		scopes := make([]any, len(ms))
		for j, m := range ms {
			scopes[j] = m.scope
		}
		driverInput[feature] = scopes
	}

	driverResult, err := d.driver.GetAll(ctx, driverInput)
	if err != nil {
		return nil, err
	}

	// Merge driver results into cache and result.
	for feature, ms := range misses {
		vals := driverResult[feature]

		for j, m := range ms {
			value := vals[j]

			key, err := d.serialize(m.scope)
			if err != nil {
				return nil, err
			}

			d.mu.Lock()
			if _, ok := d.cache[feature]; !ok {
				d.cache[feature] = make(map[string]any)
			}
			d.cache[feature][key] = value
			d.mu.Unlock()

			d.dispatch(ctx, FeatureResolved{Feature: feature, Scope: m.scope, Value: value})

			result[feature][m.idx] = value
		}
	}

	return result, nil
}

// Set stores a resolved value for the given feature and scope, bypassing the
// resolver. The cache is updated on success and FeatureUpdated is dispatched.
func (d *Decorator) Set(ctx context.Context, feature string, scope any, value any) error {
	if err := d.driver.Set(ctx, feature, scope, value); err != nil {
		return err
	}

	key, err := d.serialize(scope)
	if err != nil {
		return err
	}

	d.mu.Lock()
	if _, ok := d.cache[feature]; !ok {
		d.cache[feature] = make(map[string]any)
	}
	d.cache[feature][key] = value
	d.mu.Unlock()

	d.dispatch(ctx, FeatureUpdated{Feature: feature, Scope: scope, Value: value})

	return nil
}

// SetForAllScopes updates the resolved value for every scope that already has
// stored state for the feature. All cached entries for the feature are
// invalidated and FeatureUpdatedForAllScopes is dispatched.
func (d *Decorator) SetForAllScopes(ctx context.Context, feature string, value any) error {
	if err := d.driver.SetForAllScopes(ctx, feature, value); err != nil {
		return err
	}

	d.mu.Lock()
	delete(d.cache, feature)
	d.mu.Unlock()

	d.dispatch(ctx, FeatureUpdatedForAllScopes{Feature: feature, Value: value})

	return nil
}

// Delete removes the stored resolved value for the given feature and scope.
// The corresponding cache entry is removed on success and FeatureDeleted is
// dispatched.
func (d *Decorator) Delete(ctx context.Context, feature string, scope any) error {
	if err := d.driver.Delete(ctx, feature, scope); err != nil {
		return err
	}

	key, err := d.serialize(scope)
	if err != nil {
		return err
	}

	d.mu.Lock()
	if scopes, ok := d.cache[feature]; ok {
		delete(scopes, key)
	}
	d.mu.Unlock()

	d.dispatch(ctx, FeatureDeleted{Feature: feature, Scope: scope})

	return nil
}

// Purge removes stored state for the given features.
//   - nil: purges every feature, replaces entire cache, dispatches AllFeaturesPurged.
//   - non-empty slice: purges named features, removes them from cache, dispatches FeaturesPurged.
//   - empty non-nil slice: delegates no-op to driver, no cache changes, no events.
func (d *Decorator) Purge(ctx context.Context, features []string) error {
	if features == nil {
		if err := d.driver.Purge(ctx, nil); err != nil {
			return err
		}

		d.mu.Lock()
		d.cache = make(map[string]map[string]any)
		d.mu.Unlock()

		d.dispatch(ctx, AllFeaturesPurged{})

		return nil
	}

	if len(features) == 0 {
		return d.driver.Purge(ctx, []string{})
	}

	if err := d.driver.Purge(ctx, features); err != nil {
		return err
	}

	d.mu.Lock()
	for _, name := range features {
		delete(d.cache, name)
	}
	d.mu.Unlock()

	d.dispatch(ctx, FeaturesPurged{Features: features})

	return nil
}

// SetAll stores multiple (feature, scope, value) entries. If the underlying
// driver implements BulkFeatureSetter its SetAll is used; otherwise individual
// Set calls are made. The cache is updated for each entry.
func (d *Decorator) SetAll(ctx context.Context, entries []FeatureEntry) error {
	if bulk, ok := d.driver.(BulkFeatureSetter); ok {
		if err := bulk.SetAll(ctx, entries); err != nil {
			return err
		}

		for _, e := range entries {
			key, err := d.serialize(e.Scope)
			if err != nil {
				return err
			}

			d.mu.Lock()
			if _, ok := d.cache[e.Feature]; !ok {
				d.cache[e.Feature] = make(map[string]any)
			}
			d.cache[e.Feature][key] = e.Value
			d.mu.Unlock()
		}

		return nil
	}

	for _, e := range entries {
		if err := d.Set(ctx, e.Feature, e.Scope, e.Value); err != nil {
			return err
		}
	}

	return nil
}

// Stored returns the names of all features with persisted state. If the
// underlying driver implements StoredFeaturesLister it is delegated to;
// otherwise nil, nil is returned.
func (d *Decorator) Stored(ctx context.Context) ([]string, error) {
	if lister, ok := d.driver.(StoredFeaturesLister); ok {
		return lister.Stored(ctx)
	}

	return nil, nil
}

// FlushCache replaces the in-process cache with an empty map. No events are
// dispatched.
func (d *Decorator) FlushCache() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.cache = make(map[string]map[string]any)
}

// dispatch sends an event to the dispatcher when one is configured.
func (d *Decorator) dispatch(ctx context.Context, event Event) {
	if d.dispatcher != nil {
		d.dispatcher.Dispatch(ctx, event)
	}
}
