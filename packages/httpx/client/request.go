package client

import "net/http"

// Request wraps an outbound *http.Request with convenience accessors.
type Request struct {
	raw  *http.Request
	body []byte
}

// NewRequest creates a Request wrapper.
func NewRequest(raw *http.Request, body ...[]byte) *Request {
	r := &Request{raw: raw}

	if len(body) > 0 {
		r.body = body[0]
	}

	return r
}

// URL returns the request URL string.
func (r *Request) URL() string {
	return r.raw.URL.String()
}

// Method returns the HTTP method.
func (r *Request) Method() string {
	return r.raw.Method
}

// Header returns a request header value.
func (r *Request) Header(key string) string {
	return r.raw.Header.Get(key)
}

// Headers returns all request headers.
func (r *Request) Headers() http.Header {
	return r.raw.Header
}

// HasHeader checks whether the named header is present.
func (r *Request) HasHeader(key, value string) bool {
	return r.raw.Header.Get(key) == value
}

// Body returns the request body bytes.
func (r *Request) Body() []byte {
	return r.body
}

// Raw returns the underlying *http.Request.
func (r *Request) Raw() *http.Request {
	return r.raw
}
