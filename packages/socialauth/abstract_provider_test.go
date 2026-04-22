package socialite_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/bedrock/packages/socialauth"
)

// ── Test doubles ─────────────────────────────────────────────────────────────

// testSession is a simple in-memory session store for tests.
type testSession struct {
	data map[string]any
}

// oauthTwoTestProvider is the Go equivalent of OAuthTwoTestProviderStub.php.
type oauthTwoTestProvider struct {
	socialauth.AbstractProvider
}

// oauthTwoWithPKCETestProvider mirrors OAuthTwoWithPKCETestProviderStub.php.
type oauthTwoWithPKCETestProvider struct {
	socialauth.AbstractProvider
}

// roundTripper is a function-based http.RoundTripper for test HTTP injection.
type roundTripper func(*http.Request) (*http.Response, error)

func newTestSession() *testSession {
	return &testSession{data: make(map[string]any)}
}

func (s *testSession) Put(key string, value any) { s.data[key] = value }
func (s *testSession) Get(key string) any        { return s.data[key] }
func (s *testSession) Pull(key string) any {
	v := s.data[key]
	delete(s.data, key)

	return v
}

func newTestProvider(req *http.Request, session socialauth.Session, clientID, clientSecret, redirectURL string) *oauthTwoTestProvider {
	p := &oauthTwoTestProvider{}
	p.AbstractProvider = socialauth.NewAbstractProvider(p, req, session, clientID, clientSecret, redirectURL)

	return p
}

func (p *oauthTwoTestProvider) GetAuthURL(state *string) string {
	return p.BuildAuthURLFromBase("http://auth.url", state)
}

func (p *oauthTwoTestProvider) GetTokenURL() string { return "http://token.url" }

func (p *oauthTwoTestProvider) GetUserByToken(_ context.Context, _ string) (map[string]any, error) {
	return map[string]any{"id": "foo"}, nil
}

func (p *oauthTwoTestProvider) MapUserToObject(raw map[string]any) *socialauth.User {
	u := &socialauth.User{}
	u.ID = fmt.Sprintf("%v", raw["id"])

	return u
}

func newPKCETestProvider(req *http.Request, session socialauth.Session, clientID, clientSecret, redirectURL string) *oauthTwoWithPKCETestProvider {
	p := &oauthTwoWithPKCETestProvider{}
	p.AbstractProvider = socialauth.NewAbstractProvider(p, req, session, clientID, clientSecret, redirectURL)
	p.EnablePKCE()

	return p
}

func (p *oauthTwoWithPKCETestProvider) GetAuthURL(state *string) string {
	return p.BuildAuthURLFromBase("http://auth.url", state)
}
func (p *oauthTwoWithPKCETestProvider) GetTokenURL() string { return "http://token.url" }
func (p *oauthTwoWithPKCETestProvider) GetUserByToken(_ context.Context, _ string) (map[string]any, error) {
	return map[string]any{"id": "foo"}, nil
}
func (p *oauthTwoWithPKCETestProvider) MapUserToObject(raw map[string]any) *socialauth.User {
	u := &socialauth.User{}
	u.ID = fmt.Sprintf("%v", raw["id"])

	return u
}

func (f roundTripper) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func fakeHTTPClient(body string, status int) *http.Client {
	return &http.Client{
		Transport: roundTripper(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: status,
				Body:       io.NopCloser(strings.NewReader(body)),
			}, nil
		}),
	}
}

// ── Tests (OAuthTwoTest.php equivalents) ─────────────────────────────────────

// OAuthTwoTest::testRedirectGeneratesTheProperFrameworkRedirectResponseWithoutPKCE
// TestRedirectBuildsURLWithoutPKCE mirrors
// testRedirectGeneratesTheProperFrameworkRedirectResponseWithoutPKCE.
func TestRedirectBuildsURLWithoutPKCE(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	session := newTestSession()

	provider := newTestProvider(req, session, "client_id", "client_secret", "redirect")
	redirectURL, err := provider.Redirect(context.Background())

	if err != nil {
		t.Fatalf("Redirect() returned error: %v", err)
	}

	state, _ := session.data["state"].(string) // was Put during Redirect

	if state == "" {
		// state was already pulled — grab from URL instead
		parsed, _ := url.Parse(redirectURL)
		state = parsed.Query().Get("state")
	}

	if state == "" {
		t.Fatal("expected state to be stored in session")
	}

	want := "http://auth.url?client_id=client_id&redirect_uri=redirect&response_type=code&scope=&state=" + state

	if redirectURL != want {
		t.Errorf("redirect URL mismatch\ngot:  %s\nwant: %s", redirectURL, want)
	}
}

