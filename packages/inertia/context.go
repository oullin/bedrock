package inertia

import (
	"context"

	"github.com/bedrock/packages/inertia/protocol"
)

var (
	ctxKeyProps            = &contextKey{"props"}
	ctxKeyTemplateData     = &contextKey{"templateData"}
	ctxKeyValidationErrors = &contextKey{"validationErrors"}
	ctxKeyEncryptHistory   = &contextKey{"encryptHistory"}
	ctxKeyClearHistory     = &contextKey{"clearHistory"}
	ctxKeyHead             = &contextKey{"head"}
)

type contextKey struct{ name string }

// SetProp stores a single prop on the request context. Props set this
// way are merged into the response during Render, with higher priority
// than shared props but lower than props passed directly to Render.
func SetProp(ctx context.Context, key string, val any) context.Context {
	p := propsFromContext(ctx)
	p[key] = val

	return context.WithValue(ctx, ctxKeyProps, p)
}

// SetProps stores multiple props on the request context.
func SetProps(ctx context.Context, props protocol.Props) context.Context {
	p := propsFromContext(ctx)

	for k, v := range props {
		p[k] = v
	}

	return context.WithValue(ctx, ctxKeyProps, p)
}

// SetValidationErrors stores validation errors in the request context.
// They are automatically added to the response props under the "errors" key.
func SetValidationErrors(ctx context.Context, errors protocol.ValidationErrors) context.Context {
	return context.WithValue(ctx, ctxKeyValidationErrors, errors)
}

// SetEncryptHistory flags the response to encrypt the browser history state.
func SetEncryptHistory(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyEncryptHistory, true)
}

// SetClearHistory flags the response to clear any encrypted browser history.
func SetClearHistory(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyClearHistory, true)
}

// SetTemplateData stores additional data for the root HTML template
// used during initial (non-XHR) page visits.
func SetTemplateData(ctx context.Context, data protocol.TemplateData) context.Context {
	existing := templateDataFromContext(ctx)

	for k, v := range data {
		existing[k] = v
	}

	return context.WithValue(ctx, ctxKeyTemplateData, existing)
}

// SetTemplateDatum stores a single template data value.
func SetTemplateDatum(ctx context.Context, key string, val any) context.Context {
	d := templateDataFromContext(ctx)
	d[key] = val

	return context.WithValue(ctx, ctxKeyTemplateData, d)
}

// PropsFromContext returns the props stored in the request context.
func PropsFromContext(ctx context.Context) protocol.Props {
	if p, ok := ctx.Value(ctxKeyProps).(protocol.Props); ok {
		return p
	}

	return make(protocol.Props)
}

func propsFromContext(ctx context.Context) protocol.Props {
	return PropsFromContext(ctx)
}

func validationErrorsFromContext(ctx context.Context) protocol.ValidationErrors {
	if v, ok := ctx.Value(ctxKeyValidationErrors).(protocol.ValidationErrors); ok {
		return v
	}

	return nil
}

func encryptHistoryFromContext(ctx context.Context) bool {
	v, _ := ctx.Value(ctxKeyEncryptHistory).(bool)

	return v
}

func clearHistoryFromContext(ctx context.Context) bool {
	v, _ := ctx.Value(ctxKeyClearHistory).(bool)

	return v
}

func templateDataFromContext(ctx context.Context) protocol.TemplateData {
	if d, ok := ctx.Value(ctxKeyTemplateData).(protocol.TemplateData); ok {
		return d
	}

	return make(protocol.TemplateData)
}

// SetHead stores head elements on the request context. These are rendered
// into {{ .inertiaHead }} on initial page loads. Per-request head elements
// are merged with (and override) default head elements set via WithHead.
func SetHead(ctx context.Context, head protocol.Head) context.Context {
	existing := headFromContext(ctx)
	merged := protocol.MergeHead(existing, head)

	return context.WithValue(ctx, ctxKeyHead, merged)
}

// SetTitle is a convenience helper that sets only the <title> element
// on the request context.
func SetTitle(ctx context.Context, title string) context.Context {
	return SetHead(ctx, protocol.Head{Title: title})
}

// SetLang is a convenience helper that sets only the lang attribute
// on the request context.
func SetLang(ctx context.Context, lang string) context.Context {
	return SetHead(ctx, protocol.Head{Lang: lang})
}

// SetMeta is a convenience helper that adds meta tags to the request context.
func SetMeta(ctx context.Context, tags ...protocol.MetaTag) context.Context {
	return SetHead(ctx, protocol.Head{Meta: tags})
}

// SetLinks is a convenience helper that adds link tags to the request context.
func SetLinks(ctx context.Context, links ...protocol.LinkTag) context.Context {
	return SetHead(ctx, protocol.Head{Links: links})
}

// SetCSRFToken stores a CSRF token in the request context. When present,
// Render automatically adds <meta name="csrf-token" content="TOKEN"> to
// the head on initial page loads. Delegates to protocol.SetCSRFToken.
func SetCSRFToken(ctx context.Context, token string) context.Context {
	return protocol.SetCSRFToken(ctx, token)
}

// SetPrecognition marks the request context as a precognition request.
// Delegates to protocol.SetPrecognition.
func SetPrecognition(ctx context.Context) context.Context {
	return protocol.SetPrecognition(ctx)
}

// SetLocale stores the resolved locale in the request context. This is
// typically called by the i18n middleware, not by application code.
// Delegates to protocol.SetLocale.
func SetLocale(ctx context.Context, locale *protocol.Locale) context.Context {
	return protocol.SetLocale(ctx, locale)
}

func headFromContext(ctx context.Context) protocol.Head {
	if h, ok := ctx.Value(ctxKeyHead).(protocol.Head); ok {
		return h
	}

	return protocol.Head{}
}
