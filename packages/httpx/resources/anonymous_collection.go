package resources

import (
	"encoding/json"
	"net/http"
)

// AnonymousCollection wraps pre-transformed resource maps into a
// JSON-serialisable collection. Use it when the items have already been mapped
// or when a dedicated generic type is not available.
type AnonymousCollection struct {
	Items   []map[string]any
	WrapKey string
	With    map[string]any
}

// NewAnonymousCollection creates an AnonymousCollection from pre-mapped items.
func NewAnonymousCollection(items []map[string]any) *AnonymousCollection {
	return &AnonymousCollection{
		Items:   items,
		WrapKey: "data",
	}
}

// ToSlice returns the items, resolving any sentinel values in each map.
func (c *AnonymousCollection) ToSlice() []map[string]any {
	result := make([]map[string]any, len(c.Items))

	for i, item := range c.Items {
		result[i] = resolveMap(item)
	}

	return result
}

// ToJSON serialises the collection to JSON bytes.
func (c *AnonymousCollection) ToJSON() ([]byte, error) {
	data := c.toWrapped()

	return json.Marshal(data)
}

// Response writes the collection as a JSON HTTP response.
func (c *AnonymousCollection) Response(w http.ResponseWriter, _ *http.Request, status int) error {
	b, err := c.ToJSON()

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(b)

	return err
}

func (c *AnonymousCollection) toWrapped() any {
	data := c.ToSlice()

	if c.WrapKey != "" {
		wrapped := map[string]any{c.WrapKey: data}

		for k, v := range c.With {
			wrapped[k] = v
		}

		return wrapped
	}

	return data
}
