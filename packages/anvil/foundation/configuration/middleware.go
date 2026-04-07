package configuration

import (
	nethttp "net/http"

	apphttp "github.com/bedrock/packages/anvil/http"
)

// Middleware stores a stack of HTTP middleware.
type Middleware struct {
	stack []apphttp.Middleware
}

// NewMiddleware creates an empty middleware configuration.
func NewMiddleware() *Middleware {
	return &Middleware{stack: []apphttp.Middleware{}}
}

// Use appends middleware to the stack.
func (m *Middleware) Use(middleware apphttp.Middleware) {
	if middleware != nil {
		m.stack = append(m.stack, middleware)
	}
}

// Wrap applies middleware to a handler in registration order.
func (m *Middleware) Wrap(handler nethttp.Handler) nethttp.Handler {
	wrapped := handler

	for index := len(m.stack) - 1; index >= 0; index-- {
		wrapped = m.stack[index](wrapped)
	}

	return wrapped
}
