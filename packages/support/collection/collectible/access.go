package collectible

// Get returns the value for the given key. The second return value indicates
// whether the key was found. An optional default may be provided.
func (m *Collection[K, V]) Get(key K, defaults ...V) (V, bool) {
	if v, ok := m.items[key]; ok {
		return v, true
	}

	if len(defaults) > 0 {
		return defaults[0], false
	}

	var zero V

	return zero, false
}

// Has reports whether the given key exists in the collection.
func (m *Collection[K, V]) Has(key K) bool {
	_, ok := m.items[key]

	return ok
}

// HasAny reports whether any of the given keys exist in the collection.
func (m *Collection[K, V]) HasAny(keys ...K) bool {
	for _, k := range keys {
		if m.Has(k) {
			return true
		}
	}

	return false
}

// Contains reports whether any item in the collection satisfies the predicate.
func (m *Collection[K, V]) Contains(predicate func(V, K) bool) bool {
	for _, k := range m.keys {
		if predicate(m.items[k], k) {
			return true
		}
	}

	return false
}

// Some is an alias for Contains.
func (m *Collection[K, V]) Some(predicate func(V, K) bool) bool {
	return m.Contains(predicate)
}

// DoesntContain reports whether no item in the collection satisfies the predicate.
func (m *Collection[K, V]) DoesntContain(predicate func(V, K) bool) bool {
	return !m.Contains(predicate)
}

// HasSole reports whether exactly one item in the collection satisfies the predicate.
func (m *Collection[K, V]) HasSole(predicate func(V, K) bool) bool {
	count := 0

	for _, k := range m.keys {
		if predicate(m.items[k], k) {
			count++

			if count > 1 {
				return false
			}
		}
	}

	return count == 1
}

// Search returns the key of the first item that satisfies the predicate.
// The second return value indicates whether a match was found.
func (m *Collection[K, V]) Search(predicate func(V, K) bool) (K, bool) {
	for _, k := range m.keys {
		if predicate(m.items[k], k) {
			return k, true
		}
	}

	var zero K

	return zero, false
}

// First returns the first value that satisfies the optional predicate. If no
// predicate is given, the first item in the insertion order is returned.
func (m *Collection[K, V]) First(predicates ...func(V, K) bool) (V, bool) {
	if len(m.keys) == 0 {
		var zero V

		return zero, false
	}

	if len(predicates) == 0 || predicates[0] == nil {
		return m.items[m.keys[0]], true
	}

	predicate := predicates[0]

	for _, k := range m.keys {
		if predicate(m.items[k], k) {
			return m.items[k], true
		}
	}

	var zero V

	return zero, false
}

// Last returns the last value that satisfies the optional predicate. If no
// predicate is given, the last item in the insertion order is returned.
func (m *Collection[K, V]) Last(predicates ...func(V, K) bool) (V, bool) {
	if len(m.keys) == 0 {
		var zero V

		return zero, false
	}

	if len(predicates) == 0 || predicates[0] == nil {
		return m.items[m.keys[len(m.keys)-1]], true
	}

	predicate := predicates[0]

	for i := len(m.keys) - 1; i >= 0; i-- {
		k := m.keys[i]

		if predicate(m.items[k], k) {
			return m.items[k], true
		}
	}

	var zero V

	return zero, false
}
