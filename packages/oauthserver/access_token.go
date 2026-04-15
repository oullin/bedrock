package oauthserver

import (
	"encoding/json"
	"time"
)

// AccessToken wraps a Token and provides JSON serialization.
// It is the value associated with a user via WithAccessToken and is consulted
// when checking scopes via TokenCan/TokenCant.
//
// This mirrors Upstream OAuthServer's AccessToken class.
type AccessToken struct {
	token    *Token
	oauthserver *OAuthServer
}

// NewAccessToken constructs an AccessToken from a Token and OAuthServer config.
func NewAccessToken(t *Token, p *OAuthServer) *AccessToken {
	return &AccessToken{
		token:    t,
		oauthserver: p,
	}
}

// Can reports whether this access token has the given scope.
// Delegates to the underlying Token.Can, which respects inherited scopes.
func (a *AccessToken) Can(scope string) bool {
	if a.token == nil {
		return false
	}

	return a.token.Can(scope)
}

// Cant is the inverse of Can.
func (a *AccessToken) Cant(scope string) bool {
	return !a.Can(scope)
}

// Token returns the underlying Token record.
func (a *AccessToken) Token() *Token {
	return a.token
}

// ToArray returns a map representation of the access token for serialization.
// Keys match the Upstream OAuthServer token attribute names.
func (a *AccessToken) ToArray() map[string]any {
	if a.token == nil {
		return map[string]any{}
	}

	t := a.token

	return map[string]any{
		"id":         t.ID,
		"user_id":    t.UserID,
		"client_id":  t.ClientID,
		"name":       t.Name,
		"scopes":     t.Scopes,
		"revoked":    t.Revoked,
		"created_at": formatTime(t.CreatedAt),
		"updated_at": formatTime(t.UpdatedAt),
		"expires_at": formatTime(t.ExpiresAt),
	}
}

// MarshalJSON implements json.Marshaler using ToArray.
func (a *AccessToken) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.ToArray())
}

// formatTime returns the time as a string in RFC3339 format, or empty string
// for zero values.
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.Format(time.RFC3339)
}
