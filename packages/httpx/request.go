package httpx

import (
	"net/http"
	"net/url"
	"strings"
)

// SessionStore is the minimal session interface that httpx needs. Types such as
// session.Store satisfy it implicitly via Go structural typing.
type SessionStore interface {
	Get(key string, fallback any) any
	Put(key string, value any)
	Flash(key string, value any)
	GetOldInput(key string, fallback any) any
	HasOldInput(key string) bool
	FlashInput(values map[string]any)
	Remove(key string) any
}

// RouteResolver provides route information for the current request.
type RouteResolver interface {
	CurrentRouteName() string
	CurrentRouteAction() string
}

// Request wraps *http.Request with accessor methods for input,
// headers, content negotiation, flash data and more. It can be constructed from
// any *http.Request and optionally enriched with a session or route resolver.
type Request struct {
	raw     *http.Request
	session SessionStore
	route   RouteResolver

	// parsedInput caches the merged query+body input after the first call to
	// All(). It is lazily populated.
	parsedInput map[string]any
}

// NewRequest wraps an *http.Request.
func NewRequest(r *http.Request) *Request {
	return &Request{raw: r}
}

// SetSession attaches a session store to the request. This enables flash-data
// and old-input helpers.
func (r *Request) SetSession(s SessionStore) {
	r.session = s
}

// Session returns the attached session store, or nil.
func (r *Request) Session() SessionStore {
	return r.session
}

// SetRouteResolver attaches a route resolver to the request.
func (r *Request) SetRouteResolver(rr RouteResolver) {
	r.route = rr
}

// RouteResolver returns the attached resolver, or nil.
func (r *Request) RouteResolver() RouteResolver {
	return r.route
}

// Raw returns the underlying *http.Request.
func (r *Request) Raw() *http.Request {
	return r.raw
}

// Method returns the HTTP method (GET, POST, ...).
func (r *Request) Method() string {
	return r.raw.Method
}

// IsMethod checks whether the request method matches (case-insensitive).
func (r *Request) IsMethod(method string) bool {
	return strings.EqualFold(r.raw.Method, method)
}

// URL returns the request URL path and query string (e.g. "/foo?bar=1").
func (r *Request) URL() string {
	u := r.raw.URL.Path

	if r.raw.URL.RawQuery != "" {
		u += "?" + r.raw.URL.RawQuery
	}

	return u
}

// FullURL returns the full request URL including scheme and host.
func (r *Request) FullURL() string {
	scheme := "http"

	if r.raw.TLS != nil {
		scheme = "https"
	}

	if fwd := r.raw.Header.Get("X-Forwarded-Proto"); fwd != "" {
		scheme = fwd
	}

	return scheme + "://" + r.raw.Host + r.URL()
}

// Path returns the request URI path without the query string.
func (r *Request) Path() string {
	return r.raw.URL.Path
}

// PathInfo returns the request path used by routing validators.
func (r *Request) PathInfo() string {
	return r.Path()
}

// DecodedPath returns the URL-decoded request path used for route parameter binding.
func (r *Request) DecodedPath() string {
	path := r.Path()

	if decoded, err := url.PathUnescape(path); err == nil {
		return decoded
	}

	return path
}

// Segment returns the 1-indexed URI segment (split by /). Returns fallback if
// the index is out of range.
func (r *Request) Segment(index int, fallback ...string) string {
	segments := strings.Split(strings.Trim(r.raw.URL.Path, "/"), "/")

	if index < 1 || index > len(segments) {
		if len(fallback) > 0 {
			return fallback[0]
		}

		return ""
	}

	return segments[index-1]
}

// Segments returns all URI segments.
func (r *Request) Segments() []string {
	trimmed := strings.Trim(r.raw.URL.Path, "/")

	if trimmed == "" {
		return nil
	}

	return strings.Split(trimmed, "/")
}

// Is checks whether the request path matches any of the given patterns.
// Patterns support wildcard (*) matching.
func (r *Request) Is(patterns ...string) bool {
	path := r.Path()

	for _, pattern := range patterns {
		if matchPath(pattern, path) {
			return true
		}
	}

	return false
}

// IP returns the client's IP address.
func (r *Request) IP() string {
	if forwarded := r.raw.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.SplitN(forwarded, ",", 2)

		return strings.TrimSpace(parts[0])
	}

	host := r.raw.RemoteAddr

	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}

	return host
}

// IPs returns all client IP addresses from the X-Forwarded-For chain.
func (r *Request) IPs() []string {
	forwarded := r.raw.Header.Get("X-Forwarded-For")

	if forwarded == "" {
		return []string{r.IP()}
	}

	parts := strings.Split(forwarded, ",")
	ips := make([]string, 0, len(parts))

	for _, p := range parts {
		if ip := strings.TrimSpace(p); ip != "" {
			ips = append(ips, ip)
		}
	}

	return ips
}

// UserAgent returns the User-Agent header value.
func (r *Request) UserAgent() string {
	return r.raw.UserAgent()
}

// Secure returns true when the request was made over HTTPS.
func (r *Request) Secure() bool {
	if r.raw.TLS != nil {
		return true
	}

	return strings.EqualFold(r.raw.Header.Get("X-Forwarded-Proto"), "https")
}

// SchemeAndHost returns the scheme and host portion (e.g. "https://example.com").
func (r *Request) SchemeAndHost() string {
	scheme := "http"

	if r.Secure() {
		scheme = "https"
	}

	return scheme + "://" + r.raw.Host
}

// Host returns the request host (may include port).
func (r *Request) Host() string {
	return r.raw.Host
}

// Fingerprint returns a unique string identifying the request by its route,
// IP, and user agent.
func (r *Request) Fingerprint() string {
	parts := []string{r.Method(), r.raw.Host, r.Path(), r.IP(), r.UserAgent()}

	if r.route != nil {
		parts = append([]string{r.route.CurrentRouteName()}, parts...)
	}

	return strings.Join(parts, "|")
}

// QueryString returns the raw query string.
func (r *Request) QueryString() string {
	return r.raw.URL.RawQuery
}

// QueryValues returns the parsed query parameters.
func (r *Request) QueryValues() url.Values {
	return r.raw.URL.Query()
}

// matchPath matches a URI path against a pattern with * wildcards.
func matchPath(pattern, path string) bool {
	if pattern == "*" {
		return true
	}

	patParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	// A trailing wildcard can match any number of remaining segments.
	trailingWild := len(patParts) > 0 && patParts[len(patParts)-1] == "*"

	if !trailingWild && len(patParts) != len(pathParts) {
		return false
	}

	if trailingWild && len(pathParts) < len(patParts)-1 {
		return false
	}

	// Check all non-trailing-wildcard segments.
	limit := len(patParts)

	if trailingWild {
		limit = len(patParts) - 1
	}

	for i := 0; i < limit; i++ {
		if i >= len(pathParts) {
			return false
		}

		if patParts[i] == "*" {
			continue
		}

		if patParts[i] != pathParts[i] {
			return false
		}
	}

	return true
}
