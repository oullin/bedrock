package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/bedrock/packages/anvil/encryption"
	securitycrypto "github.com/bedrock/packages/anvil/support/crypto"
)

// ---------- test helpers ----------

type testUser struct {
	id            string
	email         string
	passwordHash  string
	rememberToken string
	verifiedAt    *time.Time
}

func (u *testUser) GetAuthIdentifierName() string { return "id" }
func (u *testUser) GetAuthIdentifier() string     { return u.id }
func (u *testUser) GetAuthPasswordName() string   { return "password" }
func (u *testUser) GetAuthPassword() string       { return u.passwordHash }
func (u *testUser) SetAuthPassword(password string) {
	u.passwordHash = password
}

func (u *testUser) GetRememberToken() string { return u.rememberToken }
func (u *testUser) SetRememberToken(token string) {
	u.rememberToken = token
}

func (u *testUser) GetRememberTokenName() string { return "remember_token" }
func (u *testUser) HasVerifiedEmail() bool       { return u.verifiedAt != nil }
func (u *testUser) MarkEmailAsVerified(at time.Time) {
	u.verifiedAt = &at
}

func (u *testUser) MarkEmailAsUnverified() {
	u.verifiedAt = nil
}

func (u *testUser) GetEmailForVerification() string  { return u.email }
func (u *testUser) GetEmailForPasswordReset() string { return u.email }

type fakeProvider struct {
	user          *testUser
	hasher        PasswordHasher
	rehashCalled  bool
	rehashEnabled bool
}

func (p *fakeProvider) RetrieveByID(_ context.Context, id string) (Authenticatable, error) {
	if p.user != nil && p.user.id == id {
		return p.user, nil
	}

	return nil, ErrUserNotFound
}

func (p *fakeProvider) RetrieveByToken(_ context.Context, id string, token string) (Authenticatable, error) {
	if p.user != nil && p.user.id == id && p.user.rememberToken == token {
		return p.user, nil
	}

	return nil, ErrUnauthorized
}

func (p *fakeProvider) RetrieveByCredentials(_ context.Context, credentials map[string]string) (Authenticatable, error) {
	if p.user == nil {
		return nil, ErrUserNotFound
	}

	if email, ok := credentials["email"]; ok && p.user.email != email {
		return nil, ErrUserNotFound
	}

	if token, ok := credentials["api_token"]; ok && token == "token-123" {
		return p.user, nil
	}

	if token, ok := credentials["hashed_api_token"]; ok && token == securitycrypto.HashString("token-123") {
		return p.user, nil
	}

	if len(credentials) == 1 {
		if _, ok := credentials["email"]; ok {
			return p.user, nil
		}
	}

	return p.user, nil
}

func (p *fakeProvider) UpdateRememberToken(_ context.Context, user Authenticatable, token string) error {
	user.SetRememberToken(token)
	return nil
}

func (p *fakeProvider) ValidateCredentials(ctx context.Context, user Authenticatable, credentials map[string]string) (bool, error) {
	password, ok := credentials["password"]
	if !ok {
		return false, nil
	}

	if err := p.hasher.Compare(ctx, user.GetAuthPassword(), password); err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (p *fakeProvider) RehashPasswordIfRequired(_ context.Context, _ Authenticatable, _ map[string]string, _ bool) error {
	if p.rehashEnabled {
		p.rehashCalled = true
	}
	return nil
}

type fakeSessionStore struct {
	sessions map[string]*Session
}

func (s *fakeSessionStore) Create(_ context.Context, session *Session) error {
	if s.sessions == nil {
		s.sessions = map[string]*Session{}
	}

	clone := *session
	s.sessions[session.ID] = &clone
	return nil
}

func (s *fakeSessionStore) FindByID(_ context.Context, id string) (*Session, error) {
	session, ok := s.sessions[id]
	if !ok {
		return nil, ErrUnauthorized
	}

	clone := *session
	return &clone, nil
}

func (s *fakeSessionStore) Update(_ context.Context, session *Session) error {
	if _, ok := s.sessions[session.ID]; !ok {
		return ErrUnauthorized
	}

	clone := *session
	s.sessions[session.ID] = &clone
	return nil
}

func (s *fakeSessionStore) Delete(_ context.Context, id string) error {
	delete(s.sessions, id)
	return nil
}

type fakeCookieManager struct {
	values  map[string]string
	deleted []string
}

func (m *fakeCookieManager) Read(_ *http.Request, name string) (string, error) {
	value, ok := m.values[name]
	if !ok || value == "" {
		return "", http.ErrNoCookie
	}

	return value, nil
}

func (m *fakeCookieManager) Write(_ http.ResponseWriter, cookie Cookie) error {
	if m.values == nil {
		m.values = map[string]string{}
	}

	m.values[cookie.Name] = cookie.Value
	return nil
}

func (m *fakeCookieManager) Delete(_ http.ResponseWriter, cookie Cookie) error {
	if m.values == nil {
		m.values = map[string]string{}
	}

	delete(m.values, cookie.Name)
	m.deleted = append(m.deleted, cookie.Name)
	return nil
}

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time { return c.now }

type fixedIDs struct {
	value string
}

func (g fixedIDs) NewID() (string, error) { return g.value, nil }

func defaultCfg() Config {
	return Config{
		DefaultGuard:     "web",
		DefaultProvider:  "users",
		SessionLifetime:  time.Hour,
		RememberLifetime: 24 * time.Hour,
		SigningKey:       []byte("test-signing-key"),
		Cookies: CookieConfig{
			SessionName:  "session",
			RememberName: "remember",
		},
	}
}

// ======================== SESSION GUARD TESTS ========================

// Upstream: testLoginStoresIdentifierInSession
func TestSessionGuardLoginStoresIdentifierInSession(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, clock, fixedIDs{value: "session-1"})
	session, _, err := guard.Login(context.Background(), httptest.NewRecorder(), user, false, false)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if session.UserID != "user-1" {
		t.Fatalf("expected session UserID 'user-1', got %q", session.UserID)
	}

	if session.ID != "session-1" {
		t.Fatalf("expected session ID 'session-1', got %q", session.ID)
	}

	// Session should be stored
	stored, err := sessions.FindByID(context.Background(), "session-1")
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if stored.UserID != "user-1" {
		t.Fatalf("stored session has wrong UserID: %q", stored.UserID)
	}
}

