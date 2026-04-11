package jsonapi

import (
	"net/http"
	"strings"
)

// SortField represents a parsed sort parameter. Descending is true when the
// original field was prefixed with "-".
type SortField struct {
	Field      string
	Descending bool
}

// Request wraps an *http.Request and provides helpers for parsing JSON:API
// query parameters (include, fields, sort, filter, page).
type Request struct {
	raw *http.Request
}

// NewRequest wraps an HTTP request for JSON:API parameter parsing.
func NewRequest(r *http.Request) *Request {
	return &Request{raw: r}
}

// Includes parses the "include" query parameter into a list of relationship
// paths. Returns nil when the parameter is absent.
func (r *Request) Includes() []string {
	raw := r.raw.URL.Query().Get("include")

	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)

		if p != "" {
			result = append(result, p)
		}
	}

	return result
}

// Fields parses the "fields[type]" query parameter for the given resource type.
// Returns the list of requested field names, or nil when the parameter is
// absent.
func (r *Request) Fields(resourceType string) []string {
	raw := r.raw.URL.Query().Get("fields[" + resourceType + "]")

	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)

		if p != "" {
			result = append(result, p)
		}
	}

	return result
}

// Sort parses the "sort" query parameter into SortField values. A leading "-"
// indicates descending order.
func (r *Request) Sort() []SortField {
	raw := r.raw.URL.Query().Get("sort")

	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	result := make([]SortField, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)

		if p == "" {
			continue
		}

		if strings.HasPrefix(p, "-") {
			result = append(result, SortField{Field: p[1:], Descending: true})
		} else {
			result = append(result, SortField{Field: p, Descending: false})
		}
	}

	return result
}

// Filter returns the value of "filter[key]" from the query string. Returns an
// empty string when absent.
func (r *Request) Filter(key string) string {
	return r.raw.URL.Query().Get("filter[" + key + "]")
}

// Page returns all "page[*]" query parameters as a map.
func (r *Request) Page() map[string]string {
	result := make(map[string]string)

	for key, values := range r.raw.URL.Query() {
		if strings.HasPrefix(key, "page[") && strings.HasSuffix(key, "]") {
			name := key[5 : len(key)-1]

			if len(values) > 0 {
				result[name] = values[0]
			}
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}
