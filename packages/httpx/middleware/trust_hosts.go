package middleware

import (
	"net/http"
	"strings"
)

// TrustHosts validates the Host header against a list of allowed hosts.
// Requests with untrusted Host values receive a 400 Bad Request.
type TrustHosts struct {
	hosts []string
}

// NewTrustHosts creates the middleware with the given allowed hosts. Each host
// may be an exact match or a wildcard pattern like "*.example.com".
func NewTrustHosts(hosts ...string) *TrustHosts {
	return &TrustHosts{hosts: hosts}
}

// Wrap returns an http.Handler that validates the Host header.
func (m *TrustHosts) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.isTrustedHost(r.Host) {
			http.Error(w, "Bad Request", http.StatusBadRequest)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *TrustHosts) isTrustedHost(host string) bool {
	// Strip port if present.
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}

	host = strings.ToLower(host)

	for _, allowed := range m.hosts {
		allowed = strings.ToLower(allowed)

		if allowed == host {
			return true
		}

		// Wildcard matching: "*.example.com" matches "foo.example.com".
		if strings.HasPrefix(allowed, "*.") {
			suffix := allowed[1:] // ".example.com"

			if strings.HasSuffix(host, suffix) {
				return true
			}
		}
	}

	return false
}
