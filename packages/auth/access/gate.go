package access

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/bedrock/packages/auth"
)

// Ability is a function that determines whether a user can perform an action.
// It receives the user, the optional model, and returns (allow bool, err error).
type Ability func(ctx context.Context, user auth.Authenticatable, model any) (bool, error)

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
	before       []func(ctx context.Context, user auth.Authenticatable, ability string, model any) (bool, bool)
	after        []func(ctx context.Context, user auth.Authenticatable, ability string, result bool, model any)
	userResolver func(ctx context.Context) auth.Authenticatable
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
func New(userResolver func(ctx context.Context) auth.Authenticatable) *Gate {
	return &Gate{
		abilities:    make(map[string]Ability),
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
func (g *Gate) Before(fn func(ctx context.Context, user auth.Authenticatable, ability string, model any) (bool, bool)) *Gate {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.before = append(g.before, fn)

	return g
}

// After registers a hook that runs after any ability check.
func (g *Gate) After(fn func(ctx context.Context, user auth.Authenticatable, ability string, result bool, model any)) *Gate {
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

// ForUser returns a new Gate that uses the given user instead of the resolver.
func (g *Gate) ForUser(user auth.Authenticatable) *Gate {
	clone := &Gate{
		abilities:    g.abilities,
		before:       g.before,
		after:        g.after,
		userResolver: func(_ context.Context) auth.Authenticatable { return user },
	}

	return clone
}

func (g *Gate) buildResponse(allowed bool, message string) Response {
	if allowed {
		return Allow(message)
	}

	return Deny(message, http.StatusForbidden)
}
