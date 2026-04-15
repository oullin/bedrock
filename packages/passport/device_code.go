package passport

import "time"

// DeviceCode is the persisted device authorization code record, corresponding
// to the oauth_device_codes table in Laravel Passport.
//
// Device codes support the Device Authorization Grant flow (RFC 8628), used by
// input-constrained devices such as smart TVs or CLI tools.
type DeviceCode struct {
	ID           string
	ClientID     string
	UserID       string
	DeviceCode   string
	UserCode     string
	Scopes       []string
	Approved     bool
	Revoked      bool
	LastPolledAt *time.Time
	ExpiresAt    time.Time
}
