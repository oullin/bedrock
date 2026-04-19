package i18n_test

import (
	"testing"

	"github.com/bedrock/packages/seo/i18n"
)

func TestDefaultI18n(t *testing.T) {
	t.Parallel()

	cfg := i18n.DefaultI18n()

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

func TestI18nConfig_Codes(t *testing.T) {
	t.Parallel()

	cfg := i18n.DefaultI18n()

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

	cfg := i18n.DefaultI18n()

	d := cfg.Default()

	if d == nil || d.Code != "en" {
		t.Error("default locale should be en")
	}
}
