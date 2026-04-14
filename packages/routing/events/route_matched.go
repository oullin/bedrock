// Package events mirrors upstream/framework/src/Framework/Routing/Events.
//
// Each struct corresponds to one PHP event class. They are dispatched by the
// router via the bedrock events package; consumers register listeners against
// the concrete struct type.
package events

// RouteMatched is dispatched after the router has resolved an incoming request
// to a specific route, before the route's middleware pipeline executes.
//
// Mirrors Framework\Routing\Events\RouteMatched.
type RouteMatched struct {
	Route   any // *routing.Route — typed as any to break import cycles.
	Request any // httpx.Request
}
