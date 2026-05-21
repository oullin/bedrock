package routing

import (
	"fmt"
	"reflect"

	"github.com/bedrock/packages/routing/contracts"
	"github.com/bedrock/packages/routing/exceptions"
)

// ImplicitRouteBinding mirrors @bedrock\Routing\ImplicitRouteBinding.
//
// In PHP it walks the route's signature parameters, picks out those whose
// type implements [contracts.UrlRoutable], and replaces the matching path
// parameter with the resolved instance (looked up via the parameter's class).
// The Go port does the same against a [DependencyContainer] keyed by
// reflect.Type.
type ImplicitRouteBinding struct{}

// ResolveForRoute resolves implicit bindings on route in place.
//
// container is used to construct fresh instances of each typed parameter so
// the resolver can call ResolveRouteBinding on them.
func (ImplicitRouteBinding) ResolveForRoute(container DependencyContainer, route *Route) error {
	parameters := route.Parameters

	if parameters == nil {
		return nil
	}

	for _, p := range route.SignatureParameters(map[string]any{
		"subClass": reflect.TypeOf((*contracts.UrlRoutable)(nil)).Elem(),
	}) {
		paramName := matchParamName(p, route.ParameterNames(), parameters)

		if paramName == "" {
			continue
		}

		raw := parameters[paramName]

		if raw == "" {
			continue
		}

		instanceAny, err := container.MakeFor(p.Type)

		if err != nil || instanceAny == nil {
			return fmt.Errorf("implicit binding: cannot construct %v: %w", p.Type, err)
		}

		instance, ok := instanceAny.(contracts.UrlRoutable)

		if !ok {
			return fmt.Errorf("type %v does not implement UrlRoutable", p.Type)
		}

		field := route.BindingFieldFor(paramName)

		var (
			model any
			rerr  error
		)
		parent := route.ParentOfParameter(paramName)

		if parent != "" && !route.PreventsScopedBindings() &&
			(route.EnforcesScopedBindings() || field != "") {
			if parentVal, ok := route.Parameters[parent]; ok {
				if parentRoutable, ok := any(parentVal).(contracts.UrlRoutable); ok {
					model, rerr = parentRoutable.ResolveChildRouteBinding(paramName, raw, field)
				}
			}
		}

		if model == nil && rerr == nil {
			model, rerr = instance.ResolveRouteBinding(raw, field)
		}

		if rerr != nil {
			return rerr
		}

		if model == nil {
			return &ModelNotFoundError{Model: p.Type.String()}
		}
		// Store as a string for parameter-map compatibility — full model
		// substitution happens via the ResolvedBindings map in M11 wiring.
		route.SetParameter(paramName, fmt.Sprintf("%v", model))
		// Best-effort: also store the typed model in a side-map for richer
		// consumers via the bindings cache.
		route.storeBoundModel(paramName, model)
	}
	// Backed enum resolution.
	for _, p := range route.SignatureParameters(map[string]any{"backedEnum": true}) {
		paramName := matchParamName(p, route.ParameterNames(), parameters)

		if paramName == "" {
			continue
		}

		raw := parameters[paramName]

		if raw == "" {
			continue
		}
		// Construct a zero value of the enum type and check tryFrom semantics
		// via the BackedEnum interface.
		zero := reflect.New(p.Type).Elem()

		if be, ok := zero.Interface().(BackedEnum); ok {
			if be.BackingValue() == raw {
				continue
			}
		}

		return &exceptions.BackedEnumCaseNotFoundException{
			Enum: p.Type.String(),
			Case: raw,
		}
	}

	return nil
}

func matchParamName(p SignatureParameter, parameterNames []string, parameters map[string]string) string {
	// In PHP this looks up the parameter by snake-casing the variable name;
	// the Go form has no parameter names from reflection so we match by
	// position. Specifically: the i-th URL parameter binds to the i-th
	// non-primitive signature parameter.
	if p.Index >= 0 && p.Index < len(parameterNames) {
		name := parameterNames[p.Index]

		if _, ok := parameters[name]; ok {
			return name
		}
	}

	return ""
}
