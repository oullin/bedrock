package auth

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

var (
	// ErrUnauthorized indicates the request does not have a valid session.
	ErrUnauthorized = errors.New("auth: unauthorized")
	// ErrUserNotFound indicates the requested user does not exist.
	ErrUserNotFound = errors.New("auth: user not found")
	// ErrUserExists indicates the target user already exists.
	ErrUserExists = errors.New("auth: user already exists")
	// ErrInvalidCredentials indicates the submitted credentials are invalid.
	ErrInvalidCredentials = errors.New("auth: invalid credentials")
	// ErrInvalidToken indicates the submitted token is invalid.
	ErrInvalidToken = errors.New("auth: invalid token")
	// ErrTokenExpired indicates the submitted token is expired.
	ErrTokenExpired = errors.New("auth: token expired")
	// ErrTwoFactorRequired indicates the login must complete a two-factor challenge.
	ErrTwoFactorRequired = errors.New("auth: two-factor challenge required")
	// ErrTwoFactorInvalid indicates the two-factor challenge failed.
	ErrTwoFactorInvalid = errors.New("auth: invalid two-factor challenge")
	// ErrPasswordConfirmationRequired indicates the user must confirm their password.
	ErrPasswordConfirmationRequired = errors.New("auth: password confirmation required")
	// ErrEmailVerificationInvalid indicates an email verification link is invalid.
	ErrEmailVerificationInvalid = errors.New("auth: invalid email verification link")
)

// ValidationError reports field-level validation failures.
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

// ThrottleError reports a fixed-window throttle rejection.
type ThrottleError struct {
	Scope      string
	RetryAfter time.Duration
}

// Error implements error.
func (e *ThrottleError) Error() string {
	return fmt.Sprintf("auth: throttled %s", e.Scope)
}
