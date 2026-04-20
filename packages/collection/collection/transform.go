package collection

import "cmp"

// Filter returns a new collection containing only items for which the callback returns true.
func (c *Collection[T]) Filter(callback func(T, int) bool) *Collection[T] {
	result := make([]T, 0)

	for i, item := range c.items {
		if callback(item, i) {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// Reject returns a new collection containing only items for which the callback returns false.
func (c *Collection[T]) Reject(callback func(T, int) bool) *Collection[T] {
	return c.Filter(func(item T, index int) bool {
		return !callback(item, index)
	})
}

// Map applies the callback to each item and returns a new collection of results.
func Map[T any, R any](c *Collection[T], callback func(T, int) R) *Collection[R] {
	result := make([]R, len(c.items))

	for i, item := range c.items {
		result[i] = callback(item, i)
	}

	return Collect(result)
}

// FlatMap applies the callback to each item, flattening the resulting slices into a single collection.
func FlatMap[T any, R any](c *Collection[T], callback func(T, int) []R) *Collection[R] {
	result := make([]R, 0)

	for i, item := range c.items {
		result = append(result, callback(item, i)...)
	}

	return Collect(result)
}

// MapInto applies the constructor to each item, returning a new collection of the mapped types.
func MapInto[T any, R any](c *Collection[T], constructor func(T) R) *Collection[R] {
	result := make([]R, len(c.items))

	for i, item := range c.items {
		result[i] = constructor(item)
	}

	return Collect(result)
}

// Reduce iterating over the collection and accumulates a single result using the callback.
func Reduce[T any, R any](c *Collection[T], callback func(R, T, int) R, initial R) R {
	result := initial

	for i, item := range c.items {
		result = callback(result, item, i)
	}

	return result
}

// Flatten returns a shallow copy of the collection.
// For non-nested typed slices this returns the items as-is.
func (c *Collection[T]) Flatten() *Collection[T] {
	return Collect(append([]T{}, c.items...))
}

// Reverse returns a new collection with items in reverse order.
func (c *Collection[T]) Reverse() *Collection[T] {
	result := make([]T, len(c.items))

	for i, j := 0, len(c.items)-1; j >= 0; i, j = i+1, j-1 {
		result[i] = c.items[j]
	}

	return Collect(result)
}

// Flip returns a new collection with the item order reversed.
// For typed Go slices this reverses the element order.
func (c *Collection[T]) Flip() *Collection[T] {
	result := make([]T, len(c.items))

	for i, j := 0, len(c.items)-1; j >= 0; i, j = i+1, j-1 {
		result[i] = c.items[j]
	}

	return Collect(result)
}

// Multiply returns a new collection with all items repeated the given number of times.
func (c *Collection[T]) Multiply(multiplier int) *Collection[T] {
	if multiplier <= 0 {
		return Empty[T]()
	}

	result := make([]T, 0, len(c.items)*multiplier)

	for i := 0; i < multiplier; i++ {
		result = append(result, c.items...)
	}

	return Collect(result)
}

// Values returns a new collection with re-indexed items (a shallow copy).
func (c *Collection[T]) Values() *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)

	return Collect(result)
}

// Keys returns a new Collection[int] containing the indices 0 through n-1.
func (c *Collection[T]) Keys() *Collection[int] {
	keys := make([]int, len(c.items))

	for i := range c.items {
		keys[i] = i
	}

	return Collect(keys)
}

// Unique returns a new collection containing only items with distinct keys
// as determined by the given key function.
func Unique[T any, K comparable](c *Collection[T], keyFunc func(T) K) *Collection[T] {
	seen := make(map[K]bool)
	result := make([]T, 0)

	for _, item := range c.items {
		key := keyFunc(item)

		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}

	return Collect(result)
}

// Duplicates returns a new collection containing all duplicate items
// as determined by the given key function.
func Duplicates[T any, K comparable](c *Collection[T], keyFunc func(T) K) *Collection[T] {
	seen := make(map[K]bool)
	result := make([]T, 0)

	for _, item := range c.items {
		key := keyFunc(item)

		if seen[key] {
			result = append(result, item)
		} else {
			seen[key] = true
		}
	}

	return Collect(result)
}

// Pluck extracts a value from each item using the given function, returning a new collection.
func Pluck[T any, V any](c *Collection[T], valueFunc func(T) V) *Collection[V] {
	result := make([]V, len(c.items))

	for i, item := range c.items {
		result[i] = valueFunc(item)
	}

	return Collect(result)
}

// Every report whether all items in the collection satisfy the given predicate.
func (c *Collection[T]) Every(callback func(T, int) bool) bool {
	for i, item := range c.items {
		if !callback(item, i) {
			return false
		}
	}

	return true
}

// Where filters items using a predicate that receives only the item (no index).
func (c *Collection[T]) Where(predicate func(T) bool) *Collection[T] {
	return c.Filter(func(item T, _ int) bool {
		return predicate(item)
	})
}

// WhereNot filters items using a negative predicate that receives only the item (no index).
func (c *Collection[T]) WhereNot(predicate func(T) bool) *Collection[T] {
	return c.Filter(func(item T, _ int) bool {
		return !predicate(item)
	})
}

// WhereNull returns items whose extracted value equals the zero value.
func WhereNull[T any, K comparable](c *Collection[T], keyFunc func(T) K) *Collection[T] {
	var zero K
	result := make([]T, 0)

	for _, item := range c.items {
		if keyFunc(item) == zero {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// WhereNotNull returns items whose extracted value is not the zero value.
func WhereNotNull[T any, K comparable](c *Collection[T], keyFunc func(T) K) *Collection[T] {
	var zero K
	result := make([]T, 0)

	for _, item := range c.items {
		if keyFunc(item) != zero {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// WhereIn returns items whose extracted key value is in the given set.
func WhereIn[T any, K comparable](c *Collection[T], keyFunc func(T) K, values []K) *Collection[T] {
	set := make(map[K]bool, len(values))

	for _, v := range values {
		set[v] = true
	}

	result := make([]T, 0)

	for _, item := range c.items {
		if set[keyFunc(item)] {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// WhereNotIn returns items whose extracted key value is not in the given set.
func WhereNotIn[T any, K comparable](c *Collection[T], keyFunc func(T) K, values []K) *Collection[T] {
	set := make(map[K]bool, len(values))

	for _, v := range values {
		set[v] = true
	}

	result := make([]T, 0)

	for _, item := range c.items {
		if !set[keyFunc(item)] {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// WhereBetween returns items whose extracted key value is between min and max (inclusive).
func WhereBetween[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K, min, max K) *Collection[T] {
	result := make([]T, 0)

	for _, item := range c.items {
		v := keyFunc(item)

		if v >= min && v <= max {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// WhereNotBetween returns items whose extracted key value is outside the range [min, max].
func WhereNotBetween[T any, K cmp.Ordered](c *Collection[T], keyFunc func(T) K, min, max K) *Collection[T] {
	result := make([]T, 0)

	for _, item := range c.items {
		v := keyFunc(item)

		if v < min || v > max {
			result = append(result, item)
		}
	}

	return Collect(result)
}
