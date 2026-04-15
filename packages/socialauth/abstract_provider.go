package socialauth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
)

// HTTPDoer is the minimal HTTP interface needed by AbstractProvider.
// *http.Client satisfies it out of the box. Tests may inject a custom
// implementation to intercept outbound requests.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// ProviderHooks declares the four methods that each concrete OAuth2 provider
// must implement. AbstractProvider calls these via its impl field, mimicking
// PHP's abstract-method dispatch without inheritance.
type ProviderHooks interface {
	// GetAuthURL constructs the provider's authorization endpoint URL.
	// state is nil when the provider operates in stateless mode.
	GetAuthURL(state *string) string

	// GetTokenURL returns the provider's token endpoint URL.
	GetTokenURL() string

	// GetUserByToken fetches the raw user payload from the provider API
	// using the given access token.
	GetUserByToken(ctx context.Context, token string) (map[string]any, error)

	// MapUserToObject converts a raw provider payload to a *User.
	MapUserToObject(raw map[string]any) *User
}

const (
	// EncodingRFC1738 is the default query-string encoding (url.QueryEscape).
	EncodingRFC1738 = iota
	// EncodingRFC3986 uses url.PathEscape-compatible percent encoding.
	EncodingRFC3986
)

// AbstractProvider is the base OAuth2 provider. Concrete providers embed it
// and set the impl field to themselves so that abstract-method calls are
// dispatched correctly.
//
// It mirrors Upstream\SocialAuth\Two\AbstractProvider.
type AbstractProvider struct {
	impl ProviderHooks

	// HTTP is the client used for all outbound OAuth requests. Set this in
	// tests to intercept requests without making real network calls.
	// Defaults to http.DefaultClient.
	HTTP HTTPDoer

	session      Session
	request      *http.Request
	clientID     string
	clientSecret string
	redirectURL  string
	scopes       []string
	scopeSep     string
	stateless    bool
	usesPKCE     bool
	parameters   map[string]string
	cachedUser   *User
	encodingType int
}

// NewAbstractProvider creates an AbstractProvider and wires impl to the
// concrete provider (self-reference pattern for Go's "abstract class").
func NewAbstractProvider(
	impl ProviderHooks,
	req *http.Request,
	session Session,
	clientID, clientSecret, redirectURL string,
) AbstractProvider {
	return AbstractProvider{
		impl:         impl,
		request:      req,
		session:      session,
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		scopeSep:     " ",
		parameters:   make(map[string]string),
		encodingType: EncodingRFC1738,
	}
}

// ── Public fluent API ────────────────────────────────────────────────────────

// Stateless disables CSRF state validation. Useful for APIs that cannot
// maintain sessions. It mirrors AbstractProvider::stateless().
func (p *AbstractProvider) Stateless() *AbstractProvider {
	p.stateless = true
	return p
}

// Scopes merges additional scopes into the existing set.
// It mirrors AbstractProvider::scopes().
func (p *AbstractProvider) Scopes(scopes []string) *AbstractProvider {
	p.scopes = append(p.scopes, scopes...)
	return p
}

// SetScopes replaces the current scope list entirely.
// It mirrors AbstractProvider::setScopes().
func (p *AbstractProvider) SetScopes(scopes []string) *AbstractProvider {
	p.scopes = scopes
	return p
}

// GetScopes returns the current scope list.
func (p *AbstractProvider) GetScopes() []string { return p.scopes }

// RedirectURL overrides the redirect URI sent to the provider.
// It mirrors AbstractProvider::redirectUrl().
func (p *AbstractProvider) RedirectURL(u string) *AbstractProvider {
	p.redirectURL = u
	return p
}

// With merges extra query-string parameters into the authorization request.
// It mirrors AbstractProvider::with().
func (p *AbstractProvider) With(params map[string]string) *AbstractProvider {
	for k, v := range params {
		p.parameters[k] = v
	}
	return p
}

// EnablePKCE activates Proof Key for Code Exchange (RFC 7636) on the next
// redirect/user cycle. It mirrors AbstractProvider::enablePKCE().
func (p *AbstractProvider) EnablePKCE() *AbstractProvider {
	p.usesPKCE = true
	return p
}

// UsesPKCE reports whether PKCE is active.
func (p *AbstractProvider) UsesPKCE() bool { return p.usesPKCE }

// IsStateless reports whether state validation is disabled.
func (p *AbstractProvider) IsStateless() bool { return p.stateless }

// SetScopeSeparator sets the character used to join multiple scopes.
func (p *AbstractProvider) SetScopeSeparator(sep string) *AbstractProvider {
	p.scopeSep = sep
	return p
}

