package authflows

// Feature constants mirror Upstream AuthFlows features.
const (
	FeatureRegistration             = "registration"
	FeatureResetPasswords           = "reset-passwords"
	FeatureEmailVerification        = "email-verification"
	FeatureUpdateProfileInformation = "update-profile-information"
	FeatureUpdatePasswords          = "update-passwords"
	FeatureTwoFactorAuthentication  = "two-factor-authentication"
)

// Features reports AuthFlows feature enablement.
type Features struct {
	config Config
}

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
	return f.config.Options[feature][option]
}
