package socialite

import (
	"context"
	"encoding/json"
	"net/http"
)

// GithubProvider handles OAuth2 authentication via GitHub.
// It mirrors upstream Socialite\Two\GithubProvider.
type GithubProvider struct {
	AbstractProvider
}

// NewGithubProvider constructs a GithubProvider.
func NewGithubProvider(req *http.Request, session Session, clientID, clientSecret, redirectURL string) *GithubProvider {
	g := &GithubProvider{}
	g.AbstractProvider = NewAbstractProvider(g, req, session, clientID, clientSecret, redirectURL)
	g.scopes = []string{"user:email"}
	g.scopeSep = " "

	return g
}

func (g *GithubProvider) GetAuthURL(state *string) string {
	return g.BuildAuthURLFromBase("https://github.com/login/oauth/authorize", state)
}

func (g *GithubProvider) GetTokenURL() string {
	return "https://github.com/login/oauth/access_token"
}

func (g *GithubProvider) GetUserByToken(ctx context.Context, token string) (map[string]any, error) {
	user, err := g.getJSON(ctx, "https://api.github.com/user", map[string]string{
		"Authorization": "token " + token,
	})

	if err != nil {
		return nil, err
	}

	// GitHub users may have no public email; fetch from the emails endpoint.
	if emailVal, _ := user["email"].(string); emailVal == "" {
		email, err := g.fetchPrimaryEmail(ctx, token)

		if err == nil && email != "" {
			user["email"] = email
		}
	}

	return user, nil
}

func (g *GithubProvider) fetchPrimaryEmail(ctx context.Context, token string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)

	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "token "+token)

	resp, err := g.getHTTPClient().Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var emails []map[string]any

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, m := range emails {
		if m["primary"] == true && m["verified"] == true {
			if email, ok := m["email"].(string); ok {
				return email, nil
			}
		}
	}

	return "", nil
}

func (g *GithubProvider) MapUserToObject(raw map[string]any) *User {
	u := &User{}
	u.ID = stringify(raw["id"])
	u.Nickname = stringify(raw["login"])
	u.Name = stringify(raw["name"])
	u.Email = stringify(raw["email"])
	u.Avatar = stringify(raw["avatar_url"])

	if v, ok := raw["node_id"]; ok {
		if u.Attributes == nil {
			u.Attributes = make(map[string]any)
		}

		u.Attributes["node_id"] = v
	}

	return u
}
