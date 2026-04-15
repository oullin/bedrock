package passport_test

import (
	"context"
	"time"

	"github.com/bedrock/packages/passport"
	cauth "github.com/bedrock/packages/contracts/auth"
)

// ---- stubProvider ------------------------------------------------------------

// stubProvider is a minimal cauth.UserProvider backed by a map.
type stubProvider struct {
	users map[string]cauth.Authenticatable
}

func (p *stubProvider) RetrieveByID(_ context.Context, id string) (cauth.Authenticatable, error) {
	if u, ok := p.users[id]; ok {
		return u, nil
	}

	return nil, nil
}

func (p *stubProvider) RetrieveByToken(_ context.Context, id, _ string) (cauth.Authenticatable, error) {
	return p.RetrieveByID(context.Background(), id)
}

func (p *stubProvider) RetrieveByCredentials(_ context.Context, creds map[string]string) (cauth.Authenticatable, error) {
	for _, u := range p.users {
		if u.GetAuthIdentifier() == creds["id"] {
			return u, nil
		}
	}

	return nil, nil
}

func (p *stubProvider) UpdateRememberToken(_ context.Context, _ cauth.Authenticatable, _ string) error {
	return nil
}

func (p *stubProvider) ValidateCredentials(_ context.Context, _ cauth.Authenticatable, _ map[string]string) (bool, error) {
	return false, nil
}

func (p *stubProvider) RehashPasswordIfRequired(_ context.Context, _ cauth.Authenticatable, _ map[string]string, _ bool) error {
	return nil
}

// ---- stubUser ----------------------------------------------------------------

// stubUser is a minimal cauth.Authenticatable implementation.
type stubUser struct {
	id       string
	password string
	token    string
}

func newStubUser(id string) *stubUser {
	return &stubUser{id: id}
}

func (u *stubUser) GetAuthIdentifierName() string { return "id" }
func (u *stubUser) GetAuthIdentifier() string     { return u.id }
func (u *stubUser) GetAuthPasswordName() string   { return "password" }
func (u *stubUser) GetAuthPassword() string       { return u.password }
func (u *stubUser) SetAuthPassword(p string)      { u.password = p }
func (u *stubUser) GetRememberToken() string      { return u.token }
func (u *stubUser) SetRememberToken(t string)     { u.token = t }
func (u *stubUser) GetRememberTokenName() string  { return "remember_token" }

// ---- helpers -----------------------------------------------------------------

// newTestToken builds an active (non-revoked, non-expired) Token.
func newTestToken(id, userID, clientID string, scopes []string) *passport.Token {
	now := time.Now()

	return &passport.Token{
		ID:        id,
		UserID:    userID,
		ClientID:  clientID,
		Name:      "test",
		Scopes:    scopes,
		Revoked:   false,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(time.Hour),
	}
}

// newExpiredToken builds a Token that has already expired.
func newExpiredToken(id, userID, clientID string, scopes []string) *passport.Token {
	past := time.Now().Add(-time.Hour)

	return &passport.Token{
		ID:        id,
		UserID:    userID,
		ClientID:  clientID,
		Scopes:    scopes,
		Revoked:   false,
		ExpiresAt: past,
	}
}

// newRevokedToken builds a Token that has been revoked.
func newRevokedToken(id, userID, clientID string, scopes []string) *passport.Token {
	t := newTestToken(id, userID, clientID, scopes)
	t.Revoked = true

	return t
}

// newPassport creates a fresh Passport with no config (suitable for unit tests).
func newPassport() *passport.Passport {
	return passport.NewPassport(nil)
}

// newPassportWithScopes creates a Passport with the given scope map registered.
func newPassportWithScopes(scopes map[string]string) *passport.Passport {
	return passport.NewPassport(nil).TokensCan(scopes)
}
