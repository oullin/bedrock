package socialauth

import (
	"context"
	"net/http"
)

// SlackProvider handles OAuth2 authentication via Slack.
type SlackProvider struct {
	AbstractProvider
	asBotUser bool
}

// NewSlackProvider constructs a SlackProvider.
func NewSlackProvider(req *http.Request, session Session, clientID, clientSecret, redirectURL string) *SlackProvider {
	s := &SlackProvider{}
	s.AbstractProvider = NewAbstractProvider(s, req, session, clientID, clientSecret, redirectURL)
	s.scopes = []string{"identity.basic", "identity.email", "identity.team", "identity.avatar"}
	s.scopeSep = ","

	return s
}

// AsBotUser configures the provider to retrieve the bot-user token rather than
// the human user token.
func (s *SlackProvider) AsBotUser() *SlackProvider {
	s.asBotUser = true

	return s
}

func (s *SlackProvider) GetAuthURL(state *string) string {
	return s.BuildAuthURLFromBase("https://slack.com/oauth/v2/authorize", state)
}

func (s *SlackProvider) GetTokenURL() string {
	return "https://slack.com/api/oauth.v2.access"
}

func (s *SlackProvider) GetUserByToken(ctx context.Context, token string) (map[string]any, error) {
	return s.getJSON(ctx, "https://slack.com/api/users.identity", map[string]string{
		"Authorization": "Bearer " + token,
	})
}

func (s *SlackProvider) MapUserToObject(raw map[string]any) *User {
	u := &User{}

	if user, ok := raw["user"].(map[string]any); ok {
		u.ID = stringify(user["id"])
		u.Name = stringify(user["name"])
		u.Email = stringify(user["email"])

		if image, ok := user["image_512"].(string); ok {
			u.Avatar = image
		}
	} else {
		u.ID = stringify(raw["sub"])
		u.Name = stringify(raw["name"])
		u.Email = stringify(raw["email"])
		u.Avatar = stringify(raw["picture"])
	}

	if team, ok := raw["team"].(map[string]any); ok {
		if u.Attributes == nil {
			u.Attributes = make(map[string]any)
		}

		u.Attributes["organization"] = stringify(team["name"])
	}

	return u
}

func (s *SlackProvider) getJSON(ctx context.Context, targetURL string, headers map[string]string) (map[string]any, error) {
	return s.AbstractProvider.getJSON(ctx, targetURL, headers)
}
