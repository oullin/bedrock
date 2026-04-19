package main

import (
	"fmt"
	"strings"

	"github.com/bedrock/packages/config"
	"github.com/bedrock/packages/seo"
	"github.com/bedrock/packages/seo/i18n"
)

// LoadI18n reads a YAML i18n config file and returns a populated I18nConfig.
// Defaults are applied first, then the file values are merged on top, and
// finally environment variable overrides (INERTIA_I18N_*) are applied.
func LoadI18n(path string) (*i18n.I18nConfig, error) {
	repo := config.New(nil)
	v := repo.Viper()

	v.SetDefault("default_locale", "en")
	v.SetDefault("url_prefix", false)
	v.SetEnvPrefix("INERTIA_I18N")
	v.AutomaticEnv()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("i18n: read config: %w", err)
	}

	cfg := i18n.DefaultI18n()

	if v.IsSet("default_locale") {
		cfg.DefaultLocale = strings.TrimSpace(v.GetString("default_locale"))
	}

	if v.IsSet("url_prefix") {
		cfg.URLPrefix = v.GetBool("url_prefix")
	}

	if v.IsSet("locales") {
		var locales map[string]*seo.Locale

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