// SetRequest replaces the underlying HTTP request (e.g. for stateless use).
func (p *AbstractProvider) SetRequest(r *http.Request) *AbstractProvider {
	p.request = r
	return p
}

// ── Core OAuth2 flow ─────────────────────────────────────────────────────────

// Redirect returns the authorization URL the user should be redirected to.
// When stateful, it stores the state (and PKCE code verifier) in the session.
// It mirrors AbstractProvider::redirect().
func (p *AbstractProvider) Redirect(_ context.Context) (string, error) {
	var statePtr *string

	if !p.stateless {
		state := generateState()
		p.session.Put("state", state)
		statePtr = &state
	}

	if p.usesPKCE {
		verifier := generateCodeVerifier()
		p.session.Put("code_verifier", verifier)
	}

	return p.impl.GetAuthURL(statePtr), nil
}

// TokenFetcher is an optional interface that concrete providers implement to
// override the default OAuth2 token exchange (e.g. to use HTTP Basic Auth
// instead of sending credentials in the POST body as Twitter requires).
type TokenFetcher interface {
	FetchAccessToken(ctx context.Context, code string) (map[string]any, error)
}

// User completes the OAuth2 callback: validates state, exchanges the code for
// a token, fetches user info, and returns a populated *User.
// It mirrors AbstractProvider::user().
func (p *AbstractProvider) User(ctx context.Context) (*User, error) {
	if p.cachedUser != nil {
		return p.cachedUser, nil
	}

	if p.hasInvalidState() {
		return nil, ErrInvalidState
	}

	code := p.request.URL.Query().Get("code")

	var tokenResp map[string]any
	var err error
	if tf, ok := p.impl.(TokenFetcher); ok {
		tokenResp, err = tf.FetchAccessToken(ctx, code)
	} else {
		tokenResp, err = p.getAccessTokenResponse(ctx, code)
	}
	if err != nil {
		return nil, fmt.Errorf("socialauth: token exchange failed: %w", err)
	}

	token := mapStr(tokenResp, "access_token")
	raw, err := p.impl.GetUserByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("socialauth: user fetch failed: %w", err)
	}

	user := p.impl.MapUserToObject(raw).SetRaw(raw).SetToken(token)

	if rt := mapStr(tokenResp, "refresh_token"); rt != "" {
		user.SetRefreshToken(rt)
	}
	if ei := mapInt(tokenResp, "expires_in"); ei > 0 {
		user.SetExpiresIn(ei)
	} else if e := mapInt(tokenResp, "expires"); e > 0 {
		user.SetExpiresIn(e)
	}
	if scopes := mapStrSlice(tokenResp, "scope"); len(scopes) > 0 {
		user.SetApprovedScopes(scopes)
	}

	p.cachedUser = user
	return user, nil
}

// UserFromToken fetches user info using a known access token, bypassing the
// full OAuth2 code-exchange flow. It mirrors AbstractProvider::userFromToken().
func (p *AbstractProvider) UserFromToken(ctx context.Context, token string) (*User, error) {
	raw, err := p.impl.GetUserByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return p.impl.MapUserToObject(raw).SetRaw(raw).SetToken(token), nil
}

// ── Authorization URL helpers ────────────────────────────────────────────────

// BuildAuthURLFromBase constructs the full authorization URL from a base URL
// and a state pointer (nil = stateless). It mirrors
// AbstractProvider::buildAuthUrlFromBase().
func (p *AbstractProvider) BuildAuthURLFromBase(base string, state *string) string {
	fields := p.getCodeFields(state)

	if p.encodingType == EncodingRFC3986 {
		vals := url.Values{}
		for k, v := range fields {
			vals.Set(k, v)
		}
		return base + "?" + strings.ReplaceAll(vals.Encode(), "+", "%20")
	}

	vals := url.Values{}
	for k, v := range fields {
		vals.Set(k, v)
	}
	return base + "?" + vals.Encode()
}

// getCodeFields returns the query parameters sent in the authorization request.
// It mirrors AbstractProvider::getCodeFields().
func (p *AbstractProvider) getCodeFields(state *string) map[string]string {
	fields := map[string]string{
		"client_id":     p.clientID,
		"redirect_uri":  p.redirectURL,
		"scope":         strings.Join(p.scopes, p.scopeSep),
		"response_type": "code",
	}

	if state != nil {
		fields["state"] = *state
	}

	if p.usesPKCE {
		verifier, _ := p.session.Get("code_verifier").(string)
		fields["code_challenge"] = generateCodeChallenge(verifier)
		fields["code_challenge_method"] = "S256"
	}

	for k, v := range p.parameters {
		fields[k] = v
	}

	return fields
}

