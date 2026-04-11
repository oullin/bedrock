package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JsonOptions controls JSON encoding behaviour.
type JsonOptions struct {
	Indent     bool
	EscapeHTML bool
}

// JsonResponse writes a JSON (or JSONP) response with configurable encoding.
type JsonResponse struct {
	writer     http.ResponseWriter
	statusCode int
	headers    http.Header
	cookies    []*http.Cookie
	data       any
	options    JsonOptions
	callback   string // JSONP callback function name
	exception  error
	original   any
}

// NewJsonResponse creates a JsonResponse that will serialise data as JSON.
func NewJsonResponse(w http.ResponseWriter, data any, status int, opts ...JsonOptions) *JsonResponse {
	var o JsonOptions

	if len(opts) > 0 {
		o = opts[0]
	}

	return &JsonResponse{
		writer:     w,
		statusCode: status,
		headers:    make(http.Header),
		data:       data,
		options:    o,
	}
}

// Status overrides the HTTP status code.
func (r *JsonResponse) Status(code int) *JsonResponse {
	r.statusCode = code

	return r
}

// GetStatusCode returns the current status code.
func (r *JsonResponse) GetStatusCode() int {
	return r.statusCode
}

// Header sets a single response header.
func (r *JsonResponse) Header(key, value string) *JsonResponse {
	r.headers.Set(key, value)

	return r
}

// WithHeaders sets multiple response headers at once.
func (r *JsonResponse) WithHeaders(headers map[string]string) *JsonResponse {
	for k, v := range headers {
		r.headers.Set(k, v)
	}

	return r
}

// Cookie appends a cookie to the response.
func (r *JsonResponse) Cookie(c *http.Cookie) *JsonResponse {
	r.cookies = append(r.cookies, c)

	return r
}

// WithCallback sets a JSONP callback function name. When set, the response
// body is wrapped in a function call and the Content-Type becomes
// "text/javascript".
func (r *JsonResponse) WithCallback(callback string) *JsonResponse {
	r.callback = callback

	return r
}

// WithException attaches an error for later inspection.
func (r *JsonResponse) WithException(err error) *JsonResponse {
	r.exception = err

	return r
}

// GetException returns the attached exception.
func (r *JsonResponse) GetException() error {
	return r.exception
}

// SetData replaces the response data.
func (r *JsonResponse) SetData(data any) *JsonResponse {
	r.data = data

	return r
}

// GetData returns the response data.
func (r *JsonResponse) GetData() any {
	return r.data
}

// SetEncodingOptions replaces the JSON encoding options.
func (r *JsonResponse) SetEncodingOptions(opts JsonOptions) *JsonResponse {
	r.options = opts

	return r
}

// Send encodes the data and writes the response.
func (r *JsonResponse) Send() error {
	body, err := r.encode()

	if err != nil {
		return err
	}

	r.original = r.data

	// Determine content type.
	ct := "application/json"

	if r.callback != "" {
		ct = "text/javascript; charset=utf-8"
		body = []byte(fmt.Sprintf("/**/ %s(%s);", r.callback, body))
	}

	if r.headers.Get("Content-Type") == "" {
		r.headers.Set("Content-Type", ct)
	}

	// Flush headers and cookies.
	dest := r.writer.Header()

	for k, vals := range r.headers {
		for _, v := range vals {
			dest.Add(k, v)
		}
	}

	for _, c := range r.cookies {
		http.SetCookie(r.writer, c)
	}

	r.writer.WriteHeader(r.statusCode)
	_, writeErr := r.writer.Write(body)

	return writeErr
}

// HasValidJSON returns true when the current data can be encoded as valid JSON.
func (r *JsonResponse) HasValidJSON() bool {
	_, err := r.encode()

	return err == nil
}

// GetOriginal returns the original data before encoding.
func (r *JsonResponse) GetOriginal() any {
	return r.original
}

// encode marshals data to JSON using the configured options.
func (r *JsonResponse) encode() ([]byte, error) {
	if r.options.Indent {
		return json.MarshalIndent(r.data, "", "    ")
	}

	return json.Marshal(r.data)
}

// FromJsonString creates a JsonResponse from a raw JSON string.
func FromJsonString(w http.ResponseWriter, jsonStr string, status int) *JsonResponse {
	resp := &JsonResponse{
		writer:     w,
		statusCode: status,
		headers:    make(http.Header),
		data:       json.RawMessage(jsonStr),
	}

	return resp
}
