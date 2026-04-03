package auth

import (
	"context"
	"net/http"
	"strings"

	securitycrypto "github.com/gollin/packages/security/crypto"
)

// TokenGuard resolves users from request tokens.
type TokenGuard struct {
	provider   UserProvider
	request    *http.Request
	inputKey   string
	storageKey string
	hash       bool
	user       Authenticatable
	resolved   bool
}

// NewTokenGuard creates a token guard.
func NewTokenGuard(provider UserProvider, request *http.Request, inputKey string, storageKey string, hash bool) *TokenGuard {
	if inputKey == "" {
		inputKey = "api_token"
	}

	if storageKey == "" {
		storageKey = inputKey
	}

	return &TokenGuard{
		provider:   provider,
		request:    request,
		inputKey:   inputKey,
		storageKey: storageKey,
		hash:       hash,
	}
}

// SetRequest updates the current request.
func (g *TokenGuard) SetRequest(request *http.Request) {
	g.request = request
	g.user = nil
	g.resolved = false
}

// User resolves the current user.
func (g *TokenGuard) User(ctx context.Context) (Authenticatable, error) {
	if g.resolved {
		return g.user, nil
	}

	g.resolved = true

	if g.request == nil {
		return nil, ErrUnauthorized
	}

	token := strings.TrimSpace(g.request.Header.Get("Authorization"))
	if token == "" {
		token = strings.TrimSpace(g.request.URL.Query().Get(g.inputKey))
	}

	if token == "" {
		return nil, ErrUnauthorized
	}

	const bearer = "Bearer "
	if strings.HasPrefix(token, bearer) {
		token = strings.TrimSpace(strings.TrimPrefix(token, bearer))
	}

	credentials := map[string]string{g.storageKey: token}
	if g.hash {
		credentials[g.storageKey] = securitycrypto.HashString(token)
	}

	user, err := g.provider.RetrieveByCredentials(ctx, credentials)
	if err != nil {
		return nil, err
	}

	g.user = user

	return user, nil
}
