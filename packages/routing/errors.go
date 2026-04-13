package routing

import "errors"

var (
	ErrRouteNotFound    = errors.New("routing: no matching route found")
	ErrMethodNotAllowed = errors.New("routing: method not allowed for this URI")
	ErrInvalidSignature = errors.New("routing: invalid URL signature")
	ErrSignatureExpired = errors.New("routing: URL signature has expired")
)
