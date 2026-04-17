package config

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// CSRFConfig holds configuration for the CSRF middleware.
type CSRFConfig struct {
	CookieName    string `json:"cookie_name"    yaml:"cookie_name"    mapstructure:"cookie_name"`
	Secure        bool   `json:"secure"         yaml:"secure"         mapstructure:"secure"`
	SameSite      string `json:"same_site"      yaml:"same_site"      mapstructure:"same_site"`
	OriginOnly    bool   `json:"origin_only"    yaml:"origin_only"    mapstructure:"origin_only"`
	AllowSameSite bool   `json:"allow_same_site" yaml:"allow_same_site" mapstructure:"allow_same_site"`
}

// DefaultCSRF returns a CSRFConfig with sensible defaults.
func DefaultCSRF() CSRFConfig {
	return CSRFConfig{
		CookieName:    "XSRF-TOKEN",
		Secure:        false,
		SameSite:      "lax",
		OriginOnly:    false,
		AllowSameSite: false,
	}
}

// LoadCSRF reads a YAML config file and returns a CSRFConfig. Defaults
// are applied first, then the file values are merged on top, and finally
// environment variable overrides (INERTIA_CSRF_*) are applied.
func LoadCSRF(path string) (CSRFConfig, error) {
	defaults := DefaultCSRF()

	v := viper.New()

	v.SetDefault("cookie_name", defaults.CookieName)
	v.SetDefault("secure", defaults.Secure)
	v.SetDefault("same_site", defaults.SameSite)
	v.SetDefault("origin_only", defaults.OriginOnly)
	v.SetDefault("allow_same_site", defaults.AllowSameSite)

	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return CSRFConfig{}, fmt.Errorf("csrf: read config: %w", err)
	}

	v.SetEnvPrefix("INERTIA_CSRF")

	v.AutomaticEnv()

	var cfg CSRFConfig

	if err := v.Unmarshal(&cfg); err != nil {
		return CSRFConfig{}, fmt.Errorf("csrf: parse config: %w", err)
	}

	return cfg, nil
}

func (c *CSRFConfig) Defaults() {
	if strings.TrimSpace(c.CookieName) == "" {
		c.CookieName = "XSRF-TOKEN"
	}

	if strings.TrimSpace(c.SameSite) == "" {
		c.SameSite = "lax"
	}
}

// SameSiteMode returns the http.SameSite value for the configured SameSite string.
func (c *CSRFConfig) SameSiteMode() http.SameSite {
	switch strings.ToLower(c.SameSite) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
