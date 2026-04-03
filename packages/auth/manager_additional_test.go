package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/foundation"
	"github.com/gollin/packages/auth/memory"
)

func TestNewManagerRejectsMissingDefaultProvider(t *testing.T) {
	t.Parallel()

	_, err := auth.NewManager(auth.Config{
		DefaultGuard:    "web",
		DefaultProvider: "users",
		Cookies: auth.CookieConfig{
			SessionName:  "session",
			RememberName: "remember",
			Path:         "/",
		},
	}, map[string]auth.UserProvider{"other": memory.NewInMemoryUserRepository()}, memory.NewInMemorySessionStore(), auth.ManagerDependencies{})

	if err == nil {
		t.Fatal("expected missing default provider error")
	}
}

func TestSessionGuardUsesRememberCookieAndExpiresSessions(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC)
	clock := memory.NewFixedClock(now)
	users := memory.NewInMemoryUserRepository()
	sessions := memory.NewInMemorySessionStore()
	hasher := auth.DefaultPasswordHasher{}

	hash, err := hasher.Hash(context.Background(), "password-123")

	if err != nil {
		t.Fatalf("Hash: %v", err)
	}

	user := &foundation.User{
		ID:           "user-1",
		Name:         "User",
		Email:        "user@example.com",
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := users.Create(context.Background(), user); err != nil {
		t.Fatalf("Create: %v", err)
	}

	manager, err := auth.NewManager(auth.Config{
		DefaultGuard:                "web",
		DefaultProvider:             "users",
		IdentifierField:             "email",
		SessionLifetime:             time.Hour,
		RememberLifetime:            24 * time.Hour,
		PasswordConfirmationTimeout: time.Hour,
		VerificationTTL:             time.Hour,
		SigningKey:                  []byte("signing-key"),
		Cookies: auth.CookieConfig{
			SessionName:  "session",
			RememberName: "remember",
			Path:         "/",
			HTTPOnly:     true,
		},
	}, map[string]auth.UserProvider{"users": users}, sessions, auth.ManagerDependencies{
		Hasher: hasher,
		Clock:  clock,
		IDs:    memory.NewSequenceIDGenerator("session"),
	})

	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	rec := httptest.NewRecorder()
	session, _, err := manager.DefaultGuard().Login(context.Background(), rec, user, false, false)

	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	clock.Advance(2 * time.Hour)
	request := httptest.NewRequest("GET", "/protected", nil)
	request.AddCookie(rec.Result().Cookies()[0])

	if _, _, err := manager.DefaultGuard().AuthenticateRequest(context.Background(), httptest.NewRecorder(), request); err != auth.ErrUnauthorized {
		t.Fatalf("expected expired session to be unauthorized, got %v", err)
	}

	rec = httptest.NewRecorder()

	if err := manager.DefaultGuard().Logout(context.Background(), rec, session, user); err != nil {
		t.Fatalf("Logout: %v", err)
	}

	rec = httptest.NewRecorder()
	_, rememberToken, err := manager.DefaultGuard().Login(context.Background(), rec, user, true, false)

	if err != nil {
		t.Fatalf("Login remember: %v", err)
	}

	if rememberToken == "" {
		t.Fatal("expected remember token")
	}

	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "remember" && strings.Contains(cookie.Value, "|") {
			t.Fatalf("expected remember cookie to be encrypted, got %q", cookie.Value)
		}
	}

	rememberRequest := httptest.NewRequest("GET", "/protected", nil)

	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "remember" {
			rememberRequest.AddCookie(cookie)
		}
	}

	authenticatedSession, authenticatedUser, err := manager.DefaultGuard().AuthenticateRequest(context.Background(), httptest.NewRecorder(), rememberRequest)

	if err != nil {
		t.Fatalf("AuthenticateRequest with remember cookie: %v", err)
	}

	if authenticatedSession == nil || authenticatedUser.GetAuthIdentifier() != user.ID {
		t.Fatal("expected remember cookie authentication to succeed")
	}

	if !manager.DefaultGuard().ViaRemember() {
		t.Fatal("expected viaRemember to be true after remember-cookie authentication")
	}
}

func TestRequestAndTokenGuards(t *testing.T) {
	t.Parallel()

	users := memory.NewInMemoryUserRepository()
	sessions := memory.NewInMemorySessionStore()
	user := &foundation.User{
		ID:        "user-1",
		Name:      "Token User",
		Email:     "token@example.com",
		APIToken:  "plain-token",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := users.Create(context.Background(), user); err != nil {
		t.Fatalf("Create: %v", err)
	}

	manager, err := auth.NewManager(auth.Config{
		DefaultGuard:    "web",
		DefaultProvider: "users",
		SigningKey:      []byte("signing-key"),
		Cookies: auth.CookieConfig{
			SessionName:  "session",
			RememberName: "remember",
			Path:         "/",
		},
	}, map[string]auth.UserProvider{"users": users}, sessions, auth.ManagerDependencies{})

	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	request := httptest.NewRequest("GET", "/profile", nil)
	requestGuard := manager.ViaRequest("request", request, func(ctx context.Context, r *http.Request, provider auth.UserProvider) (auth.Authenticatable, error) {
		return provider.RetrieveByID(ctx, "user-1")
	})

	resolved, err := requestGuard.User(context.Background())

	if err != nil {
		t.Fatalf("RequestGuard.User: %v", err)
	}

	if resolved.GetAuthIdentifier() != user.ID {
		t.Fatalf("unexpected request guard user: %s", resolved.GetAuthIdentifier())
	}

	tokenRequest := httptest.NewRequest("GET", "/api?api_token=plain-token", nil)
	tokenGuard, err := manager.RegisterTokenGuard("api", tokenRequest, "users", "api_token", "api_token", false)

	if err != nil {
		t.Fatalf("RegisterTokenGuard: %v", err)
	}

	resolved, err = tokenGuard.User(context.Background())

	if err != nil {
		t.Fatalf("TokenGuard.User: %v", err)
	}

	if resolved.GetAuthIdentifier() != user.ID {
		t.Fatalf("unexpected token guard user: %s", resolved.GetAuthIdentifier())
	}
}
