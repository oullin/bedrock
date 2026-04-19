package passport

import (
	cauth "github.com/bedrock/packages/contracts/auth"
)

// HasApiTokens is the interface that user structs implement to participate in
// Passport token authentication.
//
// In Laravel Passport this is provided via the HasApiTokens trait. In Go,
// consumers embed UserWithTokens or implement the interface directly.
type HasApiTokens interface {
	cauth.Authenticatable
	WithAccessToken(token *AccessToken) *UserWithTokens
	CurrentAccessToken() *AccessToken
	TokenCan(scope string) bool
	TokenCant(scope string) bool
}

// UserWithTokens wraps any Authenticatable with an attached AccessToken.
// It satisfies both cauth.Authenticatable and HasApiTokens.
//
// Use NewUserWithTokens to create a decorated user after resolving a bearer
// token in the TokenGuard.
type UserWithTokens struct {
	cauth.Authenticatable
	token *AccessToken
}

// NewUserWithTokens wraps user with the given access token.
func NewUserWithTokens(user cauth.Authenticatable, token *AccessToken) *UserWithTokens {
	return &UserWithTokens{
		Authenticatable: user,
		token:           token,
	}
}

// WithAccessToken returns a new UserWithTokens with the given token attached.
// Mirrors Laravel's withAccessToken($token) which returns $this.
func (u *UserWithTokens) WithAccessToken(token *AccessToken) *UserWithTokens {
	return &UserWithTokens{
		Authenticatable: u.Authenticatable,
		token:           token,
	}
}

// CurrentAccessToken returns the access token currently attached to this user.
func (u *UserWithTokens) CurrentAccessToken() *AccessToken {
	return u.token
}

// TokenCan reports whether the user's current access token has the given scope.
func (u *UserWithTokens) TokenCan(scope string) bool {
	return u.token != nil && u.token.Can(scope)
}

// TokenCant is the inverse of TokenCan.
func (u *UserWithTokens) TokenCant(scope string) bool {
	return !u.TokenCan(scope)
}
