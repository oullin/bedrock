package passport_test

import (
	"testing"

	"github.com/bedrock/packages/oauthserver"
)

func TestTokenCanReturnsTrueForPresentScope(t *testing.T) {
	tok := newTestToken("1", "u1", "c1", []string{"user:read"})
	tok = tok.WithOAuthServer(newOAuthServer())
	at := oauthserver.NewAccessToken(tok, newOAuthServer())

	user := newStubUser("u1")
	uwt := oauthserver.NewUserWithTokens(user, at)

	if !uwt.TokenCan("user:read") {
		t.Error("TokenCan should return true for a scope the token holds")
	}
}

func TestTokenCantReturnsTrueForAbsentScope(t *testing.T) {
	tok := newTestToken("1", "u1", "c1", []string{"user:read"})
	tok = tok.WithOAuthServer(newOAuthServer())
	at := oauthserver.NewAccessToken(tok, newOAuthServer())

	user := newStubUser("u1")
	uwt := oauthserver.NewUserWithTokens(user, at)

	if !uwt.TokenCant("orders:write") {
		t.Error("TokenCant should return true for a scope the token does not hold")
	}
}

func TestTokenCanReturnsFalseForAbsentScope(t *testing.T) {
	tok := newTestToken("1", "u1", "c1", []string{"user:read"})
	tok = tok.WithOAuthServer(newOAuthServer())
	at := oauthserver.NewAccessToken(tok, newOAuthServer())

	user := newStubUser("u1")
	uwt := oauthserver.NewUserWithTokens(user, at)

	if uwt.TokenCan("admin") {
		t.Error("TokenCan should return false for a scope the token does not hold")
	}
}

func TestUserWithTokensWrapsAuthenticatable(t *testing.T) {
	user := newStubUser("u42")
	uwt := oauthserver.NewUserWithTokens(user, nil)

	if uwt.GetAuthIdentifier() != "u42" {
		t.Errorf("GetAuthIdentifier = %q, want %q", uwt.GetAuthIdentifier(), "u42")
	}
}

func TestUserWithTokensCurrentAccessToken(t *testing.T) {
	tok := newTestToken("t1", "u1", "c1", nil)
	at := oauthserver.NewAccessToken(tok, newOAuthServer())
	user := newStubUser("u1")
	uwt := oauthserver.NewUserWithTokens(user, at)

	if uwt.CurrentAccessToken() != at {
		t.Error("CurrentAccessToken should return the attached access token")
	}
}

func TestUserWithTokensWithAccessToken(t *testing.T) {
	user := newStubUser("u1")
	uwt := oauthserver.NewUserWithTokens(user, nil)

	tok := newTestToken("t1", "u1", "c1", []string{"read"})
	at := oauthserver.NewAccessToken(tok, newOAuthServer())
	uwt2 := uwt.WithAccessToken(at)

	if uwt2.CurrentAccessToken() != at {
		t.Error("WithAccessToken should return a new UserWithTokens with the given token")
	}

	// Original is unchanged.
	if uwt.CurrentAccessToken() != nil {
		t.Error("WithAccessToken should not mutate the original UserWithTokens")
	}
}

func TestUserWithTokensTokenCanNilToken(t *testing.T) {
	user := newStubUser("u1")
	uwt := oauthserver.NewUserWithTokens(user, nil)

	if uwt.TokenCan("anything") {
		t.Error("TokenCan should return false when no token is attached")
	}
}

func TestHasApiTokensInterface(t *testing.T) {
	user := newStubUser("u1")
	uwt := oauthserver.NewUserWithTokens(user, nil)

	// Verify UserWithTokens satisfies the HasApiTokens interface.
	var _ oauthserver.HasApiTokens = uwt
}
