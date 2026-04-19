package protocol_test

import (
	"context"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia/protocol"
)

func TestCSRFToken_ContextRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := protocol.SetCSRFToken(context.Background(), "test-token-abc")
	got := protocol.CSRFTokenFromContext(ctx)

	if got != "test-token-abc" {
		t.Errorf("CSRFTokenFromContext() = %q, want %q", got, "test-token-abc")
	}
}

func TestCSRFTokenFromContext_Missing(t *testing.T) {
	t.Parallel()

	got := protocol.CSRFTokenFromContext(context.Background())

	if strings.TrimSpace(got) != "" {
		t.Errorf("CSRFTokenFromContext() = %q, want empty string", got)
	}
}

func TestLocale_ContextRoundTrip(t *testing.T) {
	t.Parallel()

	locale := &protocol.Locale{
		Code:      "es",
		Name:      "Español",
		Direction: "ltr",
	}

	ctx := protocol.SetLocale(context.Background(), locale)
	got := protocol.LocaleFromContext(ctx)

	if got == nil {
		t.Fatal("LocaleFromContext() = nil, want locale")
	}

	if got.Code != "es" {
		t.Errorf("Code = %q, want %q", got.Code, "es")
	}

	if got.Name != "Español" {
		t.Errorf("Name = %q, want %q", got.Name, "Español")
	}
}

func TestLocaleFromContext_Missing(t *testing.T) {
	t.Parallel()

	got := protocol.LocaleFromContext(context.Background())

	if got != nil {
		t.Errorf("LocaleFromContext() = %v, want nil", got)
	}
}

func TestHTTPPreview_ContextRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := protocol.SetHTTPPreview(context.Background())

	if !protocol.IsHTTPPreview(ctx) {
		t.Error("IsHTTPPreview() = false, want true")
	}
}

func TestIsHTTPPreview_Missing(t *testing.T) {
	t.Parallel()

	if protocol.IsHTTPPreview(context.Background()) {
		t.Error("IsHTTPPreview() = true, want false")
	}
}
