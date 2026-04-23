package socialauth

import (
	"context"
	"net/http"
	"strings"
)

// LinkedInProvider handles OAuth2 authentication via LinkedIn.
// It mirrors Upstream\SocialAuth\Two\LinkedInProvider.
type LinkedInProvider struct {
	AbstractProvider
}

// NewLinkedInProvider constructs a LinkedInProvider.
func NewLinkedInProvider(req *http.Request, session Session, clientID, clientSecret, redirectURL string) *LinkedInProvider {
	l := &LinkedInProvider{}
	l.AbstractProvider = NewAbstractProvider(l, req, session, clientID, clientSecret, redirectURL)
	l.scopes = []string{"r_liteprofile", "r_emailaddress"}
	l.scopeSep = " "

	return l
}

func (l *LinkedInProvider) GetAuthURL(state *string) string {
	return l.BuildAuthURLFromBase("https://www.linkedin.com/oauth/v2/authorization", state)
}

func (l *LinkedInProvider) GetTokenURL() string {
	return "https://www.linkedin.com/oauth/v2/accessToken"
}

func (l *LinkedInProvider) GetUserByToken(ctx context.Context, token string) (map[string]any, error) {
	profile, err := l.getJSON(ctx,
		"https://api.linkedin.com/v2/me?projection=(id,localizedFirstName,localizedLastName,profilePicture(displayImage~:playableStreams))",
		map[string]string{
			"Authorization":             "Bearer " + token,
			"X-RestLi-Protocol-Version": "2.0.0",
		},
	)

	if err != nil {
		return nil, err
	}

	emailData, err := l.getJSON(ctx,
		"https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))",
		map[string]string{
			"Authorization":             "Bearer " + token,
			"X-RestLi-Protocol-Version": "2.0.0",
		},
	)

	if err == nil {
		if elements, ok := emailData["elements"].([]any); ok && len(elements) > 0 {
			if el, ok := elements[0].(map[string]any); ok {
				if handle, ok := el["handle~"].(map[string]any); ok {
					profile["emailAddress"] = handle["emailAddress"]
				}
			}
		}
	}

	return profile, nil
}

func (l *LinkedInProvider) MapUserToObject(raw map[string]any) *User {
	u := &User{}

	u.ID = stringify(raw["id"])

	if u.ID == "" {
		u.ID = stringify(raw["sub"])
	}

	u.Name = strings.TrimSpace(stringify(raw["localizedFirstName"]) + " " + stringify(raw["localizedLastName"]))

	if u.Name == "" {
		u.Name = stringify(raw["name"])
	}

	u.Email = stringify(raw["emailAddress"])

	if u.Email == "" {
		u.Email = stringify(raw["email"])
	}

	u.Avatar = l.extractAvatar(raw)

	if u.Avatar == "" {
		u.Avatar = stringify(raw["picture"])
	}

	return u
}

func (l *LinkedInProvider) extractAvatar(raw map[string]any) string {
	pp, ok := raw["profilePicture"].(map[string]any)

	if !ok {
		return ""
	}

	di, ok := pp["displayImage~"].(map[string]any)

	if !ok {
		return ""
	}

	elements, ok := di["elements"].([]any)

	if !ok || len(elements) == 0 {
		return ""
	}

	for _, el := range elements {
		m, ok := el.(map[string]any)

		if !ok {
			continue
		}

		identifiers, ok := m["identifiers"].([]any)

		if !ok || len(identifiers) == 0 {
			continue
		}

		if id, ok := identifiers[0].(map[string]any); ok {
			return stringify(id["identifier"])
		}
	}

	return ""
}

func (l *LinkedInProvider) getJSON(ctx context.Context, targetURL string, headers map[string]string) (map[string]any, error) {
	return l.AbstractProvider.getJSON(ctx, targetURL, headers)
}