// Upstream: testLoginFiresLoginAndAuthenticatedEvents
// GAP: Bedrock does not implement event dispatching. Events would need an EventDispatcher interface.
// The session guard creates sessions and cookies but does not fire events.

// Upstream: testFailedAttemptFiresFailedEvent
// GAP: Same as above — no event dispatching in Bedrock.

// Upstream: testAuthenticateReturnsUserWhenUserIsNotNull
func TestAuthenticateRequestReturnsUserWhenSessionExists(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	sessions := &fakeSessionStore{}
	cookies := &fakeCookieManager{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, clock, fixedIDs{value: "session-1"})

	// Login first to create session
	_, _, err := guard.Login(context.Background(), httptest.NewRecorder(), user, false, false)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	// Now authenticate the request
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	_, resolved, err := guard.AuthenticateRequest(context.Background(), httptest.NewRecorder(), request)
	if err != nil {
		t.Fatalf("AuthenticateRequest: %v", err)
	}

	if resolved.GetAuthIdentifier() != "user-1" {
		t.Fatalf("expected user-1, got %q", resolved.GetAuthIdentifier())
	}
}

// Upstream: testAuthenticateThrowsWhenUserIsNull
func TestAuthenticateRequestReturnsErrorWhenNoSession(t *testing.T) {
	t.Parallel()

	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}
	clock := fixedClock{now: time.Now().UTC()}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})
	_, _, err := guard.AuthenticateRequest(context.Background(), httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

// Upstream: testNullIsReturnedForUserIfNoUserFound
func TestAuthenticateRequestReturnsErrorWhenUserNotFound(t *testing.T) {
	t.Parallel()

	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	sessions := &fakeSessionStore{}
	cookies := &fakeCookieManager{}

	// Create a session for a user that doesn't exist in the provider
	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: nil}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})

	// Manually insert a session
	sessions.Create(context.Background(), &Session{
		ID:        "s1",
		UserID:    "nonexistent",
		ExpiresAt: clock.Now().Add(time.Hour),
	})
	cookies.Write(nil, Cookie{Name: "session", Value: "s1"})

	_, _, err := guard.AuthenticateRequest(context.Background(), httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if err == nil {
		t.Fatal("expected error when user not found")
	}
}

// Upstream: testLogoutRemovesSessionTokenAndRememberMeCookie
func TestLogoutRemovesSessionAndCookies(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	sessions := &fakeSessionStore{}
	cookies := &fakeCookieManager{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})

	session, _, err := guard.Login(context.Background(), httptest.NewRecorder(), user, false, false)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	err = guard.Logout(context.Background(), httptest.NewRecorder(), session, user)
	if err != nil {
		t.Fatalf("Logout: %v", err)
	}

	// Session should be deleted
	_, err = sessions.FindByID(context.Background(), "s1")
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatal("expected session to be deleted")
	}

	// Cookies should be cleared
	if _, err := cookies.Read(nil, "session"); err == nil {
		t.Fatal("expected session cookie to be cleared")
	}
	if _, err := cookies.Read(nil, "remember"); err == nil {
		t.Fatal("expected remember cookie to be cleared")
	}
}

// Upstream: testLogoutDoesNotEnqueueRememberMeCookieForDeletionIfCookieDoesntExist
func TestLogoutClearsRememberCookieEvenIfAbsent(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	clock := fixedClock{now: time.Now().UTC()}
	sessions := &fakeSessionStore{}
	cookies := &fakeCookieManager{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})
	session, _, _ := guard.Login(context.Background(), httptest.NewRecorder(), user, false, false)

	// Logout without remember cookie set — should not error
	err := guard.Logout(context.Background(), httptest.NewRecorder(), session, user)
	if err != nil {
		t.Fatalf("Logout: %v", err)
	}
}

// Upstream: testLogoutDoesNotSetRememberTokenIfNotPreviouslySet
func TestLogoutClearsRememberToken(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com", rememberToken: "existing-token"}
	clock := fixedClock{now: time.Now().UTC()}
	sessions := &fakeSessionStore{}
	cookies := &fakeCookieManager{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})
	session, _, _ := guard.Login(context.Background(), httptest.NewRecorder(), user, false, false)

	guard.Logout(context.Background(), httptest.NewRecorder(), session, user)

	if user.GetRememberToken() != "" {
		t.Fatal("expected remember token to be cleared on logout")
	}
}

// Upstream: testLoginMethodQueuesCookieWhenRemembering
func TestLoginMethodCreatesRememberCookieWhenRemembering(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}
	cfg := defaultCfg()

	manager, err := NewManager(cfg, map[string]UserProvider{"users": &fakeProvider{user: user, hasher: hasher}}, ManagerDependencies{
		Hasher:   hasher,
		Sessions: sessions,
		Cookies:  cookies,
		Clock:    clock,
		IDs:      fixedIDs{value: "s1"},
	})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	guard := manager.DefaultGuard()
	_, rememberToken, err := guard.Login(context.Background(), httptest.NewRecorder(), user, true, false)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if rememberToken == "" {
		t.Fatal("expected remember token to be generated")
	}

	// Remember cookie should be set
	val, err := cookies.Read(nil, "remember")
	if err != nil {
		t.Fatal("expected remember cookie to be set")
	}
	if val == "" {
		t.Fatal("expected non-empty remember cookie value")
	}
}

// Upstream: testLoginMethodCreatesRememberTokenIfOneDoesntExist
func TestLoginCreatesRememberTokenIfMissing(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com", rememberToken: ""}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	encrypter, _ := encryption.New(encryption.Config{
		Key:    deriveCipherKey([]byte("test-signing-key")),
		Cipher: encryption.AES256CBC,
	})

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, encrypter, clock, fixedIDs{value: "s1"})
	_, token, err := guard.Login(context.Background(), httptest.NewRecorder(), user, true, false)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if token == "" {
		t.Fatal("expected remember token to be created")
	}
	if user.GetRememberToken() == "" {
		t.Fatal("expected user remember token to be set")
	}
}

