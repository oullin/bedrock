package conditionable

// When executes callback with target and value if value is truthy.
// If value is falsy and a defaultFn is provided, it executes defaultFn instead.
// Returns target unchanged when no branch executes.
func When[T any, V any](target T, value V, callback func(T, V) T, defaultFn ...func(T, V) T) T {
	if Truthy(value) {
		return callback(target, value)
	}

	if len(defaultFn) > 0 && defaultFn[0] != nil {
		return defaultFn[0](target, value)
	}

	return target
}

// WhenFunc resolves the condition by calling fn with target, then behaves like When.
func WhenFunc[T any, V any](target T, fn func(T) V, callback func(T, V) T, defaultFn ...func(T, V) T) T {
	value := fn(target)

	if Truthy(value) {
		return callback(target, value)
	}

	if len(defaultFn) > 0 && defaultFn[0] != nil {
		return defaultFn[0](target, value)
	}

	return target
}

// Unless executes callback with target and value if value is falsy.
// If value is truthy and a defaultFn is provided, it executes defaultFn instead.
// Returns target unchanged when no branch executes.
func Unless[T any, V any](target T, value V, callback func(T, V) T, defaultFn ...func(T, V) T) T {
	if !Truthy(value) {
		return callback(target, value)
	}

	if len(defaultFn) > 0 && defaultFn[0] != nil {
		return defaultFn[0](target, value)
	}

	return target
}

// UnlessFunc resolves the condition by calling fn with target, then behaves like Unless.
func UnlessFunc[T any, V any](target T, fn func(T) V, callback func(T, V) T, defaultFn ...func(T, V) T) T {
	value := fn(target)

	if !Truthy(value) {
		return callback(target, value)
	}

	if len(defaultFn) > 0 && defaultFn[0] != nil {
		return defaultFn[0](target, value)
	}

	return target
}
