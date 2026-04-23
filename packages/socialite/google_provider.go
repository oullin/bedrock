package socialite

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// GoogleProvider handles OAuth2 authentication via Google.
// It mirrors Laravel\Socialite\Two\GoogleProvider.
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

// UserFromIDToken maps an OpenID Connect ID token after verifying it with
// Google's tokeninfo endpoint. The endpoint validates the token signature,
// format, and expiry server-side; this method additionally enforces that the
// audience matches the configured client ID and the issuer is Google.
func (g *GoogleProvider) UserFromIDToken(ctx context.Context, idToken string) (*User, error) {
	raw, err := g.verifyIDToken(ctx, idToken)

	if err != nil {
		return nil, err
	}

	return g.MapUserToObject(raw).SetRaw(raw).SetToken(idToken), nil
}

func (g *GoogleProvider) verifyIDToken(ctx context.Context, idToken string) (map[string]any, error) {
	raw, err := g.postForm(ctx, "https://oauth2.googleapis.com/tokeninfo", nil, map[string]string{
		"id_token": idToken,
	})

	if err != nil {
		return nil, fmt.Errorf("socialite: invalid id token: %w", err)
	}

	if aud := stringify(raw["aud"]); aud != g.clientID {
		return nil, fmt.Errorf("socialite: invalid id token: audience %q does not match client id", aud)
	}

	iss := stringify(raw["iss"])

	if iss != "accounts.google.com" && iss != "https://accounts.google.com" {
		return nil, fmt.Errorf("socialite: invalid id token: unexpected issuer %q", iss)
	}

	if expired, err := idTokenExpired(raw["exp"]); err != nil {
		return nil, fmt.Errorf("socialite: invalid id token: %w", err)
	} else if expired {
		return nil, fmt.Errorf("socialite: invalid id token: token expired")
	}

	return raw, nil
}

func idTokenExpired(exp any) (bool, error) {
	const leeway = 2 * time.Minute
	now := time.Now().Add(-leeway).Unix()

	switch v := exp.(type) {
	case nil:
		return false, fmt.Errorf("missing exp claim")
	case float64:
		return now >= int64(v), nil
	case string:
		n, err := strconv.ParseInt(v, 10, 64)

		if err != nil {
			return false, fmt.Errorf("invalid exp claim: %w", err)
		}

		return now >= n, nil
	default:
		return false, fmt.Errorf("invalid exp claim type")
	}
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
