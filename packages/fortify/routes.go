package fortify

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
func (r *StdMuxRouter) Post(pattern string, handler http.Handler) {
	r.Mux.Handle("POST "+pattern, handler)
}

// Put registers a PUT route.
func (r *StdMuxRouter) Put(pattern string, handler http.Handler) {
	r.Mux.Handle("PUT "+pattern, handler)
}

// Get registers a GET route.
func (r *StdMuxRouter) Get(pattern string, handler http.Handler) {
	r.Mux.Handle("GET "+pattern, handler)
}

// Delete registers a DELETE route.
func (r *StdMuxRouter) Delete(pattern string, handler http.Handler) {
	r.Mux.Handle("DELETE "+pattern, handler)
}

// RegisterRoutes registers all enabled Fortify routes on the given router.
// The issuer parameter is used for TOTP provisioning URIs (e.g. app name).
func RegisterRoutes(router Router, f *Fortify, issuer string) {
	// Authentication
	router.Post("/login", NewLoginHandler(f))
	router.Post("/logout", NewLogoutHandler(f))

	// Registration
	if f.config.Features.Registration {
		router.Post("/register", NewRegisterHandler(f))
	}

	// Password Reset
	if f.config.Features.ResetPasswords {
		router.Post("/forgot-password", NewForgotPasswordHandler(f))
		router.Post("/reset-password", NewResetPasswordHandler(f))
	}

	// Email Verification
	if f.config.Features.EmailVerification {
		router.Post("/email/verification-notification", NewSendVerificationHandler(f))
		router.Get("/verify-email/{id}/{hash}", NewVerifyEmailHandler(f))
	}

	// Two-Factor Authentication
	if f.config.Features.TwoFactorAuthentication {
		router.Post("/user/two-factor-authentication", NewEnableTwoFactorHandler(f))
		router.Post("/user/confirmed-two-factor-authentication", NewConfirmTwoFactorHandler(f))
		router.Delete("/user/two-factor-authentication", NewDisableTwoFactorHandler(f))
		router.Get("/user/two-factor-qr-code", NewTwoFactorQRCodeHandler(f, issuer))
		router.Get("/user/two-factor-recovery-codes", NewTwoFactorRecoveryCodesHandler(f))
		router.Post("/user/two-factor-recovery-codes", NewTwoFactorRecoveryCodesHandler(f))
		router.Post("/two-factor-challenge", NewTwoFactorChallengeHandler(f))
	}

	// Password Update
	if f.config.Features.UpdatePasswords {
		router.Put("/user/password", NewUpdatePasswordHandler(f))
	}

	// Password Confirmation
	if f.config.Features.ConfirmPassword {
		router.Post("/user/confirm-password", NewConfirmPasswordHandler(f))
	}

	// Profile Information
	if f.config.Features.UpdateProfileInformation {
		router.Put("/user/profile-information", NewUpdateProfileHandler(f))
	}
}
