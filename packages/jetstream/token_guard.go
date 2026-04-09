package jetstream

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/bedrock/packages/fortify"
)

// TokenGuard authenticates requests using personal access tokens.
// It looks up the bearer token, resolves the user, and places
// both in context via the HasApiTokens interface.
type TokenGuard struct {
	tokens   TokenRepository
	provider fortify.UserProvider
}

// NewTokenGuard creates a new API token guard.
func NewTokenGuard(tokens TokenRepository, provider fortify.UserProvider) *TokenGuard {
	return &TokenGuard{tokens: tokens, provider: provider}
}

// Authenticate extracts a bearer token, validates it, and returns the user.
func (g *TokenGuard) Authenticate(ctx context.Context, r *http.Request) (fortify.Authenticatable, *PersonalAccessToken, error) {
	plain := extractBearerToken(r)
	if plain == "" {
		return nil, nil, ErrUnauthenticated
	}

	hash := HashToken(plain)

	token, err := g.tokens.FindByTokenHash(ctx, hash)
	if err != nil || token == nil {
		return nil, nil, ErrUnauthenticated
	}

	if token.IsExpired() {
		return nil, nil, ErrUnauthenticated
	}

	user, err := g.provider.RetrieveByID(ctx, token.UserID)
	if err != nil || user == nil {
		return nil, nil, ErrUnauthenticated
	}

	now := time.Now()
	token.LastUsedAt = &now
	_ = g.tokens.Update(ctx, token)

	if apiUser, ok := user.(HasApiTokens); ok {
		apiUser.SetCurrentAccessToken(token)
	}

	return user, token, nil
}

func extractBearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

// TokenCan returns middleware that checks whether the authenticated token
// has the given permission.
func TokenCan(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := TokenFromContext(r.Context())
			if token == nil {
				http.Error(w, "no access token", http.StatusUnauthorized)
				return
			}

			if !token.HasPermission(permission) {
				http.Error(w, "insufficient token permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type tokenContextKey struct{}

// WithToken returns a context with the token set.
func WithToken(ctx context.Context, token *PersonalAccessToken) context.Context {
	return context.WithValue(ctx, tokenContextKey{}, token)
}

// TokenFromContext retrieves the token from context.
func TokenFromContext(ctx context.Context) *PersonalAccessToken {
	token, _ := ctx.Value(tokenContextKey{}).(*PersonalAccessToken)
	return token
}
