package client

import (
	"errors"
	"fmt"
)

// RequestError is returned when a request completes but the server responds
// with a client or server error status code.
type RequestError struct {
	Response *Response
}

// ConnectionError is returned when the HTTP client fails to connect.
type ConnectionError struct {
	URL string
	Err error
}

var (
	ErrConnection      = errors.New("client: connection failed")
	ErrStrayRequest    = errors.New("client: unexpected request during fake mode")
	ErrBatchInProgress = errors.New("client: batch already in progress")
)

func (e *RequestError) Error() string {
	return fmt.Sprintf("client: HTTP %d", e.Response.Status())
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("client: connection to %s failed: %v", e.URL, e.Err)
}

func (e *ConnectionError) Unwrap() error {
	return e.Err
}
