package auth

import "strings"

// GenericUser is a simple map-backed authenticatable.
type GenericUser struct {
	attributes map[string]string
}

// NewGenericUser returns a new generic user.
func NewGenericUser(attributes map[string]string) *GenericUser {
	cloned := make(map[string]string, len(attributes))

	for key, value := range attributes {
		cloned[key] = value
	}

	return &GenericUser{attributes: cloned}
}

// Attributes returns a clone of the stored attributes.
func (u *GenericUser) Attributes() map[string]string {
	cloned := make(map[string]string, len(u.attributes))

	for key, value := range u.attributes {
		cloned[key] = value
	}

	return cloned
}

func (u *GenericUser) GetAuthIdentifierName() string { return "id" }
func (u *GenericUser) GetAuthIdentifier() string     { return u.attributes["id"] }
func (u *GenericUser) GetAuthPasswordName() string   { return "password" }
func (u *GenericUser) GetAuthPassword() string       { return u.attributes["password"] }
func (u *GenericUser) SetAuthPassword(password string) {
	u.attributes["password"] = password
}

func (u *GenericUser) GetRememberToken() string { return u.attributes["remember_token"] }
func (u *GenericUser) SetRememberToken(token string) {
	u.attributes["remember_token"] = token
}

func (u *GenericUser) GetRememberTokenName() string { return "remember_token" }

func (u *GenericUser) GetName() string { return u.attributes["name"] }
func (u *GenericUser) SetName(name string) {
	u.attributes["name"] = strings.TrimSpace(name)
}

func (u *GenericUser) GetEmail() string { return u.attributes["email"] }
func (u *GenericUser) SetEmail(email string) {
	u.attributes["email"] = strings.TrimSpace(email)
}

var _ Authenticatable = (*GenericUser)(nil)
var _ UserProfile = (*GenericUser)(nil)
