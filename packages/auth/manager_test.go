package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/bedrock/packages/auth"
	cauth "github.com/bedrock/packages/contracts/auth"
)

// --- Test helpers ---

// stubProvider is a test UserProvider backed by a map.
type stubProvider struct {
	users map[string]cauth.Authenticatable
}

// stubSession is a minimal in-memory SessionStore.
type stubSession struct {
	data map[string]any
}

// stubCookieManager is a minimal in-memory CookieManager.
type stubCookieManager struct {
	queued    []*http.Cookie
	forgotten []string
}

// recordingDispatcher collects dispatched events for assertions.
type recordingDispatcher struct {
	mu     sync.Mutex
	events []any
}

func (p *stubProvider) RetrieveByID(_ context.Context, id string) (cauth.Authenticatable, error) {
	return p.users[id], nil
}

func (p *stubProvider) RetrieveByToken(_ context.Context, id string, token string) (cauth.Authenticatable, error) {
	u := p.users[id]

	if u == nil || u.GetRememberToken() != token {
		return nil, nil
	}

	return u, nil
}

func (p *stubProvider) UpdateRememberToken(_ context.Context, user cauth.Authenticatable, token string) error {
	user.SetRememberToken(token)

	return nil
}

func (p *stubProvider) RetrieveByCredentials(_ context.Context, creds map[string]string) (cauth.Authenticatable, error) {
	for _, u := range p.users {
		gen, ok := u.(*auth.GenericUser)

		if !ok {
			continue
		}

		match := true

		for k, v := range creds {
			if k == "password" {
				continue
			}

			attr, _ := gen.Attributes[k].(string)

			if attr != v {
				match = false

				break
			}
		}

		if match {
			return u, nil
		}
	}

	return nil, nil
}

func (p *stubProvider) ValidateCredentials(_ context.Context, user cauth.Authenticatable, creds map[string]string) (bool, error) {
	pw := creds["password"]

	return user.GetAuthPassword() == pw, nil
}

func (p *stubProvider) RehashPasswordIfRequired(_ context.Context, _ cauth.Authenticatable, _ map[string]string, _ bool) error {
	return nil
}

func newStubSession() *stubSession {
	return &stubSession{data: make(map[string]any)}
}

func (s *stubSession) Get(key string, fallback any) any {
	if v, ok := s.data[key]; ok {
		return v
	}

	return fallback
}

func (s *stubSession) Put(key string, value any) { s.data[key] = value }

func (s *stubSession) Remove(key string) any {
	v := s.data[key]
	delete(s.data, key)

	return v
}

func (s *stubSession) Forget(keys ...string) {
	for _, k := range keys {
		delete(s.data, k)
	}
}

func (s *stubSession) Migrate(_ context.Context, _ bool) error { return nil }

func (m *stubCookieManager) Queue(cookie *http.Cookie) {
	m.queued = append(m.queued, cookie)
}

func (m *stubCookieManager) Forget(name, path, domain string) *http.Cookie {
	m.forgotten = append(m.forgotten, name)

	return &http.Cookie{Name: name, Path: path, Domain: domain, MaxAge: -1}
}

func (d *recordingDispatcher) Dispatch(_ context.Context, event any) error {
	d.mu.Lock()

	defer d.mu.Unlock()

	d.events = append(d.events, event)

	return nil
}

func (d *recordingDispatcher) has(t *testing.T, typeName string) {
	t.Helper()

	d.mu.Lock()

	defer d.mu.Unlock()

	for _, e := range d.events {
		if typeNameOf(e) == typeName {
			return
		}
	}

	t.Errorf("expected event %q to be dispatched, got %v", typeName, d.typeNames())
}

func (d *recordingDispatcher) hasNot(t *testing.T, typeName string) {
	t.Helper()

	d.mu.Lock()

	defer d.mu.Unlock()

	for _, e := range d.events {
		if typeNameOf(e) == typeName {
			t.Errorf("expected event %q NOT to be dispatched", typeName)

			return
		}
	}
}

func (d *recordingDispatcher) typeNames() []string {
	names := make([]string, len(d.events))

	for i, e := range d.events {
		names[i] = typeNameOf(e)
	}

	return names
}

func typeNameOf(v any) string {
	return fmt.Sprintf("%T", v)
}

// --- GenericUser ---

func TestGenericUser(t *testing.T) {
	u := auth.NewGenericUser(map[string]any{
		"id":       "42",
		"password": "secret",
	})

	if u.GetAuthIdentifier() != "42" {
		t.Errorf("unexpected id: %v", u.GetAuthIdentifier())
	}

	if u.GetAuthPassword() != "secret" {
		t.Errorf("unexpected password: %s", u.GetAuthPassword())
	}

	u.SetRememberToken("tok")

	if u.GetRememberToken() != "tok" {
		t.Error("remember token not stored")
	}
}

func TestGenericUserIdentifierName(t *testing.T) {
	u := auth.NewGenericUser(map[string]any{"id": "1"})

	if u.GetAuthIdentifierName() != "id" {
		t.Errorf("identifier name = %q, want %q", u.GetAuthIdentifierName(), "id")
	}
}

