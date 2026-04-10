package auth

import (
	"errors"
	"fmt"
)

// AuthenticationException is returned when authentication fails.
type AuthenticationException struct {
	Message      string
	Guards       []string
	RedirectPath string
}

// NewAuthenticationException creates an AuthenticationException.

// AuthorizationException is returned when authorization fails.
type AuthorizationException struct {
	Message    string
	StatusCode int
}

func (e *AuthenticationException) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "unauthenticated"
}

func NewAuthenticationException(guards []string, redirectPath string) *AuthenticationException {
	return &AuthenticationException{
		Message:      "unauthenticated",
		Guards:       guards,
		RedirectPath: redirectPath,
	}
}

func (e *AuthorizationException) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return fmt.Sprintf("this action is unauthorized (HTTP %d)", e.StatusCode)
}

// NewAuthorizationException creates an AuthorizationException.
func NewAuthorizationException(message string, statusCode int) *AuthorizationException {
	if statusCode == 0 {
		statusCode = 403
	}

	return &AuthorizationException{Message: message, StatusCode: statusCode}
}

var (
	// ErrInvalidProvider is returned when a user provider driver is not registered.
	ErrInvalidProvider = errors.New("auth: invalid user provider driver")
	// ErrInvalidGuard is returned when a guard driver is not registered.
	ErrInvalidGuard = errors.New("auth: invalid guard driver")
	// ErrUserNotFound is returned when no user matches the given credentials or ID.
	ErrUserNotFound = errors.New("auth: user not found")
)
