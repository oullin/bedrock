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
	collection *RouteCollection
	patterns   map[string]string
	domain     string
	namePrefix string
}

// New constructs a Router. Pass a Registry to share named routes across
// sub-routers; pass nil to create a fresh registry automatically.
func New(registry *Registry) *Router {
	if registry == nil {
		registry = NewRegistry()
	}

	return &Router{
		mux:        nethttp.NewServeMux(),
		registry:   registry,
		collection: NewRouteCollection(),
		patterns:   make(map[string]string),
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
func (r *Router) Handle(method, path string, handler HandlerFunc) *Route {
	fullPath := r.prefix + path
	mw := append([]MiddlewareFunc(nil), r.middleware...)

	methods := []string{strings.ToUpper(method)}

	if method == "" {
		methods = nil
	}

	route := newRoute(methods, fullPath, handler)
	route.prefix = r.prefix
	route.domain = r.domain
	route.namePrefix = r.namePrefix
	route.middleware = append(route.middleware, mw...)

	for param, pattern := range r.patterns {
		route.Where(param, pattern)
	}

	r.collection.Add(route)

	pattern := fullPath

	if method != "" {
		pattern = strings.ToUpper(method) + " " + fullPath
	}

	r.mux.HandleFunc(pattern, func(w nethttp.ResponseWriter, req *nethttp.Request) {
		params := extractParams(fullPath, req)

		if !route.MatchesConstraints(params) {
			nethttp.Error(w, "Not Found", nethttp.StatusNotFound)

			return
		}

		var h nethttp.Handler = nethttp.HandlerFunc(func(ww nethttp.ResponseWriter, rr *nethttp.Request) {
			if err := handler(&Context{Writer: ww, Request: rr}); err != nil {
				nethttp.Error(ww, err.Error(), nethttp.StatusInternalServerError)
			}
		})

		allMw := append(mw, route.middleware[len(mw):]...)

		for i := len(allMw) - 1; i >= 0; i-- {
			h = allMw[i](h)
		}

		h.ServeHTTP(w, req)
	})

	return route
}

// Get registers a GET route.
func (r *Router) Get(path string, handler HandlerFunc) *Route {
	return r.Handle(nethttp.MethodGet, path, handler)
}

// Post registers a POST route.
func (r *Router) Post(path string, handler HandlerFunc) *Route {
	return r.Handle(nethttp.MethodPost, path, handler)
}

// Put registers a PUT route.
func (r *Router) Put(path string, handler HandlerFunc) *Route {
	return r.Handle(nethttp.MethodPut, path, handler)
}

// Delete registers a DELETE route.
func (r *Router) Delete(path string, handler HandlerFunc) *Route {
	return r.Handle(nethttp.MethodDelete, path, handler)
}

// Patch registers a PATCH route.
func (r *Router) Patch(path string, handler HandlerFunc) *Route {
	return r.Handle(nethttp.MethodPatch, path, handler)
}

// Options registers an OPTIONS route.
func (r *Router) Options(path string, handler HandlerFunc) *Route {
	return r.Handle(nethttp.MethodOptions, path, handler)
}

// Any registers a route that responds to all HTTP methods.
func (r *Router) Any(path string, handler HandlerFunc) *Route {
	return r.Handle("", path, handler)
}

// Match registers a route that responds to the given HTTP methods.
func (r *Router) Match(methods []string, path string, handler HandlerFunc) *Route {
	var last *Route

	for _, method := range methods {
		last = r.Handle(method, path, handler)
	}

	return last
}

// Redirect registers a route that redirects from one URI to another.
func (r *Router) Redirect(from, to string, status ...int) {
	code := nethttp.StatusFound

	if len(status) > 0 {
		code = status[0]
	}

	r.Any(from, func(ctx *Context) error {
		nethttp.Redirect(ctx.Writer, ctx.Request, to, code)

		return nil
	})
}

// PermanentRedirect registers a 301 redirect from one URI to another.
func (r *Router) PermanentRedirect(from, to string) {
	r.Redirect(from, to, nethttp.StatusMovedPermanently)
}

// Fallback registers a handler for requests that don't match any other route.
func (r *Router) Fallback(handler HandlerFunc) *Route {
	route := r.Any("/", handler)
	route.SetFallback(true)

	return route
}

// Pattern registers a global parameter constraint applied to all routes
// registered after this call.
func (r *Router) Pattern(param, pattern string) *Router {
	r.patterns[param] = pattern

	return r
}

// Group registers a set of routes under a shared path prefix and middleware
// stack. The child router inherits the parent's middleware.
func (r *Router) Group(prefix string, middleware []MiddlewareFunc, register func(*Router)) {
	child := &Router{
		mux:        r.mux,
		prefix:     r.prefix + prefix,
		middleware: append(append([]MiddlewareFunc(nil), r.middleware...), middleware...),
		registry:   r.registry,
		collection: r.collection,
		patterns:   copyPatterns(r.patterns),
		domain:     r.domain,
		namePrefix: r.namePrefix,
	}

	register(child)
}

// GroupWith creates a sub-router with full attribute control. It merges the
// parent attributes with the provided GroupAttributes following Laravel
// conventions.
func (r *Router) GroupWith(attrs GroupAttributes, register func(*Router)) {
	parent := GroupAttributes{
		Prefix:     r.prefix,
		Domain:     r.domain,
		NamePrefix: r.namePrefix,
		Middleware: r.middleware,
		Where:      r.patterns,
	}

	merged := MergeGroup(parent, attrs)
	patterns := make(map[string]string, len(merged.Where))

	for k, v := range merged.Where {
		patterns[k] = v
	}

	child := &Router{
		mux:        r.mux,
		prefix:     merged.Prefix,
		middleware: merged.Middleware,
		registry:   r.registry,
		collection: r.collection,
		patterns:   patterns,
		domain:     merged.Domain,
		namePrefix: merged.NamePrefix,
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

	if handlers.Create != nil {
		r.Get(base+"/create", handlers.Create)
	}

	if handlers.Show != nil {
		r.Get(item, handlers.Show)
	}

	if handlers.Store != nil {
		r.Post(base, handlers.Store)
	}

	if handlers.Edit != nil {
		r.Get(item+"/edit", handlers.Edit)
	}

	if handlers.Update != nil {
		r.Put(item, handlers.Update)
	}

	if handlers.Destroy != nil {
		r.Delete(item, handlers.Destroy)
	}
}

// NameRoute registers a named route for URL generation. The route pattern must
// be the path portion only (without the method prefix).
func (r *Router) NameRoute(path, name string) {
	r.registry.Add(name, "", r.prefix+path)
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

// Collection returns the underlying route collection.
func (r *Router) Collection() *RouteCollection {
	return r.collection
}

// SetDomain sets the domain constraint for this router scope.
func (r *Router) SetDomain(domain string) *Router {
	r.domain = domain

	return r
}

// SetNamePrefix sets the name prefix for routes registered on this router.
func (r *Router) SetNamePrefix(prefix string) *Router {
	r.namePrefix = prefix

	return r
}

// ServeHTTP satisfies http.Handler.
func (r *Router) ServeHTTP(w nethttp.ResponseWriter, req *nethttp.Request) {
	r.mux.ServeHTTP(w, req)
}

func copyPatterns(src map[string]string) map[string]string {
	dst := make(map[string]string, len(src))

	for k, v := range src {
		dst[k] = v
	}

	return dst
}