func TestGenericUserRememberTokenName(t *testing.T) {
	u := auth.NewGenericUser(map[string]any{"id": "1"})

	if u.GetRememberTokenName() != "remember_token" {
		t.Errorf("token name = %q, want %q", u.GetRememberTokenName(), "remember_token")
	}
}

func TestGenericUserEmptyPassword(t *testing.T) {
	u := auth.NewGenericUser(map[string]any{"id": "1"})

	if u.GetAuthPassword() != "" {
		t.Errorf("expected empty password, got %q", u.GetAuthPassword())
	}
}

func TestGenericUserEmptyRememberToken(t *testing.T) {
	u := auth.NewGenericUser(map[string]any{"id": "1"})

	if u.GetRememberToken() != "" {
		t.Errorf("expected empty remember token, got %q", u.GetRememberToken())
	}
}

func TestGenericUserGetAuthPasswordName(t *testing.T) {
	u := auth.NewGenericUser(map[string]any{"id": "1"})

	if u.GetAuthPasswordName() != "password" {
		t.Errorf("password name = %q, want %q", u.GetAuthPasswordName(), "password")
	}
}

func TestGenericUserSetAuthPassword(t *testing.T) {
	u := auth.NewGenericUser(map[string]any{"id": "1"})

	u.SetAuthPassword("newpw")

	if u.GetAuthPassword() != "newpw" {
		t.Errorf("expected password %q, got %q", "newpw", u.GetAuthPassword())
	}
}

// --- Recaller ---

func TestRecallerValid(t *testing.T) {
	r := auth.NewRecaller("1|tok|hash")

	if r == nil || !r.Valid() {
		t.Fatal("expected valid recaller")
	}

	if r.ID() != "1" || r.Token() != "tok" || r.Hash() != "hash" {
		t.Errorf("unexpected parts: %s %s %s", r.ID(), r.Token(), r.Hash())
	}
}

func TestRecallerInvalidFormat(t *testing.T) {
	if auth.NewRecaller("bad") != nil {
		t.Error("expected nil for single-part value")
	}

	if auth.NewRecaller("one|two") != nil {
		t.Error("expected nil for two-part value")
	}

	if auth.NewRecaller("one|two|three|four") != nil {
		t.Error("expected nil for four-part value")
	}
}

func TestRecallerEmptySegments(t *testing.T) {
	r := auth.NewRecaller("|tok|hash")

	if r != nil && r.Valid() {
		t.Error("expected invalid recaller with empty ID")
	}

	r = auth.NewRecaller("1||hash")

	if r != nil && r.Valid() {
		t.Error("expected invalid recaller with empty token")
	}

	r = auth.NewRecaller("1|tok|")

	if r != nil && r.Valid() {
		t.Error("expected invalid recaller with empty hash")
	}
}

// --- Timebox ---

func TestTimebox(t *testing.T) {
	start := time.Now()
	auth.Timebox(50*time.Millisecond, func() {})

	if elapsed := time.Since(start); elapsed < 40*time.Millisecond {
		t.Errorf("Timebox did not wait: elapsed %v", elapsed)
	}
}

func TestTimeboxDoesNotDelayLongOperations(t *testing.T) {
	start := time.Now()
	auth.Timebox(10*time.Millisecond, func() {
		time.Sleep(50 * time.Millisecond)
	})

	elapsed := time.Since(start)

	if elapsed < 40*time.Millisecond {
		t.Errorf("Timebox should have taken at least the fn duration: %v", elapsed)
	}
}

// --- BcryptHasher ---

func TestBcryptHasher(t *testing.T) {
	h := auth.NewBcryptHasher(0)
	ctx := context.Background()

	hash, err := h.Hash(ctx, "password123")

	if err != nil {
		t.Fatal(err)
	}

	match, err := h.Check(ctx, "password123", hash)

	if err != nil {
		t.Fatal(err)
	}

	if !match {
		t.Error("Check should return true for matching password")
	}

	match, err = h.Check(ctx, "wrong", hash)

	if err != nil {
		t.Fatal(err)
	}

	if match {
		t.Error("Check should return false for wrong password")
	}
}

func TestBcryptHasherNeedsRehash(t *testing.T) {
	h4 := auth.NewBcryptHasher(4)
	h10 := auth.NewBcryptHasher(10)
	ctx := context.Background()

	hash, err := h4.Hash(ctx, "pw")

	if err != nil {
		t.Fatal(err)
	}

	if h4.NeedsRehash(hash) {
		t.Error("should not need rehash with same cost")
	}

	if !h10.NeedsRehash(hash) {
		t.Error("should need rehash with different cost")
	}
}

func TestBcryptHasherNeedsRehashInvalidHash(t *testing.T) {
	h := auth.NewBcryptHasher(0)

	if !h.NeedsRehash("not-a-valid-hash") {
		t.Error("should need rehash for invalid hash")
	}
}