// Upstream: testLoginStoresIdentifierInSession + remember=false
func TestLoginWithoutRememberDoesNotCreateRememberCookie(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})
	_, rememberToken, err := guard.Login(context.Background(), httptest.NewRecorder(), user, false, false)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if rememberToken != "" {
		t.Fatal("expected no remember token when remember=false")
	}

	if _, err := cookies.Read(nil, "remember"); err == nil {
		t.Fatal("expected no remember cookie when remember=false")
	}
}

// Upstream: testUserUsesRememberCookieIfItExists
func TestSessionGuardLoginAndRememberRestore(t *testing.T) {
	t.Parallel()

	hasher, err := NewDefaultPasswordHasher()
	if err != nil {
		t.Fatalf("NewDefaultPasswordHasher: %v", err)
	}

	passwordHash, err := hasher.Hash(context.Background(), "secret")
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}

	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}
	cfg := defaultCfg()

	manager, err := NewManager(cfg, map[string]UserProvider{"users": &fakeProvider{user: user, hasher: hasher}}, ManagerDependencies{
		Hasher:   hasher,
		Sessions: sessions,
		Cookies:  cookies,
		Clock:    clock,
		IDs:      fixedIDs{value: "session-1"},
	})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	guard := manager.DefaultGuard()
	if guard == nil {
		t.Fatal("expected default guard")
	}

	recorder := httptest.NewRecorder()
	session, rememberToken, err := guard.Login(context.Background(), recorder, user, true, false)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if rememberToken == "" {
		t.Fatal("expected remember token")
	}

	if session.ID != "session-1" {
		t.Fatalf("unexpected session id: %q", session.ID)
	}

	// Remove session cookie to simulate "new browser" — only remember cookie remains
	delete(cookies.values, cfg.Cookies.SessionName)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	restoredSession, restoredUser, err := guard.AuthenticateRequest(context.Background(), httptest.NewRecorder(), request)
	if err != nil {
		t.Fatalf("AuthenticateRequest: %v", err)
	}

	if restoredSession.UserID != user.id {
		t.Fatalf("unexpected restored session user id: %q", restoredSession.UserID)
	}

	if restoredUser.GetAuthIdentifier() != user.id {
		t.Fatalf("unexpected restored user: %q", restoredUser.GetAuthIdentifier())
	}
}

// Upstream: testSessionGuardClearsInvalidRememberCookie
func TestSessionGuardClearsInvalidRememberCookie(t *testing.T) {
	t.Parallel()

	cookies := &fakeCookieManager{values: map[string]string{"remember": "bad-value"}}
	encrypter, err := encryption.New(encryption.Config{
		Key:    deriveCipherKey([]byte("test-signing-key")),
		Cipher: encryption.AES256CBC,
	})
	if err != nil {
		t.Fatalf("encryption.New: %v", err)
	}

	guard := NewSessionGuard("web", Config{Cookies: CookieConfig{SessionName: "session", RememberName: "remember"}}, &fakeProvider{}, &fakeSessionStore{}, cookies, nil, encrypter, fixedClock{now: time.Now().UTC()}, fixedIDs{value: "session-1"})

	_, _, err = guard.AuthenticateRequest(context.Background(), httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected unauthorized, got %v", err)
	}

	if len(cookies.deleted) == 0 || cookies.deleted[0] != "remember" {
		t.Fatalf("expected remember cookie deletion, got %v", cookies.deleted)
	}
}

// Upstream: testLoginWithPendingTwoFactor
func TestLoginWithPendingTwoFactor(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1"}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})
	session, _, err := guard.Login(context.Background(), httptest.NewRecorder(), user, false, true)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if !session.PendingTwoFactor {
		t.Fatal("expected PendingTwoFactor to be true")
	}
	if session.AuthenticatedAt != nil {
		t.Fatal("expected AuthenticatedAt to be nil when pending 2FA")
	}
}

// Test session expiry detection
func TestSessionGuardExpiresOldSessions(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1"}
	now := time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)
	sessions := &fakeSessionStore{}
	cookies := &fakeCookieManager{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, fixedClock{now: now}, fixedIDs{value: "s1"})
	guard.Login(context.Background(), httptest.NewRecorder(), user, false, false)

	// Advance clock past session lifetime
	expiredClock := fixedClock{now: now.Add(2 * time.Hour)}
	guard2 := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, expiredClock, fixedIDs{value: "s2"})

	_, _, err := guard2.AuthenticateRequest(context.Background(), httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected expired session to be unauthorized, got %v", err)
	}
}

// Test session LastSeenAt update
func TestSessionGuardUpdatesLastSeenAt(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1"}
	now := time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)
	sessions := &fakeSessionStore{}
	cookies := &fakeCookieManager{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, fixedClock{now: now}, fixedIDs{value: "s1"})
	guard.Login(context.Background(), httptest.NewRecorder(), user, false, false)

	// Authenticate 30 minutes later
	later := now.Add(30 * time.Minute)
	guard2 := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, fixedClock{now: later}, fixedIDs{value: "s2"})
	session, _, err := guard2.AuthenticateRequest(context.Background(), httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if err != nil {
		t.Fatalf("AuthenticateRequest: %v", err)
	}

	if !session.LastSeenAt.Equal(later) {
		t.Fatalf("expected LastSeenAt to be updated to %v, got %v", later, session.LastSeenAt)
	}
}

// Guard name
func TestSessionGuardName(t *testing.T) {
	t.Parallel()

	guard := NewSessionGuard("api", Config{}, nil, nil, &fakeCookieManager{}, nil, nil, fixedClock{}, fixedIDs{})
	if guard.Name() != "api" {
		t.Fatalf("expected name 'api', got %q", guard.Name())
	}
}

// Logout with nil session and nil user
func TestLogoutWithNilSessionAndUser(t *testing.T) {
	t.Parallel()

	cookies := &fakeCookieManager{}
	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{}, &fakeSessionStore{}, cookies, nil, nil, fixedClock{now: time.Now().UTC()}, fixedIDs{})

	err := guard.Logout(context.Background(), httptest.NewRecorder(), nil, nil)
	if err != nil {
		t.Fatalf("Logout with nil session/user: %v", err)
	}
}

