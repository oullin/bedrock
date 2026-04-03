package lazy

import "time"

// Take returns a new lazy collection with at most limit items.
// A negative limit takes from the end (requires eager evaluation).
func (lc *Collection[T]) Take(limit int) *Collection[T] {
	if limit < 0 {
		// For negative take, we need to eagerly evaluate
		items := lc.All()
		start := len(items) + limit

		if start < 0 {
			start = 0
		}

		return From(items[start:])
	}

	return New(func(yield func(T) bool) {
		count := 0
		lc.source(func(item T) bool {
			if count >= limit {
				return false
			}

			count++

			return yield(item)
		})
	})
}

// TakeUntil returns items from the beginning until the callback returns true.
func (lc *Collection[T]) TakeUntil(callback func(T, int) bool) *Collection[T] {
	return New(func(yield func(T) bool) {
		idx := 0
		lc.source(func(item T) bool {
			if callback(item, idx) {
				return false
			}

			idx++

			return yield(item)
		})
	})
}

// TakeWhile returns items from the beginning while the callback returns true.
func (lc *Collection[T]) TakeWhile(callback func(T, int) bool) *Collection[T] {
	return lc.TakeUntil(func(item T, idx int) bool {
		return !callback(item, idx)
	})
}

// TakeUntilTimeout returns items until the given duration has elapsed.
func (lc *Collection[T]) TakeUntilTimeout(timeout time.Duration) *Collection[T] {
	return New(func(yield func(T) bool) {
		deadline := time.Now().Add(timeout)
		lc.source(func(item T) bool {
			if time.Now().After(deadline) {
				return false
			}

			return yield(item)
		})
	})
}

// Skip returns a new lazy collection that skips the first count items.
func (lc *Collection[T]) Skip(count int) *Collection[T] {
	return New(func(yield func(T) bool) {
		skipped := 0
		lc.source(func(item T) bool {
			if skipped < count {
				skipped++

				return true
			}

			return yield(item)
		})
	})
}

// SkipUntil skips items until the callback returns true, then yields the rest.
func (lc *Collection[T]) SkipUntil(callback func(T, int) bool) *Collection[T] {
	return New(func(yield func(T) bool) {
		found := false
		idx := 0
		lc.source(func(item T) bool {
			if !found && callback(item, idx) {
				found = true
			}

			if found {
				return yield(item)
			}

			idx++

			return true
		})
	})
}

// SkipWhile skips items while the callback returns true, then yields the rest.
func (lc *Collection[T]) SkipWhile(callback func(T, int) bool) *Collection[T] {
	return lc.SkipUntil(func(item T, idx int) bool {
		return !callback(item, idx)
	})
}

// Slice returns a subset of the lazy collection starting at offset with an optional length.
func (lc *Collection[T]) Slice(offset int, lengths ...int) *Collection[T] {
	result := lc.Skip(offset)

	if len(lengths) > 0 {
		result = result.Take(lengths[0])
	}

	return result
}

// Chunk eagerly evaluates the lazy collection and splits items into groups of the given size.
func (lc *Collection[T]) Chunk(size int) [][]T {
	var result [][]T
	chunk := make([]T, 0, size)
	lc.source(func(item T) bool {
		chunk = append(chunk, item)

		if len(chunk) >= size {
			result = append(result, chunk)
			chunk = make([]T, 0, size)
		}

		return true
	})

	if len(chunk) > 0 {
		result = append(result, chunk)
	}

	return result
}

// ChunkWhile splits the lazy collection into groups where consecutive items satisfy the callback.
// A new chunk is started when the callback returns false.
func (lc *Collection[T]) ChunkWhile(callback func(T, int, []T) bool) [][]T {
	var result [][]T
	current := make([]T, 0)
	idx := 0
	lc.source(func(item T) bool {
		if len(current) == 0 || callback(item, idx, current) {
			current = append(current, item)
		} else {
			result = append(result, current)
			current = []T{item}
		}

		idx++

		return true
	})

	if len(current) > 0 {
		result = append(result, current)
	}

	return result
}

// Nth returns every step-th element, starting from an optional offset.
func (lc *Collection[T]) Nth(step int, offsets ...int) *Collection[T] {
	offset := 0

	if len(offsets) > 0 {
		offset = offsets[0]
	}

	return New(func(yield func(T) bool) {
		idx := 0
		lc.source(func(item T) bool {
			if idx >= offset && (idx-offset)%step == 0 {
				if !yield(item) {
					return false
				}
			}

			idx++

			return true
		})
	})
}

// Concat returns a new lazy collection with the given items appended.
func (lc *Collection[T]) Concat(items []T) *Collection[T] {
	return New(func(yield func(T) bool) {
		lc.source(func(item T) bool {
			return yield(item)
		})

		for _, item := range items {
			if !yield(item) {
				return
			}
		}
	})
}

// Pad returns a new lazy collection padded to the given size with the specified value.
// A negative size pads at the beginning; a positive size pads at the end.
func (lc *Collection[T]) Pad(size int, value T) *Collection[T] {
	return New(func(yield func(T) bool) {
		count := 0
		absSize := size

		if absSize < 0 {
			absSize = -absSize
		}

		if size < 0 {
			// Evaluate to know count
			items := lc.All()
			padCount := absSize - len(items)

			if padCount > 0 {
				for i := 0; i < padCount; i++ {
					if !yield(value) {
						return
					}
				}
			}

			for _, item := range items {
				if !yield(item) {
					return
				}
			}
		} else {
			lc.source(func(item T) bool {
				count++

				return yield(item)
			})

			for count < size {
				if !yield(value) {
					return
				}

				count++
			}
		}
	})
}
