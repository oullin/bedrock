// Package auth is the facade for the authentication manager. It forwards
// common operations to the auth.Manager bound under "auth" in the global
// bedrock Application.
package auth

import (
	"context"
	"sync"

	"github.com/bedrock/app"
	authpkg "github.com/bedrock/packages/auth"
	cauth "github.com/bedrock/packages/contracts/auth"
)

var (
	mu     sync.Mutex
	cached *authpkg.Manager
)

// Manager returns the auth manager from the global Application. Resolved
// once per process and cached.
func Manager() *authpkg.Manager {
	mu.Lock()

	defer mu.Unlock()

	if cached == nil {
		cached = app.Resolve[*authpkg.Manager]("auth")
	}

	return cached
}

// Reset clears the cached manager. Tests must call this after reinstalling
// a different Application via app.SetApp.
func Reset() {
	mu.Lock()

	defer mu.Unlock()

	cached = nil
}

// Guard returns the named guard. Pass an empty string for the default guard.
func Guard(ctx context.Context, name string) (cauth.Guard, error) {
	return Manager().Guard(ctx, name)
}

// User resolves the currently-authenticated user via the registered
// userResolver, or returns nil if none is set.
func User(ctx context.Context) cauth.Authenticatable {
	resolver := Manager().UserResolver()

	if resolver == nil {
		return nil
	}

	return resolver(ctx)
}

// Check reports whether a user is authenticated.
func Check(ctx context.Context) bool {
	return User(ctx) != nil
}
