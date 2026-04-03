package middleware_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/foundation"
	"github.com/gollin/packages/auth/memory"
	authmw "github.com/gollin/packages/auth/middleware"
)

func TestRequireAuthenticatedAndVerifiedEmail(t *testing.T) {
	t.Parallel()

	stack, user, cookies := newMiddlewareStack(t, false)

	protected := stack.RequireAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentUser, ok := authmw.CurrentUser(r)

		if !ok || currentUser.GetAuthIdentifier() != user.ID {
			t.Fatalf("unexpected current user: %#v", currentUser)
		}

		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	rec := httptest.NewRecorder()
	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	verified := stack.RequireVerifiedEmail(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	rec = httptest.NewRecorder()
	verified.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("unexpected verified-email status: %d", rec.Code)
	}
}

func TestRequirePasswordConfirmed(t *testing.T) {
	t.Parallel()

	stack, _, cookies := newMiddlewareStack(t, true)

	protected := stack.RequirePasswordConfirmed(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	rec := httptest.NewRecorder()
	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	stack.Clock = memory.NewFixedClock(time.Date(2026, 4, 3, 6, 0, 0, 0, time.UTC))
	protected = stack.RequirePasswordConfirmed(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	rec = httptest.NewRecorder()
	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("unexpected stale confirmation status: %d", rec.Code)
	}
}

func TestWriteJSONError(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	authmw.WriteJSONError(rec, http.StatusTooManyRequests, &auth.ThrottleError{Scope: "login", RetryAfter: 45 * time.Second})

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	var body map[string]any

	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("Decode: %v", err)
	}

	if body["error"] != "throttled" {
		t.Fatalf("unexpected error body: %#v", body)
	}
}

func newMiddlewareStack(t *testing.T, confirmPassword bool) (authmw.Stack, *foundation.User, []*http.Cookie) {
	t.Helper()

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

	if confirmPassword {
		user.MarkEmailAsVerified(now)
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
		SigningKey:                  []byte("test"),
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

	if confirmPassword {
		confirmedAt := clock.Now()
		session.PasswordConfirmedAt = &confirmedAt

		if err := manager.DefaultGuard().UpdateSession(context.Background(), session); err != nil {
			t.Fatalf("UpdateSession: %v", err)
		}
	}

	return authmw.Stack{
		Guard:  manager.DefaultGuard(),
		Config: manager.Config(),
		Clock:  clock,
	}, user, rec.Result().Cookies()
}
