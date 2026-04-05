package http

import nethttp "net/http"

// Request aliases the standard library HTTP request type.
type Request = nethttp.Request

// ResponseWriter aliases the standard library response writer type.
type ResponseWriter = nethttp.ResponseWriter

// Middleware decorates an HTTP handler.
type Middleware func(nethttp.Handler) nethttp.Handler

// Capture mirrors Laravel's request capture entrypoint.
func Capture(request *nethttp.Request) *Request {
	return request
}
