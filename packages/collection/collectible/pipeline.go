package collectible

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bedrock/packages/collection/arr"
	"github.com/bedrock/packages/collection/support"
)

// Each iterates over items in insertion order, calling the callback for each
// key-value pair. Return false from the callback to stop iteration.
func (m *Collection[K, V]) Each(callback func(V, K) bool) *Collection[K, V] {
	for _, k := range m.keys {
		if !callback(m.items[k], k) {
			break
		}
	}

	return m
}

// Implode concatenates all values into a single string, separated by the
// given glue string.
func (m *Collection[K, V]) Implode(glue string) string {
	parts := make([]string, 0, len(m.keys))

	for _, k := range m.keys {
		parts = append(parts, fmt.Sprint(m.items[k]))
	}

	return strings.Join(parts, glue)
}

// Join concatenates all values into a string separated by glue. Optional
// final glue is placed between the last two items.
func (m *Collection[K, V]) Join(glue string, finalGlues ...string) string {
	parts := make([]string, 0, len(m.keys))

	for _, k := range m.keys {
		parts = append(parts, fmt.Sprint(m.items[k]))
	}

	return arr.Join(parts, glue, finalGlues...)
}

// Tap passes the collection to the callback and returns the collection,
// allowing side effects without breaking a method chain.
func (m *Collection[K, V]) Tap(callback func(*Collection[K, V])) *Collection[K, V] {
	callback(m)

	return m
}

// When applies the callback if the condition is true. An optional default
// callback is applied when the condition is false.
func (m *Collection[K, V]) When(condition bool, callback func(*Collection[K, V]) *Collection[K, V], defaults ...func(*Collection[K, V]) *Collection[K, V]) *Collection[K, V] {
	if condition {
		return callback(m)
	}

	if len(defaults) > 0 {
		return defaults[0](m)
	}

	return m
}

// Unless applies the callback unless the condition is true. An optional
// default callback is applied when the condition is true.
func (m *Collection[K, V]) Unless(condition bool, callback func(*Collection[K, V]) *Collection[K, V], defaults ...func(*Collection[K, V]) *Collection[K, V]) *Collection[K, V] {
	return m.When(!condition, callback, defaults...)
}

// ToJSON serializes the collection to JSON.
func (m *Collection[K, V]) ToJSON() ([]byte, error) {
	return json.Marshal(m.items)
}

// ToPrettyJSON serializes the collection to indented JSON.
func (m *Collection[K, V]) ToPrettyJSON() ([]byte, error) {
	return json.MarshalIndent(m.items, "", "    ")
}

// String returns the JSON representation of the collection.
func (m *Collection[K, V]) String() string {
	b, err := m.ToJSON()

	if err != nil {
		return "{}"
	}

	return string(b)
}

// MarshalJSON implements the json.Marshaler interface.
func (m *Collection[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.items)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (m *Collection[K, V]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &m.items)
}

// Copy creates a shallow copy of the collection.
func (m *Collection[K, V]) Copy() *Collection[K, V] {
	return New(m.All())
}

// Dump prints the underlying map to stdout for debugging purposes.
func (m *Collection[K, V]) Dump() *Collection[K, V] {
	fmt.Printf("%v\n", m.items)

	return m
}

// DD prints the map collection for debugging and terminates the program.
func (m *Collection[K, V]) DD() {
	m.Dump()
	os.Exit(1)
}

// ToPairs converts the collection to a slice of Pair values, preserving
// insertion order.
func (m *Collection[K, V]) ToPairs() []support.Pair[K, V] {
	result := make([]support.Pair[K, V], 0, len(m.keys))

	for _, k := range m.keys {
		result = append(result, support.Pair[K, V]{Key: k, Value: m.items[k]})
	}

	return result
}
