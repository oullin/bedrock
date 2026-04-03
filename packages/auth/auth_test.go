package auth_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/foundation"
	"github.com/gollin/packages/auth/memory"
)

func TestDefaultPasswordHasher(t *testing.T) {
	t.Parallel()

	hasher := auth.DefaultPasswordHasher{}
	encoded, err := hasher.Hash(context.Background(), "secret-pass")

	if err != nil {
		t.Fatalf("Hash: %v", err)
	}

	if err := hasher.Compare(context.Background(), encoded, "secret-pass"); err != nil {
		t.Fatalf("Compare: %v", err)
	}

	if err := hasher.Compare(context.Background(), encoded, "wrong-pass"); err == nil {
		t.Fatal("expected wrong password to fail")
	}
}

func TestSessionGuardLoginAndAuthenticateRequest(t *testing.T) {
	t.Parallel()

	clock := memory.NewFixedClock(time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC))
	ids := memory.NewSequenceIDGenerator("session")
	users := memory.NewInMemoryUserRepository()
	sessions := memory.NewInMemorySessionStore()
	hasher := auth.DefaultPasswordHasher{}

	passwordHash, err := hasher.Hash(context.Background(), "password-123")

	if err != nil {
		t.Fatalf("Hash: %v", err)
	}

	user := &foundation.User{
		ID:           "user-1",
		Name:         "Test User",
		Email:        "user@example.com",
		PasswordHash: passwordHash,
		CreatedAt:    clock.Now(),
		UpdatedAt:    clock.Now(),
	}

	if err := users.Create(context.Background(), user); err != nil {
		t.Fatalf("Create: %v", err)
	}

	manager, err := auth.NewManager(auth.Config{
		DefaultGuard:                "web",
		DefaultProvider:             "users",
		IdentifierField:             "email",
		SessionLifetime:             24 * time.Hour,
		RememberLifetime:            30 * 24 * time.Hour,
		PasswordConfirmationTimeout: 3 * time.Hour,
		VerificationTTL:             time.Hour,
		SigningKey:                  []byte("test-signing-key"),
		Cookies: auth.CookieConfig{
			SessionName:  "session",
			RememberName: "remember",
			Path:         "/",
			HTTPOnly:     true,
		},
	}, map[string]auth.UserProvider{"users": users}, sessions, auth.ManagerDependencies{
		Hasher: hasher,
		Clock:  clock,
		IDs:    ids,
	})

	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	resolved, err := manager.ValidateCredentials(context.Background(), map[string]string{
		"email":    "user@example.com",
		"password": "password-123",
	})

	if err != nil {
		t.Fatalf("ValidateCredentials: %v", err)
	}

	recorder := httptest.NewRecorder()
	session, rememberToken, err := manager.DefaultGuard().Login(context.Background(), recorder, resolved, true, false)

	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if rememberToken == "" {
		t.Fatal("expected remember token")
	}

	if len(recorder.Result().Cookies()) == 0 {
		t.Fatal("expected cookies to be written")
	}

	request := httptest.NewRequest("GET", "/protected", nil)

	for _, cookie := range recorder.Result().Cookies() {
		request.AddCookie(cookie)
	}

	authenticatedSession, authenticatedUser, err := manager.DefaultGuard().AuthenticateRequest(context.Background(), httptest.NewRecorder(), request)

	if err != nil {
		t.Fatalf("AuthenticateRequest: %v", err)
	}

	if authenticatedSession.ID != session.ID {
		t.Fatalf("unexpected session id: %s", authenticatedSession.ID)
	}

	if authenticatedUser.GetAuthIdentifier() != user.ID {
		t.Fatalf("unexpected user id: %s", authenticatedUser.GetAuthIdentifier())
	}
}
