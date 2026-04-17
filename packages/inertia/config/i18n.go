package config

import (
	"fmt"
	"strings"

	"github.com/bedrock/packages/inertia/protocol"
	"github.com/spf13/viper"
)

// I18nConfig holds the multilanguage configuration.
type I18nConfig struct {
	DefaultLocale string                   `mapstructure:"default_locale"`
	URLPrefix     bool                     `mapstructure:"url_prefix"`
	Locales       map[string]*protocol.Locale `mapstructure:"locales"`
}

// DefaultI18n returns an I18nConfig with a single English locale.
func DefaultI18n() *I18nConfig {
	cfg := &I18nConfig{
		DefaultLocale: "en",
		URLPrefix:     false,
		Locales: map[string]*protocol.Locale{
			"en": {
				Code:      "en",
				Name:      "English",
				Direction: "ltr",
			},
		},
	}

	return cfg
}

// LoadI18n reads a YAML i18n config file. Defaults are applied first,
// then the file values are merged on top, and finally environment
// variable overrides (INERTIA_I18N_*) are applied.
func LoadI18n(path string) (*I18nConfig, error) {
	v := viper.New()

	v.SetDefault("default_locale", "en")
	v.SetDefault("url_prefix", false)

	v.SetEnvPrefix("INERTIA_I18N")

	v.AutomaticEnv()

	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("i18n: read config: %w", err)
	}

	cfg := DefaultI18n()

	if v.IsSet("default_locale") {
		cfg.DefaultLocale = strings.TrimSpace(v.GetString("default_locale"))
	}

	if v.IsSet("url_prefix") {
		cfg.URLPrefix = v.GetBool("url_prefix")
	}

	if v.IsSet("locales") {
		var locales map[string]*protocol.Locale

		if err := v.UnmarshalKey("locales", &locales); err != nil {
			return nil, fmt.Errorf("i18n: parse locales: %w", err)
		}

		for code, locale := range locales {
			if locale == nil {
				return nil, fmt.Errorf("i18n: locale %q is required", code)
			}

			cfg.Locales[code] = locale
		}
	}

	for code, locale := range cfg.Locales {
		if locale == nil {
			return nil, fmt.Errorf("i18n: locale %q is required", code)
		}

		locale.Code = code
	}

	if strings.TrimSpace(cfg.DefaultLocale) == "" {
		cfg.DefaultLocale = "en"
	}

	if cfg.Default() == nil {
		return nil, fmt.Errorf("i18n: default locale %q is not configured", cfg.DefaultLocale)
	}

	return cfg, nil
}

// Lookup returns the Locale for the given code, or nil if not found.
func (cfg *I18nConfig) Lookup(code string) *protocol.Locale {
	return cfg.Locales[code]
}

// Default returns the default Locale.
func (cfg *I18nConfig) Default() *protocol.Locale {
	return cfg.Locales[cfg.DefaultLocale]
}

// Codes returns all configured locale codes.
func (cfg *I18nConfig) Codes() []string {
	codes := make([]string, 0, len(cfg.Locales))

	for code := range cfg.Locales {
		codes = append(codes, code)
	}

	return codes
}
