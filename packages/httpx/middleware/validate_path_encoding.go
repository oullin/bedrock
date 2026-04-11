package middleware

import (
	"net/http"
	"net/url"
)

// ValidatePathEncoding rejects requests whose URI path contains invalid
// percent-encoding sequences.
type ValidatePathEncoding struct{}

// NewValidatePathEncoding creates the middleware.
func NewValidatePathEncoding() *ValidatePathEncoding {
	return &ValidatePathEncoding{}
}

// Wrap returns an http.Handler that validates path encoding.
func (m *ValidatePathEncoding) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := url.PathUnescape(r.RequestURI); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)

			return
		}

		next.ServeHTTP(w, r)
	})
}
