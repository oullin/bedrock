package routing

import (
	nethttp "net/http"
	"strings"
)

// MiddlewareFunc wraps an HTTP handler with additional behavior.
type MiddlewareFunc func(nethttp.Handler) nethttp.Handler

// Router routes HTTP requests to handlers. It wraps Go 1.22+ net/http.ServeMux
// and supports method-specific patterns, middleware composition, prefix groups,
// resource routing, and a named-route registry for URL generation.
type Router struct {
	mux        *nethttp.ServeMux
	prefix     string
	middleware []MiddlewareFunc
	registry   *Registry
}

// New constructs a Router. Pass a Registry to share named routes across
// sub-routers; pass nil to create a fresh registry automatically.
func New(registry *Registry) *Router {
	if registry == nil {
		registry = NewRegistry()
	}

	return &Router{
		mux:      nethttp.NewServeMux(),
		registry: registry,
	}
}

// Use appends middleware to the router. Middleware added here applies to every
// handler registered on this router (but not to sub-routers created before
// this call).
func (r *Router) Use(mw ...MiddlewareFunc) *Router {
	r.middleware = append(r.middleware, mw...)

	return r
}

// Handle registers handler for the given HTTP method and route path.
// The pattern follows Go 1.22+ ServeMux syntax; {param} segments are
// supported for path parameter extraction via Context.Param.
func (r *Router) Handle(method, route string, handler HandlerFunc) {
	fullRoute := r.prefix + route
	mw := append([]MiddlewareFunc(nil), r.middleware...)

	pattern := fullRoute

	if method != "" {
		pattern = strings.ToUpper(method) + " " + fullRoute
	}

	r.mux.HandleFunc(pattern, func(w nethttp.ResponseWriter, req *nethttp.Request) {
		var h nethttp.Handler = nethttp.HandlerFunc(func(ww nethttp.ResponseWriter, rr *nethttp.Request) {
			if err := handler(&Context{Writer: ww, Request: rr}); err != nil {
				nethttp.Error(ww, err.Error(), nethttp.StatusInternalServerError)
			}
		})

		for i := len(mw) - 1; i >= 0; i-- {
			h = mw[i](h)
		}

		h.ServeHTTP(w, req)
	})
}

// Get registers a GET route.
func (r *Router) Get(route string, handler HandlerFunc) { r.Handle(nethttp.MethodGet, route, handler) }

// Post registers a POST route.
func (r *Router) Post(route string, handler HandlerFunc) {
	r.Handle(nethttp.MethodPost, route, handler)
}

// Put registers a PUT route.
func (r *Router) Put(route string, handler HandlerFunc) { r.Handle(nethttp.MethodPut, route, handler) }

// Delete registers a DELETE route.
func (r *Router) Delete(route string, handler HandlerFunc) {
	r.Handle(nethttp.MethodDelete, route, handler)
}

// Patch registers a PATCH route.
func (r *Router) Patch(route string, handler HandlerFunc) {
	r.Handle(nethttp.MethodPatch, route, handler)
}

// Options registers an OPTIONS route.
func (r *Router) Options(route string, handler HandlerFunc) {
	r.Handle(nethttp.MethodOptions, route, handler)
}

// Group registers a set of routes under a shared path prefix and middleware
// stack. The child router inherits the parent's middleware.
func (r *Router) Group(prefix string, middleware []MiddlewareFunc, register func(*Router)) {
	child := &Router{
		mux:        r.mux,
		prefix:     r.prefix + prefix,
		middleware: append(append([]MiddlewareFunc(nil), r.middleware...), middleware...),
		registry:   r.registry,
	}

	register(child)
}

// Resource registers standard CRUD routes for a resource.
// The resource name must not have leading/trailing slashes.
func (r *Router) Resource(name string, handlers ResourceHandlers) {
	base := "/" + strings.Trim(name, "/")
	item := base + "/{id}"

	if handlers.Index != nil {
		r.Get(base, handlers.Index)
	}

	if handlers.Show != nil {
		r.Get(item, handlers.Show)
	}

	if handlers.Store != nil {
		r.Post(base, handlers.Store)
	}

	if handlers.Update != nil {
		r.Put(item, handlers.Update)
	}

	if handlers.Destroy != nil {
		r.Delete(item, handlers.Destroy)
	}
}

// Name registers a named route for URL generation. The route pattern must
// be the path portion only (without the method prefix).
func (r *Router) Name(route, name string) {
	r.registry.Add(name, "", r.prefix+route)
}

// Route generates a URL for a named route, substituting {param} placeholders
// with the provided values. Returns empty string if the name is unknown.
func (r *Router) Route(name string, params map[string]string) string {
	return r.registry.URL(name, params)
}

// Registry returns the underlying named route registry.
func (r *Router) Registry() *Registry {
	return r.registry
}

// ServeHTTP satisfies http.Handler.
func (r *Router) ServeHTTP(w nethttp.ResponseWriter, req *nethttp.Request) {
	r.mux.ServeHTTP(w, req)
}
