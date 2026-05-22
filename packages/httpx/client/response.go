package client

import (
	"encoding/json"
	"io"
	"net/http"
)

// Response wraps an *http.Response with convenience status-checking and body
// access methods matching the upstream Http client response API.
type Response struct {
	raw   *http.Response
	body  []byte
	read  bool
	stats map[string]any
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

// Reason returns the textual reason phrase for the status code (e.g. "OK",
// "Not Found").
func (r *Response) Reason() string {
	return http.StatusText(r.raw.StatusCode)
}

// Collect decodes the JSON response body into a map. When a key is provided
// only the nested value for that key is returned.
func (r *Response) Collect(key ...string) map[string]any {
	var data map[string]any

	if err := json.Unmarshal(r.Bytes(), &data); err != nil {
		return nil
	}

	if len(key) > 0 && key[0] != "" {
		if nested, ok := data[key[0]]; ok {
			if m, ok := nested.(map[string]any); ok {
				return m
			}
		}

		return nil
	}

	return data
}

// EffectiveUri returns the final URL after following any redirects.
func (r *Response) EffectiveUri() string {
	if r.raw.Request != nil && r.raw.Request.URL != nil {
		return r.raw.Request.URL.String()
	}

	return ""
}

// ThrowUnless returns a RequestError when the condition is false and the
// response indicates failure.
func (r *Response) ThrowUnless(condition bool) error {
	if !condition {
		return r.Throw()
	}

	return nil
}

// ThrowIfStatus returns a RequestError when the response status matches the
// given code, regardless of whether the status is normally considered a failure.
func (r *Response) ThrowIfStatus(code int) error {
	if r.raw.StatusCode == code {
		return &RequestError{Response: r}
	}

	return nil
}

// ThrowUnlessStatus returns a RequestError when the response status does not
// match the given code.
func (r *Response) ThrowUnlessStatus(code int) error {
	if r.raw.StatusCode != code {
		return &RequestError{Response: r}
	}

	return nil
}

// ThrowIfClientError returns a RequestError when the response has a 4xx status.
func (r *Response) ThrowIfClientError() error {
	if r.ClientError() {
		return &RequestError{Response: r}
	}

	return nil
}

// ThrowIfServerError returns a RequestError when the response has a 5xx status.
func (r *Response) ThrowIfServerError() error {
	if r.ServerError() {
		return &RequestError{Response: r}
	}

	return nil
}

// OnError calls the given callback when the response indicates failure and
// returns the response for chaining.
func (r *Response) OnError(fn func(*Response)) *Response {
	if r.Failed() {
		fn(r)
	}

	return r
}

// Close closes the response body if it has not already been read.
func (r *Response) Close() error {
	if !r.read && r.raw.Body != nil {
		return r.raw.Body.Close()
	}

	return nil
}

// ToException returns a RequestError when the response indicates failure.
// Returns nil for successful responses.
func (r *Response) ToException() *RequestError {
	if r.Failed() {
		return &RequestError{Response: r}
	}

	return nil
}

// HandlerStats returns transfer statistics collected during the request.
// Keys may include "dns_ms", "connect_ms", "tls_ms", and "total_ms".
func (r *Response) HandlerStats() map[string]any {
	if r.stats == nil {
		return map[string]any{}
	}

	return r.stats
}

// SetStats sets the transfer statistics for this response. This is intended
// for internal use by the client.
func (r *Response) SetStats(stats map[string]any) {
	r.stats = stats
}

// Raw returns the underlying *http.Response.
func (r *Response) Raw() *http.Response {
	return r.raw
}
