package routing

// RouteRegistrar provides a fluent interface for accumulating route attributes
// before registration. It enables patterns like:
//
//	router.WithMiddleware(mw).Prefix("/api").Group(func(r *Router) { ... })
type RouteRegistrar struct {
	router *Router
	attrs  GroupAttributes
}

// WithMiddleware returns a RouteRegistrar with the given middleware.
func (r *Router) WithMiddleware(mw ...MiddlewareFunc) *RouteRegistrar {
	return &RouteRegistrar{
		router: r,
		attrs:  GroupAttributes{Middleware: mw},
	}
}

// WithoutMiddleware returns a RouteRegistrar that excludes the given middleware.
func (r *Router) WithoutMiddleware(mw ...MiddlewareFunc) *RouteRegistrar {
	return &RouteRegistrar{
		router: r,
		attrs:  GroupAttributes{WithoutMiddleware: mw},
	}
}

// WithPrefix returns a RouteRegistrar with the given prefix.
func (r *Router) WithPrefix(prefix string) *RouteRegistrar {
	return &RouteRegistrar{
		router: r,
		attrs:  GroupAttributes{Prefix: prefix},
	}
}

// WithDomain returns a RouteRegistrar with the given domain.
func (r *Router) WithDomain(domain string) *RouteRegistrar {
	return &RouteRegistrar{
		router: r,
		attrs:  GroupAttributes{Domain: domain},
	}
}

// WithName returns a RouteRegistrar with the given name prefix.
func (r *Router) WithName(name string) *RouteRegistrar {
	return &RouteRegistrar{
		router: r,
		attrs:  GroupAttributes{NamePrefix: name},
	}
}

// Middleware adds middleware to the registrar.
func (rr *RouteRegistrar) Middleware(mw ...MiddlewareFunc) *RouteRegistrar {
	rr.attrs.Middleware = append(rr.attrs.Middleware, mw...)

	return rr
}

// WithoutMiddleware adds middleware to the exclusion list.
func (rr *RouteRegistrar) WithoutMiddleware(mw ...MiddlewareFunc) *RouteRegistrar {
	rr.attrs.WithoutMiddleware = append(rr.attrs.WithoutMiddleware, mw...)

	return rr
}

// Prefix sets the URI prefix.
func (rr *RouteRegistrar) Prefix(prefix string) *RouteRegistrar {
	rr.attrs.Prefix = prefix

	return rr
}

// Domain sets the domain constraint.
func (rr *RouteRegistrar) Domain(domain string) *RouteRegistrar {
	rr.attrs.Domain = domain

	return rr
}

// Name sets the name prefix.
func (rr *RouteRegistrar) Name(name string) *RouteRegistrar {
	rr.attrs.NamePrefix = name

	return rr
}

// Where adds a parameter constraint.
func (rr *RouteRegistrar) Where(param, pattern string) *RouteRegistrar {
	if rr.attrs.Where == nil {
		rr.attrs.Where = make(map[string]string)
	}

	rr.attrs.Where[param] = pattern

	return rr
}

// Group creates a scoped sub-router with the accumulated attributes.
func (rr *RouteRegistrar) Group(fn func(*Router)) {
	rr.router.GroupWith(rr.attrs, fn)
}

// Get registers a GET route with the accumulated attributes.
func (rr *RouteRegistrar) Get(path string, handler HandlerFunc) *Route {
	var route *Route

	rr.router.GroupWith(rr.attrs, func(sub *Router) {
		route = sub.Get(path, handler)
	})

	return route
}

// Post registers a POST route with the accumulated attributes.
func (rr *RouteRegistrar) Post(path string, handler HandlerFunc) *Route {
	var route *Route

	rr.router.GroupWith(rr.attrs, func(sub *Router) {
		route = sub.Post(path, handler)
	})

	return route
}

// Put registers a PUT route with the accumulated attributes.
func (rr *RouteRegistrar) Put(path string, handler HandlerFunc) *Route {
	var route *Route

	rr.router.GroupWith(rr.attrs, func(sub *Router) {
		route = sub.Put(path, handler)
	})

	return route
}

// Delete registers a DELETE route with the accumulated attributes.
func (rr *RouteRegistrar) Delete(path string, handler HandlerFunc) *Route {
	var route *Route

	rr.router.GroupWith(rr.attrs, func(sub *Router) {
		route = sub.Delete(path, handler)
	})

	return route
}

// Patch registers a PATCH route with the accumulated attributes.
func (rr *RouteRegistrar) Patch(path string, handler HandlerFunc) *Route {
	var route *Route

	rr.router.GroupWith(rr.attrs, func(sub *Router) {
		route = sub.Patch(path, handler)
	})

	return route
}

// Options registers an OPTIONS route with the accumulated attributes.
func (rr *RouteRegistrar) Options(path string, handler HandlerFunc) *Route {
	var route *Route

	rr.router.GroupWith(rr.attrs, func(sub *Router) {
		route = sub.Options(path, handler)
	})

	return route
}

// Any registers a route for all methods with the accumulated attributes.
func (rr *RouteRegistrar) Any(path string, handler HandlerFunc) *Route {
	var route *Route

	rr.router.GroupWith(rr.attrs, func(sub *Router) {
		route = sub.Any(path, handler)
	})

	return route
}

// Match registers a route for specific methods with the accumulated attributes.
func (rr *RouteRegistrar) Match(methods []string, path string, handler HandlerFunc) *Route {
	var route *Route

	rr.router.GroupWith(rr.attrs, func(sub *Router) {
		route = sub.Match(methods, path, handler)
	})

	return route
}

// Resource registers resource routes with the accumulated attributes.
func (rr *RouteRegistrar) Resource(name string, handlers ResourceHandlers, opts ...ResourceOptions) {
	rr.router.GroupWith(rr.attrs, func(sub *Router) {
		reg := NewResourceRegistrar(sub)
		reg.Register(name, handlers, opts...)
	})
}

// APIResource registers API resource routes with the accumulated attributes.
func (rr *RouteRegistrar) APIResource(name string, handlers ResourceHandlers, opts ...ResourceOptions) {
	rr.router.GroupWith(rr.attrs, func(sub *Router) {
		reg := NewResourceRegistrar(sub)
		reg.APIResource(name, handlers, opts...)
	})
}
