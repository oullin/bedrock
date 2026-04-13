package inception

import cauth "github.com/bedrock/packages/contracts/auth"

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

// Event names dispatched by Inception.
const (
	EventLoginAttempted  = "inception.login.attempted"
	EventLoginSucceeded  = "inception.login.succeeded"
	EventLoginFailed     = "inception.login.failed"
	EventLoggedOut       = "inception.logged_out"
	EventRegistered      = "inception.registered"
	EventPasswordReset   = "inception.password.reset"
	EventPasswordUpdated = "inception.password.updated"
	EventVerified        = "inception.email.verified"
	EventProfileUpdated  = "inception.profile.updated"

	EventTwoFactorEnabled   = "inception.two_factor.enabled"
	EventTwoFactorConfirmed = "inception.two_factor.confirmed"
	EventTwoFactorDisabled  = "inception.two_factor.disabled"
	EventTwoFactorChallenge = "inception.two_factor.challenge"
	EventRecoveryCodeUsed   = "inception.two_factor.recovery_used"

	EventTeamCreated       = "inception.team.created"
	EventTeamUpdated       = "inception.team.updated"
	EventTeamDeleted       = "inception.team.deleted"
	EventTeamMemberAdded   = "inception.team.member_added"
	EventTeamMemberUpdated = "inception.team.member_updated"
	EventTeamMemberRemoved = "inception.team.member_removed"
	EventTeamMemberInvited = "inception.team.member_invited"
	EventTeamSwitched      = "inception.team.switched"
	EventTokenCreated      = "inception.token.created"
	EventTokenUpdated      = "inception.token.updated"
	EventTokenDeleted      = "inception.token.deleted"
	EventUserDeleted       = "inception.user.deleted"
)
