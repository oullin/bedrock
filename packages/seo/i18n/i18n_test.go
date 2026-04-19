package i18n_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/seo"
	"github.com/bedrock/packages/seo/i18n"
)

func newTestConfig() *i18n.I18nConfig {
	return &i18n.I18nConfig{
		DefaultLocale: "en",
		URLPrefix:     true,
		Locales: map[string]*seo.Locale{
			"en": {
				Code:      "en",
				Name:      "English",
				Direction: "ltr",
				Head: seo.Head{
					Title: "My App",
					Meta: []seo.MetaTag{
						{Name: "description", Content: "English description"},
					},
				},
			},
			"es": {
				Code:      "es",
				Name:      "Español",
				Direction: "ltr",
				Head: seo.Head{
					Title: "Mi App",
					Meta: []seo.MetaTag{
						{Name: "description", Content: "Descripción en español"},
					},
				},
			},
			"ar": {
				Code:      "ar",
				Name:      "العربية",
				Direction: "rtl",
				Head: seo.Head{
					Title: "تطبيقي",
				},
			},
		},
	}
}

func TestDefault(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig()

	d := cfg.Default()

	if d == nil || d.Code != "en" {
		t.Error("default locale should be en")
	}
}

func TestCodes(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig()

	codes := cfg.Codes()

	if len(codes) != 3 {
		t.Errorf("codes count = %d, want 3", len(codes))
	}
}

func TestMiddleware_DetectsPrefix(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig()

	var capturedPath string

	handler := i18n.Middleware(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/es/dashboard", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if capturedPath != "/dashboard" {
		t.Errorf("path = %q, want %q", capturedPath, "/dashboard")
	}
}

func TestMiddleware_DefaultLocale(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig()

	var capturedPath string

	handler := i18n.Middleware(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if capturedPath != "/dashboard" {
		t.Errorf("path = %q, want %q", capturedPath, "/dashboard")
	}
}

func TestMiddleware_StripsPrefixFromPath(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig()

	var capturedURI string

	handler := i18n.Middleware(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURI = r.RequestURI

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/ar/settings?tab=profile", nil)
	r.RequestURI = "/ar/settings?tab=profile"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if capturedURI != "/settings?tab=profile" {
		t.Errorf("RequestURI = %q, want %q", capturedURI, "/settings?tab=profile")
	}
}

func TestMiddleware_RootWithPrefix(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig()

	var capturedPath string

	handler := i18n.Middleware(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/es", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if capturedPath != "/" {
		t.Errorf("path = %q, want %q", capturedPath, "/")
	}
}

func TestMiddleware_UnknownPrefixFallsBackToDefault(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig()

	var capturedPath string

	var capturedLocale *seo.Locale

	handler := i18n.Middleware(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedLocale = seo.LocaleFromContext(r.Context())

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/zz/dashboard", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if capturedPath != "/zz/dashboard" {
		t.Errorf("path = %q, want %q (unchanged)", capturedPath, "/zz/dashboard")
	}

	if capturedLocale == nil || capturedLocale.Code != "en" {
		t.Error("expected fallback to default (en) locale")
	}
}

func TestMiddleware_HreflangTrimsTrailingSlash(t *testing.T) {
	t.Parallel()

	cfg := newTestConfig()

	var capturedLocale *seo.Locale

	handler := i18n.Middleware(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedLocale = seo.LocaleFromContext(r.Context())

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/es/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if capturedLocale == nil {
		t.Fatal("expected locale to be set")
	}

	for _, link := range capturedLocale.Head.Links {
		if link.Rel == "alternate" && strings.HasSuffix(link.Href, "/") && link.Href != "/" {
			t.Errorf("hreflang href %q has trailing slash", link.Href)
		}
	}
}
