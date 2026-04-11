package client

import (
	"errors"
	"fmt"
)

var (
	ErrConnection      = errors.New("client: connection failed")
	ErrStrayRequest    = errors.New("client: unexpected request during fake mode")
	ErrBatchInProgress = errors.New("client: batch already in progress")
)

// RequestError is returned when a request completes but the server responds
// with a client or server error status code.
type RequestError struct {
	Response *Response
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("client: HTTP %d", e.Response.Status())
}

// ConnectionError is returned when the HTTP client fails to connect.
type ConnectionError struct {
	URL string
	Err error
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("client: connection to %s failed: %v", e.URL, e.Err)
}

func (e *ConnectionError) Unwrap() error {
	return e.Err
}
