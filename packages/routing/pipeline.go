package routing

import nethttp "net/http"

// Pipeline chains middleware functions and a final handler into a single
// http.Handler. Middleware are executed in the order they were added.
type Pipeline struct {
	middleware []MiddlewareFunc
}

// NewPipeline creates a pipeline with the given middleware.
func NewPipeline(middleware ...MiddlewareFunc) *Pipeline {
	return &Pipeline{middleware: middleware}
}

// Then wraps handler with the pipeline's middleware and returns the result.
func (p *Pipeline) Then(handler nethttp.Handler) nethttp.Handler {
	for i := len(p.middleware) - 1; i >= 0; i-- {
		handler = p.middleware[i](handler)
	}

	return handler
}

// ThenFunc wraps a handler function with the pipeline's middleware.
func (p *Pipeline) ThenFunc(fn nethttp.HandlerFunc) nethttp.Handler {
	return p.Then(fn)
}
