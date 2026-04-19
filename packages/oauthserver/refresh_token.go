package oauthserver

import "time"

// RefreshToken is the persisted refresh token record, corresponding to the
// oauth_refresh_tokens table in Upstream OAuthServer.
type RefreshToken struct {
	ID            string
	AccessTokenID string
	Revoked       bool
	ExpiresAt     time.Time
}
