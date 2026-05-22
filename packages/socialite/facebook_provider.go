package socialite

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

// FacebookProvider handles OAuth2 authentication via Facebook Graph API.
type FacebookProvider struct {
	AbstractProvider
	graphVersion string
	fields       []string
}

// NewFacebookProvider constructs a FacebookProvider.
func NewFacebookProvider(req *http.Request, session Session, clientID, clientSecret, redirectURL string) *FacebookProvider {
	f := &FacebookProvider{
		graphVersion: "v3.3",
		fields:       []string{"name", "email", "gender", "verified", "link"},
	}
	f.AbstractProvider = NewAbstractProvider(f, req, session, clientID, clientSecret, redirectURL)
	f.scopes = []string{"email"}
	f.scopeSep = ","

	return f
}

// UsingGraphVersion overrides the Graph API version (default "v3.3").
func (f *FacebookProvider) UsingGraphVersion(version string) *FacebookProvider {
	f.graphVersion = version

	return f
}

// Fields sets the user fields to request from the Graph API.
func (f *FacebookProvider) Fields(fields []string) *FacebookProvider {
	f.fields = fields

	return f
}

func (f *FacebookProvider) GetAuthURL(state *string) string {
	return f.BuildAuthURLFromBase("https://www.facebook.com/dialog/oauth", state)
}

func (f *FacebookProvider) GetTokenURL() string {
	return fmt.Sprintf("https://graph.facebook.com/%s/oauth/access_token", f.graphVersion)
}

func (f *FacebookProvider) GetUserByToken(ctx context.Context, token string) (map[string]any, error) {
	appSecretProof := f.appSecretProof(token)
	endpoint := fmt.Sprintf(
		"https://graph.facebook.com/%s/me?access_token=%s&appsecret_proof=%s&fields=%s",
		f.graphVersion, token, appSecretProof, strings.Join(f.fields, ","),
	)

	return f.getJSON(ctx, endpoint, nil)
}

func (f *FacebookProvider) MapUserToObject(raw map[string]any) *User {
	u := &User{}
	u.ID = stringify(raw["id"])
	u.Nickname = stringify(raw["id"])
	u.Name = stringify(raw["name"])
	u.Email = stringify(raw["email"])
	u.Avatar = fmt.Sprintf("https://graph.facebook.com/%s/picture?type=normal", stringify(raw["id"]))

	if link, ok := raw["link"].(string); ok {
		if u.Attributes == nil {
			u.Attributes = make(map[string]any)
		}

		u.Attributes["link"] = link
	}

	return u
}

func (f *FacebookProvider) appSecretProof(token string) string {
	mac := hmac.New(sha256.New, []byte(f.clientSecret))
	mac.Write([]byte(token))

	return hex.EncodeToString(mac.Sum(nil))
}

func (f *FacebookProvider) getJSON(ctx context.Context, targetURL string, headers map[string]string) (map[string]any, error) {
	return f.AbstractProvider.getJSON(ctx, targetURL, headers)
}
