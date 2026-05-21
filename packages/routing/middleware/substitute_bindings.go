// Ref: @bedrock/code-0319
// The bedrock pipeline runs each middleware as a function that receives the
// current request and a "next" continuation. The shapes here are designed so
// they can be wrapped trivially by the bedrock pipeline package in M11.
package middleware

import "fmt"

// BindingRouter is the minimum router surface SubstituteBindings touches.
// Both [routing.Router] and any future test fake satisfy it.
type BindingRouter interface {
	SubstituteBindings(route any) error
	SubstituteImplicitBindings(route any) error
}

// SubstituteBindings is the middleware that runs explicit and implicit route
// bindings before the route handler executes.
//
// Ref: @bedrock/code-0321
type SubstituteBindings struct{ Router BindingRouter }

// New wraps a router into the middleware.
func New(router BindingRouter) *SubstituteBindings {
	return &SubstituteBindings{Router: router}
}

// Handle is the pipeline entry point: invoke the binders, then call next.
//
// route must be a *routing.Route obtained from the request's route resolver
// (M11 wires this via httpx.Request.SetRouteResolver). Returning an error
// from either binder short-circuits the pipeline; if the route declared a
// missing handler via [routing.Route.Missing] the wrapping pipeline catches
// the error and invokes it instead.
func (s *SubstituteBindings) Handle(request any, route any, next func(any) any) (any, error) {
	if route == nil {
		return nil, fmt.Errorf("substitute bindings: no route on request")
	}

	if err := s.Router.SubstituteBindings(route); err != nil {
		return nil, err
	}

	if err := s.Router.SubstituteImplicitBindings(route); err != nil {
		return nil, err
	}

	return next(request), nil
}
