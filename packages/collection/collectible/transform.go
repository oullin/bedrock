package collectible

// Only returns a new collection containing only the specified keys.
func (m *Collection[K, V]) Only(keys ...K) *Collection[K, V] {
	result := make(map[K]V)
	newKeys := make([]K, 0, len(keys))

	for _, k := range keys {
		if v, ok := m.items[k]; ok {
			result[k] = v
			newKeys = append(newKeys, k)
		}
	}

	return &Collection[K, V]{items: result, keys: newKeys}
}

// Except returns a new collection containing all items except those with the
// specified keys.
func (m *Collection[K, V]) Except(keys ...K) *Collection[K, V] {
	excludeSet := make(map[K]bool, len(keys))

	for _, k := range keys {
		excludeSet[k] = true
	}

	result := make(map[K]V)
	newKeys := make([]K, 0)

	for _, k := range m.keys {
		if !excludeSet[k] {
			result[k] = m.items[k]
			newKeys = append(newKeys, k)
		}
	}

	return &Collection[K, V]{items: result, keys: newKeys}
}

// Filter returns a new collection containing only the items for which the
// callback returns true.
func (m *Collection[K, V]) Filter(callback func(V, K) bool) *Collection[K, V] {
	result := make(map[K]V)
	newKeys := make([]K, 0)

	for _, k := range m.keys {
		v := m.items[k]

		if callback(v, k) {
			result[k] = v
			newKeys = append(newKeys, k)
		}
	}

	return &Collection[K, V]{items: result, keys: newKeys}
}

// Reject returns a new collection containing only the items for which the
// callback returns false.
func (m *Collection[K, V]) Reject(callback func(V, K) bool) *Collection[K, V] {
	return m.Filter(func(v V, k K) bool {
		return !callback(v, k)
	})
}

// MapValues transforms each value using the callback, returning a new
// Collection with the same keys and transformed values.
func MapValues[K comparable, V any, R any](m *Collection[K, V], callback func(V, K) R) *Collection[K, R] {
	result := make(map[K]R, len(m.items))
	newKeys := make([]K, len(m.keys))
	copy(newKeys, m.keys)

	for _, k := range m.keys {
		result[k] = callback(m.items[k], k)
	}

	return &Collection[K, R]{items: result, keys: newKeys}
}

// Every report whether all items in the collection satisfy the callback.
func (m *Collection[K, V]) Every(callback func(V, K) bool) bool {
	for _, k := range m.keys {
		if !callback(m.items[k], k) {
			return false
		}
	}

	return true
}

// Partition splits the collection into two: one where the callback returns
// true, and one where it returns false.
func (m *Collection[K, V]) Partition(callback func(V, K) bool) (*Collection[K, V], *Collection[K, V]) {
	pass := make(map[K]V)
	passKeys := make([]K, 0)
	fail := make(map[K]V)
	failKeys := make([]K, 0)

	for _, k := range m.keys {
		v := m.items[k]

		if callback(v, k) {
			pass[k] = v
			passKeys = append(passKeys, k)
		} else {
			fail[k] = v
			failKeys = append(failKeys, k)
		}
	}

	return &Collection[K, V]{items: pass, keys: passKeys},
		&Collection[K, V]{items: fail, keys: failKeys}
}
