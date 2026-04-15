package oauthserver

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// PersonalAccessTokenResult holds the result of creating a personal access token.
// It mirrors Upstream OAuthServer's PersonalAccessTokenResult class.
type PersonalAccessTokenResult struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int
	ExpiresAt   time.Time
}

// ToArray returns a map representation for JSON serialization.
func (r *PersonalAccessTokenResult) ToArray() map[string]any {
	return map[string]any{
		"accessToken": r.AccessToken,
		"token":       r.AccessToken,
		"type":        r.TokenType,
		"expiresIn":   r.ExpiresIn,
		"expiresAt":   r.ExpiresAt.Format(time.RFC3339),
	}
}

// MarshalJSON implements json.Marshaler using ToArray.
func (r *PersonalAccessTokenResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.ToArray())
}

// PersonalAccessTokenFactory creates personal access tokens on behalf of users.
// It mirrors Upstream OAuthServer's PersonalAccessTokenFactory class.
//
// When an AuthorizationServer is available it delegates to server.IssueToken.
// Otherwise it creates the token directly against the TokenStore (useful for
// environments without a full OAuth2 server wired up).
type PersonalAccessTokenFactory struct {
	oauthserver *OAuthServer
	server   AuthorizationServer // optional — nil-safe
	tokens   TokenStore
	events   EventDispatcher // optional — nil-safe
}

// NewPersonalAccessTokenFactory constructs the factory.
// server may be nil when using the direct-creation path.
func NewPersonalAccessTokenFactory(p *OAuthServer, server AuthorizationServer, tokens TokenStore) *PersonalAccessTokenFactory {
	return &PersonalAccessTokenFactory{
		oauthserver: p,
		server:   server,
		tokens:   tokens,
	}
}

// WithEventDispatcher attaches an event dispatcher to the factory.
func (f *PersonalAccessTokenFactory) WithEventDispatcher(d EventDispatcher) *PersonalAccessTokenFactory {
	f.events = d

	return f
}

// Create creates a new personal access token for the given user.
//
// When an AuthorizationServer is configured it uses the personal_access grant
// flow. Otherwise it generates the token record directly and saves it to the
// TokenStore.
func (f *PersonalAccessTokenFactory) Create(
	ctx context.Context,
	user cauth.Authenticatable,
	name string,
	scopes []string,
) (*PersonalAccessTokenResult, error) {
	if f.server != nil {
		return f.createViaServer(ctx, user, name, scopes)
	}

	return f.createDirect(ctx, user, name, scopes)
}

func (f *PersonalAccessTokenFactory) createViaServer(
	ctx context.Context,
	user cauth.Authenticatable,
	name string,
	scopes []string,
) (*PersonalAccessTokenResult, error) {
	req := TokenRequest{
		GrantType: GrantPersonalAccess,
		UserID:    user.GetAuthIdentifier(),
		TokenName: name,
		Scopes:    scopes,
	}

	if cfg := f.oauthserver.Config(); cfg != nil {
		req.ClientID = cfg.PersonalAccessClientID
		req.ClientSecret = cfg.PersonalAccessClientSecret
	}

	issued, err := f.server.IssueToken(ctx, req)

	if err != nil {
		return nil, fmt.Errorf("oauthserver: personal access token: %w", err)
	}

	dispatch(ctx, f.events, AccessTokenCreated{
		TokenID:  issued.TokenID,
		UserID:   user.GetAuthIdentifier(),
		ClientID: req.ClientID,
	})

	ttl := f.oauthserver.PersonalAccessTokensTTL()

	return &PersonalAccessTokenResult{
		AccessToken: issued.AccessToken,
		TokenType:   issued.TokenType,
		ExpiresIn:   issued.ExpiresIn,
		ExpiresAt:   time.Now().Add(ttl),
	}, nil
}

func (f *PersonalAccessTokenFactory) createDirect(
	ctx context.Context,
	user cauth.Authenticatable,
	name string,
	scopes []string,
) (*PersonalAccessTokenResult, error) {
	id, err := generateID()

	if err != nil {
		return nil, fmt.Errorf("oauthserver: generate token id: %w", err)
	}

	ttl := f.oauthserver.PersonalAccessTokensTTL()
	now := time.Now()

	token := &Token{
		ID:        id,
		UserID:    user.GetAuthIdentifier(),
		Name:      name,
		Scopes:    scopes,
		Revoked:   false,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(ttl),
	}

	if err := f.tokens.Save(ctx, token); err != nil {
		return nil, fmt.Errorf("oauthserver: save personal access token: %w", err)
	}

	dispatch(ctx, f.events, AccessTokenCreated{
		TokenID: id,
		UserID:  user.GetAuthIdentifier(),
	})

	return &PersonalAccessTokenResult{
		AccessToken: id,
		TokenType:   "Bearer",
		ExpiresIn:   int(ttl.Seconds()),
		ExpiresAt:   now.Add(ttl),
	}, nil
}
