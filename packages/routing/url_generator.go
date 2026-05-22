package routing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bedrock/packages/routing/exceptions"
)

// Ref: @bedrock/code-0346
// It owns the request, the route collection (for name lookups), and the
// signing key. The four big jobs:
//
//  1. Build absolute or relative URLs for arbitrary paths ([UrlGenerator.To]).
//  2. Build URLs for named routes ([UrlGenerator.Route]).
//  3. Sign and verify route URLs ([UrlGenerator.SignedRoute],
//     [UrlGenerator.HasValidSignature]).
//  4. Generate asset paths ([UrlGenerator.Asset]).
type UrlGenerator struct {
	routes        RouteCollectionInterface
	request       URLRequest
	assetRoot     string
	forcedScheme  string
	forcedRootUrl string
	key           string
}

// NewUrlGenerator constructs a generator bound to a request and the route
// collection used for name lookups.
func NewUrlGenerator(routes RouteCollectionInterface, request URLRequest, assetRoot string) *UrlGenerator {
	return &UrlGenerator{routes: routes, request: request, assetRoot: assetRoot}
}

// SetKeyResolver sets the signing key (used for signed URLs).
func (u *UrlGenerator) SetKeyResolver(key string) *UrlGenerator { u.key = key; return u }

// Ref: @bedrock/code-0346
func (u *UrlGenerator) ForceScheme(scheme string) { u.forcedScheme = scheme }

// ForceHttps is the boolean shorthand for [UrlGenerator.ForceScheme]("https").
func (u *UrlGenerator) ForceHttps(force bool) {
	if force {
		u.forcedScheme = "https"
	} else {
		u.forcedScheme = ""
	}
}

// ForceRootUrl pins the URL root used for absolute generation.
// UrlGenerator::forceRootUrl.
func (u *UrlGenerator) ForceRootUrl(root string) { u.forcedRootUrl = root }

// =====================================================================
// Path helpers
// =====================================================================

// Full returns the full URL of the current request.
func (u *UrlGenerator) Full() string {
	if u.request == nil {
		return ""
	}

	return u.request.URL()
}

// Current returns the current request URL without the query string.
func (u *UrlGenerator) Current() string {
	if u.request == nil {
		return ""
	}

	full := u.request.URL()

	if idx := strings.Index(full, "?"); idx >= 0 {
		return full[:idx]
	}

	return full
}

// To produces an absolute URL for the supplied path. extra is appended as
// path segments (URL-encoded). When secure is true the URL forces HTTPS.
//
// Ref: @bedrock/code-0346
func (u *UrlGenerator) To(path string, extra []string, secure *bool) string {
	if isAbsoluteURL(path) {
		return path
	}

	root := u.formatRoot(u.formatScheme(secure), "")
	tail := strings.TrimLeft(path, "/")

	for _, e := range extra {
		tail = strings.TrimRight(tail, "/") + "/" + url.PathEscape(e)
	}

	return root + "/" + tail
}

// Secure is a shortcut for [UrlGenerator.To] with secure=true.
func (u *UrlGenerator) Secure(path string, extra []string) string {
	t := true

	return u.To(path, extra, &t)
}

// Asset returns the URL of an asset relative to the asset root.
func (u *UrlGenerator) Asset(path string, secure *bool) string {
	if isAbsoluteURL(path) {
		return path
	}

	root := u.assetRoot

	if root == "" {
		root = u.formatRoot(u.formatScheme(secure), "")
	}

	return strings.TrimRight(root, "/") + "/" + strings.TrimLeft(path, "/")
}

// SecureAsset is a shortcut for [UrlGenerator.Asset] with secure=true.
func (u *UrlGenerator) SecureAsset(path string) string {
	t := true

	return u.Asset(path, &t)
}

func (u *UrlGenerator) formatScheme(secure *bool) string {
	if u.forcedScheme != "" {
		return u.forcedScheme
	}

	if secure != nil {
		if *secure {
			return "https"
		}

		return "http"
	}

	if u.request != nil {
		return u.request.Scheme()
	}

	return "http"
}

func (u *UrlGenerator) formatRoot(scheme, root string) string {
	if u.forcedRootUrl != "" {
		root = u.forcedRootUrl
	}

	if root == "" && u.request != nil {
		root = scheme + "://" + u.request.Host()
	}

	if root == "" {
		root = scheme + "://"
	}
	// Replace existing scheme.
	if idx := strings.Index(root, "://"); idx >= 0 {
		root = scheme + root[idx:]
	} else {
		root = scheme + "://" + root
	}

	return strings.TrimRight(root, "/")
}

// =====================================================================
// Named route URL
// =====================================================================

