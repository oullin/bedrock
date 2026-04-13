package inception

import "net/http"

// Router is a minimal interface for registering routes.
// The consuming app provides an adapter from its framework's router.
type Router interface {
	Post(pattern string, handler http.Handler)
	Put(pattern string, handler http.Handler)
	Get(pattern string, handler http.Handler)
	Delete(pattern string, handler http.Handler)
}

// StdMuxRouter adapts net/http.ServeMux to the Router interface.
type StdMuxRouter struct {
	Mux *http.ServeMux
}

// Post registers a POST route.

// Put registers a PUT route.

// Get registers a GET route.

// Delete registers a DELETE route.

// RouteConfig holds optional dependencies for route registration.
// Only the dependencies relevant to enabled features need to be set.
type RouteConfig struct {
	Tokens       TokenRepository
	Sessions     SessionRepository
	Photos       UpdatesProfilePhotos
	DeletePhotos DeletesProfilePhotos
	DeleteUser   DeletesUsers
}

func (r *StdMuxRouter) Post(pattern string, handler http.Handler) {
	r.Mux.Handle("POST "+pattern, handler)
}

func (r *StdMuxRouter) Put(pattern string, handler http.Handler) {
	r.Mux.Handle("PUT "+pattern, handler)
}

func (r *StdMuxRouter) Get(pattern string, handler http.Handler) {
	r.Mux.Handle("GET "+pattern, handler)
}

func (r *StdMuxRouter) Delete(pattern string, handler http.Handler) {
	r.Mux.Handle("DELETE "+pattern, handler)
}

// RegisterRoutes registers all enabled Inception routes on the given router.
// The issuer parameter is used for TOTP provisioning URIs (e.g. app name).
func RegisterRoutes(router Router, app *Inception, issuer string, config RouteConfig) {
	// Authentication
	router.Post("/login", NewLoginHandler(app))
	router.Post("/logout", NewLogoutHandler(app))

	// Registration
	if app.config.Features.Registration {
		router.Post("/register", NewRegisterHandler(app))
	}

	// Password Reset
	if app.config.Features.ResetPasswords {
		router.Post("/forgot-password", NewForgotPasswordHandler(app))
		router.Post("/reset-password", NewResetPasswordHandler(app))
	}

	// Email Verification
	if app.config.Features.EmailVerification {
		router.Post("/email/verification-notification", NewSendVerificationHandler(app))
		router.Get("/verify-email/{id}/{hash}", NewVerifyEmailHandler(app))
	}

	// Two-Factor Authentication
	if app.config.Features.TwoFactorAuthentication {
		router.Post("/user/two-factor-authentication", NewEnableTwoFactorHandler(app))
		router.Post("/user/confirmed-two-factor-authentication", NewConfirmTwoFactorHandler(app))
		router.Delete("/user/two-factor-authentication", NewDisableTwoFactorHandler(app))
		router.Get("/user/two-factor-qr-code", NewTwoFactorQRCodeHandler(app, issuer))
		router.Get("/user/two-factor-recovery-codes", NewTwoFactorRecoveryCodesHandler(app))
		router.Post("/user/two-factor-recovery-codes", NewTwoFactorRecoveryCodesHandler(app))
		router.Post("/two-factor-challenge", NewTwoFactorChallengeHandler(app))
	}

	// Password Update
	if app.config.Features.UpdatePasswords {
		router.Put("/user/password", NewUpdatePasswordHandler(app))
	}

	// Password Confirmation
	if app.config.Features.ConfirmPassword {
		router.Post("/user/confirm-password", NewConfirmPasswordHandler(app))
	}

	// Profile Information
	if app.config.Features.UpdateProfileInformation {
		router.Put("/user/profile-information", NewUpdateProfileHandler(app))
	}

	// Teams
	if app.config.Features.Teams {
		router.Post("/teams", NewCreateTeamHandler(app))
		router.Put("/teams/{team}", NewUpdateTeamHandler(app))
		router.Delete("/teams/{team}", NewDeleteTeamHandler(app))

		router.Post("/teams/{team}/members", NewAddTeamMemberHandler(app))
		router.Put("/teams/{team}/members/{user}", NewUpdateTeamMemberRoleHandler(app))
		router.Delete("/teams/{team}/members/{user}", NewRemoveTeamMemberHandler(app))

		router.Post("/teams/{team}/invitations", NewInviteTeamMemberHandler(app))
		router.Delete("/team-invitations/{invitation}", NewCancelInvitationHandler(app))
		router.Get("/team-invitations/{invitation}/accept", NewAcceptInvitationHandler(app))

		router.Put("/current-team", NewSwitchTeamHandler(app))
	}

	// API Tokens
	if app.config.Features.APITokens && config.Tokens != nil {
		router.Post("/user/api-tokens", NewCreateTokenHandler(app, config.Tokens))
		router.Put("/user/api-tokens/{token}", NewUpdateTokenHandler(app, config.Tokens))
		router.Delete("/user/api-tokens/{token}", NewDeleteTokenHandler(app, config.Tokens))
	}

	// Profile Photos
	if app.config.Features.ProfilePhotos && config.Photos != nil {
		router.Put("/user/profile-photo", NewUpdateProfilePhotoHandler(app, config.Photos))
	}

	if app.config.Features.ProfilePhotos && config.DeletePhotos != nil {
		router.Delete("/user/profile-photo", NewDeleteProfilePhotoHandler(app, config.DeletePhotos))
	}

	// Account Deletion
	if app.config.Features.AccountDeletion && config.DeleteUser != nil {
		router.Delete("/user", NewDeleteAccountHandler(app, config.DeleteUser))
	}

	// Browser Sessions
	if app.config.Features.BrowserSessions && config.Sessions != nil {
		router.Get("/user/sessions", NewListSessionsHandler(app, config.Sessions))
		router.Delete("/user/other-sessions", NewDeleteOtherSessionsHandler(app, config.Sessions))
	}
}
