package auth

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// ErrUnauthorized indicates the request does not have a valid session.

// ErrUserNotFound indicates the requested user does not exist.

// ErrUserExists indicates the target user already exists.

// ErrInvalidCredentials indicates the submitted credentials are invalid.

// ErrInvalidToken indicates the submitted token is invalid.

// ErrTokenExpired indicates the submitted token is expired.

// ErrTwoFactorRequired indicates the login must complete a two-factor challenge.

// ErrTwoFactorInvalid indicates the two-factor challenge failed.

// ErrPasswordConfirmationRequired indicates the user must confirm their password.

// ErrEmailVerificationInvalid indicates an email verification link is invalid.

// ValidationError reports field-level validation failures.
type ValidationError struct {
	Fields map[string]string
}

// Error implements error.

// ThrottleError reports a fixed-window throttle rejection.
type ThrottleError struct {
	Scope      string
	RetryAfter time.Duration
}

var (
	ErrUnauthorized = errors.New("auth: unauthorized")

	ErrUserNotFound = errors.New("auth: user not found")

	ErrUserExists = errors.New("auth: user already exists")

	ErrInvalidCredentials = errors.New("auth: invalid credentials")

	ErrHasherNotConfigured = errors.New("auth: password hasher is not configured")

	ErrInvalidToken = errors.New("auth: invalid token")

	ErrTokenExpired = errors.New("auth: token expired")

	ErrTwoFactorRequired = errors.New("auth: two-factor challenge required")

	ErrTwoFactorInvalid = errors.New("auth: invalid two-factor challenge")

	ErrPasswordConfirmationRequired = errors.New("auth: password confirmation required")

	ErrEmailVerificationInvalid = errors.New("auth: invalid email verification link")
)

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

// Error implements error.
func (e *ThrottleError) Error() string {
	return fmt.Sprintf("auth: throttled %s", e.Scope)
}
