package arr

import "fmt"

// First returns the first element that matches the optional callback.
// If no callback is provided, the first element of the slice is returned.
// The second return value reports whether a match was found.
func First[T any](items []T, callbacks ...func(T, int) bool) (T, bool) {
	if len(callbacks) == 0 || callbacks[0] == nil {
		if len(items) > 0 {
			return items[0], true
		}

		var zero T

		return zero, false
	}

	callback := callbacks[0]

	for i, item := range items {
		if callback(item, i) {
			return item, true
		}
	}

	var zero T

	return zero, false
}

// Last returns the last element that matches the optional callback.
// If no callback is provided, the last element of the slice is returned.
// The second return value reports whether a match was found.
func Last[T any](items []T, callbacks ...func(T, int) bool) (T, bool) {
	if len(callbacks) == 0 || callbacks[0] == nil {
		if len(items) > 0 {
			return items[len(items)-1], true
		}

		var zero T

		return zero, false
	}

	callback := callbacks[0]

	for i := len(items) - 1; i >= 0; i-- {
		if callback(items[i], i) {
			return items[i], true
		}
	}

	var zero T

	return zero, false
}

// Get returns the element at the given index.
// If the index is out of bounds, the first default value is returned,
// or the zero value of T.
func Get[T any](items []T, index int, defaults ...T) T {
	if index >= 0 && index < len(items) {
		return items[index]
	}

	if len(defaults) > 0 {
		return defaults[0]
	}

	var zero T

	return zero
}

// Sole returns the only element matching the optional callback.
// It returns an error if zero or more than one element matches.
func Sole[T any](items []T, callbacks ...func(T, int) bool) (T, error) {
	if len(callbacks) == 0 || callbacks[0] == nil {
		if len(items) == 1 {
			return items[0], nil
		}

		if len(items) == 0 {
			var zero T

			return zero, fmt.Errorf("no items found")
		}

		var zero T

		return zero, fmt.Errorf("multiple items found: %d items", len(items))
	}

	callback := callbacks[0]

	var result T
	found := 0

	for i, item := range items {
		if callback(item, i) {
			result = item
			found++

			if found > 1 {
				var zero T

				return zero, fmt.Errorf("multiple items found: %d+ items", found)
			}
		}
	}

	if found == 0 {
		var zero T

		return zero, fmt.Errorf("no items found")
	}

	return result, nil
}

// Take returns up to limit elements from the slice.
// A positive limit takes from the front; a negative limit takes from the end.
func Take[T any](items []T, limit int) []T {
	if limit < 0 {
		start := len(items) + limit

		if start < 0 {
			start = 0
		}

		result := make([]T, len(items)-start)
		copy(result, items[start:])

		return result
	}

	if limit > len(items) {
		limit = len(items)
	}

	result := make([]T, limit)
	copy(result, items[:limit])

	return result
}