// ======================== TOKEN GUARD TESTS ========================

// Upstream: testUserCanBeRetrievedByQueryStringVariable
func TestTokenGuardUserFromQueryString(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	provider := &fakeProvider{user: user}

	request := httptest.NewRequest(http.MethodGet, "/?api_token=token-123", nil)
	guard := NewTokenGuard(provider, request, "api_token", "api_token", false)

	resolved, err := guard.User(context.Background())
	if err != nil {
		t.Fatalf("User: %v", err)
	}
	if resolved.GetAuthIdentifier() != "user-1" {
		t.Fatalf("expected user-1, got %q", resolved.GetAuthIdentifier())
	}
}

// Upstream: testUserCanBeRetrievedByBearerToken
func TestTokenGuardUserFromBearerToken(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	provider := &fakeProvider{user: user}

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer token-123")
	guard := NewTokenGuard(provider, request, "api_token", "api_token", false)

	resolved, err := guard.User(context.Background())
	if err != nil {
		t.Fatalf("User: %v", err)
	}
	if resolved.GetAuthIdentifier() != "user-1" {
		t.Fatalf("expected user-1, got %q", resolved.GetAuthIdentifier())
	}
}

// Upstream: testUserCanBeRetrievedByAuthHeaders
func TestTokenGuardUserFromAuthorizationHeader(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	provider := &fakeProvider{user: user}

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "token-123")
	guard := NewTokenGuard(provider, request, "api_token", "api_token", false)

	resolved, err := guard.User(context.Background())
	if err != nil {
		t.Fatalf("User: %v", err)
	}
	if resolved.GetAuthIdentifier() != "user-1" {
		t.Fatalf("expected user-1, got %q", resolved.GetAuthIdentifier())
	}
}

// Upstream: testTokenCanBeHashed
func TestTokenGuardHashedStorageKey(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	provider := &fakeProvider{user: user}

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.URL = &url.URL{RawQuery: "token=token-123"}

	guard := NewTokenGuard(provider, request, "token", "hashed_api_token", true)
	resolved, err := guard.User(context.Background())
	if err != nil {
		t.Fatalf("User: %v", err)
	}

	if resolved.GetAuthIdentifier() != user.id {
		t.Fatalf("unexpected user: %q", resolved.GetAuthIdentifier())
	}
}

// Upstream: testValidateCanDetermineIfCredentialsAreValid / Invalid
func TestTokenGuardReturnsErrorForMissingToken(t *testing.T) {
	t.Parallel()

	provider := &fakeProvider{user: &testUser{id: "u1"}}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	guard := NewTokenGuard(provider, request, "api_token", "api_token", false)

	_, err := guard.User(context.Background())
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized for missing token, got %v", err)
	}
}

// Upstream: testValidateIfApiTokenIsEmpty
func TestTokenGuardReturnsErrorForEmptyToken(t *testing.T) {
	t.Parallel()

	provider := &fakeProvider{user: &testUser{id: "u1"}}
	request := httptest.NewRequest(http.MethodGet, "/?api_token=", nil)
	guard := NewTokenGuard(provider, request, "api_token", "api_token", false)

	_, err := guard.User(context.Background())
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized for empty token, got %v", err)
	}
}

// Upstream: testItAllowToPassCustomRequestInSetterAndUseItForValidation
func TestTokenGuardSetRequest(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	provider := &fakeProvider{user: user}

	// First request has no token
	guard := NewTokenGuard(provider, httptest.NewRequest(http.MethodGet, "/", nil), "api_token", "api_token", false)
	_, err := guard.User(context.Background())
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}

	// Set new request with token
	guard.SetRequest(httptest.NewRequest(http.MethodGet, "/?api_token=token-123", nil))
	resolved, err := guard.User(context.Background())
	if err != nil {
		t.Fatalf("User after SetRequest: %v", err)
	}
	if resolved.GetAuthIdentifier() != "user-1" {
		t.Fatalf("expected user-1, got %q", resolved.GetAuthIdentifier())
	}
}

// Upstream: testUserCanBeRetrievedByBearerTokenWithCustomKey
func TestTokenGuardCustomInputKey(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	provider := &fakeProvider{user: user}

	request := httptest.NewRequest(http.MethodGet, "/?custom_key=token-123", nil)
	guard := NewTokenGuard(provider, request, "custom_key", "api_token", false)

	resolved, err := guard.User(context.Background())
	if err != nil {
		t.Fatalf("User: %v", err)
	}
	if resolved.GetAuthIdentifier() != "user-1" {
		t.Fatalf("expected user-1, got %q", resolved.GetAuthIdentifier())
	}
}

// Token guard caches user on repeated calls
func TestTokenGuardCachesUser(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	provider := &fakeProvider{user: user}

	request := httptest.NewRequest(http.MethodGet, "/?api_token=token-123", nil)
	guard := NewTokenGuard(provider, request, "api_token", "api_token", false)

	first, _ := guard.User(context.Background())
	second, _ := guard.User(context.Background())

	if first != second {
		t.Fatal("expected cached user to be the same instance")
	}
}

// Token guard with nil request
func TestTokenGuardNilRequest(t *testing.T) {
	t.Parallel()

	guard := NewTokenGuard(&fakeProvider{}, nil, "", "", false)
	_, err := guard.User(context.Background())
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized for nil request, got %v", err)
	}
}

// Default input key
func TestTokenGuardDefaultInputKey(t *testing.T) {
	t.Parallel()

	guard := NewTokenGuard(&fakeProvider{}, nil, "", "", false)
	if guard.inputKey != "api_token" {
		t.Fatalf("expected default inputKey 'api_token', got %q", guard.inputKey)
	}
	if guard.storageKey != "api_token" {
		t.Fatalf("expected default storageKey 'api_token', got %q", guard.storageKey)
	}
}

// Upstream: testTokenGuardValidate
func TestTokenGuardValidateSuccess(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}

	guard := NewTokenGuard(&fakeProvider{user: user, hasher: hasher}, nil, "", "", false)

	ok, err := guard.Validate(context.Background(), map[string]string{
		"email":    "user@example.com",
		"password": "secret",
	})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if !ok {
		t.Fatal("expected Validate to return true for valid credentials")
	}
}

