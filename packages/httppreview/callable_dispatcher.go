package httppreview

import (
	"reflect"

	"github.com/bedrock/packages/routing"
)

// CallableDispatcher resolves a callable's dependencies (triggering validation)
// and then short-circuits with a [SuccessResponse] panic instead of executing
// the callable.
//
// Mirrors Framework\Foundation\Routing\HTTPPreviewCallableDispatcher.
type CallableDispatcher struct {
	routing.ResolvesRouteDependencies
	container routing.DependencyContainer
}

// NewCallableDispatcher creates a httppreview callable dispatcher.

// Dispatch resolves the callable's parameters (which may trigger form request
// validation via the container) and then panics with [SuccessResponse] instead
// of invoking the callable.
//
// The route parameter is typed as any to satisfy the
// [contracts.CallableDispatcher] interface.

// routeAccessor is a minimal interface for reading route parameters. It avoids
// a hard dependency on *routing.Route which would create an import cycle
// through the contracts package.
type routeAccessor interface {
	ParametersWithoutNulls() map[string]string
	ParameterNames() []string
}

func NewCallableDispatcher(container routing.DependencyContainer) *CallableDispatcher {
	d := &CallableDispatcher{container: container}
	d.ResolvesRouteDependencies.Bind(container)

	return d
}

func (d *CallableDispatcher) Dispatch(route any, callable any) (any, error) {
	r, ok := route.(routeAccessor)

	if ok {
		rv := reflect.ValueOf(callable)

		if rv.Kind() == reflect.Func {
			d.ResolveMethodDependencies(
				r.ParametersWithoutNulls(),
				rv.Type(),
				r.ParameterNames(),
			)
		}
	}

	panic(SuccessResponse{})
}
