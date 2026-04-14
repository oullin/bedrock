package routing

import "github.com/bedrock/packages/routing/matching"

// RouteCollectionInterface mirrors
// Illuminate\Routing\RouteCollectionInterface.
//
// Both [RouteCollection] and [CompiledRouteCollection] satisfy this interface,
// so consumers (notably [Router]) can hold a single field and switch
// implementations in production via route caching.
type RouteCollectionInterface interface {
	// Add registers a Route in the collection and returns the same Route.
	Add(route *Route) *Route
	// RefreshNameLookups rebuilds the name → Route lookup table.
	RefreshNameLookups()
	// RefreshActionLookups rebuilds the action → Route lookup table.
	RefreshActionLookups()
	// Match returns the Route that matches the supplied request, or returns
	// an error wrapping [ErrRouteNotFound] / [MethodNotAllowedError].
	Match(request matching.MatchableRequest) (*Route, error)
	// Get returns the routes registered for the given HTTP method, or all
	// routes when method is empty.
	Get(method string) []*Route
	// HasNamedRoute reports whether a route with the given name is registered.
	HasNamedRoute(name string) bool
	// GetByName returns the named route or nil.
	GetByName(name string) *Route
	// GetByAction returns the route for the given controller action or nil.
	GetByAction(action string) *Route
	// GetRoutes returns every route in registration order.
	GetRoutes() []*Route
	// GetRoutesByMethod returns the routes-by-method map.
	GetRoutesByMethod() map[string][]*Route
	// GetRoutesByName returns the name → Route map.
	GetRoutesByName() map[string]*Route
}
