package passport_test

import (
	"testing"
	"time"

	"github.com/bedrock/packages/passport"
)

func TestTokensCanRegistersScopes(t *testing.T) {
	p := newPassportWithScopes(map[string]string{
		"user:read":    "Read user data",
		"orders:write": "Write orders",
	})

	scopes := p.Scopes()

	if len(scopes) != 2 {
		t.Errorf("Scopes() len = %d, want 2", len(scopes))
	}
}

func TestFindScopeReturnsScope(t *testing.T) {
	p := newPassportWithScopes(map[string]string{
		"user:read": "Read user data",
	})

	s := p.FindScope("user:read")

	if s == nil {
		t.Fatal("FindScope returned nil for a registered scope")
	}

	if s.ID != "user:read" {
		t.Errorf("scope ID = %q, want %q", s.ID, "user:read")
	}

	if s.Description != "Read user data" {
		t.Errorf("scope Description = %q, want %q", s.Description, "Read user data")
	}
}

func TestFindScopeReturnsNilForMissing(t *testing.T) {
	p := newPassport()

	if p.FindScope("nonexistent") != nil {
		t.Error("FindScope should return nil for an unregistered scope")
	}
}

func TestHasScopeReturnsTrueForRegistered(t *testing.T) {
	p := newPassportWithScopes(map[string]string{"admin": "Admin access"})

	if !p.HasScope("admin") {
		t.Error("HasScope should return true for a registered scope")
	}
}

func TestHasScopeReturnsFalseForUnregistered(t *testing.T) {
	p := newPassport()

	if p.HasScope("missing") {
		t.Error("HasScope should return false for an unregistered scope")
	}
}

func TestScopeIDsReturnsAllIDs(t *testing.T) {
	p := newPassportWithScopes(map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	})

	ids := p.ScopeIDs()

	if len(ids) != 3 {
		t.Errorf("ScopeIDs() len = %d, want 3", len(ids))
	}

	idSet := make(map[string]bool)

	for _, id := range ids {
		idSet[id] = true
	}

	for _, want := range []string{"a", "b", "c"} {
		if !idSet[want] {
			t.Errorf("ScopeIDs missing %q", want)
		}
	}
}

func TestInheritedScopesEnabledDisabled(t *testing.T) {
	p := newPassport()

	if p.InheritedScopesEnabled() {
		t.Error("inherited scopes should be disabled by default")
	}

	p.UseInheritedScopes(true)

	if !p.InheritedScopesEnabled() {
		t.Error("inherited scopes should be enabled after UseInheritedScopes(true)")
	}

	p.UseInheritedScopes(false)

	if p.InheritedScopesEnabled() {
		t.Error("inherited scopes should be disabled after UseInheritedScopes(false)")
	}
}

func TestGrantEnableDisable(t *testing.T) {
	p := newPassport()

	grants := []struct {
		name    string
		enable  func()
		disable func()
		id      string
	}{
		{"auth_code", func() { p.EnableAuthorizationCodeGrant() }, func() { p.DisableAuthorizationCodeGrant() }, passport.GrantAuthorizationCode},
		{"password", func() { p.EnablePasswordGrant() }, func() { p.DisablePasswordGrant() }, passport.GrantPassword},
		{"client_credentials", func() { p.EnableClientCredentialsGrant() }, func() { p.DisableClientCredentialsGrant() }, passport.GrantClientCredentials},
		{"implicit", func() { p.EnableImplicitGrant() }, func() { p.DisableImplicitGrant() }, passport.GrantImplicit},
		{"device_code", func() { p.EnableDeviceCodeGrant() }, func() { p.DisableDeviceCodeGrant() }, passport.GrantDeviceCode},
		{"refresh_token", func() { p.EnableRefreshTokenGrant() }, func() { p.DisableRefreshTokenGrant() }, passport.GrantRefreshToken},
	}

	for _, g := range grants {
		if p.IsGrantEnabled(g.id) {
			t.Errorf("%s grant should be disabled by default", g.name)
		}

		g.enable()

		if !p.IsGrantEnabled(g.id) {
			t.Errorf("%s grant should be enabled after Enable()", g.name)
		}

		g.disable()

		if p.IsGrantEnabled(g.id) {
			t.Errorf("%s grant should be disabled after Disable()", g.name)
		}
	}
}

func TestTokensExpireIn(t *testing.T) {
	p := newPassport()

	p.TokensExpireIn(30 * time.Minute)

	if p.TokensTTL() != 30*time.Minute {
		t.Errorf("TokensTTL = %v, want %v", p.TokensTTL(), 30*time.Minute)
	}
}

func TestRefreshTokensExpireIn(t *testing.T) {
	p := newPassport()

	p.RefreshTokensExpireIn(7 * 24 * time.Hour)

	if p.RefreshTokensTTL() != 7*24*time.Hour {
		t.Errorf("RefreshTokensTTL = %v, want %v", p.RefreshTokensTTL(), 7*24*time.Hour)
	}
}

func TestPersonalAccessTokensExpireIn(t *testing.T) {
	p := newPassport()

	p.PersonalAccessTokensExpireIn(6 * 30 * 24 * time.Hour)

	if p.PersonalAccessTokensTTL() != 6*30*24*time.Hour {
		t.Errorf("PersonalAccessTokensTTL = %v, want %v", p.PersonalAccessTokensTTL(), 6*30*24*time.Hour)
	}
}

func TestActingAsAndClearActing(t *testing.T) {
	p := newPassport()
	user := newStubUser("u1")
	tok := newTestToken("t1", "u1", "c1", []string{"read"})
	at := passport.NewAccessToken(tok, p)

	p.ActingAs(user, at, []string{"read"})

	guard := passport.NewTokenGuard(p,
		passport.NewMemoryTokenStore(),
		passport.NewMemoryClientStore(),
		&stubProvider{},
	)

	got, _ := guard.User(nil)

	if got == nil || got.GetAuthIdentifier() != "u1" {
		t.Errorf("ActingAs: got user %v, want u1", got)
	}

	p.ClearActing()

	// After clearing, no request → nil user.
	guard2 := passport.NewTokenGuard(p,
		passport.NewMemoryTokenStore(),
		passport.NewMemoryClientStore(),
		&stubProvider{},
	)

	got2, _ := guard2.User(nil)

	if got2 != nil {
		t.Error("expected nil user after ClearActing")
	}
}

func TestTokensCanIsChainable(t *testing.T) {
	p := passport.NewPassport(nil).
		TokensCan(map[string]string{"read": "Read"}).
		TokensCan(map[string]string{"write": "Write"}).
		UseInheritedScopes(true)

	if !p.HasScope("read") || !p.HasScope("write") {
		t.Error("TokensCan should be chainable and additive")
	}

	if !p.InheritedScopesEnabled() {
		t.Error("UseInheritedScopes should be chainable")
	}
}
