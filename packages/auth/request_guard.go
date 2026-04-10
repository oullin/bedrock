package auth

import (
	"context"
	"net/http"
	"sync"
)

// RequestGuard authenticates requests via a custom callback function.
type RequestGuard struct {
	mu       sync.RWMutex
	callback RequestCallback
	request  *http.Request
	user     Authenticatable
}

// NewRequestGuard creates a RequestGuard using the given callback.
func NewRequestGuard(callback RequestCallback) *RequestGuard {
	return &RequestGuard{callback: callback}
}

// SetUser sets the authenticated user on the guard.
func (g *RequestGuard) SetUser(user Authenticatable) {
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
func (g *RequestGuard) User(ctx context.Context) (Authenticatable, error) {
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
