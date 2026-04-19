package seo

import "context"

type ctxKey struct{ name string }

var ctxKeyLocale = &ctxKey{"locale"}

// SetLocale stores the resolved locale in the request context.
func SetLocale(ctx context.Context, locale *Locale) context.Context {
	return context.WithValue(ctx, ctxKeyLocale, locale)
}

// LocaleFromContext returns the locale stored in context, or nil.
func LocaleFromContext(ctx context.Context) *Locale {
	l, _ := ctx.Value(ctxKeyLocale).(*Locale)

	return l
}
