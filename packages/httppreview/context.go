package httppreview

import (
	"context"
	"net/http"
)

// contextKey is an unexported type used for context values to avoid collisions
// with other packages.
type contextKey struct{}

// precognitiveKey is the context key that marks a request as precognitive.
var precognitiveKey = contextKey{}

// MarkPrecognitive returns a shallow copy of r with the precognitive context
// value set. This is called by the middleware after it determines that the
// request is attempting httppreview.
func MarkPrecognitive(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), precognitiveKey, true)

	return r.WithContext(ctx)
}

// IsPrecognitive reports whether the middleware has marked this request as
// precognitive via [MarkPrecognitive].
// middleware — distinct from the header-based check.
func IsPrecognitive(r *http.Request) bool {
	v, _ := r.Context().Value(precognitiveKey).(bool)

	return v
}

// IsAttemptingHTTPPreview reports whether the request carries a
// HTTPPreview header with the value "true". This checks the client's intent
// to make a precognitive request and matches the header exactly.
func IsAttemptingHTTPPreview(r *http.Request) bool {
	return r.Header.Get("HTTPPreview") == "true"
}
