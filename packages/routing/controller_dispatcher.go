package routing

import (
	"github.com/bedrock/packages/routing/controllers"
)

// ControllerDispatcher resolves and invokes a controller method against a
// matched route, applying reflection-based parameter resolution.
//
// Mirrors Illuminate\Routing\ControllerDispatcher.
type ControllerDispatcher struct {
	ResolvesRouteDependencies
	FiltersControllerMiddleware
	container DependencyContainer
}

// NewControllerDispatcher wires a dispatcher to a container.
func NewControllerDispatcher(container DependencyContainer) *ControllerDispatcher {
	d := &ControllerDispatcher{container: container}
	d.ResolvesRouteDependencies.Bind(container)

	return d
}

// Dispatch invokes method on controller with the parameters resolved from the
// route. controller must be a value or pointer whose type exposes the method.
func (d *ControllerDispatcher) Dispatch(route *Route, controller any, method string) (any, error) {
	in, m, _, err := d.ResolveClassMethodDependencies(
		route.ParametersWithoutNulls(),
		controller,
		method,
		route.ParameterNames(),
	)

	if err != nil {
		return nil, err
	}

	out := m.Func.Call(in)

	return packReturn(out), nil
}

// GetMiddleware returns the controller's middleware filtered by the
// only/except options for method.
//
// If the controller satisfies [controllers.HasMiddleware] its declared list
// is returned; otherwise the controller's [Controller.GetMiddleware] entries
// (if any) are flattened.
func (d *ControllerDispatcher) GetMiddleware(controller any, method string) []any {
	var entries []controllers.Middleware

	if hm, ok := controller.(controllers.HasMiddleware); ok {
		entries = hm.Middleware()
	}

	out := make([]any, 0, len(entries))

	for _, e := range entries {
		opts := map[string]any{}

		if e.Only != nil {
			opts["only"] = e.Only
		}

		if e.Except != nil {
			opts["except"] = e.Except
		}

		if MethodExcludedByOptions(method, opts) {
			continue
		}

		out = append(out, e.Middleware)
	}

	return out
}
