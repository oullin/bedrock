package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// CorsOptions configures Cross-Origin Resource Sharing behaviour.
type CorsOptions struct {
	AllowedOrigins   []string // "*" to allow all
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	MaxAge           int // preflight cache duration in seconds
	AllowCredentials bool
}

// HandleCors implements CORS preflight and simple request header injection.
type HandleCors struct {
	opts CorsOptions
}

// NewHandleCors creates the CORS middleware.
func NewHandleCors(opts CorsOptions) *HandleCors {
	return &HandleCors{opts: opts}
}

// Wrap returns an http.Handler that handles CORS.
func (m *HandleCors) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "" {
			next.ServeHTTP(w, r)

			return
		}

		if !m.isAllowedOrigin(origin) {
			next.ServeHTTP(w, r)

			return
		}

		// Set CORS headers.
		m.setOriginHeader(w, origin)

		if len(m.opts.ExposedHeaders) > 0 {
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(m.opts.ExposedHeaders, ", "))
		}

		if m.opts.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight.
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			m.handlePreflight(w, r)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *HandleCors) handlePreflight(w http.ResponseWriter, r *http.Request) {
	if len(m.opts.AllowedMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.opts.AllowedMethods, ", "))
	}

	requestHeaders := r.Header.Get("Access-Control-Request-Headers")

	if requestHeaders != "" {
		if m.allowsAllHeaders() {
			w.Header().Set("Access-Control-Allow-Headers", requestHeaders)
		} else {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(m.opts.AllowedHeaders, ", "))
		}
	}

	if m.opts.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(m.opts.MaxAge))
	}

	w.WriteHeader(http.StatusNoContent)
}

func (m *HandleCors) isAllowedOrigin(origin string) bool {
	for _, allowed := range m.opts.AllowedOrigins {
		if allowed == "*" {
			return true
		}

		if strings.EqualFold(allowed, origin) {
			return true
		}
	}

	return false
}

func (m *HandleCors) setOriginHeader(w http.ResponseWriter, origin string) {
	for _, allowed := range m.opts.AllowedOrigins {
		if allowed == "*" && !m.opts.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Origin", "*")

			return
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Add("Vary", "Origin")
}

func (m *HandleCors) allowsAllHeaders() bool {
	for _, h := range m.opts.AllowedHeaders {
		if h == "*" {
			return true
		}
	}

	return false
}
