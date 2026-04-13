package access

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"sync"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// Ability is a function that determines whether a user can perform an action.
// It receives the user, the optional model, and returns (allow bool, err error).
type Ability func(ctx context.Context, user cauth.Authenticatable, model any) (bool, error)

// Policy is a struct implementing named ability methods.
// Methods must have the signature: func(ctx, user, model) (bool, error) or func(ctx, user) (bool, error).
type Policy any

// Response represents the result of a gate check.
type Response struct {
	Allowed    bool
	Message    string
	StatusCode int
}

// Allow returns an allowing Response.

// Deny returns a denying Response.

// AuthorizationException is returned when a gate check fails.
type AuthorizationException struct {
	Response Response
}

// Gate manages abilities and policies for authorization.
type Gate struct {
	mu           sync.RWMutex
	abilities    map[string]Ability
	policies     map[string]any
	before       []func(ctx context.Context, user cauth.Authenticatable, ability string, model any) (bool, bool)
	after        []func(ctx context.Context, user cauth.Authenticatable, ability string, result bool, model any)
	userResolver func(ctx context.Context) cauth.Authenticatable
}

func Allow(message string) Response {
	return Response{Allowed: true, Message: message, StatusCode: http.StatusOK}
}

func Deny(message string, statusCode int) Response {
	if statusCode == 0 {
		statusCode = http.StatusForbidden
	}

	return Response{Allowed: false, Message: message, StatusCode: statusCode}
}

func (e *AuthorizationException) Error() string {
	if e.Response.Message != "" {
		return e.Response.Message
	}

	return fmt.Sprintf("this action is unauthorized (HTTP %d)", e.Response.StatusCode)
}

// New creates a Gate with the given user resolver.
func New(userResolver func(ctx context.Context) cauth.Authenticatable) *Gate {
	return &Gate{
		abilities:    make(map[string]Ability),
		policies:     make(map[string]any),
		userResolver: userResolver,
	}
}

// Define registers an ability by name.
func (g *Gate) Define(ability string, fn Ability) *Gate {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.abilities[ability] = fn

	return g
}

// Before registers a hook that runs before any ability check.
// If the hook returns (result, true) the ability check is skipped.
func (g *Gate) Before(fn func(ctx context.Context, user cauth.Authenticatable, ability string, model any) (bool, bool)) *Gate {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.before = append(g.before, fn)

	return g
}

// After registers a hook that runs after any ability check.
func (g *Gate) After(fn func(ctx context.Context, user cauth.Authenticatable, ability string, result bool, model any)) *Gate {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.after = append(g.after, fn)

	return g
}

// Inspect evaluates an ability and returns a Response.
func (g *Gate) Inspect(ctx context.Context, ability string, model any) Response {
	user := g.userResolver(ctx)

	g.mu.RLock()
	befores := g.before
	afters := g.after
	fn, hasFn := g.abilities[ability]
	g.mu.RUnlock()

	// Before hooks.
	for _, hook := range befores {
		if result, handled := hook(ctx, user, ability, model); handled {
			return g.buildResponse(result, "")
		}
	}

	// Check policies before raw abilities.
	if model != nil {
		if policyResult, handled := g.callPolicyMethod(ctx, user, ability, model); handled {
			resp := g.buildResponse(policyResult, "")

			for _, hook := range afters {
				hook(ctx, user, ability, policyResult, model)
			}

			return resp
		}
	}

	if !hasFn {
		resp := Deny(fmt.Sprintf("ability %q is not defined", ability), http.StatusForbidden)

		for _, hook := range afters {
			hook(ctx, user, ability, false, model)
		}

		return resp
	}

	allowed, err := fn(ctx, user, model)

	if err != nil {
		resp := Deny(err.Error(), http.StatusInternalServerError)

		for _, hook := range afters {
			hook(ctx, user, ability, false, model)
		}

		return resp
	}

	resp := g.buildResponse(allowed, "")

	for _, hook := range afters {
		hook(ctx, user, ability, allowed, model)
	}

	return resp
}

// Check returns true if the user can perform the ability.
func (g *Gate) Check(ctx context.Context, ability string, model any) bool {
	return g.Inspect(ctx, ability, model).Allowed
}

// Any returns true if the user can perform any of the abilities.
func (g *Gate) Any(ctx context.Context, abilities []string, model any) bool {
	for _, a := range abilities {
		if g.Check(ctx, a, model) {
			return true
		}
	}

	return false
}