// --- Manager ---

func TestManagerGuardResolvesDefaultGuard(t *testing.T) {
	m := auth.NewManager("web")
	m.Extend("session", func(name string, config map[string]any, provider cauth.UserProvider) (cauth.Guard, error) {
		return auth.NewSessionGuard(name, provider, newStubSession(), nil, nil), nil
	})
	m.SetConfig("web", map[string]any{"driver": "session"})

	g, err := m.Guard(context.Background(), "")

	if err != nil {
		t.Fatal(err)
	}

	if g == nil {
		t.Error("expected non-nil guard")
	}
}

func TestManagerGuardResolvesNamedGuard(t *testing.T) {
	m := auth.NewManager("web")
	m.Extend("token", func(name string, config map[string]any, provider cauth.UserProvider) (cauth.Guard, error) {
		return auth.NewTokenGuard(name, provider), nil
	})
	m.SetConfig("api", map[string]any{"driver": "token"})

	g, err := m.Guard(context.Background(), "api")

	if err != nil {
		t.Fatal(err)
	}

	if g == nil {
		t.Error("expected non-nil guard")
	}
}

func TestManagerGuardCachesInstances(t *testing.T) {
	m := auth.NewManager("web")
	m.Extend("session", func(name string, config map[string]any, provider cauth.UserProvider) (cauth.Guard, error) {
		return auth.NewSessionGuard(name, provider, newStubSession(), nil, nil), nil
	})
	m.SetConfig("web", map[string]any{"driver": "session"})

	g1, _ := m.Guard(context.Background(), "web")
	g2, _ := m.Guard(context.Background(), "web")

	if g1 != g2 {
		t.Error("Guard should cache and return same instance")
	}
}

func TestManagerGuardReturnsErrorForUnknownDriver(t *testing.T) {
	m := auth.NewManager("web")
	m.SetConfig("web", map[string]any{"driver": "unknown"})

	_, err := m.Guard(context.Background(), "web")

	if err == nil {
		t.Error("expected error for unknown driver")
	}
}

func TestManagerViaRequest(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1"})
	m := auth.NewManager("custom")
	m.ViaRequest("custom", func(_ context.Context, _ *http.Request) (cauth.Authenticatable, error) {
		return user, nil
	})

	g, err := m.Guard(context.Background(), "custom")

	if err != nil {
		t.Fatal(err)
	}

	if g == nil {
		t.Error("expected ViaRequest guard")
	}
}

func TestManagerSetRequestPropagates(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": "1", "api_token": "tok"})
	provider := &stubProvider{users: map[string]cauth.Authenticatable{"1": user}}

	m := auth.NewManager("api")
	m.Extend("token", func(name string, config map[string]any, p cauth.UserProvider) (cauth.Guard, error) {
		return auth.NewTokenGuard(name, p), nil
	})
	m.Provider("users", func(config map[string]any) (cauth.UserProvider, error) {
		return provider, nil
	})
	m.SetConfig("api", map[string]any{"driver": "token", "provider": "users"})

	// Resolve the guard first.
	g, _ := m.Guard(context.Background(), "api")

	// Then set request.
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer tok")
	m.SetRequest(req)

	u, _ := g.User(context.Background())

	if u == nil {
		t.Error("expected user after SetRequest propagation")
	}
}

// --- Errors ---

func TestAuthenticationException(t *testing.T) {
	e := auth.NewAuthenticationException([]string{"web", "api"}, "/login")

	if e.Error() != "unauthenticated" {
		t.Errorf("Error() = %q, want %q", e.Error(), "unauthenticated")
	}

	if len(e.Guards) != 2 {
		t.Errorf("expected 2 guards, got %d", len(e.Guards))
	}

	if e.RedirectPath != "/login" {
		t.Errorf("RedirectPath = %q, want %q", e.RedirectPath, "/login")
	}
}

func TestAuthorizationException(t *testing.T) {
	e := auth.NewAuthorizationException("forbidden", 403)

	if e.Error() != "forbidden" {
		t.Errorf("Error() = %q, want %q", e.Error(), "forbidden")
	}

	if e.StatusCode != 403 {
		t.Errorf("StatusCode = %d, want 403", e.StatusCode)
	}
}

func TestAuthorizationExceptionDefaultStatus(t *testing.T) {
	e := auth.NewAuthorizationException("nope", 0)

	if e.StatusCode != 403 {
		t.Errorf("StatusCode = %d, want 403 (default)", e.StatusCode)
	}
}

func TestAuthenticationExceptionCustomMessage(t *testing.T) {
	e := &auth.AuthenticationException{Message: "custom msg"}

	if e.Error() != "custom msg" {
		t.Errorf("Error() = %q, want %q", e.Error(), "custom msg")
	}
}

func TestAuthorizationExceptionEmptyMessage(t *testing.T) {
	e := &auth.AuthorizationException{StatusCode: 403}

	if e.Error() == "" {
		t.Error("expected non-empty default message")
	}
}
