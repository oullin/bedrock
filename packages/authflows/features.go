package authflows

// Features controls which AuthFlows capabilities are active.
type Features struct {
	Registration             bool
	ResetPasswords           bool
	EmailVerification        bool
	UpdateProfileInformation bool
	UpdatePasswords          bool
	TwoFactorAuthentication  bool
	ConfirmPassword          bool
}

// DefaultFeatures returns a Features with all capabilities enabled.
func DefaultFeatures() Features {
	return Features{
		Registration:             true,
		ResetPasswords:           true,
		EmailVerification:        true,
		UpdateProfileInformation: true,
		UpdatePasswords:          true,
		TwoFactorAuthentication:  true,
		ConfirmPassword:          true,
	}
}