// OAuthTwoTest::testRedirectGeneratesTheProperFrameworkRedirectResponseWithPKCE
// TestRedirectBuildsURLWithPKCE mirrors
// testRedirectGeneratesTheProperFrameworkRedirectResponseWithPKCE.
func TestRedirectBuildsURLWithPKCE(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	session := newTestSession()

	provider := newPKCETestProvider(req, session, "client_id", "client_secret", "redirect")
	redirectURL, err := provider.Redirect(context.Background())

	if err != nil {
		t.Fatalf("Redirect() returned error: %v", err)
	}

	parsed, _ := url.Parse(redirectURL)
	q := parsed.Query()

	if q.Get("code_challenge_method") != "S256" {
		t.Errorf("expected code_challenge_method=S256, got %q", q.Get("code_challenge_method"))
	}

	if q.Get("code_challenge") == "" {
		t.Error("expected code_challenge to be set")
	}

	if q.Get("state") == "" {
		t.Error("expected state to be set")
	}
}

// OAuthTwoTest::testTokenRequestIncludesPKCECodeVerifier
// TestTokenRequestIncludesPKCECodeVerifier mirrors
// testTokenRequestIncludesPKCECodeVerifier.
func TestTokenRequestIncludesPKCECodeVerifier(t *testing.T) {
	state := strings.Repeat("A", 40)
	rawURL := "http://example.com/callback?state=" + state + "&code=code"
	req, _ := http.NewRequest(http.MethodGet, rawURL, nil)

	session := newTestSession()
	// Simulate what redirect stored: state and code_verifier.
	session.Put("state", state)
	codeVerifier := "test-verifier-value-that-is-at-least-32-chars-long-xxxx"
	session.Put("code_verifier", codeVerifier)

	var capturedBody string
	provider := newPKCETestProvider(req, session, "client_id", "client_secret", "redirect_uri")
	provider.HTTP = &http.Client{
		Transport: roundTripper(func(r *http.Request) (*http.Response, error) {
			b, _ := io.ReadAll(r.Body)
			capturedBody = string(b)
			body := `{"access_token":"access_token","refresh_token":"refresh_token","expires_in":3600}`

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
			}, nil
		}),
	}

	user, err := provider.User(context.Background())

	if err != nil {
		t.Fatalf("User() returned error: %v", err)
	}

	if user.ID != "foo" {
		t.Errorf("expected user ID 'foo', got %q", user.ID)
	}

	if user.Token != "access_token" {
		t.Errorf("expected token 'access_token', got %q", user.Token)
	}

	if user.RefreshToken != "refresh_token" {
		t.Errorf("expected refresh_token 'refresh_token', got %q", user.RefreshToken)
	}

	if user.ExpiresIn != 3600 {
		t.Errorf("expected expires_in 3600, got %d", user.ExpiresIn)
	}

	// Verify code_verifier was included in the POST body.
	vals, _ := url.ParseQuery(capturedBody)

	if vals.Get("code_verifier") != codeVerifier {
		t.Errorf("expected code_verifier %q in token request, got %q", codeVerifier, vals.Get("code_verifier"))
	}
}

