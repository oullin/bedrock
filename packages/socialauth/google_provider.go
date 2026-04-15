package socialauth

import (
	"context"
	"net/http"
)

// GoogleProvider handles OAuth2 authentication via Google.
// It mirrors Upstream\SocialAuth\Two\GoogleProvider.
type GoogleProvider struct {
	AbstractProvider
}

// NewGoogleProvider constructs a GoogleProvider.
func NewGoogleProvider(req *http.Request, session Session, clientID, clientSecret, redirectURL string) *GoogleProvider {
	g := &GoogleProvider{}
	g.AbstractProvider = NewAbstractProvider(g, req, session, clientID, clientSecret, redirectURL)
	g.scopes = []string{"openid", "profile", "email"}
	g.scopeSep = " "
	return g
}

func (g *GoogleProvider) GetAuthURL(state *string) string {
	return g.BuildAuthURLFromBase("https://accounts.google.com/o/oauth2/auth", state)
}

func (g *GoogleProvider) GetTokenURL() string {
	return "https://www.googleapis.com/oauth2/v4/token"
}

func (g *GoogleProvider) GetUserByToken(ctx context.Context, token string) (map[string]any, error) {
	return g.getWithBearer(ctx, "https://www.googleapis.com/oauth2/v3/userinfo", token)
}

func (g *GoogleProvider) MapUserToObject(raw map[string]any) *User {
	u := &User{}
	u.ID = stringify(raw["sub"])
	u.Nickname = stringify(raw["nickname"])
	u.Name = stringify(raw["name"])
	u.Email = stringify(raw["email"])
	u.Avatar = stringify(raw["picture"])
	if u.Attributes == nil {
		u.Attributes = make(map[string]any)
	}
	u.Attributes["avatar_original"] = stringify(raw["picture"])
	u.Attributes["verified_email"] = raw["email_verified"]
	u.Attributes["link"] = raw["profile"]
	return u
}
