package jsonapi

import (
	"encoding/json"
	"net/http"
)

// Collection wraps a slice of items and produces a JSON:API compliant document
// with a "data" array and optional "included" sideloading.
type Collection[T any] struct {
	Items      []T
	ResourceFn func(T) *JsonApiResource[T]
	Included   []Resource
	Links      map[string]any
	Meta       map[string]any
}

// NewCollection creates a JSON:API collection from items and a function that
// wraps each item in a JsonApiResource.
func NewCollection[T any](items []T, resourceFn func(T) *JsonApiResource[T]) *Collection[T] {
	return &Collection[T]{
		Items:      items,
		ResourceFn: resourceFn,
	}
}

// WithIncluded adds resources to the "included" section of the document.
// Duplicate entries (same type+id) are automatically removed.
func (c *Collection[T]) WithIncluded(resources ...Resource) *Collection[T] {
	c.Included = append(c.Included, resources...)

	return c
}

// WithLinks sets the top-level "links" object.
func (c *Collection[T]) WithLinks(links map[string]any) *Collection[T] {
	c.Links = links

	return c
}

// WithMeta sets the top-level "meta" object.
func (c *Collection[T]) WithMeta(meta map[string]any) *Collection[T] {
	c.Meta = meta

	return c
}

// ToDocument produces the JSON:API document map with "data" and optional
// "included", "links", and "meta" keys.
func (c *Collection[T]) ToDocument(req *http.Request) map[string]any {
	data := make([]map[string]any, len(c.Items))

	for i, item := range c.Items {
		resource := c.ResourceFn(item)
		data[i] = resource.ToResourceObject(req)
	}

	doc := map[string]any{"data": data}

	if len(c.Included) > 0 {
		doc["included"] = c.deduplicateIncluded(req)
	}

	if len(c.Links) > 0 {
		doc["links"] = c.Links
	}

	if len(c.Meta) > 0 {
		doc["meta"] = c.Meta
	}

	return doc
}

// ToJSON serialises the collection as a JSON:API document.
func (c *Collection[T]) ToJSON(req *http.Request) ([]byte, error) {
	return json.Marshal(c.ToDocument(req))
}

// Response writes the collection as a JSON:API HTTP response.
func (c *Collection[T]) Response(w http.ResponseWriter, req *http.Request, status int) error {
	b, err := c.ToJSON(req)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)
	_, err = w.Write(b)

	return err
}

// deduplicateIncluded resolves included resources and removes duplicates by
// type+id.
func (c *Collection[T]) deduplicateIncluded(req *http.Request) []map[string]any {
	seen := make(map[string]struct{})

	var result []map[string]any

	for _, inc := range c.Included {
		obj := inc.ToResourceObject(req)

		key := ""

		if t, ok := obj["type"].(string); ok {
			key += t
		}

		key += ":"

		if id, ok := obj["id"].(string); ok {
			key += id
		}

		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		result = append(result, obj)
	}

	return result
}
