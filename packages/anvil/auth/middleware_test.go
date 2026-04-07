package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// ======================== AUTHENTICATE MIDDLEWARE TESTS ========================

// Laravel: AuthenticateMiddlewareTest
func TestAuthenticateMiddlewareAllowsAuthenticated(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})

	// Log the user in first.
	guard.Login(context.Background(), httptest.NewRecorder(), user, false, false)
	cookies.values["session"] = "s1"

	handler := Authenticate(guard)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := UserFromContext(r.Context())
		if !ok {
			t.Fatal("expected user in context")
		}
		if u.GetAuthIdentifier() != "user-1" {
			t.Fatalf("expected user-1, got %q", u.GetAuthIdentifier())
		}
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/dashboard", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestAuthenticateMiddlewareRejectsUnauthenticated(t *testing.T) {
	t.Parallel()

	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: nil}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})

	handlerCalled := false
	handler := Authenticate(guard)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/dashboard", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	if handlerCalled {
		t.Fatal("expected handler NOT to be called for unauthenticated request")
	}
}

// ======================== AUTHORIZE MIDDLEWARE TESTS ========================

// Laravel: AuthorizeMiddlewareTest
func TestAuthorizeMiddlewareAllowsAuthorized(t *testing.T) {
	t.Parallel()

	authorizeFn := AuthorizeFunc(func(_ context.Context, user Authenticatable, ability string, _ ...any) error {
		if user.GetAuthIdentifier() == "admin" && ability == "edit" {
			return nil
		}
		return fmt.Errorf("denied")
	})

	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := Authorize(authorizeFn, "edit")(inner)

	user := &testUser{id: "admin", email: "admin@example.com"}
	req := httptest.NewRequest(http.MethodGet, "/edit", nil)
	req = req.WithContext(WithUser(req.Context(), user))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestAuthorizeMiddlewareDeniesUnauthorized(t *testing.T) {
	t.Parallel()

	authorizeFn := AuthorizeFunc(func(_ context.Context, _ Authenticatable, _ string, _ ...any) error {
		return fmt.Errorf("denied")
	})

	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := Authorize(authorizeFn, "edit")(inner)

	user := &testUser{id: "reader", email: "reader@example.com"}
	req := httptest.NewRequest(http.MethodGet, "/edit", nil)
	req = req.WithContext(WithUser(req.Context(), user))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestAuthorizeMiddlewareRejectsNoUser(t *testing.T) {
	t.Parallel()

	authorizeFn := AuthorizeFunc(func(_ context.Context, _ Authenticatable, _ string, _ ...any) error {
		return nil
	})

	handler := Authorize(authorizeFn, "edit")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/edit", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when no user in context, got %d", rec.Code)
	}
}

// ======================== ENSURE EMAIL IS VERIFIED TESTS ========================

// Laravel: EnsureEmailIsVerifiedTest
func TestEnsureEmailIsVerifiedAllowsVerified(t *testing.T) {
	t.Parallel()

	now := time.Now()
	user := &testUser{id: "user-1", email: "user@example.com", verifiedAt: &now}

	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := EnsureEmailIsVerified(inner)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req = req.WithContext(WithUser(req.Context(), user))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestEnsureEmailIsVerifiedRejectsUnverified(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com", verifiedAt: nil}

	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := EnsureEmailIsVerified(inner)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req = req.WithContext(WithUser(req.Context(), user))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rec.Code)
	}
}

func TestEnsureEmailIsVerifiedRejectsNoUser(t *testing.T) {
	t.Parallel()

	handler := EnsureEmailIsVerified(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/profile", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

// ======================== CONTEXT HELPERS TESTS ========================

func TestUserContextRoundtrip(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	ctx := WithUser(context.Background(), user)

	got, ok := UserFromContext(ctx)
	if !ok {
		t.Fatal("expected user in context")
	}
	if got.GetAuthIdentifier() != "user-1" {
		t.Fatalf("expected user-1, got %q", got.GetAuthIdentifier())
	}
}

func TestUserFromContextReturnsFalseWhenMissing(t *testing.T) {
	t.Parallel()

	_, ok := UserFromContext(context.Background())
	if ok {
		t.Fatal("expected no user in empty context")
	}
}
