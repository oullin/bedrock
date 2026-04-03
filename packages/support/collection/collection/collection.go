// Package collection provides a fluent, generic wrapper for working with slices of data.
package collection

import "iter"

// Collection wraps a slice and provides a fluent API for working with arrays of data.
type Collection[T any] struct {
	items []T
}

// New creates a new Collection from the given items.
func New[T any](items ...T) *Collection[T] {
	if items == nil {
		items = make([]T, 0)
	}

	return &Collection[T]{items: items}
}

// Collect creates a new Collection from a slice.
func Collect[T any](items []T) *Collection[T] {
	if items == nil {
		items = make([]T, 0)
	}

	return &Collection[T]{items: items}
}

// Empty creates an empty Collection.
func Empty[T any]() *Collection[T] {
	return &Collection[T]{items: make([]T, 0)}
}

// Wrap wraps the given value in a collection if it is not already one.
func Wrap[T any](value any) *Collection[T] {
	switch v := value.(type) {
	case *Collection[T]:
		return v
	case []T:
		return Collect(v)
	default:
		if item, ok := value.(T); ok {
			return New(item)
		}

		return Empty[T]()
	}
}

// Unwrap returns the underlying items from a Collection, or the value itself if it is a slice.
func Unwrap[T any](value any) []T {
	if c, ok := value.(*Collection[T]); ok {
		return c.All()
	}

	if items, ok := value.([]T); ok {
		return items
	}

	return nil
}

// Times create a new collection by invoking the callback a given number of times.
func Times[T any](number int, callback func(int) T) *Collection[T] {
	if number < 1 {
		return Empty[T]()
	}

	items := make([]T, number)

	for i := 0; i < number; i++ {
		items[i] = callback(i + 1)
	}

	return Collect(items)
}

// Range creates a collection of consecutive integers from start to end (inclusive).
func Range(from, to int) *Collection[int] {
	if from > to {
		items := make([]int, 0, from-to+1)

		for i := from; i >= to; i-- {
			items = append(items, i)
		}

		return Collect(items)
	}

	items := make([]int, 0, to-from+1)

	for i := from; i <= to; i++ {
		items = append(items, i)
	}

	return Collect(items)
}

// All returns all items in the collection as a slice.
func (c *Collection[T]) All() []T {
	return c.items
}

// Count returns the total number of items in the collection.
func (c *Collection[T]) Count() int {
	return len(c.items)
}

// IsEmpty reports whether the collection contains no items.
func (c *Collection[T]) IsEmpty() bool {
	return len(c.items) == 0
}

// IsNotEmpty reports whether the collection contains at least one item.
func (c *Collection[T]) IsNotEmpty() bool {
	return len(c.items) > 0
}

// ContainsOneItem reports whether the collection contains exactly one item.
func (c *Collection[T]) ContainsOneItem() bool {
	return len(c.items) == 1
}

// ContainsManyItems reports whether the collection contains more than one item.
func (c *Collection[T]) ContainsManyItems() bool {
	return len(c.items) > 1
}

// HasMany is an alias for ContainsManyItems.
func (c *Collection[T]) HasMany() bool {
	return c.ContainsManyItems()
}

// Len returns the number of items, implementing sort.Interface.
func (c *Collection[T]) Len() int {
	return len(c.items)
}

// ToBase returns the collection itself.
func (c *Collection[T]) ToBase() *Collection[T] {
	return c
}

// Iter returns an iter.Seq[T] that yields each item in the collection.
func (c *Collection[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range c.items {
			if !yield(item) {
				return
			}
		}
	}
}

// PairIter returns an iter.Seq2[int, T] that yields each index-item pair in the collection.
func (c *Collection[T]) PairIter() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i, item := range c.items {
			if !yield(i, item) {
				return
			}
		}
	}
}
