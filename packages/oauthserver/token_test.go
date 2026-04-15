package passport_test

import (
	"testing"

	"github.com/bedrock/packages/oauthserver"
)

func TestTokenCanReturnsTrueForMatchingScope(t *testing.T) {
	tok := newTestToken("1", "u1", "c1", []string{"user:read", "orders:write"})
	tok = tok.WithOAuthServer(newOAuthServer())

	if !tok.Can("user:read") {
		t.Error("Can should return true for a scope the token holds")
	}
}

func TestTokenCantReturnsTrueForMissingScope(t *testing.T) {
	tok := newTestToken("1", "u1", "c1", []string{"user:read"})
	tok = tok.WithOAuthServer(newOAuthServer())

	if !tok.Cant("orders:write") {
		t.Error("Cant should return true for a scope the token does not hold")
	}
}

func TestTokenWildcardInScopesGrantsAllPermissions(t *testing.T) {
	tok := newTestToken("1", "u1", "c1", []string{"*"})
	tok = tok.WithOAuthServer(newOAuthServer())

	for _, scope := range []string{"user:read", "orders:write", "admin:delete", "any:scope"} {
		if !tok.Can(scope) {
			t.Errorf("Can(%q) should return true for a wildcard token", scope)
		}
	}
}

func TestTokenCanReturnsFalseForWildcardScope(t *testing.T) {
	tok := newTestToken("1", "u1", "c1", []string{"user:read"})
	tok = tok.WithOAuthServer(newOAuthServer())

	// Can("*") is always false — you cannot check for the literal wildcard.
	if tok.Can("*") {
		t.Error("Can(\"*\") should always return false")
	}
}

func TestTokenWildcardScopeCanReturnsFalseForWildcardCheck(t *testing.T) {
	// Even a wildcard token cannot return true for Can("*").
	tok := newTestToken("1", "u1", "c1", []string{"*"})
	tok = tok.WithOAuthServer(newOAuthServer())

	if tok.Can("*") {
		t.Error("Can(\"*\") should always return false even for a wildcard token")
	}
}

func TestTokenInheritedScopesEnabled(t *testing.T) {
	p := newOAuthServer().UseInheritedScopes(true)

	tok := newTestToken("1", "u1", "c1", []string{"user"})
	tok = tok.WithOAuthServer(p)

	if !tok.Can("user:read") {
		t.Error("Can(\"user:read\") should be true when token has \"user\" and inherited scopes are enabled")
	}

	if !tok.Can("user") {
		t.Error("Can(\"user\") should still be true for direct scope match")
	}
}

func TestTokenInheritedScopeResolution(t *testing.T) {
	// "admin:webhooks:read" should resolve to ["admin", "admin:webhooks", "admin:webhooks:read"]
	// so a token with just ["admin"] grants "admin:webhooks:read" when inherited is on.
	p := newOAuthServer().UseInheritedScopes(true)

	cases := []struct {
		tokenScopes []string
		check       string
		want        bool
	}{
		{[]string{"admin"}, "admin:webhooks:read", true},
		{[]string{"admin:webhooks"}, "admin:webhooks:read", true},
		{[]string{"admin:webhooks:read"}, "admin:webhooks:read", true},
		{[]string{"orders"}, "admin:webhooks:read", false},
		{[]string{"admin"}, "admin", true},
	}

	for _, tc := range cases {
		tok := newTestToken("1", "u1", "c1", tc.tokenScopes)
		tok = tok.WithOAuthServer(p)

		got := tok.Can(tc.check)

		if got != tc.want {
			t.Errorf("token scopes=%v Can(%q) = %v, want %v", tc.tokenScopes, tc.check, got, tc.want)
		}
	}
}

func TestTokenInheritedScopesDisabledByDefault(t *testing.T) {
	// Without UseInheritedScopes, ancestor scopes do NOT grant children.
	tok := newTestToken("1", "u1", "c1", []string{"user"})
	tok = tok.WithOAuthServer(newOAuthServer())

	if tok.Can("user:read") {
		t.Error("Can(\"user:read\") should be false when inherited scopes are disabled")
	}
}

func TestTokenIsRevoked(t *testing.T) {
	tok := newRevokedToken("1", "u1", "c1", nil)

	if !tok.IsRevoked() {
		t.Error("IsRevoked should return true for a revoked token")
	}
}

func TestTokenRevoke(t *testing.T) {
	tok := newTestToken("1", "u1", "c1", nil)

	if tok.IsRevoked() {
		t.Fatal("token should not be revoked initially")
	}

	tok.Revoke()

	if !tok.IsRevoked() {
		t.Error("IsRevoked should return true after Revoke()")
	}
}

func TestTokenCantIsInverseOfCan(t *testing.T) {
	tok := newTestToken("1", "u1", "c1", []string{"read"})
	tok = tok.WithOAuthServer(newOAuthServer())

	if tok.Can("read") == tok.Cant("read") {
		t.Error("Can and Cant should always return opposite values")
	}

	if tok.Can("write") == tok.Cant("write") {
		t.Error("Can and Cant should always return opposite values for missing scope")
	}
}

func TestTokenWithoutOAuthServerUsesDirectScopeMatch(t *testing.T) {
	tok := &oauthserver.Token{
		Scopes: []string{"user"},
	}

	if !tok.Can("user") {
		t.Error("Can should still work without a oauthserver attached")
	}

	if tok.Can("user:read") {
		t.Error("Can(\"user:read\") should be false without a oauthserver (no inherited scope resolution)")
	}
}
