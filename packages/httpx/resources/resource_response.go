package resources

import (
	"encoding/json"
	"net/http"
)

// PaginationMeta holds pagination metadata matching Laravel's conventions.
type PaginationMeta struct {
	CurrentPage int    `json:"current_page"`
	LastPage    int    `json:"last_page"`
	PerPage     int    `json:"per_page"`
	Total       int    `json:"total"`
	From        int    `json:"from"`
	To          int    `json:"to"`
	Path        string `json:"path"`
}

// PaginationLinks holds pagination link URLs.
type PaginationLinks struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
}

// PaginatedResponse wraps a collection with pagination metadata.
type PaginatedResponse[T any] struct {
	Collection *Collection[T]
	Meta       PaginationMeta
	Links      PaginationLinks
}

// NewPaginatedResponse creates a PaginatedResponse from a collection and
// pagination metadata.
func NewPaginatedResponse[T any](collection *Collection[T], meta PaginationMeta, links ...PaginationLinks) *PaginatedResponse[T] {
	pr := &PaginatedResponse[T]{
		Collection: collection,
		Meta:       meta,
	}

	if len(links) > 0 {
		pr.Links = links[0]
	}

	return pr
}

// ToJSON serialises the paginated response including data, meta and links.
func (p *PaginatedResponse[T]) ToJSON(req *http.Request) ([]byte, error) {
	data := p.Collection.ToSlice(req)

	result := map[string]any{
		"data":  data,
		"meta":  p.Meta,
		"links": p.Links,
	}

	// Merge additional top-level data from the collection.
	for k, v := range p.Collection.With {
		result[k] = v
	}

	return json.Marshal(result)
}

// Response writes the paginated collection as a JSON HTTP response.
func (p *PaginatedResponse[T]) Response(w http.ResponseWriter, req *http.Request, status int) error {
	b, err := p.ToJSON(req)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(b)

	return err
}
