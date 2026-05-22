package socialite

import (
	"context"
	"net/http"
)

// BitbucketProvider handles OAuth2 authentication via Bitbucket.
type BitbucketProvider struct {
	AbstractProvider
}

// NewBitbucketProvider constructs a BitbucketProvider.
func NewBitbucketProvider(req *http.Request, session Session, clientID, clientSecret, redirectURL string) *BitbucketProvider {
	b := &BitbucketProvider{}
	b.AbstractProvider = NewAbstractProvider(b, req, session, clientID, clientSecret, redirectURL)
	b.scopes = []string{"email"}
	b.scopeSep = " "

	return b
}

func (b *BitbucketProvider) GetAuthURL(state *string) string {
	return b.BuildAuthURLFromBase("https://bitbucket.org/site/oauth2/authorize", state)
}

func (b *BitbucketProvider) GetTokenURL() string {
	return "https://bitbucket.org/site/oauth2/access_token"
}

func (b *BitbucketProvider) GetUserByToken(ctx context.Context, token string) (map[string]any, error) {
	user, err := b.getWithBearer(ctx, "https://api.bitbucket.org/2.0/user", token)

	if err != nil {
		return nil, err
	}

	// Fetch primary confirmed email.
	emails, err := b.getWithBearer(ctx, "https://api.bitbucket.org/2.0/user/emails", token)

	if err == nil {
		if values, ok := emails["values"].([]any); ok {
			for _, v := range values {
				m, ok := v.(map[string]any)

				if !ok {
					continue
				}

				if m["is_primary"] == true && m["is_confirmed"] == true {
					user["email"] = m["email"]

					break
				}
			}
		}
	}

	return user, nil
}

func (b *BitbucketProvider) MapUserToObject(raw map[string]any) *User {
	u := &User{}
	u.ID = stringify(raw["account_id"])
	u.Nickname = stringify(raw["username"])
	u.Name = stringify(raw["display_name"])
	u.Email = stringify(raw["email"])

	if links, ok := raw["links"].(map[string]any); ok {
		if avatar, ok := links["avatar"].(map[string]any); ok {
			u.Avatar = stringify(avatar["href"])
		}
	}

	return u
}

func (b *BitbucketProvider) getWithBearer(ctx context.Context, targetURL, token string) (map[string]any, error) {
	return b.AbstractProvider.getWithBearer(ctx, targetURL, token)
}
