package eloquent

import "github.com/bedrock/packages/database/query"

// Scope applies a constraint to a query builder.
type Scope interface {
	Apply(builder *query.Builder)
}

// ScopeFunc is a function that acts as a Scope.
type ScopeFunc func(builder *query.Builder)

// Apply satisfies the Scope interface.

// GlobalScopes manages the global scope registry for a model.
type GlobalScopes struct {
	scopes map[string]Scope
}

func (f ScopeFunc) Apply(builder *query.Builder) { f(builder) }

// InitScopes initializes the scope map.
func (gs *GlobalScopes) InitScopes() {
	if gs.scopes == nil {
		gs.scopes = make(map[string]Scope)
	}
}

// AddGlobalScope registers a named global scope.
func (gs *GlobalScopes) AddGlobalScope(name string, scope Scope) {
	gs.InitScopes()
	gs.scopes[name] = scope
}

// RemoveGlobalScope removes a named global scope.
func (gs *GlobalScopes) RemoveGlobalScope(name string) {
	delete(gs.scopes, name)
}

// GetGlobalScopes returns all registered global scopes.
func (gs *GlobalScopes) GetGlobalScopes() map[string]Scope {
	return gs.scopes
}

// HasGlobalScope checks if a named global scope is registered.
func (gs *GlobalScopes) HasGlobalScope(name string) bool {
	_, ok := gs.scopes[name]

	return ok
}

// ApplyScopes applies all global scopes to the given builder.
func (gs *GlobalScopes) ApplyScopes(builder *query.Builder) {
	for _, scope := range gs.scopes {
		scope.Apply(builder)
	}
}
