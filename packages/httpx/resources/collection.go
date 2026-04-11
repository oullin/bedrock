package resources

import (
	"encoding/json"
	"net/http"
)

// Collection wraps a slice of items and a mapper function to produce a
// JSON-serialisable array of resources.
type Collection[T any] struct {
	Items   []T
	MapFunc func(T, *http.Request) map[string]any
	WrapKey string
	With    map[string]any
}

// NewCollection creates a Collection from a slice and a mapping function.
func NewCollection[T any](items []T, mapFn func(T, *http.Request) map[string]any) *Collection[T] {
	return &Collection[T]{
		Items:   items,
		MapFunc: mapFn,
		WrapKey: "data",
	}
}

// ToSlice serialises each item through the mapper and returns the result.
func (c *Collection[T]) ToSlice(req *http.Request) []map[string]any {
	result := make([]map[string]any, len(c.Items))

	for i, item := range c.Items {
		result[i] = resolveMap(c.MapFunc(item, req))
	}

	return result
}

// ToJSON serialises the collection to JSON bytes.
func (c *Collection[T]) ToJSON(req *http.Request) ([]byte, error) {
	data := c.toWrapped(req)

	return json.Marshal(data)
}

// Response writes the collection as a JSON HTTP response.
func (c *Collection[T]) Response(w http.ResponseWriter, req *http.Request, status int) error {
	b, err := c.ToJSON(req)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(b)

	return err
}

func (c *Collection[T]) toWrapped(req *http.Request) any {
	data := c.ToSlice(req)

	if c.WrapKey != "" {
		wrapped := map[string]any{c.WrapKey: data}

		for k, v := range c.With {
			wrapped[k] = v
		}

		return wrapped
	}

	return data
}
