package oauthserver

import (
	"context"
	"encoding/json"
	"net/http"
)

type contextKey string

const tokenUserKey contextKey = "passport_user"

// WithUser stores an authenticated user in the request context.
func WithUser(r *http.Request, user interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), tokenUserKey, user))
}

// UserFromContext retrieves the authenticated user from the request context.
// Returns nil when no user has been set.
func UserFromContext(ctx context.Context) interface{} {
	return ctx.Value(tokenUserKey)
}

// CheckToken builds middleware that requires ALL listed scopes to be present
// on the authenticated token.
//
//   - Returns 401 JSON when no valid token is present.
//   - Returns 403 JSON when one or more scopes are missing.
//   - Stores the authenticated user in context via WithUser.
func CheckToken(guard *TokenGuard, scopes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := guard.User(r.Context())

			if err != nil || user == nil {
				writeJSON(w, http.StatusUnauthorized, "Unauthenticated.")

				return
			}

			r = WithUser(r, user)

			// No scope requirements — pass straight through.
			if len(scopes) == 0 {
				next.ServeHTTP(w, r)

				return
			}

			uwt, ok := user.(*UserWithTokens)

			if !ok {
				writeJSON(w, http.StatusForbidden, "Invalid scope(s) provided.")

				return
			}

			for _, scope := range scopes {
				if uwt.TokenCant(scope) {
					writeJSON(w, http.StatusForbidden, "Invalid scope(s) provided.")

					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CheckTokenForAnyScope builds middleware that requires at least ONE of the
// listed scopes to be present on the authenticated token.
func CheckTokenForAnyScope(guard *TokenGuard, scopes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := guard.User(r.Context())

			if err != nil || user == nil {
				writeJSON(w, http.StatusUnauthorized, "Unauthenticated.")

				return
			}

			r = WithUser(r, user)

			if len(scopes) == 0 {
				next.ServeHTTP(w, r)

				return
			}

			uwt, ok := user.(*UserWithTokens)

			if !ok {
				writeJSON(w, http.StatusForbidden, "Invalid scope(s) provided.")

				return
			}

			for _, scope := range scopes {
				if uwt.TokenCan(scope) {
					next.ServeHTTP(w, r)

					return
				}
			}

			writeJSON(w, http.StatusForbidden, "Invalid scope(s) provided.")
		})
	}
}

// CheckClientCredentials builds middleware that validates client credentials
// tokens and optionally checks that the required scopes are present.
func CheckClientCredentials(guard *TokenGuard, scopes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			client, err := guard.Client(r.Context())

			if err != nil || client == nil {
				writeJSON(w, http.StatusUnauthorized, "Unauthenticated.")

				return
			}

			for _, scope := range scopes {
				ok, err := clientCredentialHasScope(r.Context(), guard, client, scope)

				if err != nil || !ok {
					writeJSON(w, http.StatusForbidden, "Invalid scope(s) provided.")

					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CheckClientCredentialsForAnyScope builds middleware that validates client
// credentials tokens and requires at least one listed scope to be present.
func CheckClientCredentialsForAnyScope(guard *TokenGuard, scopes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			client, err := guard.Client(r.Context())

			if err != nil || client == nil {
				writeJSON(w, http.StatusUnauthorized, "Unauthenticated.")

				return
			}

			if len(scopes) == 0 {
				next.ServeHTTP(w, r)

				return
			}

			for _, scope := range scopes {
				ok, err := clientCredentialHasScope(r.Context(), guard, client, scope)

				if err == nil && ok {
					next.ServeHTTP(w, r)

					return
				}
			}

			writeJSON(w, http.StatusForbidden, "Invalid scope(s) provided.")
		})
	}
}

func clientCredentialHasScope(ctx context.Context, guard *TokenGuard, client *Client, scope string) (bool, error) {
	token, err := guard.Token(ctx)

	if err != nil {
		return false, err
	}

	if token != nil {
		return token.Can(scope), nil
	}

	return client.HasScope(scope), nil
}

// CreateFreshApiToken builds middleware that issues a fresh encrypted cookie
// token for SPA (single-page application) requests.
//
// This implementation is a no-op stub; SPA cookie signing requires an
// encryption layer that callers inject via the OAuthServer config.
func CreateFreshApiToken(_ *TokenGuard, _ *OAuthServer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

// writeJSON writes a JSON response with the given HTTP status and message.
// Body format: {"message":"<msg>"}
func writeJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	body, _ := json.Marshal(map[string]string{"message": message})

	_, _ = w.Write(body)
}
