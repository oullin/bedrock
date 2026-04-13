package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/bedrock/packages/auth"
	cauth "github.com/bedrock/packages/contracts/auth"
)

func TestTokenGuardUserCanBeRetrievedByQueryStringVariable(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "api_token": "foo"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	guard := auth.NewTokenGuard("api", provider)

	req := httptest.NewRequest(http.MethodGet, "/?api_token=foo", nil)
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got == nil || got.GetAuthIdentifier() != "1" {
		t.Error("expected user with id \"1\"")
	}

	if !guard.Check(context.Background()) {
		t.Error("Check should return true")
	}

	if guard.Guest(context.Background()) {
		t.Error("Guest should return false")
	}

	if guard.ID(context.Background()) != "1" {
		t.Errorf("ID = %v, want \"1\"", guard.ID(context.Background()))
	}
}

func TestTokenGuardUserCanBeRetrievedByBearerToken(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "api_token": "foo"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	guard := auth.NewTokenGuard("api", provider)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer foo")
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got == nil || got.GetAuthIdentifier() != "1" {
		t.Error("expected user with id \"1\" from Bearer token")
	}
}

func TestTokenGuardUserCanBeRetrievedByAuthHeaders(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "api_token": "foo"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	guard := auth.NewTokenGuard("api", provider)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetBasicAuth("anything", "foo")
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got == nil || got.GetAuthIdentifier() != "1" {
		t.Error("expected user with id \"1\" from Basic Auth password")
	}
}

func TestTokenGuardUserCanBeRetrievedByFormField(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "api_token": "foo"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	guard := auth.NewTokenGuard("api", provider)

	form := url.Values{"api_token": {"foo"}}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got == nil || got.GetAuthIdentifier() != "1" {
		t.Error("expected user with id \"1\" from form field")
	}
}

func TestTokenGuardReturnsNilWhenTokenIsEmpty(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	guard := auth.NewTokenGuard("api", provider)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Error("expected nil user when no token")
	}
}

func TestTokenGuardReturnsNilWhenUserNotFound(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	guard := auth.NewTokenGuard("api", provider)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Error("expected nil user for invalid token")
	}
}

func TestTokenGuardCachesUserPerRequest(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "api_token": "foo"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	guard := auth.NewTokenGuard("api", provider)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer foo")
	guard.SetRequest(req)

	ctx := context.Background()
	u1, _ := guard.User(ctx)
	u2, _ := guard.User(ctx)

	if u1 != u2 {
		t.Error("User() should return cached user on subsequent calls")
	}
}

func TestTokenGuardSetRequestClearsCache(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "api_token": "foo"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	guard := auth.NewTokenGuard("api", provider)

	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.Header.Set("Authorization", "Bearer foo")
	guard.SetRequest(req1)

	u1, _ := guard.User(context.Background())

	if u1 == nil {
		t.Fatal("expected user from first request")
	}

	req2 := httptest.NewRequest(http.MethodGet, "/", nil) // No token.
	guard.SetRequest(req2)

	u2, _ := guard.User(context.Background())

	if u2 != nil {
		t.Error("expected nil user after SetRequest clears cache")
	}
}

func TestTokenGuardCustomKeys(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "custom_token": "bar"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	guard := auth.NewTokenGuard("api", provider)
	guard.SetInputKey("custom_token")
	guard.SetStorageKey("custom_token")

	req := httptest.NewRequest(http.MethodGet, "/?custom_token=bar", nil)
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got == nil || got.GetAuthIdentifier() != "1" {
		t.Error("expected user with custom token key")
	}
}

func TestTokenGuardCustomKeyBearerToken(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "custom_token": "baz"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	guard := auth.NewTokenGuard("api", provider)
	guard.SetInputKey("custom_token")
	guard.SetStorageKey("custom_token")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer baz")
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got == nil || got.GetAuthIdentifier() != "1" {
		t.Error("expected user with custom key and Bearer token")
	}
}

