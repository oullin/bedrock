package routing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	nethttp "net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// UrlGenerator generates URLs from named routes and handles signed URLs.
type UrlGenerator struct {
	registry *Registry
	rootURL  string
	signKey  []byte
	defaults map[string]string
	request  *nethttp.Request
}

// NewUrlGenerator creates a URL generator.
func NewUrlGenerator(registry *Registry, rootURL string, signKey []byte) *UrlGenerator {
	return &UrlGenerator{
		registry: registry,
		rootURL:  strings.TrimRight(rootURL, "/"),
		signKey:  signKey,
		defaults: make(map[string]string),
	}
}

// To generates an absolute URL from a path.
func (g *UrlGenerator) To(path string, query ...url.Values) string {
	if isValidURL(path) {
		return path
	}

	base := g.rootURL + "/" + strings.TrimLeft(path, "/")

	if len(query) > 0 && len(query[0]) > 0 {
		base += "?" + query[0].Encode()
	}

	return base
}

// Secure generates an HTTPS URL.
func (g *UrlGenerator) Secure(path string) string {
	u := g.To(path)

	return strings.Replace(u, "http://", "https://", 1)
}

// Route generates a URL for a named route, substituting parameters.
func (g *UrlGenerator) Route(name string, params map[string]string, query ...url.Values) string {
	merged := make(map[string]string, len(g.defaults)+len(params))

	for k, v := range g.defaults {
		merged[k] = v
	}

	for k, v := range params {
		merged[k] = v
	}

	pattern := g.registry.URL(name, merged)

	if strings.HasPrefix(pattern, "#!") {
		return pattern
	}

	result := g.rootURL + pattern

	if len(query) > 0 && len(query[0]) > 0 {
		result += "?" + query[0].Encode()
	}

	return result
}

// Current returns the URL of the current request.
func (g *UrlGenerator) Current() string {
	if g.request == nil {
		return g.rootURL + "/"
	}

	return g.rootURL + g.request.URL.Path
}

// Previous returns the Referer URL or the fallback.
func (g *UrlGenerator) Previous(fallback ...string) string {
	if g.request != nil {
		ref := g.request.Header.Get("Referer")

		if ref != "" {
			return ref
		}
	}

	if len(fallback) > 0 {
		return fallback[0]
	}

	return g.rootURL + "/"
}

// PreviousPath returns the path component of the Referer URL.
func (g *UrlGenerator) PreviousPath(fallback ...string) string {
	prev := g.Previous(fallback...)

	parsed, err := url.Parse(prev)

	if err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		}

		return "/"
	}

	return parsed.Path
}

// Asset generates an asset URL.
func (g *UrlGenerator) Asset(path string) string {
	return g.rootURL + "/" + strings.TrimLeft(path, "/")
}

// SignedRoute generates a signed URL for a named route.
func (g *UrlGenerator) SignedRoute(name string, params map[string]string, expiration ...time.Duration) (string, error) {
	routeURL := g.Route(name, params)

	parsed, err := url.Parse(routeURL)

	if err != nil {
		return "", err
	}

	q := parsed.Query()

	if len(expiration) > 0 && expiration[0] != 0 {
		expires := time.Now().Add(expiration[0]).Unix()
		q.Set("expires", strconv.FormatInt(expires, 10))
	}

	parsed.RawQuery = q.Encode()

	sig := g.createSignature(parsed.Path, parsed.RawQuery)
	q.Set("signature", sig)
	parsed.RawQuery = q.Encode()

	return parsed.String(), nil
}

// TemporarySignedRoute generates a signed URL with an expiration.
func (g *UrlGenerator) TemporarySignedRoute(name string, expiration time.Duration, params map[string]string) (string, error) {
	return g.SignedRoute(name, params, expiration)
}

// HasValidSignature checks if a request's URL has a valid signature
// and has not expired.
func (g *UrlGenerator) HasValidSignature(req *nethttp.Request) bool {
	return g.HasCorrectSignature(req) && g.SignatureHasNotExpired(req)
}

// HasCorrectSignature checks if a request's URL signature is valid,
// ignoring expiration.
func (g *UrlGenerator) HasCorrectSignature(req *nethttp.Request) bool {
	signature := req.URL.Query().Get("signature")

	if signature == "" {
		return false
	}

	q := req.URL.Query()
	q.Del("signature")

	expected := g.createSignature(req.URL.Path, q.Encode())

	return hmac.Equal([]byte(signature), []byte(expected))
}

// SignatureHasNotExpired checks if a signed URL has not expired.
func (g *UrlGenerator) SignatureHasNotExpired(req *nethttp.Request) bool {
	expires := req.URL.Query().Get("expires")

	if expires == "" {
		return true
	}

	ts, err := strconv.ParseInt(expires, 10, 64)

	if err != nil {
		return false
	}

	return time.Now().Unix() <= ts
}

// SetRootURL updates the root URL.
func (g *UrlGenerator) SetRootURL(rootURL string) {
	g.rootURL = strings.TrimRight(rootURL, "/")
}

// SetRequest sets the current request for URL generation context.
func (g *UrlGenerator) SetRequest(req *nethttp.Request) {
	g.request = req
}

// SetDefaults sets default parameter values used in URL generation.
func (g *UrlGenerator) SetDefaults(defaults map[string]string) {
	g.defaults = defaults
}

func (g *UrlGenerator) createSignature(path, rawQuery string) string {
	data := path

	if rawQuery != "" {
		data += "?" + sortedQuery(rawQuery)
	}

	mac := hmac.New(sha256.New, g.signKey)
	mac.Write([]byte(data))

	return hex.EncodeToString(mac.Sum(nil))
}

func sortedQuery(raw string) string {
	values, err := url.ParseQuery(raw)

	if err != nil {
		return raw
	}

	keys := make([]string, 0, len(values))

	for k := range values {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var parts []string

	for _, k := range keys {
		for _, v := range values[k] {
			parts = append(parts, url.QueryEscape(k)+"="+url.QueryEscape(v))
		}
	}

	return strings.Join(parts, "&")
}

func isValidURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}
