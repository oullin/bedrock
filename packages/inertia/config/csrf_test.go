package config_test

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/inertia/config"
)

func TestDefaultCSRF(t *testing.T) {
	t.Parallel()

	cfg := config.DefaultCSRF()

	if cfg.CookieName != "XSRF-TOKEN" {
		t.Errorf("CookieName = %q, want %q", cfg.CookieName, "XSRF-TOKEN")
	}

	if cfg.Secure {
		t.Error("Secure should be false by default")
	}

	if cfg.SameSite != "lax" {
		t.Errorf("SameSite = %q, want %q", cfg.SameSite, "lax")
	}
}

func TestLoadCSRF(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "csrf.yml")

	content := `
secure: true
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadCSRF(path)

	if err != nil {
		t.Fatal(err)
	}

	if !cfg.Secure {
		t.Error("Secure should be true (from file)")
	}

	// Defaults should be applied for fields not in the file.
	if cfg.CookieName != "XSRF-TOKEN" {
		t.Errorf("CookieName = %q, want %q (default)", cfg.CookieName, "XSRF-TOKEN")
	}

	if cfg.SameSite != "lax" {
		t.Errorf("SameSite = %q, want %q (default)", cfg.SameSite, "lax")
	}
}

func TestLoadCSRF_EnvOverride(t *testing.T) {
	t.Setenv("INERTIA_CSRF_COOKIE_NAME", "MY-TOKEN")

	dir := t.TempDir()
	path := filepath.Join(dir, "csrf.yml")

	content := `
cookie_name: "FILE-TOKEN"
`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.LoadCSRF(path)

	if err != nil {
		t.Fatal(err)
	}

	if cfg.CookieName != "MY-TOKEN" {
		t.Errorf("CookieName = %q, want %q (env override)", cfg.CookieName, "MY-TOKEN")
	}
}

func TestLoadCSRF_FileNotFound(t *testing.T) {
	t.Parallel()

	_, err := config.LoadCSRF("/nonexistent/csrf.yml")

	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestCSRFConfig_Defaults_FillsEmptyFields(t *testing.T) {
	t.Parallel()

	cfg := config.CSRFConfig{}

	cfg.Defaults()

	if cfg.CookieName != "XSRF-TOKEN" {
		t.Errorf("CookieName = %q, want %q", cfg.CookieName, "XSRF-TOKEN")
	}

	if cfg.SameSite != "lax" {
		t.Errorf("SameSite = %q, want %q", cfg.SameSite, "lax")
	}
}

func TestCSRFConfig_Defaults_PreservesExistingValues(t *testing.T) {
	t.Parallel()

	cfg := config.CSRFConfig{
		CookieName: "CUSTOM-TOKEN",
		SameSite:   "strict",
	}

	cfg.Defaults()

	if cfg.CookieName != "CUSTOM-TOKEN" {
		t.Errorf("CookieName = %q, want %q", cfg.CookieName, "CUSTOM-TOKEN")
	}

	if cfg.SameSite != "strict" {
		t.Errorf("SameSite = %q, want %q", cfg.SameSite, "strict")
	}
}

func TestCSRFConfig_SameSiteMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  http.SameSite
	}{
		{"strict", http.SameSiteStrictMode},
		{"Strict", http.SameSiteStrictMode},
		{"none", http.SameSiteNoneMode},
		{"None", http.SameSiteNoneMode},
		{"lax", http.SameSiteLaxMode},
		{"Lax", http.SameSiteLaxMode},
		{"", http.SameSiteLaxMode},
		{"unknown", http.SameSiteLaxMode},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cfg := config.CSRFConfig{SameSite: tt.input}
			got := cfg.SameSiteMode()

			if got != tt.want {
				t.Errorf("SameSiteMode(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
