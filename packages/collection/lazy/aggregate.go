package lazy

import (
	"cmp"

	"github.com/bedrock/packages/collection/support"
)

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

// Sum returns the sum of all numeric values in the lazy collection.
func Sum[T support.Numeric](lc *Collection[T]) T {
	var total T
	lc.source(func(item T) bool {
		total += item

		return true
	})

	return total
}

// Avg returns the arithmetic mean of all numeric values in the lazy collection.
func Avg[T support.Numeric](lc *Collection[T]) float64 {
	var total T
	count := 0
	lc.source(func(item T) bool {
		total += item
		count++

		return true
	})
	if count == 0 {
		return 0
	}

	return float64(total) / float64(count)
}

// Min returns the minimum value in the lazy collection.
func Min[T cmp.Ordered](lc *Collection[T]) (T, bool) {
	var result T
	found := false
	lc.source(func(item T) bool {
		if !found || item < result {
			result = item
			found = true
		}

		return true
	})

	return result, found
}

// Max returns the maximum value in the lazy collection.
func Max[T cmp.Ordered](lc *Collection[T]) (T, bool) {
	var result T
	found := false
	lc.source(func(item T) bool {
		if !found || item > result {
			result = item
			found = true
		}

		return true
	})

	return result, found
}
