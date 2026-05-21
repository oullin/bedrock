package watchers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bedrock/packages/telescope"
)

// 64 KB

// RequestWatcher monitors HTTP requests/responses and records them as
// Telescope entries. It mirrors the the underlying behavior class.
//
// Options:
//   - "size_limit" (int): max response body bytes to store (default 64 KB).
//   - "ignore_http_methods" ([]string): HTTP methods to skip (e.g., "OPTIONS").
//   - "ignore_status_codes" ([]int): HTTP status codes to skip.
type RequestWatcher struct {
	telescope.BaseWatcher
}

// NewRequestWatcher creates a RequestWatcher with the given options.

// Register is a no-op for RequestWatcher; use Middleware() to integrate with
// an http.Handler chain.

// Middleware wraps an http.Handler, recording each request/response pair.

// record builds and stores a request entry.

// Headers (masked).

// Payload (masked, body read up to limit).

// Response.

// ShouldIgnoreMethod reports whether the HTTP method is in the ignore list.

// shouldIgnoreStatus reports whether the HTTP status code is in the ignore list.

// sizeLimit returns the configured response size limit in bytes.

// ─── Helpers ─────────────────────────────────────────────────────────────────

// Strip port.

// ─── Response recorder ───────────────────────────────────────────────────────

type responseRecorder struct {
	http.ResponseWriter
	statusCode  int
	body        bytes.Buffer
	headers     http.Header
	contentType string
}

const defaultRequestSizeLimit = 64 * 1024

func NewRequestWatcher(t *telescope.Telescope, options map[string]any) *RequestWatcher {
	w := &RequestWatcher{}
	w.SetTelescope(t)
	w.Options = options

	return w
}

func (w *RequestWatcher) Register(_ any) error { return nil }

func (w *RequestWatcher) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := newResponseRecorder(rw)
		next.ServeHTTP(rec, r)

		duration := time.Since(start)

		if !w.shouldIgnoreMethod(r.Method) && !w.shouldIgnoreStatus(rec.statusCode) {
			w.record(r, rec, duration)
		}
	})
}

func (w *RequestWatcher) record(r *http.Request, rec *responseRecorder, duration time.Duration) {
	sizeLimit := w.sizeLimit()
	scope := w.Scope()

	requestHeaders := maskHeaders(r.Header, scope.HiddenRequestHeaders())

	payload := readRequestBody(r, sizeLimit)
	payload = maskMap(payload, scope.HiddenRequestParameters())

	responseBody := truncate(rec.body.Bytes(), sizeLimit)
	responseHeaders := maskHeaders(rec.headers, scope.HiddenResponseParameters())

	content := map[string]any{
		"method":           r.Method,
		"uri":              requestURI(r),
		"headers":          requestHeaders,
		"payload":          payload,
		"response_status":  rec.statusCode,
		"response_headers": responseHeaders,
		"response":         formatResponseBody(rec.contentType, responseBody),
		"duration":         float64(duration.Milliseconds()),
		"ip_address":       clientIP(r),
	}

	entry := telescope.NewEntry(telescope.EntryTypeRequest, content)

	scope.RecordRequest(entry)
}

func (w *RequestWatcher) shouldIgnoreMethod(method string) bool {
	for _, m := range w.StringsOption("ignore_http_methods") {
		if strings.EqualFold(m, method) {
			return true
		}
	}

	return false
}

func (w *RequestWatcher) shouldIgnoreStatus(code int) bool {
	if raw, ok := w.Options["ignore_status_codes"]; ok {
		switch v := raw.(type) {
		case []int:
			for _, s := range v {
				if s == code {
					return true
				}
			}
		case []any:
			for _, s := range v {
				if n, ok := s.(int); ok && n == code {
					return true
				}
			}
		}
	}

	return false
}

func (w *RequestWatcher) sizeLimit() int {
	if v, ok := w.Options["size_limit"]; ok {
		switch n := v.(type) {
		case int:
			return n
		case int64:
			return int(n)
		case float64:
			return int(n)
		}
	}

	return defaultRequestSizeLimit
}

func requestURI(r *http.Request) string {
	uri := r.RequestURI

	if uri == "" {
		uri = r.URL.RequestURI()
	}

	return uri
}

func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.SplitN(ip, ",", 2)[0]
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	addr := r.RemoteAddr

	if idx := strings.LastIndex(addr, ":"); idx > 0 {
		addr = addr[:idx]
	}

	return addr
}

func readRequestBody(r *http.Request, limit int) map[string]any {
	if r.Body == nil {
		return parseQuery(r.URL.RawQuery)
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, int64(limit)))
	r.Body = io.NopCloser(bytes.NewReader(body))

	if err != nil || len(body) == 0 {
		return parseQuery(r.URL.RawQuery)
	}

	ct := r.Header.Get("Content-Type")

	if strings.Contains(ct, "application/json") {
		var m map[string]any

		if json.Unmarshal(body, &m) == nil {
			return m
		}
	}

	if strings.Contains(ct, "multipart/form-data") {
		if payload, ok := parseMultipartBody(body, ct); ok {
			return payload
		}
	}

	if strings.Contains(ct, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err == nil {
			m := make(map[string]any, len(r.PostForm))

			for k, v := range r.PostForm {
				if len(v) == 1 {
					m[k] = v[0]
				} else {
					m[k] = v
				}
			}

			return m
		}
	}

	return map[string]any{"raw": string(body)}
}

func parseQuery(rawQuery string) map[string]any {
	if rawQuery == "" {
		return nil
	}

	result := make(map[string]any)

	for _, part := range strings.Split(rawQuery, "&") {
		kv := strings.SplitN(part, "=", 2)

		if len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}

	return result
}

func maskHeaders(headers http.Header, hidden []string) map[string]string {
	set := make(map[string]struct{}, len(hidden))

	for _, h := range hidden {
		set[strings.ToLower(h)] = struct{}{}
	}

	out := make(map[string]string, len(headers))

	for k, v := range headers {
		if _, ok := set[strings.ToLower(k)]; ok {
			out[k] = "********"
		} else if len(v) > 0 {
			out[k] = v[0]
		}
	}

	return out
}

func maskMap(m map[string]any, hidden []string) map[string]any {
	set := make(map[string]struct{}, len(hidden))

	for _, h := range hidden {
		set[h] = struct{}{}
	}

	out := make(map[string]any, len(m))

	for k, v := range m {
		if _, ok := set[k]; ok {
			out[k] = "********"
		} else {
			out[k] = v
		}
	}

	return out
}

func formatResponseBody(contentType string, body []byte) any {
	if strings.Contains(contentType, "application/json") {
		var v any

		if json.Unmarshal(body, &v) == nil {
			return v
		}
	}

	return string(body)
}

func truncate(b []byte, limit int) []byte {
	if len(b) <= limit {
		return b
	}

	return b[:limit]
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		headers:        make(http.Header),
	}
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)

	if r.contentType == "" {
		r.contentType = r.ResponseWriter.Header().Get("Content-Type")
	}

	for k, v := range r.ResponseWriter.Header() {
		r.headers[k] = v
	}

	return r.ResponseWriter.Write(b)
}
