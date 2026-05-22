package routegen

import "strings"

// Param holds the metadata for a single URI parameter.
type Param struct {
	// Name is the raw parameter name from the URI (e.g. "post", "user").
	Name string

	// Optional is true when the URI uses the {param?} syntax.
	Optional bool

	// Key is the custom binding field declared via {model:key} syntax.
	// Empty string means the default key ("id") is used.
	Key string

	// Default is the string form of the default value supplied via route
	// defaults or middleware URL::defaults(). Empty when no default exists.
	Default string
}

// Placeholder returns the URI placeholder string for this parameter.
// Optional params use {name?}; required params use {name}.
func (p Param) Placeholder() string {
	if p.Optional {
		return "{" + p.Name + "?}"
	}

	return "{" + p.Name + "}"
}

// SafeName returns a TypeScript-safe version of the parameter name.
func (p Param) SafeName() string {
	return SafeMethod(p.Name, "Param")
}

// TSTypes returns the TypeScript type union for this parameter.
// Without model binding information only string|number is emitted,
// matching the PHP fallback behaviour.
func (p Param) TSTypes() string {
	// Go routes do not carry PHP reflection data so we default to
	// string | number for all route parameters — identical to the PHP
	// fallback when no ReflectionParameter is available.
	return "string | number"
}

// resolveParamSafeNames returns a set of names that are already in use for the
// generated function to avoid collisions with the standard variable names
// (args, options, parsedArgs).
func resolveParamSafeNames(method string) map[string]string {
	reserved := map[string]string{
		"args":       "routeArgs",
		"options":    "routeOptions",
		"parsedArgs": "routeParsedArgs",
	}

	result := make(map[string]string, len(reserved))

	for k, v := range reserved {
		if method == k {
			result[k] = v
		} else {
			result[k] = k
		}
	}

	return result
}

// argsVar returns the variable name for the positional arguments object,
// accounting for method-name collisions.
func argsVar(method string) string {
	names := resolveParamSafeNames(method)

	return names["args"]
}

// optionsVar returns the variable name for the query-options argument.
func optionsVar(method string) string {
	names := resolveParamSafeNames(method)

	return names["options"]
}

// parsedArgsVar returns the variable name for the resolved parameter object.
func parsedArgsVar(method string) string {
	names := resolveParamSafeNames(method)

	return names["parsedArgs"]
}

// hasOptional reports whether any parameter in the slice is optional.
func hasOptional(params []Param) bool {
	for _, p := range params {
		if p.Optional {
			return true
		}
	}

	return false
}

// allOptional reports whether every parameter in the slice is optional.
func allOptional(params []Param) bool {
	if len(params) == 0 {
		return false
	}

	for _, p := range params {
		if !p.Optional {
			return false
		}
	}

	return true
}

// optionalNames returns the names of optional parameters in order.
func optionalNames(params []Param) []string {
	names := make([]string, 0)

	for _, p := range params {
		if p.Optional {
			names = append(names, p.Name)
		}
	}

	return names
}

// jsonStringSlice renders a Go string slice as a JSON array literal
// (e.g. ["get","head"]).
func jsonStringSlice(ss []string) string {
	var b strings.Builder

	b.WriteByte('[')

	for i, s := range ss {
		if i > 0 {
			b.WriteByte(',')
		}

		b.WriteByte('"')
		b.WriteString(s)
		b.WriteByte('"')
	}

	b.WriteByte(']')

	return b.String()
}

// jsonString renders a Go string as a JSON string literal (e.g. "/posts/{post}").
func jsonString(s string) string {
	// Minimal JSON encoding — escape backslash and double quote.
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)

	return `"` + s + `"`
}
