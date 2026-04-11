package middleware

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

// CheckResponseForModifications handles conditional requests using ETag and
// If-Modified-Since headers. When the response has not been modified, a 304
// Not Modified is returned instead of the full body.
type CheckResponseForModifications struct{}

// NewCheckResponseForModifications creates the middleware.
func NewCheckResponseForModifications() *CheckResponseForModifications {
	return &CheckResponseForModifications{}
}

// Wrap returns an http.Handler that checks ETags and modification times.
func (m *CheckResponseForModifications) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only apply to safe methods.
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			next.ServeHTTP(w, r)

			return
		}

		// Capture the response.
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		result := rec.Result()
		body := rec.Body.Bytes()

		// Generate ETag if not present.
		etag := result.Header.Get("ETag")

		if etag == "" && len(body) > 0 {
			hash := sha256.Sum256(body)
			etag = fmt.Sprintf(`"%x"`, hash[:8])
			rec.Header().Set("ETag", etag)
		}

		// Check If-None-Match.
		if ifNoneMatch := r.Header.Get("If-None-Match"); ifNoneMatch != "" && etag != "" {
			if ifNoneMatch == etag {
				copyHeaders(w, rec)
				w.WriteHeader(http.StatusNotModified)

				return
			}
		}

		// Check If-Modified-Since.
		if ims := r.Header.Get("If-Modified-Since"); ims != "" {
			if lastMod := result.Header.Get("Last-Modified"); lastMod != "" {
				imsTime, err1 := time.Parse(http.TimeFormat, ims)
				lmTime, err2 := time.Parse(http.TimeFormat, lastMod)

				if err1 == nil && err2 == nil && !lmTime.After(imsTime) {
					copyHeaders(w, rec)
					w.WriteHeader(http.StatusNotModified)

					return
				}
			}
		}

		// Forward the full response.
		copyHeaders(w, rec)
		w.WriteHeader(rec.Code)
		w.Write(body)
	})
}

func copyHeaders(w http.ResponseWriter, rec *httptest.ResponseRecorder) {
	for k, vals := range rec.Header() {
		for _, v := range vals {
			w.Header().Add(k, v)
		}
	}
}
