package authflows

// Event names dispatched by AuthFlows.

// LoginAttemptedPayload is dispatched when a login attempt begins.
type LoginAttemptedPayload struct {
	Identifier string
	Remember   bool
}

// LoginSucceededPayload is dispatched after a successful login.
type LoginSucceededPayload struct {
	User     Authenticatable
	Remember bool
}

// LoginFailedPayload is dispatched after a failed login.
type LoginFailedPayload struct {
	Identifier string
}

// LoggedOutPayload is dispatched after logout.
type LoggedOutPayload struct {
	User Authenticatable
}

// RegisteredPayload is dispatched after a new user registers.
type RegisteredPayload struct {
	User Authenticatable
}

const (
	EventLoginAttempted  = "authflows.login.attempted"
	EventLoginSucceeded  = "authflows.login.succeeded"
	EventLoginFailed     = "authflows.login.failed"
	EventLoggedOut       = "authflows.logged_out"
	EventRegistered      = "authflows.registered"
	EventPasswordReset   = "authflows.password.reset"
	EventPasswordUpdated = "authflows.password.updated"
	EventVerified        = "authflows.email.verified"

	EventTwoFactorEnabled   = "authflows.two_factor.enabled"
	EventTwoFactorConfirmed = "authflows.two_factor.confirmed"
	EventTwoFactorDisabled  = "authflows.two_factor.disabled"
	EventTwoFactorChallenge = "authflows.two_factor.challenge"
	EventRecoveryCodeUsed   = "authflows.two_factor.recovery_used"
)
