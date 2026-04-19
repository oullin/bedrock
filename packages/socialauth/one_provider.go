package socialauth

import (
	"context"
	"net/http"
)

// TemporaryCredentials holds the OAuth1 request-token pair returned by the
// provider during the redirect phase.
type TemporaryCredentials struct {
	Identifier string
	Secret     string
}

// TokenCredentials holds the OAuth1 access-token pair obtained after the user
// authorizes the application.
type TokenCredentials struct {
	Identifier string
	Secret     string
}

// OAuth1Server is the interface that concrete OAuth1 providers implement.
// It mirrors the League\OAuth1\Client\Server\Server contract used by
// Upstream\SocialAuth\One\AbstractProvider.
type OAuth1Server interface {
	// GetTemporaryCredentials fetches a request token from the provider.
	GetTemporaryCredentials(ctx context.Context) (*TemporaryCredentials, error)

	// GetAuthorizationURL builds the URL to which the user is redirected.
	GetAuthorizationURL(temp *TemporaryCredentials) string

	// GetTokenCredentials exchanges a request token + verifier for an access
	// token.
	GetTokenCredentials(ctx context.Context, temp *TemporaryCredentials, oauthToken, verifier string) (*TokenCredentials, error)

	// GetUserDetails fetches the raw user payload using the access token.
	GetUserDetails(ctx context.Context, token *TokenCredentials) (map[string]any, error)

	// MapUserToObject converts the raw payload and token into a *User.
	MapUserToObject(raw map[string]any, token *TokenCredentials) *User
}

// OneAbstractProvider is the base OAuth1 provider.
// It mirrors Upstream\SocialAuth\One\AbstractProvider.
type OneAbstractProvider struct {
	server  OAuth1Server
	request *http.Request
	session Session
}

// NewOneAbstractProvider creates an OAuth1 provider backed by the given server.
func NewOneAbstractProvider(server OAuth1Server, req *http.Request, session Session) *OneAbstractProvider {
	return &OneAbstractProvider{
		server:  server,
		request: req,
		session: session,
	}
}

// Redirect fetches temporary credentials, stores them in the session, and
// returns the authorization URL. It mirrors One\AbstractProvider::redirect().
func (p *OneAbstractProvider) Redirect(ctx context.Context) (string, error) {
	temp, err := p.server.GetTemporaryCredentials(ctx)

	if err != nil {
		return "", err
	}

	p.session.Put("oauth.temp", temp)

	return p.server.GetAuthorizationURL(temp), nil
}

// User validates the OAuth1 callback parameters, fetches token credentials,
// retrieves user details, and returns a populated *User.
// It mirrors One\AbstractProvider::user().
func (p *OneAbstractProvider) User(ctx context.Context) (*User, error) {
	oauthToken := p.request.URL.Query().Get("oauth_token")
	verifier := p.request.URL.Query().Get("oauth_verifier")

	if verifier == "" {
		return nil, ErrMissingVerifier
	}

	temp, ok := p.session.Get("oauth.temp").(*TemporaryCredentials)

	if !ok || temp == nil {
		return nil, ErrMissingTemporaryCredentials
	}

	tokenCreds, err := p.server.GetTokenCredentials(ctx, temp, oauthToken, verifier)

	if err != nil {
		return nil, err
	}

	raw, err := p.server.GetUserDetails(ctx, tokenCreds)

	if err != nil {
		return nil, err
	}

	return p.server.MapUserToObject(raw, tokenCreds), nil
}
