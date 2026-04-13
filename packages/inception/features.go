package inception

// Features controls which Inception capabilities are active.
type Features struct {
	Registration             bool
	ResetPasswords           bool
	EmailVerification        bool
	UpdateProfileInformation bool
	UpdatePasswords          bool
	TwoFactorAuthentication  bool
	ConfirmPassword          bool
	Teams                    bool
	APITokens                bool
	ProfilePhotos            bool
	AccountDeletion          bool
	BrowserSessions          bool
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
		Teams:                    true,
		APITokens:                true,
		ProfilePhotos:            true,
		AccountDeletion:          true,
		BrowserSessions:          true,
	}
}