// Every returns true if the user can perform all of the abilities.
func (g *Gate) Every(ctx context.Context, abilities []string, model any) bool {
	for _, a := range abilities {
		if !g.Check(ctx, a, model) {
			return false
		}
	}

	return true
}

// Authorize checks the ability and returns an error if unauthorized.
func (g *Gate) Authorize(ctx context.Context, ability string, model any) error {
	resp := g.Inspect(ctx, ability, model)

	if !resp.Allowed {
		return &AuthorizationException{Response: resp}
	}

	return nil
}

// Denies returns true if the user cannot perform the ability.
func (g *Gate) Denies(ctx context.Context, ability string, model any) bool {
	return !g.Check(ctx, ability, model)
}

// Has reports whether an ability has been defined.
func (g *Gate) Has(ability string) bool {
	g.mu.RLock()

	defer g.mu.RUnlock()

	_, ok := g.abilities[ability]

	return ok
}

// None returns true if the user cannot perform any of the given abilities.
func (g *Gate) None(ctx context.Context, abilities []string, model any) bool {
	return !g.Any(ctx, abilities, model)
}

// AllowIf returns an Allow response if condition is true, otherwise Deny.
func AllowIf(condition bool, message string, code int) Response {
	if condition {
		return Allow(message)
	}

	return Deny(message, code)
}

// DenyIf returns a Deny response if condition is true, otherwise Allow.
func DenyIf(condition bool, message string, code int) Response {
	if condition {
		return Deny(message, code)
	}

	return Allow(message)
}

// Abilities returns the names of all defined abilities.
func (g *Gate) Abilities() []string {
	g.mu.RLock()

	defer g.mu.RUnlock()

	names := make([]string, 0, len(g.abilities))

	for k := range g.abilities {
		names = append(names, k)
	}

	sort.Strings(names)

	return names
}

// RegisterPolicy registers a policy for the given model type name.
func (g *Gate) RegisterPolicy(modelType string, policy any) *Gate {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.policies[modelType] = policy

	return g
}

// GetPolicyFor returns the policy registered for the given model's type.
func (g *Gate) GetPolicyFor(model any) (any, bool) {
	g.mu.RLock()

	defer g.mu.RUnlock()

	if model == nil {
		return nil, false
	}

	typeName := reflect.TypeOf(model).String()
	policy, ok := g.policies[typeName]

	return policy, ok
}

// ForUser returns a new Gate that uses the given user instead of the resolver.
func (g *Gate) ForUser(user cauth.Authenticatable) *Gate {
	clone := &Gate{
		abilities:    g.abilities,
		policies:     g.policies,
		before:       g.before,
		after:        g.after,
		userResolver: func(_ context.Context) cauth.Authenticatable { return user },
	}

	return clone
}

// callPolicyMethod looks up and invokes a policy method for the given ability.
// Returns (result, true) if a policy method was found, (false, false) otherwise.
func (g *Gate) callPolicyMethod(ctx context.Context, user cauth.Authenticatable, ability string, model any) (bool, bool) {
	g.mu.RLock()
	typeName := reflect.TypeOf(model).String()
	policy, ok := g.policies[typeName]
	g.mu.RUnlock()

	if !ok {
		return false, false
	}

	pv := reflect.ValueOf(policy)

	// Check for a Before method first.
	if before := pv.MethodByName("Before"); before.IsValid() {
		results := before.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(user),
			reflect.ValueOf(ability),
		})

		if len(results) == 2 {
			result := results[0].Bool()
			handled := results[1].Bool()

			if handled {
				return result, true
			}
		}
	}

	// Capitalize the first letter of the ability for the method name.
	methodName := capitalizeFirst(ability)
	method := pv.MethodByName(methodName)

	if !method.IsValid() {
		return false, false
	}

	mt := method.Type()

	var results []reflect.Value

	switch mt.NumIn() {
	case 3: // (ctx, user, model)
		results = method.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(user),
			reflect.ValueOf(model),
		})
	case 2: // (ctx, user)
		results = method.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(user),
		})
	default:
		return false, false
	}

	if len(results) >= 1 && results[0].Kind() == reflect.Bool {
		return results[0].Bool(), true
	}

	return false, false
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}

	b := []byte(s)

	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 'a' - 'A'
	}

	return string(b)
}

func (g *Gate) buildResponse(allowed bool, message string) Response {
	if allowed {
		return Allow(message)
	}

	return Deny(message, http.StatusForbidden)
}
