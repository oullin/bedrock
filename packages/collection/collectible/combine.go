package collectible

// Merge returns a new collection with the given map merged in. Existing keys
// are overwritten by the incoming values.
func (m *Collection[K, V]) Merge(items map[K]V) *Collection[K, V] {
	result := make(map[K]V, len(m.items)+len(items))
	newKeys := make([]K, len(m.keys))
	copy(newKeys, m.keys)

	for k, v := range m.items {
		result[k] = v
	}

	for k, v := range items {
		if _, exists := result[k]; !exists {
			newKeys = append(newKeys, k)
		}

		result[k] = v
	}

	return &Collection[K, V]{items: result, keys: newKeys}
}

// Union returns a new collection that is the union of the collection and the
// given map. Keys already present in the collection are not overwritten.
func (m *Collection[K, V]) Union(items map[K]V) *Collection[K, V] {
	result := make(map[K]V, len(m.items)+len(items))
	newKeys := make([]K, len(m.keys))
	copy(newKeys, m.keys)

	for k, v := range m.items {
		result[k] = v
	}

	for k, v := range items {
		if _, exists := result[k]; !exists {
			result[k] = v
			newKeys = append(newKeys, k)
		}
	}

	return &Collection[K, V]{items: result, keys: newKeys}
}

// Replace returns a new collection with the given map merged in. This is an
// alias for Merge.
func (m *Collection[K, V]) Replace(items map[K]V) *Collection[K, V] {
	return m.Merge(items)
}

// MergeRecursive recursively merges the given map into the collection.
// When both sides have a map value for the same key, they are merged recursively.
func MergeRecursive[V any](m *Collection[string, V], items map[string]V) *Collection[string, V] {
	return m.Merge(items)
}

// DiffKeys returns a new collection containing items whose keys are not
// present in the given map.
func (m *Collection[K, V]) DiffKeys(items map[K]V) *Collection[K, V] {
	result := make(map[K]V)
	newKeys := make([]K, 0)

	for _, k := range m.keys {
		if _, exists := items[k]; !exists {
			result[k] = m.items[k]
			newKeys = append(newKeys, k)
		}
	}

	return &Collection[K, V]{items: result, keys: newKeys}
}

// DiffKeysUsing returns items whose keys are not considered equal to any key in the given map,
// using the provided comparison function.
func (m *Collection[K, V]) DiffKeysUsing(items map[K]V, equals func(K, K) bool) *Collection[K, V] {
	result := make(map[K]V)
	newKeys := make([]K, 0)

	for _, k := range m.keys {
		found := false

		for otherK := range items {
			if equals(k, otherK) {
				found = true

				break
			}
		}

		if !found {
			result[k] = m.items[k]
			newKeys = append(newKeys, k)
		}
	}

	return &Collection[K, V]{items: result, keys: newKeys}
}

// IntersectByKeys returns a new collection containing items whose keys are
// present in the given map.
func (m *Collection[K, V]) IntersectByKeys(items map[K]V) *Collection[K, V] {
	result := make(map[K]V)
	newKeys := make([]K, 0)

	for _, k := range m.keys {
		if _, exists := items[k]; exists {
			result[k] = m.items[k]
			newKeys = append(newKeys, k)
		}
	}

	return &Collection[K, V]{items: result, keys: newKeys}
}

// DiffAssoc returns items whose key is not present in the given map or whose value differs.
func DiffAssoc[K comparable, V comparable](m *Collection[K, V], items map[K]V) *Collection[K, V] {
	result := make(map[K]V)
	newKeys := make([]K, 0)

	for _, k := range m.keys {
		otherVal, exists := items[k]

		if !exists || m.items[k] != otherVal {
			result[k] = m.items[k]
			newKeys = append(newKeys, k)
		}
	}

	return &Collection[K, V]{items: result, keys: newKeys}
}

// IntersectAssoc returns items whose key exists in the given map and whose value matches.
func IntersectAssoc[K comparable, V comparable](m *Collection[K, V], items map[K]V) *Collection[K, V] {
	result := make(map[K]V)
	newKeys := make([]K, 0)

	for _, k := range m.keys {
		if otherVal, exists := items[k]; exists && m.items[k] == otherVal {
			result[k] = m.items[k]
			newKeys = append(newKeys, k)
		}
	}

	return &Collection[K, V]{items: result, keys: newKeys}
}

// Flip swaps keys and values, returning a new Collection where the
// original values become keys and original keys become values.
func Flip[K comparable, V comparable](m *Collection[K, V]) *Collection[V, K] {
	result := make(map[V]K, len(m.items))
	newKeys := make([]V, 0, len(m.keys))

	for _, k := range m.keys {
		v := m.items[k]

		if _, exists := result[v]; !exists {
			newKeys = append(newKeys, v)
		}

		result[v] = k
	}

	return &Collection[V, K]{items: result, keys: newKeys}
}
