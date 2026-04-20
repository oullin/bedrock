package collectible

import "sort"

// SortKeys returns a new collection with string keys sorted in ascending
// order.
func SortKeys[V any](m *Collection[string, V]) *Collection[string, V] {
	result := make(map[string]V, len(m.items))

	for k, v := range m.items {
		result[k] = v
	}

	newKeys := make([]string, len(m.keys))
	copy(newKeys, m.keys)

	sort.Strings(newKeys)

	return &Collection[string, V]{items: result, keys: newKeys}
}

// SortKeysDesc returns a new collection with string keys sorted in
// descending order.
func SortKeysDesc[V any](m *Collection[string, V]) *Collection[string, V] {
	result := make(map[string]V, len(m.items))

	for k, v := range m.items {
		result[k] = v
	}

	newKeys := make([]string, len(m.keys))
	copy(newKeys, m.keys)

	sort.Sort(sort.Reverse(sort.StringSlice(newKeys)))

	return &Collection[string, V]{items: result, keys: newKeys}
}

// SortKeysUsing returns a new collection with keys sorted using the provided
// comparison function.
func (m *Collection[K, V]) SortKeysUsing(less func(K, K) bool) *Collection[K, V] {
	result := make(map[K]V, len(m.items))

	for k, v := range m.items {
		result[k] = v
	}

	newKeys := make([]K, len(m.keys))
	copy(newKeys, m.keys)

	sort.SliceStable(newKeys, func(i, j int) bool {
		return less(newKeys[i], newKeys[j])
	})

	return &Collection[K, V]{items: result, keys: newKeys}
}
