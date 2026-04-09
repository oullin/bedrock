package authflows

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// --- test user ---

type testUser struct {
	id       string
	email    string
	password string
	twoFA    bool
}

func (u *testUser) GetAuthIdentifierName() string        { return "id" }
func (u *testUser) GetAuthIdentifier() string            { return u.id }
func (u *testUser) GetAuthPasswordName() string          { return "password" }
func (u *testUser) GetAuthPassword() string              { return u.password }
func (u *testUser) SetAuthPassword(p string)             { u.password = p }
func (u *testUser) GetRememberToken() string             { return "" }
func (u *testUser) SetRememberToken(_ string)            {}
func (u *testUser) GetRememberTokenName() string         { return "remember_token" }
func (u *testUser) IsTwoFactorEnabled() bool             { return u.twoFA }
func (u *testUser) SetTwoFactorEnabled(e bool)           { u.twoFA = e }
func (u *testUser) GetTwoFactorSecret() string           { return "secret" }
func (u *testUser) SetTwoFactorSecret(_ string)          {}
func (u *testUser) GetTwoFactorRecoveryCodes() []string  { return nil }
func (u *testUser) SetTwoFactorRecoveryCodes(_ []string) {}
func (u *testUser) GetTwoFactorConfirmedAt() *time.Time {
	if u.twoFA {
		t := time.Now()
		return &t
	}
	return nil
}
func (u *testUser) SetTwoFactorConfirmedAt(_ *time.Time) {}

// --- test guard ---

type testGuard struct {
	loggedIn          Authenticatable
	loggedInRemember  bool
	pendingTwoFactor  bool
	loggedOut         bool
	authenticatedUser Authenticatable
}

func (g *testGuard) Name() string { return "test" }
func (g *testGuard) AuthenticateRequest(_ context.Context, _ http.ResponseWriter, _ *http.Request) (Authenticatable, error) {
	return g.authenticatedUser, nil
}
func (g *testGuard) Login(_ context.Context, _ http.ResponseWriter, user Authenticatable, remember bool) error {
	g.loggedIn = user
	g.loggedInRemember = remember
	return nil
}
func (g *testGuard) LoginWithPendingTwoFactor(_ context.Context, _ http.ResponseWriter, user Authenticatable) error {
	g.pendingTwoFactor = true
	g.loggedIn = user
	return nil
}
func (g *testGuard) Logout(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
	g.loggedOut = true
	return nil
}

// --- test provider ---

type testProvider struct {
	user           Authenticatable
	validPassword  string
	validateCalled bool
}

func (p *testProvider) RetrieveByID(_ context.Context, _ string) (Authenticatable, error) {
	return p.user, nil
}
func (p *testProvider) RetrieveByToken(_ context.Context, _ string, _ string) (Authenticatable, error) {
	return p.user, nil
}
func (p *testProvider) RetrieveByCredentials(_ context.Context, _ map[string]string) (Authenticatable, error) {
	return p.user, nil
}
func (p *testProvider) UpdateRememberToken(_ context.Context, _ Authenticatable, _ string) error {
	return nil
}
func (p *testProvider) ValidateCredentials(_ context.Context, _ Authenticatable, creds map[string]string) (bool, error) {
	p.validateCalled = true
	return creds["password"] == p.validPassword, nil
}
func (p *testProvider) RehashPasswordIfRequired(_ context.Context, _ Authenticatable, _ map[string]string, _ bool) error {
	return nil
}

// --- test event dispatcher ---

type testEvents struct {
	dispatched []Event
}

func (e *testEvents) Dispatch(_ context.Context, event Event) error {
	e.dispatched = append(e.dispatched, event)
	return nil
}

// --- test responder ---

type testResponder struct {
	loginCalled              bool
	logoutCalled             bool
	twoFactorChallengeCalled bool
}

func (r *testResponder) LoginResponse(w http.ResponseWriter, _ *http.Request) {
	r.loginCalled = true
	w.WriteHeader(http.StatusOK)
}
func (r *testResponder) LogoutResponse(w http.ResponseWriter, _ *http.Request) {
	r.logoutCalled = true
	w.WriteHeader(http.StatusOK)
}
func (r *testResponder) RegisterResponse(_ http.ResponseWriter, _ *http.Request)                  {}
func (r *testResponder) PasswordResetLinkSentResponse(_ http.ResponseWriter, _ *http.Request)     {}
func (r *testResponder) PasswordResetResponse(_ http.ResponseWriter, _ *http.Request)             {}
func (r *testResponder) PasswordUpdateResponse(_ http.ResponseWriter, _ *http.Request)            {}
func (r *testResponder) PasswordConfirmResponse(_ http.ResponseWriter, _ *http.Request)           {}
func (r *testResponder) ProfileInformationUpdatedResponse(_ http.ResponseWriter, _ *http.Request) {}
func (r *testResponder) EmailVerificationSentResponse(_ http.ResponseWriter, _ *http.Request)     {}
func (r *testResponder) TwoFactorChallengeResponse(w http.ResponseWriter, _ *http.Request) {
	r.twoFactorChallengeCalled = true
	w.WriteHeader(http.StatusOK)
}
func (r *testResponder) TwoFactorEnabledResponse(_ http.ResponseWriter, _ *http.Request)  {}
func (r *testResponder) TwoFactorDisabledResponse(_ http.ResponseWriter, _ *http.Request) {}

