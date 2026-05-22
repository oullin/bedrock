package resources

import (
	"encoding/json"
	"net/http"
)

// PaginationMeta holds pagination metadata matching the upstream conventions.
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

// ToJSON serialises the paginated response including data, meta and links.

// Merge additional top-level data from the collection.

// Response writes the paginated collection as a JSON HTTP response.

// ResourceResponse wraps a Resource and provides response-level features such
// as automatic 201 status for recently created resources, custom headers, and
// a response callback.
type ResourceResponse struct {
	resource         Resource
	recentlyCreated  bool
	headers          map[string]string
	responseCallback func(http.ResponseWriter, *http.Request)
}

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

func (p *PaginatedResponse[T]) ToJSON(req *http.Request) ([]byte, error) {
	data := p.Collection.ToSlice(req)

	result := map[string]any{
		"data":  data,
		"meta":  p.Meta,
		"links": p.Links,
	}

	for k, v := range p.Collection.With {
		result[k] = v
	}

	return json.Marshal(result)
}

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

// NewResourceResponse wraps a resource for HTTP response rendering.
func NewResourceResponse(resource Resource) *ResourceResponse {
	return &ResourceResponse{
		resource: resource,
		headers:  make(map[string]string),
	}
}

// RecentlyCreated marks the resource as recently created so that the response
// automatically uses 201 instead of 200 when no explicit status is given.
func (rr *ResourceResponse) RecentlyCreated() *ResourceResponse {
	rr.recentlyCreated = true

	return rr
}

// WithHeader adds a header to the response.
func (rr *ResourceResponse) WithHeader(key, value string) *ResourceResponse {
	rr.headers[key] = value

	return rr
}

// WithResponse registers a callback invoked after headers are written but
// before the body is sent.
func (rr *ResourceResponse) WithResponse(fn func(http.ResponseWriter, *http.Request)) *ResourceResponse {
	rr.responseCallback = fn

	return rr
}

// calculateStatus returns 201 when the resource was recently created,
// otherwise 200.
func (rr *ResourceResponse) calculateStatus() int {
	if rr.recentlyCreated {
		return http.StatusCreated
	}

	return http.StatusOK
}

// Response writes the resource as a JSON HTTP response. The status code
// defaults to the calculated status (200 or 201).
func (rr *ResourceResponse) Response(w http.ResponseWriter, req *http.Request) error {
	data := rr.resource.ToMap(req)
	b, err := json.Marshal(data)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	for k, v := range rr.headers {
		w.Header().Set(k, v)
	}

	if rr.responseCallback != nil {
		rr.responseCallback(w, req)
	}

	w.WriteHeader(rr.calculateStatus())
	_, err = w.Write(b)

	return err
}
