package auth

// GenericUser is a simple in-memory implementation of Authenticatable, useful
// for testing or when users are loaded from arbitrary maps.
type GenericUser struct {
	Attributes map[string]any
}

// NewGenericUser creates a GenericUser from a map of attributes.
func NewGenericUser(attrs map[string]any) *GenericUser {
	return &GenericUser{Attributes: attrs}
}

func (u *GenericUser) GetAuthIdentifierName() string { return "id" }

func (u *GenericUser) GetAuthIdentifier() string {
	if v, ok := u.Attributes["id"].(string); ok {
		return v
	}

	return ""
}

// GetAuthIdentifierForBroadcasting returns the identifier used for private
// broadcast channel authorization.
func (u *GenericUser) GetAuthIdentifierForBroadcasting() string {
	return u.GetAuthIdentifier()
}

func (u *GenericUser) GetAuthPasswordName() string { return "password" }

func (u *GenericUser) GetAuthPassword() string {
	if v, ok := u.Attributes["password"].(string); ok {
		return v
	}

	return ""
}

func (u *GenericUser) SetAuthPassword(password string) {
	u.Attributes["password"] = password
}

func (u *GenericUser) GetRememberToken() string {
	if v, ok := u.Attributes["remember_token"].(string); ok {
		return v
	}

	return ""
}

func (u *GenericUser) SetRememberToken(token string) {
	u.Attributes["remember_token"] = token
}

func (u *GenericUser) GetRememberTokenName() string { return "remember_token" }

// Get returns the value for a given attribute key.
func (u *GenericUser) Get(key string) (any, bool) {
	v, ok := u.Attributes[key]

	return v, ok
}

// Set sets a value for a given attribute key.
func (u *GenericUser) Set(key string, value any) {
	u.Attributes[key] = value
}

// Has reports whether the attribute key exists.
func (u *GenericUser) Has(key string) bool {
	_, ok := u.Attributes[key]

	return ok
}

// Delete removes an attribute by key.
func (u *GenericUser) Delete(key string) {
	delete(u.Attributes, key)
}
