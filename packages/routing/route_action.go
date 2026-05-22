package routing

import (
	"fmt"
	"reflect"
	"strings"
)

// Action is the parsed form of a route action.
//
// In the upstream framework this is a free-form associative array; in Go it is a struct so
// callers get type safety without losing parity. The most-used keys map to
// fields directly:
//
//   - Uses: the actual handler. May be a Go func value, a "Controller@method"
//     string, an invokable controller type, or nil for fluent registration.
//   - Controller: the canonical "Type@Method" form when known.
//     of the same name.
//   - Domain, Prefix, As, Namespace: group-style attributes that may be
//     attached to the action.
//   - Extras: a catch-all map for any non-standard keys passed through by
//     callers that want to round-trip arbitrary action metadata.
//
// Ref: @bedrock/code-0333
type Action struct {
	Uses       any
	Controller string
	Middleware []any
	Where      map[string]string
	Defaults   map[string]any
	Domain     string
	Prefix     string
	As         string
	Namespace  string
	Extras     map[string]any
}

// ParseAction normalizes a user-supplied action into an [*Action].
//
// The accepted shapes mirror the upstream RouteAction::parse:
//   - nil: produces a placeholder action whose Uses returns an error when
//     invoked, equivalent to the upstream missingAction closure.
//   - a Go func: stored in Uses as-is.
//   - a string of the form "Controller@method": stored in Uses and Controller.
//   - a struct/pointer with an Invoke method: stored in Uses; Controller is
//     populated using the Go type name + "@__invoke".
//   - a map[string]any: passed through, with the "uses" key resolved as above
//     and other keys distributed to the matching fields or Extras.
//
// Ref: @bedrock/code-0333
func ParseAction(uri string, action any) (*Action, error) {
	if action == nil {
		return missingAction(uri), nil
	}

	switch a := action.(type) {
	case func():
		return &Action{Uses: a}, nil
	case map[string]any:
		return parseMapAction(uri, a)
	case string:
		return parseStringAction(a)
	}

	rv := reflect.ValueOf(action)

	if rv.Kind() == reflect.Func {
		return &Action{Uses: action}, nil
	}

	if isInvokable(rv) {
		t := rv.Type()

		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		return &Action{Uses: action, Controller: t.String() + "@Invoke"}, nil
	}

	return nil, fmt.Errorf("invalid route action: %T", action)
}

func missingAction(uri string) *Action {
	return &Action{
		Uses: func() error { return fmt.Errorf("route for [%s] has no action", uri) },
	}
}

func parseMapAction(uri string, m map[string]any) (*Action, error) {
	a := &Action{Extras: map[string]any{}}

	if v, ok := m["uses"]; ok {
		sub, err := ParseAction(uri, v)

		if err != nil {
			return nil, err
		}

		a.Uses = sub.Uses
		a.Controller = sub.Controller
	}

	for k, v := range m {
		switch k {
		case "uses":
			continue
		case "controller":
			if s, ok := v.(string); ok {
				a.Controller = s
			}
		case "middleware":
			if mw, ok := v.([]any); ok {
				a.Middleware = mw
			}
		case "where":
			if w, ok := v.(map[string]string); ok {
				a.Where = w
			}
		case "defaults":
			if d, ok := v.(map[string]any); ok {
				a.Defaults = d
			}
		case "domain":
			if s, ok := v.(string); ok {
				a.Domain = s
			}
		case "prefix":
			if s, ok := v.(string); ok {
				a.Prefix = s
			}
		case "as":
			if s, ok := v.(string); ok {
				a.As = s
			}
		case "namespace":
			if s, ok := v.(string); ok {
				a.Namespace = s
			}
		default:
			a.Extras[k] = v
		}
	}

	if a.Uses == nil {
		// the upstream findCallable: scan numeric-keyed entries for a callable.
		// In Go we only have string-keyed maps here, so fall back to a missing
		// action so the error surfaces lazily at dispatch time.
		a.Uses = missingAction("").Uses
	}

	return a, nil
}

func parseStringAction(s string) (*Action, error) {
	if !strings.Contains(s, "@") {
		// Treat as invokable controller name: "Pkg.Type" → "Pkg.Type@Invoke".
		return &Action{Uses: s + "@Invoke", Controller: s + "@Invoke"}, nil
	}

	return &Action{Uses: s, Controller: s}, nil
}

// isInvokable reports whether v has an exported "Invoke" method, the Go
// equivalent of PHP's __invoke magic method.
func isInvokable(v reflect.Value) bool {
	if !v.IsValid() {
		return false
	}

	t := v.Type()

	if _, ok := t.MethodByName("Invoke"); ok {
		return true
	}

	if t.Kind() != reflect.Ptr {
		if _, ok := reflect.PtrTo(t).MethodByName("Invoke"); ok {
			return true
		}
	}

	return false
}

// ContainsSerializedClosure is the parity stub for
// RouteAction::containsSerializedClosure. Go has no equivalent of PHP's
// serialized closures, so this always reports false. The function is retained
// so [RouteSignatureParameters] can call it without conditionals.
func ContainsSerializedClosure(_ *Action) bool { return false }
