package auth

import (
	"context"
	"net/http"
)

type contextKey string

const (
	userContextKey    contextKey = "auth.user"
	sessionContextKey contextKey = "auth.session"
)

// AuthorizeFunc checks whether a user is authorized for an ability.
// Returns nil if authorized, an error otherwise.
type AuthorizeFunc func(ctx context.Context, user Authenticatable, ability string, arguments ...any) error

// WithUser stores the authenticated user in the request context.
func WithUser(ctx context.Context, user Authenticatable) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext retrieves the authenticated user from the request context.
func UserFromContext(ctx context.Context) (Authenticatable, bool) {
	user, ok := ctx.Value(userContextKey).(Authenticatable)
	return user, ok
}

// WithSession stores the session in the request context.
func WithSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, session)
}

// Authenticate is HTTP middleware that rejects unauthenticated requests with 401.
// When multiple guards are provided, the first one that successfully authenticates
// the request wins.
func Authenticate(guards ...*SessionGuard) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, guard := range guards {
				session, user, err := guard.AuthenticateRequest(r.Context(), w, r)
				if err == nil {
					ctx := WithUser(r.Context(), user)
					ctx = WithSession(ctx, session)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}

			guardNames := make([]string, len(guards))
			for i, guard := range guards {
				guardNames[i] = guard.Name()
			}

			http.Error(w, "Unauthenticated.", http.StatusUnauthorized)
		})
	}
}

// Authorize is HTTP middleware that checks a gate ability for the authenticated user.
func Authorize(authorizeFn AuthorizeFunc, ability string, arguments ...any) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthenticated.", http.StatusUnauthorized)
				return
			}

			if err := authorizeFn(r.Context(), user, ability, arguments...); err != nil {
				http.Error(w, "Forbidden.", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// EnsureEmailIsVerified is HTTP middleware that rejects users who have not verified their email.
func EnsureEmailIsVerified(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthenticated.", http.StatusUnauthorized)
			return
		}

		verifiable, ok := user.(MustVerifyEmail)
		if !ok || !verifiable.HasVerifiedEmail() {
			http.Error(w, "Your email address is not verified.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RedirectIfAuthenticated is HTTP middleware that redirects authenticated users.
func RedirectIfAuthenticated(guard *SessionGuard, redirectTo string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _, err := guard.AuthenticateRequest(r.Context(), w, r)
			if err == nil {
				http.Redirect(w, r, redirectTo, http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
