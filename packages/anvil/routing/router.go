package routing

import (
	"fmt"
	nethttp "net/http"

	viewpkg "github.com/bedrock/packages/anvil/view"
)

// HandlerFunc handles a routed request.
type HandlerFunc func(*Context) error

// Context holds request-scoped routing state.
type Context struct {
	Writer   nethttp.ResponseWriter
	Request  *nethttp.Request
	Renderer *viewpkg.Renderer
}

// Render renders a named view with a 200 status.
func (c *Context) Render(name string, data any) error {
	if c.Renderer == nil {
		return fmt.Errorf("routing: renderer is not configured")
	}

	return c.Renderer.WriteHTML(c.Writer, nethttp.StatusOK, name, data)
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

// Router routes HTTP requests to handlers.
type Router struct {
	mux      *nethttp.ServeMux
	renderer *viewpkg.Renderer
}

// New constructs a router bound to a view renderer.
func New(renderer *viewpkg.Renderer) *Router {
	return &Router{
		mux:      nethttp.NewServeMux(),
		renderer: renderer,
	}
}

// Handle registers a method-specific route.
func (r *Router) Handle(method string, route string, handler HandlerFunc) {
	r.mux.HandleFunc(route, func(writer nethttp.ResponseWriter, request *nethttp.Request) {
		if method != "" && request.Method != method {
			writer.Header().Set("Allow", method)
			writer.WriteHeader(nethttp.StatusMethodNotAllowed)
			return
		}

		if err := handler(&Context{
			Writer:   writer,
			Request:  request,
			Renderer: r.renderer,
		}); err != nil {
			nethttp.Error(writer, err.Error(), nethttp.StatusInternalServerError)
		}
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

// ServeHTTP serves the configured routes.
func (r *Router) ServeHTTP(writer nethttp.ResponseWriter, request *nethttp.Request) {
	r.mux.ServeHTTP(writer, request)
}
