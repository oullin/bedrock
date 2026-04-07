package access_test

import (
	"context"
	"errors"
	"testing"

	auth "github.com/bedrock/packages/anvil/auth"
	authaccess "github.com/bedrock/packages/anvil/auth/access"
)

// ---------- test helpers ----------

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

type post struct {
	AuthorID string
}

var ctx = context.Background()

// ---------- Upstream: testBasicClosuresCanBeDefined ----------

func TestBasicClosuresCanBeDefined(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("open-gate", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})
	gate.Define("closed-gate", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Deny("denied")
	})

	user := gateUser{id: "u1"}

	if !gate.Check(ctx, user, "open-gate") {
		t.Fatal("expected open-gate to allow")
	}
	if gate.Check(ctx, user, "closed-gate") {
		t.Fatal("expected closed-gate to deny")
	}
}

// ---------- Upstream: testBeforeCallbacksCanOverrideResultIfNecessary ----------

func TestBeforeCallbacksCanOverrideResultIfNecessary(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("open-gate", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})
	gate.Before(func(_ context.Context, _ auth.Authenticatable, _ string, _ ...any) (authaccess.Response, bool) {
		return authaccess.Deny("blocked by before"), true
	})

	user := gateUser{id: "u1"}
	if gate.Check(ctx, user, "open-gate") {
		t.Fatal("expected before callback to override and deny")
	}
}

// ---------- Upstream: testBeforeCallbacksDontInterruptGateCheckIfNoValueIsReturned ----------

func TestBeforeCallbacksDontInterruptGateCheckIfNoValueIsReturned(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("open-gate", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})
	gate.Before(func(_ context.Context, _ auth.Authenticatable, _ string, _ ...any) (authaccess.Response, bool) {
		return authaccess.Response{}, false // not handled
	})

	user := gateUser{id: "u1"}
	if !gate.Check(ctx, user, "open-gate") {
		t.Fatal("expected before callback to pass through when not handled")
	}
}

// ---------- Upstream: testAfterCallbacksAreCalledWithResult ----------

func TestAfterCallbacksAreCalledWithResult(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("view", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})

	var receivedAllowed bool
	gate.After(func(_ context.Context, _ auth.Authenticatable, _ string, result authaccess.Response, _ ...any) authaccess.Response {
		receivedAllowed = result.Allowed
		return result
	})

	user := gateUser{id: "u1"}
	gate.Check(ctx, user, "view")

	if !receivedAllowed {
		t.Fatal("expected after callback to receive allowed result")
	}
}

// ---------- Upstream: testAfterCallbacksCanAllowIfNull ----------

func TestAfterCallbacksCanAllowUndefined(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	// No ability defined — the gate would deny by default.
	gate.After(func(_ context.Context, _ auth.Authenticatable, _ string, result authaccess.Response, _ ...any) authaccess.Response {
		if !result.Allowed {
			return authaccess.Allow() // override the denial
		}
		return result
	})

	user := gateUser{id: "u1"}
	if !gate.Check(ctx, user, "undefined-ability") {
		t.Fatal("expected after callback to override undefined ability to allow")
	}
}

// ---------- Upstream: testAfterCallbacksDoNotOverridePreviousResult ----------

func TestAfterCallbacksDoNotOverridePreviousResult(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("deny-gate", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Deny("nope")
	})
	gate.After(func(_ context.Context, _ auth.Authenticatable, _ string, result authaccess.Response, _ ...any) authaccess.Response {
		// Return the result as-is; don't override
		return result
	})

	user := gateUser{id: "u1"}
	if gate.Check(ctx, user, "deny-gate") {
		t.Fatal("expected after callback not to override denied result")
	}
}

// ---------- Upstream: testAfterCallbacksDoNotOverrideEachOther ----------

func TestAfterCallbacksDoNotOverrideEachOther(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("view", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})

	var first, second bool
	gate.After(func(_ context.Context, _ auth.Authenticatable, _ string, result authaccess.Response, _ ...any) authaccess.Response {
		first = true
		return result
	})
	gate.After(func(_ context.Context, _ auth.Authenticatable, _ string, result authaccess.Response, _ ...any) authaccess.Response {
		second = true
		return result
	})

	user := gateUser{id: "u1"}
	gate.Check(ctx, user, "view")

	if !first || !second {
		t.Fatal("expected both after callbacks to run")
	}
}

