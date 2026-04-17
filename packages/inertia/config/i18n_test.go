package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/inertia/config"
)

func TestDefaultI18n(t *testing.T) {
	t.Parallel()

	cfg := config.DefaultI18n()

	if cfg.DefaultLocale != "en" {
		t.Errorf("DefaultLocale = %q, want %q", cfg.DefaultLocale, "en")
	}

	if cfg.URLPrefix {
		t.Error("URLPrefix should be false by default")
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

	if en.Direction != "ltr" {
		t.Errorf("en.Direction = %q, want %q", en.Direction, "ltr")
	}
}

func TestLoadI18n(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "i18n.yml")

	content := `
url_prefix: true

locales:
  es:
    name: "Español"
    direction: "ltr"
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadI18n(path)

	if err != nil {
		t.Fatal(err)
	}

	// Default should apply for fields not in the file.
	if cfg.DefaultLocale != "en" {
		t.Errorf("DefaultLocale = %q, want %q (default)", cfg.DefaultLocale, "en")
	}

	if !cfg.URLPrefix {
		t.Error("URLPrefix should be true (from file)")
	}

	es := cfg.Lookup("es")

	if es == nil {
		t.Fatal("es locale not found")
	}

	if es.Code != "es" {
		t.Errorf("es.Code = %q, want %q", es.Code, "es")
	}

	if es.Name != "Español" {
		t.Errorf("es.Name = %q, want %q", es.Name, "Español")
	}
}

func TestLoadI18n_EnvOverride(t *testing.T) {
	t.Setenv("INERTIA_I18N_DEFAULT_LOCALE", "es")

	dir := t.TempDir()
	path := filepath.Join(dir, "i18n.yml")

	content := `
default_locale: "en"
locales:
  en:
    name: "English"
    direction: "ltr"
  es:
    name: "Español"
    direction: "ltr"
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadI18n(path)

	if err != nil {
		t.Fatal(err)
	}

	if cfg.DefaultLocale != "es" {
		t.Errorf("DefaultLocale = %q, want %q (env override)", cfg.DefaultLocale, "es")
	}
}

func TestLoadI18n_InvalidDefaultLocale(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "i18n.yml")

	content := `
default_locale: "es"
locales:
  en:
    name: "English"
    direction: "ltr"
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := config.LoadI18n(path)

	if err == nil {
		t.Fatal("expected error for missing default locale")
	}
}

func TestLoadI18n_FileNotFound(t *testing.T) {
	_, err := config.LoadI18n("/nonexistent/i18n.yml")

	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestI18nConfig_Codes(t *testing.T) {
	t.Parallel()

	cfg := config.DefaultI18n()

	codes := cfg.Codes()

	if len(codes) != 1 {
		t.Errorf("codes count = %d, want 1", len(codes))
	}

	if codes[0] != "en" {
		t.Errorf("codes[0] = %q, want %q", codes[0], "en")
	}
}

func TestI18nConfig_Default(t *testing.T) {
	t.Parallel()

	cfg := config.DefaultI18n()

	d := cfg.Default()

	if d == nil || d.Code != "en" {
		t.Error("default locale should be en")
	}
}
