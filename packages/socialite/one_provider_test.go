package socialite_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/bedrock/packages/socialite"
)

// ── OAuth1 test doubles ───────────────────────────────────────────────────────

// mockOAuth1Server is the Go equivalent of OAuthOneTestProviderStub.php,
// which mocked League\OAuth1\Client\Server\Twitter.
type mockOAuth1Server struct {
	tempCreds  *socialite.TemporaryCredentials
	tokenCreds *socialite.TokenCredentials
	authURL    string
	rawUser    map[string]any
	// record calls
	gotTempCreds  bool
	gotTokenCreds bool
}

func (m *mockOAuth1Server) GetTemporaryCredentials(_ context.Context) (*socialite.TemporaryCredentials, error) {
	m.gotTempCreds = true

	return m.tempCreds, nil
}

func (m *mockOAuth1Server) GetAuthorizationURL(temp *socialite.TemporaryCredentials) string {
	return m.authURL
}

func (m *mockOAuth1Server) GetTokenCredentials(_ context.Context, temp *socialite.TemporaryCredentials, oauthToken, verifier string) (*socialite.TokenCredentials, error) {
	m.gotTokenCreds = true

	return m.tokenCreds, nil
}

func (m *mockOAuth1Server) GetUserDetails(_ context.Context, _ *socialite.TokenCredentials) (map[string]any, error) {
	return m.rawUser, nil
}

func (m *mockOAuth1Server) MapUserToObject(raw map[string]any, token *socialite.TokenCredentials) *socialite.User {
	u := &socialite.User{}
	u.ID = stringify(raw["uid"])
	u.Email = stringify(raw["email"])
	u.Raw = raw["extra"].(map[string]any)
	u.SetOAuthToken(token.Identifier, token.Secret)

	return u
}

// stringify converts any to string (duplicate of package helper, for test use).
func stringify(v any) string {
	if v == nil {
		return ""
	}

	if s, ok := v.(string); ok {
		return s
	}

	return ""
}

// ── Tests (OAuthOneTest.php equivalents) ─────────────────────────────────────

// TestOAuth1RedirectGeneratesURL mirrors
// testRedirectGeneratesTheProperIlluminateRedirectResponse.
func TestOAuth1RedirectGeneratesURL(t *testing.T) {
	server := &mockOAuth1Server{
		tempCreds: &socialite.TemporaryCredentials{Identifier: "id", Secret: "secret"},
		authURL:   "http://auth.url",
	}
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	session := newTestSession()

	provider := socialite.NewOneAbstractProvider(server, req, session)
	redirectURL, err := provider.Redirect(context.Background())

	if err != nil {
		t.Fatalf("Redirect() returned error: %v", err)
	}

	if redirectURL != "http://auth.url" {
		t.Errorf("expected redirect to http://auth.url, got %q", redirectURL)
	}

	if !server.gotTempCreds {
		t.Error("expected GetTemporaryCredentials to be called")
	}
	// Temporary credentials must be stored in session.
	stored := session.Get("oauth.temp")

	if stored == nil {
		t.Error("expected temporary credentials to be stored in session under 'oauth.temp'")
	}
}

// TestOAuth1UserReturnsAuthenticatedUser mirrors
// testUserReturnsAUserInstanceForTheAuthenticatedRequest.
func TestOAuth1UserReturnsAuthenticatedUser(t *testing.T) {
	rawURL := "http://example.com/callback?oauth_token=oauth_token&oauth_verifier=oauth_verifier"
	req, _ := http.NewRequest(http.MethodGet, rawURL, nil)
	session := newTestSession()

	temp := &socialite.TemporaryCredentials{Identifier: "temp_id", Secret: "temp_secret"}
	session.Put("oauth.temp", temp)

	server := &mockOAuth1Server{
		tokenCreds: &socialite.TokenCredentials{Identifier: "identifier", Secret: "secret"},
		rawUser: map[string]any{
			"uid":   "uid",
			"email": "foo@bar.com",
			"extra": map[string]any{"extra": "extra"},
		},
	}

	provider := socialite.NewOneAbstractProvider(server, req, session)
	user, err := provider.User(context.Background())

	if err != nil {
		t.Fatalf("User() returned error: %v", err)
	}

	if user.ID != "uid" {
		t.Errorf("expected ID 'uid', got %q", user.ID)
	}

	if user.Email != "foo@bar.com" {
		t.Errorf("expected email 'foo@bar.com', got %q", user.Email)
	}

	if user.Token != "identifier" {
		t.Errorf("expected Token 'identifier', got %q", user.Token)
	}

	if user.TokenSecret != "secret" {
		t.Errorf("expected TokenSecret 'secret', got %q", user.TokenSecret)
	}
}

// TestOAuth1ErrorsOnMissingVerifier mirrors
// testExceptionIsThrownWhenVerifierIsMissing.
func TestOAuth1ErrorsOnMissingVerifier(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/callback", nil)
	session := newTestSession()
	server := &mockOAuth1Server{}

	provider := socialite.NewOneAbstractProvider(server, req, session)
	_, err := provider.User(context.Background())

	if err != socialite.ErrMissingVerifier {
		t.Errorf("expected ErrMissingVerifier, got %v", err)
	}
}

// TestOAuth1ErrorsOnMissingTemporaryCredentials mirrors
// testExceptionIsThrownWhenTemporaryCredentialsAreMissing.
func TestOAuth1ErrorsOnMissingTemporaryCredentials(t *testing.T) {
	rawURL := "http://example.com/callback?oauth_token=oauth_token&oauth_verifier=oauth_verifier"
	req, _ := http.NewRequest(http.MethodGet, rawURL, nil)
	session := newTestSession() // no "oauth.temp" stored
	server := &mockOAuth1Server{}

	provider := socialite.NewOneAbstractProvider(server, req, session)
	_, err := provider.User(context.Background())

	if err != socialite.ErrMissingTemporaryCredentials {
		t.Errorf("expected ErrMissingTemporaryCredentials, got %v", err)
	}
}
