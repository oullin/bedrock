package api

import "net/http"

// AuthCallback decides whether an HTTP request is allowed to reach a Horizon
// dashboard handler. Laravel Horizon exposes the same extension point via
// Horizon::auth(fn ($request) => bool).
type AuthCallback func(*http.Request) bool

// AllowAll permits every request to reach the dashboard; Horizon uses this in
// local environments.
func AllowAll(*http.Request) bool { return true }

// RequireAuth gates HTTP handling behind an AuthCallback. Failed checks return
// 403 Forbidden — the same status Laravel Horizon emits when
// Gate::check('viewHorizon') fails.
func RequireAuth(callback AuthCallback) func(http.Handler) http.Handler {
	if callback == nil {
		callback = AllowAll
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !callback(r) {
				http.Error(w, "Forbidden", http.StatusForbidden)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