// ---------- Upstream: testCurrentUserThatIsOnGateAlwaysInjectedIntoClosureCallbacks ----------

func TestCurrentUserIsInjectedIntoClosureCallbacks(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	var receivedID string
	gate.Define("check-user", func(_ context.Context, user auth.Authenticatable, _ ...any) authaccess.Response {
		receivedID = user.GetAuthIdentifier()
		return authaccess.Allow()
	})

	user := gateUser{id: "user-42"}
	gate.Check(ctx, user, "check-user")

	if receivedID != "user-42" {
		t.Fatalf("expected user id 'user-42', got %q", receivedID)
	}
}

// ---------- Upstream: testASingleArgumentCanBePassedWhenCheckingAbilities ----------

func TestSingleArgumentCanBePassedWhenCheckingAbilities(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("view-account", func(_ context.Context, user auth.Authenticatable, args ...any) authaccess.Response {
		acct := args[0].(account)
		if user.GetAuthIdentifier() == acct.OwnerID {
			return authaccess.Allow()
		}
		return authaccess.Deny("not owner")
	})

	user := gateUser{id: "u1"}
	if !gate.Check(ctx, user, "view-account", account{OwnerID: "u1"}) {
		t.Fatal("expected owner to be allowed")
	}
	if gate.Check(ctx, user, "view-account", account{OwnerID: "u2"}) {
		t.Fatal("expected non-owner to be denied")
	}
}

// ---------- Upstream: testMultipleArgumentsCanBePassedWhenCheckingAbilities ----------

func TestMultipleArgumentsCanBePassedWhenCheckingAbilities(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("multi-arg", func(_ context.Context, _ auth.Authenticatable, args ...any) authaccess.Response {
		if len(args) == 3 && args[0] == "a" && args[1] == "b" && args[2] == "c" {
			return authaccess.Allow()
		}
		return authaccess.Deny("wrong args")
	})

	user := gateUser{id: "u1"}
	if !gate.Check(ctx, user, "multi-arg", "a", "b", "c") {
		t.Fatal("expected multiple arguments to be passed correctly")
	}
}

// ---------- Upstream: testPolicyClassesCanBeDefinedToHandleChecksForGivenType ----------

func TestPolicyClassesCanBeDefinedToHandleChecksForGivenType(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Policy(account{}, func(_ context.Context, user auth.Authenticatable, ability string, arguments ...any) (authaccess.Response, bool) {
		acct := arguments[0].(account)
		if ability == "view" && user.GetAuthIdentifier() == acct.OwnerID {
			return authaccess.Allow(), true
		}
		return authaccess.Deny("policy denied"), true
	})

	user := gateUser{id: "u1"}
	if !gate.Check(ctx, user, "view", account{OwnerID: "u1"}) {
		t.Fatal("expected policy to allow owner to view")
	}
	if gate.Check(ctx, user, "view", account{OwnerID: "other"}) {
		t.Fatal("expected policy to deny non-owner")
	}
}

// ---------- Upstream: testPolicyDefaultToFalseIfMethodDoesNotExistAndGateDoesNotExist ----------

func TestPolicyDefaultToFalseIfMethodDoesNotExistAndGateDoesNotExist(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	// No abilities or policies defined.

	user := gateUser{id: "u1"}
	if gate.Check(ctx, user, "nonexistent-ability") {
		t.Fatal("expected undefined ability to default to deny")
	}
}

// ---------- Upstream: testPoliciesAlwaysOverrideClosuresWithSameName ----------

func TestPoliciesOverrideClosuresForSameResource(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()

	// Define a closure that would allow
	gate.Define("view", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})

	// Register a policy that denies
	gate.Policy(account{}, func(_ context.Context, _ auth.Authenticatable, ability string, _ ...any) (authaccess.Response, bool) {
		if ability == "view" {
			return authaccess.Deny("policy overrides"), true
		}
		return authaccess.Response{}, false
	})

	user := gateUser{id: "u1"}
	// When checking with a resource argument, policy should take precedence
	if gate.Check(ctx, user, "view", account{OwnerID: "u1"}) {
		t.Fatal("expected policy to override the closure and deny")
	}
}

// ---------- Upstream: testPoliciesDeferToGatesIfMethodDoesNotExist ----------

