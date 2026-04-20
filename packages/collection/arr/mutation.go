package arr

// Prepend inserts an element at the beginning of the slice and returns the new slice.
func Prepend[T any](items []T, value T) []T {
	return append([]T{value}, items...)
}

// Push appends one or more elements to the end of the slice and returns the new slice.
func Push[T any](items []T, values ...T) []T {
	return append(items, values...)
}

// Set returns a new slice with the element at the given index replaced.
// If the index is out of bounds, the original slice is returned unchanged.
func Set[T any](items []T, index int, value T) []T {
	if index < 0 || index >= len(items) {
		result := make([]T, len(items))
		copy(result, items)

		return result
	}

	result := make([]T, len(items))
	copy(result, items)
	result[index] = value

	return result
}

// Forget returns a new slice with the element at the given index removed.
// If the index is out of bounds, the original slice is returned unchanged.
func Forget[T any](items []T, index int) []T {
	if index < 0 || index >= len(items) {
		result := make([]T, len(items))
		copy(result, items)

		return result
	}

	result := make([]T, 0, len(items)-1)
	result = append(result, items[:index]...)
	result = append(result, items[index+1:]...)

	return result
}

// Pull removes the element at the given index from the slice and returns both
// the removed element and the remaining slice. If the index is out of bounds,
// the zero value and the original slice are returned.
func Pull[T any](items []T, index int) (T, []T) {
	if index < 0 || index >= len(items) {
		var zero T
		result := make([]T, len(items))
		copy(result, items)

		return zero, result
	}

	value := items[index]
	result := make([]T, 0, len(items)-1)
	result = append(result, items[:index]...)
	result = append(result, items[index+1:]...)

	return value, result
}