// ── Token exchange ───────────────────────────────────────────────────────────

// getAccessTokenResponse POSTs the authorization code to the token endpoint
// and returns the decoded JSON response. It mirrors
// AbstractProvider::getAccessTokenResponse().
func (p *AbstractProvider) getAccessTokenResponse(ctx context.Context, code string) (map[string]any, error) {
	fields := p.getTokenFields(code)
	return p.postForm(ctx, p.impl.GetTokenURL(), map[string]string{"Accept": "application/json"}, fields)
}

// GetTokenFields returns the form parameters for the token exchange POST.
// It mirrors AbstractProvider::getTokenFields().
func (p *AbstractProvider) GetTokenFields(code string) map[string]string {
	return p.getTokenFields(code)
}

func (p *AbstractProvider) getTokenFields(code string) map[string]string {
	fields := map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     p.clientID,
		"client_secret": p.clientSecret,
		"code":          code,
		"redirect_uri":  p.redirectURL,
	}

	if p.usesPKCE {
		verifier, _ := p.session.Pull("code_verifier").(string)
		fields["code_verifier"] = verifier
	}

	return fields
}

// ── State / PKCE ─────────────────────────────────────────────────────────────

// hasInvalidState reports whether the state parameter in the callback request
// does not match the value stored in the session. Returns false when stateless.
// It mirrors AbstractProvider::hasInvalidState().
func (p *AbstractProvider) hasInvalidState() bool {
	if p.stateless {
		return false
	}
	stored, _ := p.session.Pull("state").(string)
	incoming := p.request.URL.Query().Get("state")
	return stored == "" || stored != incoming
}

// generateState creates a 40-character random alphanumeric string.
// It mirrors Str::random(40) used for CSRF state tokens.
func generateState() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 40)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[n.Int64()]
	}
	return string(b)
}

// generateCodeVerifier creates a 96-character base64url-encoded random string
// for use as a PKCE code verifier. It mirrors AbstractProvider::getCodeVerifier().
func generateCodeVerifier() string {
	b := make([]byte, 72)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("socialauth: cannot generate PKCE code verifier: " + err.Error())
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// generateCodeChallenge computes the S256 PKCE code challenge from a verifier.
// It mirrors AbstractProvider::getCodeChallenge().
func generateCodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// ── HTTP helpers ─────────────────────────────────────────────────────────────

// postForm performs an HTTP POST with application/x-www-form-urlencoded body
// and returns the decoded JSON response body.
func (p *AbstractProvider) postForm(ctx context.Context, targetURL string, headers, formParams map[string]string) (map[string]any, error) {
	vals := url.Values{}
	for k, v := range formParams {
		vals.Set(k, v)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, strings.NewReader(vals.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := p.getHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("socialauth: provider returned HTTP %d: %s", resp.StatusCode, bytes.TrimSpace(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("socialauth: cannot decode token response: %w", err)
	}
	return result, nil
}

// getWithBearer performs a GET request with an Authorization: Bearer header.
func (p *AbstractProvider) getWithBearer(ctx context.Context, targetURL, token string) (map[string]any, error) {
	return p.getJSON(ctx, targetURL, map[string]string{"Authorization": "Bearer " + token})
}

// getJSON performs a GET request and returns the decoded JSON response.
func (p *AbstractProvider) getJSON(ctx context.Context, targetURL string, headers map[string]string) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := p.getHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("socialauth: provider returned HTTP %d: %s", resp.StatusCode, bytes.TrimSpace(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("socialauth: cannot decode API response: %w", err)
	}
	return result, nil
}

func (p *AbstractProvider) getHTTPClient() HTTPDoer {
	if p.HTTP != nil {
		return p.HTTP
	}
	return http.DefaultClient
}

// ── JSON helpers ─────────────────────────────────────────────────────────────

func mapStr(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func mapInt(m map[string]any, key string) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		}
	}
	return 0
}

func mapStrSlice(m map[string]any, key string) []string {
	v, ok := m[key]
	if !ok {
		return nil
	}
	switch s := v.(type) {
	case string:
		if s == "" {
			return nil
		}
		return strings.Fields(s)
	case []any:
		out := make([]string, 0, len(s))
		for _, item := range s {
			if str, ok := item.(string); ok {
				out = append(out, str)
			}
		}
		return out
	}
	return nil
}
