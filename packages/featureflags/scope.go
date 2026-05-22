package featureflags

import (
	"fmt"
	"reflect"
	"strconv"
)

// Scopeable is implemented by types that control their own scope identifier.
type Scopeable interface {
	FeatureScopeIdentifier() string
}

// NullScope is the serialized form of a nil scope, matching the upstream "__laravel_null" convention adapted for this package.
const NullScope = "__null"

// SerializeScope converts an arbitrary Go value to a stable string key
// suitable for use as a storage key.
//
// Rules applied in order:
//
//	nil                       → "__null"
//	string                    → as-is
//	bool                      → "true" | "false"
//	int, int8, int16, …       → fmt.Sprintf("%d", v)
//	uint, uint8, uint16, …    → fmt.Sprintf("%d", v)
//	Scopeable                 → v.FeatureScopeIdentifier()
//	fmt.Stringer              → v.String()
//	struct or *struct         → reflect: looks for exported field "ID" or "Id";
//	                            produces "pkg.TypeName|<id>"
//	otherwise                 → ErrUnserializableScope
func SerializeScope(scope any) (string, error) {
	if scope == nil {
		return NullScope, nil
	}

	switch v := scope.(type) {
	case string:
		return v, nil
	case bool:
		if v {
			return "true", nil
		}

		return "false", nil
	case int:
		return fmt.Sprintf("%d", v), nil
	case int8:
		return fmt.Sprintf("%d", v), nil
	case int16:
		return fmt.Sprintf("%d", v), nil
	case int32:
		return fmt.Sprintf("%d", v), nil
	case int64:
		return fmt.Sprintf("%d", v), nil
	case uint:
		return fmt.Sprintf("%d", v), nil
	case uint8:
		return fmt.Sprintf("%d", v), nil
	case uint16:
		return fmt.Sprintf("%d", v), nil
	case uint32:
		return fmt.Sprintf("%d", v), nil
	case uint64:
		return fmt.Sprintf("%d", v), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case Scopeable:
		return v.FeatureScopeIdentifier(), nil
	case fmt.Stringer:
		return v.String(), nil
	}

	return serializeReflect(scope)
}

// serializeReflect handles struct and *struct types by inspecting for an
// exported "ID" or "Id" field to produce a "pkg.TypeName|<id>" key.
func serializeReflect(scope any) (string, error) {
	rv := reflect.ValueOf(scope)

	// Dereference pointer.
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return NullScope, nil
		}

		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return "", fmt.Errorf("%w: %T", ErrUnserializableScope, scope)
	}

	rt := rv.Type()

	// Look for exported field named "ID" or "Id".
	for _, name := range []string{"ID", "Id"} {
		field := rv.FieldByName(name)

		if !field.IsValid() {
			continue
		}

		// Ensure the field is exported.
		sf, ok := rt.FieldByName(name)

		if !ok || !sf.IsExported() {
			continue
		}

		typeName := fmt.Sprintf("%s.%s", rt.PkgPath(), rt.Name())

		return fmt.Sprintf("%s|%v", typeName, field.Interface()), nil
	}

	return "", fmt.Errorf("%w: %T has no exported ID or Id field", ErrUnserializableScope, scope)
}
