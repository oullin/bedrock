package featureflags

import "context"

// Driver is the low-level feature-flag backend contract. Implementations
// persist and retrieve resolved feature values. They do not add caching or
// event dispatch — that is the responsibility of the Decorator.
type Driver interface {
	// Define registers a resolver for a named feature. Calling Define a second
	// time for the same name overwrites the previous resolver.
	Define(name string, resolver func(ctx context.Context, scope any) (any, error))

	// Defined returns the names of all features with registered resolvers.
	Defined() []string

	// GetAll resolves multiple features for multiple scopes in a single call.
	// The input map keys are feature names; the values are slices of scopes.
	// The returned map has the same keys; each value slice is parallel-indexed
	// to the input — result[feature][i] corresponds to features[feature][i].
	GetAll(ctx context.Context, features map[string][]any) (map[string][]any, error)

	// Get resolves a single feature for a single scope. If the feature has
	// stored state, that value is returned. Otherwise the registered resolver is
	// invoked and the result stored. Returns ErrFeatureNotDefined if neither
	// stored state nor a resolver exists.
	Get(ctx context.Context, feature string, scope any) (any, error)

	// Set stores a resolved value for the given feature and scope, bypassing the
	// resolver.
	Set(ctx context.Context, feature string, scope any, value any) error

	// SetForAllScopes updates the resolved value for every scope that already
	// has stored state for the feature.
	SetForAllScopes(ctx context.Context, feature string, value any) error

	// Delete removes the stored resolved value for the given feature and scope.
	// The next Get call will re-invoke the resolver.
	Delete(ctx context.Context, feature string, scope any) error

	// Purge removes all stored state for the given features.
	// Passing nil purges every feature. Passing an empty non-nil slice is a
	// no-op.
	Purge(ctx context.Context, features []string) error
}

// FeatureEntry is a single (feature, scope, value) triple used by
// BulkFeatureSetter.
type FeatureEntry struct {
	Feature string
	Scope   any
	Value   any
}

// StoredFeaturesLister is optionally implemented by drivers that can
// enumerate all features with persisted state in the backend.
type StoredFeaturesLister interface {
	Stored(ctx context.Context) ([]string, error)
}

// BulkFeatureSetter is optionally implemented by drivers that support atomic
// bulk writes.
type BulkFeatureSetter interface {
	SetAll(ctx context.Context, entries []FeatureEntry) error
}

// CacheFlusher is optionally implemented by drivers (or decorators) that
// maintain an in-process cache layer.
type CacheFlusher interface {
	FlushCache()
}
