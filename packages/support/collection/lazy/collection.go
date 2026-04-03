package lazy

import "iter"

// Collection provides lazy evaluation for collection operations.
// Items are only computed when needed, making it memory-efficient for large datasets.
type Collection[T any] struct {
	source iter.Seq[T]
}

// New creates a new Collection from an iterator function.
// The source receives a yield function; call yield for each item and
// stop generating if the yield returns false.
func New[T any](source iter.Seq[T]) *Collection[T] {
	return &Collection[T]{source: source}
}

// From creates a Collection from a slice.
func From[T any](items []T) *Collection[T] {
	return New(func(yield func(T) bool) {
		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	})
}

// Empty creates an empty Collection.
func Empty[T any]() *Collection[T] {
	return New[T](func(yield func(T) bool) {})
}

// Range creates a lazy collection of sequential integers from start to end (inclusive).
// If from > to, the sequence counts downward.
func Range(from, to int) *Collection[int] {
	return New(func(yield func(int) bool) {
		if from <= to {
			for i := from; i <= to; i++ {
				if !yield(i) {
					return
				}
			}
		} else {
			for i := from; i >= to; i-- {
				if !yield(i) {
					return
				}
			}
		}
	})
}

// Times create a lazy collection by invoking the callback n times.
// The callback receives a 1-based index.
func Times[T any](number int, callback func(int) T) *Collection[T] {
	return New(func(yield func(T) bool) {
		for i := 1; i <= number; i++ {
			if !yield(callback(i)) {
				return
			}
		}
	})
}

// Iter returns the underlying iterator for use with range loops.
func (lc *Collection[T]) Iter() iter.Seq[T] {
	return lc.source
}

// All eagerly evaluates the lazy collection and returns all items as a slice.
func (lc *Collection[T]) All() []T {
	result := make([]T, 0)
	lc.source(func(item T) bool {
		result = append(result, item)

		return true
	})

	return result
}

// Eager eagerly evaluates the lazy collection and returns all items as a slice.
func (lc *Collection[T]) Eager() []T {
	return lc.All()
}

// Collect converts the lazy collection to a slice.
func (lc *Collection[T]) Collect() []T {
	return lc.All()
}

// Count returns the total number of items by eagerly evaluating the collection.
func (lc *Collection[T]) Count() int {
	count := 0
	lc.source(func(_ T) bool {
		count++

		return true
	})

	return count
}

// IsEmpty reports whether the lazy collection contains no items.
func (lc *Collection[T]) IsEmpty() bool {
	empty := true
	lc.source(func(_ T) bool {
		empty = false

		return false
	})

	return empty
}

// IsNotEmpty reports whether the lazy collection contains at least one item.
func (lc *Collection[T]) IsNotEmpty() bool {
	return !lc.IsEmpty()
}

// ContainsOneItem reports whether the lazy collection contains exactly one item.
func (lc *Collection[T]) ContainsOneItem() bool {
	return lc.Count() == 1
}

// ContainsManyItems reports whether the lazy collection contains more than one item.
func (lc *Collection[T]) ContainsManyItems() bool {
	return lc.Count() > 1
}
