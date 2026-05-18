package multisteps

import (
	"fmt"
	"reflect"
)

// argKind discriminates Arg values.
type argKind int

// Arg is a declarative argument binding for a job — a runtime Variable, a
// Response reference to a prior job's output, or a Literal value.
type Arg struct {
	kind    argKind
	name    string
	field   string
	literal any
}

// A is a convenience alias for the Args map shape, keeping call sites compact.
//
//	multisteps.Args(multisteps.A{"userId": multisteps.Variable("uid")})
type A = map[string]Arg

const (
	argLiteral argKind = iota
	argVariable
	argResponse
)

// Variable references a runtime variable bound at Run() time by name.
func Variable(name string) Arg {
	return Arg{kind: argVariable, name: name}
}

// Response references a prior job's output. If field is "" the whole result is
// returned; otherwise field is resolved via reflection (struct field or map key).
func Response(job, field string) Arg {
	return Arg{kind: argResponse, name: job, field: field}
}

// Literal wraps a constant value so it can be passed through the Args map alongside variables and responses.
func Literal(v any) Arg {
	return Arg{kind: argLiteral, literal: v}
}

// Kind reports whether the Arg is Literal/Variable/Response.
func (a Arg) Kind() string {
	switch a.kind {
	case argLiteral:
		return "literal"
	case argVariable:
		return "variable"
	case argResponse:
		return "response"
	}

	return "unknown"
}

// DependsOnJob returns the job name a Response Arg depends on, or "" otherwise.
func (a Arg) DependsOnJob() string {
	if a.kind == argResponse {
		return a.name
	}

	return ""
}

// resolve produces the runtime value of an Arg using the provided variables
// and previous-job responses.
func (a Arg) resolve(vars map[string]any, responses map[string]any) (any, error) {
	switch a.kind {
	case argLiteral:
		return a.literal, nil
	case argVariable:
		value, ok := vars[a.name]

		if !ok {
			return nil, fmt.Errorf("variable %q not provided", a.name)
		}

		return value, nil
	case argResponse:
		raw, ok := responses[a.name]

		if !ok {
			return nil, fmt.Errorf("response %q not available", a.name)
		}

		if a.field == "" {
			return raw, nil
		}

		return resolveField(raw, a.field)
	}

	return nil, fmt.Errorf("unknown arg kind")
}

func resolveField(raw any, field string) (any, error) {
	if raw == nil {
		return nil, &UnresolvedResponseError{Field: field}
	}

	switch m := raw.(type) {
	case map[string]any:
		value, ok := m[field]

		if !ok {
			return nil, &UnresolvedResponseError{Field: field}
		}

		return value, nil
	}

	v := reflect.ValueOf(raw)

	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil, &UnresolvedResponseError{Field: field}
		}

		v = v.Elem()
	}

	if v.Kind() == reflect.Map {
		key := reflect.ValueOf(field)

		if !key.Type().AssignableTo(v.Type().Key()) {
			return nil, &UnresolvedResponseError{Field: field}
		}

		found := v.MapIndex(key)

		if !found.IsValid() {
			return nil, &UnresolvedResponseError{Field: field}
		}

		return found.Interface(), nil
	}

	if v.Kind() == reflect.Struct {
		fv := v.FieldByName(field)

		if !fv.IsValid() {
			return nil, &UnresolvedResponseError{Field: field}
		}

		return fv.Interface(), nil
	}

	return nil, &UnresolvedResponseError{Field: field}
}

// As is a typed helper for extracting a field from a Result.Responses entry.
func As[T any](r Result, job, field string) (T, error) {
	var zero T

	raw, ok := r.Responses[job]

	if !ok {
		return zero, fmt.Errorf("response %q not available", job)
	}

	if field == "" {
		out, ok := raw.(T)

		if !ok {
			return zero, fmt.Errorf("response %q is %T, not %T", job, raw, zero)
		}

		return out, nil
	}

	value, err := resolveField(raw, field)

	if err != nil {
		return zero, err
	}

	out, ok := value.(T)

	if !ok {
		return zero, fmt.Errorf("response %q field %q is %T, not %T", job, field, value, zero)
	}

	return out, nil
}
