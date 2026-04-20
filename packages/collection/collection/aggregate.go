package collection

import (
	"cmp"
	"slices"

	"github.com/bedrock/packages/collection/support"
)

// Median returns the median value of a float64 collection.
func Median(c *Collection[float64]) float64 {
	if len(c.items) == 0 {
		return 0
	}

	sorted := make([]float64, len(c.items))
	copy(sorted, c.items)

	slices.Sort(sorted)

	mid := len(sorted) / 2

	if len(sorted)%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}

	return sorted[mid]
}

// MedianBy returns the median value extracted from each item by the given function.
func MedianBy[T any](c *Collection[T], valueFunc func(T) float64) float64 {
	return Median(Map(c, func(item T, _ int) float64 {
		return valueFunc(item)
	}))
}

// Mode returns the most frequently occurring values in the collection.
func Mode[T comparable](c *Collection[T]) []T {
	if len(c.items) == 0 {
		return nil
	}

	counts := make(map[T]int)
	maxCount := 0

	for _, item := range c.items {
		counts[item]++

		if counts[item] > maxCount {
			maxCount = counts[item]
		}
	}

	result := make([]T, 0)

	for item, count := range counts {
		if count == maxCount {
			result = append(result, item)
		}
	}

	return result
}

// Sum returns the sum of all items in a numeric collection.
func Sum[T support.Numeric](c *Collection[T]) T {
	var total T

	for _, item := range c.items {
		total += item
	}

	return total
}

// SumBy returns the sum of values extracted from each item by the given function.
func SumBy[T any, N support.Numeric](c *Collection[T], valueFunc func(T) N) N {
	var total N

	for _, item := range c.items {
		total += valueFunc(item)
	}

	return total
}

// Avg returns the arithmetic mean of all items in a numeric collection.
func Avg[T support.Numeric](c *Collection[T]) float64 {
	if len(c.items) == 0 {
		return 0
	}

	return float64(Sum(c)) / float64(len(c.items))
}

// AvgBy returns the arithmetic mean of values extracted from each item by the given function.
func AvgBy[T any, N support.Numeric](c *Collection[T], valueFunc func(T) N) float64 {
	if len(c.items) == 0 {
		return 0
	}

	return float64(SumBy(c, valueFunc)) / float64(len(c.items))
}

// Average is an alias for Avg.
func Average[T support.Numeric](c *Collection[T]) float64 {
	return Avg(c)
}

// Min returns the minimum value in an ordered collection.
// The second return value indicates whether the collection was non-empty.
func Min[T cmp.Ordered](c *Collection[T]) (T, bool) {
	if len(c.items) == 0 {
		var zero T

		return zero, false
	}

	result := c.items[0]

	for _, item := range c.items[1:] {
		if item < result {
			result = item
		}
	}

	return result, true
}

// MinBy returns the item with the minimum key as determined by the given function.
// The second return value indicates whether the collection was non-empty.
func MinBy[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K) (T, bool) {
	if len(c.items) == 0 {
		var zero T

		return zero, false
	}

	result := c.items[0]
	minKey := keyFunc(result)

	for _, item := range c.items[1:] {
		k := keyFunc(item)

		if k < minKey {
			minKey = k
			result = item
		}
	}

	return result, true
}

// Max returns the maximum value in an ordered collection.
// The second return value indicates whether the collection was non-empty.
func Max[T cmp.Ordered](c *Collection[T]) (T, bool) {
	if len(c.items) == 0 {
		var zero T

		return zero, false
	}

	result := c.items[0]

	for _, item := range c.items[1:] {
		if item > result {
			result = item
		}
	}

	return result, true
}

// MaxBy returns the item with the maximum key as determined by the given function.
// The second return value indicates whether the collection was non-empty.
func MaxBy[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K) (T, bool) {
	if len(c.items) == 0 {
		var zero T

		return zero, false
	}

	result := c.items[0]
	maxKey := keyFunc(result)

	for _, item := range c.items[1:] {
		k := keyFunc(item)

		if k > maxKey {
			maxKey = k
			result = item
		}
	}

	return result, true
}
