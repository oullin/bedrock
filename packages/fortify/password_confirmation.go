package fortify

import (
	"context"
	"net/http"
	"time"
)

type contextKey string

const passwordConfirmedAtKey contextKey = "fortify.password_confirmed_at"

// SetPasswordConfirmedAt marks the password as confirmed in the request context.
// In practice, the consuming app should persist this to the session. This
// context value serves as the in-request marker.
func SetPasswordConfirmedAt(r *http.Request) {
	ctx := context.WithValue(r.Context(), passwordConfirmedAtKey, time.Now())
	*r = *r.WithContext(ctx)
}

// GetPasswordConfirmedAt retrieves the password confirmation timestamp from context.
func GetPasswordConfirmedAt(r *http.Request) *time.Time {
	v, ok := r.Context().Value(passwordConfirmedAtKey).(time.Time)
	if !ok {
		return nil
	}

	return &v
}

// PasswordConfirmedAtFromContext retrieves the confirmation time from context.
// This is used by the middleware when the app stores the value in context
// (e.g., loaded from session during request lifecycle).
func PasswordConfirmedAtFromContext(ctx context.Context) *time.Time {
	v, ok := ctx.Value(passwordConfirmedAtKey).(time.Time)
	if !ok {
		return nil
	}

	return &v
}

// WithPasswordConfirmedAt returns a context with the confirmation timestamp set.
// The consuming app uses this to inject the session-stored value into context.
func WithPasswordConfirmedAt(ctx context.Context, at time.Time) context.Context {
	return context.WithValue(ctx, passwordConfirmedAtKey, at)
}

// EnsurePasswordIsConfirmed returns middleware that requires recent password confirmation.
func EnsurePasswordIsConfirmed(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			confirmedAt := PasswordConfirmedAtFromContext(r.Context())

			if confirmedAt == nil || time.Since(*confirmedAt) > timeout {
				http.Error(w, "password confirmation required", http.StatusLocked)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
