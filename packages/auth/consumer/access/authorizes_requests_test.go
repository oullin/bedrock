package access_test

import (
	"context"
	"testing"

	auth "github.com/bedrock/packages/auth"
	authaccess "github.com/bedrock/packages/auth/access"
	consumeraccess "github.com/bedrock/packages/auth/consumer/access"
)

type helperUser struct{}

func (helperUser) GetAuthIdentifierName() string { return "id" }
func (helperUser) GetAuthIdentifier() string     { return "user-1" }
func (helperUser) GetAuthPasswordName() string   { return "password" }
func (helperUser) GetAuthPassword() string       { return "" }
func (helperUser) SetAuthPassword(string)        {}
func (helperUser) GetRememberToken() string      { return "" }
func (helperUser) SetRememberToken(string)       {}
func (helperUser) GetRememberTokenName() string  { return "remember_token" }

func TestAuthorizesRequestsHelpers(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("view", func(context.Context, auth.Authenticatable, ...any) authaccess.Response {
		return authaccess.Allow()
	})

	helper := consumeraccess.AuthorizesRequests{Gate: gate}
	ability, args := helper.ParseAbilityAndArguments("", "show", "account")
	if ability != "view" {
		t.Fatalf("unexpected normalized ability: %q", ability)
	}

	if len(args) != 1 || args[0] != "account" {
		t.Fatalf("unexpected arguments: %#v", args)
	}

	if err := helper.Authorize(context.Background(), helperUser{}, "view"); err != nil {
		t.Fatalf("Authorize: %v", err)
	}

	if helper.ResourceAbilityMap()["destroy"] != "delete" {
		t.Fatal("expected destroy to map to delete")
	}
}
