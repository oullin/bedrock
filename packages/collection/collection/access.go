package collection

import "github.com/bedrock/packages/collection/support"

// First returns the first element matching the optional predicate.
// If no predicate is provided, the first element is returned.
// The second return value indicates whether a match was found.
func (c *Collection[T]) First(predicates ...func(T, int) bool) (T, bool) {
	if len(c.items) == 0 {
		var zero T

		return zero, false
	}

	if len(predicates) == 0 || predicates[0] == nil {
		return c.items[0], true
	}

	predicate := predicates[0]

	for i, item := range c.items {
		if predicate(item, i) {
			return item, true
		}
	}

	var zero T

	return zero, false
}

// FirstOrFail returns the first element matching the optional predicate,
// or an ItemNotFoundError if no match is found.
func (c *Collection[T]) FirstOrFail(predicates ...func(T, int) bool) (T, error) {
	item, ok := c.First(predicates...)

	if !ok {
		var zero T

		return zero, &support.ItemNotFoundError{}
	}

	return item, nil
}

// Last returns the last element matching the optional predicate.
// If no predicate is provided, the last element is returned.
// The second return value indicates whether a match was found.
func (c *Collection[T]) Last(predicates ...func(T, int) bool) (T, bool) {
	if len(c.items) == 0 {
		var zero T

		return zero, false
	}

	if len(predicates) == 0 || predicates[0] == nil {
		return c.items[len(c.items)-1], true
	}

	predicate := predicates[0]

	for i := len(c.items) - 1; i >= 0; i-- {
		if predicate(c.items[i], i) {
			return c.items[i], true
		}
	}

	var zero T

	return zero, false
}

// Sole returns the only element matching the optional predicate.
// It returns an ItemNotFoundError if no items match, or a MultipleItemsFoundError
// if more than one item matches.
func (c *Collection[T]) Sole(predicates ...func(T, int) bool) (T, error) {
	var filtered *Collection[T]

	if len(predicates) == 0 || predicates[0] == nil {
		filtered = c
	} else {
		filtered = c.Filter(predicates[0])
	}

	if filtered.Count() == 0 {
		var zero T

		return zero, &support.ItemNotFoundError{}
	}

	if filtered.Count() > 1 {
		var zero T

		return zero, &support.MultipleItemsFoundError{Count: filtered.Count()}
	}

	return filtered.items[0], nil
}

// HasSole reports whether exactly one item matches the optional predicate.
func (c *Collection[T]) HasSole(predicates ...func(T, int) bool) bool {
	var filtered *Collection[T]

	if len(predicates) == 0 || predicates[0] == nil {
		filtered = c
	} else {
		filtered = c.Filter(predicates[0])
	}

	return filtered.Count() == 1
}

// Get returns the item at the given index.
// Negative indices count from the end. The second return value indicates
// whether the index was within bounds. An optional default may be provided.
func (c *Collection[T]) Get(index int, defaults ...T) (T, bool) {
	if index < 0 {
		index = len(c.items) + index
	}

	if index >= 0 && index < len(c.items) {
		return c.items[index], true
	}

	if len(defaults) > 0 {
		return defaults[0], false
	}

	var zero T

	return zero, false
}

// GetOrPut returns the item at the given index if it exists.
// Otherwise, it appends the value to the collection and returns it.
func (c *Collection[T]) GetOrPut(index int, value T) T {
	if index >= 0 && index < len(c.items) {
		return c.items[index]
	}

	c.items = append(c.items, value)

	return value
}

// Contains reports whether any item in the collection satisfies the predicate.
func (c *Collection[T]) Contains(predicate func(T, int) bool) bool {
	for i, item := range c.items {
		if predicate(item, i) {
			return true
		}
	}

	return false
}

// Some is an alias for Contains.
func (c *Collection[T]) Some(predicate func(T, int) bool) bool {
	return c.Contains(predicate)
}

// DoesntContain reports whether no item in the collection satisfies the predicate.
func (c *Collection[T]) DoesntContain(predicate func(T, int) bool) bool {
	return !c.Contains(predicate)
}

// Search returns the index of the first item satisfying the predicate.
// The second return value indicates whether a match was found.
func (c *Collection[T]) Search(predicate func(T, int) bool) (int, bool) {
	for i, item := range c.items {
		if predicate(item, i) {
			return i, true
		}
	}

	return -1, false
}

// Before returns the item immediately before the first item matching the predicate.
func (c *Collection[T]) Before(predicate func(T, int) bool) (T, bool) {
	for i, item := range c.items {
		if predicate(item, i) {
			if i == 0 {
				var zero T

				return zero, false
			}

			return c.items[i-1], true
		}
	}

	var zero T

	return zero, false
}

// After returns the item immediately after the first item matching the predicate.
func (c *Collection[T]) After(predicate func(T, int) bool) (T, bool) {
	for i, item := range c.items {
		if predicate(item, i) {
			if i >= len(c.items)-1 {
				var zero T

				return zero, false
			}

			return c.items[i+1], true
		}
	}

	var zero T

	return zero, false
}

// Has reports whether the given index exists in the collection.
// Negative indices count from the end.
func (c *Collection[T]) Has(index int) bool {
	if index < 0 {
		index = len(c.items) + index
	}

	return index >= 0 && index < len(c.items)
}

// HasAny reports whether any of the given indices exist in the collection.
func (c *Collection[T]) HasAny(indices ...int) bool {
	for _, idx := range indices {
		if c.Has(idx) {
			return true
		}
	}

	return false
}

// Only returns a new collection containing only items at the given indices.
func (c *Collection[T]) Only(indices ...int) *Collection[T] {
	result := make([]T, 0, len(indices))

	for _, idx := range indices {
		if idx >= 0 && idx < len(c.items) {
			result = append(result, c.items[idx])
		}
	}

	return Collect(result)
}

// Except returns a new collection excluding items at the given indices.
func (c *Collection[T]) Except(indices ...int) *Collection[T] {
	excludeSet := make(map[int]bool, len(indices))

	for _, idx := range indices {
		excludeSet[idx] = true
	}

	result := make([]T, 0)

	for i, item := range c.items {
		if !excludeSet[i] {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// Dot returns a shallow copy of the collection.
// For typed Go slices, no dot-notation expansion is possible.
func (c *Collection[T]) Dot() *Collection[T] {
	return c.Copy()
}

// Undot returns a shallow copy of the collection.
// For typed Go slices, no dot-notation expansion is possible.
func (c *Collection[T]) Undot() *Collection[T] {
	return c.Copy()
}
