package support

// Optional wraps a potentially-nil value and provides safe access,
// preventing nil pointer dereferences when accessing methods or properties.
type Optional[T any] struct {
	value *T
}

// Some wraps a value in an Optional, indicating the value is present.
func Some[T any](value T) Optional[T] {
	return Optional[T]{value: &value}
}

// None returns an empty Optional with no value.
func None[T any]() Optional[T] {
	return Optional[T]{value: nil}
}

// Opt wraps an existing pointer in an Optional.
// If ptr is nil, the Optional will be empty.
func Opt[T any](ptr *T) Optional[T] {
	return Optional[T]{value: ptr}
}

// Get returns the wrapped value and whether it was present.
func (o Optional[T]) Get() (T, bool) {
	if o.value == nil {
		var zero T

		return zero, false
	}

	return *o.value, true
}

// OrElse returns the wrapped value if present, or the given default.
func (o Optional[T]) OrElse(def T) T {
	if o.value == nil {
		return def
	}

	return *o.value
}

// OrElseGet returns the wrapped value if present, or calls the given
// function to produce a default.
func (o Optional[T]) OrElseGet(fn func() T) T {
	if o.value == nil {
		return fn()
	}

	return *o.value
}

// IsPresent reports whether the Optional contains a value.
func (o Optional[T]) IsPresent() bool {
	return o.value != nil
}

// IsEmpty reports whether the Optional does not contain a value.
func (o Optional[T]) IsEmpty() bool {
	return o.value == nil
}

// IfPresent calls the given function with the wrapped value if it is present.
func (o Optional[T]) IfPresent(fn func(T)) {
	if o.value != nil {
		fn(*o.value)
	}
}

// IfPresentOrElse calls the given function with the wrapped value if present,
// or calls the empty function if absent.
func (o Optional[T]) IfPresentOrElse(fn func(T), empty func()) {
	if o.value != nil {
		fn(*o.value)
	} else if empty != nil {
		empty()
	}
}

// Map applies the given function to the wrapped value if present,
// returning a new Optional with the result.
func (o Optional[T]) Map(fn func(T) T) Optional[T] {
	if o.value == nil {
		return None[T]()
	}

	result := fn(*o.value)

	return Some(result)
}

// Filter returns the Optional if the value is present and the predicate
// returns true; otherwise returns None.
func (o Optional[T]) Filter(predicate func(T) bool) Optional[T] {
	if o.value != nil && predicate(*o.value) {
		return o
	}

	return None[T]()
}

// Pointer returns the underlying pointer. May be nil if empty.
func (o Optional[T]) Pointer() *T {
	return o.value
}
