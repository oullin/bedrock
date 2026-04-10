package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedrock/packages/auth"
)

// --- EnsureAuthenticated ---

func TestEnsureAuthenticatedRejects(t *testing.T) {
	provider := &stubProvider{users: map[any]auth.Authenticatable{}}
	guard := auth.NewTokenGuard("api", provider)

	mw := auth.EnsureAuthenticated(guard)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req)
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

func TestEnsureAuthenticatedAllows(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": 1, "api_token": "tok"})
	provider := &stubProvider{users: map[any]auth.Authenticatable{1: user}}
	guard := auth.NewTokenGuard("api", provider)

	mw := auth.EnsureAuthenticated(guard)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer tok")
	guard.SetRequest(req)
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestEnsureAuthenticatedSetsUserInContext(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": 1, "api_token": "tok"})
	provider := &stubProvider{users: map[any]auth.Authenticatable{1: user}}
	guard := auth.NewTokenGuard("api", provider)

	mw := auth.EnsureAuthenticated(guard)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := auth.UserFromContext(r.Context())
		if u == nil {
			t.Error("expected user in context")
		}

		if u.GetAuthIdentifier() != 1 {
			t.Errorf("user id = %v, want 1", u.GetAuthIdentifier())
		}

		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer tok")
	guard.SetRequest(req)
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

// --- RedirectIfAuthenticated ---

func TestRedirectIfAuthenticatedRedirectsWhenLoggedIn(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": 1, "password": "pw"})
	provider := &stubProvider{users: map[any]auth.Authenticatable{1: user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	ctx := context.Background()
	_ = guard.Login(ctx, user, false)

	mw := auth.RedirectIfAuthenticated(guard, "/dashboard")
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Errorf("expected 302, got %d", rr.Code)
	}

	if rr.Header().Get("Location") != "/dashboard" {
		t.Errorf("expected redirect to /dashboard, got %s", rr.Header().Get("Location"))
	}
}

func TestRedirectIfAuthenticatedPassesThroughWhenGuest(t *testing.T) {
	provider := &stubProvider{users: map[any]auth.Authenticatable{}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	mw := auth.RedirectIfAuthenticated(guard, "/dashboard")
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

// --- EnsureEmailIsVerified ---

type verifiedUser struct {
	auth.GenericUser
	verified bool
}

func (u *verifiedUser) HasVerifiedEmail() bool                { return u.verified }
func (u *verifiedUser) MarkEmailAsVerified() error            { return nil }
func (u *verifiedUser) SendEmailVerificationNotification()    {}
func (u *verifiedUser) GetEmailForVerification() string       { return "test@example.com" }

func TestEnsureEmailIsVerifiedRejects(t *testing.T) {
	user := &verifiedUser{
		GenericUser: *auth.NewGenericUser(map[string]any{"id": 1, "password": "pw"}),
		verified:    false,
	}
	provider := &stubProvider{users: map[any]auth.Authenticatable{1: user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	_ = guard.Login(context.Background(), user, false)

	mw := auth.EnsureEmailIsVerified(guard)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", rr.Code)
	}
}

func TestEnsureEmailIsVerifiedAllowsVerifiedUser(t *testing.T) {
	user := &verifiedUser{
		GenericUser: *auth.NewGenericUser(map[string]any{"id": 1, "password": "pw"}),
		verified:    true,
	}
	provider := &stubProvider{users: map[any]auth.Authenticatable{1: user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	_ = guard.Login(context.Background(), user, false)

	mw := auth.EnsureEmailIsVerified(guard)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestEnsureEmailIsVerifiedAllowsUserWithoutMustVerifyEmail(t *testing.T) {
	// User doesn't implement MustVerifyEmail — should pass through.
	user := auth.NewGenericUser(map[string]any{"id": 1, "password": "pw"})
	provider := &stubProvider{users: map[any]auth.Authenticatable{1: user}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)
	_ = guard.Login(context.Background(), user, false)

	mw := auth.EnsureEmailIsVerified(guard)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestEnsureEmailIsVerifiedRejectsUnauthenticatedUser(t *testing.T) {
	provider := &stubProvider{users: map[any]auth.Authenticatable{}}
	sess := newStubSession()
	guard := auth.NewSessionGuard("web", provider, sess, nil, nil)

	mw := auth.EnsureEmailIsVerified(guard)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

// --- RequirePassword ---

func TestRequirePasswordRedirectsWhenNotConfirmed(t *testing.T) {
	sess := newStubSession()

	mw := auth.RequirePassword(sess, 30*time.Minute, "/confirm-password")
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Errorf("expected 302, got %d", rr.Code)
	}

	if rr.Header().Get("Location") != "/confirm-password" {
		t.Errorf("expected redirect to /confirm-password, got %s", rr.Header().Get("Location"))
	}
}

func TestRequirePasswordAllowsRecentlyConfirmed(t *testing.T) {
	sess := newStubSession()
	sess.Put("auth.password_confirmed_at", time.Now().Unix())

	mw := auth.RequirePassword(sess, 30*time.Minute, "/confirm-password")
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestRequirePasswordRedirectsWhenExpired(t *testing.T) {
	sess := newStubSession()
	sess.Put("auth.password_confirmed_at", time.Now().Add(-1*time.Hour).Unix())

	mw := auth.RequirePassword(sess, 30*time.Minute, "/confirm-password")
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/settings", nil)
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusFound {
		t.Errorf("expected 302, got %d", rr.Code)
	}
}

// --- AuthenticateWithBasicAuth ---

func TestAuthenticateWithBasicAuthRejectsNoCredentials(t *testing.T) {
	provider := &stubProvider{users: map[any]auth.Authenticatable{}}
	hasher := auth.NewBcryptHasher(0)

	mw := auth.AuthenticateWithBasicAuth(provider, hasher)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}

	if rr.Header().Get("WWW-Authenticate") == "" {
		t.Error("expected WWW-Authenticate header")
	}
}

func TestAuthenticateWithBasicAuthAcceptsValid(t *testing.T) {
	hasher := auth.NewBcryptHasher(4)
	hash, _ := hasher.Hash("secret")
	user := auth.NewGenericUser(map[string]any{"id": 1, "email": "user@test.com", "password": hash})
	provider := &stubProvider{users: map[any]auth.Authenticatable{1: user}}

	mw := auth.AuthenticateWithBasicAuth(provider, hasher)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := auth.UserFromContext(r.Context())
		if u == nil {
			t.Error("expected user in context after basic auth")
		}

		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetBasicAuth("user@test.com", "secret")
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestAuthenticateWithBasicAuthRejectsWrongPassword(t *testing.T) {
	hasher := auth.NewBcryptHasher(4)
	hash, _ := hasher.Hash("secret")
	user := auth.NewGenericUser(map[string]any{"id": 1, "email": "user@test.com", "password": hash})
	provider := &stubProvider{users: map[any]auth.Authenticatable{1: user}}

	mw := auth.AuthenticateWithBasicAuth(provider, hasher)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetBasicAuth("user@test.com", "wrong")
	mw(inner).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}

// --- WithUser / UserFromContext ---

func TestWithUserAndUserFromContext(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": 1, "password": "pw"})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = auth.WithUser(req, user)

	got := auth.UserFromContext(req.Context())
	if got == nil {
		t.Fatal("expected user from context")
	}

	if got.GetAuthIdentifier() != 1 {
		t.Errorf("user id = %v, want 1", got.GetAuthIdentifier())
	}
}

func TestUserFromContextReturnsNilWhenNotSet(t *testing.T) {
	got := auth.UserFromContext(context.Background())
	if got != nil {
		t.Error("expected nil when no user in context")
	}
}
