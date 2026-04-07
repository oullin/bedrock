package access_test

import (
	"context"
	"testing"

	auth "github.com/bedrock/packages/auth"
	authaccess "github.com/bedrock/packages/auth/access"
	consumeraccess "github.com/bedrock/packages/auth/consumer/access"
)

// ---------- test helpers ----------

type helperUser struct{}

func (helperUser) GetAuthIdentifierName() string { return "id" }
func (helperUser) GetAuthIdentifier() string     { return "user-1" }
func (helperUser) GetAuthPasswordName() string   { return "password" }
func (helperUser) GetAuthPassword() string       { return "" }
func (helperUser) SetAuthPassword(string)        {}
func (helperUser) GetRememberToken() string      { return "" }
func (helperUser) SetRememberToken(string)       {}
func (helperUser) GetRememberTokenName() string  { return "remember_token" }

// ======================== TESTS ========================

func TestAuthorizesRequestsAuthorize(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("view", func(context.Context, auth.Authenticatable, ...any) authaccess.Response {
		return authaccess.Allow()
	})

	helper := consumeraccess.AuthorizesRequests{Gate: gate}
	if err := helper.Authorize(context.Background(), helperUser{}, "view"); err != nil {
		t.Fatalf("Authorize: %v", err)
	}
}

func TestAuthorizesRequestsAuthorizeReturnsError(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("edit", func(context.Context, auth.Authenticatable, ...any) authaccess.Response {
		return authaccess.Deny("not allowed")
	})

	helper := consumeraccess.AuthorizesRequests{Gate: gate}
	err := helper.Authorize(context.Background(), helperUser{}, "edit")
	if err == nil {
		t.Fatal("expected error for denied ability")
	}
}

func TestParseAbilityAndArgumentsWithEmptyAbility(t *testing.T) {
	t.Parallel()

	helper := consumeraccess.AuthorizesRequests{}

	// When ability is empty, first argument becomes the ability (normalized)
	ability, args := helper.ParseAbilityAndArguments("", "show", "account")
	if ability != "view" {
		t.Fatalf("expected normalized ability 'view', got %q", ability)
	}
	if len(args) != 1 || args[0] != "account" {
		t.Fatalf("expected args [account], got %v", args)
	}
}

func TestParseAbilityAndArgumentsWithDestroyAction(t *testing.T) {
	t.Parallel()

	helper := consumeraccess.AuthorizesRequests{}

	ability, args := helper.ParseAbilityAndArguments("", "destroy", "resource")
	if ability != "delete" {
		t.Fatalf("expected normalized ability 'delete', got %q", ability)
	}
	if len(args) != 1 || args[0] != "resource" {
		t.Fatalf("expected args [resource], got %v", args)
	}
}

func TestParseAbilityAndArgumentsWithExplicitAbility(t *testing.T) {
	t.Parallel()

	helper := consumeraccess.AuthorizesRequests{}

	ability, args := helper.ParseAbilityAndArguments("create", "arg1")
	if ability != "create" {
		t.Fatalf("expected ability 'create', got %q", ability)
	}
	if len(args) != 1 || args[0] != "arg1" {
		t.Fatalf("expected args [arg1], got %v", args)
	}
}

func TestParseAbilityAndArgumentsWithNoArguments(t *testing.T) {
	t.Parallel()

	helper := consumeraccess.AuthorizesRequests{}

	ability, args := helper.ParseAbilityAndArguments("edit")
	if ability != "edit" {
		t.Fatalf("expected ability 'edit', got %q", ability)
	}
	if len(args) != 0 {
		t.Fatalf("expected no args, got %v", args)
	}
}

func TestResourceAbilityMap(t *testing.T) {
	t.Parallel()

	helper := consumeraccess.AuthorizesRequests{}
	m := helper.ResourceAbilityMap()

	expected := map[string]string{
		"index":   "viewAny",
		"show":    "view",
		"create":  "create",
		"store":   "create",
		"edit":    "update",
		"update":  "update",
		"destroy": "delete",
	}

	for k, v := range expected {
		if m[k] != v {
			t.Fatalf("expected %q -> %q, got %q", k, v, m[k])
		}
	}
}

func TestResourceMethodsWithoutModels(t *testing.T) {
	t.Parallel()

	helper := consumeraccess.AuthorizesRequests{}
	methods := helper.ResourceMethodsWithoutModels()

	expected := map[string]bool{"index": true, "create": true, "store": true}
	for _, m := range methods {
		if !expected[m] {
			t.Fatalf("unexpected method without model: %q", m)
		}
	}
	if len(methods) != 3 {
		t.Fatalf("expected 3 methods, got %d", len(methods))
	}
}