func TestPoliciesDeferToGatesIfMethodDoesNotExist(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("special", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})

	// Policy only handles "view", not "special"
	gate.Policy(account{}, func(_ context.Context, _ auth.Authenticatable, ability string, _ ...any) (authaccess.Response, bool) {
		if ability == "view" {
			return authaccess.Allow(), true
		}
		return authaccess.Response{}, false // not handled → defer to gate
	})

	user := gateUser{id: "u1"}
	if !gate.Check(ctx, user, "special", account{OwnerID: "u1"}) {
		t.Fatal("expected gate closure to handle when policy defers")
	}
}

// ---------- Upstream: testAuthorizeThrowsUnauthorizedException ----------

func TestAuthorizeThrowsUnauthorizedException(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("denied", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Deny("not allowed")
	})

	user := gateUser{id: "u1"}
	err := gate.Authorize(ctx, user, "denied")
	if err == nil {
		t.Fatal("expected Authorize to return an error")
	}

	var authErr authaccess.AuthorizationException
	if !errors.As(err, &authErr) {
		t.Fatalf("expected AuthorizationException, got %T: %v", err, err)
	}
	if authErr.Ability != "denied" {
		t.Fatalf("expected ability 'denied', got %q", authErr.Ability)
	}
	if authErr.Message != "not allowed" {
		t.Fatalf("expected message 'not allowed', got %q", authErr.Message)
	}
}

// ---------- Upstream: testAuthorizeReturnsAllowedResponse ----------

func TestAuthorizeReturnsNilForAllowedAbility(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("allowed", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})

	user := gateUser{id: "u1"}
	if err := gate.Authorize(ctx, user, "allowed"); err != nil {
		t.Fatalf("expected Authorize to succeed, got %v", err)
	}
}

// ---------- Upstream: testResponseReturnsResponseWhenAbilityGranted ----------

func TestInspectReturnsResponseWhenAbilityGranted(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("view", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})

	user := gateUser{id: "u1"}
	response := gate.Inspect(ctx, user, "view")

	if !response.Allowed {
		t.Fatal("expected response to be allowed")
	}
	if response.Message != "" {
		t.Fatalf("expected empty message for allowed response, got %q", response.Message)
	}
}

// ---------- Upstream: testResponseReturnsResponseWhenAbilityDenied ----------

func TestInspectReturnsResponseWhenAbilityDenied(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("edit", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Deny("you cannot edit")
	})

	user := gateUser{id: "u1"}
	response := gate.Inspect(ctx, user, "edit")

	if response.Allowed {
		t.Fatal("expected response to be denied")
	}
	if response.Message != "you cannot edit" {
		t.Fatalf("expected message 'you cannot edit', got %q", response.Message)
	}
}

// ---------- Upstream: testAnyAbilityCheckPassesIfAllPass ----------

func TestAnyAbilityCheckPassesIfAllPass(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("a", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})
	gate.Define("b", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})

	user := gateUser{id: "u1"}
	if !gate.Any(ctx, user, []string{"a", "b"}) {
		t.Fatal("expected Any to pass when all abilities pass")
	}
}

// ---------- Upstream: testAnyAbilityCheckPassesIfAtLeastOnePasses ----------

func TestAnyAbilityCheckPassesIfAtLeastOnePasses(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("a", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})
	gate.Define("b", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Deny("denied")
	})

	user := gateUser{id: "u1"}
	if !gate.Any(ctx, user, []string{"a", "b"}) {
		t.Fatal("expected Any to pass when at least one ability passes")
	}
}

// ---------- Upstream: testAnyAbilityCheckFailsIfNonePass ----------

func TestAnyAbilityCheckFailsIfNonePass(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("a", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Deny("denied")
	})
	gate.Define("b", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Deny("denied")
	})

	user := gateUser{id: "u1"}
	if gate.Any(ctx, user, []string{"a", "b"}) {
		t.Fatal("expected Any to fail when none pass")
	}
}

// ---------- Upstream: testDenies ----------

func TestDeniesMethod(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("open", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})
	gate.Define("closed", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Deny("no")
	})

	user := gateUser{id: "u1"}
	if gate.Denies(ctx, user, "open") {
		t.Fatal("expected Denies to return false for allowed ability")
	}
	if !gate.Denies(ctx, user, "closed") {
		t.Fatal("expected Denies to return true for denied ability")
	}
}

// ---------- Upstream: testPoliciesMayHaveBeforeMethodsToOverrideChecks ----------

