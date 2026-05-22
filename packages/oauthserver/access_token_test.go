package oauthserver_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/bedrock/packages/oauthserver"
)

func TestAccessTokenCanScope(t *testing.T) {
	tok := newTestToken("1", "u1", "c1", []string{"user:read"})
	tok = tok.WithOAuthServer(newOAuthServer())
	at := oauthserver.NewAccessToken(tok, newOAuthServer())

	if !at.Can("user:read") {
		t.Error("Can should return true for a scope the token holds")
	}
}

func TestAccessTokenCantScope(t *testing.T) {
	tok := newTestToken("1", "u1", "c1", []string{"user:read"})
	tok = tok.WithOAuthServer(newOAuthServer())
	at := oauthserver.NewAccessToken(tok, newOAuthServer())

	if !at.Cant("orders:write") {
		t.Error("Cant should return true for a scope the token does not hold")
	}
}

func TestAccessTokenWildcardScopes(t *testing.T) {
	tok := newTestToken("1", "u1", "c1", []string{"*"})
	tok = tok.WithOAuthServer(newOAuthServer())
	at := oauthserver.NewAccessToken(tok, newOAuthServer())

	if !at.Can("any:scope") {
		t.Error("Can should return true for a wildcard token")
	}

	if at.Can("*") {
		t.Error("Can(\"*\") should always return false")
	}
}

func TestAccessTokenInheritedScopes(t *testing.T) {
	p := newOAuthServer().UseInheritedScopes(true)

	tok := newTestToken("1", "u1", "c1", []string{"admin"})
	tok = tok.WithOAuthServer(p)
	at := oauthserver.NewAccessToken(tok, p)

	if !at.Can("admin:webhooks:read") {
		t.Error("Can(\"admin:webhooks:read\") should be true for [\"admin\"] token with inherited scopes")
	}
}

func TestAccessTokenToArray(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	tok := &oauthserver.Token{
		ID:        "tok-123",
		UserID:    "user-1",
		ClientID:  "client-1",
		Name:      "My Token",
		Scopes:    []string{"read"},
		Revoked:   false,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(time.Hour),
	}

	at := oauthserver.NewAccessToken(tok, newOAuthServer())
	arr := at.ToArray()

	fields := []string{"id", "user_id", "client_id", "name", "scopes", "revoked", "created_at", "updated_at", "expires_at"}

	for _, f := range fields {
		if _, ok := arr[f]; !ok {
			t.Errorf("ToArray missing field %q", f)
		}
	}

	if arr["id"] != "tok-123" {
		t.Errorf("id = %v, want %q", arr["id"], "tok-123")
	}
}

func TestAccessTokenMarshalJSON(t *testing.T) {
	tok := newTestToken("t1", "u1", "c1", []string{"read"})
	at := oauthserver.NewAccessToken(tok, newOAuthServer())

	b, err := json.Marshal(at)

	if err != nil {
		t.Fatalf("MarshalJSON error: %v", err)
	}

	var got map[string]any

	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if got["id"] != "t1" {
		t.Errorf("JSON id = %v, want %q", got["id"], "t1")
	}
}

func TestAccessTokenTokenReturnsUnderlying(t *testing.T) {
	tok := newTestToken("t1", "u1", "c1", nil)
	at := oauthserver.NewAccessToken(tok, newOAuthServer())

	if at.Token() != tok {
		t.Error("Token() should return the underlying Token")
	}
}

func TestAccessTokenNilToken(t *testing.T) {
	at := oauthserver.NewAccessToken(nil, newOAuthServer())

	if at.Can("anything") {
		t.Error("Can should return false for a nil token")
	}

	if !at.Cant("anything") {
		t.Error("Cant should return true for a nil token")
	}
}
