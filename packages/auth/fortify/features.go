package fortify

import "strconv"

// Feature constants mirror Laravel Fortify features.

// Features reports Fortify feature enablement.
type Features struct {
	config Config
}

const (
	FeatureRegistration             = "registration"
	FeatureResetPasswords           = "reset-passwords"
	FeatureEmailVerification        = "email-verification"
	FeatureUpdateProfileInformation = "update-profile-information"
	FeatureUpdatePasswords          = "update-passwords"
	FeatureTwoFactorAuthentication  = "two-factor-authentication"
)

// NewFeatures creates a feature set from config.
func NewFeatures(config Config) Features {
	return Features{config: config}
}

// Enabled reports whether a feature is enabled.
func (f Features) Enabled(feature string) bool {
	for _, enabled := range f.config.Features {
		if enabled == feature {
			return true
		}
	}

	return false
}

// OptionEnabled reports whether a feature option is enabled.
func (f Features) OptionEnabled(feature string, option string) bool {
	if !f.Enabled(feature) {
		return false
	}

	switch value := f.config.Options[feature][option].(type) {
	case bool:
		return value
	case string:
		parsed, err := strconv.ParseBool(value)

		return err == nil && parsed
	default:
		return false
	}
}

// OptionInt returns an integer feature option or the supplied fallback.
func (f Features) OptionInt(feature string, option string, fallback int) int {
	if !f.Enabled(feature) {
		return fallback
	}

	switch value := f.config.Options[feature][option].(type) {
	case int:
		return value
	case int64:
		return int(value)
	case float64:
		return int(value)
	case string:
		parsed, err := strconv.Atoi(value)

		if err == nil {
			return parsed
		}
	}

	return fallback
}
