package passport

import (
	"sync"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// Passport is the central configuration and registrar for the passport package.
// It mirrors the static API of Passport Passport facade.
//
// Create one via NewPassport and pass it throughout your application to configure
// scopes, expiry times, and grant type availability.
type Passport struct {
	mu     sync.RWMutex
	config *PassportConfig

	scopes          map[string]Scope
	inheritedScopes bool

	tokensTTL               time.Duration
	refreshTokensTTL        time.Duration
	personalAccessTokensTTL time.Duration
	clientCredentialsTTL    time.Duration

	// Grant type toggles — all disabled by default (explicit opt-in required).
	authCodeGrantEnabled          bool
	passwordGrantEnabled          bool
	clientCredentialsGrantEnabled bool
	implicitGrantEnabled          bool
	deviceCodeGrantEnabled        bool
	refreshTokenGrantEnabled      bool

	// Test overrides set by ActingAs / ActingAsClient.
	actingAsUser   cauth.Authenticatable
	actingAsToken  *AccessToken
	actingAsClient *Client
	actingAsScopes []string
}

const (
	defaultTokensTTL               = 365 * 24 * time.Hour
	defaultRefreshTokensTTL        = 30 * 24 * time.Hour
	defaultPersonalAccessTokensTTL = 365 * 24 * time.Hour
	defaultClientCredentialsTTL    = 365 * 24 * time.Hour
)

// NewPassport creates a Passport with the given config and sensible defaults.
// When cfg is nil, an empty PassportConfig is used (useful for tests).
func NewPassport(cfg *PassportConfig) *Passport {
	if cfg == nil {
		cfg = &PassportConfig{Guard: "web"}
	}

	return &Passport{
		config:                  cfg,
		scopes:                  make(map[string]Scope),
		tokensTTL:               defaultTokensTTL,
		refreshTokensTTL:        defaultRefreshTokensTTL,
		personalAccessTokensTTL: defaultPersonalAccessTokensTTL,
		clientCredentialsTTL:    defaultClientCredentialsTTL,
	}
}

// ---- Scope registration -------------------------------------------------------

// TokensCan registers the available OAuth2 scopes for the application.
// The map key is the scope ID and the value is its human-readable description.
// Mirrors Passport::tokensCan().
func (p *Passport) TokensCan(scopes map[string]string) *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	for id, desc := range scopes {
		p.scopes[id] = Scope{ID: id, Description: desc}
	}

	return p
}

// Scopes returns all registered scopes as a slice.
func (p *Passport) Scopes() []Scope {
	p.mu.RLock()

	defer p.mu.RUnlock()

	out := make([]Scope, 0, len(p.scopes))

	for _, s := range p.scopes {
		out = append(out, s)
	}

	return out
}

// ScopeIDs returns the IDs of all registered scopes.
func (p *Passport) ScopeIDs() []string {
	p.mu.RLock()

	defer p.mu.RUnlock()

	ids := make([]string, 0, len(p.scopes))

	for id := range p.scopes {
		ids = append(ids, id)
	}

	return ids
}

// FindScope returns the Scope for the given ID, or nil if not registered.
// Mirrors Passport::scopes() in the upstream framework (which returns a Collection keyed by ID).
func (p *Passport) FindScope(id string) *Scope {
	p.mu.RLock()

	defer p.mu.RUnlock()

	if s, ok := p.scopes[id]; ok {
		return &s
	}

	return nil
}

// HasScope reports whether the given scope ID is registered.
func (p *Passport) HasScope(id string) bool {
	p.mu.RLock()

	defer p.mu.RUnlock()

	_, ok := p.scopes[id]

	return ok
}

// UseInheritedScopes enables or disables hierarchical scope resolution.
// When enabled, Can("user:read") returns true for a token that carries "user".
func (p *Passport) UseInheritedScopes(enabled bool) *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.inheritedScopes = enabled

	return p
}

// InheritedScopesEnabled reports whether hierarchical scope resolution is active.
func (p *Passport) InheritedScopesEnabled() bool {
	p.mu.RLock()

	defer p.mu.RUnlock()

	return p.inheritedScopes
}

// ---- Token expiry -------------------------------------------------------------

// TokensExpireIn sets the TTL for access tokens.
func (p *Passport) TokensExpireIn(d time.Duration) *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.tokensTTL = d

	return p
}

// TokensTTL returns the current access token TTL.
func (p *Passport) TokensTTL() time.Duration {
	p.mu.RLock()

	defer p.mu.RUnlock()

	return p.tokensTTL
}

// RefreshTokensExpireIn sets the TTL for refresh tokens.
func (p *Passport) RefreshTokensExpireIn(d time.Duration) *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.refreshTokensTTL = d

	return p
}

// RefreshTokensTTL returns the current refresh token TTL.
func (p *Passport) RefreshTokensTTL() time.Duration {
	p.mu.RLock()

	defer p.mu.RUnlock()

	return p.refreshTokensTTL
}

// PersonalAccessTokensExpireIn sets the TTL for personal access tokens.
func (p *Passport) PersonalAccessTokensExpireIn(d time.Duration) *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.personalAccessTokensTTL = d

	return p
}

// PersonalAccessTokensTTL returns the current personal access token TTL.
func (p *Passport) PersonalAccessTokensTTL() time.Duration {
	p.mu.RLock()

	defer p.mu.RUnlock()

	return p.personalAccessTokensTTL
}

