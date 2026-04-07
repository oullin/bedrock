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

// ======================== MULTI-GUARD AUTHENTICATE TESTS ========================

// Laravel: testDefaultUnauthenticatedThrowsWithGuards
func TestAuthenticateWithGuardsRejectsUnauthenticated(t *testing.T) {
	t.Parallel()

	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	guard1 := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: nil}, &fakeSessionStore{}, &fakeCookieManager{}, nil, nil, clock, fixedIDs{value: "s1"})
	guard2 := NewSessionGuard("api", defaultCfg(), &fakeProvider{user: nil}, &fakeSessionStore{}, &fakeCookieManager{}, nil, nil, clock, fixedIDs{value: "s2"})

	handler := Authenticate(guard1, guard2)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 when all guards fail, got %d", rec.Code)
	}
}

// Laravel: testSecondaryAuthenticatedUpdatesDefaultDriver
func TestSecondaryAuthenticatedUpdatesDefaultDriver(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}

	// First guard has no sessions (will fail).
	guard1 := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: nil}, &fakeSessionStore{}, &fakeCookieManager{}, nil, nil, clock, fixedIDs{value: "s1"})

	// Second guard has a valid session.
	sessions2 := &fakeSessionStore{}
	cookies2 := &fakeCookieManager{}
	guard2 := NewSessionGuard("api", defaultCfg(), &fakeProvider{user: user}, sessions2, cookies2, nil, nil, clock, fixedIDs{value: "s2"})
	guard2.Login(context.Background(), httptest.NewRecorder(), user, false, false)
	cookies2.values["session"] = "s2"

	var contextUser Authenticatable
	handler := Authenticate(guard1, guard2)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextUser, _ = UserFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 from second guard, got %d", rec.Code)
	}
	if contextUser == nil || contextUser.GetAuthIdentifier() != "user-1" {
		t.Fatal("expected user from second guard in context")
	}
}

// Laravel: testMultipleDriversUnauthenticatedThrows
func TestMultipleDriversUnauthenticatedThrows(t *testing.T) {
	t.Parallel()

	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	guard1 := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: nil}, &fakeSessionStore{}, &fakeCookieManager{}, nil, nil, clock, fixedIDs{value: "s1"})
	guard2 := NewSessionGuard("admin", defaultCfg(), &fakeProvider{user: nil}, &fakeSessionStore{}, &fakeCookieManager{}, nil, nil, clock, fixedIDs{value: "s2"})

	handlerCalled := false
	handler := Authenticate(guard1, guard2)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	if handlerCalled {
		t.Fatal("expected handler NOT to be called")
	}
}

// Laravel: testMultipleDriversUnauthenticatedThrowsWithGuards
func TestMultipleDriversUnauthenticatedThrowsWithGuards(t *testing.T) {
	t.Parallel()

	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	guard1 := NewSessionGuard("session", defaultCfg(), &fakeProvider{user: nil}, &fakeSessionStore{}, &fakeCookieManager{}, nil, nil, clock, fixedIDs{value: "s1"})
	guard2 := NewSessionGuard("token", defaultCfg(), &fakeProvider{user: nil}, &fakeSessionStore{}, &fakeCookieManager{}, nil, nil, clock, fixedIDs{value: "s2"})

	handler := Authenticate(guard1, guard2)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 with named guards, got %d", rec.Code)
	}
}

// Laravel: testMultipleDriversAuthenticatedUpdatesDefault
func TestMultipleDriversAuthenticatedUpdatesDefault(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "admin-1", email: "admin@example.com"}
	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}

	sessions1 := &fakeSessionStore{}
	cookies1 := &fakeCookieManager{}
	guard1 := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions1, cookies1, nil, nil, clock, fixedIDs{value: "s1"})
	guard1.Login(context.Background(), httptest.NewRecorder(), user, false, false)
	cookies1.values["session"] = "s1"

	guard2 := NewSessionGuard("api", defaultCfg(), &fakeProvider{user: nil}, &fakeSessionStore{}, &fakeCookieManager{}, nil, nil, clock, fixedIDs{value: "s2"})

	var contextUser Authenticatable
	handler := Authenticate(guard1, guard2)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contextUser, _ = UserFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if contextUser == nil || contextUser.GetAuthIdentifier() != "admin-1" {
		t.Fatal("expected first guard's user in context")
	}
}

// ======================== AUTHORIZE WITH PARAMETERS TESTS ========================

