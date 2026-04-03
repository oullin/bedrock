package arr

// Map applies a callback to each element and returns a slice of the results.
func Map[T any, R any](items []T, callback func(T, int) R) []R {
	result := make([]R, len(items))

	for i, item := range items {
		result[i] = callback(item, i)
	}

	return result
}

// MapWithKeys applies a callback to each element, producing key-value pairs
// that are collected into a map.
func MapWithKeys[T any, K comparable, V any](items []T, callback func(T) (K, V)) map[K]V {
	result := make(map[K]V)

	for _, item := range items {
		key, value := callback(item)
		result[key] = value
	}

	return result
}

// MapSpread applies a callback to each element and returns a slice of the results.
// In Go this is identical to [Map].
func MapSpread[T any, R any](items []T, callback func(T, int) R) []R {
	return Map(items, callback)
}

// KeyBy indexes the slice elements by the key returned from keyFunc.
func KeyBy[T any, K comparable](items []T, keyFunc func(T) K) map[K]T {
	result := make(map[K]T)

	for _, item := range items {
		result[keyFunc(item)] = item
	}

	return result
}

// Pluck extracts a value from each element using valueFunc and returns
// the collected values as a slice.
func Pluck[T any, V any](items []T, valueFunc func(T) V) []V {
	result := make([]V, len(items))

	for i, item := range items {
		result[i] = valueFunc(item)
	}

	return result
}

// Divide returns two slices: one containing the indices and one containing the values.
func Divide[T any](items []T) ([]int, []T) {
	keys := make([]int, len(items))
	values := make([]T, len(items))

	for i, item := range items {
		keys[i] = i
		values[i] = item
	}

	return keys, values
}
