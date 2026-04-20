package collection

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bedrock/packages/collection/arr"
)

// Each iterates over the items, passing each item and its index to the callback.
// Return false from the callback to stop iterating.
func (c *Collection[T]) Each(callback func(T, int) bool) *Collection[T] {
	for i, item := range c.items {
		if !callback(item, i) {
			break
		}
	}

	return c
}

// EachSpread iterates over the collection's items, passing each item value
// into the given callback. In Go this operates the same as Each.
func (c *Collection[T]) EachSpread(callback func(T, int) bool) *Collection[T] {
	return c.Each(callback)
}

// Tap passes the collection to the given callback and returns the collection unchanged.
func (c *Collection[T]) Tap(callback func(*Collection[T])) *Collection[T] {
	callback(c)

	return c
}

// TapEach calls the given callback on each item for side effects, returning the original collection.
func (c *Collection[T]) TapEach(callback func(T, int)) *Collection[T] {
	for i, item := range c.items {
		callback(item, i)
	}

	return c
}

// Pipe passes the collection to the given callback and returns the callback's result.
func Pipe[T any, R any](c *Collection[T], callback func(*Collection[T]) R) R {
	return callback(c)
}

// PipeInto passes the collection to the given constructor and returns the result.
func PipeInto[T any, R any](c *Collection[T], constructor func(*Collection[T]) R) R {
	return constructor(c)
}

// PipeThrough passes the collection through a series of callbacks, returning the final result.
func PipeThrough[T any](c *Collection[T], callbacks ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	result := c

	for _, cb := range callbacks {
		result = cb(result)
	}

	return result
}

// When applies the callback if the condition is true; otherwise it applies the optional default callback.
func (c *Collection[T]) When(condition bool, callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	if condition {
		return callback(c)
	}

	if len(defaults) > 0 {
		return defaults[0](c)
	}

	return c
}

// WhenEmpty applies the callback when the collection is empty.
func (c *Collection[T]) WhenEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.When(c.IsEmpty(), callback, defaults...)
}

// WhenNotEmpty applies the callback when the collection is not empty.
func (c *Collection[T]) WhenNotEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.When(c.IsNotEmpty(), callback, defaults...)
}

// Unless applies the callback unless the condition is true.
func (c *Collection[T]) Unless(condition bool, callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.When(!condition, callback, defaults...)
}

// UnlessEmpty applies the callback unless the collection is empty.
func (c *Collection[T]) UnlessEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.WhenNotEmpty(callback, defaults...)
}

// UnlessNotEmpty applies the callback unless the collection is not empty.
func (c *Collection[T]) UnlessNotEmpty(callback func(*Collection[T]) *Collection[T], defaults ...func(*Collection[T]) *Collection[T]) *Collection[T] {
	return c.WhenEmpty(callback, defaults...)
}

// Implode joins elements into a string using the given glue, converting each item via fmt.Sprint.
func (c *Collection[T]) Implode(glue string) string {
	parts := make([]string, len(c.items))

	for i, item := range c.items {
		parts[i] = fmt.Sprint(item)
	}

	return strings.Join(parts, glue)
}

// Join joins elements into a string using the given glue,
// with an optional final separator between the last two items.
func (c *Collection[T]) Join(glue string, finalGlues ...string) string {
	parts := make([]string, len(c.items))

	for i, item := range c.items {
		parts[i] = fmt.Sprint(item)
	}

	return arr.Join(parts, glue, finalGlues...)
}

// ToSlice returns a copy of the underlying slice.
func (c *Collection[T]) ToSlice() []T {
	result := make([]T, len(c.items))
	copy(result, c.items)

	return result
}

// ToJSON serializes the collection to JSON bytes.
func (c *Collection[T]) ToJSON() ([]byte, error) {
	return json.Marshal(c.items)
}

// ToPrettyJSON serializes the collection to indented JSON bytes.
func (c *Collection[T]) ToPrettyJSON() ([]byte, error) {
	return json.MarshalIndent(c.items, "", "    ")
}

// String returns the JSON string representation of the collection.
func (c *Collection[T]) String() string {
	b, err := c.ToJSON()

	if err != nil {
		return "[]"
	}

	return string(b)
}

// MarshalJSON implements the json.Marshaler interface.
func (c *Collection[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.items)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *Collection[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &c.items)
}

// Copy creates a shallow copy of the collection.
func (c *Collection[T]) Copy() *Collection[T] {
	result := make([]T, len(c.items))
	copy(result, c.items)

	return Collect(result)
}

// Dump prints the collection items to stdout for debugging.
func (c *Collection[T]) Dump() *Collection[T] {
	fmt.Printf("%v\n", c.items)

	return c
}

// DD prints the collection items for debugging and terminates the program.
func (c *Collection[T]) DD() {
	c.Dump()
	os.Exit(1)
}
