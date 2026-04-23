package support

// When returns the callback result when condition is true, otherwise value or
// the optional default callback result.
func When[T any](value T, condition bool, callback func(T) T, defaults ...func(T) T) T {
	if condition {
		return callback(value)
	}

	if len(defaults) > 0 && defaults[0] != nil {
		return defaults[0](value)
	}

	return value
}

// Unless returns the callback result when condition is false, otherwise value
// or the optional default callback result.
func Unless[T any](value T, condition bool, callback func(T) T, defaults ...func(T) T) T {
	return When(value, !condition, callback, defaults...)
}

// TapValue invokes callback with value and returns value.
func TapValue[T any](value T, callback ...func(T)) T {
	for _, cb := range callback {
		if cb != nil {
			cb(value)
		}
	}

	return value
}
