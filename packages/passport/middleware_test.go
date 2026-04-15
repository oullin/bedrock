package passport_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/passport"
	cauth "github.com/bedrock/packages/contracts/auth"
)

// okHandler is a test HTTP handler that always returns 200 OK.
var okHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
})

func setupMiddlewareGuard(p *passport.Passport, tok *passport.Token, user cauth.Authenticatable) (*passport.TokenGuard, string) {
	tokens := passport.NewMemoryTokenStore().WithPassport(p)
	clients := passport.NewMemoryClientStore()
	provider := &stubProvider{users: map[string]cauth.Authenticatable{user.GetAuthIdentifier(): user}}
	guard := passport.NewTokenGuard(p, tokens, clients, provider)

	_ = tokens.Save(context.Background(), tok)

	return guard, tok.ID
}

func TestCheckTokenPassesWithMatchingScope(t *testing.T) {
	p := newPassport()
	user := newStubUser("u1")
	tok := newTestToken("tok-1", "u1", "c1", []string{"user:read"})
	guard, tokenID := setupMiddlewareGuard(p, tok, user)

	mw := passport.CheckToken(guard, "user:read")
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
	p := newPassport()
	user := newStubUser("u1")
	tok := newTestToken("tok-1", "u1", "c1", []string{"user:read"})
	guard, tokenID := setupMiddlewareGuard(p, tok, user)

	mw := passport.CheckToken(guard, "admin")
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
	p := newPassport()
	guard := newGuard(p, passport.NewMemoryTokenStore(), passport.NewMemoryClientStore(), nil)

	mw := passport.CheckToken(guard, "read")
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
	p := newPassport()
	user := newStubUser("u1")
	tok := newTestToken("tok-1", "u1", "c1", []string{"some:scope"})
	guard, tokenID := setupMiddlewareGuard(p, tok, user)

	// No scope argument — any authenticated request passes.
	mw := passport.CheckToken(guard)
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
	p := newPassport()
	user := newStubUser("u1")
	tok := newTestToken("tok-1", "u1", "c1", []string{"orders:read"})
	guard, tokenID := setupMiddlewareGuard(p, tok, user)

	mw := passport.CheckTokenForAnyScope(guard, "user:read", "orders:read")
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
	p := newPassport()
	user := newStubUser("u1")
	tok := newTestToken("tok-1", "u1", "c1", []string{"user:read"})
	guard, tokenID := setupMiddlewareGuard(p, tok, user)

	mw := passport.CheckTokenForAnyScope(guard, "admin", "superuser")
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
	p := newPassport()
	user := newStubUser("u1")
	tok := newTestToken("tok-1", "u1", "c1", []string{"user:read", "orders:read"})
	guard, tokenID := setupMiddlewareGuard(p, tok, user)

	// Both scopes present — passes.
	mw := passport.CheckToken(guard, "user:read", "orders:read")
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
	p := newPassport()
	user := newStubUser("u1")
	// Token only has user:read, not admin.
	tok := newTestToken("tok-1", "u1", "c1", []string{"user:read"})
	guard, tokenID := setupMiddlewareGuard(p, tok, user)

	mw := passport.CheckToken(guard, "user:read", "admin")
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenID)
	guard.SetRequest(req)
	mw(okHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d (one scope missing)", rr.Code, http.StatusForbidden)
	}
}
