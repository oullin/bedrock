package client

import "net/http"

// RoundTripFunc performs a single HTTP round trip.
type RoundTripFunc func(req *http.Request) (*http.Response, error)

// Middleware intercepts an HTTP request before it is sent. It receives the
// request and a next function to call the next middleware (or the actual
// transport). It returns the response and any error.
type Middleware func(req *http.Request, next RoundTripFunc) (*http.Response, error)
