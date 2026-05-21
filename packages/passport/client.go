package passport

// Client represents an OAuth2 client application, corresponding to the
// oauth_clients table in Passport.
type Client struct {
	ID                   string
	UserID               string
	Name                 string
	Secret               string
	Provider             string
	RedirectURIs         []string
	GrantTypes           []string
	Scopes               []string
	PersonalAccessClient bool
	PasswordClient       bool
	Revoked              bool
}

// FirstParty reports whether this is a first-party client.
// A client is first-party if it is a personal access client or a password client.
func (c *Client) FirstParty() bool {
	return c.PersonalAccessClient || c.PasswordClient
}

// Confidential reports whether this client has a non-empty secret.
// Confidential clients (web server apps) can keep a secret; public clients cannot.
func (c *Client) Confidential() bool {
	return c.Secret != ""
}

// HasGrantType reports whether the given grant type is registered for this client.
func (c *Client) HasGrantType(grantType string) bool {
	for _, g := range c.GrantTypes {
		if g == grantType {
			return true
		}
	}

	return false
}

// HasScope reports whether the given scope is in the client's allowed scope list.
func (c *Client) HasScope(scope string) bool {
	for _, s := range c.Scopes {
		if s == scope {
			return true
		}
	}

	return false
}