func TestTokenGuardValidateFailure(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}

	guard := NewTokenGuard(&fakeProvider{user: user, hasher: hasher}, nil, "", "", false)

	ok, err := guard.Validate(context.Background(), map[string]string{
		"email":    "user@example.com",
		"password": "wrong",
	})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if ok {
		t.Fatal("expected Validate to return false for invalid credentials")
	}
}

func TestTokenGuardValidateUserNotFound(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	guard := NewTokenGuard(&fakeProvider{user: nil, hasher: hasher}, nil, "", "", false)

	ok, err := guard.Validate(context.Background(), map[string]string{
		"email":    "nobody@example.com",
		"password": "secret",
	})
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if ok {
		t.Fatal("expected Validate to return false when user not found")
	}
}

// ======================== REQUEST GUARD TESTS ========================

func TestRequestGuardResolvesUser(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	provider := &fakeProvider{user: user}

	guard := NewRequestGuard(func(_ context.Context, r *http.Request, provider UserProvider) (Authenticatable, error) {
		if r.URL.Path != "/me" {
			return nil, ErrUnauthorized
		}
		return provider.RetrieveByID(context.Background(), user.id)
	}, httptest.NewRequest(http.MethodGet, "/me", nil), provider)

	resolved, err := guard.User(context.Background())
	if err != nil {
		t.Fatalf("User: %v", err)
	}
	if resolved.GetAuthIdentifier() != "user-1" {
		t.Fatalf("expected user-1, got %q", resolved.GetAuthIdentifier())
	}
}

func TestRequestGuardReturnsErrorForInvalidPath(t *testing.T) {
	t.Parallel()

	guard := NewRequestGuard(func(_ context.Context, r *http.Request, _ UserProvider) (Authenticatable, error) {
		if r.URL.Path != "/valid" {
			return nil, ErrUnauthorized
		}
		return nil, nil
	}, httptest.NewRequest(http.MethodGet, "/invalid", nil), &fakeProvider{})

	_, err := guard.User(context.Background())
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestRequestGuardCachesUser(t *testing.T) {
	t.Parallel()

	callCount := 0
	user := &testUser{id: "u1"}
	guard := NewRequestGuard(func(_ context.Context, _ *http.Request, _ UserProvider) (Authenticatable, error) {
		callCount++
		return user, nil
	}, httptest.NewRequest(http.MethodGet, "/", nil), &fakeProvider{})

	guard.User(context.Background())
	guard.User(context.Background())

	if callCount != 1 {
		t.Fatalf("expected callback to be called once, got %d", callCount)
	}
}

func TestRequestGuardSetRequestResetsCache(t *testing.T) {
	t.Parallel()

	callCount := 0
	user := &testUser{id: "u1"}
	guard := NewRequestGuard(func(_ context.Context, _ *http.Request, _ UserProvider) (Authenticatable, error) {
		callCount++
		return user, nil
	}, httptest.NewRequest(http.MethodGet, "/", nil), &fakeProvider{})

	guard.User(context.Background())
	guard.SetRequest(httptest.NewRequest(http.MethodGet, "/new", nil))
	guard.User(context.Background())

	if callCount != 2 {
		t.Fatalf("expected callback to be called twice after SetRequest, got %d", callCount)
	}
}

func TestRequestGuardNilCallback(t *testing.T) {
	t.Parallel()

	guard := NewRequestGuard(nil, httptest.NewRequest(http.MethodGet, "/", nil), &fakeProvider{})
	_, err := guard.User(context.Background())
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized for nil callback, got %v", err)
	}
}

func TestRequestGuardNilRequest(t *testing.T) {
	t.Parallel()

	guard := NewRequestGuard(func(_ context.Context, _ *http.Request, _ UserProvider) (Authenticatable, error) {
		return nil, nil
	}, nil, &fakeProvider{})
	_, err := guard.User(context.Background())
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized for nil request, got %v", err)
	}
}

// ======================== MANAGER TESTS ========================

// Upstream: testAttemptCallsRetrieveByCredentials
func TestManagerValidateCredentials(t *testing.T) {
	t.Parallel()

	hasher, err := NewDefaultPasswordHasher()
	if err != nil {
		t.Fatalf("NewDefaultPasswordHasher: %v", err)
	}

	passwordHash, err := hasher.Hash(context.Background(), "secret")
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}

	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	manager, err := NewManager(Config{DefaultProvider: "users"}, map[string]UserProvider{
		"users": &fakeProvider{user: user, hasher: hasher},
	}, ManagerDependencies{})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	resolved, err := manager.ValidateCredentials(context.Background(), map[string]string{
		"email":    "user@example.com",
		"password": "secret",
	})
	if err != nil {
		t.Fatalf("ValidateCredentials: %v", err)
	}

	if resolved.GetAuthIdentifier() != user.id {
		t.Fatalf("unexpected validated user: %q", resolved.GetAuthIdentifier())
	}
}

