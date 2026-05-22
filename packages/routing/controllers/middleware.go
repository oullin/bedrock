// Ref: @bedrock/code-0298
// In the upstream framework 11+ this package supplies the [HasMiddleware] interface that
// controllers implement to declare their middleware (the modern replacement
// for PHP 8 attribute-based declarations). Both forms remain in 13.x, and
// this Go port treats the interface form as the canonical entry point.
package controllers

// Middleware is a value object describing a single piece of controller
// middleware along with the only/except filters that scope it to specific
// methods.
//
// Ref: @bedrock/code-0300
type Middleware struct {
	Middleware any // string class name, slice of class names, or func value
	Only       []string
	Except     []string
}

// NewMiddleware constructs a Middleware value with no method filter.
func NewMiddleware(m any) Middleware { return Middleware{Middleware: m} }

// WithOnly returns a copy that applies only to the named methods.
func (m Middleware) WithOnly(methods ...string) Middleware {
	m.Only = append([]string(nil), methods...)

	return m
}

// WithExcept returns a copy that excludes the named methods.
func (m Middleware) WithExcept(methods ...string) Middleware {
	m.Except = append([]string(nil), methods...)

	return m
}
