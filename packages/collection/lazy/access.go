package lazy

import "github.com/bedrock/packages/collection/support"

// First returns the first element matching the optional predicate.
// If no predicate is given, it returns the first element.
// The second return value indicates whether a matching element was found.
func (lc *Collection[T]) First(predicates ...func(T, int) bool) (T, bool) {
	var result T
	found := false
	idx := 0

	if len(predicates) == 0 || predicates[0] == nil {
		lc.source(func(item T) bool {
			result = item
			found = true

			return false
		})
	} else {
		predicate := predicates[0]
		lc.source(func(item T) bool {
			if predicate(item, idx) {
				result = item
				found = true

				return false
			}

			idx++

			return true
		})
	}

	return result, found
}

// FirstOrFail returns the first element matching the optional predicate,
// or an ItemNotFoundError if no element is found.
func (lc *Collection[T]) FirstOrFail(predicates ...func(T, int) bool) (T, error) {
	item, ok := lc.First(predicates...)

	if !ok {
		var zero T

		return zero, &support.ItemNotFoundError{}
	}

	return item, nil
}

// Last returns the last element matching the optional predicate.
// The second return value indicates whether a matching element was found.
func (lc *Collection[T]) Last(predicates ...func(T, int) bool) (T, bool) {
	var result T
	found := false
	idx := 0

	if len(predicates) == 0 || predicates[0] == nil {
		lc.source(func(item T) bool {
			result = item
			found = true

			return true
		})
	} else {
		predicate := predicates[0]
		lc.source(func(item T) bool {
			if predicate(item, idx) {
				result = item
				found = true
			}

			idx++

			return true
		})
	}

	return result, found
}

// Sole returns the only element matching the optional predicate.
// It returns an error if zero or more than one element matches.
func (lc *Collection[T]) Sole(predicates ...func(T, int) bool) (T, error) {
	var result T
	count := 0
	idx := 0

	if len(predicates) == 0 || predicates[0] == nil {
		lc.source(func(item T) bool {
			result = item
			count++

			return count < 2
		})
	} else {
		predicate := predicates[0]
		lc.source(func(item T) bool {
			if predicate(item, idx) {
				result = item
				count++

				if count > 1 {
					return false
				}
			}

			idx++

			return true
		})
	}

	if count == 0 {
		var zero T

		return zero, &support.ItemNotFoundError{}
	}

	if count > 1 {
		var zero T

		return zero, &support.MultipleItemsFoundError{Count: count}
	}

	return result, nil
}

// Get returns the item at the given zero-based index.
// The second return value indicates whether the index exists.
func (lc *Collection[T]) Get(index int) (T, bool) {
	var result T
	found := false
	idx := 0
	lc.source(func(item T) bool {
		if idx == index {
			result = item
			found = true

			return false
		}

		idx++

		return true
	})

	return result, found
}

// Contains reports whether any item satisfies the predicate.
func (lc *Collection[T]) Contains(predicate func(T, int) bool) bool {
	found := false
	idx := 0
	lc.source(func(item T) bool {
		if predicate(item, idx) {
			found = true

			return false
		}

		idx++

		return true
	})

	return found
}

// Some is an alias for Contains.
func (lc *Collection[T]) Some(predicate func(T, int) bool) bool {
	return lc.Contains(predicate)
}

// DoesntContain reports whether no item satisfies the predicate.
func (lc *Collection[T]) DoesntContain(predicate func(T, int) bool) bool {
	return !lc.Contains(predicate)
}

// Search returns the index of the first item satisfying the predicate.
// The second return value indicates whether a match was found.
func (lc *Collection[T]) Search(predicate func(T, int) bool) (int, bool) {
	resultIdx := -1
	found := false
	idx := 0
	lc.source(func(item T) bool {
		if predicate(item, idx) {
			resultIdx = idx
			found = true

			return false
		}

		idx++

		return true
	})

	return resultIdx, found
}

// Before returns the item immediately before the first item matching the predicate.
// The second return value indicates whether such an item exists.
func (lc *Collection[T]) Before(predicate func(T, int) bool) (T, bool) {
	var prev T
	hasPrev := false
	found := false
	idx := 0
	lc.source(func(item T) bool {
		if predicate(item, idx) {
			found = true

			return false
		}

		prev = item
		hasPrev = true
		idx++

		return true
	})

	if found && hasPrev {
		return prev, true
	}

	var zero T

	return zero, false
}

// After returns the item immediately after the first item matching the predicate.
// The second return value indicates whether such an item exists.
func (lc *Collection[T]) After(predicate func(T, int) bool) (T, bool) {
	matched := false

	var result T
	found := false
	idx := 0
	lc.source(func(item T) bool {
		if matched {
			result = item
			found = true

			return false
		}

		if predicate(item, idx) {
			matched = true
		}

		idx++

		return true
	})

	return result, found
}

// Has reports whether the given zero-based index exists in the collection.
func (lc *Collection[T]) Has(index int) bool {
	_, ok := lc.Get(index)

	return ok
}

// HasAny reports whether any of the given indices exist in the collection.
func (lc *Collection[T]) HasAny(indices ...int) bool {
	for _, idx := range indices {
		if lc.Has(idx) {
			return true
		}
	}

	return false
}

// HasSole reports whether exactly one item matches the optional predicate.
// If no predicate is given, it checks whether the collection has exactly one item.
func (lc *Collection[T]) HasSole(predicates ...func(T, int) bool) bool {
	count := 0
	idx := 0

	if len(predicates) == 0 || predicates[0] == nil {
		lc.source(func(_ T) bool {
			count++

			return count < 2
		})
	} else {
		predicate := predicates[0]
		lc.source(func(item T) bool {
			if predicate(item, idx) {
				count++

				if count > 1 {
					return false
				}
			}

			idx++

			return true
		})
	}

	return count == 1
}
