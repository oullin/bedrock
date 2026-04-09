package jetstream

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedrock/packages/fortify"
)

// --- test token repo ---

type testTokenRepo struct {
	tokens  map[string]*PersonalAccessToken
	byHash  map[string]*PersonalAccessToken
	created bool
	deleted bool
}

func newTestTokenRepo(tokens ...*PersonalAccessToken) *testTokenRepo {
	r := &testTokenRepo{
		tokens: make(map[string]*PersonalAccessToken),
		byHash: make(map[string]*PersonalAccessToken),
	}
	for _, t := range tokens {
		r.tokens[t.ID] = t
		r.byHash[t.TokenHash] = t
	}
	return r
}

func (r *testTokenRepo) Create(_ context.Context, t *PersonalAccessToken) error {
	r.created = true
	if t.ID == "" {
		t.ID = "generated-id"
	}
	r.tokens[t.ID] = t
	r.byHash[t.TokenHash] = t
	return nil
}
func (r *testTokenRepo) FindByID(_ context.Context, id string) (*PersonalAccessToken, error) {
	return r.tokens[id], nil
}
func (r *testTokenRepo) FindByTokenHash(_ context.Context, hash string) (*PersonalAccessToken, error) {
	return r.byHash[hash], nil
}
func (r *testTokenRepo) FindByUser(_ context.Context, _ string) ([]PersonalAccessToken, error) {
	return nil, nil
}
func (r *testTokenRepo) Update(_ context.Context, t *PersonalAccessToken) error {
	r.tokens[t.ID] = t
	return nil
}
func (r *testTokenRepo) Delete(_ context.Context, id string) error {
	r.deleted = true
	delete(r.tokens, id)
	return nil
}

// --- token type tests ---

func TestPersonalAccessTokenHasPermission(t *testing.T) {
	token := &PersonalAccessToken{Permissions: []string{"read", "write"}}

	if !token.HasPermission("read") {
		t.Fatal("expected read permission")
	}

	if token.HasPermission("delete") {
		t.Fatal("should not have delete permission")
	}
}

func TestPersonalAccessTokenWildcardPermission(t *testing.T) {
	token := &PersonalAccessToken{Permissions: []string{"*"}}

	if !token.HasPermission("anything") {
		t.Fatal("wildcard should match any permission")
	}
}

func TestPersonalAccessTokenIsExpired(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	future := time.Now().Add(time.Hour)

	expired := &PersonalAccessToken{ExpiresAt: &past}
	if !expired.IsExpired() {
		t.Fatal("token with past expiry should be expired")
	}

	valid := &PersonalAccessToken{ExpiresAt: &future}
	if valid.IsExpired() {
		t.Fatal("token with future expiry should not be expired")
	}

	noExpiry := &PersonalAccessToken{ExpiresAt: nil}
	if noExpiry.IsExpired() {
		t.Fatal("token with no expiry should not be expired")
	}
}

func TestGeneratePlainToken(t *testing.T) {
	plain, hash, err := GeneratePlainToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(plain) != 80 {
		t.Fatalf("expected 80 char plain token, got %d", len(plain))
	}

	if hash != HashToken(plain) {
		t.Fatal("hash should match HashToken(plain)")
	}

	plain2, _, _ := GeneratePlainToken()
	if plain == plain2 {
		t.Fatal("tokens should be unique")
	}
}

// --- create token handler tests ---

func TestCreateTokenSuccess(t *testing.T) {
	user := &testTeamUser{id: "1"}
	tokenRepo := newTestTokenRepo()
	js, events := buildTestJetstream(user, newTestTeamRepo(), newTestInvitationRepo())

	handler := NewCreateTokenHandler(js, tokenRepo)
	w := httptest.NewRecorder()
	r := postForm("/user/api-tokens", "name=My+Token&permissions=read,write")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var result NewTokenResult
	_ = json.NewDecoder(w.Body).Decode(&result)

	if result.PlainText == "" {
		t.Fatal("expected plain text token in response")
	}

	if result.Token.Name != "My Token" {
		t.Fatalf("expected token name 'My Token', got %s", result.Token.Name)
	}

	if len(result.Token.Permissions) != 2 {
		t.Fatalf("expected 2 permissions, got %d", len(result.Token.Permissions))
	}

	if !tokenRepo.created {
		t.Fatal("expected token to be created in repo")
	}

	if len(events.dispatched) != 1 || events.dispatched[0].Name != EventTokenCreated {
		t.Fatal("expected TokenCreated event")
	}
}

