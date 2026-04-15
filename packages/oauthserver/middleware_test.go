package passport_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/oauthserver"
	cauth "github.com/bedrock/packages/contracts/auth"
)

// okHandler is a test HTTP handler that always returns 200 OK.
var okHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
})

func setupMiddlewareGuard(p *oauthserver.OAuthServer, tok *oauthserver.Token, user cauth.Authenticatable) (*oauthserver.TokenGuard, string) {
	tokens := oauthserver.NewMemoryTokenStore().WithOAuthServer(p)
	clients := oauthserver.NewMemoryClientStore()
	provider := &stubProvider{users: map[string]cauth.Authenticatable{user.GetAuthIdentifier(): user}}
	guard := oauthserver.NewTokenGuard(p, tokens, clients, provider)

	_ = tokens.Save(context.Background(), tok)

	return guard, tok.ID
}

func TestCheckTokenPassesWithMatchingScope(t *testing.T) {
	p := newOAuthServer()
	user := newStubUser("u1")
	tok := newTestToken("tok-1", "u1", "c1", []string{"user:read"})
	guard, tokenID := setupMiddlewareGuard(p, tok, user)

	mw := oauthserver.CheckToken(guard, "user:read")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenID)
	guard.SetRequest(req)
	mw(okHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestCheckTokenFailsWithMissingScope(t *testing.T) {
	p := newOAuthServer()
	user := newStubUser("u1")
	tok := newTestToken("tok-1", "u1", "c1", []string{"user:read"})
	guard, tokenID := setupMiddlewareGuard(p, tok, user)

	mw := oauthserver.CheckToken(guard, "admin")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenID)
	guard.SetRequest(req)
	mw(okHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusForbidden)
	}

	var body map[string]string

	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}

	if body["message"] != "Invalid scope(s) provided." {
		t.Errorf("message = %q, want %q", body["message"], "Invalid scope(s) provided.")
	}
}

func TestCheckTokenReturns401WhenUnauthenticated(t *testing.T) {
	p := newOAuthServer()
	guard := newGuard(p, oauthserver.NewMemoryTokenStore(), oauthserver.NewMemoryClientStore(), nil)

	mw := oauthserver.CheckToken(guard, "read")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil) // No Authorization header.
	guard.SetRequest(req)
	mw(okHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}

	var body map[string]string

	_ = json.Unmarshal(rr.Body.Bytes(), &body)

	if body["message"] != "Unauthenticated." {
		t.Errorf("message = %q, want %q", body["message"], "Unauthenticated.")
	}
}

func TestCheckTokenPassesWithNoScopesRequired(t *testing.T) {
	p := newOAuthServer()
	user := newStubUser("u1")
	tok := newTestToken("tok-1", "u1", "c1", []string{"some:scope"})
	guard, tokenID := setupMiddlewareGuard(p, tok, user)

	// No scope argument — any authenticated request passes.
	mw := oauthserver.CheckToken(guard)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenID)
	guard.SetRequest(req)
	mw(okHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestCheckTokenForAnyScopePassesIfAnyMatch(t *testing.T) {
	p := newOAuthServer()
	user := newStubUser("u1")
	tok := newTestToken("tok-1", "u1", "c1", []string{"orders:read"})
	guard, tokenID := setupMiddlewareGuard(p, tok, user)

	mw := oauthserver.CheckTokenForAnyScope(guard, "user:read", "orders:read")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenID)
	guard.SetRequest(req)
	mw(okHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestCheckTokenForAnyScopeFailsIfNoneMatch(t *testing.T) {
	p := newOAuthServer()
	user := newStubUser("u1")
	tok := newTestToken("tok-1", "u1", "c1", []string{"user:read"})
	guard, tokenID := setupMiddlewareGuard(p, tok, user)

	mw := oauthserver.CheckTokenForAnyScope(guard, "admin", "superuser")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenID)
	guard.SetRequest(req)
	mw(okHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusForbidden)
	}
}

func TestCheckTokenMultipleScopesAllRequired(t *testing.T) {
	p := newOAuthServer()
	user := newStubUser("u1")
	tok := newTestToken("tok-1", "u1", "c1", []string{"user:read", "orders:read"})
	guard, tokenID := setupMiddlewareGuard(p, tok, user)

	// Both scopes present — passes.
	mw := oauthserver.CheckToken(guard, "user:read", "orders:read")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenID)
	guard.SetRequest(req)
	mw(okHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d (both scopes present)", rr.Code, http.StatusOK)
	}
}

func TestCheckTokenMultipleScopesOneMissing(t *testing.T) {
	p := newOAuthServer()
	user := newStubUser("u1")
	// Token only has user:read, not admin.
	tok := newTestToken("tok-1", "u1", "c1", []string{"user:read"})
	guard, tokenID := setupMiddlewareGuard(p, tok, user)

	mw := oauthserver.CheckToken(guard, "user:read", "admin")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenID)
	guard.SetRequest(req)
	mw(okHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d (one scope missing)", rr.Code, http.StatusForbidden)
	}
}
