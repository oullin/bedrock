package authflows

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/authflows/actions"
	"github.com/gollin/packages/auth/authflows/contracts"
	"github.com/gollin/packages/auth/authflows/responses"
	"github.com/gollin/packages/auth/foundation"
	authmw "github.com/gollin/packages/auth/middleware"
	"github.com/gollin/packages/auth/passwords"
	"github.com/gollin/packages/auth/support/ratelimit"
	configpkg "github.com/gollin/packages/config"
)

// Dependencies wires AuthFlows services together.
type Dependencies struct {
	AuthManager                         *auth.Manager
	Users                               auth.UserRepository
	PasswordBroker                      *passwords.Broker
	Verification                        *foundation.VerificationService
	CreateUsers                         contracts.CreatesNewUsers
	ResetUserPasswords                  contracts.ResetsUserPasswords
	UpdateUserPasswords                 contracts.UpdatesUserPasswords
	UpdateUserProfileInformation        contracts.UpdatesUserProfileInformation
	TwoFactorProvider                   contracts.TwoFactorAuthenticationProvider
	RedirectsIfTwoFactorAuthenticatable contracts.RedirectsIfTwoFactorAuthenticatable
	Responses                           *responses.Registry
	Limiter                             *ratelimit.Limiter
	Clock                               auth.Clock
}

// Server is the AuthFlows application surface.
type Server struct {
	config          Config
	features        Features
	paths           RoutePath
	authManager     *auth.Manager
	users           auth.UserRepository
	passwordBroker  *passwords.Broker
	verification    *foundation.VerificationService
	createUsers     contracts.CreatesNewUsers
	resetUsers      contracts.ResetsUserPasswords
	updatePasswords contracts.UpdatesUserPasswords
	updateProfiles  contracts.UpdatesUserProfileInformation
	twoFactor       contracts.TwoFactorAuthenticationProvider
	redirects       contracts.RedirectsIfTwoFactorAuthenticatable
	responses       *responses.Registry
	limiter         *ratelimit.Limiter
	clock           auth.Clock
}

// NewServerFromRepository loads AuthFlows config and creates a server.
func NewServerFromRepository(repo *configpkg.Repository, deps Dependencies) (*Server, error) {
	cfg, err := ConfigFromRepository(repo)
	if err != nil {
		return nil, err
	}
	return NewServer(cfg, deps)
}

// NewServer creates a AuthFlows server.
func NewServer(cfg Config, deps Dependencies) (*Server, error) {
	if deps.AuthManager == nil {
		return nil, fmt.Errorf("authflows: auth manager is required")
	}
	if deps.Users == nil {
		return nil, fmt.Errorf("authflows: user repository is required")
	}
	if deps.PasswordBroker == nil {
		return nil, fmt.Errorf("authflows: password broker is required")
	}
	if deps.Verification == nil {
		return nil, fmt.Errorf("authflows: verification service is required")
	}
	if deps.Clock == nil {
		deps.Clock = auth.SystemClock{}
	}
	if deps.Limiter == nil {
		deps.Limiter = ratelimit.New()
	}
	if deps.CreateUsers == nil {
		deps.CreateUsers = actions.CreateUser{
			Users:  deps.Users,
			Hasher: deps.AuthManager.Hasher(),
			IDs:    auth.RandomIDGenerator{},
			Clock:  deps.Clock,
		}
	}
	if deps.ResetUserPasswords == nil {
		deps.ResetUserPasswords = actions.ResetUserPassword{
			Users:  deps.Users,
			Hasher: deps.AuthManager.Hasher(),
			Clock:  deps.Clock,
		}
	}
	if deps.UpdateUserPasswords == nil {
		deps.UpdateUserPasswords = actions.UpdateUserPassword{
			Users:  deps.Users,
			Hasher: deps.AuthManager.Hasher(),
			Clock:  deps.Clock,
		}
	}
	if deps.UpdateUserProfileInformation == nil {
		deps.UpdateUserProfileInformation = actions.UpdateUserProfileInformation{
			Users: deps.Users,
		}
	}
	if deps.TwoFactorProvider == nil {
		deps.TwoFactorProvider = DefaultTwoFactorProvider{}
	}
	if deps.RedirectsIfTwoFactorAuthenticatable == nil {
		deps.RedirectsIfTwoFactorAuthenticatable = StaticTwoFactorRedirect{
			Path: cfg.Prefix + NewRoutePath(cfg.Paths).For("two-factor.login", "/two-factor-challenge"),
		}
	}
	if deps.Responses == nil {
		deps.Responses = responses.DefaultRegistry()
	}

	return &Server{
		config:          cfg,
		features:        NewFeatures(cfg),
		paths:           NewRoutePath(cfg.Paths),
		authManager:     deps.AuthManager,
		users:           deps.Users,
		passwordBroker:  deps.PasswordBroker,
		verification:    deps.Verification,
		createUsers:     deps.CreateUsers,
		resetUsers:      deps.ResetUserPasswords,
		updatePasswords: deps.UpdateUserPasswords,
		updateProfiles:  deps.UpdateUserProfileInformation,
		twoFactor:       deps.TwoFactorProvider,
		redirects:       deps.RedirectsIfTwoFactorAuthenticatable,
		responses:       deps.Responses,
		limiter:         deps.Limiter,
		clock:           deps.Clock,
	}, nil
}

