package collection

// Value returns the given value as-is.
func Value[T any](value T) T {
	return value
}

// ValueFunc calls the given callback and returns its result.
func ValueFunc[T any](callback func() T) T {
	return callback()
}

// Head returns the first element of a slice and true,
// or the zero value and false if the slice is empty.
func Head[T any](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T

		return zero, false
	}

	return items[0], true
}

// Last returns the last element of a slice and true,
// or the zero value and false if the slice is empty.
func Last[T any](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T

		return zero, false
	}

	return items[len(items)-1], true
}

// WhenValue returns value if condition is true, otherwise returns the first
// default or the zero value of T.
func WhenValue[T any](condition bool, value T, defaults ...T) T {
	if condition {
		return value
	}

	if len(defaults) > 0 {
		return defaults[0]
	}

	var zero T

	return zero
}

// WhenFunc calls callback and returns its result if condition is true,
// otherwise calls the first default callback or returns the zero value of T.
func WhenFunc[T any](condition bool, callback func() T, defaults ...func() T) T {
	if condition {
		return callback()
	}

	if len(defaults) > 0 {
		return defaults[0]()
	}

	var zero T

	return zero
}