func TestPoliciesBeforeHooksOverrideChecks(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Before(func(_ context.Context, user auth.Authenticatable, ability string, _ ...any) (authaccess.Response, bool) {
		if user.GetAuthIdentifier() == "admin" {
			return authaccess.Allow(), true // admin bypasses all
		}
		return authaccess.Response{}, false
	})

	gate.Define("manage", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Deny("only admin")
	})

	admin := gateUser{id: "admin"}
	regular := gateUser{id: "user"}

	if !gate.Check(ctx, admin, "manage") {
		t.Fatal("expected admin to be allowed by before hook")
	}
	if gate.Check(ctx, regular, "manage") {
		t.Fatal("expected regular user to be denied")
	}
}

// ---------- Upstream: testAuthorizationExceptionError ----------

func TestAuthorizationExceptionErrorString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		exc      authaccess.AuthorizationException
		expected string
	}{
		{
			name:     "with message",
			exc:      authaccess.AuthorizationException{Ability: "edit", Message: "custom message"},
			expected: "custom message",
		},
		{
			name:     "with ability only",
			exc:      authaccess.AuthorizationException{Ability: "edit"},
			expected: `auth: authorization denied for "edit"`,
		},
		{
			name:     "empty",
			exc:      authaccess.AuthorizationException{},
			expected: "auth: authorization denied",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.exc.Error(); got != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

// ---------- Upstream: testAuthorizeReturnsAnAllowedResponseForATruthyReturn ----------

func TestAuthorizeReturnsNilForTruthyPolicyReturn(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Policy(account{}, func(_ context.Context, _ auth.Authenticatable, ability string, _ ...any) (authaccess.Response, bool) {
		if ability == "view" {
			return authaccess.Allow(), true
		}
		return authaccess.Response{}, false
	})

	user := gateUser{id: "u1"}
	if err := gate.Authorize(ctx, user, "view", account{OwnerID: "u1"}); err != nil {
		t.Fatalf("expected Authorize to succeed, got %v", err)
	}
}

// ---------- Upstream: testAuthorizeWithPolicyThatReturnsDeniedResponseObjectThrowsException ----------

func TestAuthorizeWithPolicyDeniedResponseThrowsException(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Policy(account{}, func(_ context.Context, _ auth.Authenticatable, ability string, _ ...any) (authaccess.Response, bool) {
		if ability == "delete" {
			return authaccess.Deny("cannot delete"), true
		}
		return authaccess.Response{}, false
	})

	user := gateUser{id: "u1"}
	err := gate.Authorize(ctx, user, "delete", account{OwnerID: "u1"})
	if err == nil {
		t.Fatal("expected Authorize to return an error")
	}

	var authErr authaccess.AuthorizationException
	if !errors.As(err, &authErr) {
		t.Fatalf("expected AuthorizationException, got %T", err)
	}
	if authErr.Message != "cannot delete" {
		t.Fatalf("expected message 'cannot delete', got %q", authErr.Message)
	}
}

// ---------- Multiple policies for different types ----------

func TestMultiplePoliciesForDifferentTypes(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Policy(account{}, func(_ context.Context, user auth.Authenticatable, ability string, arguments ...any) (authaccess.Response, bool) {
		acct := arguments[0].(account)
		if ability == "view" && user.GetAuthIdentifier() == acct.OwnerID {
			return authaccess.Allow(), true
		}
		return authaccess.Deny("account denied"), true
	})
	gate.Policy(post{}, func(_ context.Context, user auth.Authenticatable, ability string, arguments ...any) (authaccess.Response, bool) {
		p := arguments[0].(post)
		if ability == "edit" && user.GetAuthIdentifier() == p.AuthorID {
			return authaccess.Allow(), true
		}
		return authaccess.Deny("post denied"), true
	})

	user := gateUser{id: "u1"}
	if !gate.Check(ctx, user, "view", account{OwnerID: "u1"}) {
		t.Fatal("expected account policy to allow owner")
	}
	if gate.Check(ctx, user, "view", account{OwnerID: "u2"}) {
		t.Fatal("expected account policy to deny non-owner")
	}
	if !gate.Check(ctx, user, "edit", post{AuthorID: "u1"}) {
		t.Fatal("expected post policy to allow author")
	}
	if gate.Check(ctx, user, "edit", post{AuthorID: "u2"}) {
		t.Fatal("expected post policy to deny non-author")
	}
}

// ---------- Before and After ordering ----------

func TestBeforeAndAfterCallbackOrdering(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	var order []string

	gate.Before(func(_ context.Context, _ auth.Authenticatable, _ string, _ ...any) (authaccess.Response, bool) {
		order = append(order, "before1")
		return authaccess.Response{}, false
	})
	gate.Before(func(_ context.Context, _ auth.Authenticatable, _ string, _ ...any) (authaccess.Response, bool) {
		order = append(order, "before2")
		return authaccess.Response{}, false
	})
	gate.Define("test", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		order = append(order, "ability")
		return authaccess.Allow()
	})
	gate.After(func(_ context.Context, _ auth.Authenticatable, _ string, result authaccess.Response, _ ...any) authaccess.Response {
		order = append(order, "after1")
		return result
	})
	gate.After(func(_ context.Context, _ auth.Authenticatable, _ string, result authaccess.Response, _ ...any) authaccess.Response {
		order = append(order, "after2")
		return result
	})

	user := gateUser{id: "u1"}
	gate.Check(ctx, user, "test")

	expected := []string{"before1", "before2", "ability", "after1", "after2"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d callbacks, got %d: %v", len(expected), len(order), order)
	}
	for i, v := range expected {
		if order[i] != v {
			t.Fatalf("expected order[%d] = %q, got %q", i, v, order[i])
		}
	}
}

// ---------- Before callback short-circuits remaining befores ----------

func TestBeforeCallbackShortCircuits(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	var secondCalled bool

	gate.Before(func(_ context.Context, _ auth.Authenticatable, _ string, _ ...any) (authaccess.Response, bool) {
		return authaccess.Allow(), true // short-circuit
	})
	gate.Before(func(_ context.Context, _ auth.Authenticatable, _ string, _ ...any) (authaccess.Response, bool) {
		secondCalled = true
		return authaccess.Response{}, false
	})
	gate.Define("test", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		t.Fatal("ability should not be called when before short-circuits")
		return authaccess.Deny("should not reach")
	})

	user := gateUser{id: "u1"}
	if !gate.Check(ctx, user, "test") {
		t.Fatal("expected before callback to allow")
	}
	if secondCalled {
		t.Fatal("expected second before callback to not be called")
	}
}