// RegisterRoutes attaches AuthFlows routes to a mux.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	stack := authmw.Stack{
		Guard:  s.authManager.DefaultGuard(),
		Config: s.authManager.Config(),
		Clock:  s.clock,
	}

	if s.config.Views {
		mux.Handle("GET "+s.fullPath(s.paths.For("login", "/login")), s.guest(http.HandlerFunc(s.loginView)))
	}
	mux.Handle("POST "+s.fullPath(s.paths.For("login", "/login")), s.guest(http.HandlerFunc(s.login)))
	mux.Handle("POST "+s.fullPath(s.paths.For("logout", "/logout")), stack.RequireAuthenticated(http.HandlerFunc(s.logout)))

	if s.features.Enabled(FeatureResetPasswords) {
		if s.config.Views {
			mux.Handle("GET "+s.fullPath(s.paths.For("password.request", "/forgot-password")), s.guest(http.HandlerFunc(s.passwordResetRequestView)))
			mux.Handle("GET "+s.fullPath(s.paths.For("password.reset", "/reset-password/{token}")), s.guest(http.HandlerFunc(s.newPasswordView)))
		}
		mux.Handle("POST "+s.fullPath(s.paths.For("password.email", "/forgot-password")), s.guest(http.HandlerFunc(s.sendPasswordResetLink)))
		mux.Handle("POST "+s.fullPath(s.paths.For("password.update", "/reset-password")), s.guest(http.HandlerFunc(s.resetPassword)))
	}

	if s.features.Enabled(FeatureRegistration) {
		if s.config.Views {
			mux.Handle("GET "+s.fullPath(s.paths.For("register", "/register")), s.guest(http.HandlerFunc(s.registerView)))
		}
		mux.Handle("POST "+s.fullPath(s.paths.For("register", "/register")), s.guest(http.HandlerFunc(s.register)))
	}

	if s.features.Enabled(FeatureEmailVerification) {
		if s.config.Views {
			mux.Handle("GET "+s.fullPath(s.paths.For("verification.notice", "/email/verify")), stack.RequireAuthenticated(http.HandlerFunc(s.emailVerificationPrompt)))
		}
		mux.Handle("GET "+s.fullPath(s.paths.For("verification.verify", "/email/verify/{id}/{hash}")), stack.RequireAuthenticated(http.HandlerFunc(s.verifyEmail)))
		mux.Handle("POST "+s.fullPath(s.paths.For("verification.send", "/email/verification-notification")), stack.RequireAuthenticated(http.HandlerFunc(s.sendEmailVerification)))
	}

	if s.features.Enabled(FeatureUpdateProfileInformation) {
		mux.Handle("PUT "+s.fullPath(s.paths.For("user-profile-information.update", "/user/profile-information")), stack.RequireAuthenticated(http.HandlerFunc(s.updateProfileInformation)))
	}
	if s.features.Enabled(FeatureUpdatePasswords) {
		mux.Handle("PUT "+s.fullPath(s.paths.For("user-password.update", "/user/password")), stack.RequireAuthenticated(http.HandlerFunc(s.updatePassword)))
	}

	if s.config.Views {
		mux.Handle("GET "+s.fullPath(s.paths.For("password.confirm", "/user/confirm-password")), stack.RequireAuthenticated(http.HandlerFunc(s.confirmPasswordView)))
	}
	mux.Handle("GET "+s.fullPath(s.paths.For("password.confirmation", "/user/confirmed-password-status")), stack.RequireAuthenticated(http.HandlerFunc(s.confirmedPasswordStatus)))
	mux.Handle("POST "+s.fullPath(s.paths.For("password.confirm", "/user/confirm-password")), stack.RequireAuthenticated(http.HandlerFunc(s.confirmPassword)))

	if s.features.Enabled(FeatureTwoFactorAuthentication) {
		if s.config.Views {
			mux.Handle("GET "+s.fullPath(s.paths.For("two-factor.login", "/two-factor-challenge")), s.guest(http.HandlerFunc(s.twoFactorChallengeView)))
		}
		mux.Handle("POST "+s.fullPath(s.paths.For("two-factor.login", "/two-factor-challenge")), s.guest(http.HandlerFunc(s.completeTwoFactorLogin)))

		protected := stack.RequireAuthenticated(http.HandlerFunc(s.enableTwoFactor))
		if s.features.OptionEnabled(FeatureTwoFactorAuthentication, "confirmPassword") {
			protected = stack.RequirePasswordConfirmed(http.HandlerFunc(s.enableTwoFactor))
		}
		mux.Handle("POST "+s.fullPath(s.paths.For("two-factor.enable", "/user/two-factor-authentication")), protected)
		mux.Handle("POST "+s.fullPath(s.paths.For("two-factor.confirm", "/user/confirmed-two-factor-authentication")), stack.RequireAuthenticated(http.HandlerFunc(s.confirmTwoFactor)))
		mux.Handle("DELETE "+s.fullPath(s.paths.For("two-factor.disable", "/user/two-factor-authentication")), stack.RequireAuthenticated(http.HandlerFunc(s.disableTwoFactor)))
		mux.Handle("GET "+s.fullPath(s.paths.For("two-factor.qr-code", "/user/two-factor-qr-code")), stack.RequireAuthenticated(http.HandlerFunc(s.twoFactorQRCode)))
		mux.Handle("GET "+s.fullPath(s.paths.For("two-factor.secret-key", "/user/two-factor-secret-key")), stack.RequireAuthenticated(http.HandlerFunc(s.twoFactorSecretKey)))
		mux.Handle("GET "+s.fullPath(s.paths.For("two-factor.recovery-codes", "/user/two-factor-recovery-codes")), stack.RequireAuthenticated(http.HandlerFunc(s.recoveryCodes)))
		mux.Handle("POST "+s.fullPath(s.paths.For("two-factor.recovery-codes", "/user/two-factor-recovery-codes")), stack.RequireAuthenticated(http.HandlerFunc(s.regenerateRecoveryCodes)))
	}
}

