package middleware

import (
	"net/http"
	"strings"

	"github.com/bedrock/packages/inertia/protocol"
)

// Config holds the configuration for the Inertia middleware.
type Config struct {
	// Version is the current asset version string. When the client's
	// X-Inertia-Version header does not match, the middleware responds
	// with 409 Conflict to trigger a full page reload.
	Version string
}

// New returns an HTTP middleware that implements the Inertia.js
// server-side protocol: Vary header, asset version checking, and
// redirect status conversion.
//
// The returned middleware has the standard signature
// func(http.Handler) http.Handler and is compatible with any
// router that accepts this pattern (chi, stdlib, alice, etc.).
func New(cfg Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Always set Vary so caches differentiate Inertia vs
			// regular responses for the same URL.
			w.Header().Set("Vary", protocol.HeaderInertia)

			if !protocol.IsInertiaRequest(r) {
				next.ServeHTTP(w, r)

				return
			}

			// Version check (GET requests only).
			if r.Method == http.MethodGet {
				clientVersion := r.Header.Get(protocol.HeaderVersion)

				if strings.TrimSpace(clientVersion) != "" && clientVersion != cfg.Version {
					w.Header().Set(protocol.HeaderLocation, r.RequestURI)

					w.WriteHeader(http.StatusConflict)

					return
				}
			}

			// Wrap the writer so 302 → 303 conversion happens
			// transparently for PUT/PATCH/DELETE.
			si := &statusInterceptor{
				ResponseWriter: w,
				request:        r,
			}

			next.ServeHTTP(si, r)
		})
	}
}