// OAuthTwoTest::testUserReturnsAUserInstanceForTheAuthenticatedRequest
// OAuthTwoTest::testUserRefreshesToken
// TestUserReturnsAuthenticatedUser mirrors
// testUserReturnsAUserInstanceForTheAuthenticatedRequest.
func TestUserReturnsAuthenticatedUser(t *testing.T) {
	state := strings.Repeat("A", 40)
	rawURL := "http://example.com/callback?state=" + state + "&code=code"
	req, _ := http.NewRequest(http.MethodGet, rawURL, nil)

	session := newTestSession()
	session.Put("state", state)

	provider := newTestProvider(req, session, "client_id", "client_secret", "redirect_uri")
	provider.HTTP = fakeHTTPClient(
		`{"access_token":"access_token","refresh_token":"refresh_token","expires_in":3600}`,
		http.StatusOK,
	)

	user, err := provider.User(context.Background())

	if err != nil {
		t.Fatalf("User() returned error: %v", err)
	}

	if user.ID != "foo" {
		t.Errorf("expected ID 'foo', got %q", user.ID)
	}

	if user.Token != "access_token" {
		t.Errorf("expected token 'access_token', got %q", user.Token)
	}

	if user.RefreshToken != "refresh_token" {
		t.Errorf("expected refresh_token 'refresh_token', got %q", user.RefreshToken)
	}

	if user.ExpiresIn != 3600 {
		t.Errorf("expected expires_in 3600, got %d", user.ExpiresIn)
	}

	// Second call should return the cached user.
	user2, _ := provider.User(context.Background())

	if user2 != user {
		t.Error("expected second User() call to return cached instance")
	}
}

// OAuthTwoTest::testExceptionIsThrownIfStateIsInvalid
// TestUserErrorsOnInvalidState mirrors testExceptionIsThrownIfStateIsInvalid.
func TestUserErrorsOnInvalidState(t *testing.T) {
	rawURL := "http://example.com/callback?state=" + strings.Repeat("B", 40) + "&code=code"
	req, _ := http.NewRequest(http.MethodGet, rawURL, nil)

	session := newTestSession()
	session.Put("state", strings.Repeat("A", 40)) // different from request

	provider := newTestProvider(req, session, "client_id", "client_secret", "redirect")
	_, err := provider.User(context.Background())

	if err != socialauth.ErrInvalidState {
		t.Errorf("expected ErrInvalidState, got %v", err)
	}
}

// OAuthTwoTest::testExceptionIsThrownIfStateIsNotSet
// TestUserErrorsOnMissingState mirrors testExceptionIsThrownIfStateIsNotSet.
func TestUserErrorsOnMissingState(t *testing.T) {
	rawURL := "http://example.com/callback?state=somestate&code=code"
	req, _ := http.NewRequest(http.MethodGet, rawURL, nil)

	session := newTestSession() // no state stored

	provider := newTestProvider(req, session, "client_id", "client_secret", "redirect")
	_, err := provider.User(context.Background())

	if err != socialauth.ErrInvalidState {
		t.Errorf("expected ErrInvalidState, got %v", err)
	}
}

// OAuthTwoTest::testCanGetAuthUrl
// TestGetAuthURL mirrors testCanGetAuthUrl.
func TestGetAuthURL(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	session := newTestSession()

	provider := newTestProvider(req, session, "client_id", "client_secret", "redirect")
	got := provider.GetAuthURL(nil)
	want := "http://auth.url?client_id=client_id&redirect_uri=redirect&response_type=code&scope="

	if got != want {
		t.Errorf("GetAuthURL mismatch\ngot:  %s\nwant: %s", got, want)
	}
}

// OAuthTwoTest::testCanGetStatelessAuthUrl
// TestGetStatelessAuthURL mirrors testCanGetStatelessAuthUrl.
func TestGetStatelessAuthURL(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	session := newTestSession()

	provider := newTestProvider(req, session, "client_id", "client_secret", "redirect")
	provider.Stateless()
	got := provider.GetAuthURL(nil)
	want := "http://auth.url?client_id=client_id&redirect_uri=redirect&response_type=code&scope="

	if got != want {
		t.Errorf("stateless GetAuthURL mismatch\ngot:  %s\nwant: %s", got, want)
	}
}

// OAuthTwoTest::testCanGetStatelessAuthUrl
// TestStatelessRedirectOmitsState verifies that Redirect() in stateless mode
// does not include a state parameter in the URL.
func TestStatelessRedirectOmitsState(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	session := newTestSession()

	provider := newTestProvider(req, session, "client_id", "client_secret", "redirect")
	provider.Stateless()
	redirectURL, err := provider.Redirect(context.Background())

	if err != nil {
		t.Fatalf("Redirect() returned error: %v", err)
	}

	parsed, _ := url.Parse(redirectURL)

	if parsed.Query().Get("state") != "" {
		t.Error("stateless redirect should not include state parameter")
	}

	if session.data["state"] != nil {
		t.Error("stateless redirect should not store state in session")
	}
}
