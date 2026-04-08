package http

import (
	"encoding/json"
	"fmt"
	nethttp "net/http"
)

// JsonResponse represents a JSON HTTP response with status code, data,
// and encoding options. It mirrors Upstream's JsonResponse class.
type JsonResponse struct {
	statusCode      int
	data            any
	originalContent any
	indent          string
}

// NewJsonResponse creates a JSON response with the given data and status code.
func NewJsonResponse(data any, status int) *JsonResponse {
	return &JsonResponse{
		statusCode:      status,
		data:            data,
		originalContent: data,
	}
}

// NewJsonResponseFromString creates a JSON response from a raw JSON string.
// It returns an error if the string is not valid JSON.
func NewJsonResponseFromString(raw string, status int) (*JsonResponse, error) {
	var parsed any
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil, fmt.Errorf("invalid JSON string: %w", err)
	}

	return &JsonResponse{
		statusCode:      status,
		data:            parsed,
		originalContent: raw,
	}, nil
}

// SetData replaces the response data. Returns an error if the data
// cannot be marshalled to JSON.
func (jr *JsonResponse) SetData(data any) error {
	if _, err := json.Marshal(data); err != nil {
		return fmt.Errorf("json response: %w", err)
	}

	jr.data = data
	jr.originalContent = data

	return nil
}

// GetData returns the current response data.
func (jr *JsonResponse) GetData() any {
	return jr.data
}

// GetOriginalContent returns the original content set on the response.
func (jr *JsonResponse) GetOriginalContent() any {
	return jr.originalContent
}

// SetStatusCode sets the HTTP status code.
func (jr *JsonResponse) SetStatusCode(code int) {
	jr.statusCode = code
}

// GetStatusCode returns the HTTP status code.
func (jr *JsonResponse) GetStatusCode() int {
	return jr.statusCode
}

// SetIndent sets the JSON indentation string (e.g. "  " for pretty-print).
// An empty string disables indentation.
func (jr *JsonResponse) SetIndent(indent string) {
	jr.indent = indent
}

// GetIndent returns the current indentation string.
func (jr *JsonResponse) GetIndent() string {
	return jr.indent
}

// WriteTo writes the JSON response to the given ResponseWriter.
func (jr *JsonResponse) WriteTo(w ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(jr.statusCode)

	enc := json.NewEncoder(w)
	if jr.indent != "" {
		enc.SetIndent("", jr.indent)
	}

	return enc.Encode(jr.data)
}

// Send is a convenience that calls WriteTo. It satisfies a common
// "send the response" pattern.
func (jr *JsonResponse) Send(w nethttp.ResponseWriter) error {
	return jr.WriteTo(w)
}
