package routing

import (
	"fmt"
	"strings"
)

// BindingContainer is the minimum container surface RouteBinding needs to
// resolve a class-based binding. The bedrock packages/container Container
// satisfies this via its Make method.
type BindingContainer interface {
	Make(abstract string) (any, error)
}

// BindingResolver is a callback that converts a raw URL value (and the
type BindingResolver func(value string, route any) (any, error)

// ModelInstance is the surface a user-defined model type must expose so
// [ForModel] can resolve it.
// Model that route-model-binding actually depends on.
type ModelInstance interface {
	ResolveRouteBinding(value, field string) (any, error)
	ResolveSoftDeletableRouteBinding(value, field string) (any, error)
	IsSoftDeletable() bool
}

// RouteBinding is a parity-named wrapper for the static helpers below. PHP
// callers reach the helpers through `RouteBinding::forCallback($c, $binder)`;
// Go callers can use either the package functions directly or the methods on
// this empty struct, both produce the same result.
//
// Ref: @bedrock/code-0334
type RouteBinding struct{}

// ForCallback returns the binding resolver for the given binder. If binder is
// a [BindingResolver] it is returned unchanged. If it is a string it is
// treated as a "Class@method" reference and resolved through the container.
//
// Ref: @bedrock/code-0334

// ForCallback is the method form of the package-level helper.

// ForModel returns a binding resolver that fetches a model instance via the
// container and asks it to resolve the URL value. If the model returns nil and
// fallback is supplied, fallback is invoked; otherwise an error is returned.
//
// Ref: @bedrock/code-0334
// ModelNotFoundException; the Go version returns a typed error so callers can
// test for it with errors.Is.

type ModelNotFoundError struct{ Model string }

func ForCallback(container BindingContainer, binder any) BindingResolver {
	if b, ok := binder.(BindingResolver); ok {
		return b
	}

	if s, ok := binder.(string); ok {
		return createClassBinding(container, s)
	}

	return func(_ string, _ any) (any, error) {
		return nil, fmt.Errorf("invalid binder: %T", binder)
	}
}

func (RouteBinding) ForCallback(container BindingContainer, binder any) BindingResolver {
	return ForCallback(container, binder)
}

func createClassBinding(container BindingContainer, binding string) BindingResolver {
	return func(value string, route any) (any, error) {
		class, method := parseCallback(binding, "bind")
		instance, err := container.Make(class)

		if err != nil {
			return nil, err
		}

		bound, ok := instance.(interface {
			Bind(value string, route any) (any, error)
		})

		if ok && method == "bind" {
			return bound.Bind(value, route)
		}

		dispatcher, ok := instance.(interface {
			Call(method string, value string, route any) (any, error)
		})

		if !ok {
			return nil, fmt.Errorf("class binding %q does not expose %s()", class, method)
		}

		return dispatcher.Call(method, value, route)
	}
}

func ForModel(container BindingContainer, class string, fallback BindingResolver) BindingResolver {
	return func(value string, route any) (any, error) {
		if value == "" {
			return nil, nil
		}

		raw, err := container.Make(class)

		if err != nil {
			return nil, err
		}

		instance, ok := raw.(ModelInstance)

		if !ok {
			return nil, fmt.Errorf("class %q is not a ModelInstance", class)
		}

		method := instance.ResolveRouteBinding

		if r, ok := route.(interface{ AllowsTrashedBindings() bool }); ok && r.AllowsTrashedBindings() && instance.IsSoftDeletable() {
			method = instance.ResolveSoftDeletableRouteBinding
		}

		model, err := method(value, "")

		if err != nil {
			return nil, err
		}

		if model != nil {
			return model, nil
		}

		if fallback != nil {
			return fallback(value, route)
		}

		return nil, &ModelNotFoundError{Model: class}
	}
}

func (e *ModelNotFoundError) Error() string { return "no query results for model [" + e.Model + "]" }

// parseCallback splits a "Class@method" string into its components, falling
// back to the supplied default method when no "@" is present.
func parseCallback(callback, defaultMethod string) (class, method string) {
	if idx := strings.Index(callback, "@"); idx >= 0 {
		return callback[:idx], callback[idx+1:]
	}

	return callback, defaultMethod
}
