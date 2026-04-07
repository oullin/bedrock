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

// Capture mirrors Laravel's request capture entrypoint.
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
