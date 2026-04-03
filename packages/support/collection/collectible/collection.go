package collectible

import (
	"iter"

	"github.com/gollin/packages/support/collection/support"
)

// Collection wraps a map and provides a fluent API for working with
// key-value pairs while maintaining insertion order.
type Collection[K comparable, V any] struct {
	items map[K]V
	keys  []K // maintains insertion order
}

// New creates a new Collection from a map.
func New[K comparable, V any](items map[K]V) *Collection[K, V] {
	if items == nil {
		items = make(map[K]V)
	}

	keys := make([]K, 0, len(items))

	for k := range items {
		keys = append(keys, k)
	}

	return &Collection[K, V]{items: items, keys: keys}
}

// FromPairs creates a Collection from key-value pairs.
func FromPairs[K comparable, V any](pairs ...support.Pair[K, V]) *Collection[K, V] {
	items := make(map[K]V, len(pairs))
	keys := make([]K, 0, len(pairs))

	for _, p := range pairs {
		if _, exists := items[p.Key]; !exists {
			keys = append(keys, p.Key)
		}

		items[p.Key] = p.Value
	}

	return &Collection[K, V]{items: items, keys: keys}
}

// All returns a shallow copy of the underlying map.
func (m *Collection[K, V]) All() map[K]V {
	result := make(map[K]V, len(m.items))

	for k, v := range m.items {
		result[k] = v
	}

	return result
}

// Keys returns all keys as a slice, preserving insertion order.
func (m *Collection[K, V]) Keys() []K {
	result := make([]K, len(m.keys))
	copy(result, m.keys)

	return result
}

// Values returns all values as a slice, preserving insertion order.
func (m *Collection[K, V]) Values() []V {
	result := make([]V, 0, len(m.keys))

	for _, k := range m.keys {
		result = append(result, m.items[k])
	}

	return result
}

// Count returns the number of items in the collection.
func (m *Collection[K, V]) Count() int {
	return len(m.items)
}

// IsEmpty reports whether the collection contains no items.
func (m *Collection[K, V]) IsEmpty() bool {
	return len(m.items) == 0
}

// IsNotEmpty reports whether the collection contains at least one item.
func (m *Collection[K, V]) IsNotEmpty() bool {
	return len(m.items) > 0
}

// ContainsOneItem reports whether the collection contains exactly one item.
func (m *Collection[K, V]) ContainsOneItem() bool {
	return len(m.items) == 1
}

// ContainsManyItems reports whether the collection contains more than one item.
func (m *Collection[K, V]) ContainsManyItems() bool {
	return len(m.items) > 1
}

// Iter returns an iterator over the collection's key-value pairs in insertion
// order.
func (m *Collection[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, k := range m.keys {
			if !yield(k, m.items[k]) {
				return
			}
		}
	}
}
