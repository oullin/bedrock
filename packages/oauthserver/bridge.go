package oauthserver

import "context"

// AuthorizationServer issues OAuth2 tokens and handles authorization requests.
// Implementations wrap a Go OAuth2 backend (e.g. ory/fosite, go-oauth2/oauth2)
// and expose the behavior OAuthServer delegates to League OAuth2 Server.
type AuthorizationServer interface {
	// IssueToken handles a token endpoint request for the given grant type.
	IssueToken(ctx context.Context, req TokenRequest) (*IssuedToken, error)

	// ValidateAuthorizationRequest validates an OAuth2 GET /oauth/authorize request.
	ValidateAuthorizationRequest(ctx context.Context, req AuthorizationRequest) (*ValidatedAuthRequest, error)

	// CompleteAuthorizationRequest completes the authorization flow after
	// the user approves or denies the request.
	CompleteAuthorizationRequest(ctx context.Context, req *ValidatedAuthRequest, approved bool) (*AuthorizationResponse, error)
}

// ResourceServer validates bearer tokens on incoming API requests.
// Implementations verify token signatures and look up token records.
type ResourceServer interface {
	// ValidateToken parses and validates a bearer token string.
	// Returns the validated Token record or an error.
	ValidateToken(ctx context.Context, tokenString string) (*Token, error)
}

// TokenRequest represents a token endpoint POST request.
// All grant-specific fields are present but only the relevant ones are populated
// depending on the GrantType value.
type TokenRequest struct {
	// GrantType identifies the OAuth2 grant (one of the Grant* constants).
	GrantType string

	// Common OAuth2 client authentication fields.
	ClientID     string
	ClientSecret string
	Scopes       []string

	// Authorization Code grant fields.
	Code        string
	RedirectURI string

	// PKCE field (Authorization Code + PKCE grant).
	CodeVerifier string

	// Password grant fields.
	Username string
	Password string

	// Refresh Token grant field.
	RefreshToken string

	// Personal Access grant fields.
	UserID    string
	TokenName string

	// Device Authorization grant field.
	DeviceCode string

	// Provider selects the auth provider for multi-provider applications.
	Provider string
}

// IssuedToken is the result of a successful token issuance from the AuthorizationServer.
type IssuedToken struct {
	TokenType    string
	AccessToken  string
	ExpiresIn    int
	RefreshToken string // empty for grants that do not issue refresh tokens
	TokenID      string // database row ID
}

// AuthorizationRequest represents the query parameters of GET /oauth/authorize.
type AuthorizationRequest struct {
	ClientID            string
	RedirectURI         string
	State               string
	Scopes              []string
	ResponseType        string
	CodeChallenge       string
	CodeChallengeMethod string
}

// ValidatedAuthRequest holds the validated state that bridges the GET /oauth/authorize
// step and the POST /oauth/authorize approval step.
type ValidatedAuthRequest struct {
	Client      *Client
	Scopes      []Scope
	User        interface{ GetAuthIdentifier() string }
	State       string
	RedirectURI string

	// BackendState is opaque backend-specific state serialised to the session
	// between the authorize and approve/deny steps.
	BackendState []byte
}

// AuthorizationResponse is returned by CompleteAuthorizationRequest.
// The RedirectURI includes the authorization code or error parameters.
type AuthorizationResponse struct {
	// RedirectURI is the fully-formed redirect URL (includes code= or error= params).
	RedirectURI string
}
