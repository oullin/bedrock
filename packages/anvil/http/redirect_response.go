package http

import (
	nethttp "net/http"
	"net/url"
	"strings"
)

// RedirectResponse represents an HTTP redirect with status code, target URL,
// headers, and optional URL fragment. It mirrors Upstream's RedirectResponse.
type RedirectResponse struct {
	statusCode int
	targetUrl  string
	fragment   string
	headers    nethttp.Header
}

// NewRedirectResponse creates a redirect response to the given URL with the
// specified status code.
func NewRedirectResponse(targetUrl string, status int) *RedirectResponse {
	return &RedirectResponse{
		statusCode: status,
		targetUrl:  targetUrl,
		headers:    nethttp.Header{},
	}
}

// GetTargetUrl returns the redirect target URL including any fragment.
func (rr *RedirectResponse) GetTargetUrl() string {
	u := rr.targetUrl
	if rr.fragment != "" {
		// Strip any existing fragment before appending the new one.
		if i := strings.IndexByte(u, '#'); i >= 0 {
			u = u[:i]
		}
		u += "#" + rr.fragment
	}
	return u
}

// SetTargetUrl sets the redirect target URL.
func (rr *RedirectResponse) SetTargetUrl(target string) {
	rr.targetUrl = target
}

// GetStatusCode returns the HTTP status code.
func (rr *RedirectResponse) GetStatusCode() int {
	return rr.statusCode
}

// SetStatusCode sets the HTTP status code.
func (rr *RedirectResponse) SetStatusCode(code int) {
	rr.statusCode = code
}

// SetHeader sets a single header value.
func (rr *RedirectResponse) SetHeader(key, value string) *RedirectResponse {
	rr.headers.Set(key, value)
	return rr
}

// GetHeader returns a header value.
func (rr *RedirectResponse) GetHeader(key string) string {
	return rr.headers.Get(key)
}

// WithFragment sets the URL fragment (without the leading "#").
func (rr *RedirectResponse) WithFragment(fragment string) *RedirectResponse {
	rr.fragment = fragment
	return rr
}

// WithoutFragment removes any URL fragment.
func (rr *RedirectResponse) WithoutFragment() *RedirectResponse {
	rr.fragment = ""
	return rr
}

// SameOriginOption configures same-origin enforcement.
type SameOriginOption func(*sameOriginConfig)

type sameOriginConfig struct {
	checkScheme bool
	checkPort   bool
}

// WithoutSchemeValidation disables scheme comparison during same-origin checks.
func WithoutSchemeValidation() SameOriginOption {
	return func(c *sameOriginConfig) {
		c.checkScheme = false
	}
}

// WithoutPortValidation disables port comparison during same-origin checks.
func WithoutPortValidation() SameOriginOption {
	return func(c *sameOriginConfig) {
		c.checkPort = false
	}
}

// IsSameOrigin reports whether the redirect target has the same origin as
// the given request. By default it compares scheme, hostname, and port.
func (rr *RedirectResponse) IsSameOrigin(r *Request, opts ...SameOriginOption) bool {
	cfg := sameOriginConfig{checkScheme: true, checkPort: true}
	for _, o := range opts {
		o(&cfg)
	}

	targetURL, err := url.Parse(rr.GetTargetUrl())
	if err != nil {
		return false
	}

	reqScheme := Scheme(r)
	reqHost := r.Host

	// Compare scheme.
	if cfg.checkScheme {
		if !strings.EqualFold(targetURL.Scheme, reqScheme) {
			return false
		}
	}

	// Split host and port.
	targetHostname := targetURL.Hostname()
	targetPort := targetURL.Port()
	reqHostname, reqPort := splitHostPort(reqHost)

	// Default ports. When scheme validation is disabled, use the request
	// scheme for both sides so that default-port differences caused solely
	// by differing schemes do not trigger a mismatch.
	portScheme := reqScheme
	if cfg.checkScheme {
		if targetPort == "" {
			targetPort = defaultPort(targetURL.Scheme)
		}
	} else if targetPort == "" {
		targetPort = defaultPort(portScheme)
	}
	if reqPort == "" {
		reqPort = defaultPort(portScheme)
	}

	// Compare hostname.
	if !strings.EqualFold(targetHostname, reqHostname) {
		return false
	}

	// Compare port.
	if cfg.checkPort && targetPort != reqPort {
		return false
	}

	return true
}

// WriteTo writes the redirect response to the given ResponseWriter.
func (rr *RedirectResponse) WriteTo(w ResponseWriter, r *Request) {
	for k, vals := range rr.headers {
		for _, v := range vals {
			w.Header().Add(k, v)
		}
	}

	nethttp.Redirect(w, r, rr.GetTargetUrl(), rr.statusCode)
}

func splitHostPort(host string) (hostname, port string) {
	if i := strings.LastIndexByte(host, ':'); i >= 0 {
		// Check it's not an IPv6 bracket.
		if i > 0 && host[0] == '[' {
			if j := strings.IndexByte(host, ']'); j > i {
				return host, ""
			}
		}
		return host[:i], host[i+1:]
	}
	return host, ""
}

func defaultPort(scheme string) string {
	switch strings.ToLower(scheme) {
	case "https":
		return "443"
	case "http":
		return "80"
	default:
		return ""
	}
}
