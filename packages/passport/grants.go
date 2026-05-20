package passport

import "context"

// Grant type identifiers — matches Passport grant names exactly.
// These correspond to the grant type strings used in token requests.

// GrantAuthorizationCode is the standard OAuth2 Authorization Code grant.

// GrantAuthorizationCodePKCE is the Authorization Code grant with PKCE
// (Proof Key for Code Exchange), used for public clients like SPAs and mobile apps.

// GrantPassword is the Resource Owner Password Credentials grant.
// Use only for highly trusted first-party applications.

// GrantClientCredentials is the Client Credentials grant for machine-to-machine auth.

// GrantImplicit is the (deprecated) Implicit grant for browser-based apps.

// GrantRefreshToken is the Refresh Token grant for obtaining new access tokens.

// GrantDeviceCode is the Device Authorization grant (RFC 8628) for input-constrained devices.

// GrantPersonalAccess is the custom Passport grant for personal access tokens.
// It is not a standard OAuth2 grant type.

// GrantDriver handles a specific OAuth2 grant flow.
// The AuthorizationServer delegates token issuance to registered GrantDrivers.
//
// This mirrors the League\OAuth2\Server\Grant\AbstractGrant contract.
type GrantDriver interface {
	// GrantType returns the identifier for this grant (one of the Grant* constants).
	GrantType() string

	// RespondToAccessTokenRequest processes a token endpoint request for this grant.
	RespondToAccessTokenRequest(ctx context.Context, req TokenRequest) (*IssuedToken, error)
}

// AuthorizationCodeDriver extends GrantDriver for flows that involve a browser
// redirect: Authorization Code, Authorization Code + PKCE, and Implicit.
type AuthorizationCodeDriver interface {
	GrantDriver

	// ValidateAuthorizationRequest validates the initial GET /oauth/authorize request.
	ValidateAuthorizationRequest(ctx context.Context, req AuthorizationRequest) (*ValidatedAuthRequest, error)

	// CompleteAuthorizationRequest finishes the flow after the user approves or denies.
	CompleteAuthorizationRequest(ctx context.Context, req *ValidatedAuthRequest, approved bool) (*AuthorizationResponse, error)
}

// UserValidator validates user credentials for the Password grant.
// Implementations call the UserProvider and verify the password.
type UserValidator interface {
	ValidateCredentials(ctx context.Context, username, password, provider string) (string, error)
}

// DeviceAuthorizationDriver extends GrantDriver for the Device Authorization flow.
type DeviceAuthorizationDriver interface {
	GrantDriver

	// CreateDeviceAuthorization creates a new device code and user code pair.
	CreateDeviceAuthorization(ctx context.Context, clientID string, scopes []string) (*DeviceCode, error)

	// ApproveDeviceAuthorization marks a device code as approved by the given user.
	ApproveDeviceAuthorization(ctx context.Context, userCode string, userID string) error

	// DenyDeviceAuthorization marks a device code as denied.
	DenyDeviceAuthorization(ctx context.Context, userCode string) error
}

const (
	GrantAuthorizationCode = "authorization_code"

	GrantAuthorizationCodePKCE = "authorization_code_pkce"

	GrantPassword = "password"

	GrantClientCredentials = "client_credentials"

	GrantImplicit = "implicit"

	GrantRefreshToken = "refresh_token"

	GrantDeviceCode = "urn:ietf:params:oauth:grant-type:device_code"

	GrantPersonalAccess = "personal_access"
)
