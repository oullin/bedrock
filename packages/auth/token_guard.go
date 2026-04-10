package auth

import (
	"context"
	"net/http"
	"strings"
	"sync"
)

// TokenGuard authenticates requests via a bearer token found in the query
// string, form body, request header, or HTTP Basic Auth.
type TokenGuard struct {
	mu         sync.RWMutex
	name       string
	provider   UserProvider
	request    *http.Request
	user       Authenticatable
	inputKey   string // query/form field name (default: "api_token")
	storageKey string // user attribute key (default: "api_token")
	hashable   bool   // whether tokens are hashed in storage
}

// NewTokenGuard creates a TokenGuard with default key names.
func NewTokenGuard(name string, provider UserProvider) *TokenGuard {
	return &TokenGuard{
		name:       name,
		provider:   provider,
		inputKey:   "api_token",
		storageKey: "api_token",
	}
}

// SetInputKey sets the query/form parameter name used to read the token.
func (g *TokenGuard) SetInputKey(key string) { g.inputKey = key }

// SetStorageKey sets the user attribute key for the stored token.
func (g *TokenGuard) SetStorageKey(key string) { g.storageKey = key }

// SetRequest attaches the incoming HTTP request.
func (g *TokenGuard) SetRequest(r *http.Request) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.request = r
	g.user = nil
}

// SetUser sets the authenticated user on the guard.
func (g *TokenGuard) SetUser(user Authenticatable) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.user = user
}

// HasUser reports whether the guard has a resolved user without triggering resolution.
func (g *TokenGuard) HasUser() bool {
	g.mu.RLock()

	defer g.mu.RUnlock()

	return g.user != nil
}

// ForgetUser clears the resolved user, forcing re-resolution on the next User() call.
func (g *TokenGuard) ForgetUser() {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.user = nil
}

// User returns the authenticated user, or nil if the token is absent/invalid.
func (g *TokenGuard) User(ctx context.Context) (Authenticatable, error) {
	g.mu.Lock()

	defer g.mu.Unlock()

	if g.user != nil {
		return g.user, nil
	}

	token := g.getTokenFromRequest()

	if token == "" {
		return nil, nil
	}

	user, err := g.provider.RetrieveByCredentials(ctx, map[string]any{
		g.storageKey: token,
	})

	if err != nil || user == nil {
		return nil, err
	}

	g.user = user

	return user, nil
}

func (g *TokenGuard) Check(ctx context.Context) bool {
	u, _ := g.User(ctx)

	return u != nil
}

func (g *TokenGuard) Guest(ctx context.Context) bool { return !g.Check(ctx) }

func (g *TokenGuard) ID(ctx context.Context) any {
	u, _ := g.User(ctx)

	if u == nil {
		return nil
	}

	return u.GetAuthIdentifier()
}

func (g *TokenGuard) getTokenFromRequest() string {
	if g.request == nil {
		return ""
	}

	// 1. Query string.
	if t := g.request.URL.Query().Get(g.inputKey); t != "" {
		return t
	}

	// 2. Form field.
	if t := g.request.FormValue(g.inputKey); t != "" {
		return t
	}

	// 3. Authorization header (Bearer).
	if auth := g.request.Header.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	// 4. HTTP Basic Auth (password field).
	if _, password, ok := g.request.BasicAuth(); ok && password != "" {
		return password
	}

	return ""
}
