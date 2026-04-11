package httpx

import (
	"errors"
	"fmt"
	"net/http"
)

// HttpResponseError is an error that carries an HTTP status code, headers, and
// an optional response payload. Middleware and handlers can inspect it to render
// a proper HTTP response.
type HttpResponseError struct {
	StatusCode int
	Message    string
	Headers    http.Header
	Response   any
}

// NewHttpResponseError creates an HttpResponseError with the given status and
// message.

// ThrottleRequestsError extends HttpResponseError with rate-limit metadata.
type ThrottleRequestsError struct {
	*HttpResponseError
	RetryAfter int // seconds until the client may retry
}

var (
	ErrPostTooLarge   = errors.New("httpx: post body exceeds the allowed size")
	ErrMalformedURL   = errors.New("httpx: malformed URL")
	ErrOriginMismatch = errors.New("httpx: origin does not match the request")
	ErrThrottle       = errors.New("httpx: too many requests")
)

func (e *HttpResponseError) Error() string {
	return fmt.Sprintf("httpx: HTTP %d: %s", e.StatusCode, e.Message)
}

func NewHttpResponseError(status int, message string) *HttpResponseError {
	return &HttpResponseError{
		StatusCode: status,
		Message:    message,
		Headers:    make(http.Header),
	}
}

// NewThrottleRequestsError creates a 429 error with retry-after metadata.
func NewThrottleRequestsError(retryAfter int) *ThrottleRequestsError {
	e := NewHttpResponseError(http.StatusTooManyRequests, "Too Many Requests")
	e.Headers.Set("Retry-After", fmt.Sprintf("%d", retryAfter))

	return &ThrottleRequestsError{
		HttpResponseError: e,
		RetryAfter:        retryAfter,
	}
}
