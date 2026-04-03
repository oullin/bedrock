package access

import (
	"context"
	"fmt"
	"sync"

	auth "github.com/gollin/packages/auth"
)

// Response represents an authorization decision.
type Response struct {
	Allowed bool
	Message string
}

// Allow returns an allowed response.

// Deny returns a denied response.

// AuthorizationException reports an authorization failure.
type AuthorizationException struct {
	Ability string
	Message string
}

// Error implements error.

// AbilityFunc evaluates an ability for a user.
type AbilityFunc func(ctx context.Context, user auth.Authenticatable, arguments ...any) Response

// BeforeFunc can short-circuit an ability check.
type BeforeFunc func(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) (Response, bool)

// AfterFunc can adjust the computed result.
type AfterFunc func(ctx context.Context, user auth.Authenticatable, ability string, result Response, arguments ...any) Response

// Authorizer is the interface required by Foundation/Auth helpers.
type Authorizer interface {
	Inspect(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) Response
	Check(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) bool
	Any(ctx context.Context, user auth.Authenticatable, abilities []string, arguments ...any) bool
	Denies(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) bool
	Authorize(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) error
}

// Gate stores defined authorization abilities.
type Gate struct {
	mu        sync.RWMutex
	abilities map[string]AbilityFunc
	before    []BeforeFunc
	after     []AfterFunc
}

func Allow() Response {
	return Response{Allowed: true}
}

func Deny(message string) Response {
	return Response{Allowed: false, Message: message}
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
	return &Gate{abilities: map[string]AbilityFunc{}}
}

// Define registers an ability callback.
func (g *Gate) Define(ability string, callback AbilityFunc) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.abilities[ability] = callback
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
	callback, ok := g.abilities[ability]
	g.mu.RUnlock()

	for _, hook := range before {
		if result, handled := hook(ctx, user, ability, arguments...); handled {
			return result
		}
	}

	result := Deny("auth: ability is not defined")

	if ok {
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

// Denies reports whether an ability is denied.
func (g *Gate) Denies(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) bool {
	return !g.Check(ctx, user, ability, arguments...)
}

// Authorize authorizes an ability or returns an exception.
func (g *Gate) Authorize(ctx context.Context, user auth.Authenticatable, ability string, arguments ...any) error {
	response := g.Inspect(ctx, user, ability, arguments...)

	if response.Allowed {
		return nil
	}

	return AuthorizationException{Ability: ability, Message: response.Message}
}

var _ Authorizer = (*Gate)(nil)
