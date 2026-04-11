package client

import (
	"encoding/json"
	"io"
	"net/http"
)

// Response wraps an *http.Response with convenience status-checking and body
// access methods matching Upstream's Http client response API.
type Response struct {
	raw  *http.Response
	body []byte
	read bool
}

// NewResponse creates a Response from a raw *http.Response. The body is read
// lazily on first access.
func NewResponse(raw *http.Response) *Response {
	return &Response{raw: raw}
}

// Status returns the HTTP status code.
func (r *Response) Status() int {
	return r.raw.StatusCode
}

// Header returns a response header value.
func (r *Response) Header(key string) string {
	return r.raw.Header.Get(key)
}

// Headers returns all response headers.
func (r *Response) Headers() http.Header {
	return r.raw.Header
}

// Body returns the response body as a string.
func (r *Response) Body() string {
	return string(r.Bytes())
}

// Bytes returns the response body as bytes.
func (r *Response) Bytes() []byte {
	if !r.read {
		r.read = true

		if r.raw.Body != nil {
			r.body, _ = io.ReadAll(r.raw.Body)
			r.raw.Body.Close()
		}
	}

	return r.body
}

// JSON decodes the response body into the given value.
func (r *Response) JSON(v any) error {
	return json.Unmarshal(r.Bytes(), v)
}

// Ok returns true when the status code is 200.
func (r *Response) Ok() bool {
	return r.raw.StatusCode == http.StatusOK
}

// Created returns true when the status code is 201.
func (r *Response) Created() bool {
	return r.raw.StatusCode == http.StatusCreated
}

// Accepted returns true when the status code is 202.
func (r *Response) Accepted() bool {
	return r.raw.StatusCode == http.StatusAccepted
}

// NoContent returns true when the status code is 204.
func (r *Response) NoContent() bool {
	return r.raw.StatusCode == http.StatusNoContent
}

// MovedPermanently returns true when the status code is 301.
func (r *Response) MovedPermanently() bool {
	return r.raw.StatusCode == http.StatusMovedPermanently
}

// Found returns true when the status code is 302.
func (r *Response) Found() bool {
	return r.raw.StatusCode == http.StatusFound
}

// NotModified returns true when the status code is 304.
func (r *Response) NotModified() bool {
	return r.raw.StatusCode == http.StatusNotModified
}

// BadRequest returns true when the status code is 400.
func (r *Response) BadRequest() bool {
	return r.raw.StatusCode == http.StatusBadRequest
}

// Unauthorized returns true when the status code is 401.
func (r *Response) Unauthorized() bool {
	return r.raw.StatusCode == http.StatusUnauthorized
}

// PaymentRequired returns true when the status code is 402.
func (r *Response) PaymentRequired() bool {
	return r.raw.StatusCode == http.StatusPaymentRequired
}

// Forbidden returns true when the status code is 403.
func (r *Response) Forbidden() bool {
	return r.raw.StatusCode == http.StatusForbidden
}

// NotFound returns true when the status code is 404.
func (r *Response) NotFound() bool {
	return r.raw.StatusCode == http.StatusNotFound
}

// RequestTimeout returns true when the status code is 408.
func (r *Response) RequestTimeout() bool {
	return r.raw.StatusCode == http.StatusRequestTimeout
}

// Conflict returns true when the status code is 409.
func (r *Response) Conflict() bool {
	return r.raw.StatusCode == http.StatusConflict
}

// UnprocessableEntity returns true when the status code is 422.
func (r *Response) UnprocessableEntity() bool {
	return r.raw.StatusCode == http.StatusUnprocessableEntity
}

// TooManyRequests returns true when the status code is 429.
func (r *Response) TooManyRequests() bool {
	return r.raw.StatusCode == http.StatusTooManyRequests
}

// Successful returns true when the status is in the 2xx range.
func (r *Response) Successful() bool {
	return r.raw.StatusCode >= 200 && r.raw.StatusCode < 300
}

// Redirect returns true when the status is in the 3xx range.
func (r *Response) Redirect() bool {
	return r.raw.StatusCode >= 300 && r.raw.StatusCode < 400
}

// ClientError returns true when the status is in the 4xx range.
func (r *Response) ClientError() bool {
	return r.raw.StatusCode >= 400 && r.raw.StatusCode < 500
}

// ServerError returns true when the status is in the 5xx range.
func (r *Response) ServerError() bool {
	return r.raw.StatusCode >= 500 && r.raw.StatusCode < 600
}

// Failed returns true when the status is >= 400.
func (r *Response) Failed() bool {
	return r.raw.StatusCode >= 400
}

// Throw returns a RequestError if the response has a client or server error
// status code. Otherwise it returns nil.
func (r *Response) Throw() error {
	if r.Failed() {
		return &RequestError{Response: r}
	}

	return nil
}

// ThrowIf returns a RequestError if the given condition is true and the
// response indicates failure.
func (r *Response) ThrowIf(condition bool) error {
	if condition {
		return r.Throw()
	}

	return nil
}

// Cookies returns all cookies set on the response.
func (r *Response) Cookies() []*http.Cookie {
	return r.raw.Cookies()
}

// Raw returns the underlying *http.Response.
func (r *Response) Raw() *http.Response {
	return r.raw
}
