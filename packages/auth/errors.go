package auth

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

var (
	// ErrUnauthorized indicates that the request is missing a valid authenticated session.
	ErrUnauthorized = errors.New("auth: unauthorized")
	// ErrUserNotFound indicates that no user matched the requested identifier.
	ErrUserNotFound = errors.New("auth: user not found")
	// ErrUserExists indicates that a new user conflicts with an existing identifier.
	ErrUserExists = errors.New("auth: user already exists")
	// ErrInvalidCredentials indicates that an identifier or password pair is invalid.
	ErrInvalidCredentials = errors.New("auth: invalid credentials")
	// ErrInvalidToken indicates that a signed or reset token is invalid.
	ErrInvalidToken = errors.New("auth: invalid token")
	// ErrTokenExpired indicates that a signed or reset token has expired.
	ErrTokenExpired = errors.New("auth: token expired")
	// ErrTwoFactorRequired indicates that a login requires an additional two-factor challenge.
	ErrTwoFactorRequired = errors.New("auth: two-factor challenge required")
	// ErrTwoFactorInvalid indicates that a two-factor code or recovery code is invalid.
	ErrTwoFactorInvalid = errors.New("auth: invalid two-factor challenge")
	// ErrPasswordConfirmationRequired indicates that the user must confirm their password first.
	ErrPasswordConfirmationRequired = errors.New("auth: password confirmation required")
	// ErrEmailVerificationInvalid indicates that an email verification link is invalid.
	ErrEmailVerificationInvalid = errors.New("auth: invalid email verification link")
)

// ValidationError captures field-level validation failures.
type ValidationError struct {
	Fields map[string]string
}

// Error implements error.
func (e *ValidationError) Error() string {
	if len(e.Fields) == 0 {
		return "auth: validation failed"
	}

	keys := make([]string, 0, len(e.Fields))
	for key := range e.Fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", key, e.Fields[key]))
	}

	return "auth: validation failed: " + strings.Join(parts, ", ")
}

// ThrottleError reports that a caller exceeded an auth throttle.
type ThrottleError struct {
	Scope      string
	RetryAfter time.Duration
}

// Error implements error.
func (e *ThrottleError) Error() string {
	return fmt.Sprintf("auth: throttled %s", e.Scope)
}
