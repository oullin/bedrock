package collectible

// GetOrPut returns the value for the given key if it exists. Otherwise, it
// stores and returns the provided default value.
func (m *Collection[K, V]) GetOrPut(key K, value V) V {
	if v, ok := m.items[key]; ok {
		return v
	}

	m.Put(key, value)

	return value
}

// Put sets the given key-value pair in the collection.
func (m *Collection[K, V]) Put(key K, value V) *Collection[K, V] {
	if _, exists := m.items[key]; !exists {
		m.keys = append(m.keys, key)
	}

	m.items[key] = value

	return m
}

// Pull removes an item by key and returns its value. The second return value
// indicates whether the key was found.
func (m *Collection[K, V]) Pull(key K) (V, bool) {
	v, ok := m.items[key]

	if ok {
		delete(m.items, key)

		for i, k := range m.keys {
			if k == key {
				m.keys = append(m.keys[:i], m.keys[i+1:]...)

				break
			}
		}
	}

	return v, ok
}

// Forget removes one or more items by key.
func (m *Collection[K, V]) Forget(keys ...K) *Collection[K, V] {
	for _, key := range keys {
		m.Pull(key)
	}

	return m
}