func (s *Server) fullPath(path string) string {
	return s.config.Prefix + path
}

func (s *Server) guest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _, err := s.authManager.DefaultGuard().AuthenticateRequest(r.Context(), w, r)
		if err == nil && session != nil && !session.PendingTwoFactor {
			authmw.WriteJSONError(w, http.StatusConflict, auth.ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) loginView(w http.ResponseWriter, r *http.Request) {
	s.responses.LoginViewResponse.ToResponse(w, r, map[string]any{"view": "login"})
}

func (s *Server) registerView(w http.ResponseWriter, r *http.Request) {
	s.responses.RegisterViewResponse.ToResponse(w, r, map[string]any{"view": "register"})
}

func (s *Server) passwordResetRequestView(w http.ResponseWriter, r *http.Request) {
	s.responses.RequestPasswordResetLinkViewResponse.ToResponse(w, r, map[string]any{"view": "forgot-password"})
}

func (s *Server) newPasswordView(w http.ResponseWriter, r *http.Request) {
	s.responses.ResetPasswordViewResponse.ToResponse(w, r, map[string]any{"view": "reset-password", "token": r.PathValue("token")})
}

func (s *Server) confirmPasswordView(w http.ResponseWriter, r *http.Request) {
	s.responses.ConfirmPasswordViewResponse.ToResponse(w, r, map[string]any{"view": "confirm-password"})
}

func (s *Server) twoFactorChallengeView(w http.ResponseWriter, r *http.Request) {
	s.responses.TwoFactorChallengeViewResponse.ToResponse(w, r, map[string]any{"view": "two-factor-challenge"})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	body, err := decodeBody(r)
	if err != nil {
		s.responses.LockoutResponse.ToResponse(w, r, err)
		return
	}

	username := stringValue(body[s.config.Username])
	password := stringValue(body["password"])
	remember := boolValue(body["remember"])
	if s.config.LowercaseUsernames {
		username = strings.ToLower(username)
	}
	if username == "" || password == "" {
		s.responses.LockoutResponse.ToResponse(w, r, &auth.ValidationError{Fields: map[string]string{
			s.config.Username: "username is required",
			"password":        "password is required",
		}})
		return
	}

	if err := s.checkLimiter("login", username+"|"+clientIP(r)); err != nil {
		s.responses.LockoutResponse.ToResponse(w, r, err)
		return
	}

	user, err := s.authManager.ValidateCredentials(r.Context(), map[string]string{
		s.config.Username: username,
		"password":        password,
	})
	if err != nil {
		s.responses.LockoutResponse.ToResponse(w, r, err)
		return
	}

	if twoFactorUser, ok := user.(auth.TwoFactorAuthenticatable); ok && twoFactorUser.IsTwoFactorEnabled() {
		session, _, err := s.authManager.DefaultGuard().Login(r.Context(), w, user, remember, true)
		if err != nil {
			s.responses.LockoutResponse.ToResponse(w, r, err)
			return
		}
		s.responses.TwoFactorLoginResponse.ToResponse(w, r, map[string]any{
			"status":            "two_factor_required",
			"twoFactorRedirect": s.redirects.Redirect(user),
			"session":           session,
		})
		return
	}

	session, rememberToken, err := s.authManager.DefaultGuard().Login(r.Context(), w, user, remember, false)
	if err != nil {
		s.responses.LockoutResponse.ToResponse(w, r, err)
		return
	}
	s.responses.LoginResponse.ToResponse(w, r, map[string]any{
		"status":        "authenticated",
		"user":          user,
		"session":       session,
		"rememberToken": rememberToken,
	})
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	user, _ := authmw.CurrentUser(r)
	session, _ := authmw.CurrentSession(r)
	if err := s.authManager.DefaultGuard().Logout(r.Context(), w, session, user); err != nil {
		s.responses.LogoutResponse.ToResponse(w, r, err)
		return
	}
	s.responses.LogoutResponse.ToResponse(w, r, map[string]any{"status": "logged_out"})
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var input contracts.RegisterInput
	if err := decodeJSON(r, &input); err != nil {
		s.responses.RegisterResponse.ToResponse(w, r, err)
		return
	}

	user, err := s.createUsers.Create(r.Context(), input)
	if err != nil {
		s.responses.RegisterResponse.ToResponse(w, r, err)
		return
	}
	if s.features.Enabled(FeatureEmailVerification) {
		if err := s.verification.Send(r.Context(), user); err != nil {
			s.responses.RegisterResponse.ToResponse(w, r, err)
			return
		}
	}
	s.responses.RegisterResponse.ToResponse(w, r, map[string]any{"status": "registered", "user": user})
}

func (s *Server) sendPasswordResetLink(w http.ResponseWriter, r *http.Request) {
	var input contracts.ForgotPasswordInput
	if err := decodeJSON(r, &input); err != nil {
		s.responses.FailedPasswordResetLinkRequestResponse.ToResponse(w, r, err)
		return
	}
	if _, err := s.passwordBroker.SendResetLink(r.Context(), input.Email); err != nil {
		s.responses.FailedPasswordResetLinkRequestResponse.ToResponse(w, r, err)
		return
	}
	s.responses.SuccessfulPasswordResetLinkRequestResponse.ToResponse(w, r, map[string]any{"status": "password_reset_link_sent"})
}

func (s *Server) resetPassword(w http.ResponseWriter, r *http.Request) {
	var input contracts.ResetPasswordInput
	if err := decodeJSON(r, &input); err != nil {
		s.responses.FailedPasswordResetResponse.ToResponse(w, r, err)
		return
	}
	if input.PasswordConfirmation != input.Password {
		s.responses.FailedPasswordResetResponse.ToResponse(w, r, &auth.ValidationError{Fields: map[string]string{"password_confirmation": "password confirmation must match"}})
		return
	}
	user, err := s.passwordBroker.Reset(r.Context(), input.Email, input.Token, input.Password, resetAdapter{s.resetUsers})
	if err != nil {
		s.responses.FailedPasswordResetResponse.ToResponse(w, r, err)
		return
	}
	s.responses.PasswordResetResponse.ToResponse(w, r, map[string]any{"status": "password_reset", "user": user})
}

func (s *Server) confirmedPasswordStatus(w http.ResponseWriter, r *http.Request) {
	session, _ := authmw.CurrentSession(r)
	s.responses.PasswordConfirmedResponse.ToResponse(w, r, map[string]any{
		"confirmed": session.PasswordConfirmedAt != nil && !auth.ExpiredPasswordConfirmation(s.authManager.Config(), session.PasswordConfirmedAt, s.clock.Now()),
	})
}

func (s *Server) confirmPassword(w http.ResponseWriter, r *http.Request) {
	var input contracts.ConfirmPasswordInput
	if err := decodeJSON(r, &input); err != nil {
		s.responses.FailedPasswordConfirmationResponse.ToResponse(w, r, err)
		return
	}
	user, _ := authmw.CurrentUser(r)
	session, _ := authmw.CurrentSession(r)
	if err := s.authManager.Hasher().Compare(r.Context(), user.GetAuthPassword(), input.Password); err != nil {
		s.responses.FailedPasswordConfirmationResponse.ToResponse(w, r, auth.ErrInvalidCredentials)
		return
	}
	now := s.clock.Now()
	session.PasswordConfirmedAt = &now
	session.LastSeenAt = now
	if err := s.authManager.DefaultGuard().UpdateSession(r.Context(), session); err != nil {
		s.responses.FailedPasswordConfirmationResponse.ToResponse(w, r, err)
		return
	}
	s.responses.PasswordConfirmedResponse.ToResponse(w, r, map[string]any{"status": "password_confirmed"})
}

func (s *Server) emailVerificationPrompt(w http.ResponseWriter, r *http.Request) {
	s.responses.VerifyEmailViewResponse.ToResponse(w, r, map[string]any{"view": "verify-email"})
}

func (s *Server) sendEmailVerification(w http.ResponseWriter, r *http.Request) {
	user, _ := authmw.CurrentUser(r)
	if err := s.verification.Send(r.Context(), user); err != nil {
		s.responses.EmailVerificationNotificationSentResponse.ToResponse(w, r, err)
		return
	}
	s.responses.EmailVerificationNotificationSentResponse.ToResponse(w, r, map[string]any{"status": "verification_link_sent"})
}

func (s *Server) verifyEmail(w http.ResponseWriter, r *http.Request) {
	expiresAt, err := strconv.ParseInt(r.URL.Query().Get("expires"), 10, 64)
	if err != nil {
		s.responses.VerifyEmailResponse.ToResponse(w, r, auth.ErrInvalidToken)
		return
	}
	user, err := s.verification.Verify(r.Context(), r.PathValue("id"), r.PathValue("hash"), expiresAt, r.URL.Query().Get("signature"))
	if err != nil {
		s.responses.VerifyEmailResponse.ToResponse(w, r, err)
		return
	}
	s.responses.VerifyEmailResponse.ToResponse(w, r, map[string]any{"status": "email_verified", "user": user})
}

func (s *Server) completeTwoFactorLogin(w http.ResponseWriter, r *http.Request) {
	var input contracts.TwoFactorChallengeInput
	if err := decodeJSON(r, &input); err != nil {
		s.responses.FailedTwoFactorLoginResponse.ToResponse(w, r, err)
		return
	}

	session, user, err := s.authManager.DefaultGuard().AuthenticateRequest(r.Context(), w, r)
	if err != nil || session == nil || !session.PendingTwoFactor {
		s.responses.FailedTwoFactorLoginResponse.ToResponse(w, r, auth.ErrUnauthorized)
		return
	}
	twoFactorUser, ok := user.(auth.TwoFactorAuthenticatable)
	if !ok || !twoFactorUser.IsTwoFactorEnabled() {
		s.responses.FailedTwoFactorLoginResponse.ToResponse(w, r, auth.ErrUnauthorized)
		return
	}
	if err := s.checkLimiter("two-factor", session.ID); err != nil {
		s.responses.LockoutResponse.ToResponse(w, r, err)
		return
	}

	valid := false
	if input.Code != "" {
		valid = s.twoFactor.Validate(twoFactorUser.GetTwoFactorSecret(), input.Code, s.clock.Now(), 1)
	}
	if !valid && input.RecoveryCode != "" {
		codes := twoFactorUser.GetTwoFactorRecoveryCodes()
		for index, code := range codes {
			if code == strings.TrimSpace(input.RecoveryCode) {
				codes = append(codes[:index], codes[index+1:]...)
				twoFactorUser.SetTwoFactorRecoveryCodes(codes)
				if err := s.users.Update(r.Context(), user); err != nil {
					s.responses.FailedTwoFactorLoginResponse.ToResponse(w, r, err)
					return
				}
				valid = true
				break
			}
		}
	}
	if !valid {
		s.responses.FailedTwoFactorLoginResponse.ToResponse(w, r, auth.ErrTwoFactorInvalid)
		return
	}

	rememberToken, err := s.authManager.DefaultGuard().CompleteTwoFactor(r.Context(), w, session, user)
	if err != nil {
		s.responses.FailedTwoFactorLoginResponse.ToResponse(w, r, err)
		return
	}
	s.responses.TwoFactorLoginResponse.ToResponse(w, r, map[string]any{
		"status":        "authenticated",
		"user":          user,
		"session":       session,
		"rememberToken": rememberToken,
	})
}

func (s *Server) enableTwoFactor(w http.ResponseWriter, r *http.Request) {
	user, _ := authmw.CurrentUser(r)
	profile, ok := user.(auth.UserProfile)
	if !ok {
		s.responses.TwoFactorEnabledResponse.ToResponse(w, r, auth.ErrUnauthorized)
		return
	}
	twoFactorUser, ok := user.(auth.TwoFactorAuthenticatable)
	if !ok {
		s.responses.TwoFactorEnabledResponse.ToResponse(w, r, auth.ErrUnauthorized)
		return
	}

	secret, err := s.twoFactor.GenerateSecret()
	if err != nil {
		s.responses.TwoFactorEnabledResponse.ToResponse(w, r, err)
		return
	}
	recoveryCodes, err := s.twoFactor.GenerateRecoveryCodes()
	if err != nil {
		s.responses.TwoFactorEnabledResponse.ToResponse(w, r, err)
		return
	}
	twoFactorUser.SetTwoFactorSecret(secret)
	twoFactorUser.SetTwoFactorRecoveryCodes(recoveryCodes)
	if !s.features.OptionEnabled(FeatureTwoFactorAuthentication, "confirm") {
		now := s.clock.Now()
		twoFactorUser.SetTwoFactorConfirmedAt(&now)
	}
	if err := s.users.Update(r.Context(), user); err != nil {
		s.responses.TwoFactorEnabledResponse.ToResponse(w, r, err)
		return
	}
	s.responses.TwoFactorEnabledResponse.ToResponse(w, r, map[string]any{
		"status":        "two_factor_enabled",
		"secret":        secret,
		"recoveryCodes": recoveryCodes,
		"otpAuthUrl":    s.twoFactor.OTPAuthURL("gollin", profile.GetEmail(), secret),
	})
}

func (s *Server) confirmTwoFactor(w http.ResponseWriter, r *http.Request) {
	var input contracts.TwoFactorChallengeInput
	if err := decodeJSON(r, &input); err != nil {
		s.responses.TwoFactorConfirmedResponse.ToResponse(w, r, err)
		return
	}
	user, _ := authmw.CurrentUser(r)
	twoFactorUser, ok := user.(auth.TwoFactorAuthenticatable)
	if !ok || !twoFactorUser.IsTwoFactorEnabled() {
		s.responses.TwoFactorConfirmedResponse.ToResponse(w, r, auth.ErrUnauthorized)
		return
	}
	if !s.twoFactor.Validate(twoFactorUser.GetTwoFactorSecret(), input.Code, s.clock.Now(), 1) {
		s.responses.TwoFactorConfirmedResponse.ToResponse(w, r, auth.ErrTwoFactorInvalid)
		return
	}
	now := s.clock.Now()
	twoFactorUser.SetTwoFactorConfirmedAt(&now)
	if err := s.users.Update(r.Context(), user); err != nil {
		s.responses.TwoFactorConfirmedResponse.ToResponse(w, r, err)
		return
	}
	s.responses.TwoFactorConfirmedResponse.ToResponse(w, r, map[string]any{"status": "two_factor_confirmed"})
}

func (s *Server) disableTwoFactor(w http.ResponseWriter, r *http.Request) {
	user, _ := authmw.CurrentUser(r)
	twoFactorUser, ok := user.(auth.TwoFactorAuthenticatable)
	if !ok {
		s.responses.TwoFactorDisabledResponse.ToResponse(w, r, auth.ErrUnauthorized)
		return
	}
	twoFactorUser.SetTwoFactorEnabled(false)
	if err := s.users.Update(r.Context(), user); err != nil {
		s.responses.TwoFactorDisabledResponse.ToResponse(w, r, err)
		return
	}
	s.responses.TwoFactorDisabledResponse.ToResponse(w, r, map[string]any{"status": "two_factor_disabled"})
}

func (s *Server) twoFactorQRCode(w http.ResponseWriter, r *http.Request) {
	user, _ := authmw.CurrentUser(r)
	profile, _ := user.(auth.UserProfile)
	twoFactorUser, ok := user.(auth.TwoFactorAuthenticatable)
	if !ok {
		s.responses.TwoFactorEnabledResponse.ToResponse(w, r, auth.ErrUnauthorized)
		return
	}
	s.responses.TwoFactorEnabledResponse.ToResponse(w, r, map[string]any{
		"otpAuthUrl": s.twoFactor.OTPAuthURL("gollin", profile.GetEmail(), twoFactorUser.GetTwoFactorSecret()),
	})
}

func (s *Server) twoFactorSecretKey(w http.ResponseWriter, r *http.Request) {
	user, _ := authmw.CurrentUser(r)
	twoFactorUser, ok := user.(auth.TwoFactorAuthenticatable)
	if !ok {
		s.responses.TwoFactorEnabledResponse.ToResponse(w, r, auth.ErrUnauthorized)
		return
	}
	s.responses.TwoFactorEnabledResponse.ToResponse(w, r, map[string]any{"secret": twoFactorUser.GetTwoFactorSecret()})
}

func (s *Server) recoveryCodes(w http.ResponseWriter, r *http.Request) {
	user, _ := authmw.CurrentUser(r)
	twoFactorUser, ok := user.(auth.TwoFactorAuthenticatable)
	if !ok {
		s.responses.RecoveryCodesGeneratedResponse.ToResponse(w, r, auth.ErrUnauthorized)
		return
	}
	s.responses.RecoveryCodesGeneratedResponse.ToResponse(w, r, map[string]any{"recoveryCodes": twoFactorUser.GetTwoFactorRecoveryCodes()})
}

func (s *Server) regenerateRecoveryCodes(w http.ResponseWriter, r *http.Request) {
	user, _ := authmw.CurrentUser(r)
	twoFactorUser, ok := user.(auth.TwoFactorAuthenticatable)
	if !ok {
		s.responses.RecoveryCodesGeneratedResponse.ToResponse(w, r, auth.ErrUnauthorized)
		return
	}
	codes, err := s.twoFactor.GenerateRecoveryCodes()
	if err != nil {
		s.responses.RecoveryCodesGeneratedResponse.ToResponse(w, r, err)
		return
	}
	twoFactorUser.SetTwoFactorRecoveryCodes(codes)
	if err := s.users.Update(r.Context(), user); err != nil {
		s.responses.RecoveryCodesGeneratedResponse.ToResponse(w, r, err)
		return
	}
	s.responses.RecoveryCodesGeneratedResponse.ToResponse(w, r, map[string]any{"status": "recovery_codes_generated", "recoveryCodes": codes})
}

func (s *Server) updateProfileInformation(w http.ResponseWriter, r *http.Request) {
	var input contracts.UpdateProfileInformationInput
	if err := decodeJSON(r, &input); err != nil {
		s.responses.ProfileInformationUpdatedResponse.ToResponse(w, r, err)
		return
	}
	user, _ := authmw.CurrentUser(r)
	if err := s.updateProfiles.Update(r.Context(), user, input); err != nil {
		s.responses.ProfileInformationUpdatedResponse.ToResponse(w, r, err)
		return
	}
	s.responses.ProfileInformationUpdatedResponse.ToResponse(w, r, map[string]any{"status": "profile_information_updated", "user": user})
}

func (s *Server) updatePassword(w http.ResponseWriter, r *http.Request) {
	var input contracts.UpdatePasswordInput
	if err := decodeJSON(r, &input); err != nil {
		s.responses.PasswordUpdateResponse.ToResponse(w, r, err)
		return
	}
	user, _ := authmw.CurrentUser(r)
	if err := s.updatePasswords.Update(r.Context(), user, input); err != nil {
		s.responses.PasswordUpdateResponse.ToResponse(w, r, err)
		return
	}
	s.responses.PasswordUpdateResponse.ToResponse(w, r, map[string]any{"status": "password_updated"})
}

func (s *Server) checkLimiter(scope string, key string) error {
	spec := s.config.Limiters[scope]
	limit, window, err := parseLimiter(spec)
	if err != nil {
		return err
	}
	allowed, retryAfter := s.limiter.Allow(scope+":"+key, limit, window, s.clock.Now())
	if !allowed {
		return &auth.ThrottleError{Scope: scope, RetryAfter: retryAfter}
	}
	return nil
}

func decodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return &auth.ValidationError{Fields: map[string]string{"body": "request body is required"}}
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(dst); err != nil {
		return &auth.ValidationError{Fields: map[string]string{"body": err.Error()}}
	}
	return nil
}

func decodeBody(r *http.Request) (map[string]any, error) {
	body := map[string]any{}
	if err := decodeJSON(r, &body); err != nil {
		return nil, err
	}
	return body, nil
}

func stringValue(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	default:
		return ""
	}
}

func boolValue(value any) bool {
	typed, _ := value.(bool)
	return typed
}

func clientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

type resetAdapter struct {
	action contracts.ResetsUserPasswords
}

func (a resetAdapter) Reset(ctx context.Context, user auth.Authenticatable, password string) error {
	return a.action.Reset(ctx, user, password)
}
