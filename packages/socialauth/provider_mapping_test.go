package socialite_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/bedrock/packages/socialauth"
)

func unsignedJWT(t *testing.T, claims map[string]any) string {
	t.Helper()

	header, err := json.Marshal(map[string]any{"alg": "none", "typ": "JWT"})

	if err != nil {
		t.Fatal(err)
	}

	payload, err := json.Marshal(claims)

	if err != nil {
		t.Fatal(err)
	}

	return base64.RawURLEncoding.EncodeToString(header) + "." +
		base64.RawURLEncoding.EncodeToString(payload) + ".signature"
}

func fakeHTTPClientFunc(fn func(*http.Request) (*http.Response, error)) *http.Client {
	return &http.Client{Transport: roundTripper(fn)}
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

// GoogleProviderTest::test_it_can_map_a_user_from_an_access_token
func TestGoogleProviderMapsUserFromAccessToken(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	provider := socialauth.NewGoogleProvider(req, session, "client-id", "client-secret", "http://callback")
	provider.HTTP = fakeHTTPClient(`{"sub":"google-123","nickname":"g-user","name":"Google User","email":"user@example.com","picture":"https://example.com/google.png","email_verified":true,"profile":"https://profiles.example.com/google-123"}`, http.StatusOK)

	user, err := provider.UserFromToken(context.Background(), "access-token")

	if err != nil {
		t.Fatalf("UserFromToken() returned error: %v", err)
	}

	if user.ID != "google-123" {
		t.Fatalf("unexpected user ID: got %q want %q", user.ID, "google-123")
	}

	if user.Email != "user@example.com" {
		t.Fatalf("unexpected email: got %q want %q", user.Email, "user@example.com")
	}

	if got := user.Attributes["avatar_original"]; got != "https://example.com/google.png" {
		t.Fatalf("unexpected avatar_original: got %v", got)
	}
}

// GoogleProviderIdTokenTest::test_it_can_detect_jwt_tokens
func TestGoogleProviderCanDetectJWTTokens(t *testing.T) {
	token := unsignedJWT(t, map[string]any{"sub": "google-123"})

	if !socialauth.IsJWT(token) {
		t.Fatal("expected token to be detected as JWT")
	}

	if socialauth.IsJWT("access-token") {
		t.Fatal("expected opaque access token not to be detected as JWT")
	}
}

// GoogleProviderIdTokenTest::test_it_uses_jwt_verification_for_id_tokens
// GoogleProviderIdTokenTest::test_user_mapping_works_with_id_token_format
func TestGoogleProviderMapsUserFromIDToken(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	provider := socialauth.NewGoogleProvider(req, session, "client-id", "client-secret", "http://callback")
	idToken := unsignedJWT(t, map[string]any{
		"sub":            "google-123",
		"name":           "Google User",
		"email":          "user@example.com",
		"picture":        "https://example.com/google.png",
		"email_verified": true,
		"profile":        "https://profiles.example.com/google-123",
	})

	user, err := provider.UserFromIDToken(context.Background(), idToken)

	if err != nil {
		t.Fatalf("UserFromIDToken() returned error: %v", err)
	}

	if user.ID != "google-123" || user.Email != "user@example.com" {
		t.Fatalf("unexpected user from ID token: %+v", user)
	}

	if got := user.Attributes["verified_email"]; got != true {
		t.Fatalf("verified_email = %#v, want true", got)
	}
}

// GoogleProviderIdTokenTest::test_it_falls_back_to_api_call_for_access_tokens
func TestGoogleProviderFallsBackToAPIForAccessTokens(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	provider := socialauth.NewGoogleProvider(req, session, "client-id", "client-secret", "http://callback")
	called := false
	provider.HTTP = fakeHTTPClientFunc(func(r *http.Request) (*http.Response, error) {
		called = true

		return jsonResponse(http.StatusOK, `{"sub":"google-123","name":"Google User"}`), nil
	})

	user, err := provider.UserFromToken(context.Background(), "access-token")

	if err != nil {
		t.Fatalf("UserFromToken() returned error: %v", err)
	}

	if !called {
		t.Fatal("expected opaque access token to call userinfo API")
	}

	if user.ID != "google-123" {
		t.Fatalf("unexpected user ID: %q", user.ID)
	}
}

// GoogleProviderIdTokenTest::test_it_handles_invalid_jwt_tokens
func TestGoogleProviderHandlesInvalidJWTTokens(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	provider := socialauth.NewGoogleProvider(req, session, "client-id", "client-secret", "http://callback")

	if _, err := provider.UserFromIDToken(context.Background(), "not.a.valid-jwt"); err == nil {
		t.Fatal("expected invalid JWT to return an error")
	}
}

// OAuthTwoTest::testUserRefreshesToken
// OAuthTwoTest::testUserRefreshesGoogleToken
func TestGoogleProviderRefreshesToken(t *testing.T) {
	state := strings.Repeat("A", 40)
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/callback?state="+state+"&code=code", nil)
	session := newTestSession()
	session.Put("state", state)

	provider := socialauth.NewGoogleProvider(req, session, "client-id", "client-secret", "http://callback")
	provider.HTTP = fakeHTTPClientFunc(func(r *http.Request) (*http.Response, error) {
		switch {
		case r.Method == http.MethodPost && r.URL.String() == provider.GetTokenURL():
			return jsonResponse(http.StatusOK, `{"access_token":"access-token","refresh_token":"refresh-token","expires_in":3600,"scope":"openid profile email"}`), nil
		case r.Method == http.MethodGet && r.URL.String() == "https://www.googleapis.com/oauth2/v3/userinfo":
			return jsonResponse(http.StatusOK, `{"sub":"google-123","nickname":"g-user","name":"Google User","email":"user@example.com","picture":"https://example.com/google.png","email_verified":true}`), nil
		default:
			return nil, fmt.Errorf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})

	user, err := provider.User(context.Background())

	if err != nil {
		t.Fatalf("User() returned error: %v", err)
	}

	if user.RefreshToken != "refresh-token" {
		t.Fatalf("unexpected refresh token: got %q want %q", user.RefreshToken, "refresh-token")
	}

	if user.ExpiresIn != 3600 {
		t.Fatalf("unexpected expires_in: got %d want %d", user.ExpiresIn, 3600)
	}

	if len(user.ApprovedScopes) != 3 {
		t.Fatalf("unexpected scopes: got %v", user.ApprovedScopes)
	}
}

// LinkedInProviderTest::test_it_can_map_a_user_without_an_email_address
func TestLinkedInProviderMapsUserWithoutEmailAddress(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	provider := socialauth.NewLinkedInProvider(req, session, "client-id", "client-secret", "http://callback")

	user := provider.MapUserToObject(map[string]any{
		"id":                 "linkedin-123",
		"localizedFirstName": "Linked",
		"localizedLastName":  "In",
		"profilePicture": map[string]any{
			"displayImage~": map[string]any{
				"elements": []any{
					map[string]any{
						"identifiers": []any{
							map[string]any{"identifier": "https://example.com/linkedin.png"},
						},
					},
				},
			},
		},
	})

	if user.Email != "" {
		t.Fatalf("unexpected email: got %q want empty", user.Email)
	}

	if user.Avatar != "https://example.com/linkedin.png" {
		t.Fatalf("unexpected avatar: got %q want %q", user.Avatar, "https://example.com/linkedin.png")
	}
}

// LinkedInOpenIdProviderTest::test_response
func TestLinkedInOpenIDProviderMapsUser(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	provider := socialauth.NewLinkedInProvider(req, session, "client-id", "client-secret", "http://callback")

	user := provider.MapUserToObject(map[string]any{
		"sub":     "linkedin-openid-123",
		"name":    "LinkedIn OpenID",
		"email":   "openid@example.com",
		"picture": "https://example.com/linkedin-openid.png",
	})

	if user.ID != "linkedin-openid-123" || user.Email != "openid@example.com" || user.Avatar == "" {
		t.Fatalf("unexpected OpenID user: %+v", user)
	}
}

// LinkedInOpenIdProviderTest::test_missing_email_and_avatar
func TestLinkedInOpenIDProviderMapsUserWithoutEmailAndAvatar(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	provider := socialauth.NewLinkedInProvider(req, session, "client-id", "client-secret", "http://callback")

	user := provider.MapUserToObject(map[string]any{
		"sub":  "linkedin-openid-123",
		"name": "LinkedIn OpenID",
	})

	if user.ID != "linkedin-openid-123" {
		t.Fatalf("unexpected ID: %q", user.ID)
	}

	if user.Email != "" || user.Avatar != "" {
		t.Fatalf("expected missing email/avatar to map to empty strings: %+v", user)
	}
}

// OAuthTwoTest::testUserReturnsAUserInstanceForTheAuthenticatedFacebookRequest
func TestFacebookProviderReturnsAuthenticatedUser(t *testing.T) {
	state := strings.Repeat("B", 40)
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/callback?state="+state+"&code=code", nil)
	session := newTestSession()
	session.Put("state", state)

	provider := socialauth.NewFacebookProvider(req, session, "client-id", "client-secret", "http://callback")
	provider.HTTP = fakeHTTPClientFunc(func(r *http.Request) (*http.Response, error) {
		switch {
		case r.Method == http.MethodPost && r.URL.String() == provider.GetTokenURL():
			return jsonResponse(http.StatusOK, `{"access_token":"access-token","refresh_token":"refresh-token","expires_in":3600}`), nil
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.String(), "https://graph.facebook.com/"):
			return jsonResponse(http.StatusOK, `{"id":"facebook-123","name":"Facebook User","email":"user@example.com","link":"https://facebook.com/facebook-user"}`), nil
		default:
			return nil, fmt.Errorf("unexpected request: %s %s", r.Method, r.URL.String())
		}
	})

	user, err := provider.User(context.Background())

	if err != nil {
		t.Fatalf("User() returned error: %v", err)
	}

	if user.ID != "facebook-123" {
		t.Fatalf("unexpected user ID: got %q want %q", user.ID, "facebook-123")
	}

	if user.Token != "access-token" {
		t.Fatalf("unexpected token: got %q want %q", user.Token, "access-token")
	}

	if user.RefreshToken != "refresh-token" {
		t.Fatalf("unexpected refresh token: got %q want %q", user.RefreshToken, "refresh-token")
	}
}

// SlackOpenIdProviderTest::test_response
func TestSlackOpenIDProviderMapsUser(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	provider := socialauth.NewSlackProvider(req, session, "client-id", "client-secret", "http://callback")

	user := provider.MapUserToObject(map[string]any{
		"sub":     "slack-openid-123",
		"name":    "Slack OpenID",
		"email":   "slack@example.com",
		"picture": "https://example.com/slack-openid.png",
	})

	if user.ID != "slack-openid-123" || user.Email != "slack@example.com" || user.Avatar == "" {
		t.Fatalf("unexpected OpenID user: %+v", user)
	}
}

// SlackOpenIdProviderTest::test_missing_email_and_avatar
func TestSlackOpenIDProviderMapsUserWithoutEmailAndAvatar(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	session := newTestSession()
	provider := socialauth.NewSlackProvider(req, session, "client-id", "client-secret", "http://callback")

	user := provider.MapUserToObject(map[string]any{
		"sub":  "slack-openid-123",
		"name": "Slack OpenID",
	})

	if user.ID != "slack-openid-123" {
		t.Fatalf("unexpected ID: %q", user.ID)
	}

	if user.Email != "" || user.Avatar != "" {
		t.Fatalf("expected missing email/avatar to map to empty strings: %+v", user)
	}
}
