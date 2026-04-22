package oauthserver

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// TokenGuard authenticates requests by resolving a Bearer token from the
// Authorization header, looking it up in the TokenStore, and retrieving the
// associated user via the UserProvider.
//
// It implements contracts/auth.Guard and mirrors Upstream OAuthServer's TokenGuard.
type TokenGuard struct {
	mu       sync.RWMutex
	oauthserver *OAuthServer
	tokens   TokenStore
	clients  ClientStore
	provider cauth.UserProvider
	request  *http.Request
	user     cauth.Authenticatable // cached *UserWithTokens
	token    *Token                // cached access token
	client   *Client               // cached client
}

// NewTokenGuard creates a TokenGuard.
func NewTokenGuard(
	p *OAuthServer,
	tokens TokenStore,
	clients ClientStore,
	provider cauth.UserProvider,
) *TokenGuard {
	return &TokenGuard{
		oauthserver: p,
		tokens:   tokens,
		clients:  clients,
		provider: provider,
	}
}

// SetRequest attaches the HTTP request, clearing any previously cached user
// and token so that the next User() call resolves fresh state.
func (g *TokenGuard) SetRequest(r *http.Request) {
	g.mu.Lock()

	defer g.mu.Unlock()

	g.request = r
	g.user = nil
	g.token = nil
	g.client = nil
}

// User returns the authenticated user for the current request.
//
// Resolution order:
//  1. ActingAs override (for tests)
//  2. Cached result from a previous call on this request
//  3. Bearer token from Authorization header → TokenStore.Find → UserProvider
func (g *TokenGuard) User(ctx context.Context) (cauth.Authenticatable, error) {
	g.mu.Lock()

	defer g.mu.Unlock()

	// 1. Test override.
	if actingUser, actingToken, _, _, ok := g.oauthserver.actingAsState(); ok && actingUser != nil {
		if g.user == nil {
			g.user = NewUserWithTokens(actingUser, actingToken)
		}

		return g.user, nil
	}

	// 2. Cached result.
	if g.user != nil {
		return g.user, nil
	}

	token := g.token
	if token == nil {
		var err error
		token, err = g.resolveToken(ctx)
		if err != nil || token == nil {
			return nil, err
		}

		g.token = token
	}

	// Machine-to-machine tokens (client credentials) have no user.
	if token.UserID == "" {
		return nil, nil
	}

	user, err := g.provider.RetrieveByID(ctx, token.UserID)

	if err != nil || user == nil {
		return nil, err
	}

	accessToken := NewAccessToken(token, g.oauthserver)

	g.user = NewUserWithTokens(user, accessToken)
	g.token = token

	return g.user, nil
}

// Client returns the OAuth client associated with the current request's token.
func (g *TokenGuard) Client(ctx context.Context) (*Client, error) {
	g.mu.Lock()

	defer g.mu.Unlock()

	// Test override: return acting-as client.
	if _, _, actingClient, scopes, ok := g.oauthserver.actingAsState(); ok && actingClient != nil {
		client := *actingClient
		if len(scopes) > 0 {
			client.Scopes = scopes
		}

		return &client, nil
	}

	if g.client != nil {
		return g.client, nil
	}

	if g.token == nil {
		token, err := g.resolveToken(ctx)
		if err != nil || token == nil {
			return nil, err
		}

		g.token = token
	}

	client, err := g.clients.FindActive(ctx, g.token.ClientID)
	if err != nil || client == nil {
		return nil, err
	}

	g.client = client

	return g.client, nil
}

// Token returns the current request's bearer token, if one is present and valid.
func (g *TokenGuard) Token(ctx context.Context) (*Token, error) {
	g.mu.Lock()

	defer g.mu.Unlock()

	if g.token == nil {
		token, err := g.resolveToken(ctx)
		if err != nil || token == nil {
			return nil, err
		}

		g.token = token
	}

	copy := *g.token

	return &copy, nil
}

func (g *TokenGuard) resolveToken(ctx context.Context) (*Token, error) {
	bearer := g.bearerToken()

	if bearer == "" {
		return nil, nil
	}

	token, err := g.tokens.Find(ctx, bearer)

	if err != nil {
		return nil, err
	}

	if token == nil || token.IsRevoked() {
		return nil, nil
	}

	token = token.WithOAuthServer(g.oauthserver)

	if !token.ExpiresAt.IsZero() && token.ExpiresAt.Before(time.Now()) {
		return nil, nil
	}

	return token, nil
}

// Check reports whether the request is authenticated.
func (g *TokenGuard) Check(ctx context.Context) bool {
	u, _ := g.User(ctx)

	return u != nil
}

// Guest reports whether the request is unauthenticated.
func (g *TokenGuard) Guest(ctx context.Context) bool {
	return !g.Check(ctx)
}

// ID returns the authenticated user's identifier, or nil if unauthenticated.
func (g *TokenGuard) ID(ctx context.Context) any {
	u, _ := g.User(ctx)

	if u == nil {
		return nil
	}

	return u.GetAuthIdentifier()
}

// bearerToken extracts the Bearer token from the Authorization header.
// Returns an empty string when the header is absent or malformed.
func (g *TokenGuard) bearerToken() string {
	if g.request == nil {
		return ""
	}

	header := g.request.Header.Get("Authorization")

	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}

	return ""
}
