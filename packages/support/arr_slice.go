package support

import (
	"cmp"
	"errors"
	"math"
	"math/rand/v2"
	"slices"
)

// ArrCollapse merges a slice of slices into a single flat slice.
func ArrCollapse[T any](items [][]T) []T {
	total := 0

	for _, s := range items {
		total += len(s)
	}

	result := make([]T, 0, total)

	for _, s := range items {
		result = append(result, s...)
	}

	return result
}

// ArrFirst returns the first element matching the predicate.
// If no predicate is given, returns the first element.
// Returns (value, true) if found, (zero, false) otherwise.
func ArrFirst[T any](items []T, predicate ...func(T, int) bool) (T, bool) {
	if len(items) == 0 {
		var zero T

		return zero, false
	}

	if len(predicate) == 0 || predicate[0] == nil {
		return items[0], true
	}

	fn := predicate[0]

	for i, item := range items {
		if fn(item, i) {
			return item, true
		}
	}

	var zero T

	return zero, false
}

// ArrLast returns the last element matching the predicate.
// If no predicate is given, returns the last element.
// Returns (value, true) if found, (zero, false) otherwise.
func ArrLast[T any](items []T, predicate ...func(T, int) bool) (T, bool) {
	if len(items) == 0 {
		var zero T

		return zero, false
	}

	if len(predicate) == 0 || predicate[0] == nil {
		return items[len(items)-1], true
	}

	fn := predicate[0]

	for i := len(items) - 1; i >= 0; i-- {
		if fn(items[i], i) {
			return items[i], true
		}
	}

	var zero T

	return zero, false
}

// ArrFlatten flattens nested []any slices to the given depth.
// Default depth is infinite (math.MaxInt).
func ArrFlatten(items []any, depth ...int) []any {
	d := math.MaxInt

	if len(depth) > 0 {
		d = depth[0]
	}

	return flattenAny(items, d)
}

func flattenAny(items []any, depth int) []any {
	var result []any

	for _, item := range items {
		if sub, ok := item.([]any); ok && depth > 0 {
			result = append(result, flattenAny(sub, depth-1)...)
		} else {
			result = append(result, item)
		}
	}

	return result
}

// ArrPrepend inserts a value at the beginning of a slice.
func ArrPrepend[T any](items []T, value T) []T {
	return append([]T{value}, items...)
}

// ArrRandom returns random elements from a slice.
// Without a count, returns a single-element slice.
// With count, returns that many unique random elements.
func ArrRandom[T any](items []T, count ...int) ([]T, error) {
	if len(items) == 0 {
		return nil, errors.New("cannot get random element from empty slice")
	}

	n := 1

	if len(count) > 0 {
		n = count[0]
	}

	if n > len(items) {
		return nil, errors.New("requested count exceeds slice length")
	}

	if n <= 0 {
		return nil, errors.New("requested count must be positive")
	}

	indices := rand.Perm(len(items))
	result := make([]T, n)

	for i := 0; i < n; i++ {
		result[i] = items[indices[i]]
	}

	return result, nil
}

// ArrSort returns a sorted copy of the slice.
func ArrSort[T cmp.Ordered](items []T) []T {
	sorted := make([]T, len(items))
	copy(sorted, items)

	slices.Sort(sorted)

	return sorted
}

// ArrSortFunc returns a sorted copy using a key-extraction function.
func ArrSortFunc[T any, K cmp.Ordered](items []T, fn func(T) K) []T {
	sorted := make([]T, len(items))
	copy(sorted, items)

	slices.SortStableFunc(sorted, func(a, b T) int {
		return cmp.Compare(fn(a), fn(b))
	})

	return sorted
}

// ArrWhere filters a slice using a predicate function.
func ArrWhere[T any](items []T, fn func(T, int) bool) []T {
	var result []T

	for i, item := range items {
		if fn(item, i) {
			result = append(result, item)
		}
	}

	return result
}

// ArrWrap ensures the value is a slice. If it is already a []T, returns it.
// If nil, returns an empty slice. Otherwise wraps in a single-element slice.
func ArrWrap[T any](value any) []T {
	if value == nil {
		return []T{}
	}

	if s, ok := value.([]T); ok {
		return s
	}

	if v, ok := value.(T); ok {
		return []T{v}
	}

	return []T{}
}
