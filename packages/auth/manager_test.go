package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedrock/packages/auth"
)

func TestGenericUser(t *testing.T) {
	u := auth.NewGenericUser(map[string]any{
		"id":       42,
		"password": "secret",
	})

	if u.GetAuthIdentifier() != 42 {
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

func TestRecaller(t *testing.T) {
	r := auth.NewRecaller("1|tok|hash")
	if r == nil || !r.Valid() {
		t.Fatal("expected valid recaller")
	}

	if r.ID() != "1" || r.Token() != "tok" || r.Hash() != "hash" {
		t.Errorf("unexpected parts: %s %s %s", r.ID(), r.Token(), r.Hash())
	}

	if auth.NewRecaller("bad") != nil {
		t.Error("expected nil for invalid cookie value")
	}
}

func TestTimebox(t *testing.T) {
	start := time.Now()
	auth.Timebox(50*time.Millisecond, func() {})

	if elapsed := time.Since(start); elapsed < 40*time.Millisecond {
		t.Errorf("Timebox did not wait: elapsed %v", elapsed)
	}
}

func TestBcryptHasher(t *testing.T) {
	h := auth.NewBcryptHasher(0)

	hash, err := h.Hash("password123")
	if err != nil {
		t.Fatal(err)
	}

	if !h.Check("password123", hash) {
		t.Error("Check should return true for matching password")
	}

	if h.Check("wrong", hash) {
		t.Error("Check should return false for wrong password")
	}
}

// stubProvider is a test UserProvider backed by a map.
type stubProvider struct {
	users map[any]auth.Authenticatable
}

func (p *stubProvider) RetrieveByID(_ context.Context, id any) (auth.Authenticatable, error) {
	return p.users[id], nil
}

func (p *stubProvider) RetrieveByToken(_ context.Context, id any, token string) (auth.Authenticatable, error) {
	u := p.users[id]
	if u == nil || u.GetRememberToken() != token {
		return nil, nil
	}

	return u, nil
}

func (p *stubProvider) UpdateRememberToken(_ context.Context, user auth.Authenticatable, token string) error {
	user.SetRememberToken(token)

	return nil
}

func (p *stubProvider) RetrieveByCredentials(_ context.Context, creds map[string]any) (auth.Authenticatable, error) {
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

			if gen.Attributes[k] != v {
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

func (p *stubProvider) ValidateCredentials(_ context.Context, user auth.Authenticatable, creds map[string]any) bool {
	pw, _ := creds["password"].(string)

	return user.GetAuthPassword() == pw
}

func (p *stubProvider) RehashPasswordIfRequired(_ context.Context, _ auth.Authenticatable, _ map[string]any, _ bool) error {
	return nil
}

// stubSession is a minimal in-memory SessionStore.
type stubSession struct {
	data map[string]any
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

func TestSessionGuardLoginLogout(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": 1, "password": "pw"})
	provider := &stubProvider{users: map[any]auth.Authenticatable{1: user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	ctx := context.Background()

	if guard.Check(ctx) {
		t.Fatal("should not be authenticated initially")
	}

	if err := guard.Login(ctx, user, false); err != nil {
		t.Fatal(err)
	}

	if !guard.Check(ctx) {
		t.Error("should be authenticated after login")
	}

	if err := guard.Logout(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestTokenGuardBearerHeader(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": 1, "api_token": "mytoken"})
	provider := &stubProvider{users: map[any]auth.Authenticatable{1: user}}

	guard := auth.NewTokenGuard(provider)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer mytoken")
	guard.SetRequest(req)

	got, err := guard.User(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if got == nil {
		t.Error("expected authenticated user")
	}
}

func TestEnsureAuthenticatedMiddleware(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": 1, "api_token": "tok"})
	provider := &stubProvider{users: map[any]auth.Authenticatable{1: user}}

	guard := auth.NewTokenGuard(provider)
	mw := auth.EnsureAuthenticated(guard)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Unauthenticated request.
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req)
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}

	// Authenticated request.
	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.Header.Set("Authorization", "Bearer tok")
	guard.SetRequest(req2)
	mw(inner).ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr2.Code)
	}
}