// ---------- After callback can modify the result ----------

func TestAfterCallbackCanModifyResult(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("test", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Deny("initially denied")
	})
	gate.After(func(_ context.Context, _ auth.Authenticatable, _ string, _ authaccess.Response, _ ...any) authaccess.Response {
		return authaccess.Allow() // override to allow
	})

	user := gateUser{id: "u1"}
	if !gate.Check(ctx, user, "test") {
		t.Fatal("expected after callback to override denial to allow")
	}
}

// ---------- Policy with pointer receiver ----------

func TestPolicyWithPointerToResource(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Policy(&account{}, func(_ context.Context, user auth.Authenticatable, ability string, arguments ...any) (authaccess.Response, bool) {
		acct := arguments[0].(*account)
		if ability == "view" && user.GetAuthIdentifier() == acct.OwnerID {
			return authaccess.Allow(), true
		}
		return authaccess.Deny("denied"), true
	})

	user := gateUser{id: "u1"}
	if !gate.Check(ctx, user, "view", &account{OwnerID: "u1"}) {
		t.Fatal("expected pointer policy to allow owner")
	}
}

// ---------- Undefined ability returns deny ----------

func TestUndefinedAbilityReturnsDeny(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	user := gateUser{id: "u1"}

	response := gate.Inspect(ctx, user, "nonexistent")
	if response.Allowed {
		t.Fatal("expected undefined ability to return deny")
	}
	if response.Message != "auth: ability is not defined" {
		t.Fatalf("expected 'auth: ability is not defined', got %q", response.Message)
	}
}

// ---------- Multiple abilities on same gate (isolation) ----------

func TestMultipleAbilitiesOnSameGate(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("read", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})
	gate.Define("write", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Deny("no write")
	})
	gate.Define("admin", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})

	user := gateUser{id: "u1"}
	if !gate.Check(ctx, user, "read") {
		t.Fatal("expected read to be allowed")
	}
	if gate.Check(ctx, user, "write") {
		t.Fatal("expected write to be denied")
	}
	if !gate.Check(ctx, user, "admin") {
		t.Fatal("expected admin to be allowed")
	}
}

// ---------- Allow and Deny constructors ----------

