package socialauth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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
	if IsJWT(token) {
		return DecodeJWTClaims(token)
	}

	return g.getWithBearer(ctx, "https://www.googleapis.com/oauth2/v3/userinfo", token)
}

// UserFromIDToken maps an OpenID Connect ID token without calling userinfo.
func (g *GoogleProvider) UserFromIDToken(_ context.Context, idToken string) (*User, error) {
	raw, err := DecodeJWTClaims(idToken)

	if err != nil {
		return nil, err
	}

	return g.MapUserToObject(raw).SetRaw(raw).SetToken(idToken), nil
}

// UserFromToken accepts either an access token or an OpenID Connect ID token.
func (g *GoogleProvider) UserFromToken(ctx context.Context, token string) (*User, error) {
	if IsJWT(token) {
		return g.UserFromIDToken(ctx, token)
	}

	return g.AbstractProvider.UserFromToken(ctx, token)
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

// IsJWT reports whether token has the three base64url segments of a JWT.
func IsJWT(token string) bool {
	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return false
	}

	for _, part := range parts {
		if part == "" {
			return false
		}
	}

	return true
}

// DecodeJWTClaims decodes JWT claims without performing provider signature
// verification. Production verification belongs at the provider boundary; this
// helper only mirrors SocialAuth's testable mapping path for already trusted
// OpenID Connect ID tokens.
func DecodeJWTClaims(token string) (map[string]any, error) {
	if !IsJWT(token) {
		return nil, errors.New("socialauth: token is not a jwt")
	}

	parts := strings.Split(token, ".")
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])

	if err != nil {
		return nil, err
	}

	claims := map[string]any{}

	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	return claims, nil
}
