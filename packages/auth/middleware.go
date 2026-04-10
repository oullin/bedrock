package auth

import (
	"context"
	"net/http"
	"time"
)

type contextKey string

const userContextKey contextKey = "auth_user"

// WithUser stores an authenticated user in the request context.
func WithUser(r *http.Request, user Authenticatable) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), userContextKey, user))
}

// UserFromContext retrieves the authenticated user from the request context.
func UserFromContext(ctx context.Context) Authenticatable {
	u, _ := ctx.Value(userContextKey).(Authenticatable)

	return u
}

// EnsureAuthenticated rejects unauthenticated requests with 401.
func EnsureAuthenticated(guard Guard) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := guard.User(r.Context())
			if err != nil || user == nil {
				http.Error(w, "unauthenticated", http.StatusUnauthorized)

				return
			}

			next.ServeHTTP(w, WithUser(r, user))
		})
	}
}

// RedirectIfAuthenticated redirects authenticated users to the given path.
func RedirectIfAuthenticated(guard Guard, redirectTo string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if guard.Check(r.Context()) {
				http.Redirect(w, r, redirectTo, http.StatusFound)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// EnsureEmailIsVerified rejects users whose email is not verified.
func EnsureEmailIsVerified(guard Guard) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := guard.User(r.Context())
			if err != nil || user == nil {
				http.Error(w, "unauthenticated", http.StatusUnauthorized)

				return
			}

			if mv, ok := user.(MustVerifyEmail); ok && !mv.HasVerifiedEmail() {
				http.Error(w, "email not verified", http.StatusForbidden)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePassword redirects the user to the password confirmation page if they
// have not recently confirmed their password within the given timeout.
func RequirePassword(session SessionStore, timeout time.Duration, confirmPath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const key = "auth.password_confirmed_at"
			confirmedAt, _ := session.Get(key, int64(0)).(int64)
			threshold := time.Now().Add(-timeout).Unix()

			if confirmedAt < threshold {
				http.Redirect(w, r, confirmPath, http.StatusFound)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AuthenticateWithBasicAuth attempts HTTP Basic authentication on each request.
func AuthenticateWithBasicAuth(provider UserProvider, hasher PasswordHasher) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
				http.Error(w, "unauthorized", http.StatusUnauthorized)

				return
			}

			user, err := provider.RetrieveByCredentials(r.Context(), map[string]any{"email": username})
			if err != nil || user == nil || !hasher.Check(password, user.GetAuthPassword()) {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
				http.Error(w, "unauthorized", http.StatusUnauthorized)

				return
			}

			next.ServeHTTP(w, WithUser(r, user))
		})
	}
}
