package socialite

import (
	"context"
	"fmt"
	"net/http"
)

// GitlabProvider handles OAuth2 authentication via GitLab (cloud or self-hosted).
// It mirrors Laravel\Socialite\Two\GitlabProvider.
type GitlabProvider struct {
	AbstractProvider
	host string
}

// NewGitlabProvider constructs a GitlabProvider targeting gitlab.com.
func NewGitlabProvider(req *http.Request, session Session, clientID, clientSecret, redirectURL string) *GitlabProvider {
	g := &GitlabProvider{host: "https://gitlab.com"}
	g.AbstractProvider = NewAbstractProvider(g, req, session, clientID, clientSecret, redirectURL)
	g.scopes = []string{"read_user"}
	g.scopeSep = " "
	return g
}

// SetHost overrides the GitLab instance URL (e.g. for self-hosted instances).
func (g *GitlabProvider) SetHost(host string) *GitlabProvider {
	g.host = host
	return g
}

func (g *GitlabProvider) GetAuthURL(state *string) string {
	return g.BuildAuthURLFromBase(fmt.Sprintf("%s/oauth/authorize", g.host), state)
}

func (g *GitlabProvider) GetTokenURL() string {
	return fmt.Sprintf("%s/oauth/token", g.host)
}

func (g *GitlabProvider) GetUserByToken(ctx context.Context, token string) (map[string]any, error) {
	return g.getWithBearer(ctx, fmt.Sprintf("%s/api/v4/user", g.host), token)
}

func (g *GitlabProvider) MapUserToObject(raw map[string]any) *User {
	u := &User{}
	u.ID = stringify(raw["id"])
	u.Nickname = stringify(raw["username"])
	u.Name = stringify(raw["name"])
	u.Email = stringify(raw["email"])
	u.Avatar = stringify(raw["avatar_url"])
	return u
}

func (g *GitlabProvider) getWithBearer(ctx context.Context, targetURL, token string) (map[string]any, error) {
	return g.AbstractProvider.getWithBearer(ctx, targetURL, token)
}
