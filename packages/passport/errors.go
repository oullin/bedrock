package passport

import "errors"

var (
	// ErrTokenNotFound is returned when a token cannot be found in the store.
	ErrTokenNotFound = errors.New("passport: token not found")

	// ErrTokenRevoked is returned when a token has been revoked.
	ErrTokenRevoked = errors.New("passport: token has been revoked")

	// ErrTokenExpired is returned when a token has expired.
	ErrTokenExpired = errors.New("passport: token has expired")

	// ErrClientNotFound is returned when a client cannot be found in the store.
	ErrClientNotFound = errors.New("passport: client not found")

	// ErrClientRevoked is returned when a client has been revoked.
	ErrClientRevoked = errors.New("passport: client has been revoked")

	// ErrInvalidScope is returned when a token does not have the required scope.
	ErrInvalidScope = errors.New("passport: invalid scope(s) provided")

	// ErrUnauthenticated is returned when no valid token is present.
	ErrUnauthenticated = errors.New("passport: unauthenticated")

	// ErrInvalidGrant is returned when an unsupported grant type is requested.
	ErrInvalidGrant = errors.New("passport: invalid grant type")

	// ErrInvalidRequest is returned when a malformed request is received.
	ErrInvalidRequest = errors.New("passport: invalid request")
)