// Upstream: testAttemptReturnsFalseIfUserNotGiven
func TestManagerValidateCredentialsReturnsErrorForInvalidPassword(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}

	manager, _ := NewManager(Config{DefaultProvider: "users"}, map[string]UserProvider{
		"users": &fakeProvider{user: user, hasher: hasher},
	}, ManagerDependencies{})

	_, err := manager.ValidateCredentials(context.Background(), map[string]string{
		"email":    "user@example.com",
		"password": "wrong-password",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestManagerValidateCredentialsReturnsErrorForUnknownUser(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	manager, _ := NewManager(Config{DefaultProvider: "users"}, map[string]UserProvider{
		"users": &fakeProvider{user: nil, hasher: hasher},
	}, ManagerDependencies{})

	_, err := manager.ValidateCredentials(context.Background(), map[string]string{
		"email":    "missing@example.com",
		"password": "secret",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestManagerRequiresAtLeastOneProvider(t *testing.T) {
	t.Parallel()

	_, err := NewManager(Config{DefaultProvider: "users"}, map[string]UserProvider{}, ManagerDependencies{})
	if err == nil {
		t.Fatal("expected error for empty providers")
	}
}

func TestManagerRequiresRegisteredDefaultProvider(t *testing.T) {
	t.Parallel()

	_, err := NewManager(Config{DefaultProvider: "nonexistent"}, map[string]UserProvider{
		"users": &fakeProvider{},
	}, ManagerDependencies{})
	if err == nil {
		t.Fatal("expected error for unregistered default provider")
	}
}

func TestManagerProvider(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	fp := &fakeProvider{hasher: hasher}
	manager, _ := NewManager(Config{DefaultProvider: "users"}, map[string]UserProvider{
		"users": fp,
	}, ManagerDependencies{})

	p, ok := manager.Provider("users")
	if !ok || p != fp {
		t.Fatal("expected to get registered provider")
	}

	_, ok = manager.Provider("missing")
	if ok {
		t.Fatal("expected missing provider to return false")
	}
}

func TestManagerDefaultGuardIsNilWithoutSessionDeps(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	manager, _ := NewManager(Config{DefaultGuard: "web", DefaultProvider: "users"}, map[string]UserProvider{
		"users": &fakeProvider{hasher: hasher},
	}, ManagerDependencies{})

	if manager.DefaultGuard() != nil {
		t.Fatal("expected nil default guard when sessions/cookies not provided")
	}
}

func TestManagerRegisterSessionGuardRequiresDeps(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	manager, _ := NewManager(Config{DefaultProvider: "users"}, map[string]UserProvider{
		"users": &fakeProvider{hasher: hasher},
	}, ManagerDependencies{})

	_, err := manager.RegisterSessionGuard("web", "users")
	if err == nil {
		t.Fatal("expected error when sessions missing")
	}
}

func TestManagerRegisterTokenGuard(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	manager, _ := NewManager(Config{DefaultProvider: "users"}, map[string]UserProvider{
		"users": &fakeProvider{hasher: hasher},
	}, ManagerDependencies{})

	guard, err := manager.RegisterTokenGuard("api", httptest.NewRequest(http.MethodGet, "/", nil), "", "", "", false)
	if err != nil {
		t.Fatalf("RegisterTokenGuard: %v", err)
	}
	if guard == nil {
		t.Fatal("expected token guard")
	}
}

func TestManagerRegisterTokenGuardUnknownProvider(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	manager, _ := NewManager(Config{DefaultProvider: "users"}, map[string]UserProvider{
		"users": &fakeProvider{hasher: hasher},
	}, ManagerDependencies{})

	_, err := manager.RegisterTokenGuard("api", nil, "nonexistent", "", "", false)
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestManagerViaRequest(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "u1"}
	hasher, _ := NewDefaultPasswordHasher()
	manager, _ := NewManager(Config{DefaultProvider: "users"}, map[string]UserProvider{
		"users": &fakeProvider{user: user, hasher: hasher},
	}, ManagerDependencies{})

	guard := manager.ViaRequest("custom", httptest.NewRequest(http.MethodGet, "/", nil), func(_ context.Context, _ *http.Request, p UserProvider) (Authenticatable, error) {
		return p.RetrieveByID(context.Background(), "u1")
	})

	resolved, err := guard.User(context.Background())
	if err != nil {
		t.Fatalf("User: %v", err)
	}
	if resolved.GetAuthIdentifier() != "u1" {
		t.Fatalf("expected u1, got %q", resolved.GetAuthIdentifier())
	}
}

// ======================== ERROR TESTS ========================

func TestAuthenticationExceptionError(t *testing.T) {
	t.Parallel()

	exc := AuthenticationException{Message: "custom error", Guards: []string{"web"}}
	expected := "custom error (guards=web)"
	if exc.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, exc.Error())
	}

	exc2 := AuthenticationException{}
	if exc2.Error() != "auth: authentication failed" {
		t.Fatalf("expected 'auth: authentication failed', got %q", exc2.Error())
	}

	exc3 := AuthenticationException{Message: "denied", Guards: []string{"web", "api"}}
	expected3 := "denied (guards=web,api)"
	if exc3.Error() != expected3 {
		t.Fatalf("expected %q, got %q", expected3, exc3.Error())
	}

	exc4 := AuthenticationException{Message: "no guards"}
	if exc4.Error() != "no guards" {
		t.Fatalf("expected 'no guards', got %q", exc4.Error())
	}
}

func TestSentinelErrors(t *testing.T) {
	t.Parallel()

	errors := []error{ErrUnauthorized, ErrUserNotFound, ErrUserExists, ErrInvalidCredentials, ErrInvalidToken, ErrTokenExpired}
	for _, e := range errors {
		if e == nil {
			t.Fatal("sentinel error should not be nil")
		}
		if e.Error() == "" {
			t.Fatal("sentinel error should have non-empty message")
		}
	}
}

// ======================== DEFAULTS TESTS ========================

func TestSystemClockReturnsCurrentTime(t *testing.T) {
	t.Parallel()

	clock := SystemClock{}
	now := clock.Now()
	if time.Since(now) > time.Second {
		t.Fatal("SystemClock.Now() should return current time")
	}
}

func TestRandomIDGeneratorProducesUniqueIDs(t *testing.T) {
	t.Parallel()

	gen := RandomIDGenerator{}
	id1, err := gen.NewID()
	if err != nil {
		t.Fatalf("NewID: %v", err)
	}
	id2, _ := gen.NewID()

	if id1 == "" || id2 == "" {
		t.Fatal("expected non-empty IDs")
	}
	if id1 == id2 {
		t.Fatal("expected unique IDs")
	}
}

// ======================== RECALLER TESTS ========================

func TestRecallerParsing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     string
		wantValid bool
		wantID    string
		wantToken string
		wantHash  string
	}{
		{"valid", "user-1|token|hash", true, "user-1", "token", "hash"},
		{"empty", "", false, "", "", ""},
		{"one segment", "single", false, "single", "", ""},
		{"two segments", "a|b", false, "a", "b", ""},
		{"empty id", "|token|hash", false, "", "token", "hash"},
		{"empty token", "id||hash", false, "id", "", "hash"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRecaller(tc.value)
			if r.Valid() != tc.wantValid {
				t.Fatalf("Valid() = %v, want %v", r.Valid(), tc.wantValid)
			}
			if r.ID() != tc.wantID {
				t.Fatalf("ID() = %q, want %q", r.ID(), tc.wantID)
			}
			if r.Token() != tc.wantToken {
				t.Fatalf("Token() = %q, want %q", r.Token(), tc.wantToken)
			}
			if r.Hash() != tc.wantHash {
				t.Fatalf("Hash() = %q, want %q", r.Hash(), tc.wantHash)
			}
		})
	}
}

