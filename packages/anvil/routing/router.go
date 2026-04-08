package routing

import (
	nethttp "net/http"
	"strings"
)

// Middleware decorates an HTTP handler.
type Middleware func(nethttp.Handler) nethttp.Handler

// HandlerFunc handles a routed request.
type HandlerFunc func(*Context) error

// Context holds request-scoped routing state.
type Context struct {
	Writer  nethttp.ResponseWriter
	Request *nethttp.Request
}

// Param returns a route parameter value by name.
func (c *Context) Param(name string) string {
	return c.Request.PathValue(name)
}

// Text writes a plain-text response.
func (c *Context) Text(status int, body string) error {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteHeader(status)
	_, _ = c.Writer.Write([]byte(body))

	return nil
}

// ResourceHandlers holds handler functions for standard CRUD routes.
type ResourceHandlers struct {
	Index   HandlerFunc // GET /resource
	Show    HandlerFunc // GET /resource/{id}
	Store   HandlerFunc // POST /resource
	Update  HandlerFunc // PUT /resource/{id}
	Destroy HandlerFunc // DELETE /resource/{id}
}

// Router routes HTTP requests to handlers.
type Router struct {
	mux        *nethttp.ServeMux
	prefix     string
	middleware []Middleware
	names      map[string]string
}

// New constructs a router.
func New() *Router {
	return &Router{
		mux:   nethttp.NewServeMux(),
		names: map[string]string{},
	}
}

// Handle registers a method-specific route.
func (r *Router) Handle(method string, route string, handler HandlerFunc) {
	fullRoute := r.prefix + route
	mw := append([]Middleware(nil), r.middleware...)

	// Go 1.22+ ServeMux supports "METHOD /path" patterns for method routing.
	pattern := fullRoute
	if method != "" {
		pattern = method + " " + fullRoute
	}

	r.mux.HandleFunc(pattern, func(writer nethttp.ResponseWriter, request *nethttp.Request) {
		var h nethttp.Handler = nethttp.HandlerFunc(func(w nethttp.ResponseWriter, req *nethttp.Request) {
			if err := handler(&Context{
				Writer:  w,
				Request: req,
			}); err != nil {
				nethttp.Error(w, err.Error(), nethttp.StatusInternalServerError)
			}
		})

		for i := len(mw) - 1; i >= 0; i-- {
			h = mw[i](h)
		}

		h.ServeHTTP(writer, request)
	})
}

// Get registers a GET route.
func (r *Router) Get(route string, handler HandlerFunc) {
	r.Handle(nethttp.MethodGet, route, handler)
}

// Post registers a POST route.
func (r *Router) Post(route string, handler HandlerFunc) {
	r.Handle(nethttp.MethodPost, route, handler)
}

// Put registers a PUT route.
func (r *Router) Put(route string, handler HandlerFunc) {
	r.Handle(nethttp.MethodPut, route, handler)
}

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
	r.Handle("OPTIONS", route, handler)
}

// Group registers routes under a shared prefix and middleware stack.
func (r *Router) Group(prefix string, middleware []Middleware, register func(*Router)) {
	child := &Router{
		mux:        r.mux,
		prefix:     r.prefix + prefix,
		middleware: append(append([]Middleware(nil), r.middleware...), middleware...),
		names:      r.names,
	}
	register(child)
}

// Resource registers standard CRUD routes for a resource name.
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

// Name associates a name with a route pattern for URL generation.
func (r *Router) Name(route string, name string) {
	r.names[name] = r.prefix + route
}

// Route generates a URL for a named route, substituting path parameters.
func (r *Router) Route(name string, params map[string]string) string {
	pattern, ok := r.names[name]
	if !ok {
		return ""
	}
	for k, v := range params {
		pattern = strings.ReplaceAll(pattern, "{"+k+"}", v)
	}
	return pattern
}

// ServeHTTP serves the configured routes.
func (r *Router) ServeHTTP(writer nethttp.ResponseWriter, request *nethttp.Request) {
	r.mux.ServeHTTP(writer, request)
}
