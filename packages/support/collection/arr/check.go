package arr

// Accessible reports whether the given value is non-nil.
func Accessible(value any) bool {
	return value != nil
}

// IsList reports whether the given slice is a list.
// In Go, all slices are sequential lists, so this always returns true.
func IsList[T any](items []T) bool {
	return true
}

// Exists -> it reports whether the given index is valid for the slice.
func Exists[T any](items []T, index int) bool {
	return index >= 0 && index < len(items)
}

// Has reports whether all the given indices are valid for the slice.
func Has[T any](items []T, indices ...int) bool {
	for _, idx := range indices {
		if idx < 0 || idx >= len(items) {
			return false
		}
	}

	return true
}

// HasAny reports whether at least one of the given indices is valid for the slice.
func HasAny[T any](items []T, indices ...int) bool {
	for _, idx := range indices {
		if idx >= 0 && idx < len(items) {
			return true
		}
	}

	return false
}
