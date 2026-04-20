package collection

import "github.com/bedrock/packages/collection/support"

// Diff returns the items in the collection that are not present in the given slice.
func Diff[T comparable](c *Collection[T], items []T) *Collection[T] {
	lookup := make(map[T]bool, len(items))

	for _, item := range items {
		lookup[item] = true
	}

	result := make([]T, 0)

	for _, item := range c.items {
		if !lookup[item] {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// DiffUsing returns items not present in the given slice, using a custom equality function.
func (c *Collection[T]) DiffUsing(items []T, equals func(T, T) bool) *Collection[T] {
	result := make([]T, 0)

	for _, item := range c.items {
		found := false

		for _, other := range items {
			if equals(item, other) {
				found = true

				break
			}
		}

		if !found {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// Intersect returns the items present in both the collection and the given slice.
func Intersect[T comparable](c *Collection[T], items []T) *Collection[T] {
	lookup := make(map[T]bool, len(items))

	for _, item := range items {
		lookup[item] = true
	}

	result := make([]T, 0)

	for _, item := range c.items {
		if lookup[item] {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// IntersectUsing returns items present in both the collection and the given slice,
// using a custom equality function.
func (c *Collection[T]) IntersectUsing(items []T, equals func(T, T) bool) *Collection[T] {
	result := make([]T, 0)

	for _, item := range c.items {
		for _, other := range items {
			if equals(item, other) {
				result = append(result, item)

				break
			}
		}
	}

	return Collect(result)
}

// Zip merges the collection with each of the given slices element-by-element.
func Zip[T any](c *Collection[T], others ...[]T) *Collection[[]T] {
	maxLen := len(c.items)

	for _, o := range others {
		if len(o) > maxLen {
			maxLen = len(o)
		}
	}

	result := make([][]T, maxLen)

	for i := 0; i < maxLen; i++ {
		group := make([]T, 0, 1+len(others))

		if i < len(c.items) {
			group = append(group, c.items[i])
		} else {
			var zero T

			group = append(group, zero)
		}

		for _, o := range others {
			if i < len(o) {
				group = append(group, o[i])
			} else {
				var zero T

				group = append(group, zero)
			}
		}

		result[i] = group
	}

	return Collect(result)
}

// CrossJoin returns the cross-product of the collection with the given slices.
func CrossJoin[T any](c *Collection[T], others ...[]T) *Collection[[]T] {
	results := [][]T{{}}
	allLists := append([][]T{c.items}, others...)

	for _, list := range allLists {
		var newResults [][]T

		for _, result := range results {
			for _, item := range list {
				newResult := make([]T, len(result)+1)
				copy(newResult, result)
				newResult[len(result)] = item
				newResults = append(newResults, newResult)
			}
		}

		results = newResults
	}

	return Collect(results)
}

// Combine pairs of keys from this collection with values from the given slice,
// returning a collection of Pair values.
func Combine[K any, V any](keys *Collection[K], values []V) *Collection[support.Pair[K, V]] {
	minLen := len(keys.items)

	if len(values) < minLen {
		minLen = len(values)
	}

	result := make([]support.Pair[K, V], minLen)

	for i := 0; i < minLen; i++ {
		result[i] = support.Pair[K, V]{Key: keys.items[i], Value: values[i]}
	}

	return Collect(result)
}

// Collapse flattens a collection of slices into a single, flat collection.
func Collapse[T any](c *Collection[[]T]) *Collection[T] {
	result := make([]T, 0)

	for _, items := range c.items {
		result = append(result, items...)
	}

	return Collect(result)
}

// GroupBy groups the collection's items by a key returned from the given function.
func GroupBy[T any, K comparable](c *Collection[T], keyFunc func(T) K) map[K]*Collection[T] {
	groups := make(map[K]*Collection[T])

	for _, item := range c.items {
		key := keyFunc(item)

		if _, ok := groups[key]; !ok {
			groups[key] = Empty[T]()
		}

		groups[key].Push(item)
	}

	return groups
}

// KeyBy indexes the collection by a key returned from the given function.
func KeyBy[T any, K comparable](c *Collection[T], keyFunc func(T) K) map[K]T {
	result := make(map[K]T)

	for _, item := range c.items {
		result[keyFunc(item)] = item
	}

	return result
}

// CountBy counts how many items produce each key from the given function.
func CountBy[T any, K comparable](c *Collection[T], keyFunc func(T) K) map[K]int {
	result := make(map[K]int)

	for _, item := range c.items {
		result[keyFunc(item)]++
	}

	return result
}

// MapToDictionary maps each item to a key-value pair and groups values by key.
func MapToDictionary[T any, K comparable, V any](c *Collection[T], callback func(T) (K, V)) map[K][]V {
	result := make(map[K][]V)

	for _, item := range c.items {
		key, value := callback(item)
		result[key] = append(result[key], value)
	}

	return result
}

// MapToGroups is an alias for MapToDictionary.
func MapToGroups[T any, K comparable, V any](c *Collection[T], callback func(T) (K, V)) map[K][]V {
	return MapToDictionary(c, callback)
}

// MapWithKeys maps each item to a key-value pair, returning a map.
func MapWithKeys[T any, K comparable, V any](c *Collection[T], callback func(T) (K, V)) map[K]V {
	result := make(map[K]V)

	for _, item := range c.items {
		key, value := callback(item)
		result[key] = value
	}

	return result
}
