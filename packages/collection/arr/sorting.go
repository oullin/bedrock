package arr

import (
	"math/rand/v2"
	"sort"
)

// Sort returns a new slice sorted using the provided less function.
// The sort is stable.
func Sort[T any](items []T, less func(a, b T) bool) []T {
	result := make([]T, len(items))
	copy(result, items)

	sort.SliceStable(result, func(i, j int) bool {
		return less(result[i], result[j])
	})

	return result
}

// SortDesc returns a new slice sorted in descending order using the provided less function.
func SortDesc[T any](items []T, less func(a, b T) bool) []T {
	return Sort(items, func(a, b T) bool {
		return less(b, a)
	})
}

// SortRecursive sorts a slice using the provided less function.
// For flat slices this behaves identically to [Sort].
func SortRecursive[T any](items []T, less func(a, b T) bool) []T {
	return Sort(items, less)
}

// SortRecursiveDesc sorts a slice in descending order using the provided less function.
func SortRecursiveDesc[T any](items []T, less func(a, b T) bool) []T {
	return SortDesc(items, less)
}

// Shuffle returns a new slice with the elements in random order.
func Shuffle[T any](items []T) []T {
	result := make([]T, len(items))
	copy(result, items)

	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	return result
}

// Random returns a new slice containing count random elements.
// If the count is omitted, it defaults to 1.
func Random[T any](items []T, counts ...int) []T {
	count := 1

	if len(counts) > 0 {
		count = counts[0]
	}

	shuffled := Shuffle(items)

	if count >= len(shuffled) {
		return shuffled
	}

	result := make([]T, count)
	copy(result, shuffled[:count])

	return result
}
