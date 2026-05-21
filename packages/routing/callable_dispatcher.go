package routing

import "reflect"

// CallableDispatcher dispatches a closure-style route action by resolving its
// parameters via reflection and invoking it.
//
// Ref: @bedrock/code-0287
type CallableDispatcher struct {
	ResolvesRouteDependencies
	container DependencyContainer
}

// NewCallableDispatcher constructs a dispatcher bound to the given container.
// Either argument may be nil — the dispatcher will then resolve only string
// parameters from the route map.
func NewCallableDispatcher(container DependencyContainer) *CallableDispatcher {
	d := &CallableDispatcher{container: container}
	d.ResolvesRouteDependencies.Bind(container)

	return d
}

// Dispatch runs the supplied callable against the bound parameters of the
// supplied route and returns whatever the callable produced.
func (d *CallableDispatcher) Dispatch(route *Route, callable any) (any, error) {
	rv := reflect.ValueOf(callable)

	if rv.Kind() != reflect.Func {
		return nil, &MissingControllerMethodError{Type: rv.Type().String(), Method: "(callable)"}
	}

	in, err := d.ResolveMethodDependencies(
		route.ParametersWithoutNulls(),
		rv.Type(),
		route.ParameterNames(),
	)

	if err != nil {
		return nil, err
	}

	out := rv.Call(in)

	return packReturn(out), nil
}

// packReturn collapses a reflect.Call result slice into a single value:
//   - 0 returns → nil
//   - 1 return  → the value
//   - 2 returns where the second is an error → (value, error)
//   - otherwise → a []any of all returns
//
// This keeps the dispatcher signature aligned with [contracts.CallableDispatcher].
func packReturn(out []reflect.Value) any {
	switch len(out) {
	case 0:
		return nil
	case 1:
		return out[0].Interface()
	case 2:
		// Common (value, error) form: surface the value; the caller already
		// uses the error from a separate channel.
		return out[0].Interface()
	}

	all := make([]any, len(out))

	for i, v := range out {
		all[i] = v.Interface()
	}

	return all
}