// Laravel: testSimpleAbilityWithStringParameter
func TestAuthorizeWithStringParameter(t *testing.T) {
	t.Parallel()

	var receivedArgs []any
	authorizeFn := AuthorizeFunc(func(_ context.Context, _ Authenticatable, _ string, arguments ...any) error {
		receivedArgs = arguments
		return nil
	})

	handler := Authorize(authorizeFn, "edit", "post-123")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(WithUser(req.Context(), &testUser{id: "u1"}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if len(receivedArgs) != 1 || receivedArgs[0] != "post-123" {
		t.Fatalf("expected string parameter 'post-123', got %v", receivedArgs)
	}
}

// Laravel: testSimpleAbilityWithNullParameter
func TestAuthorizeWithNilParameter(t *testing.T) {
	t.Parallel()

	var receivedArgs []any
	authorizeFn := AuthorizeFunc(func(_ context.Context, _ Authenticatable, _ string, arguments ...any) error {
		receivedArgs = arguments
		return nil
	})

	handler := Authorize(authorizeFn, "edit", nil)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(WithUser(req.Context(), &testUser{id: "u1"}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if len(receivedArgs) != 1 || receivedArgs[0] != nil {
		t.Fatalf("expected nil parameter, got %v", receivedArgs)
	}
}

// Laravel: testSimpleAbilityWithOptionalParameter
func TestAuthorizeWithOptionalParameter(t *testing.T) {
	t.Parallel()

	var receivedArgs []any
	authorizeFn := AuthorizeFunc(func(_ context.Context, _ Authenticatable, _ string, arguments ...any) error {
		receivedArgs = arguments
		return nil
	})

	handler := Authorize(authorizeFn, "edit")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(WithUser(req.Context(), &testUser{id: "u1"}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if len(receivedArgs) != 0 {
		t.Fatalf("expected no parameters, got %v", receivedArgs)
	}
}

// Laravel: testSimpleAbilityWithStringParameterFromRouteParameter
func TestAuthorizeWithStringParameterFromRouteParameter(t *testing.T) {
	t.Parallel()

	var receivedArgs []any
	authorizeFn := AuthorizeFunc(func(_ context.Context, _ Authenticatable, _ string, arguments ...any) error {
		receivedArgs = arguments
		return nil
	})

	handler := Authorize(authorizeFn, "edit", "route-param-value")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/posts/route-param-value", nil)
	req = req.WithContext(WithUser(req.Context(), &testUser{id: "u1"}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if len(receivedArgs) != 1 || receivedArgs[0] != "route-param-value" {
		t.Fatalf("expected route param value, got %v", receivedArgs)
	}
}

// Laravel: testSimpleAbilityWithStringParameter0FromRouteParameter
func TestAuthorizeWithStringParameter0FromRouteParameter(t *testing.T) {
	t.Parallel()

	var receivedArgs []any
	authorizeFn := AuthorizeFunc(func(_ context.Context, _ Authenticatable, _ string, arguments ...any) error {
		receivedArgs = arguments
		return nil
	})

	handler := Authorize(authorizeFn, "edit", "0")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(WithUser(req.Context(), &testUser{id: "u1"}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if len(receivedArgs) != 1 || receivedArgs[0] != "0" {
		t.Fatalf("expected '0' parameter, got %v", receivedArgs)
	}
}

// Laravel: testModelInstanceAsParameter (using multiple params)
func TestAuthorizeWithMultipleParameters(t *testing.T) {
	t.Parallel()

	var receivedArgs []any
	authorizeFn := AuthorizeFunc(func(_ context.Context, _ Authenticatable, _ string, arguments ...any) error {
		receivedArgs = arguments
		return nil
	})

	handler := Authorize(authorizeFn, "edit", "post-1", "comment-2")(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(WithUser(req.Context(), &testUser{id: "u1"}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if len(receivedArgs) != 2 || receivedArgs[0] != "post-1" || receivedArgs[1] != "comment-2" {
		t.Fatalf("expected two parameters, got %v", receivedArgs)
	}
}

// ======================== REDIRECT IF AUTHENTICATED TESTS ========================

// Laravel: RedirectIfAuthenticatedMiddlewareTest
func TestRedirectIfAuthenticatedRedirectsAuthenticatedUser(t *testing.T) {
	t.Parallel()

	user := &testUser{id: "user-1", email: "user@example.com"}
	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	cookies := &fakeCookieManager{}
	sessions := &fakeSessionStore{}

	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: user}, sessions, cookies, nil, nil, clock, fixedIDs{value: "s1"})
	guard.Login(context.Background(), httptest.NewRecorder(), user, false, false)
	cookies.values["session"] = "s1"

	handlerCalled := false
	handler := RedirectIfAuthenticated(guard, "/dashboard")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/login", nil))

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302 redirect, got %d", rec.Code)
	}
	if rec.Header().Get("Location") != "/dashboard" {
		t.Fatalf("expected redirect to /dashboard, got %q", rec.Header().Get("Location"))
	}
	if handlerCalled {
		t.Fatal("expected handler NOT to be called for authenticated user")
	}
}

func TestRedirectIfAuthenticatedAllowsGuest(t *testing.T) {
	t.Parallel()

	clock := fixedClock{now: time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)}
	guard := NewSessionGuard("web", defaultCfg(), &fakeProvider{user: nil}, &fakeSessionStore{}, &fakeCookieManager{}, nil, nil, clock, fixedIDs{value: "s1"})

	handlerCalled := false
	handler := RedirectIfAuthenticated(guard, "/dashboard")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/login", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if !handlerCalled {
		t.Fatal("expected handler to be called for guest")
	}
}