func TestTokenGuardSetUser(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "api_token": "foo"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	guard := auth.NewTokenGuard("api", provider)

	if guard.HasUser() {
		t.Error("HasUser should be false initially")
	}

	guard.SetUser(user)

	if !guard.HasUser() {
		t.Error("HasUser should be true after SetUser")
	}

	got, _ := guard.User(context.Background())

	if got != user {
		t.Error("User() should return the user set via SetUser")
	}
}

func TestTokenGuardForgetUser(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "api_token": "foo"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	guard := auth.NewTokenGuard("api", provider)

	guard.SetUser(user)
	guard.ForgetUser()

	if guard.HasUser() {
		t.Error("HasUser should be false after ForgetUser")
	}
}

func TestTokenGuardCheckAndGuestForUnauthenticatedUser(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	guard := auth.NewTokenGuard("api", provider)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req)

	if guard.Check(context.Background()) {
		t.Error("Check should be false with no token")
	}

	if !guard.Guest(context.Background()) {
		t.Error("Guest should be true with no token")
	}
}

func TestTokenGuardIDReturnsNilWhenNoUser(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	guard := auth.NewTokenGuard("api", provider)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req)

	if guard.ID(context.Background()) != nil {
		t.Error("ID should return nil when no user")
	}
}

// --- TokenGuard: Validate ---

func TestTokenGuardValidateReturnsTrueForValidToken(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "api_token": "valid-token"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	guard := auth.NewTokenGuard("api", provider)

	ok := guard.Validate(context.Background(), map[string]string{"api_token": "valid-token"})

	if !ok {
		t.Error("Validate should return true for valid token")
	}
}

func TestTokenGuardValidateReturnsFalseForInvalidToken(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	guard := auth.NewTokenGuard("api", provider)

	ok := guard.Validate(context.Background(), map[string]string{"api_token": "invalid"})

	if ok {
		t.Error("Validate should return false for invalid token")
	}
}

// --- TokenGuard: GetTokenForRequest ---

func TestTokenGuardGetTokenForRequestFromQuery(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	guard := auth.NewTokenGuard("api", provider)

	req := httptest.NewRequest(http.MethodGet, "/?api_token=my-token", nil)
	guard.SetRequest(req)

	if guard.GetTokenForRequest() != "my-token" {
		t.Errorf("GetTokenForRequest() = %q, want %q", guard.GetTokenForRequest(), "my-token")
	}
}

func TestTokenGuardGetTokenForRequestFromBearer(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	guard := auth.NewTokenGuard("api", provider)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer bearer-token")
	guard.SetRequest(req)

	if guard.GetTokenForRequest() != "bearer-token" {
		t.Errorf("GetTokenForRequest() = %q, want %q", guard.GetTokenForRequest(), "bearer-token")
	}
}

func TestTokenGuardGetTokenForRequestReturnsEmptyWhenNoToken(t *testing.T) {
	provider := &stubProvider{users: map[string]cauth.Authenticatable{}}
	guard := auth.NewTokenGuard("api", provider)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req)

	if guard.GetTokenForRequest() != "" {
		t.Error("GetTokenForRequest should return empty when no token")
	}
}

// --- TokenGuard: SHA256 hashing ---

func TestTokenGuardHashesTokenWithSHA256(t *testing.T) {
	// SHA256 of "plain-token"
	hashedToken := "23fb79e20d37abf2418d78115eb0cc8c74b52f4ed8b91dda7fc03a1d41fc15e3"
	user := auth.NewGenericUser(map[string]any{"id": "1", "api_token": hashedToken})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}
	guard := auth.NewTokenGuard("api", provider)
	guard.SetHash(true)

	req := httptest.NewRequest(http.MethodGet, "/?api_token=plain-token", nil)
	guard.SetRequest(req)

	u, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if u == nil {
		t.Error("expected user when token hashes match")
	}
}
