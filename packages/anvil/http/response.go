package http

import (
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"strings"
)

// Response represents an HTTP response with status code, headers, and content.
// It mirrors Laravel's Illuminate\Http\Response class.
type Response struct {
	statusCode      int
	headers         nethttp.Header
	content         []byte
	originalContent any
}

// NewResponse creates a response with the given content and status code.
func NewResponse(content any, status int) (*Response, error) {
	r := &Response{
		statusCode: status,
		headers:    nethttp.Header{},
	}

	if err := r.SetContent(content); err != nil {
		return nil, err
	}

	return r, nil
}

// SetContent sets the response body. Strings and []byte are used directly.
// Other types are JSON-encoded and the Content-Type header is set to
// "application/json".
func (r *Response) SetContent(data any) error {
	r.originalContent = data

	switch v := data.(type) {
	case nil:
		r.content = nil
	case string:
		r.content = []byte(v)
	case []byte:
		r.content = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("response: %w", err)
		}
		r.content = b
		r.headers.Set("Content-Type", "application/json")
	}

	return nil
}

// GetContent returns the response body as bytes.
func (r *Response) GetContent() []byte {
	return r.content
}

// GetOriginalContent returns the original content before encoding.
func (r *Response) GetOriginalContent() any {
	return r.originalContent
}

// SetStatusCode sets the HTTP status code.
func (r *Response) SetStatusCode(code int) {
	r.statusCode = code
}

// GetStatusCode returns the HTTP status code.
func (r *Response) GetStatusCode() int {
	return r.statusCode
}

// StatusText returns the text phrase for the current status code.
func (r *Response) StatusText() string {
	return nethttp.StatusText(r.statusCode)
}

// SetHeader sets a single header value (replaces any existing value).
func (r *Response) SetHeader(key, value string) *Response {
	r.headers.Set(key, value)
	return r
}

// WithHeaders sets multiple headers from a map.
func (r *Response) WithHeaders(headers map[string]string) *Response {
	for k, v := range headers {
		r.headers.Set(k, v)
	}
	return r
}

// WithoutHeader removes one or more headers by key.
func (r *Response) WithoutHeader(keys ...string) *Response {
	for _, k := range keys {
		r.headers.Del(k)
	}
	return r
}

// GetHeader returns a header value, or the default if not set.
func (r *Response) GetHeader(key string, defaultValue ...string) string {
	val := r.headers.Get(key)
	if val != "" {
		return val
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// Headers returns the underlying header map.
func (r *Response) Headers() nethttp.Header {
	return r.headers
}

// WriteTo writes the response to the given ResponseWriter.
func (r *Response) WriteTo(w ResponseWriter) error {
	for k, vals := range r.headers {
		for _, v := range vals {
			w.Header().Add(k, v)
		}
	}

	// If no explicit Content-Type was set, try to detect.
	if w.Header().Get("Content-Type") == "" {
		if json.Valid(r.content) && len(r.content) > 0 && (r.content[0] == '{' || r.content[0] == '[') {
			w.Header().Set("Content-Type", "application/json")
		} else if len(r.content) > 0 && isLikelyHTML(r.content) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		} else {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		}
	}

	w.WriteHeader(r.statusCode)

	if r.content != nil {
		_, err := w.Write(r.content)
		return err
	}

	return nil
}

func isLikelyHTML(b []byte) bool {
	s := strings.TrimSpace(string(b))
	return strings.HasPrefix(s, "<") && strings.Contains(s, ">")
}
