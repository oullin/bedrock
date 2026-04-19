package featureflags

import "context"

// Event is the marker interface for all featureflags events.
type Event interface {
	FeatureFlagsEvent()
}

// EventDispatcher dispatches featureflags events to registered listeners.
// Defined locally so the featureflags package does not depend on the events package.
type EventDispatcher interface {
	Dispatch(ctx context.Context, event Event)
}

// FeatureResolved is dispatched when a feature value is resolved for a scope
// (either from storage or by invoking the resolver).
type FeatureResolved struct {
	Feature string
	Scope   any
	Value   any
}

// UnknownFeatureResolved is dispatched when Get is called for a feature that
// has no registered resolver and no stored state.
type UnknownFeatureResolved struct {
	Feature string
	Scope   any
}

// FeatureUpdated is dispatched when a feature value is explicitly set for a
// specific scope.
type FeatureUpdated struct {
	Feature string
	Scope   any
	Value   any
}

// FeatureDeleted is dispatched when a feature's stored state is removed for a
// specific scope.
type FeatureDeleted struct {
	Feature string
	Scope   any
}

// FeatureUpdatedForAllScopes is dispatched when a feature value is set across
// all stored scopes via SetForAllScopes.
type FeatureUpdatedForAllScopes struct {
	Feature string
	Value   any
}

// FeaturesPurged is dispatched when specific named features are purged.
type FeaturesPurged struct {
	Features []string
}

// AllFeaturesPurged is dispatched when all feature state is purged (Purge with
// a nil features list).
type AllFeaturesPurged struct{}

func (FeatureResolved) FeatureFlagsEvent()            {}
func (UnknownFeatureResolved) FeatureFlagsEvent()     {}
func (FeatureUpdated) FeatureFlagsEvent()             {}
func (FeatureDeleted) FeatureFlagsEvent()             {}
func (FeatureUpdatedForAllScopes) FeatureFlagsEvent() {}
func (FeaturesPurged) FeatureFlagsEvent()             {}
func (AllFeaturesPurged) FeatureFlagsEvent()          {}

var _ Event = FeatureResolved{}
var _ Event = UnknownFeatureResolved{}
var _ Event = FeatureUpdated{}
var _ Event = FeatureDeleted{}
var _ Event = FeatureUpdatedForAllScopes{}
var _ Event = FeaturesPurged{}
var _ Event = AllFeaturesPurged{}
