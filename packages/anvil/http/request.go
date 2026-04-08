package http

import (
	"encoding/json"
	nethttp "net/http"
	"strconv"
	"strings"
)

// Request aliases the standard library HTTP request type.
type Request = nethttp.Request

// ResponseWriter aliases the standard library response writer type.
type ResponseWriter = nethttp.ResponseWriter

// Middleware decorates an HTTP handler.
type Middleware func(nethttp.Handler) nethttp.Handler

// Capture mirrors Upstream's request capture entrypoint.
func Capture(request *nethttp.Request) *Request {
	return request
}

// Input returns a form/query value by key with an optional default.
func Input(r *Request, key string, defaultValue ...string) string {
	if r.Form != nil {
		if val := r.FormValue(key); val != "" {
			return val
		}
	}

	if val := r.URL.Query().Get(key); val != "" {
		return val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}

// Query returns a query string value by key.
func Query(r *Request, key string, defaultValue ...string) string {
	val := r.URL.Query().Get(key)
	if val != "" {
		return val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}

// Boolean returns a query/form value interpreted as a boolean.
func Boolean(r *Request, key string) bool {
	val := strings.TrimSpace(Input(r, key))
	if val == "" {
		return false
	}

	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}

	return parsed
}

// Integer returns a query/form value interpreted as an integer.
func Integer(r *Request, key string, defaultValue ...int) int {
	val := strings.TrimSpace(Input(r, key))
	if val == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0
	}

	parsed, err := strconv.Atoi(val)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0
	}

	return parsed
}

// JSON writes a JSON response with the given status code.
func JSON(w ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

// Redirect sends an HTTP redirect response.
func Redirect(w ResponseWriter, r *Request, url string, status int) {
	nethttp.Redirect(w, r, url, status)
}

// ---------- Path/URL methods ----------

// Path returns the request path without the query string.
func Path(r *Request) string {
	return r.URL.Path
}

// Url returns the URL without the query string.
func Url(r *Request) string {
	return Scheme(r) + "://" + r.Host + r.URL.Path
}

// FullUrl returns the full URL including the query string.
func FullUrl(r *Request) string {
	u := Scheme(r) + "://" + r.Host + r.URL.Path
	if r.URL.RawQuery != "" {
		u += "?" + r.URL.RawQuery
	}
	return u
}

// FullUrlWithQuery returns the full URL with additional query parameters merged in.
func FullUrlWithQuery(r *Request, merge map[string]string) string {
	q := r.URL.Query()
	for k, v := range merge {
		q.Set(k, v)
	}
	return Scheme(r) + "://" + r.Host + r.URL.Path + "?" + q.Encode()
}

// Host returns the host from the request (may include port).
func Host(r *Request) string {
	return r.Host
}

// Scheme returns the request scheme ("https" or "http").
func Scheme(r *Request) string {
	if r.TLS != nil {
		return "https"
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		return strings.ToLower(proto)
	}
	return "http"
}

// Method returns the HTTP method.
func Method(r *Request) string {
	return r.Method
}

// IsMethod checks if the request method matches (case-insensitive).
func IsMethod(r *Request, method string) bool {
	return strings.EqualFold(r.Method, method)
}

// ---------- Content negotiation ----------

// ContentType returns the request Content-Type header value.
func ContentType(r *Request) string {
	ct := r.Header.Get("Content-Type")
	if i := strings.IndexByte(ct, ';'); i >= 0 {
		ct = ct[:i]
	}
	return strings.TrimSpace(ct)
}

// IsJson reports whether the request Content-Type indicates JSON.
func IsJson(r *Request) bool {
	return strings.Contains(ContentType(r), "json")
}

// WantsJson reports whether the client prefers a JSON response.
func WantsJson(r *Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "/json") || strings.Contains(accept, "+json")
}

// Accepts reports whether the client accepts any of the given content types.
func Accepts(r *Request, types []string) bool {
	accept := r.Header.Get("Accept")
	if accept == "" || accept == "*/*" {
		return true
	}
	for _, t := range types {
		if strings.Contains(accept, t) {
			return true
		}
	}
	return false
}

// Prefers returns the first content type from the list that the client accepts.
// Returns an empty string if none match.
func Prefers(r *Request, types []string) string {
	accept := r.Header.Get("Accept")
	if accept == "" || accept == "*/*" {
		if len(types) > 0 {
			return types[0]
		}
		return ""
	}
	for _, t := range types {
		if strings.Contains(accept, t) {
			return t
		}
	}
	return ""
}

// ---------- Request data methods ----------

// All returns all query and form values merged into a single map.
func All(r *Request) map[string]string {
	_ = r.ParseForm()
	result := map[string]string{}
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	for k, v := range r.Form {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}
	return result
}

// Has reports whether all the given keys are present in the request input.
func Has(r *Request, keys ...string) bool {
	all := All(r)
	for _, k := range keys {
		if _, ok := all[k]; !ok {
			return false
		}
	}
	return true
}

// Missing reports whether any of the given keys are absent from the request input.
func Missing(r *Request, keys ...string) bool {
	return !Has(r, keys...)
}

// Only returns a map containing only the specified keys from the request input.
func Only(r *Request, keys ...string) map[string]string {
	all := All(r)
	result := map[string]string{}
	for _, k := range keys {
		if v, ok := all[k]; ok {
			result[k] = v
		}
	}
	return result
}

// Except returns a map with all request input except the specified keys.
func Except(r *Request, keys ...string) map[string]string {
	all := All(r)
	for _, k := range keys {
		delete(all, k)
	}
	return all
}

// Header returns a request header value with an optional default.
func Header(r *Request, key string, defaultValue ...string) string {
	val := r.Header.Get(key)
	if val != "" {
		return val
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

// BearerToken extracts the bearer token from the Authorization header.
func BearerToken(r *Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return auth[7:]
	}
	return ""
}

// Ip returns the client IP address.
func Ip(r *Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		if i := strings.IndexByte(forwarded, ','); i >= 0 {
			return strings.TrimSpace(forwarded[:i])
		}
		return strings.TrimSpace(forwarded)
	}
	host := r.RemoteAddr
	if i := strings.LastIndexByte(host, ':'); i >= 0 {
		return host[:i]
	}
	return host
}
