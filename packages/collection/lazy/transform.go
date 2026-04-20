package lazy

// Filter returns a new lazy collection containing only items for which the callback returns true.
func (lc *Collection[T]) Filter(callback func(T, int) bool) *Collection[T] {
	return New(func(yield func(T) bool) {
		idx := 0
		lc.source(func(item T) bool {
			if callback(item, idx) {
				if !yield(item) {
					return false
				}
			}

			idx++

			return true
		})
	})
}

// Reject returns a new lazy collection excluding items for which the callback returns true.
func (lc *Collection[T]) Reject(callback func(T, int) bool) *Collection[T] {
	return lc.Filter(func(item T, index int) bool {
		return !callback(item, index)
	})
}

// Map transforms each item using the callback, returning a new lazy collection.
func Map[T any, R any](lc *Collection[T], callback func(T, int) R) *Collection[R] {
	return New(func(yield func(R) bool) {
		idx := 0
		lc.source(func(item T) bool {
			if !yield(callback(item, idx)) {
				return false
			}

			idx++

			return true
		})
	})
}

// FlatMap transforms each item into a slice and flattens the results into a new lazy collection.
func FlatMap[T any, R any](lc *Collection[T], callback func(T, int) []R) *Collection[R] {
	return New(func(yield func(R) bool) {
		idx := 0
		lc.source(func(item T) bool {
			for _, r := range callback(item, idx) {
				if !yield(r) {
					return false
				}
			}

			idx++

			return true
		})
	})
}

// Every report whether all items satisfy the callback.
func (lc *Collection[T]) Every(callback func(T, int) bool) bool {
	result := true
	idx := 0
	lc.source(func(item T) bool {
		if !callback(item, idx) {
			result = false

			return false
		}

		idx++

		return true
	})

	return result
}

// Unique returns a lazy collection containing only unique items as determined by the key function.
func Unique[T any, K comparable](lc *Collection[T], keyFunc func(T) K) *Collection[T] {
	return New(func(yield func(T) bool) {
		seen := make(map[K]bool)
		lc.source(func(item T) bool {
			key := keyFunc(item)

			if !seen[key] {
				seen[key] = true

				return yield(item)
			}

			return true
		})
	})
}

// Pluck extracts a value from each item using the given function, returning a new lazy collection.
func Pluck[T any, V any](lc *Collection[T], valueFunc func(T) V) *Collection[V] {
	return Map(lc, func(item T, _ int) V {
		return valueFunc(item)
	})
}
