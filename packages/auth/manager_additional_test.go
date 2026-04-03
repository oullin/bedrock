package auth_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/foundation"
	"github.com/gollin/packages/auth/memory"
	"github.com/gollin/packages/security/encryption"
)

type failingIDGenerator struct {
	err error
}

type failingSessionStore struct {
	auth.SessionStore
	createErr error
	updateErr error
}

func (g failingIDGenerator) NewID() (string, error) {
	return "", g.err
}

func (s failingSessionStore) Create(ctx context.Context, session *auth.Session) error {
	if s.createErr != nil {
		return s.createErr
	}

	return s.SessionStore.Create(ctx, session)
}

func (s failingSessionStore) Update(ctx context.Context, session *auth.Session) error {
	if s.updateErr != nil {
		return s.updateErr
	}

	return s.SessionStore.Update(ctx, session)
}

func TestNewManagerRejectsMissingDefaultProvider(t *testing.T) {
	t.Parallel()

	users := newInMemoryUserRepository(t)

	_, err := auth.NewManager(auth.Config{
		DefaultGuard:    "web",
		DefaultProvider: "users",
		Cookies: auth.CookieConfig{
			SessionName:  "session",
			RememberName: "remember",
			Path:         "/",
		},
	}, map[string]auth.UserProvider{"other": users}, memory.NewInMemorySessionStore(), auth.ManagerDependencies{})

	if err == nil {
		t.Fatal("expected missing default provider error")
	}
}

func TestSessionGuardUsesRememberCookieAndExpiresSessions(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC)
	clock := memory.NewFixedClock(now)
	users := newInMemoryUserRepository(t)
	sessions := memory.NewInMemorySessionStore()
	hasher := newDefaultPasswordHasher(t)

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

	users := newInMemoryUserRepository(t)
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

