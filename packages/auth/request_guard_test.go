package auth_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/auth"
)

func TestRequestGuardResolvesUserViaCallback(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": 1, "password": "pw"})
	guard := auth.NewRequestGuard(func(_ context.Context, _ *http.Request) (auth.Authenticatable, error) {
		return user, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got != user {
		t.Error("expected callback to return user")
	}
}

func TestRequestGuardCachesResultPerRequest(t *testing.T) {
	callCount := 0
	user := auth.NewGenericUser(map[string]any{"id": 1, "password": "pw"})
	guard := auth.NewRequestGuard(func(_ context.Context, _ *http.Request) (auth.Authenticatable, error) {
		callCount++

		return user, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req)

	ctx := context.Background()
	_, _ = guard.User(ctx)
	_, _ = guard.User(ctx)

	if callCount != 1 {
		t.Errorf("callback called %d times, expected 1 (cached)", callCount)
	}
}

func TestRequestGuardReturnsNilWhenCallbackReturnsNil(t *testing.T) {
	guard := auth.NewRequestGuard(func(_ context.Context, _ *http.Request) (auth.Authenticatable, error) {
		return nil, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Error("expected nil user")
	}
}

func TestRequestGuardReturnsErrorFromCallback(t *testing.T) {
	guard := auth.NewRequestGuard(func(_ context.Context, _ *http.Request) (auth.Authenticatable, error) {
		return nil, errors.New("auth error")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if got != nil {
		t.Error("expected nil user on error")
	}

	if err == nil {
		t.Error("expected error to be propagated")
	}

	if err.Error() != "auth error" {
		t.Errorf("error = %q, want %q", err.Error(), "auth error")
	}
}

func TestRequestGuardReturnsNilWithNoRequest(t *testing.T) {
	guard := auth.NewRequestGuard(func(_ context.Context, _ *http.Request) (auth.Authenticatable, error) {
		return auth.NewGenericUser(map[string]any{"id": 1}), nil
	})

	// No SetRequest called.
	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Error("expected nil when no request set")
	}
}

func TestRequestGuardReturnsNilWithNilCallback(t *testing.T) {
	guard := auth.NewRequestGuard(nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req)

	got, err := guard.User(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if got != nil {
		t.Error("expected nil with nil callback")
	}
}

func TestRequestGuardSetRequestClearsCache(t *testing.T) {
	callCount := 0
	user := auth.NewGenericUser(map[string]any{"id": 1, "password": "pw"})
	guard := auth.NewRequestGuard(func(_ context.Context, _ *http.Request) (auth.Authenticatable, error) {
		callCount++

		return user, nil
	})

	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req1)
	_, _ = guard.User(context.Background())

	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req2)
	_, _ = guard.User(context.Background())

	if callCount != 2 {
		t.Errorf("callback called %d times, expected 2 (cache cleared on SetRequest)", callCount)
	}
}

func TestRequestGuardCheckAndGuest(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": 1, "password": "pw"})
	guard := auth.NewRequestGuard(func(_ context.Context, _ *http.Request) (auth.Authenticatable, error) {
		return user, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req)

	ctx := context.Background()

	if !guard.Check(ctx) {
		t.Error("Check should return true when user is returned")
	}

	if guard.Guest(ctx) {
		t.Error("Guest should return false when user is returned")
	}
}

func TestRequestGuardSetUser(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": 1, "password": "pw"})
	guard := auth.NewRequestGuard(nil)

	guard.SetUser(user)

	if !guard.HasUser() {
		t.Error("HasUser should be true after SetUser")
	}

	got, _ := guard.User(context.Background())

	if got != user {
		t.Error("User() should return user set via SetUser")
	}
}

func TestRequestGuardForgetUser(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": 1, "password": "pw"})
	guard := auth.NewRequestGuard(nil)

	guard.SetUser(user)
	guard.ForgetUser()

	if guard.HasUser() {
		t.Error("HasUser should be false after ForgetUser")
	}
}

func TestRequestGuardHasUserReturnsFalseInitially(t *testing.T) {
	guard := auth.NewRequestGuard(nil)

	if guard.HasUser() {
		t.Error("HasUser should be false initially")
	}
}

func TestRequestGuardIDReturnsIdentifier(t *testing.T) {
	user := auth.NewGenericUser(map[string]any{"id": 42, "password": "pw"})
	guard := auth.NewRequestGuard(func(_ context.Context, _ *http.Request) (auth.Authenticatable, error) {
		return user, nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	guard.SetRequest(req)

	if guard.ID(context.Background()) != 42 {
		t.Errorf("ID = %v, want 42", guard.ID(context.Background()))
	}
}
