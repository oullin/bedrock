package protocol

import (
	"context"

	"github.com/bedrock/packages/seo"
)

type ctxKey struct{ name string }

var (
	ctxKeyCSRFToken    = &ctxKey{"csrfToken"}
	ctxKeyHTTPPreview = &ctxKey{"httppreview"}
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

// SetHTTPPreview marks the request context as a httppreview request.
func SetHTTPPreview(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyHTTPPreview, true)
}

// IsHTTPPreview reports whether the request context was marked as
// a httppreview request by the httppreview middleware.
func IsHTTPPreview(ctx context.Context) bool {
	v, _ := ctx.Value(ctxKeyHTTPPreview).(bool)

	return v
}
