package i18n_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/packages/inertia/i18n"
)

const testConfig = `
default_locale: "en"
url_prefix: true

locales:
  en:
    name: "English"
    direction: "ltr"
    head:
      title: "My App"
      meta:
        - name: "description"
          content: "English description"
  es:
    name: "Español"
    direction: "ltr"
    head:
      title: "Mi App"
      meta:
        - name: "description"
          content: "Descripción en español"
  ar:
    name: "العربية"
    direction: "rtl"
    head:
      title: "تطبيقي"
`

func writeConfigFile(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "i18n.yml")

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	return path
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, testConfig)
	cfg, err := i18n.LoadConfig(path)

	if err != nil {
		t.Fatal(err)
	}

	if cfg.DefaultLocale != "en" {
		t.Errorf("default locale = %q, want %q", cfg.DefaultLocale, "en")
	}

	if len(cfg.Locales) != 3 {
		t.Errorf("locales count = %d, want 3", len(cfg.Locales))
	}

	en := cfg.Lookup("en")

	if en == nil {
		t.Fatal("en locale not found")
	}

	if en.Code != "en" {
		t.Errorf("en.Code = %q, want %q", en.Code, "en")
	}

	if en.Name != "English" {
		t.Errorf("en.Name = %q, want %q", en.Name, "English")
	}

	if en.Head.Title != "My App" {
		t.Errorf("en.Head.Title = %q, want %q", en.Head.Title, "My App")
	}

	ar := cfg.Lookup("ar")

	if ar == nil {
		t.Fatal("ar locale not found")
	}

	if ar.Direction != "rtl" {
		t.Errorf("ar.Direction = %q, want %q", ar.Direction, "rtl")
	}
}

func TestLoadConfig_EnvOverride(t *testing.T) {
	t.Setenv("INERTIA_I18N_DEFAULT_LOCALE", "es")

	path := writeConfigFile(t, testConfig)
	cfg, err := i18n.LoadConfig(path)

	if err != nil {
		t.Fatal(err)
	}

	if cfg.DefaultLocale != "es" {
		t.Errorf("default locale = %q, want %q", cfg.DefaultLocale, "es")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, "default_locale: [\ninvalid")
	_, err := i18n.LoadConfig(path)

	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	t.Parallel()

	_, err := i18n.LoadConfig("/nonexistent/i18n.yml")

	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestDefault(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, testConfig)
	cfg, _ := i18n.LoadConfig(path)

	d := cfg.Default()

	if d == nil || d.Code != "en" {
		t.Error("default locale should be en")
	}
}

func TestCodes(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, testConfig)
	cfg, _ := i18n.LoadConfig(path)

	codes := cfg.Codes()

	if len(codes) != 3 {
		t.Errorf("codes count = %d, want 3", len(codes))
	}
}

func TestMiddleware_DetectsPrefix(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, testConfig)
	cfg, _ := i18n.LoadConfig(path)

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

	path := writeConfigFile(t, testConfig)
	cfg, _ := i18n.LoadConfig(path)

	var capturedPath string

	handler := i18n.Middleware(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	// No prefix detected, path should remain unchanged.
	if capturedPath != "/dashboard" {
		t.Errorf("path = %q, want %q", capturedPath, "/dashboard")
	}
}

func TestMiddleware_StripsPrefixFromPath(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, testConfig)
	cfg, _ := i18n.LoadConfig(path)

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
		t.Errorf("requestURI = %q, want %q", capturedURI, "/settings?tab=profile")
	}
}

func TestMiddleware_RootWithPrefix(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, testConfig)
	cfg, _ := i18n.LoadConfig(path)

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

	path := writeConfigFile(t, testConfig)
	cfg, _ := i18n.LoadConfig(path)

	var capturedPath string

	var capturedLocale *protocol.Locale

	handler := i18n.Middleware(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedLocale = protocol.LocaleFromContext(r.Context())

		w.WriteHeader(http.StatusOK)
	}))

	// "xx" is not a configured locale, so path should remain unchanged
	// and the default locale ("en") should be used.
	r := httptest.NewRequest(http.MethodGet, "/xx/dashboard", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if capturedPath != "/xx/dashboard" {
		t.Errorf("path = %q, want %q", capturedPath, "/xx/dashboard")
	}

	if capturedLocale == nil || capturedLocale.Code != "en" {
		t.Errorf("locale = %v, want default locale with code %q", capturedLocale, "en")
	}
}

func TestMiddleware_HreflangTrimsTrailingSlash(t *testing.T) {
	t.Parallel()

	path := writeConfigFile(t, testConfig)
	cfg, _ := i18n.LoadConfig(path)

	var capturedLocale *protocol.Locale

	handler := i18n.Middleware(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedLocale = protocol.LocaleFromContext(r.Context())

		w.WriteHeader(http.StatusOK)
	}))

	// Request /es/ — the trailing slash on the clean path ("/") combined
	// with the prefix should produce "/es" not "/es/".
	r := httptest.NewRequest(http.MethodGet, "/es/admin/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if capturedLocale == nil {
		t.Fatal("locale not set in context")
	}

	for _, link := range capturedLocale.Head.Links {
		if link.Rel == "alternate" && strings.HasSuffix(link.Href, "/") && link.Href != "/" {
			t.Errorf("hreflang href %q should not have trailing slash", link.Href)
		}
	}
}