func TestSessionGuardLoginReturnsIDGenerationError(t *testing.T) {
	t.Parallel()

	users := newInMemoryUserRepository(t)
	user := &foundation.User{
		ID:        "user-1",
		Name:      "User",
		Email:     "user@example.com",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := users.Create(context.Background(), user); err != nil {
		t.Fatalf("Create: %v", err)
	}

	guard := auth.NewSessionGuard("web", auth.Config{
		SessionLifetime:  time.Hour,
		RememberLifetime: 24 * time.Hour,
		Cookies: auth.CookieConfig{
			SessionName:  "session",
			RememberName: "remember",
			Path:         "/",
		},
	}, users, memory.NewInMemorySessionStore(), newDefaultPasswordHasher(t), mustEncrypter(t), []byte("hash-key"), memory.NewFixedClock(time.Now().UTC()), failingIDGenerator{err: errBoom}, auth.NoopLogger{})

	if _, _, err := guard.Login(context.Background(), httptest.NewRecorder(), user, false, false); err != errBoom {
		t.Fatalf("expected id error, got %v", err)
	}
}

func TestSessionGuardLoginRestoresRememberTokenWhenSessionCreateFails(t *testing.T) {
	t.Parallel()

	users := newInMemoryUserRepository(t)
	user := &foundation.User{
		ID:            "user-1",
		Name:          "User",
		Email:         "user@example.com",
		RememberToken: "original-token",
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if err := users.Create(context.Background(), user); err != nil {
		t.Fatalf("Create: %v", err)
	}

	guard := auth.NewSessionGuard("web", auth.Config{
		SessionLifetime:  time.Hour,
		RememberLifetime: 24 * time.Hour,
		Cookies: auth.CookieConfig{
			SessionName:  "session",
			RememberName: "remember",
			Path:         "/",
		},
	}, users, failingSessionStore{SessionStore: memory.NewInMemorySessionStore(), createErr: errBoom}, newDefaultPasswordHasher(t), mustEncrypter(t), []byte("hash-key"), memory.NewFixedClock(time.Now().UTC()), memory.NewSequenceIDGenerator("session"), auth.NoopLogger{})

	rec := httptest.NewRecorder()

	if _, _, err := guard.Login(context.Background(), rec, user, true, false); err == nil {
		t.Fatal("expected login error")
	}

	if len(rec.Result().Cookies()) != 0 {
		t.Fatalf("expected no cookies on failed login, got %d", len(rec.Result().Cookies()))
	}

	stored, err := users.RetrieveByID(context.Background(), user.ID)

	if err != nil {
		t.Fatalf("RetrieveByID: %v", err)
	}

	if stored.GetRememberToken() != "original-token" {
		t.Fatalf("expected remember token to be restored, got %q", stored.GetRememberToken())
	}
}

func TestSessionGuardLoginRestoresRememberTokenWhenRememberCookieEncryptionFails(t *testing.T) {
	t.Parallel()

	users := newInMemoryUserRepository(t)
	sessions := memory.NewInMemorySessionStore()
	user := &foundation.User{
		ID:            "user-1",
		Name:          "User",
		Email:         "user@example.com",
		RememberToken: "original-token",
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	if err := users.Create(context.Background(), user); err != nil {
		t.Fatalf("Create: %v", err)
	}

	guard := auth.NewSessionGuard("web", auth.Config{
		SessionLifetime:  time.Hour,
		RememberLifetime: 24 * time.Hour,
		Cookies: auth.CookieConfig{
			SessionName:  "session",
			RememberName: "remember",
			Path:         "/",
		},
	}, users, sessions, newDefaultPasswordHasher(t), &encryption.Encrypter{}, []byte("hash-key"), memory.NewFixedClock(time.Now().UTC()), memory.NewSequenceIDGenerator("session"), auth.NoopLogger{})

	rec := httptest.NewRecorder()

	if _, _, err := guard.Login(context.Background(), rec, user, true, false); err == nil {
		t.Fatal("expected login error")
	}

	if len(rec.Result().Cookies()) != 0 {
		t.Fatalf("expected no cookies on failed login, got %d", len(rec.Result().Cookies()))
	}

	if _, err := sessions.FindByID(context.Background(), "session-1"); err != auth.ErrUnauthorized {
		t.Fatalf("expected no persisted session, got %v", err)
	}

	stored, err := users.RetrieveByID(context.Background(), user.ID)

	if err != nil {
		t.Fatalf("RetrieveByID: %v", err)
	}

	if stored.GetRememberToken() != "original-token" {
		t.Fatalf("expected remember token to be restored, got %q", stored.GetRememberToken())
	}
}

func TestSessionGuardCompleteTwoFactorRestoresStateWhenRememberUpdateFails(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	clock := memory.NewFixedClock(now)
	users := newInMemoryUserRepository(t)
	sessionStore := memory.NewInMemorySessionStore()
	user := &foundation.User{
		ID:            "user-1",
		Name:          "User",
		Email:         "user@example.com",
		RememberToken: "original-token",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := users.Create(context.Background(), user); err != nil {
		t.Fatalf("Create: %v", err)
	}

	session := &auth.Session{
		ID:               "session-1",
		UserID:           user.ID,
		PendingTwoFactor: true,
		PendingRemember:  true,
		LastSeenAt:       now,
		CreatedAt:        now,
		ExpiresAt:        now.Add(time.Hour),
	}

	if err := sessionStore.Create(context.Background(), session); err != nil {
		t.Fatalf("Create session: %v", err)
	}

	guard := auth.NewSessionGuard("web", auth.Config{
		SessionLifetime:  time.Hour,
		RememberLifetime: 24 * time.Hour,
		Cookies: auth.CookieConfig{
			SessionName:  "session",
			RememberName: "remember",
			Path:         "/",
		},
	}, users, failingSessionStore{SessionStore: sessionStore, updateErr: errBoom}, newDefaultPasswordHasher(t), mustEncrypter(t), []byte("hash-key"), clock, memory.NewSequenceIDGenerator("session"), auth.NoopLogger{})

	rec := httptest.NewRecorder()

	if _, err := guard.CompleteTwoFactor(context.Background(), rec, session, user); err == nil {
		t.Fatal("expected complete-two-factor error")
	}

	if len(rec.Result().Cookies()) != 0 {
		t.Fatalf("expected no cookies on failed completion, got %d", len(rec.Result().Cookies()))
	}

	if !session.PendingTwoFactor || !session.PendingRemember || session.AuthenticatedAt != nil {
		t.Fatalf("expected session state restored, got %#v", session)
	}

	stored, err := users.RetrieveByID(context.Background(), user.ID)

	if err != nil {
		t.Fatalf("RetrieveByID: %v", err)
	}

	if stored.GetRememberToken() != "original-token" {
		t.Fatalf("expected remember token to be restored, got %q", stored.GetRememberToken())
	}
}

var errBoom = errors.New("boom")

func mustEncrypter(t *testing.T) *encryption.Encrypter {
	t.Helper()

	encrypter, err := encryption.New(encryption.Config{
		Key:    []byte("0123456789abcdef0123456789abcdef"),
		Cipher: encryption.AES256CBC,
	})

	if err != nil {
		t.Fatalf("encryption.New: %v", err)
	}

	return encrypter
}
