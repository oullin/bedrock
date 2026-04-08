package access

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"

	auth "github.com/bedrock/packages/anvil/auth"
)

// Response represents an authorization decision.
type Response struct {
	Allowed bool
	Message string
	Code    int
	Status  int
}

// AuthorizationException reports an authorization failure.
type AuthorizationException struct {
	Ability string
	Message string
}

// AbilityFunc evaluates an ability for a user.
type AbilityFunc func(ctx context.Context, user auth.Authenticatable, arguments ...any) Response

// PolicyFunc evaluates an ability for a registered resource or model target.
type PolicyFunc func(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) (Response, bool)

// BeforeFunc can short-circuit an ability check.
type BeforeFunc func(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) (Response, bool)

// AfterFunc can adjust the computed result.
type AfterFunc func(ctx context.Context, user auth.Authenticatable, ability string, result Response, arguments ...any) Response

// Authorizer is the interface required by consumer helpers.
type Authorizer interface {
	Inspect(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) Response
	Check(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) bool
	Any(ctx context.Context, user auth.Authenticatable, abilities []string, arguments ...any) bool
	Every(ctx context.Context, user auth.Authenticatable, abilities []string, arguments ...any) bool
	Denies(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) bool
	Authorize(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) error
}

// Gate stores defined abilities and policies.
type Gate struct {
	mu                    sync.RWMutex
	abilities             map[string]AbilityFunc
	policies              map[string]PolicyFunc
	before                []BeforeFunc
	after                 []AfterFunc
	defaultDenialResponse *Response
}

// Allow returns an allowed response.
func Allow() Response {
	return Response{Allowed: true}
}

// Deny returns a denied response.
func Deny(message string) Response {
	return Response{Allowed: false, Message: message}
}

// DenyWithCode returns a denied response carrying a custom application code and HTTP status.
func DenyWithCode(message string, code int, status int) Response {
	return Response{Allowed: false, Message: message, Code: code, Status: status}
}

// ToMap returns the response as a map.
func (r Response) ToMap() map[string]any {
	return map[string]any{
		"allowed": r.Allowed,
		"message": r.Message,
		"code":    r.Code,
		"status":  r.Status,
	}
}

func (e AuthorizationException) Error() string {
	if e.Message != "" {
		return e.Message
	}

	if e.Ability == "" {
		return "auth: authorization denied"
	}

	return fmt.Sprintf("auth: authorization denied for %q", e.Ability)
}

// NewGate returns an empty gate.
func NewGate() *Gate {
	return &Gate{
		abilities: map[string]AbilityFunc{},
		policies:  map[string]PolicyFunc{},
	}
}

// NewGateWithDenialResponse returns a gate with a default denial response.
func NewGateWithDenialResponse(response Response) *Gate {
	return &Gate{
		abilities:             map[string]AbilityFunc{},
		policies:              map[string]PolicyFunc{},
		defaultDenialResponse: &response,
	}
}

// SetDefaultDenialResponse sets a default denial response used by Authorize.
func (g *Gate) SetDefaultDenialResponse(response Response) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.defaultDenialResponse = &response
}

// Define registers an ability callback.
func (g *Gate) Define(ability string, callback AbilityFunc) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.abilities[strings.TrimSpace(ability)] = callback
}

// Policy registers a policy callback for a resource target.
func (g *Gate) Policy(target any, callback PolicyFunc) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.policies[policyKey(target)] = callback
}

// Before registers a before callback.
func (g *Gate) Before(callback BeforeFunc) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.before = append(g.before, callback)
}

// After registers an after callback.
func (g *Gate) After(callback AfterFunc) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.after = append(g.after, callback)
}

// Inspect evaluates an ability and returns the full response.
func (g *Gate) Inspect(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) Response {
	g.mu.RLock()
	before := append([]BeforeFunc(nil), g.before...)
	after := append([]AfterFunc(nil), g.after...)
	callback, ok := g.abilities[strings.TrimSpace(ability)]
	policies := make(map[string]PolicyFunc, len(g.policies))
	for key, value := range g.policies {
		policies[key] = value
	}
	g.mu.RUnlock()

	for _, hook := range before {
		if result, handled := hook(ctx, user, ability, arguments...); handled {
			return result
		}
	}

	result := Deny("auth: ability is not defined")

	if len(arguments) > 0 {
		if policy, ok := policies[policyKey(arguments[0])]; ok {
			if policyResult, handled := policy(ctx, user, ability, arguments...); handled {
				result = policyResult
			}
		}
	}

	if !result.Allowed && result.Message == "auth: ability is not defined" && ok {
		result = callback(ctx, user, arguments...)
	}

	for _, hook := range after {
		result = hook(ctx, user, ability, result, arguments...)
	}

	return result
}

