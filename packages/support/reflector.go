package support

import "reflect"

// TypeName returns the concrete type name for a value.
func TypeName(value any) string {
	if value == nil {
		return ""
	}

	t := reflect.TypeOf(value)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Name() != "" {
		return t.Name()
	}

	return t.String()
}

// IsCallable reports whether value is a function.
func IsCallable(value any) bool {
	if value == nil {
		return false
	}

	return reflect.TypeOf(value).Kind() == reflect.Func
}

// Implements reports whether value implements the target interface pointer.
func Implements(value any, target any) bool {
	if value == nil || target == nil {
		return false
	}

	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Pointer || targetType.Elem().Kind() != reflect.Interface {
		return false
	}

	valueType := reflect.TypeOf(value)

	return valueType.Implements(targetType.Elem())
}
