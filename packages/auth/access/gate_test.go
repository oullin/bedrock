package access_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/auth/access"
	cauth "github.com/bedrock/packages/contracts/auth"
)

type testUser struct {
	id string
}

type testModel struct {
	Name string
}

type testPolicy struct{}

type testPolicyWithBefore struct{}

func (u *testUser) GetAuthIdentifierName() string { return "id" }
func (u *testUser) GetAuthIdentifier() string     { return u.id }
func (u *testUser) GetAuthPasswordName() string   { return "password" }
func (u *testUser) GetAuthPassword() string       { return "" }
func (u *testUser) SetAuthPassword(_ string)      {}
func (u *testUser) GetRememberToken() string      { return "" }
func (u *testUser) SetRememberToken(_ string)     {}
func (u *testUser) GetRememberTokenName() string  { return "remember_token" }

func (p *testPolicy) View(_ context.Context, _ cauth.Authenticatable, _ any) bool {
	return true
}

func (p *testPolicy) Delete(_ context.Context, _ cauth.Authenticatable, _ any) bool {
	return false
}

func (p *testPolicyWithBefore) Before(_ context.Context, _ cauth.Authenticatable, ability string) (bool, bool) {
	if ability == "admin-only" {
		return false, true
	}

	return false, false
}

func (p *testPolicyWithBefore) View(_ context.Context, _ cauth.Authenticatable, _ any) bool {
	return true
}

func userResolver(user cauth.Authenticatable) func(context.Context) cauth.Authenticatable {
	return func(_ context.Context) cauth.Authenticatable { return user }
}

// --- Gate: Has ---

func TestGateHasReturnsTrueForDefinedAbility(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))
	g.Define("edit", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return true, nil
	})

	if !g.Has("edit") {
		t.Error("Has should return true for defined ability")
	}
}

func TestGateHasReturnsFalseForUndefinedAbility(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))

	if g.Has("edit") {
		t.Error("Has should return false for undefined ability")
	}
}

// --- Gate: None ---

func TestGateNoneReturnsTrueWhenAllDenied(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))
	g.Define("edit", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return false, nil
	})
	g.Define("delete", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return false, nil
	})

	if !g.None(context.Background(), []string{"edit", "delete"}, nil) {
		t.Error("None should return true when all abilities are denied")
	}
}

func TestGateNoneReturnsFalseWhenAnyAllowed(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))
	g.Define("edit", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return true, nil
	})
	g.Define("delete", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return false, nil
	})

	if g.None(context.Background(), []string{"edit", "delete"}, nil) {
		t.Error("None should return false when any ability is allowed")
	}
}

// --- Gate: AllowIf / DenyIf ---

func TestAllowIfReturnsAllowWhenTrue(t *testing.T) {
	resp := access.AllowIf(true, "allowed", 0)

	if !resp.Allowed {
		t.Error("AllowIf(true) should return allowed response")
	}
}

func TestAllowIfReturnsDenyWhenFalse(t *testing.T) {
	resp := access.AllowIf(false, "denied", 403)

	if resp.Allowed {
		t.Error("AllowIf(false) should return denied response")
	}

	if resp.StatusCode != 403 {
		t.Errorf("StatusCode = %d, want 403", resp.StatusCode)
	}
}

func TestDenyIfReturnsDenyWhenTrue(t *testing.T) {
	resp := access.DenyIf(true, "denied", 403)

	if resp.Allowed {
		t.Error("DenyIf(true) should return denied response")
	}
}

func TestDenyIfReturnsAllowWhenFalse(t *testing.T) {
	resp := access.DenyIf(false, "ok", 0)

	if !resp.Allowed {
		t.Error("DenyIf(false) should return allowed response")
	}
}

// --- Gate: Abilities ---

func TestGateAbilitiesReturnsDefinedAbilities(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))
	g.Define("edit", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return true, nil
	})
	g.Define("delete", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return false, nil
	})

	abilities := g.Abilities()

	if len(abilities) != 2 {
		t.Fatalf("expected 2 abilities, got %d", len(abilities))
	}

	// Sorted alphabetically.
	if abilities[0] != "delete" || abilities[1] != "edit" {
		t.Errorf("abilities = %v, want [delete, edit]", abilities)
	}
}

func TestGateAbilitiesReturnsEmptyWhenNoneDefined(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))

	if len(g.Abilities()) != 0 {
		t.Error("Abilities should return empty for new gate")
	}
}

// --- Gate: Policy ---

func TestGatePolicyMethodIsCalledForModel(t *testing.T) {
	user := &testUser{id: "1"}
	g := access.New(userResolver(user))
	g.RegisterPolicy("*access_test.testModel", &testPolicy{})

	model := &testModel{Name: "test"}

	if !g.Check(context.Background(), "view", model) {
		t.Error("policy View should allow")
	}

	if g.Check(context.Background(), "delete", model) {
		t.Error("policy Delete should deny")
	}
}

func TestGatePolicyBeforeMethodOverrides(t *testing.T) {
	user := &testUser{id: "1"}
	g := access.New(userResolver(user))
	g.RegisterPolicy("*access_test.testModel", &testPolicyWithBefore{})

	model := &testModel{Name: "test"}

	// "admin-only" is intercepted by Before and denied.
	if g.Check(context.Background(), "admin-only", model) {
		t.Error("Before should deny admin-only")
	}

	// "view" falls through Before and uses the View method.
	if !g.Check(context.Background(), "view", model) {
		t.Error("View should be allowed when Before does not handle")
	}
}