// ======================== CONFIG TESTS ========================

func TestExpiredHelper(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 5, 12, 0, 0, 0, time.UTC)
	past := now.Add(-2 * time.Hour)
	recent := now.Add(-30 * time.Minute)

	if !Expired(past, time.Hour, now) {
		t.Fatal("expected past time to be expired")
	}
	if Expired(recent, time.Hour, now) {
		t.Fatal("expected recent time to not be expired")
	}
}

// ======================== STATEFUL GUARD INTERFACE ========================

func TestSessionGuardImplementsStatefulGuard(t *testing.T) {
	t.Parallel()

	var _ StatefulGuard = (*SessionGuard)(nil)
}

// ======================== ATTEMPT / LOGIN USING ID / ONCE TESTS ========================

// Upstream: testAttemptCallsRetrieveByCredentials + testAttemptReturnsTrue
func TestSessionGuardAttemptSuccess(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user, hasher: hasher}, sessions, cookies, hasher, nil, clock, fixedIDs{value: "s1"})

	ok, err := guard.Attempt(context.Background(), httptest.NewRecorder(), map[string]string{
		"email":    "user@example.com",
		"password": "secret",
	}, false)
	if err != nil {
		t.Fatalf("Attempt: %v", err)
	}
	if !ok {
		t.Fatal("expected Attempt to return true for valid credentials")
	}

	// A session should have been created.
	if len(sessions.sessions) == 0 {
		t.Fatal("expected a session to be created after successful attempt")
	}
}

// Upstream: testAttemptReturnsFalse
func TestSessionGuardAttemptFailure(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user, hasher: hasher}, sessions, cookies, hasher, nil, clock, fixedIDs{value: "s1"})

	ok, err := guard.Attempt(context.Background(), httptest.NewRecorder(), map[string]string{
		"email":    "user@example.com",
		"password": "wrong-password",
	}, false)
	if err != nil {
		t.Fatalf("Attempt: %v", err)
	}
	if ok {
		t.Fatal("expected Attempt to return false for invalid credentials")
	}

	if len(sessions.sessions) != 0 {
		t.Fatal("expected no session to be created after failed attempt")
	}
}

// Upstream: testAttemptReturnsFalseWhenUserNotFound
func TestSessionGuardAttemptUserNotFound(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: nil, hasher: hasher}, sessions, cookies, hasher, nil, clock, fixedIDs{value: "s1"})

	ok, err := guard.Attempt(context.Background(), httptest.NewRecorder(), map[string]string{
		"email":    "nobody@example.com",
		"password": "secret",
	}, false)
	if err != nil {
		t.Fatalf("Attempt: %v", err)
	}
	if ok {
		t.Fatal("expected Attempt to return false when user not found")
	}
}

// Upstream: testLoginUsingId
func TestSessionGuardLoginUsingId(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-42", email: "user@example.com"}
	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})

	session, authedUser, err := guard.LoginUsingId(context.Background(), httptest.NewRecorder(), "user-42", false)
	if err != nil {
		t.Fatalf("LoginUsingId: %v", err)
	}
	if session == nil {
		t.Fatal("expected session to be created")
	}
	if authedUser.GetAuthIdentifier() != "user-42" {
		t.Fatalf("expected user ID 'user-42', got %q", authedUser.GetAuthIdentifier())
	}
}

// Upstream: testLoginUsingIdFailsWhenUserNotFound
func TestSessionGuardLoginUsingIdNotFound(t *testing.T) {
	t.Parallel()

	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: nil}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})

	_, _, err := guard.LoginUsingId(context.Background(), httptest.NewRecorder(), "nonexistent", false)
	if err == nil {
		t.Fatal("expected error when user not found")
	}
}

// Upstream: testOnceUsingId
func TestSessionGuardOnceUsingId(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-42", email: "user@example.com"}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})

	authedUser, err := guard.OnceUsingId(context.Background(), "user-42")
	if err != nil {
		t.Fatalf("OnceUsingId: %v", err)
	}
	if authedUser.GetAuthIdentifier() != "user-42" {
		t.Fatalf("expected user ID 'user-42', got %q", authedUser.GetAuthIdentifier())
	}

	// No session should be created for stateless auth.
	if len(sessions.sessions) != 0 {
		t.Fatal("expected no session for stateless OnceUsingId")
	}
}

// Upstream: testOnce
func TestSessionGuardOnce(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user, hasher: hasher}, sessions, cookies, hasher, nil, clock, fixedIDs{value: "s1"})

	authedUser, err := guard.Once(context.Background(), map[string]string{
		"email":    "user@example.com",
		"password": "secret",
	})
	if err != nil {
		t.Fatalf("Once: %v", err)
	}
	if authedUser.GetAuthIdentifier() != "user-1" {
		t.Fatalf("expected user ID 'user-1', got %q", authedUser.GetAuthIdentifier())
	}

	// No session should be created for stateless auth.
	if len(sessions.sessions) != 0 {
		t.Fatal("expected no session for stateless Once")
	}
}

// Upstream: testOnceFailsWithInvalidCredentials
func TestSessionGuardOnceFailsWithInvalidCredentials(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user, hasher: hasher}, sessions, cookies, hasher, nil, clock, fixedIDs{value: "s1"})

	_, err := guard.Once(context.Background(), map[string]string{
		"email":    "user@example.com",
		"password": "wrong",
	})
	if err == nil {
		t.Fatal("expected Once to fail with invalid credentials")
	}
}

// ======================== BASIC AUTH TESTS ========================

// Upstream: testBasicReturnsNullOnValidAttempt
func TestBasicReturnsNilOnValidAttempt(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user, hasher: hasher}, sessions, cookies, hasher, nil, clock, fixedIDs{value: "s1"})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic "+basicAuth("user@example.com", "secret"))
	rec := httptest.NewRecorder()

	err := guard.Basic(context.Background(), rec, req, "email", nil)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