func TestAllowAndDenyConstructors(t *testing.T) {
	t.Parallel()

	allow := authaccess.Allow()
	if !allow.Allowed {
		t.Fatal("Allow() should return Allowed=true")
	}
	if allow.Message != "" {
		t.Fatal("Allow() should have empty message")
	}

	deny := authaccess.Deny("reason")
	if deny.Allowed {
		t.Fatal("Deny() should return Allowed=false")
	}
	if deny.Message != "reason" {
		t.Fatalf("expected message 'reason', got %q", deny.Message)
	}

	denyEmpty := authaccess.Deny("")
	if denyEmpty.Message != "" {
		t.Fatal("Deny('') should have empty message")
	}
}

// ---------- Authorizer interface compliance ----------

func TestGateImplementsAuthorizer(t *testing.T) {
	t.Parallel()

	var _ authaccess.Authorizer = authaccess.NewGate()
}

// ---------- Ability name trimming ----------

func TestAbilityNameTrimming(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("  spaced  ", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})

	user := gateUser{id: "u1"}
	if !gate.Check(ctx, user, "spaced") {
		t.Fatal("expected trimmed ability name to match")
	}
	if !gate.Check(ctx, user, "  spaced  ") {
		t.Fatal("expected untrimmed check to match trimmed definition")
	}
}

// ---------- Before callback allows even when ability is undefined ----------

func TestBeforeAllowsUndefinedAbility(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Before(func(_ context.Context, user auth.Authenticatable, _ string, _ ...any) (authaccess.Response, bool) {
		if user.GetAuthIdentifier() == "super" {
			return authaccess.Allow(), true
		}
		return authaccess.Response{}, false
	})

	super := gateUser{id: "super"}
	if !gate.Check(ctx, super, "anything") {
		t.Fatal("expected before to allow even for undefined ability")
	}

	regular := gateUser{id: "regular"}
	if gate.Check(ctx, regular, "anything") {
		t.Fatal("expected regular user to be denied for undefined ability")
	}
}

// ---------- Context passed through to all callbacks ----------

func TestContextPassedThroughToCallbacks(t *testing.T) {
	t.Parallel()

	type ctxKey string
	gate := authaccess.NewGate()

	gate.Before(func(c context.Context, _ auth.Authenticatable, _ string, _ ...any) (authaccess.Response, bool) {
		if c.Value(ctxKey("key")) != "value" {
			t.Fatal("expected context to contain key in before callback")
		}
		return authaccess.Response{}, false
	})
	gate.Define("test", func(c context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		if c.Value(ctxKey("key")) != "value" {
			t.Fatal("expected context to contain key in ability callback")
		}
		return authaccess.Allow()
	})
	gate.After(func(c context.Context, _ auth.Authenticatable, _ string, result authaccess.Response, _ ...any) authaccess.Response {
		if c.Value(ctxKey("key")) != "value" {
			t.Fatal("expected context to contain key in after callback")
		}
		return result
	})

	user := gateUser{id: "u1"}
	c := context.WithValue(ctx, ctxKey("key"), "value")
	gate.Check(c, user, "test")
}

// ---------- Concurrent access is safe ----------

func TestGateConcurrentAccess(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("concurrent", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})

	user := gateUser{id: "u1"}
	done := make(chan bool, 100)

	for i := 0; i < 100; i++ {
		go func() {
			gate.Check(ctx, user, "concurrent")
			done <- true
		}()
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}

// ---------- Policy key resolution for string target ----------

func TestPolicyKeyResolutionForStringTarget(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Policy("custom-resource", func(_ context.Context, _ auth.Authenticatable, ability string, _ ...any) (authaccess.Response, bool) {
		if ability == "view" {
			return authaccess.Allow(), true
		}
		return authaccess.Response{}, false
	})

	user := gateUser{id: "u1"}
	if !gate.Check(ctx, user, "view", "custom-resource") {
		t.Fatal("expected string policy target to work")
	}
}

// ---------- Any with empty abilities ----------

func TestAnyWithEmptyAbilities(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	user := gateUser{id: "u1"}

	if gate.Any(ctx, user, []string{}) {
		t.Fatal("expected Any with empty list to return false")
	}
}

// ---------- Authorize with undefined ability ----------

func TestAuthorizeWithUndefinedAbility(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	user := gateUser{id: "u1"}

	err := gate.Authorize(ctx, user, "nonexistent")
	if err == nil {
		t.Fatal("expected Authorize for undefined ability to fail")
	}

	var authErr authaccess.AuthorizationException
	if !errors.As(err, &authErr) {
		t.Fatalf("expected AuthorizationException, got %T", err)
	}
}