// Check reports whether a user is authorized.
func (g *Gate) Check(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) bool {
	return g.Inspect(ctx, user, ability, arguments...).Allowed
}

// Any reports whether any ability is authorized.
func (g *Gate) Any(ctx context.Context, user auth.Authenticatable, abilities []string, arguments ...any) bool {
	for _, ability := range abilities {
		if g.Check(ctx, user, ability, arguments...) {
			return true
		}
	}

	return false
}

// Every reports whether all given abilities are authorized.
func (g *Gate) Every(ctx context.Context, user auth.Authenticatable, abilities []string, arguments ...any) bool {
	for _, ability := range abilities {
		if !g.Check(ctx, user, ability, arguments...) {
			return false
		}
	}

	return true
}

// Denies reports whether an ability is denied.
func (g *Gate) Denies(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) bool {
	return !g.Check(ctx, user, ability, arguments...)
}

// Has reports whether an ability has been defined.
func (g *Gate) Has(ability string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	_, ok := g.abilities[strings.TrimSpace(ability)]
	return ok
}

// None reports whether all given abilities are denied.
func (g *Gate) None(ctx context.Context, user auth.Authenticatable, abilities []string, arguments ...any) bool {
	return !g.Any(ctx, user, abilities, arguments...)
}

// AllowIf grants access when the condition is true, otherwise denies with the given message.
func (g *Gate) AllowIf(condition bool, message string) Response {
	if condition {
		return Allow()
	}

	return Deny(message)
}

// DenyIf denies access when the condition is true, otherwise allows.
func (g *Gate) DenyIf(condition bool, message string) Response {
	if condition {
		return Deny(message)
	}

	return Allow()
}

// ForUser returns a copy of the gate with all definitions but no user binding
// (Go gates are stateless; this clones abilities/policies/hooks so the caller
// can layer per-user overrides without mutating the original).
func (g *Gate) ForUser() *Gate {
	g.mu.RLock()
	defer g.mu.RUnlock()

	child := &Gate{
		abilities: make(map[string]AbilityFunc, len(g.abilities)),
		policies:  make(map[string]PolicyFunc, len(g.policies)),
		before:    append([]BeforeFunc(nil), g.before...),
		after:     append([]AfterFunc(nil), g.after...),
	}

	for k, v := range g.abilities {
		child.abilities[k] = v
	}

	for k, v := range g.policies {
		child.policies[k] = v
	}

	return child
}

// Resource bulk-registers CRUD abilities (viewAny, view, create, update, delete)
// for a named resource using the supplied policy.
func (g *Gate) Resource(resource string, policy PolicyFunc) {
	for _, ability := range []string{"viewAny", "view", "create", "update", "delete"} {
		name := resource + "." + ability
		localAbility := ability

		g.Define(name, func(ctx context.Context, user auth.Authenticatable, arguments ...any) Response {
			result, handled := policy(ctx, user, localAbility, arguments...)
			if handled {
				return result
			}

			return Deny("auth: ability is not defined")
		})
	}
}

// Authorize authorizes an ability or returns an exception.
func (g *Gate) Authorize(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) error {
	response := g.Inspect(ctx, user, ability, arguments...)
	if response.Allowed {
		return nil
	}

	message := response.Message
	if message == "" || message == "auth: ability is not defined" {
		g.mu.RLock()
		def := g.defaultDenialResponse
		g.mu.RUnlock()

		if def != nil {
			message = def.Message
		}
	}

	return AuthorizationException{Ability: ability, Message: message}
}

func policyKey(target any) string {
	switch value := target.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(value)
	}

	typ := reflect.TypeOf(target)
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	if typ.PkgPath() == "" {
		return typ.Name()
	}

	return typ.PkgPath() + "." + typ.Name()
}

var _ Authorizer = (*Gate)(nil)
