package passport

import "time"

// AuthCode is the persisted authorization code record, corresponding to the
// oauth_auth_codes table in Laravel Passport.
//
// Authorization codes are short-lived tokens exchanged for access tokens in
// the Authorization Code Grant flow.
type AuthCode struct {
	ID        string
	UserID    string
	ClientID  string
	Scopes    []string
	Revoked   bool
	ExpiresAt time.Time
}
