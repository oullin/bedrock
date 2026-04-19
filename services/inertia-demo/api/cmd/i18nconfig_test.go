package main

import (
	"os"
	"path/filepath"
	"testing"
)

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

	cfg, err := LoadI18n(path)

	if err != nil {
		t.Fatal(err)
	}

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

	cfg, err := LoadI18n(path)

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

	_, err := LoadI18n(path)

	if err == nil {
		t.Fatal("expected error for missing default locale")
	}
}

func TestLoadI18n_FileNotFound(t *testing.T) {
	_, err := LoadI18n("/nonexistent/i18n.yml")

	if err == nil {
		t.Error("expected error for missing file")
	}
}
