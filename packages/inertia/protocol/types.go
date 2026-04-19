package protocol

import "errors"

// ErrNotFound indicates the requested resource could not be found.

// Props holds page properties passed to the frontend component.
type Props map[string]any

// TemplateData holds additional data passed to the root HTML template
// during initial (non-XHR) page visits.
type TemplateData map[string]any

// ValidationErrors holds field-level validation errors to be shared
// with the frontend via the "errors" prop.
type ValidationErrors map[string]any

// JSONMarshaler abstracts JSON encoding/decoding so callers can swap
// in a faster implementation (e.g. github.com/goccy/go-json) without
// changing any Inertia code.
type JSONMarshaler interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

// Logger is a minimal logging interface compatible with the standard
// library's log.Logger and most structured logging packages.
type Logger interface {
	Printf(format string, v ...any)
}

var ErrNotFound = errors.New("not found")
