package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

// CacheOptions configures cache-control directives.
type CacheOptions struct {
	Public     bool
	Private    bool
	NoCache    bool
	NoStore    bool
	MaxAge     int // seconds, -1 to omit
	SMaxAge    int // seconds, -1 to omit
	ETag       bool
	MustRevali bool
}

// SetCacheHeaders adds Cache-Control (and optionally ETag) headers.
type SetCacheHeaders struct {
	opts CacheOptions
}

// NewSetCacheHeaders creates the middleware with the given options.
func NewSetCacheHeaders(opts CacheOptions) *SetCacheHeaders {
	return &SetCacheHeaders{opts: opts}
}

// Wrap returns an http.Handler that sets cache headers after the response is
// produced.
func (m *SetCacheHeaders) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		var directives []string

		if m.opts.Public {
			directives = append(directives, "public")
		}

		if m.opts.Private {
			directives = append(directives, "private")
		}

		if m.opts.NoCache {
			directives = append(directives, "no-cache")
		}

		if m.opts.NoStore {
			directives = append(directives, "no-store")
		}

		if m.opts.MaxAge >= 0 {
			directives = append(directives, fmt.Sprintf("max-age=%d", m.opts.MaxAge))
		}

		if m.opts.SMaxAge >= 0 {
			directives = append(directives, fmt.Sprintf("s-maxage=%d", m.opts.SMaxAge))
		}

		if m.opts.MustRevali {
			directives = append(directives, "must-revalidate")
		}

		if len(directives) > 0 {
			w.Header().Set("Cache-Control", strings.Join(directives, ", "))
		}
	})
}
