package oauthserver

import "errors"

var (
	// ErrTokenNotFound is returned when a token cannot be found in the store.
	ErrTokenNotFound = errors.New("oauthserver: token not found")

	// ErrTokenRevoked is returned when a token has been revoked.
	ErrTokenRevoked = errors.New("oauthserver: token has been revoked")

	// ErrTokenExpired is returned when a token has expired.
	ErrTokenExpired = errors.New("oauthserver: token has expired")

	// ErrClientNotFound is returned when a client cannot be found in the store.
	ErrClientNotFound = errors.New("oauthserver: client not found")

	// ErrClientRevoked is returned when a client has been revoked.
	ErrClientRevoked = errors.New("oauthserver: client has been revoked")

	// ErrInvalidScope is returned when a token does not have the required scope.
	ErrInvalidScope = errors.New("oauthserver: invalid scope(s) provided")

	// ErrUnauthenticated is returned when no valid token is present.
	ErrUnauthenticated = errors.New("oauthserver: unauthenticated")

	// ErrInvalidGrant is returned when an unsupported grant type is requested.
	ErrInvalidGrant = errors.New("oauthserver: invalid grant type")

	// ErrInvalidRequest is returned when a malformed request is received.
	ErrInvalidRequest = errors.New("oauthserver: invalid request")
)
