package protocol

import "context"

var (
	ctxKeyCSRFToken    = &ctxKey{"csrfToken"}
	ctxKeyLocale       = &ctxKey{"locale"}
	ctxKeyHTTPPreview = &ctxKey{"httppreview"}
)

type ctxKey struct{ name string }

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

// SetLocale stores the resolved locale in the request context.
func SetLocale(ctx context.Context, locale *Locale) context.Context {
	return context.WithValue(ctx, ctxKeyLocale, locale)
}

// LocaleFromContext returns the locale stored in context, or nil.
func LocaleFromContext(ctx context.Context) *Locale {
	l, _ := ctx.Value(ctxKeyLocale).(*Locale)

	return l
}

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