// Route generates a URL for the named route.
//
// Ref: @bedrock/code-0346
func (u *UrlGenerator) Route(name string, parameters map[string]any, absolute bool) (string, error) {
	route := u.routes.GetByName(name)

	if route == nil {
		return "", fmt.Errorf("route [%s] not defined", name)
	}

	return u.ToRoute(route, parameters, absolute)
}

// Ref: @bedrock/code-0346
func (u *UrlGenerator) ToRoute(route *Route, parameters map[string]any, absolute bool) (string, error) {
	parameters = mergeRouteDefaults(route, parameters)

	for _, name := range route.ParameterNames() {
		if _, hasParam := parameters[name]; !hasParam {
			if !route.HasDefault(name) && !isOptionalParam(route.Uri, name) {
				return "", exceptions.ForMissingParameters(route.GetName(), []string{name})
			}
		}
	}

	gen := NewRouteUrlGenerator(u, u.request)

	return gen.To(route, parameters, absolute), nil
}

func mergeRouteDefaults(route *Route, parameters map[string]any) map[string]any {
	if parameters == nil {
		parameters = map[string]any{}
	}

	if len(route.DefaultValues) == 0 {
		return parameters
	}

	merged := make(map[string]any, len(route.DefaultValues)+len(parameters))

	for k, v := range route.DefaultValues {
		merged[k] = v
	}

	for k, v := range parameters {
		merged[k] = v
	}

	return merged
}

func isOptionalParam(uri, name string) bool {
	return strings.Contains(uri, "{"+name+"?}")
}

// =====================================================================
// Signed URLs
// =====================================================================

// SignedRoute produces a signed URL for the named route.
//
// expiration may be 0 (no expiration) or a positive number of seconds from
// now until the signature must be considered expired.
// UrlGenerator::signedRoute.
func (u *UrlGenerator) SignedRoute(name string, parameters map[string]any, expiration int64, absolute bool) (string, error) {
	if parameters == nil {
		parameters = map[string]any{}
	}

	if _, ok := parameters["signature"]; ok {
		return "", fmt.Errorf("\"signature\" is a reserved parameter for signed routes")
	}

	if _, ok := parameters["expires"]; ok {
		return "", fmt.Errorf("\"expires\" is a reserved parameter for signed routes")
	}

	if expiration > 0 {
		parameters["expires"] = strconv.FormatInt(time.Now().Unix()+expiration, 10)
	}
	// Build the canonical URL (without the signature) for HMAC input.
	base, err := u.Route(name, parameters, absolute)

	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, []byte(u.key))
	mac.Write([]byte(base))
	signature := hex.EncodeToString(mac.Sum(nil))
	parameters["signature"] = signature

	return u.Route(name, parameters, absolute)
}

// TemporarySignedRoute is the parity-named alias matching the upstream two-arg
// shortcut.
func (u *UrlGenerator) TemporarySignedRoute(name string, expiration int64, parameters map[string]any, absolute bool) (string, error) {
	return u.SignedRoute(name, parameters, expiration, absolute)
}

// HasValidSignature reports whether request carries a signature that matches
// its URL and has not expired.
func (u *UrlGenerator) HasValidSignature(request URLRequest, absolute bool) bool {
	return u.HasCorrectSignature(request, absolute) && u.SignatureHasNotExpired(request)
}

// HasCorrectSignature performs the HMAC comparison without the expiry check.
func (u *UrlGenerator) HasCorrectSignature(request URLRequest, absolute bool) bool {
	urlStr := request.URL()

	if !absolute {
		urlStr = "/" + strings.TrimLeft(request.Path(), "/")
	}

	if idx := strings.Index(urlStr, "?"); idx >= 0 {
		urlStr = urlStr[:idx]
	}

	queryString := stripSignatureFromQuery(request.QueryString())
	original := urlStr

	if queryString != "" {
		original += "?" + queryString
	}

	expected := request.Query("signature")
	mac := hmac.New(sha256.New, []byte(u.key))
	mac.Write([]byte(original))
	got := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(got), []byte(expected))
}

// SignatureHasNotExpired reports whether the request's "expires" query
// parameter has not yet passed.
func (u *UrlGenerator) SignatureHasNotExpired(request URLRequest) bool {
	expires := request.Query("expires")

	if expires == "" {
		return true
	}

	exp, err := strconv.ParseInt(expires, 10, 64)

	if err != nil {
		return false
	}

	return time.Now().Unix() <= exp
}

func stripSignatureFromQuery(qs string) string {
	if qs == "" {
		return ""
	}

	parts := strings.Split(qs, "&")
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		if strings.HasPrefix(p, "signature=") {
			continue
		}

		out = append(out, p)
	}

	return strings.Join(out, "&")
}

func isAbsoluteURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "//") || strings.HasPrefix(path, "mailto:") ||
		strings.HasPrefix(path, "tel:")
}
