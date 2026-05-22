package routing

import (
	"reflect"
	"strconv"
)

// same name. It performs reflection-based parameter resolution for both
// callable handlers (used by [CallableDispatcher]) and controller methods
// (used by [ControllerDispatcher]).
//
// The PHP version uses ReflectionParameter to inspect type hints, then
// either fetches an instance from the container, falls back to the route's
// parameter map (string keys → parameter values), or uses the parameter's
// default value. The Go version follows the same algorithm using
// [reflect.Type.In] and a [DependencyContainer] for type-keyed lookups.
type ResolvesRouteDependencies struct {
	container DependencyContainer
}

// DependencyContainer is the surface a container must expose so the resolver
// can look up typed parameters by their reflect.Type. The bedrock container
// has [DependencyContainer.Make] but lacks reflection-based binding; the
// concrete adapter in M11 will wrap container.Container.
type DependencyContainer interface {
	// MakeFor resolves an instance of the requested type. It is the Go
	// counterpart to PHP's $container->make($className), keyed by reflect.Type
	// rather than a string class name.
	MakeFor(t reflect.Type) (any, error)
}

// Bind wires the resolver to a container.

// ResolveMethodDependencies merges the route parameter map into the supplied
// callable's parameter list, filling any non-string-typed parameters via the
// container.
//
// Returns the final argument list as a slice of [reflect.Value] ready to feed
// to [reflect.Value.Call].

// String/numeric primitives consume route parameters in declaration order.

// Otherwise resolve from the container (may be nil for tests/M5).

// ResolveClassMethodDependencies is the controller-method counterpart of
// ResolveMethodDependencies. It locates the method on the supplied receiver
// via reflection and resolves its arguments the same way.

// reflect.Method.Type includes the receiver as In(0); skip it.

// MissingControllerMethodError is returned by
// [ResolveClassMethodDependencies] when the named method is absent from the
// supplied receiver.
type MissingControllerMethodError struct {
	Type   string
	Method string
}

func (r *ResolvesRouteDependencies) Bind(container DependencyContainer) {
	r.container = container
}

func (r *ResolvesRouteDependencies) ResolveMethodDependencies(
	parameters map[string]string,
	fnType reflect.Type,
	parameterNames []string,
) ([]reflect.Value, error) {
	in := make([]reflect.Value, fnType.NumIn())
	stringConsumed := 0

	for i := 0; i < fnType.NumIn(); i++ {
		paramType := fnType.In(i)

		if isPrimitive(paramType) {
			if stringConsumed < len(parameterNames) {
				name := parameterNames[stringConsumed]

				if val, ok := parameters[name]; ok {
					in[i] = convertString(val, paramType)
					stringConsumed++

					continue
				}
			}

			in[i] = reflect.Zero(paramType)
			stringConsumed++

			continue
		}

		if r.container != nil {
			if v, err := r.container.MakeFor(paramType); err == nil && v != nil {
				in[i] = reflect.ValueOf(v)

				continue
			}
		}

		in[i] = reflect.Zero(paramType)
	}

	return in, nil
}

func (r *ResolvesRouteDependencies) ResolveClassMethodDependencies(
	parameters map[string]string,
	receiver any,
	method string,
	parameterNames []string,
) ([]reflect.Value, *reflect.Method, reflect.Value, error) {
	rv := reflect.ValueOf(receiver)
	rvType := rv.Type()
	m, ok := rvType.MethodByName(method)

	if !ok {
		return nil, nil, rv, &MissingControllerMethodError{Type: rvType.String(), Method: method}
	}

	fnType := m.Type
	in := make([]reflect.Value, fnType.NumIn())
	in[0] = rv
	stringConsumed := 0

	for i := 1; i < fnType.NumIn(); i++ {
		pt := fnType.In(i)

		if isPrimitive(pt) {
			if stringConsumed < len(parameterNames) {
				name := parameterNames[stringConsumed]

				if val, ok := parameters[name]; ok {
					in[i] = convertString(val, pt)
					stringConsumed++

					continue
				}
			}

			in[i] = reflect.Zero(pt)
			stringConsumed++

			continue
		}

		if r.container != nil {
			if v, err := r.container.MakeFor(pt); err == nil && v != nil {
				in[i] = reflect.ValueOf(v)

				continue
			}
		}

		in[i] = reflect.Zero(pt)
	}

	return in, &m, rv, nil
}

func (e *MissingControllerMethodError) Error() string {
	return "controller " + e.Type + " has no method " + e.Method
}

// isPrimitive reports whether t is one of the kinds that a route parameter
// value (always a string at the wire level) can be converted into.
func isPrimitive(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.String, reflect.Int, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		return true
	}

	return false
}

// convertString parses a string parameter into the requested primitive type,
func convertString(s string, t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.String:
		return reflect.ValueOf(s).Convert(t)
	case reflect.Int, reflect.Int32, reflect.Int64:
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			return reflect.ValueOf(i).Convert(t)
		}
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		if u, err := strconv.ParseUint(s, 10, 64); err == nil {
			return reflect.ValueOf(u).Convert(t)
		}
	case reflect.Float32, reflect.Float64:
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return reflect.ValueOf(f).Convert(t)
		}
	case reflect.Bool:
		if b, err := strconv.ParseBool(s); err == nil {
			return reflect.ValueOf(b)
		}
	}

	return reflect.Zero(t)
}
