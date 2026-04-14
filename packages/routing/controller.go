package routing

import "github.com/bedrock/packages/routing/controllers"

// Controller is the embeddable base class for routing controllers.
//
// In Go the recommended pattern is to compose a [Controller] into your own
// struct and override [Controller.Middleware] (via the [HasMiddleware]
// interface) to declare per-method middleware. The base type retains the
// PHP class's public surface (Middleware, GetMiddleware, CallAction) so
// dispatchers can probe for it uniformly.
//
// Mirrors Framework\Routing\Controller.
type Controller struct {
	middleware []controllerMiddlewareEntry
}

type controllerMiddlewareEntry struct {
	Middleware any
	Options    map[string]any
}

// MiddlewareOptions is the chainable filter helper returned by Use.
//
// Mirrors Framework\Routing\ControllerMiddlewareOptions.
type MiddlewareOptions struct{ options map[string]any }

// Only constrains the middleware to the named methods.
func (o *MiddlewareOptions) Only(methods ...string) *MiddlewareOptions {
	o.options["only"] = methods

	return o
}

// Except excludes the middleware from the named methods.
func (o *MiddlewareOptions) Except(methods ...string) *MiddlewareOptions {
	o.options["except"] = methods

	return o
}

// Use registers middleware on the controller and returns a chainable options
// object (Only/Except).
//
// Mirrors Controller::middleware. Renamed to Use in Go because Middleware is
// already the [HasMiddleware] method name.
func (c *Controller) Use(middleware any) *MiddlewareOptions {
	options := map[string]any{}
	c.middleware = append(c.middleware, controllerMiddlewareEntry{
		Middleware: middleware,
		Options:    options,
	})

	return &MiddlewareOptions{options: options}
}

// GetMiddleware returns the middleware registered via [Controller.Use].
func (c *Controller) GetMiddleware() []controllerMiddlewareEntry { return c.middleware }

// Middleware satisfies [controllers.HasMiddleware] for controllers that only
// register middleware imperatively via Use. Override this on your own
// embedding type to return declarative entries.
func (c *Controller) Middleware() []controllers.Middleware {
	out := make([]controllers.Middleware, 0, len(c.middleware))

	for _, e := range c.middleware {
		m := controllers.Middleware{Middleware: e.Middleware}

		if only, ok := e.Options["only"].([]string); ok {
			m.Only = only
		}

		if except, ok := e.Options["except"].([]string); ok {
			m.Except = except
		}

		out = append(out, m)
	}

	return out
}
