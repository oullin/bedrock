package socialite

import (
	"fmt"
	"strconv"
)

// User holds all fields returned by an OAuth provider after a successful
// authentication cycle. OAuth2-specific fields (Token, RefreshToken, etc.)
// are populated for OAuth2 providers; OAuth1 providers set Token and
// TokenSecret instead.
type User struct {
	// Core identity fields (all providers).
	ID       string
	Nickname string
	Name     string
	Email    string
	Avatar   string

	// Raw is the unmodified JSON payload returned by the provider's user
	// info endpoint.
	Raw map[string]any

	// Attributes holds additional provider-specific fields mapped via Map().
	Attributes map[string]any

	// OAuth2 fields.
	Token          string
	RefreshToken   string
	ExpiresIn      int
	ApprovedScopes []string

	// OAuth1 fields.
	TokenSecret string
}

// Map applies a flat attribute map to the user's core fields, returning the
// receiver for chaining. Unknown keys are stored in Attributes.
func (u *User) Map(attrs map[string]any) *User {
	if v, ok := attrs["id"]; ok {
		u.ID = stringify(v)
	}

	if v, ok := attrs["nickname"]; ok {
		u.Nickname = stringify(v)
	}

	if v, ok := attrs["name"]; ok {
		u.Name = stringify(v)
	}

	if v, ok := attrs["email"]; ok {
		u.Email = stringify(v)
	}

	if v, ok := attrs["avatar"]; ok {
		u.Avatar = stringify(v)
	}

	if v, ok := attrs["avatar_original"]; ok {
		if u.Attributes == nil {
			u.Attributes = make(map[string]any)
		}

		u.Attributes["avatar_original"] = v
	}

	return u
}

// SetRaw stores the raw provider response.
func (u *User) SetRaw(raw map[string]any) *User {
	u.Raw = raw

	return u
}

// GetRaw returns the raw provider response.
func (u *User) GetRaw() map[string]any { return u.Raw }

// Get returns a field value by key, checking core fields first, then
// Attributes, then Raw.
func (u *User) Get(key string) any {
	switch key {
	case "id":
		return u.ID
	case "nickname":
		return u.Nickname
	case "name":
		return u.Name
	case "email":
		return u.Email
	case "avatar":
		return u.Avatar
	case "token":
		return u.Token
	case "refreshToken":
		return u.RefreshToken
	case "expiresIn":
		return u.ExpiresIn
	case "approvedScopes":
		return u.ApprovedScopes
	case "tokenSecret":
		return u.TokenSecret
	}

	if u.Attributes != nil {
		if v, ok := u.Attributes[key]; ok {
			return v
		}
	}

	if u.Raw != nil {
		return u.Raw[key]
	}

	return nil
}

// contracts/socialite.User interface implementation.

func (u *User) GetID() string       { return u.ID }
func (u *User) GetNickname() string { return u.Nickname }
func (u *User) GetName() string     { return u.Name }
func (u *User) GetEmail() string    { return u.Email }
func (u *User) GetAvatar() string   { return u.Avatar }

// SetToken sets OAuth2 access token.
func (u *User) SetToken(token string) *User { u.Token = token; return u }

// SetRefreshToken sets the OAuth2 refresh token.
func (u *User) SetRefreshToken(token string) *User { u.RefreshToken = token; return u }

// SetExpiresIn sets seconds until the access token expires.
func (u *User) SetExpiresIn(secs int) *User { u.ExpiresIn = secs; return u }

// SetApprovedScopes sets the list of scopes approved by the user.
func (u *User) SetApprovedScopes(scopes []string) *User { u.ApprovedScopes = scopes; return u }

// SetOAuthToken sets OAuth1 token and secret.
func (u *User) SetOAuthToken(token, secret string) *User {
	u.Token = token
	u.TokenSecret = secret

	return u
}

// stringify converts any value to its string representation.
func stringify(v any) string {
	if v == nil {
		return ""
	}

	switch s := v.(type) {
	case string:
		return s
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", s)
	}
}
