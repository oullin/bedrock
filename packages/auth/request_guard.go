package auth

import (
	"context"
	"net/http"
)

// RequestGuardCallback resolves a user from the current request.
type RequestGuardCallback func(ctx context.Context, r *http.Request, provider UserProvider) (Authenticatable, error)

// RequestGuard mirrors Laravel's callback-based request guard.
type RequestGuard struct {
	callback RequestGuardCallback
	provider UserProvider
	request  *http.Request
	user     Authenticatable
	resolved bool
}

// NewRequestGuard creates a new request guard.
func NewRequestGuard(callback RequestGuardCallback, request *http.Request, provider UserProvider) *RequestGuard {
	return &RequestGuard{
		callback: callback,
		provider: provider,
		request:  request,
	}
}

// SetRequest updates the current request.
func (g *RequestGuard) SetRequest(request *http.Request) {
	g.request = request
	g.resolved = false
	g.user = nil
}

// User resolves the current user.
func (g *RequestGuard) User(ctx context.Context) (Authenticatable, error) {
	if g.resolved {
		return g.user, nil
	}

	g.resolved = true

	if g.callback == nil || g.request == nil {
		return nil, ErrUnauthorized
	}

	user, err := g.callback(ctx, g.request, g.provider)

	if err != nil {
		return nil, err
	}

	g.user = user

	return user, nil
}