// ClientCredentialsTokensExpireIn sets the TTL for client credentials tokens.
func (p *Passport) ClientCredentialsTokensExpireIn(d time.Duration) *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.clientCredentialsTTL = d

	return p
}

// ClientCredentialsTTL returns the current client credentials token TTL.
func (p *Passport) ClientCredentialsTTL() time.Duration {
	p.mu.RLock()

	defer p.mu.RUnlock()

	return p.clientCredentialsTTL
}

// ---- Grant type toggles -------------------------------------------------------

// EnableAuthorizationCodeGrant activates the Authorization Code grant.
func (p *Passport) EnableAuthorizationCodeGrant() *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.authCodeGrantEnabled = true

	return p
}

// DisableAuthorizationCodeGrant deactivates the Authorization Code grant.
func (p *Passport) DisableAuthorizationCodeGrant() *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.authCodeGrantEnabled = false

	return p
}

// EnablePasswordGrant activates the Resource Owner Password Credentials grant.
func (p *Passport) EnablePasswordGrant() *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.passwordGrantEnabled = true

	return p
}

// DisablePasswordGrant deactivates the Password grant.
func (p *Passport) DisablePasswordGrant() *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.passwordGrantEnabled = false

	return p
}

// EnableClientCredentialsGrant activates the Client Credentials grant.
func (p *Passport) EnableClientCredentialsGrant() *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.clientCredentialsGrantEnabled = true

	return p
}

// DisableClientCredentialsGrant deactivates the Client Credentials grant.
func (p *Passport) DisableClientCredentialsGrant() *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.clientCredentialsGrantEnabled = false

	return p
}

// EnableImplicitGrant activates the (deprecated) Implicit grant.
func (p *Passport) EnableImplicitGrant() *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.implicitGrantEnabled = true

	return p
}

// DisableImplicitGrant deactivates the Implicit grant.
func (p *Passport) DisableImplicitGrant() *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.implicitGrantEnabled = false

	return p
}

// EnableDeviceCodeGrant activates the Device Authorization grant (RFC 8628).
func (p *Passport) EnableDeviceCodeGrant() *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.deviceCodeGrantEnabled = true

	return p
}

// DisableDeviceCodeGrant deactivates the Device Authorization grant.
func (p *Passport) DisableDeviceCodeGrant() *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.deviceCodeGrantEnabled = false

	return p
}

// EnableRefreshTokenGrant activates the Refresh Token grant.
func (p *Passport) EnableRefreshTokenGrant() *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.refreshTokenGrantEnabled = true

	return p
}

// DisableRefreshTokenGrant deactivates the Refresh Token grant.
func (p *Passport) DisableRefreshTokenGrant() *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.refreshTokenGrantEnabled = false

	return p
}

// IsGrantEnabled reports whether the given grant type is currently enabled.
// grantType should be one of the Grant* constants defined in grants.go.
func (p *Passport) IsGrantEnabled(grantType string) bool {
	p.mu.RLock()

	defer p.mu.RUnlock()

	switch grantType {
	case GrantAuthorizationCode, GrantAuthorizationCodePKCE:
		return p.authCodeGrantEnabled
	case GrantPassword:
		return p.passwordGrantEnabled
	case GrantClientCredentials:
		return p.clientCredentialsGrantEnabled
	case GrantImplicit:
		return p.implicitGrantEnabled
	case GrantDeviceCode:
		return p.deviceCodeGrantEnabled
	case GrantRefreshToken:
		return p.refreshTokenGrantEnabled
	default:
		return false
	}
}

// Config returns the PassportConfig used to create this Passport.
func (p *Passport) Config() *PassportConfig {
	p.mu.RLock()

	defer p.mu.RUnlock()

	return p.config
}

// ---- Test helpers ------------------------------------------------------------

// ActingAs sets a test override so that the TokenGuard returns user with the
// given access token and scopes without consulting the token store.
// Mirrors Passport::actingAs().
func (p *Passport) ActingAs(user cauth.Authenticatable, token *AccessToken, scopes []string) *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.actingAsUser = user
	p.actingAsToken = token
	p.actingAsScopes = scopes

	return p
}

// ActingAsClient sets a test override so that the TokenGuard returns the given
// client (machine-to-machine authentication) without consulting the token store.
// Mirrors Passport::actingAsClient().
func (p *Passport) ActingAsClient(client *Client, scopes []string) *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.actingAsClient = client
	p.actingAsScopes = scopes

	return p
}

// ClearActing removes all test overrides.
func (p *Passport) ClearActing() *Passport {
	p.mu.Lock()

	defer p.mu.Unlock()

	p.actingAsUser = nil
	p.actingAsToken = nil
	p.actingAsClient = nil
	p.actingAsScopes = nil

	return p
}

// actingAs returns the acting-as state under a read lock.
// Returns (user, token, client, scopes, isSet).
func (p *Passport) actingAsState() (cauth.Authenticatable, *AccessToken, *Client, []string, bool) {
	p.mu.RLock()

	defer p.mu.RUnlock()

	if p.actingAsUser != nil || p.actingAsClient != nil {
		return p.actingAsUser, p.actingAsToken, p.actingAsClient, p.actingAsScopes, true
	}

	return nil, nil, nil, nil, false
}
