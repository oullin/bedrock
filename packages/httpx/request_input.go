package httpx

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Header returns a request header value, or the optional fallback.
func (r *Request) Header(key string, fallback ...string) string {
	v := r.raw.Header.Get(key)

	if v == "" && len(fallback) > 0 {
		return fallback[0]
	}

	return v
}

// HasHeader checks whether the named header is present.
func (r *Request) HasHeader(key string) bool {
	return r.raw.Header.Get(key) != ""
}

// BearerToken extracts the token from an "Authorization: Bearer <token>" header.
func (r *Request) BearerToken() string {
	auth := r.raw.Header.Get("Authorization")

	if len(auth) > 7 && strings.EqualFold(auth[:7], "bearer ") {
		return auth[7:]
	}

	return ""
}

// Cookie returns the value of the named cookie, or the optional fallback.
func (r *Request) Cookie(name string, fallback ...string) string {
	c, err := r.raw.Cookie(name)

	if err != nil || c.Value == "" {
		if len(fallback) > 0 {
			return fallback[0]
		}

		return ""
	}

	return c.Value
}

// HasCookie checks whether the named cookie is present.
func (r *Request) HasCookie(name string) bool {
	_, err := r.raw.Cookie(name)

	return err == nil
}

// All returns all input data merged from query parameters and the request body
// (form or JSON). Query parameters take lower precedence; body values override
// query values with the same key.
func (r *Request) All() map[string]any {
	if r.parsedInput != nil {
		return r.parsedInput
	}

	merged := make(map[string]any)

	// Query parameters first (lower precedence).
	for k, v := range r.raw.URL.Query() {
		if len(v) == 1 {
			merged[k] = v[0]
		} else {
			merged[k] = v
		}
	}

	// Body input (higher precedence).
	ct := r.raw.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(ct, "application/json"):
		body, err := r.readBody()

		if err == nil && len(body) > 0 {
			var data map[string]any

			if json.Unmarshal(body, &data) == nil {
				for k, v := range data {
					merged[k] = v
				}
			}
		}
	default:
		// Form-encoded or multipart.
		_ = r.raw.ParseForm()

		if r.raw.PostForm != nil {
			for k, v := range r.raw.PostForm {
				if len(v) == 1 {
					merged[k] = v[0]
				} else {
					merged[k] = v
				}
			}
		}
	}

	r.parsedInput = merged

	return merged
}

// Input returns an input value by key (body takes precedence over query). If
// the key is absent the optional fallback is returned.
func (r *Request) Input(key string, fallback ...string) string {
	all := r.All()

	if v, ok := all[key]; ok {
		return toString(v)
	}

	if len(fallback) > 0 {
		return fallback[0]
	}

	return ""
}

// Query returns a query-string parameter value.
func (r *Request) Query(key string, fallback ...string) string {
	v := r.raw.URL.Query().Get(key)

	if v == "" && len(fallback) > 0 {
		return fallback[0]
	}

	return v
}

// Post returns a form-posted value.
func (r *Request) Post(key string, fallback ...string) string {
	_ = r.raw.ParseForm()

	v := r.raw.PostFormValue(key)

	if v == "" && len(fallback) > 0 {
		return fallback[0]
	}

	return v
}

// Boolean returns the input value as a boolean. Values "1", "true", "on" and
// "yes" are truthy; everything else is false.
func (r *Request) Boolean(key string) bool {
	v := strings.ToLower(r.Input(key))

	switch v {
	case "1", "true", "on", "yes":
		return true
	}

	return false
}

// Integer returns the input value as an int. Returns the optional fallback (or
// 0) when the value is absent or not numeric.
func (r *Request) Integer(key string, fallback ...int) int {
	v := r.Input(key)

	if v == "" {
		if len(fallback) > 0 {
			return fallback[0]
		}

		return 0
	}

	n, err := strconv.Atoi(v)

	if err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		}

		return 0
	}

	return n
}

// Float returns the input value as a float64.
func (r *Request) Float(key string, fallback ...float64) float64 {
	v := r.Input(key)

	if v == "" {
		if len(fallback) > 0 {
			return fallback[0]
		}

		return 0
	}

	f, err := strconv.ParseFloat(v, 64)

	if err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		}

		return 0
	}

	return f
}

