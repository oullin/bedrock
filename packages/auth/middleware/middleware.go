package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	auth "github.com/gollin/packages/auth"
)

type contextKey string

const (
	userContextKey    contextKey = "auth:user"
	sessionContextKey contextKey = "auth:session"
)

// CurrentUser returns the authenticated user from request context.
func CurrentUser(r *http.Request) (auth.Authenticatable, bool) {
	user, ok := r.Context().Value(userContextKey).(auth.Authenticatable)
	return user, ok
}

// CurrentSession returns the session from request context.
func CurrentSession(r *http.Request) (*auth.Session, bool) {
	session, ok := r.Context().Value(sessionContextKey).(*auth.Session)
	return session, ok
}

// Stack provides auth middleware helpers.
type Stack struct {
	Guard  auth.StatefulGuard
	Config auth.Config
	Clock  auth.Clock
}

// RequireAuthenticated rejects requests without an authenticated session.
func (s Stack) RequireAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, user, err := s.Guard.AuthenticateRequest(r.Context(), w, r)
		if err != nil || session == nil || session.PendingTwoFactor {
			WriteJSONError(w, http.StatusUnauthorized, auth.ErrUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		ctx = context.WithValue(ctx, sessionContextKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireVerifiedEmail rejects users with unverified email.
func (s Stack) RequireVerifiedEmail(next http.Handler) http.Handler {
	return s.RequireAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _ := CurrentUser(r)
		verifiable, ok := user.(auth.MustVerifyEmail)
		if !ok || !verifiable.HasVerifiedEmail() {
			WriteJSONError(w, http.StatusForbidden, auth.ErrEmailVerificationInvalid)
			return
		}
		next.ServeHTTP(w, r)
	}))
}

// RequirePasswordConfirmed rejects requests with stale password confirmation.
func (s Stack) RequirePasswordConfirmed(next http.Handler) http.Handler {
	return s.RequireAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := CurrentSession(r)
		if auth.ExpiredPasswordConfirmation(s.Config, session.PasswordConfirmedAt, s.Clock.Now()) {
			WriteJSONError(w, http.StatusForbidden, auth.ErrPasswordConfirmationRequired)
			return
		}
		next.ServeHTTP(w, r)
	}))
}

// WriteJSONError writes a standard auth error body.
func WriteJSONError(w http.ResponseWriter, status int, err error) {
	code := "auth_error"
	message := "Authentication failed."
	fields := map[string]string(nil)
	retryAfter := time.Duration(0)

	switch typed := err.(type) {
	case *auth.ValidationError:
		code = "validation_error"
		message = typed.Error()
		fields = typed.Fields
	case *auth.ThrottleError:
		code = "throttled"
		message = "Too many attempts."
		retryAfter = typed.RetryAfter
	default:
		switch {
		case errors.Is(err, auth.ErrUnauthorized):
			code = "unauthorized"
			message = "Authentication required."
		case errors.Is(err, auth.ErrInvalidCredentials):
			code = "invalid_credentials"
			message = "The provided credentials are incorrect."
		case errors.Is(err, auth.ErrInvalidToken):
			code = "invalid_token"
			message = "The token is invalid."
		case errors.Is(err, auth.ErrTokenExpired):
			code = "token_expired"
			message = "The token has expired."
		case errors.Is(err, auth.ErrTwoFactorRequired):
			code = "two_factor_required"
			message = "A two-factor challenge is required."
		case errors.Is(err, auth.ErrTwoFactorInvalid):
			code = "two_factor_invalid"
			message = "The two-factor challenge is invalid."
		case errors.Is(err, auth.ErrPasswordConfirmationRequired):
			code = "password_confirmation_required"
			message = "Password confirmation is required."
		case errors.Is(err, auth.ErrUserExists):
			code = "user_exists"
			message = "A user with that email already exists."
		case errors.Is(err, auth.ErrEmailVerificationInvalid):
			code = "email_verification_invalid"
			message = "The email verification link is invalid."
		default:
			message = err.Error()
		}
	}

	body := struct {
		Error      string            `json:"error"`
		Message    string            `json:"message"`
		Fields     map[string]string `json:"fields,omitempty"`
		RetryAfter int64             `json:"retryAfter,omitempty"`
	}{
		Error:   code,
		Message: message,
		Fields:  fields,
	}
	if retryAfter > 0 {
		body.RetryAfter = int64(retryAfter.Seconds())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
