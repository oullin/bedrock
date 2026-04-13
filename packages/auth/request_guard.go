package auth

import (
	"context"
	"net/http"
	"sync"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// RequestGuard authenticates requests via a custom callback function.
type RequestGuard struct {
	mu       sync.RWMutex
	callback RequestCallback
	request  *http.Request
	user     cauth.Authenticatable
	provider cauth.UserProvider
}

// NewRequestGuard creates a RequestGuard using the given callback.
func NewRequestGuard(callback RequestCallback) *RequestGuard {
	return &RequestGuard{callback: callback}
}

// SetUser sets the authenticated user on the guard.
func (g *RequestGuard) SetUser(user cauth.Authenticatable) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.user = user
}

// HasUser reports whether the guard has a resolved user without triggering resolution.
func (g *RequestGuard) HasUser() bool {
	g.mu.RLock()

	defer g.mu.RUnlock()

	return g.user != nil
}

// ForgetUser clears the resolved user, forcing re-resolution on the next User() call.
func (g *RequestGuard) ForgetUser() {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.user = nil
}

// SetRequest attaches the incoming HTTP request.
func (g *RequestGuard) SetRequest(r *http.Request) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.request = r
	g.user = nil
}

// User resolves the user via the callback.
func (g *RequestGuard) User(ctx context.Context) (cauth.Authenticatable, error) {
	g.mu.Lock()

	defer g.mu.Unlock()

	if g.user != nil {
		return g.user, nil
	}

	if g.callback == nil || g.request == nil {
		return nil, nil
	}

	user, err := g.callback(ctx, g.request)

	if err != nil || user == nil {
		return nil, err
	}

	g.user = user

	return user, nil
}

func (g *RequestGuard) Check(ctx context.Context) bool {
	u, _ := g.User(ctx)

	return u != nil
}

func (g *RequestGuard) Guest(ctx context.Context) bool { return !g.Check(ctx) }

func (g *RequestGuard) ID(ctx context.Context) any {
	u, _ := g.User(ctx)

	if u == nil {
		return nil
	}

	return u.GetAuthIdentifier()
}

// Validate always returns false for request guards (no credential-based auth).
func (g *RequestGuard) Validate(_ context.Context, _ map[string]string) bool {
	return false
}

// GetProvider returns the user provider.
func (g *RequestGuard) GetProvider() cauth.UserProvider {
	g.mu.RLock()

	defer g.mu.RUnlock()

	return g.provider
}

// SetProvider sets the user provider.
func (g *RequestGuard) SetProvider(p cauth.UserProvider) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.provider = p
}
