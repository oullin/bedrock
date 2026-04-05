package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gollin/packages/illuminate/encryption"
	securitycrypto "github.com/gollin/packages/illuminate/support/crypto"
)

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
	user   *testUser
	hasher PasswordHasher
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

func (p *fakeProvider) RehashPasswordIfRequired(context.Context, Authenticatable, map[string]string, bool) error {
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
	cfg := Config{
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

func TestRequestAndTokenGuards(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	provider := &fakeProvider{user: user}

	request := httptest.NewRequest(http.MethodGet, "/?api_token=token-123", nil)
	tokenGuard := NewTokenGuard(provider, request, "api_token", "api_token", false)
	resolved, err := tokenGuard.User(context.Background())
	if err != nil {
		t.Fatalf("TokenGuard.User query: %v", err)
	}

	if resolved.GetAuthIdentifier() != user.id {
		t.Fatalf("unexpected token guard user: %q", resolved.GetAuthIdentifier())
	}

	request = httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer token-123")
	tokenGuard.SetRequest(request)

	resolved, err = tokenGuard.User(context.Background())
	if err != nil {
		t.Fatalf("TokenGuard.User bearer: %v", err)
	}

	if resolved.GetAuthIdentifier() != user.id {
		t.Fatalf("unexpected bearer user: %q", resolved.GetAuthIdentifier())
	}

	requestGuard := NewRequestGuard(func(_ context.Context, r *http.Request, provider UserProvider) (Authenticatable, error) {
		if r.URL.Path != "/me" {
			return nil, ErrUnauthorized
		}

		return provider.RetrieveByID(context.Background(), user.id)
	}, httptest.NewRequest(http.MethodGet, "/me", nil), provider)

	resolved, err = requestGuard.User(context.Background())
	if err != nil {
		t.Fatalf("RequestGuard.User: %v", err)
	}

	if resolved.GetAuthIdentifier() != user.id {
		t.Fatalf("unexpected request guard user: %q", resolved.GetAuthIdentifier())
	}
}

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