func TestCreateTokenMissingName(t *testing.T) {
	user := &testTeamUser{id: "1"}
	js, _ := buildTestJetstream(user, newTestTeamRepo(), newTestInvitationRepo())

	handler := NewCreateTokenHandler(js, newTestTokenRepo())
	w := httptest.NewRecorder()
	r := postForm("/user/api-tokens", "permissions=read")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestCreateTokenDefaultPermissions(t *testing.T) {
	user := &testTeamUser{id: "1"}
	tokenRepo := newTestTokenRepo()
	js, _ := buildTestJetstream(user, newTestTeamRepo(), newTestInvitationRepo())

	handler := NewCreateTokenHandler(js, tokenRepo)
	w := httptest.NewRecorder()
	r := postForm("/user/api-tokens", "name=Default+Token")

	handler.ServeHTTP(w, r)

	var result NewTokenResult
	_ = json.NewDecoder(w.Body).Decode(&result)

	if len(result.Token.Permissions) != 1 || result.Token.Permissions[0] != "*" {
		t.Fatalf("expected wildcard permission by default, got %v", result.Token.Permissions)
	}
}

// --- update token handler tests ---

func TestUpdateTokenSuccess(t *testing.T) {
	token := &PersonalAccessToken{ID: "tok-1", UserID: "1", Permissions: []string{"read"}}
	user := &testTeamUser{id: "1"}
	tokenRepo := newTestTokenRepo(token)
	js, events := buildTestJetstream(user, newTestTeamRepo(), newTestInvitationRepo())

	mux := http.NewServeMux()
	mux.Handle("PUT /user/api-tokens/{token}", NewUpdateTokenHandler(js, tokenRepo))

	w := httptest.NewRecorder()
	r := putForm("/user/api-tokens/tok-1", "permissions=read,write,delete")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	updated := tokenRepo.tokens["tok-1"]
	if len(updated.Permissions) != 3 {
		t.Fatalf("expected 3 permissions, got %d", len(updated.Permissions))
	}

	if len(events.dispatched) != 1 || events.dispatched[0].Name != EventTokenUpdated {
		t.Fatal("expected TokenUpdated event")
	}
}

func TestUpdateTokenForbiddenForOtherUser(t *testing.T) {
	token := &PersonalAccessToken{ID: "tok-1", UserID: "other"}
	user := &testTeamUser{id: "1"}
	tokenRepo := newTestTokenRepo(token)
	js, _ := buildTestJetstream(user, newTestTeamRepo(), newTestInvitationRepo())

	mux := http.NewServeMux()
	mux.Handle("PUT /user/api-tokens/{token}", NewUpdateTokenHandler(js, tokenRepo))

	w := httptest.NewRecorder()
	r := putForm("/user/api-tokens/tok-1", "permissions=*")
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

// --- delete token handler tests ---

func TestDeleteTokenSuccess(t *testing.T) {
	token := &PersonalAccessToken{ID: "tok-1", UserID: "1"}
	user := &testTeamUser{id: "1"}
	tokenRepo := newTestTokenRepo(token)
	js, events := buildTestJetstream(user, newTestTeamRepo(), newTestInvitationRepo())

	mux := http.NewServeMux()
	mux.Handle("DELETE /user/api-tokens/{token}", NewDeleteTokenHandler(js, tokenRepo))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/user/api-tokens/tok-1", nil)
	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if !tokenRepo.deleted {
		t.Fatal("expected token to be deleted")
	}

	if len(events.dispatched) != 1 || events.dispatched[0].Name != EventTokenDeleted {
		t.Fatal("expected TokenDeleted event")
	}
}

// --- token guard tests ---

func TestTokenGuardAuthenticate(t *testing.T) {
	plain, hash, _ := GeneratePlainToken()
	token := &PersonalAccessToken{ID: "tok-1", UserID: "1", TokenHash: hash}
	user := &testTeamUser{id: "1"}
	tokenRepo := newTestTokenRepo(token)
	provider := &stubUserProvider{user: user}

	guard := NewTokenGuard(tokenRepo, provider)

	r := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	r.Header.Set("Authorization", "Bearer "+plain)

	authUser, authToken, err := guard.Authenticate(r.Context(), r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if authUser.GetAuthIdentifier() != "1" {
		t.Fatal("expected user id 1")
	}

	if authToken.ID != "tok-1" {
		t.Fatal("expected token tok-1")
	}
}

func TestTokenGuardRejectsExpired(t *testing.T) {
	plain, hash, _ := GeneratePlainToken()
	past := time.Now().Add(-time.Hour)
	token := &PersonalAccessToken{ID: "tok-1", UserID: "1", TokenHash: hash, ExpiresAt: &past}
	tokenRepo := newTestTokenRepo(token)

	guard := NewTokenGuard(tokenRepo, &stubUserProvider{})

	r := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	r.Header.Set("Authorization", "Bearer "+plain)

	_, _, err := guard.Authenticate(r.Context(), r)
	if err != ErrUnauthenticated {
		t.Fatalf("expected unauthenticated error, got %v", err)
	}
}

func TestTokenGuardRejectsMissingBearer(t *testing.T) {
	guard := NewTokenGuard(newTestTokenRepo(), &stubUserProvider{})

	r := httptest.NewRequest(http.MethodGet, "/api/test", nil)

	_, _, err := guard.Authenticate(r.Context(), r)
	if err != ErrUnauthenticated {
		t.Fatalf("expected unauthenticated error, got %v", err)
	}
}

// --- TokenCan middleware tests ---

func TestTokenCanAllows(t *testing.T) {
	called := false
	handler := TokenCan("read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	token := &PersonalAccessToken{Permissions: []string{"read", "write"}}
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r = r.WithContext(WithToken(r.Context(), token))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if !called {
		t.Fatal("expected handler to be called")
	}
}

func TestTokenCanDenies(t *testing.T) {
	called := false
	handler := TokenCan("delete")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	token := &PersonalAccessToken{Permissions: []string{"read"}}
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r = r.WithContext(WithToken(r.Context(), token))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if called {
		t.Fatal("handler should not be called without permission")
	}

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestTokenCanDeniesNoToken(t *testing.T) {
	handler := TokenCan("read")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not be called")
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

// --- stub user provider ---

type stubUserProvider struct {
	user *testTeamUser
}

func (p *stubUserProvider) RetrieveByID(_ context.Context, _ string) (fortify.Authenticatable, error) {
	if p.user == nil {
		return nil, nil
	}
	return p.user, nil
}
func (p *stubUserProvider) RetrieveByToken(_ context.Context, _ string, _ string) (fortify.Authenticatable, error) {
	return nil, nil
}
func (p *stubUserProvider) RetrieveByCredentials(_ context.Context, _ map[string]string) (fortify.Authenticatable, error) {
	return nil, nil
}
func (p *stubUserProvider) UpdateRememberToken(_ context.Context, _ fortify.Authenticatable, _ string) error {
	return nil
}
func (p *stubUserProvider) ValidateCredentials(_ context.Context, _ fortify.Authenticatable, _ map[string]string) (bool, error) {
	return false, nil
}
func (p *stubUserProvider) RehashPasswordIfRequired(_ context.Context, _ fortify.Authenticatable, _ map[string]string, _ bool) error {
	return nil
}
