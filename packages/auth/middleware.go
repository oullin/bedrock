package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type contextKey string

const (
	userContextKey    contextKey = "auth:user"
	sessionContextKey contextKey = "auth:session"
)

// CurrentUser returns the authenticated user attached to a request.
func CurrentUser(r *http.Request) (*User, bool) {
	user, ok := r.Context().Value(userContextKey).(*User)
	return user, ok
}

// CurrentSession returns the current session attached to a request.
func CurrentSession(r *http.Request) (*Session, bool) {
	session, ok := r.Context().Value(sessionContextKey).(*Session)
	return session, ok
}

// RequireAuthenticated ensures that a request has an authenticated session.
func (m *Manager) RequireAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, user, err := m.Authenticate(r.Context(), w, r)
		if err != nil || session.PendingTwoFactor {
			WriteHTTPError(w, http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		ctx = context.WithValue(ctx, sessionContextKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireVerifiedEmail ensures that the authenticated user has a verified email.
func (m *Manager) RequireVerifiedEmail(next http.Handler) http.Handler {
	return m.RequireAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := CurrentUser(r)
		if user.EmailVerifiedAt == nil {
			WriteHTTPError(w, http.StatusForbidden, ErrEmailVerificationInvalid)
			return
		}
		next.ServeHTTP(w, r)
	}))
}

// RequirePasswordConfirmed ensures that the session has a recent password confirmation timestamp.
func (m *Manager) RequirePasswordConfirmed(next http.Handler) http.Handler {
	return m.RequireAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := CurrentSession(r)
		if session == nil || session.PasswordConfirmedAt == nil || m.clock.Now().After(session.PasswordConfirmedAt.Add(m.config.PasswordConfirmationTimeout)) {
			WriteHTTPError(w, http.StatusForbidden, ErrPasswordConfirmationRequired)
			return
		}
		next.ServeHTTP(w, r)
	}))
}

// WriteHTTPError writes a standard JSON auth error body.
func WriteHTTPError(w http.ResponseWriter, status int, err error) {
	code := "auth_error"
	message := "Authentication failed."
	fields := map[string]string(nil)
	retryAfter := time.Duration(0)

	switch typed := err.(type) {
	case *ValidationError:
		code = "validation_error"
		message = typed.Error()
		fields = typed.Fields
	case *ThrottleError:
		code = "throttled"
		message = "Too many attempts."
		retryAfter = typed.RetryAfter
	default:
		switch {
		case errors.Is(err, ErrUnauthorized):
			code = "unauthorized"
			message = "Authentication required."
		case errors.Is(err, ErrInvalidCredentials):
			code = "invalid_credentials"
			message = "The provided credentials are incorrect."
		case errors.Is(err, ErrInvalidToken):
			code = "invalid_token"
			message = "The token is invalid."
		case errors.Is(err, ErrTokenExpired):
			code = "token_expired"
			message = "The token has expired."
		case errors.Is(err, ErrTwoFactorRequired):
			code = "two_factor_required"
			message = "A two-factor challenge is required."
		case errors.Is(err, ErrTwoFactorInvalid):
			code = "two_factor_invalid"
			message = "The two-factor challenge is invalid."
		case errors.Is(err, ErrPasswordConfirmationRequired):
			code = "password_confirmation_required"
			message = "Password confirmation is required."
		case errors.Is(err, ErrUserExists):
			code = "user_exists"
			message = "A user with that email already exists."
		case errors.Is(err, ErrEmailVerificationInvalid):
			code = "email_verification_invalid"
			message = "The email verification link is invalid."
		default:
			message = err.Error()
		}
	}

	type errorBody struct {
		Error      string            `json:"error"`
		Message    string            `json:"message"`
		Fields     map[string]string `json:"fields,omitempty"`
		RetryAfter int64             `json:"retryAfter,omitempty"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	body := errorBody{
		Error:   code,
		Message: message,
		Fields:  fields,
	}
	if retryAfter > 0 {
		body.RetryAfter = int64(retryAfter.Seconds())
	}
	_ = json.NewEncoder(w).Encode(body)
}
