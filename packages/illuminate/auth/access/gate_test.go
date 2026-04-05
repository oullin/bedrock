package access_test

import (
	"context"
	"errors"
	"testing"

	auth "github.com/gollin/packages/framework/auth"
	authaccess "github.com/gollin/packages/framework/auth/access"
)

type gateUser struct {
	id string
}

func (u gateUser) GetAuthIdentifierName() string { return "id" }
func (u gateUser) GetAuthIdentifier() string     { return u.id }
func (u gateUser) GetAuthPasswordName() string   { return "password" }
func (u gateUser) GetAuthPassword() string       { return "" }
func (u gateUser) SetAuthPassword(string)        {}
func (u gateUser) GetRememberToken() string      { return "" }
func (u gateUser) SetRememberToken(string)       {}
func (u gateUser) GetRememberTokenName() string  { return "remember_token" }

type account struct {
	OwnerID string
}

func TestGateAuthorizePolicyAndHooks(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("delete-users", func(context.Context, auth.Authenticatable, ...any) authaccess.Response {
		return authaccess.Deny("nope")
	})
	gate.Policy(account{}, func(_ context.Context, user auth.Authenticatable, ability string, arguments ...any) (authaccess.Response, bool) {
		record := arguments[0].(account)
		if ability == "view" && user.GetAuthIdentifier() == record.OwnerID {
			return authaccess.Allow(), true
		}

		return authaccess.Deny("policy denied"), true
	})

	afterCalled := false
	gate.After(func(_ context.Context, _ auth.Authenticatable, ability string, result authaccess.Response, _ ...any) authaccess.Response {
		if ability == "delete-users" {
			afterCalled = true
		}

		return result
	})

	user := gateUser{id: "user-1"}
	if !gate.Check(context.Background(), user, "view", account{OwnerID: "user-1"}) {
		t.Fatal("expected policy to allow resource owner")
	}

	if gate.Check(context.Background(), user, "delete-users") {
		t.Fatal("expected delete-users to be denied")
	}

	if !afterCalled {
		t.Fatal("expected after hook to run")
	}

	err := gate.Authorize(context.Background(), user, "delete-users")
	var authErr authaccess.AuthorizationException
	if !errors.As(err, &authErr) {
		t.Fatalf("expected AuthorizationException, got %v", err)
	}
}
