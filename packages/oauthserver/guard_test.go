package passport_test

import (
	"context"
	"net/http/httptest"
	"testing"

	cauth "github.com/bedrock/packages/contracts/auth"
	"github.com/bedrock/packages/oauthserver"
)

func newGuard(p *oauthserver.OAuthServer, tokens oauthserver.TokenStore, clients oauthserver.ClientStore, users map[string]cauth.Authenticatable) *oauthserver.TokenGuard {
	provider := &stubProvider{users: users}

	return oauthserver.NewTokenGuard(p, tokens, clients, provider)
}

func TestTokenGuardReturnsUserForBearerToken(t *testing.T) {
	p := newOAuthServer()
	user := newStubUser("u1")
	tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)
	clients := oauthserver.NewMemoryClientStore()

	tok := newTestToken("tok-abc", "u1", "c1", []string{"read"})
	_ = tokens.Save(context.Background(), tok)

	guard := newGuard(p, tokens, clients, map[string]cauth.Authenticatable{"u1": user})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer tok-abc")
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got == nil {
		t.Fatal("expected authenticated user, got nil")
	}

	if got.GetAuthIdentifier() != "u1" {
		t.Errorf("user id = %q, want %q", got.GetAuthIdentifier(), "u1")
	}
}

func TestTokenGuardReturnsNilForRevokedToken(t *testing.T) {
	p := newOAuthServer()
	tokens := oauthserver.NewMemoryTokenStore()
	clients := oauthserver.NewMemoryClientStore()

	tok := newRevokedToken("tok-rev", "u1", "c1", nil)
	_ = tokens.Save(context.Background(), tok)

	guard := newGuard(p, tokens, clients, map[string]cauth.Authenticatable{})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer tok-rev")
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Error("expected nil user for revoked token")
	}
}

func TestTokenGuardReturnsNilForExpiredToken(t *testing.T) {
	p := newOAuthServer()
	tokens := oauthserver.NewMemoryTokenStore()
	clients := oauthserver.NewMemoryClientStore()

	tok := newExpiredToken("tok-exp", "u1", "c1", nil)
	_ = tokens.Save(context.Background(), tok)

	guard := newGuard(p, tokens, clients, map[string]cauth.Authenticatable{})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer tok-exp")
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Error("expected nil user for expired token")
	}
}

func TestTokenGuardReturnsNilWhenNoToken(t *testing.T) {
	p := newOAuthServer()
	guard := newGuard(p, oauthserver.NewMemoryTokenStore(), oauthserver.NewMemoryClientStore(), nil)

	req := httptest.NewRequest("GET", "/", nil)
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Error("expected nil user when no Authorization header")
	}
}

func TestTokenGuardClientReturnsClientForToken(t *testing.T) {
	p := newOAuthServer()
	user := newStubUser("u1")
	tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)
	clients := oauthserver.NewMemoryClientStore()

	client := &oauthserver.Client{ID: "c1", Name: "MyApp", Secret: "s"}
	_ = clients.Create(context.Background(), client)

	tok := newTestToken("tok-1", "u1", "c1", nil)
	_ = tokens.Save(context.Background(), tok)

	guard := newGuard(p, tokens, clients, map[string]cauth.Authenticatable{"u1": user})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer tok-1")
	guard.SetRequest(req)

	// Resolve user first to populate the cached token.
	_, _ = guard.User(context.Background())

	got, err := guard.Client(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got == nil {
		t.Fatal("expected client, got nil")
	}

	if got.ID != "c1" {
		t.Errorf("client id = %q, want %q", got.ID, "c1")
	}
}

func TestTokenGuardClientReturnsNilForRevokedToken(t *testing.T) {
	p := newOAuthServer()
	tokens := oauthserver.NewMemoryTokenStore()
	clients := oauthserver.NewMemoryClientStore()

	tok := newRevokedToken("tok-rev", "u1", "c1", nil)
	_ = tokens.Save(context.Background(), tok)

	guard := newGuard(p, tokens, clients, map[string]cauth.Authenticatable{})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer tok-rev")
	guard.SetRequest(req)

	got, err := guard.Client(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Error("expected nil client for revoked token")
	}
}

