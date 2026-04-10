package events

import (
	"net/http"

	contracts "github.com/bedrock/packages/contracts/auth"
)

// Attempting is dispatched when an authentication attempt begins.
type Attempting struct {
	Guard       string
	Credentials map[string]any
	Remember    bool
}

// Authenticated is dispatched when a user is resolved from the session or token.
type Authenticated struct {
	Guard string
	User  contracts.Authenticatable
}

// CurrentDeviceLogout is dispatched when the current device is logged out.
type CurrentDeviceLogout struct {
	Guard string
	User  contracts.Authenticatable
}

// Failed is dispatched when an authentication attempt fails.
type Failed struct {
	Guard       string
	User        contracts.Authenticatable // May be nil if user was not found.
	Credentials map[string]any
}

// Lockout is dispatched when a login attempt is throttled.
type Lockout struct {
	Request *http.Request
}

// Login is dispatched after a successful login.
type Login struct {
	Guard    string
	User     contracts.Authenticatable
	Remember bool
}

// Logout is dispatched after the user logs out.
type Logout struct {
	Guard string
	User  contracts.Authenticatable
}

// OtherDeviceLogout is dispatched when other devices are logged out.
type OtherDeviceLogout struct {
	Guard string
	User  contracts.Authenticatable
}

// PasswordReset is dispatched after a password is successfully reset.
type PasswordReset struct {
	User contracts.Authenticatable
}

// PasswordResetLinkSent is dispatched after a password reset link is sent.
type PasswordResetLinkSent struct {
	User contracts.CanResetPassword
}

// Registered is dispatched after a new user registers.
type Registered struct {
	User contracts.Authenticatable
}

// Validated is dispatched when credentials are validated but before login.
type Validated struct {
	Guard string
	User  contracts.Authenticatable
}

// Verified is dispatched when a user's email is verified.
type Verified struct {
	User contracts.MustVerifyEmail
}