// ---------- Before callback receives correct ability name ----------

func TestBeforeCallbackReceivesAbilityName(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	var receivedAbility string

	gate.Before(func(_ context.Context, _ auth.Authenticatable, ability string, _ ...any) (authaccess.Response, bool) {
		receivedAbility = ability
		return authaccess.Response{}, false
	})
	gate.Define("specific-ability", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})

	user := gateUser{id: "u1"}
	gate.Check(ctx, user, "specific-ability")

	if receivedAbility != "specific-ability" {
		t.Fatalf("expected before callback to receive 'specific-ability', got %q", receivedAbility)
	}
}

// ---------- After callback receives correct ability and arguments ----------

func TestAfterCallbackReceivesAbilityAndArguments(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	var receivedAbility string
	var receivedArgs []any

	gate.Define("view", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})
	gate.After(func(_ context.Context, _ auth.Authenticatable, ability string, result authaccess.Response, args ...any) authaccess.Response {
		receivedAbility = ability
		receivedArgs = args
		return result
	})

	user := gateUser{id: "u1"}
	gate.Check(ctx, user, "view", "arg1", "arg2")

	if receivedAbility != "view" {
		t.Fatalf("expected ability 'view', got %q", receivedAbility)
	}
	if len(receivedArgs) != 2 || receivedArgs[0] != "arg1" || receivedArgs[1] != "arg2" {
		t.Fatalf("expected args ['arg1', 'arg2'], got %v", receivedArgs)
	}
}

// ---------- Policy and ability for same ability name ----------

func TestPolicyAndAbilityCoexist(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()

	gate.Define("view", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Deny("ability says no")
	})

	gate.Policy(account{}, func(_ context.Context, _ auth.Authenticatable, ability string, _ ...any) (authaccess.Response, bool) {
		if ability == "view" {
			return authaccess.Allow(), true
		}
		return authaccess.Response{}, false
	})

	user := gateUser{id: "u1"}

	// With resource argument → policy wins
	if !gate.Check(ctx, user, "view", account{OwnerID: "u1"}) {
		t.Fatal("expected policy to win when resource argument provided")
	}

	// Without resource argument → ability used
	if gate.Check(ctx, user, "view") {
		t.Fatal("expected ability closure to deny when no resource argument")
	}
}

// ---------- Before allows, after skipped (GAP: Upstream runs after even on before short-circuit) ----------

// GAP: In Upstream, after callbacks run even when a before callback short-circuits.
// In Bedrock's Go implementation, after callbacks are skipped when before returns handled=true.
// This is a known behavioral difference documented for future alignment.
func TestAfterSkippedWhenBeforeShortCircuits(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	afterCalled := false

	gate.Before(func(_ context.Context, _ auth.Authenticatable, _ string, _ ...any) (authaccess.Response, bool) {
		return authaccess.Allow(), true
	})
	gate.After(func(_ context.Context, _ auth.Authenticatable, _ string, result authaccess.Response, _ ...any) authaccess.Response {
		afterCalled = true
		return result
	})

	user := gateUser{id: "u1"}
	result := gate.Check(ctx, user, "anything")

	if !result {
		t.Fatal("expected before callback to allow")
	}
	// NOTE: In Upstream, afterCalled would be true.
	// In Bedrock, before short-circuits completely, skipping after callbacks.
	if afterCalled {
		t.Fatal("in current implementation, after should NOT run when before short-circuits")
	}
}

// ---------- Upstream: testAllowIfRespondsTrue / testAllowIfRespondsFalse ----------

func TestAllowIfAndDenyIf(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()

	t.Run("AllowIf true returns allowed", func(t *testing.T) {
		resp := gate.AllowIf(true, "should not see this")
		if !resp.Allowed {
			t.Fatal("expected AllowIf(true) to allow")
		}
	})

	t.Run("AllowIf false returns denied with message", func(t *testing.T) {
		resp := gate.AllowIf(false, "nope")
		if resp.Allowed {
			t.Fatal("expected AllowIf(false) to deny")
		}
		if resp.Message != "nope" {
			t.Fatalf("expected message 'nope', got %q", resp.Message)
		}
	})

	t.Run("DenyIf true returns denied with message", func(t *testing.T) {
		resp := gate.DenyIf(true, "blocked")
		if resp.Allowed {
			t.Fatal("expected DenyIf(true) to deny")
		}
		if resp.Message != "blocked" {
			t.Fatalf("expected message 'blocked', got %q", resp.Message)
		}
	})

	t.Run("DenyIf false returns allowed", func(t *testing.T) {
		resp := gate.DenyIf(false, "should not see this")
		if !resp.Allowed {
			t.Fatal("expected DenyIf(false) to allow")
		}
	})
}

