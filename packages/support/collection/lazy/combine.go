package lazy

// GroupBy groups items by the key returned by the given function.
// This requires eager evaluation to build the groups.
func GroupBy[T any, K comparable](lc *Collection[T], keyFunc func(T) K) map[K]*Collection[T] {
	groups := make(map[K][]T)
	lc.source(func(item T) bool {
		key := keyFunc(item)
		groups[key] = append(groups[key], item)

		return true
	})
	result := make(map[K]*Collection[T])

	for key, items := range groups {
		result[key] = From(items)
	}

	return result
}

// KeyBy indexes items by the key returned by the given function.
// Duplicate keys cause the later value to overwrite the earlier one.
func KeyBy[T any, K comparable](lc *Collection[T], keyFunc func(T) K) map[K]T {
	result := make(map[K]T)
	lc.source(func(item T) bool {
		result[keyFunc(item)] = item

		return true
	})

	return result
}

// CountBy counts occurrences of each key returned by the given function.
func CountBy[T any, K comparable](lc *Collection[T], keyFunc func(T) K) map[K]int {
	result := make(map[K]int)
	lc.source(func(item T) bool {
		result[keyFunc(item)]++

		return true
	})

	return result
}
