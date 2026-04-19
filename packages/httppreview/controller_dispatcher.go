package httppreview

import (
	"fmt"
	"reflect"

	"github.com/bedrock/packages/routing"
	"github.com/bedrock/packages/routing/controllers"
)

// ControllerDispatcher resolves a controller method's dependencies (triggering
// validation) and then short-circuits with a [SuccessResponse] panic instead
// of executing the method.
//
// Mirrors Framework\Foundation\Routing\HTTPPreviewControllerDispatcher.
type ControllerDispatcher struct {
	routing.ResolvesRouteDependencies
	routing.FiltersControllerMiddleware
	container routing.DependencyContainer
}

// NewControllerDispatcher creates a httppreview controller dispatcher.
func NewControllerDispatcher(container routing.DependencyContainer) *ControllerDispatcher {
	d := &ControllerDispatcher{container: container}
	d.ResolvesRouteDependencies.Bind(container)

	return d
}

// Dispatch ensures the method exists on the controller, resolves its
// parameters (which may trigger form request validation), and then panics
// with [SuccessResponse] instead of invoking the method.
//
// The route parameter is typed as any to satisfy the
// [contracts.ControllerDispatcher] interface.
func (d *ControllerDispatcher) Dispatch(route any, controller any, method string) (any, error) {
	d.ensureMethodExists(controller, method)

	r, ok := route.(routeAccessor)

	if ok {
		d.ResolveClassMethodDependencies(
			r.ParametersWithoutNulls(),
			controller,
			method,
			r.ParameterNames(),
		)
	}

	panic(SuccessResponse{})
}

// GetMiddleware returns the controller's middleware filtered by only/except
// options for the given method.
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

		if routing.MethodExcludedByOptions(method, opts) {
			continue
		}

		out = append(out, e.Middleware)
	}

	return out
}

// ensureMethodExists panics if the controller does not have the named method.
// This mirrors HTTPPreviewControllerDispatcher::ensureMethodExists which
// throws a RuntimeException.
func (d *ControllerDispatcher) ensureMethodExists(controller any, method string) {
	rv := reflect.ValueOf(controller)
	rt := rv.Type()

	if _, ok := rt.MethodByName(method); !ok {
		panic(fmt.Sprintf(
			"httppreview: attempting to predict the outcome of the [%s::%s()] method but the method is not defined",
			rt.String(), method,
		))
	}
}
