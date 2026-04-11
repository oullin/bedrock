package client

import "net/http"

// RequestSending is dispatched before an HTTP request is sent.
type RequestSending struct {
	Request *http.Request
}

// ResponseReceived is dispatched after an HTTP response is received.
type ResponseReceived struct {
	Request  *http.Request
	Response *Response
}

// ConnectionFailed is dispatched when a connection attempt fails.
type ConnectionFailed struct {
	Request *http.Request
	Err     error
}
