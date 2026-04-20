package collection

import "math"

// Chunk breaks the collection into multiple slices of the given size.
func (c *Collection[T]) Chunk(size int) [][]T {
	if size <= 0 {
		return nil
	}

	chunks := make([][]T, 0)

	for i := 0; i < len(c.items); i += size {
		end := i + size

		if end > len(c.items) {
			end = len(c.items)
		}

		chunk := make([]T, end-i)
		copy(chunk, c.items[i:end])
		chunks = append(chunks, chunk)
	}

	return chunks
}

// ChunkWhile breaks the collection into groups as long as the callback returns true.
// A new group is started each time the callback returns false.
func (c *Collection[T]) ChunkWhile(callback func(T, int, []T) bool) [][]T {
	if len(c.items) == 0 {
		return nil
	}

	chunks := make([][]T, 0)
	current := []T{c.items[0]}

	for i := 1; i < len(c.items); i++ {
		if callback(c.items[i], i, current) {
			current = append(current, c.items[i])
		} else {
			chunks = append(chunks, current)
			current = []T{c.items[i]}
		}
	}

	chunks = append(chunks, current)

	return chunks
}

// Split breaks the collection into the given number of groups.
func (c *Collection[T]) Split(numberOfGroups int) [][]T {
	if len(c.items) == 0 || numberOfGroups <= 0 {
		return nil
	}

	groups := make([][]T, 0, numberOfGroups)
	groupSize := float64(len(c.items)) / float64(numberOfGroups)

	for i := 0; i < numberOfGroups; i++ {
		start := int(math.Round(float64(i) * groupSize))
		end := int(math.Round(float64(i+1) * groupSize))

		if start >= len(c.items) {
			break
		}

		if end > len(c.items) {
			end = len(c.items)
		}

		chunk := make([]T, end-start)
		copy(chunk, c.items[start:end])
		groups = append(groups, chunk)
	}

	return groups
}

// SplitIn splits the collection into groups, filling non-terminal groups.
func (c *Collection[T]) SplitIn(numberOfGroups int) [][]T {
	size := int(math.Ceil(float64(len(c.items)) / float64(numberOfGroups)))

	return c.Chunk(size)
}

// Sliding returns a sliding window view of the collection with the given window size and step.
func (c *Collection[T]) Sliding(size int, steps ...int) [][]T {
	step := 1

	if len(steps) > 0 {
		step = steps[0]
	}

	if size <= 0 || step <= 0 || len(c.items) == 0 {
		return nil
	}

	chunks := make([][]T, 0)

	for i := 0; i+size <= len(c.items); i += step {
		chunk := make([]T, size)
		copy(chunk, c.items[i:i+size])
		chunks = append(chunks, chunk)
	}

	return chunks
}

// Slice extracts a portion of the collection starting at the given offset.
// An optional length limits how many items are returned.
func (c *Collection[T]) Slice(offset int, lengths ...int) *Collection[T] {
	if offset < 0 {
		offset = len(c.items) + offset

		if offset < 0 {
			offset = 0
		}
	}

	if offset >= len(c.items) {
		return Empty[T]()
	}

	if len(lengths) > 0 {
		end := offset + lengths[0]

		if end > len(c.items) {
			end = len(c.items)
		}

		result := make([]T, end-offset)
		copy(result, c.items[offset:end])

		return Collect(result)
	}

	result := make([]T, len(c.items)-offset)
	copy(result, c.items[offset:])

	return Collect(result)
}

// Splice removes and returns a slice of items starting at the given offset.
// An optional length limits the number of items removed.
func (c *Collection[T]) Splice(offset int, lengths ...int) *Collection[T] {
	length := len(c.items) - offset

	if len(lengths) > 0 {
		length = lengths[0]
	}

	if offset < 0 {
		offset = len(c.items) + offset

		if offset < 0 {
			offset = 0
		}
	}

	if offset >= len(c.items) {
		return Empty[T]()
	}

	end := offset + length

	if end > len(c.items) {
		end = len(c.items)
	}

	removed := make([]T, end-offset)
	copy(removed, c.items[offset:end])
	c.items = append(c.items[:offset], c.items[end:]...)

	return Collect(removed)
}

