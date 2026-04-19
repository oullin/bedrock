package socialauth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
)

// TwitterProvider handles OAuth2 authentication via Twitter API v2.
// Twitter deprecated OAuth1 for new apps; this implementation uses OAuth2
// with PKCE (Proof Key for Code Exchange), mirroring
// Upstream\SocialAuth\Two\TwitterProvider.
type TwitterProvider struct {
	AbstractProvider
}

// NewTwitterProvider constructs a TwitterProvider with PKCE enabled by default.
func NewTwitterProvider(req *http.Request, session Session, clientID, clientSecret, redirectURL string) *TwitterProvider {
	t := &TwitterProvider{}
	t.AbstractProvider = NewAbstractProvider(t, req, session, clientID, clientSecret, redirectURL)
	t.scopes = []string{"users.read", "tweet.read"}
	t.scopeSep = " "
	t.usesPKCE = true
	t.encodingType = EncodingRFC3986

	return t
}

func (t *TwitterProvider) GetAuthURL(state *string) string {
	return t.BuildAuthURLFromBase("https://twitter.com/i/oauth2/authorize", state)
}

func (t *TwitterProvider) GetTokenURL() string {
	return "https://api.twitter.com/2/oauth2/token"
}

// GetUserByToken fetches the Twitter user via API v2 with bearer auth.
func (t *TwitterProvider) GetUserByToken(ctx context.Context, token string) (map[string]any, error) {
	resp, err := t.getJSON(ctx,
		"https://api.twitter.com/2/users/me?user.fields=profile_image_url",
		map[string]string{"Authorization": "Bearer " + token},
	)

	if err != nil {
		return nil, err
	}
	// Twitter wraps the payload in a "data" key.
	if data, ok := resp["data"].(map[string]any); ok {
		return data, nil
	}

	return resp, nil
}

func (t *TwitterProvider) MapUserToObject(raw map[string]any) *User {
	u := &User{}
	u.ID = stringify(raw["id"])
	u.Nickname = stringify(raw["username"])
	u.Name = stringify(raw["name"])
	u.Avatar = stringify(raw["profile_image_url"])

	return u
}

// FetchAccessToken overrides the default token exchange to use HTTP Basic Auth,
// as required by Twitter's API v2. It satisfies the TokenFetcher interface.
func (t *TwitterProvider) FetchAccessToken(ctx context.Context, code string) (map[string]any, error) {
	fields := t.GetTokenFields(code)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.GetTokenURL(), formBody(fields))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	creds := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", t.clientID, t.clientSecret)))
	req.Header.Set("Authorization", "Basic "+creds)

	resp, err := t.getHTTPClient().Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return decodeJSONBody(resp)
}

func (t *TwitterProvider) getJSON(ctx context.Context, targetURL string, headers map[string]string) (map[string]any, error) {
	return t.AbstractProvider.getJSON(ctx, targetURL, headers)
}