// ---------- Upstream: testGateHasAbility ----------

func TestGateHas(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("edit", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})

	if !gate.Has("edit") {
		t.Fatal("expected Has('edit') to be true")
	}
	if gate.Has("delete") {
		t.Fatal("expected Has('delete') to be false")
	}
	if !gate.Has("  edit  ") {
		t.Fatal("expected Has with extra whitespace to match after trimming")
	}
}

// ---------- Upstream: testNoneAbilityCheckPassesIfNonePass ----------

func TestNoneAbilityCheck(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Define("open", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})
	gate.Define("closed", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Deny("no")
	})

	user := gateUser{id: "u1"}

	if gate.None(ctx, user, []string{"closed"}) != true {
		t.Fatal("expected None to return true when all abilities are denied")
	}
	if gate.None(ctx, user, []string{"open", "closed"}) != false {
		t.Fatal("expected None to return false when at least one ability is allowed")
	}
	if gate.None(ctx, user, []string{"open"}) != false {
		t.Fatal("expected None to return false when ability is allowed")
	}
	if gate.None(ctx, user, []string{}) != true {
		t.Fatal("expected None with empty abilities to return true")
	}
}

// ---------- Upstream: testResponseReturnsWithCode ----------

func TestResponseWithCode(t *testing.T) {
	t.Parallel()

	resp := authaccess.DenyWithCode("forbidden", 403, 403)
	if resp.Allowed {
		t.Fatal("expected DenyWithCode to deny")
	}
	if resp.Code != 403 {
		t.Fatalf("expected Code 403, got %d", resp.Code)
	}
	if resp.Status != 403 {
		t.Fatalf("expected Status 403, got %d", resp.Status)
	}
	if resp.Message != "forbidden" {
		t.Fatalf("expected message 'forbidden', got %q", resp.Message)
	}

	allowed := authaccess.Allow()
	if allowed.Code != 0 || allowed.Status != 0 {
		t.Fatal("expected Allow() to have zero Code and Status")
	}
}

// ---------- Upstream: testForUser ----------

func TestGateForUser(t *testing.T) {
	t.Parallel()

	parent := authaccess.NewGate()
	parent.Define("edit", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})
	parent.Before(func(_ context.Context, _ auth.Authenticatable, _ string, _ ...any) (authaccess.Response, bool) {
		return authaccess.Response{}, false
	})

	child := parent.ForUser()

	user := gateUser{id: "u1"}
	if !child.Check(ctx, user, "edit") {
		t.Fatal("expected child gate to inherit abilities from parent")
	}

	// Mutating child does not affect parent.
	child.Define("delete", func(_ context.Context, _ auth.Authenticatable, _ ...any) authaccess.Response {
		return authaccess.Allow()
	})

	if parent.Has("delete") {
		t.Fatal("expected parent gate to be unaffected by child mutation")
	}
}

// ---------- Upstream: testResourceGates ----------

func TestGateResource(t *testing.T) {
	t.Parallel()

	gate := authaccess.NewGate()
	gate.Resource("posts", func(_ context.Context, user auth.Authenticatable, ability string, _ ...any) (authaccess.Response, bool) {
		if ability == "create" && user.GetAuthIdentifier() == "admin" {
			return authaccess.Allow(), true
		}
		return authaccess.Deny("denied"), true
	})

	admin := gateUser{id: "admin"}
	reader := gateUser{id: "reader"}

	if !gate.Check(ctx, admin, "posts.create") {
		t.Fatal("expected admin to be allowed posts.create")
	}
	if gate.Check(ctx, reader, "posts.create") {
		t.Fatal("expected reader to be denied posts.create")
	}

	// All CRUD abilities should be registered.
	for _, ability := range []string{"posts.viewAny", "posts.view", "posts.create", "posts.update", "posts.delete"} {
		if !gate.Has(ability) {
			t.Fatalf("expected resource ability %q to be registered", ability)
		}
	}
}
