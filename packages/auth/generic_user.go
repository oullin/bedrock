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

func (u *GenericUser) GetAuthIdentifier() any { return u.Attributes["id"] }

func (u *GenericUser) GetAuthPassword() string {
	if v, ok := u.Attributes["password"].(string); ok {
		return v
	}

	return ""
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
