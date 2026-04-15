package websockets

import (
	"net/url"
	"strings"
)

// ValidateOrigin reports whether the given origin matches any pattern in the
// allowedOrigins list.
//
// Rules:
//   - An empty allowedOrigins slice means all origins are permitted.
//   - The pattern "*" matches any origin.
//   - Wildcard patterns of the form "*.host.tld" match any subdomain of host.tld.
//   - Exact string matches are also supported.
//
// The origin is compared as a hostname (scheme and port are stripped).
func ValidateOrigin(origin string, allowedOrigins []string) bool {
	if len(allowedOrigins) == 0 {
		return true
	}

	host := extractHost(origin)

	for _, pattern := range allowedOrigins {
		if pattern == "*" {
			return true
		}
		if matchPattern(host, pattern) {
			return true
		}
	}
	return false
}

// extractHost parses the hostname from an origin string.
// Falls back to the raw string when not a valid URL.
func extractHost(origin string) string {
	u, err := url.Parse(origin)
	if err != nil || u.Host == "" {
		return origin
	}
	// Strip port.
	host := u.Hostname()
	return host
}

// matchPattern checks whether host matches a single pattern.
// Supports exact matches and wildcard prefix patterns ("*.example.com").
func matchPattern(host, pattern string) bool {
	if strings.EqualFold(host, pattern) {
		return true
	}

	// Wildcard pattern: "*.example.com"
	if strings.HasPrefix(pattern, "*.") {
		suffix := pattern[1:] // ".example.com"
		return strings.HasSuffix(strings.ToLower(host), strings.ToLower(suffix))
	}

	return false
}