// Upstream: testBasicReturnsNullWhenAlreadyLoggedIn
func TestBasicReturnsNilWhenAlreadyLoggedIn(t *testing.T) {
	t.Parallel()

	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}
	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})

	user := &testUser{id: "user-1", email: "user@example.com"}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(WithUser(req.Context(), user))
	rec := httptest.NewRecorder()

	err := guard.Basic(context.Background(), rec, req, "email", nil)
	if err != nil {
		t.Fatalf("expected nil when already logged in, got %v", err)
	}
}

// Upstream: testBasicReturnsResponseOnFailure
func TestBasicReturnsResponseOnFailure(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user, hasher: hasher}, sessions, cookies, hasher, nil, clock, fixedIDs{value: "s1"})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic "+basicAuth("user@example.com", "wrong"))
	rec := httptest.NewRecorder()

	err := guard.Basic(context.Background(), rec, req, "email", nil)
	if err == nil {
		t.Fatal("expected error on invalid credentials")
	}
	if rec.Header().Get("WWW-Authenticate") != "Basic" {
		t.Fatal("expected WWW-Authenticate header")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

// Upstream: testBasicWithExtraConditions
func TestBasicWithExtraConditions(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user, hasher: hasher}, sessions, cookies, hasher, nil, clock, fixedIDs{value: "s1"})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic "+basicAuth("user@example.com", "secret"))
	rec := httptest.NewRecorder()

	err := guard.Basic(context.Background(), rec, req, "email", map[string]string{"active": "1"})
	if err != nil {
		t.Fatalf("expected nil with extra conditions, got %v", err)
	}
}

// Upstream: testBasicWithExtraArrayConditions
func TestBasicWithExtraArrayConditions(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user, hasher: hasher}, sessions, cookies, hasher, nil, clock, fixedIDs{value: "s1"})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic "+basicAuth("user@example.com", "secret"))
	rec := httptest.NewRecorder()

	err := guard.Basic(context.Background(), rec, req, "email", map[string]string{"active": "1", "role": "admin"})
	if err != nil {
		t.Fatalf("expected nil with multiple extra conditions, got %v", err)
	}
}

// ======================== REHASHING TESTS ========================

// Upstream: testAttemptRehashesPasswordWhenRequired
func TestAttemptRehashesPasswordWhenRequired(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}
	provider := &fakeProvider{user: user, hasher: hasher, rehashEnabled: true}

	guard := NewSessionGuard("web", defaultCfg(), provider, sessions, cookies, hasher, nil, clock, fixedIDs{value: "s1"})

	ok, err := guard.Attempt(context.Background(), httptest.NewRecorder(), map[string]string{
		"email":    "user@example.com",
		"password": "secret",
	}, false)
	if err != nil {
		t.Fatalf("Attempt: %v", err)
	}
	if !ok {
		t.Fatal("expected Attempt to succeed")
	}
	if !provider.rehashCalled {
		t.Fatal("expected RehashPasswordIfRequired to be called")
	}
}

// Upstream: testAttemptDoesntRehashPasswordWhenDisabled
func TestAttemptDoesntRehashPasswordWhenDisabled(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}
	provider := &fakeProvider{user: user, hasher: hasher, rehashEnabled: false}

	guard := NewSessionGuard("web", defaultCfg(), provider, sessions, cookies, hasher, nil, clock, fixedIDs{value: "s1"})

	ok, err := guard.Attempt(context.Background(), httptest.NewRecorder(), map[string]string{
		"email":    "user@example.com",
		"password": "secret",
	}, false)
	if err != nil {
		t.Fatalf("Attempt: %v", err)
	}
	if !ok {
		t.Fatal("expected Attempt to succeed")
	}
	if provider.rehashCalled {
		t.Fatal("expected RehashPasswordIfRequired NOT to be called")
	}
}

// ======================== FORGET USER / ATTEMPT CALLBACKS ========================

// Upstream: testForgetUserSetsUserToNull
func TestForgetUserSetsUserToNull(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})

	rec := httptest.NewRecorder()
	session, _, err := guard.Login(context.Background(), rec, user, false, false)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	// Logout clears the session — next auth attempt should fail.
	if err := guard.Logout(context.Background(), rec, session, user); err != nil {
		t.Fatalf("Logout: %v", err)
	}

	cookies.values["session"] = session.ID
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	_, _, err = guard.AuthenticateRequest(context.Background(), httptest.NewRecorder(), req)
	if err == nil {
		t.Fatal("expected auth to fail after logout/forget")
	}
}

// Upstream: testAttemptAndWithCallbacks
func TestAttemptAndWithCallbacks(t *testing.T) {
	t.Parallel()

	hasher, _ := NewDefaultPasswordHasher()
	passwordHash, _ := hasher.Hash(context.Background(), "secret")
	user := &testUser{id: "user-1", email: "user@example.com", passwordHash: passwordHash}
	clock := fixedClock{now: time.Now().UTC()}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user, hasher: hasher}, sessions, cookies, hasher, nil, clock, fixedIDs{value: "s1"})

	// Callback rejects the attempt.
	ok, err := guard.AttemptWhen(context.Background(), httptest.NewRecorder(), map[string]string{
		"email":    "user@example.com",
		"password": "secret",
	}, false, func(u Authenticatable) bool {
		return u.GetAuthIdentifier() == "admin" // rejects user-1
	})
	if err != nil {
		t.Fatalf("AttemptWhen: %v", err)
	}
	if ok {
		t.Fatal("expected AttemptWhen to return false when callback rejects")
	}

	// Callback allows the attempt.
	ok, err = guard.AttemptWhen(context.Background(), httptest.NewRecorder(), map[string]string{
		"email":    "user@example.com",
		"password": "secret",
	}, false, func(u Authenticatable) bool {
		return u.GetAuthIdentifier() == "user-1"
	})
	if err != nil {
		t.Fatalf("AttemptWhen: %v", err)
	}
	if !ok {
		t.Fatal("expected AttemptWhen to succeed when callback allows")
	}
}

func basicAuth(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}
