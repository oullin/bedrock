package pennant

import "context"

// ScopedFeatureInteraction provides a fluent scoped API for evaluating and
// managing feature flags against a fixed set of scopes. It is the Go
// equivalent of the upstream Pennant PendingScopedFeatureInteraction.
type ScopedFeatureInteraction struct {
	decorator *Decorator
	scopes    []any
}

// NewScopedFeatureInteraction creates a ScopedFeatureInteraction bound to the
// given Decorator and scopes.
func NewScopedFeatureInteraction(decorator *Decorator, scopes ...any) *ScopedFeatureInteraction {
	return &ScopedFeatureInteraction{
		decorator: decorator,
		scopes:    scopes,
	}
}

// For returns a new ScopedFeatureInteraction that merges the current scopes
// with the supplied ones. The receiver is never mutated.
func (s *ScopedFeatureInteraction) For(scopes ...any) *ScopedFeatureInteraction {
	merged := make([]any, 0, len(s.scopes)+len(scopes))
	merged = append(merged, s.scopes...)
	merged = append(merged, scopes...)

	return NewScopedFeatureInteraction(s.decorator, merged...)
}

// resolveScopes returns the scopes slice. When empty it returns []any{nil} so
// that subsequent operations use the global (nil) scope.
func (s *ScopedFeatureInteraction) resolveScopes() []any {
	if len(s.scopes) == 0 {
		return []any{nil}
	}

	return s.scopes
}

// isTruthy reports whether v counts as an active/truthy feature value.
// false (bool) and nil are the only falsy values.
func isTruthy(v any) bool {
	if v == nil {
		return false
	}

	if b, ok := v.(bool); ok {
		return b
	}

	return true
}

// Active reports whether feature evaluates to a truthy value for the first
// scope. Returns false on any error.
func (s *ScopedFeatureInteraction) Active(ctx context.Context, feature string) bool {
	scope := s.resolveScopes()[0]

	val, err := s.decorator.Get(ctx, feature, scope)

	if err != nil {
		return false
	}

	return isTruthy(val)
}

// Inactive is the inverse of Active.
func (s *ScopedFeatureInteraction) Inactive(ctx context.Context, feature string) bool {
	return !s.Active(ctx, feature)
}

// Value returns the raw resolved value for the feature using the only scope.
// It returns ErrMultipleScopes when the interaction has more than one scope.
func (s *ScopedFeatureInteraction) Value(ctx context.Context, feature string) (any, error) {
	scopes := s.resolveScopes()

	if len(scopes) > 1 {
		return nil, ErrMultipleScopes
	}

	return s.decorator.Get(ctx, feature, scopes[0])
}

// Values returns resolved values for multiple features using the first scope.
func (s *ScopedFeatureInteraction) Values(ctx context.Context, features []string) (map[string]any, error) {
	scope := s.resolveScopes()[0]

	input := make(map[string][]any, len(features))

	for _, f := range features {
		input[f] = []any{scope}
	}

	all, err := s.decorator.GetAll(ctx, input)

	if err != nil {
		return nil, err
	}

	out := make(map[string]any, len(features))

	for _, f := range features {
		if vals, ok := all[f]; ok && len(vals) > 0 {
			out[f] = vals[0]
		}
	}

	return out, nil
}

// AllAreActive reports whether all named features are active for ALL scopes.
func (s *ScopedFeatureInteraction) AllAreActive(ctx context.Context, features []string) bool {
	scopes := s.resolveScopes()

	for _, feature := range features {
		for _, scope := range scopes {
			val, err := s.decorator.Get(ctx, feature, scope)

			if err != nil || !isTruthy(val) {
				return false
			}
		}
	}

	return true
}

// SomeAreActive reports whether at least one (feature, scope) pair is active.
func (s *ScopedFeatureInteraction) SomeAreActive(ctx context.Context, features []string) bool {
	scopes := s.resolveScopes()

	for _, feature := range features {
		for _, scope := range scopes {
			val, err := s.decorator.Get(ctx, feature, scope)

			if err == nil && isTruthy(val) {
				return true
			}
		}
	}

	return false
}

// AllAreInactive reports whether all named features are inactive for ALL scopes.
func (s *ScopedFeatureInteraction) AllAreInactive(ctx context.Context, features []string) bool {
	scopes := s.resolveScopes()

	for _, feature := range features {
		for _, scope := range scopes {
			val, err := s.decorator.Get(ctx, feature, scope)

			if err == nil && isTruthy(val) {
				return false
			}
		}
	}

	return true
}

