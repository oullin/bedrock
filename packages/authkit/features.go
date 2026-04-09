package authkit

// Features controls which AuthKit capabilities are active.
type Features struct {
	Teams           bool
	APITokens       bool
	ProfilePhotos   bool
	AccountDeletion bool
	BrowserSessions bool
}

// DefaultFeatures returns a Features with all capabilities enabled.
func DefaultFeatures() Features {
	return Features{
		Teams:           true,
		APITokens:       true,
		ProfilePhotos:   true,
		AccountDeletion: true,
		BrowserSessions: true,
	}
}