func TestTokenGuardCachesUserPerRequest(t *testing.T) {
	p := newOAuthServer()
	user := newStubUser("u1")
	tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)
	clients := oauthserver.NewMemoryClientStore()

	tok := newTestToken("tok-1", "u1", "c1", nil)
	_ = tokens.Save(context.Background(), tok)

	guard := newGuard(p, tokens, clients, map[string]cauth.Authenticatable{"u1": user})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer tok-1")
	guard.SetRequest(req)

	ctx := context.Background()
	u1, _ := guard.User(ctx)
	u2, _ := guard.User(ctx)

	if u1 != u2 {
		t.Error("User() should return the same cached instance on subsequent calls")
	}
}

func TestTokenGuardSetRequestClearsCache(t *testing.T) {
	p := newOAuthServer()
	user := newStubUser("u1")
	tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)
	clients := oauthserver.NewMemoryClientStore()

	tok := newTestToken("tok-1", "u1", "c1", nil)
	_ = tokens.Save(context.Background(), tok)

	guard := newGuard(p, tokens, clients, map[string]cauth.Authenticatable{"u1": user})

	req1 := httptest.NewRequest("GET", "/", nil)
	req1.Header.Set("Authorization", "Bearer tok-1")
	guard.SetRequest(req1)
	u1, _ := guard.User(context.Background())

	if u1 == nil {
		t.Fatal("expected user from first request")
	}

	// New request without a token.
	req2 := httptest.NewRequest("GET", "/", nil)
	guard.SetRequest(req2)
	u2, _ := guard.User(context.Background())

	if u2 != nil {
		t.Error("expected nil user after SetRequest clears the cache")
	}
}

func TestTokenGuardActingAsOverride(t *testing.T) {
	p := newOAuthServer()
	user := newStubUser("u99")

	tok := newTestToken("t-acting", "u99", "c1", []string{"admin"})
	at := oauthserver.NewAccessToken(tok, p)
	p.ActingAs(user, at, []string{"admin"})

	defer p.ClearActing()

	guard := newGuard(p, oauthserver.NewMemoryTokenStore(), oauthserver.NewMemoryClientStore(), nil)

	// No request needed — ActingAs bypasses the request entirely.
	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got == nil {
		t.Fatal("expected user from ActingAs override")
	}

	if got.GetAuthIdentifier() != "u99" {
		t.Errorf("user id = %q, want %q", got.GetAuthIdentifier(), "u99")
	}
}

func TestTokenGuardActingAsClientOverride(t *testing.T) {
	p := newOAuthServer()
	client := &oauthserver.Client{ID: "c99", Name: "TestClient"}
	p.ActingAsClient(client, []string{"read"})

	defer p.ClearActing()

	guard := newGuard(p, oauthserver.NewMemoryTokenStore(), oauthserver.NewMemoryClientStore(), nil)

	got, err := guard.Client(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got == nil {
		t.Fatal("expected client from ActingAsClient override")
	}

	if got.ID != "c99" {
		t.Errorf("client id = %q, want %q", got.ID, "c99")
	}
}

func TestTokenGuardCheckAndGuest(t *testing.T) {
	p := newOAuthServer()
	user := newStubUser("u1")
	tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)
	clients := oauthserver.NewMemoryClientStore()

	tok := newTestToken("tok-1", "u1", "c1", nil)
	_ = tokens.Save(context.Background(), tok)

	guard := newGuard(p, tokens, clients, map[string]cauth.Authenticatable{"u1": user})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer tok-1")
	guard.SetRequest(req)

	if !guard.Check(context.Background()) {
		t.Error("Check should return true for authenticated request")
	}

	if guard.Guest(context.Background()) {
		t.Error("Guest should return false for authenticated request")
	}
}

func TestTokenGuardIDReturnsIdentifier(t *testing.T) {
	p := newOAuthServer()
	user := newStubUser("u5")
	tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)
	clients := oauthserver.NewMemoryClientStore()

	tok := newTestToken("tok-5", "u5", "c1", nil)
	_ = tokens.Save(context.Background(), tok)

	guard := newGuard(p, tokens, clients, map[string]cauth.Authenticatable{"u5": user})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer tok-5")
	guard.SetRequest(req)

	id := guard.ID(context.Background())

	if id != "u5" {
		t.Errorf("ID = %v, want %q", id, "u5")
	}
}