// SomeAreInactive reports whether at least one (feature, scope) pair is inactive.
func (s *ScopedFeatureInteraction) SomeAreInactive(ctx context.Context, features []string) bool {
	scopes := s.resolveScopes()

	for _, feature := range features {
		for _, scope := range scopes {
			val, err := s.decorator.Get(ctx, feature, scope)

			if err != nil || !isTruthy(val) {
				return true
			}
		}
	}

	return false
}

// Activate sets feature(s) to true for all scopes.
func (s *ScopedFeatureInteraction) Activate(ctx context.Context, features []string) error {
	return s.ActivateWithValue(ctx, features, true)
}

// ActivateWithValue sets feature(s) to a specific value for all scopes.
func (s *ScopedFeatureInteraction) ActivateWithValue(ctx context.Context, features []string, value any) error {
	scopes := s.resolveScopes()

	for _, feature := range features {
		for _, scope := range scopes {
			if err := s.decorator.Set(ctx, feature, scope, value); err != nil {
				return err
			}
		}
	}

	return nil
}

// Deactivate sets feature(s) to false for all scopes.
func (s *ScopedFeatureInteraction) Deactivate(ctx context.Context, features []string) error {
	return s.ActivateWithValue(ctx, features, false)
}

// Forget removes the stored state for feature(s) for all scopes.
func (s *ScopedFeatureInteraction) Forget(ctx context.Context, features []string) error {
	scopes := s.resolveScopes()

	for _, feature := range features {
		for _, scope := range scopes {
			if err := s.decorator.Delete(ctx, feature, scope); err != nil {
				return err
			}
		}
	}

	return nil
}

// Purge removes all stored state for the named features. nil purges everything.
func (s *ScopedFeatureInteraction) Purge(ctx context.Context, features []string) error {
	return s.decorator.Purge(ctx, features)
}

// When executes whenActive if the feature is active for the first scope,
// otherwise executes whenInactive. Either callback may be nil (no-op).
func (s *ScopedFeatureInteraction) When(
	ctx context.Context,
	feature string,
	whenActive func(value any) (any, error),
	whenInactive func(value any) (any, error),
) (any, error) {
	val, err := s.Value(ctx, feature)

	if err != nil && whenInactive != nil {
		return whenInactive(nil)
	}

	if isTruthy(val) {
		if whenActive != nil {
			return whenActive(val)
		}

		return nil, nil
	}

	if whenInactive != nil {
		return whenInactive(val)
	}

	return nil, nil
}

// Unless is the inverse of When: executes whenInactive if the feature is
// active, and whenActive if the feature is inactive.
func (s *ScopedFeatureInteraction) Unless(
	ctx context.Context,
	feature string,
	whenInactive func(value any) (any, error),
	whenActive func(value any) (any, error),
) (any, error) {
	return s.When(ctx, feature, whenActive, whenInactive)
}

// Load eagerly resolves the named features for all scopes, populating the
// Decorator cache. Results are discarded; only the side-effect matters.
func (s *ScopedFeatureInteraction) Load(ctx context.Context, features []string) error {
	scopes := s.resolveScopes()

	input := make(map[string][]any, len(features))

	for _, f := range features {
		input[f] = scopes
	}

	_, err := s.decorator.GetAll(ctx, input)

	return err
}

// LoadMissing resolves only the features not already present in the Decorator
// cache. Because GetAll itself skips cached entries, forwarding to Load is
// sufficient.
func (s *ScopedFeatureInteraction) LoadMissing(ctx context.Context, features []string) error {
	return s.Load(ctx, features)
}

// LoadAll eagerly resolves every registered feature for all scopes.
func (s *ScopedFeatureInteraction) LoadAll(ctx context.Context) error {
	return s.Load(ctx, s.decorator.Defined())
}

// For returns a ScopedFeatureInteraction bound to the default driver.
func (d *Decorator) For(scopes ...any) *ScopedFeatureInteraction {
	return NewScopedFeatureInteraction(d, scopes...)
}

// For returns a ScopedFeatureInteraction backed by the Manager's default driver.
func (m *Manager) For(scopes ...any) (*ScopedFeatureInteraction, error) {
	if len(scopes) == 0 {
		m.mu.RLock()
		resolver := m.scopeResolver
		m.mu.RUnlock()

		if resolver != nil {
			scope, err := resolver(context.Background())

			if err != nil {
				return nil, err
			}

			scopes = []any{scope}
		}
	}

	dec, err := m.DefaultDecorator()

	if err != nil {
		return nil, err
	}

	return NewScopedFeatureInteraction(dec, scopes...), nil
}
