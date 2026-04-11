package middleware

import "net/http"

// FrameGuardMode defines the X-Frame-Options header value.

// FrameGuard sets the X-Frame-Options header to protect against clickjacking.
type FrameGuard struct {
	mode string
}

const (
	FrameDeny      = "DENY"
	FrameSameOrgin = "SAMEORIGIN"
)

// NewFrameGuard creates a FrameGuard middleware. mode should be "DENY" or
// "SAMEORIGIN"; it defaults to "SAMEORIGIN".
func NewFrameGuard(mode ...string) *FrameGuard {
	m := FrameSameOrgin

	if len(mode) > 0 && mode[0] != "" {
		m = mode[0]
	}

	return &FrameGuard{mode: m}
}

// Wrap returns an http.Handler that sets X-Frame-Options before delegating.
func (fg *FrameGuard) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", fg.mode)
		next.ServeHTTP(w, r)
	})
}
