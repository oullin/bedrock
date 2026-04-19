package oauthserver

import "time"

// Token is the persisted OAuth2 access token record, corresponding to the
// oauth_access_tokens table in Upstream OAuthServer.
type Token struct {
	ID        string
	UserID    string
	ClientID  string
	Name      string
	Scopes    []string
	Revoked   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time

	oauthserver *OAuthServer
}

// WithOAuthServer attaches the OAuthServer config to this token so that scope
// resolution (including inherited scopes) can be performed via Can/Cant.
func (t *Token) WithOAuthServer(p *OAuthServer) *Token {
	t.oauthserver = p

	return t
}

// Can reports whether this token's scopes include the requested scope.
//
// Rules (matching Upstream OAuthServer Token::can):
//   - Can("*") always returns false — you cannot check for the wildcard literally.
//   - If the token carries the "*" wildcard scope it grants every non-wildcard check.
//   - When inherited scopes are enabled via OAuthServer.UseInheritedScopes, ancestor
//     scopes are resolved: Can("user:read") returns true for a ["user"] token.
func (t *Token) Can(scope string) bool {
	if scope == "*" {
		return false
	}

	for _, s := range t.Scopes {
		if s == "*" {
			return true
		}
	}

	var inherited bool

	if t.oauthserver != nil {
		inherited = t.oauthserver.InheritedScopesEnabled()
	}

	return scopeExistsIn(scope, t.Scopes, inherited)
}

// Cant is the inverse of Can.
func (t *Token) Cant(scope string) bool {
	return !t.Can(scope)
}

// Revoke marks the token as revoked in memory.
// To persist the revocation, call TokenStore.Revoke.
func (t *Token) Revoke() bool {
	t.Revoked = true

	return true
}

// IsRevoked reports whether the token has been revoked.
func (t *Token) IsRevoked() bool {
	return t.Revoked
}