// SpliceReplace removes a portion at the given offset and replaces it with the provided items.
// It returns the removed items.
func (c *Collection[T]) SpliceReplace(offset, length int, replacement []T) *Collection[T] {
	if offset < 0 {
		offset = len(c.items) + offset

		if offset < 0 {
			offset = 0
		}
	}

	end := offset + length

	if end > len(c.items) {
		end = len(c.items)
	}

	removed := make([]T, end-offset)
	copy(removed, c.items[offset:end])
	newItems := make([]T, 0, len(c.items)-length+len(replacement))
	newItems = append(newItems, c.items[:offset]...)
	newItems = append(newItems, replacement...)
	newItems = append(newItems, c.items[end:]...)
	c.items = newItems

	return Collect(removed)
}

// Take returns a new collection with the specified number of items from the front.
// A negative limit takes from the end.
func (c *Collection[T]) Take(limit int) *Collection[T] {
	if limit < 0 {
		return c.Slice(limit)
	}

	return c.Slice(0, limit)
}

// TakeUntil returns items from the start until the callback returns true.
func (c *Collection[T]) TakeUntil(callback func(T, int) bool) *Collection[T] {
	result := make([]T, 0)

	for i, item := range c.items {
		if callback(item, i) {
			break
		}

		result = append(result, item)
	}

	return Collect(result)
}

// TakeWhile returns items from the start as long as the callback returns true.
func (c *Collection[T]) TakeWhile(callback func(T, int) bool) *Collection[T] {
	return c.TakeUntil(func(item T, index int) bool {
		return !callback(item, index)
	})
}

// Skip returns a new collection with the first n items removed.
func (c *Collection[T]) Skip(count int) *Collection[T] {
	return c.Slice(count)
}

// SkipUntil skips items until the callback returns true, then returns the rest.
func (c *Collection[T]) SkipUntil(callback func(T, int) bool) *Collection[T] {
	result := make([]T, 0)
	found := false

	for i, item := range c.items {
		if !found && callback(item, i) {
			found = true
		}

		if found {
			result = append(result, item)
		}
	}

	return Collect(result)
}

// SkipWhile skips items as long as the callback returns true, then returns the rest.
func (c *Collection[T]) SkipWhile(callback func(T, int) bool) *Collection[T] {
	return c.SkipUntil(func(item T, index int) bool {
		return !callback(item, index)
	})
}

// Nth returns a new collection containing every n-th element, starting at an optional offset.
func (c *Collection[T]) Nth(step int, offsets ...int) *Collection[T] {
	offset := 0

	if len(offsets) > 0 {
		offset = offsets[0]
	}

	result := make([]T, 0)

	for i := offset; i < len(c.items); i += step {
		result = append(result, c.items[i])
	}

	return Collect(result)
}

// ForPage returns a subset of items for the given page number and page size.
func (c *Collection[T]) ForPage(page, perPage int) *Collection[T] {
	offset := (page - 1) * perPage

	return c.Slice(offset, perPage)
}

// Partition splits the collection into two: items that pass the predicate and items that do not.
func (c *Collection[T]) Partition(callback func(T, int) bool) (*Collection[T], *Collection[T]) {
	pass := make([]T, 0)
	fail := make([]T, 0)

	for i, item := range c.items {
		if callback(item, i) {
			pass = append(pass, item)
		} else {
			fail = append(fail, item)
		}
	}

	return Collect(pass), Collect(fail)
}

// Concat returns a new collection with the given items appended.
func (c *Collection[T]) Concat(items []T) *Collection[T] {
	result := make([]T, len(c.items)+len(items))
	copy(result, c.items)
	copy(result[len(c.items):], items)

	return Collect(result)
}

// Merge returns a new collection with the given items merged in (appended).
func (c *Collection[T]) Merge(items []T) *Collection[T] {
	return c.Concat(items)
}

// Pad pads the collection to the specified length with the given value.
// A negative size pad on the left; a positive size pad on the right.
func (c *Collection[T]) Pad(size int, value T) *Collection[T] {
	absSize := size

	if absSize < 0 {
		absSize = -absSize
	}

	if len(c.items) >= absSize {
		result := make([]T, len(c.items))
		copy(result, c.items)

		return Collect(result)
	}

	padCount := absSize - len(c.items)
	padding := make([]T, padCount)

	for i := range padding {
		padding[i] = value
	}

	if size < 0 {
		result := make([]T, absSize)
		copy(result, padding)
		copy(result[padCount:], c.items)

		return Collect(result)
	}

	result := make([]T, absSize)
	copy(result, c.items)
	copy(result[len(c.items):], padding)

	return Collect(result)
}
