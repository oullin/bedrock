package i18n

import (
	"github.com/bedrock/packages/seo"
)

// I18nConfig holds the multilanguage configuration consumed by Middleware.
// Construct it directly, via DefaultI18n(), or by deserialising a YAML/JSON
// document into it from your application's config layer.
type I18nConfig struct {
	DefaultLocale string                 `mapstructure:"default_locale" yaml:"default_locale" json:"default_locale"`
	URLPrefix     bool                   `mapstructure:"url_prefix"     yaml:"url_prefix"     json:"url_prefix"`
	Locales       map[string]*seo.Locale `mapstructure:"locales"        yaml:"locales"        json:"locales"`
}

// DefaultI18n returns an I18nConfig with a single English locale.
func DefaultI18n() *I18nConfig {
	return &I18nConfig{
		DefaultLocale: "en",
		URLPrefix:     false,
		Locales: map[string]*seo.Locale{
			"en": {
				Code:      "en",
				Name:      "English",
				Direction: "ltr",
			},
		},
	}
}

// Lookup returns the Locale for the given code, or nil if not found.
func (cfg *I18nConfig) Lookup(code string) *seo.Locale {
	return cfg.Locales[code]
}

// Default returns the default Locale.
func (cfg *I18nConfig) Default() *seo.Locale {
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