// Date parses an input value as a time.Time using the provided layout.
func (r *Request) Date(key string, layout string) (time.Time, error) {
	return time.Parse(layout, r.Input(key))
}

// Only returns a subset of the input containing only the specified keys.
func (r *Request) Only(keys ...string) map[string]any {
	all := r.All()
	result := make(map[string]any, len(keys))

	for _, k := range keys {
		if v, ok := all[k]; ok {
			result[k] = v
		}
	}

	return result
}

// Except returns all input except the specified keys.
func (r *Request) Except(keys ...string) map[string]any {
	all := r.All()
	result := make(map[string]any, len(all))

	skip := make(map[string]struct{}, len(keys))

	for _, k := range keys {
		skip[k] = struct{}{}
	}

	for k, v := range all {
		if _, ok := skip[k]; !ok {
			result[k] = v
		}
	}

	return result
}

// Has returns true when all of the given keys are present in the input.
func (r *Request) Has(keys ...string) bool {
	all := r.All()

	for _, k := range keys {
		if _, ok := all[k]; !ok {
			return false
		}
	}

	return true
}

// HasAny returns true when at least one of the given keys is present.
func (r *Request) HasAny(keys ...string) bool {
	all := r.All()

	for _, k := range keys {
		if _, ok := all[k]; ok {
			return true
		}
	}

	return false
}

// Filled returns true when all of the given keys are present and non-empty.
func (r *Request) Filled(keys ...string) bool {
	all := r.All()

	for _, k := range keys {
		v, ok := all[k]

		if !ok {
			return false
		}

		if toString(v) == "" {
			return false
		}
	}

	return true
}

// Missing returns true when any of the given keys is absent from the input.
func (r *Request) Missing(keys ...string) bool {
	return !r.Has(keys...)
}

// Keys returns all input keys.
func (r *Request) Keys() []string {
	all := r.All()
	keys := make([]string, 0, len(all))

	for k := range all {
		keys = append(keys, k)
	}

	return keys
}

// File returns the uploaded file for the given input name, or nil.
func (r *Request) File(key string) *UploadedFile {
	if r.raw.MultipartForm == nil {
		_ = r.raw.ParseMultipartForm(32 << 20) // 32 MB
	}

	if r.raw.MultipartForm == nil {
		return nil
	}

	files, ok := r.raw.MultipartForm.File[key]

	if !ok || len(files) == 0 {
		return nil
	}

	return NewUploadedFile(files[0])
}

// AllFiles returns all uploaded files keyed by input name.
func (r *Request) AllFiles() map[string][]*UploadedFile {
	if r.raw.MultipartForm == nil {
		_ = r.raw.ParseMultipartForm(32 << 20)
	}

	if r.raw.MultipartForm == nil {
		return nil
	}

	result := make(map[string][]*UploadedFile, len(r.raw.MultipartForm.File))

	for name, headers := range r.raw.MultipartForm.File {
		files := make([]*UploadedFile, len(headers))

		for i, h := range headers {
			files[i] = NewUploadedFile(h)
		}

		result[name] = files
	}

	return result
}

// HasFile checks whether a file was uploaded with the given input name.
func (r *Request) HasFile(key string) bool {
	return r.File(key) != nil
}

// readBody reads and returns the request body. The body is replaced with a
// no-op reader so it can be consumed only once.
func (r *Request) readBody() ([]byte, error) {
	if r.raw.Body == nil {
		return nil, nil
	}

	body, err := io.ReadAll(r.raw.Body)
	r.raw.Body.Close()
	r.raw.Body = http.NoBody

	return body, err
}

// toString converts a value to its string representation.
func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case []string:
		if len(val) > 0 {
			return val[0]
		}

		return ""
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	case json.Number:
		return val.String()
	case nil:
		return ""
	default:
		return ""
	}
}

// multipartFileHeader is extracted for testability. It matches
// *multipart.FileHeader from the stdlib.
var _ *multipart.FileHeader
