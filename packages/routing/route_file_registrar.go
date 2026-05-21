package routing

// RouteFileRegistrar is the Go counterpart of the PHP class of the same name.
//
// In PHP it requires a routes file (which executes top-level Route::*
// statements against $router). Go has no equivalent of `require` injecting
// variables into a file's scope; instead, this Go form accepts a Go function
// that receives the router and registers routes against it. The helper exists
// so the Router::Group machinery can treat closures and "route file" loaders
// uniformly.
//
// Ref: @bedrock/code-0337
type RouteFileRegistrar struct{ router *Router }

// NewRouteFileRegistrar wraps the router so the resulting value can be used
// as a parity-named alternative to a closure in [Router.Group].
func NewRouteFileRegistrar(router *Router) *RouteFileRegistrar {
	return &RouteFileRegistrar{router: router}
}

// Register invokes the supplied route loader against the wrapped router.
func (r *RouteFileRegistrar) Register(loader func(*Router)) { loader(r.router) }
