package protocol

import (
	"context"

	"github.com/bedrock/packages/seo"
)

type ctxKey struct{ name string }

var (
	ctxKeyCSRFToken    = &ctxKey{"csrfToken"}
	ctxKeyPrecognition = &ctxKey{"precognition"}
)

// SetCSRFToken stores a CSRF token in the request context. When present,
// Render automatically adds <meta name="csrf-token" content="TOKEN"> to
// the head on initial page loads.
func SetCSRFToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, ctxKeyCSRFToken, token)
}

// CSRFTokenFromContext returns the CSRF token stored in context, or "".
func CSRFTokenFromContext(ctx context.Context) string {
	s, _ := ctx.Value(ctxKeyCSRFToken).(string)

	return s
}

// SetLocale re-exports seo.SetLocale.
var SetLocale = seo.SetLocale

// LocaleFromContext re-exports seo.LocaleFromContext.
var LocaleFromContext = seo.LocaleFromContext

// SetPrecognition marks the request context as a precognition request.
func SetPrecognition(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyPrecognition, true)
}

// IsPrecognition reports whether the request context was marked as
// a precognition request by the precognition middleware.
func IsPrecognition(ctx context.Context) bool {
	v, _ := ctx.Value(ctxKeyPrecognition).(bool)

	return v
}
