package arr

// Only returns the elements at the given indices.
func Only[T any](items []T, indices []int) []T {
	result := make([]T, 0, len(indices))

	for _, idx := range indices {
		if idx >= 0 && idx < len(items) {
			result = append(result, items[idx])
		}
	}

	return result
}

// Except returns all elements except those at the given indices.
func Except[T any](items []T, indices []int) []T {
	excludeSet := make(map[int]bool, len(indices))

	for _, idx := range indices {
		excludeSet[idx] = true
	}

	result := make([]T, 0)

	for i, item := range items {
		if !excludeSet[i] {
			result = append(result, item)
		}
	}

	return result
}

// Where returns the elements for which the callback returns true.
func Where[T any](items []T, callback func(T, int) bool) []T {
	result := make([]T, 0)

	for i, item := range items {
		if callback(item, i) {
			result = append(result, item)
		}
	}

	return result
}

// WhereNotNull returns all elements that are not the zero value of their type.
func WhereNotNull[T comparable](items []T) []T {
	var zero T
	result := make([]T, 0)

	for _, item := range items {
		if item != zero {
			result = append(result, item)
		}
	}

	return result
}

// Reject returns the elements for which the callback returns false.
func Reject[T any](items []T, callback func(T, int) bool) []T {
	return Where(items, func(item T, index int) bool {
		return !callback(item, index)
	})
}

// Partition splits elements into two slices: those that pass the callback
// and those that do not.
func Partition[T any](items []T, callback func(T, int) bool) ([]T, []T) {
	pass := make([]T, 0)
	fail := make([]T, 0)

	for i, item := range items {
		if callback(item, i) {
			pass = append(pass, item)
		} else {
			fail = append(fail, item)
		}
	}

	return pass, fail
}
