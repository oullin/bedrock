package precognition

import (
	"context"
	"net/http"
)

// contextKey is an unexported type used for context values to avoid collisions
// with other packages.
type contextKey struct{}

// precognitiveKey is the context key that marks a request as precognitive.
// This is the Go equivalent of the upstream // $request->attributes->set('precognitive', true).
var precognitiveKey = contextKey{}

// MarkPrecognitive returns a shallow copy of r with the precognitive context
// value set. This is called by the middleware after it determines that the
// request is attempting precognition.
func MarkPrecognitive(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), precognitiveKey, true)

	return r.WithContext(ctx)
}

// IsPrecognitive reports whether the middleware has marked this request as
// precognitive via [MarkPrecognitive]. This is the Go equivalent of the upstream // $request->isPrecognitive() which checks the request attribute set by the
// middleware — distinct from the header-based check.
func IsPrecognitive(r *http.Request) bool {
	v, _ := r.Context().Value(precognitiveKey).(bool)

	return v
}

// IsAttemptingPrecognition reports whether the request carries a
// Precognition header with the value "true". This checks the client's intent
// to make a precognitive request, matching the upstream // $request->isAttemptingPrecognition() which checks the header exactly.
func IsAttemptingPrecognition(r *http.Request) bool {
	return r.Header.Get("Precognition") == "true"
}
