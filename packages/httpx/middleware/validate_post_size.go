package middleware

import (
	"net/http"
	"strconv"
)

// ValidatePostSize rejects requests whose Content-Length exceeds a maximum.
type ValidatePostSize struct {
	maxBytes int64
}

// NewValidatePostSize creates the middleware with the given byte limit.
func NewValidatePostSize(maxBytes int64) *ValidatePostSize {
	return &ValidatePostSize{maxBytes: maxBytes}
}

// Wrap returns an http.Handler that checks Content-Length before delegating.
func (m *ValidatePostSize) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > m.maxBytes {
			http.Error(w, "Post body too large", http.StatusRequestEntityTooLarge)

			return
		}

		// Also check the Content-Length header string for cases where
		// ContentLength might not be set but the header is present.
		if cl := r.Header.Get("Content-Length"); cl != "" {
			if size, err := strconv.ParseInt(cl, 10, 64); err == nil && size > m.maxBytes {
				http.Error(w, "Post body too large", http.StatusRequestEntityTooLarge)

				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
