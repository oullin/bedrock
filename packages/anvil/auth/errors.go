package auth

import (
	"errors"
	"fmt"
	"strings"
)

// AuthenticationException mirrors Laravel's authentication exception intent.
type AuthenticationException struct {
	Message    string
	Guards     []string
	RedirectTo string
}

var (
	ErrUnauthorized       = errors.New("auth: unauthorized")
	ErrUserNotFound       = errors.New("auth: user not found")
	ErrUserExists         = errors.New("auth: user already exists")
	ErrInvalidCredentials = errors.New("auth: invalid credentials")
	ErrInvalidToken       = errors.New("auth: invalid token")
	ErrTokenExpired       = errors.New("auth: token expired")
	ErrThrottled          = errors.New("auth: too many requests")
)

// Error implements error.
func (e AuthenticationException) Error() string {
	message := strings.TrimSpace(e.Message)

	if message == "" {
		message = "auth: authentication failed"
	}

	if len(e.Guards) == 0 {
		return message
	}

	return fmt.Sprintf("%s (guards=%s)", message, strings.Join(e.Guards, ","))
}
