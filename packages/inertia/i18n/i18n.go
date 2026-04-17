package i18n

import (
	"github.com/bedrock/packages/inertia/config"
)

// LoadConfig reads a YAML i18n config file. Defaults are applied first,
// then the file values are merged on top, and finally env var overrides
// (INERTIA_I18N_*) are applied.
func LoadConfig(path string) (*config.I18nConfig, error) {
	return config.LoadI18n(path)
}