func TestGateGetPolicyForReturnsPolicy(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))
	policy := &testPolicy{}
	g.RegisterPolicy("*access_test.testModel", policy)

	got, ok := g.GetPolicyFor(&testModel{})

	if !ok {
		t.Fatal("GetPolicyFor should return true for registered model type")
	}

	if got != policy {
		t.Error("GetPolicyFor should return the registered policy")
	}
}

func TestGateGetPolicyForReturnsFalseForUnregistered(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))

	_, ok := g.GetPolicyFor(&testModel{})

	if ok {
		t.Error("GetPolicyFor should return false for unregistered type")
	}
}

// --- Gate: ForUser with policies ---

func TestGateForUserSharesPolicies(t *testing.T) {
	user1 := &testUser{id: "1"}
	user2 := &testUser{id: "2"}
	g := access.New(userResolver(user1))
	g.RegisterPolicy("*access_test.testModel", &testPolicy{})

	forked := g.ForUser(user2)

	model := &testModel{Name: "test"}

	if !forked.Check(context.Background(), "view", model) {
		t.Error("ForUser gate should share policies")
	}
}

// --- Gate: Existing tests for Check/Any/Every/Denies/Authorize ---

func TestGateCheckReturnsTrueForAllowedAbility(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))
	g.Define("edit", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return true, nil
	})

	if !g.Check(context.Background(), "edit", nil) {
		t.Error("Check should return true for allowed ability")
	}
}

func TestGateCheckReturnsFalseForDeniedAbility(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))
	g.Define("edit", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return false, nil
	})

	if g.Check(context.Background(), "edit", nil) {
		t.Error("Check should return false for denied ability")
	}
}

func TestGateBeforeHookOverridesResult(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))
	g.Define("edit", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return false, nil
	})
	g.Before(func(_ context.Context, _ cauth.Authenticatable, _ string, _ any) (bool, bool) {
		return true, true
	})

	if !g.Check(context.Background(), "edit", nil) {
		t.Error("Before hook should override to allow")
	}
}

func TestGateAfterHookIsCalled(t *testing.T) {
	called := false
	g := access.New(userResolver(&testUser{id: "1"}))
	g.Define("edit", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return true, nil
	})
	g.After(func(_ context.Context, _ cauth.Authenticatable, _ string, _ bool, _ any) {
		called = true
	})

	g.Check(context.Background(), "edit", nil)

	if !called {
		t.Error("After hook should be called")
	}
}

func TestGateAnyReturnsTrueIfAnyAllowed(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))
	g.Define("edit", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return false, nil
	})
	g.Define("view", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return true, nil
	})

	if !g.Any(context.Background(), []string{"edit", "view"}, nil) {
		t.Error("Any should return true when at least one ability is allowed")
	}
}

func TestGateEveryReturnsFalseIfAnyDenied(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))
	g.Define("edit", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return false, nil
	})
	g.Define("view", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return true, nil
	})

	if g.Every(context.Background(), []string{"edit", "view"}, nil) {
		t.Error("Every should return false when any ability is denied")
	}
}

func TestGateAuthorizeReturnsErrorWhenDenied(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))
	g.Define("edit", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return false, nil
	})

	err := g.Authorize(context.Background(), "edit", nil)

	if err == nil {
		t.Error("Authorize should return error when denied")
	}
}

func TestGateAuthorizeReturnsNilWhenAllowed(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))
	g.Define("edit", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return true, nil
	})

	err := g.Authorize(context.Background(), "edit", nil)

	if err != nil {
		t.Errorf("Authorize should return nil when allowed, got %v", err)
	}
}

func TestGateDeniesReturnsTrueForDeniedAbility(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))
	g.Define("edit", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return false, nil
	})

	if !g.Denies(context.Background(), "edit", nil) {
		t.Error("Denies should return true for denied ability")
	}
}

func TestGateForUserUsesSpecifiedUser(t *testing.T) {
	user1 := &testUser{id: "1"}
	user2 := &testUser{id: "2"}

	g := access.New(userResolver(user1))
	g.Define("edit", func(_ context.Context, u cauth.Authenticatable, _ any) (bool, error) {
		return u.GetAuthIdentifier() == "2", nil
	})

	if g.Check(context.Background(), "edit", nil) {
		t.Error("user1 should not be allowed")
	}

	forked := g.ForUser(user2)

	if !forked.Check(context.Background(), "edit", nil) {
		t.Error("user2 should be allowed via ForUser")
	}
}

func TestGateInspectReturnsResponse(t *testing.T) {
	g := access.New(userResolver(&testUser{id: "1"}))
	g.Define("edit", func(_ context.Context, _ cauth.Authenticatable, _ any) (bool, error) {
		return true, nil
	})

	resp := g.Inspect(context.Background(), "edit", nil)

	if !resp.Allowed {
		t.Error("Inspect should return allowed response")
	}
}

func TestAllowAndDenyResponses(t *testing.T) {
	allow := access.Allow("good")

	if !allow.Allowed {
		t.Error("Allow should set Allowed to true")
	}

	if allow.Message != "good" {
		t.Errorf("Allow.Message = %q, want %q", allow.Message, "good")
	}

	deny := access.Deny("bad", 403)

	if deny.Allowed {
		t.Error("Deny should set Allowed to false")
	}

	if deny.StatusCode != 403 {
		t.Errorf("Deny.StatusCode = %d, want 403", deny.StatusCode)
	}
}