// --- helpers ---

func buildTestAuthFlows(guard *testGuard, provider *testProvider, events *testEvents, responder *testResponder) *AuthFlows {
	config := DefaultConfig()
	config.Features = Features{}

	return &AuthFlows{
		config:    config,
		guard:     guard,
		provider:  provider,
		hasher:    &stubHasher{},
		events:    events,
		responder: responder,
	}
}

func loginRequest(email string, password string) *http.Request {
	body := strings.NewReader("email=" + email + "&password=" + password)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.RemoteAddr = "127.0.0.1:1234"

	return req
}

func jsonLoginRequest(email string, password string) *http.Request {
	body := strings.NewReader(`{"email":"` + email + `","password":"` + password + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "127.0.0.1:1234"

	return req
}

// --- tests ---

func TestLoginHandlerSuccess(t *testing.T) {
	user := &testUser{id: "1", email: "user@example.com", password: "hashed"}
	guard := &testGuard{}
	provider := &testProvider{user: user, validPassword: "secret"}
	events := &testEvents{}
	responder := &testResponder{}
	f := buildTestAuthFlows(guard, provider, events, responder)

	handler := NewLoginHandler(f)
	w := httptest.NewRecorder()
	r := loginRequest("user@example.com", "secret")

	handler.ServeHTTP(w, r)

	if !responder.loginCalled {
		t.Fatal("expected login response to be called")
	}

	if guard.loggedIn == nil {
		t.Fatal("expected user to be logged in")
	}

	if guard.loggedIn.GetAuthIdentifier() != "1" {
		t.Fatal("expected user id 1")
	}

	if len(events.dispatched) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events.dispatched))
	}

	if events.dispatched[0].Name != EventLoginAttempted {
		t.Fatalf("expected LoginAttempted event, got %s", events.dispatched[0].Name)
	}

	if events.dispatched[1].Name != EventLoginSucceeded {
		t.Fatalf("expected LoginSucceeded event, got %s", events.dispatched[1].Name)
	}
}

func TestLoginHandlerJSONInput(t *testing.T) {
	user := &testUser{id: "1", email: "user@example.com", password: "hashed"}
	guard := &testGuard{}
	provider := &testProvider{user: user, validPassword: "secret"}
	responder := &testResponder{}
	f := buildTestAuthFlows(guard, provider, &testEvents{}, responder)

	handler := NewLoginHandler(f)
	w := httptest.NewRecorder()
	r := jsonLoginRequest("user@example.com", "secret")

	handler.ServeHTTP(w, r)

	if !responder.loginCalled {
		t.Fatal("expected login response for JSON input")
	}

	if guard.loggedIn == nil {
		t.Fatal("expected user to be logged in via JSON")
	}
}

func TestLoginHandlerInvalidCredentials(t *testing.T) {
	user := &testUser{id: "1", email: "user@example.com", password: "hashed"}
	guard := &testGuard{}
	provider := &testProvider{user: user, validPassword: "correct"}
	events := &testEvents{}
	responder := &testResponder{}
	f := buildTestAuthFlows(guard, provider, events, responder)

	handler := NewLoginHandler(f)
	w := httptest.NewRecorder()
	r := loginRequest("user@example.com", "wrong")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}

	if guard.loggedIn != nil {
		t.Fatal("user should not be logged in")
	}

	if len(events.dispatched) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events.dispatched))
	}

	if events.dispatched[1].Name != EventLoginFailed {
		t.Fatalf("expected LoginFailed event, got %s", events.dispatched[1].Name)
	}
}

func TestLoginHandlerRateLimiting(t *testing.T) {
	user := &testUser{id: "1", email: "user@example.com"}
	guard := &testGuard{}
	provider := &testProvider{user: user, validPassword: "correct"}
	responder := &testResponder{}
	f := buildTestAuthFlows(guard, provider, &testEvents{}, responder)
	f.config.LoginRateLimit = 2
	f.config.LoginRateDecay = time.Minute
	f.limiter = &stubLimiter{tooMany: true}

	handler := NewLoginHandler(f)
	w := httptest.NewRecorder()
	r := loginRequest("user@example.com", "correct")

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w.Code)
	}
}

func TestLoginHandlerTwoFactorChallenge(t *testing.T) {
	user := &testUser{id: "1", email: "user@example.com", password: "hashed", twoFA: true}
	guard := &testGuard{}
	provider := &testProvider{user: user, validPassword: "secret"}
	responder := &testResponder{}
	f := buildTestAuthFlows(guard, provider, &testEvents{}, responder)

	handler := NewLoginHandler(f)
	w := httptest.NewRecorder()
	r := loginRequest("user@example.com", "secret")

	handler.ServeHTTP(w, r)

	if !responder.twoFactorChallengeCalled {
		t.Fatal("expected two factor challenge response")
	}

	if !guard.pendingTwoFactor {
		t.Fatal("expected pending two factor login")
	}
}

// --- stub limiter for tests ---

type stubLimiter struct {
	tooMany bool
	hits    int
	cleared bool
}

func (l *stubLimiter) TooManyAttempts(_ string, _ int) bool { return l.tooMany }
func (l *stubLimiter) Hit(_ string, _ time.Duration) int    { l.hits++; return l.hits }
func (l *stubLimiter) Clear(_ string)                       { l.cleared = true }
func (l *stubLimiter) AvailableIn(_ string) time.Duration   { return time.Minute }
