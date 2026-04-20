package lazy

// Reduce reduces the lazy collection to a single value by applying the callback
// to an accumulator and each item in sequence.
func Reduce[T any, R any](lc *Collection[T], callback func(R, T, int) R, initial R) R {
	result := initial
	idx := 0
	lc.source(func(item T) bool {
		result = callback(result, item, idx)
		idx++

		return true
	})

	return result
}
