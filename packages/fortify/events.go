package fortify

import cauth "github.com/bedrock/packages/contracts/auth"

// Event names dispatched by Fortify.

// LoginAttemptedPayload is dispatched when a login attempt begins.
type LoginAttemptedPayload struct {
	Identifier string
	Remember   bool
}

// LoginSucceededPayload is dispatched after a successful login.
type LoginSucceededPayload struct {
	User     cauth.Authenticatable
	Remember bool
}

// LoginFailedPayload is dispatched after a failed login.
type LoginFailedPayload struct {
	Identifier string
}

// LoggedOutPayload is dispatched after logout.
type LoggedOutPayload struct {
	User cauth.Authenticatable
}

// RegisteredPayload is dispatched after a new user registers.
type RegisteredPayload struct {
	User cauth.Authenticatable
}

const (
	EventLoginAttempted  = "fortify.login.attempted"
	EventLoginSucceeded  = "fortify.login.succeeded"
	EventLoginFailed     = "fortify.login.failed"
	EventLoggedOut       = "fortify.logged_out"
	EventRegistered      = "fortify.registered"
	EventPasswordReset   = "fortify.password.reset"
	EventPasswordUpdated = "fortify.password.updated"
	EventVerified        = "fortify.email.verified"

	EventTwoFactorEnabled   = "fortify.two_factor.enabled"
	EventTwoFactorConfirmed = "fortify.two_factor.confirmed"
	EventTwoFactorDisabled  = "fortify.two_factor.disabled"
	EventTwoFactorChallenge = "fortify.two_factor.challenge"
	EventRecoveryCodeUsed   = "fortify.two_factor.recovery_used"
)
