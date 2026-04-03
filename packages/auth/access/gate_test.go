package access_test

import (
	"context"
	"errors"
	"testing"

	auth "github.com/gollin/packages/auth"
	authaccess "github.com/gollin/packages/auth/access"
)

func TestGateAuthorizeAndAny(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("view-dashboard", func(context.Context, auth.Authenticatable, ...any) authaccess.Response {
		return authaccess.Allow()
	})
	gate.Define("delete-users", func(context.Context, auth.Authenticatable, ...any) authaccess.Response {
		return authaccess.Deny("nope")
	})

	user := auth.NewGenericUser(map[string]string{"id": "user-1"})

	if !gate.Check(context.Background(), user, "view-dashboard") {
		t.Fatal("expected ability to be allowed")
	}

	if !gate.Any(context.Background(), user, []string{"delete-users", "view-dashboard"}) {
		t.Fatal("expected Any to allow one ability")
	}

	err := gate.Authorize(context.Background(), user, "delete-users")
	var authErr authaccess.AuthorizationException
	if !errors.As(err, &authErr) {
		t.Fatalf("expected AuthorizationException, got %v", err)
	}
}
