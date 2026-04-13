package conditionable

import "reflect"

// Truthy reports whether value is truthy using PHP-compatible semantics.
// The following values are considered falsy: nil, false, 0 (any numeric type),
// and the empty string. Everything else is truthy.
func Truthy(value any) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case int8:
		return v != 0
	case int16:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case uint:
		return v != 0
	case uint8:
		return v != 0
	case uint16:
		return v != 0
	case uint32:
		return v != 0
	case uint64:
		return v != 0
	case float32:
		return v != 0
	case float64:
		return v != 0
	case string:
		return v != ""
	}

	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return !rv.IsNil()
	}

	return true
}
