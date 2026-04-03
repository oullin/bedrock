package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/gollin/packages/auth/internal/secure"
	"github.com/gollin/packages/auth/internal/throttle"
	"github.com/gollin/packages/auth/internal/totp"
)

const defaultGuardName = "web"

// Dependencies wires concrete stores and helpers into a Manager.
type Dependencies struct {
	UserStore          UserStore
	SessionStore       SessionStore
	PasswordResetStore PasswordResetStore
	TwoFactorStore     TwoFactorStore
	Mailer             Mailer
	PasswordHasher     PasswordHasher
	LinkSigner         LinkSigner
	Clock              Clock
	IDGenerator        IDGenerator
	Logger             Logger
	Observer           Observer
}

// Manager coordinates guards, stores, and AuthFlows-style auth flows.
type Manager struct {
	config                   Config
	users                    UserStore
	sessions                 SessionStore
	passwordResets           PasswordResetStore
	twoFactor                TwoFactorStore
	mailer                   Mailer
	hasher                   PasswordHasher
	signer                   LinkSigner
	clock                    Clock
	ids                      IDGenerator
	logger                   Logger
	observer                 Observer
	limiter                  *throttle.Limiter
	guards                   map[string]*SessionGuard
	defaultGuard             string
	authenticateUsing        AuthenticateUsingFunc
	createUsersUsing         CreateUsersUsingFunc
	resetUserPasswordsUsing  ResetUserPasswordsUsingFunc
	updateUserPasswordsUsing UpdateUserPasswordsUsingFunc
}

// NewManager creates a new auth manager with a default session guard.
func NewManager(config Config, deps Dependencies) (*Manager, error) {
	cfg := config.withDefaults()

	if deps.UserStore == nil {
		deps.UserStore = NewInMemoryUserStore()
	}
	if deps.SessionStore == nil {
		deps.SessionStore = NewInMemorySessionStore()
	}
	if deps.PasswordResetStore == nil {
		deps.PasswordResetStore = NewInMemoryPasswordResetStore()
	}
	if deps.TwoFactorStore == nil {
		deps.TwoFactorStore = NewInMemoryTwoFactorStore()
	}
	if deps.Mailer == nil {
		deps.Mailer = &InMemoryMailer{}
	}
	if deps.PasswordHasher == nil {
		deps.PasswordHasher = DefaultPasswordHasher{}
	}
	if deps.LinkSigner == nil {
		deps.LinkSigner = HMACLinkSigner{
			Key:   cfg.SigningKey,
			Clock: deps.Clock,
		}
	}
	if deps.Clock == nil {
		deps.Clock = SystemClock{}
	}
	if deps.IDGenerator == nil {
		deps.IDGenerator = RandomIDGenerator{}
	}
	if deps.Logger == nil {
		deps.Logger = NoopLogger{}
	}
	if deps.Observer == nil {
		deps.Observer = NoopObserver{}
	}

	manager := &Manager{
		config:         cfg,
		users:          deps.UserStore,
		sessions:       deps.SessionStore,
		passwordResets: deps.PasswordResetStore,
		twoFactor:      deps.TwoFactorStore,
		mailer:         deps.Mailer,
		hasher:         deps.PasswordHasher,
		signer:         deps.LinkSigner,
		clock:          deps.Clock,
		ids:            deps.IDGenerator,
		logger:         deps.Logger,
		observer:       deps.Observer,
		limiter:        throttle.NewLimiter(),
		guards:         make(map[string]*SessionGuard),
		defaultGuard:   defaultGuardName,
	}

	manager.guards[defaultGuardName] = NewSessionGuard(cfg, deps.UserStore, deps.SessionStore, deps.IDGenerator, deps.Clock, deps.Logger)
	manager.authenticateUsing = manager.defaultAuthenticate
	manager.createUsersUsing = manager.defaultCreateUser
	manager.resetUserPasswordsUsing = manager.defaultUpdatePassword
	manager.updateUserPasswordsUsing = manager.defaultUpdatePassword

	return manager, nil
}

// Config returns the active auth configuration.
func (m *Manager) Config() Config {
	return m.config
}

// DefaultGuard returns the default session guard.
func (m *Manager) DefaultGuard() *SessionGuard {
	return m.guards[m.defaultGuard]
}

// Guard returns a named guard when present.
func (m *Manager) Guard(name string) (*SessionGuard, bool) {
	guard, ok := m.guards[name]
	return guard, ok
}

// SetAuthenticateUsing replaces the login callback.
func (m *Manager) SetAuthenticateUsing(fn AuthenticateUsingFunc) {
	if fn != nil {
		m.authenticateUsing = fn
	}
}

// SetCreateUsersUsing replaces the user registration callback.
func (m *Manager) SetCreateUsersUsing(fn CreateUsersUsingFunc) {
	if fn != nil {
		m.createUsersUsing = fn
	}
}

// SetResetUserPasswordsUsing replaces the password reset callback.
func (m *Manager) SetResetUserPasswordsUsing(fn ResetUserPasswordsUsingFunc) {
	if fn != nil {
		m.resetUserPasswordsUsing = fn
	}
}

// SetUpdateUserPasswordsUsing replaces the password update callback.
func (m *Manager) SetUpdateUserPasswordsUsing(fn UpdateUserPasswordsUsingFunc) {
	if fn != nil {
		m.updateUserPasswordsUsing = fn
	}
}

// Register creates a new user and dispatches an email verification link.
func (m *Manager) Register(ctx context.Context, input RegisterInput) (*User, error) {
	if err := validateRegisterInput(input); err != nil {
		return nil, err
	}

	user, err := m.createUsersUsing(ctx, m, input)
	if err != nil {
		return nil, err
	}

	if err := m.SendEmailVerification(ctx, user); err != nil {
		return nil, err
	}

	m.observer.Record(ctx, "auth.registered", map[string]string{"user_id": user.ID})
	return user, nil
}

// Login verifies credentials and creates either an authenticated or pending session.
func (m *Manager) Login(ctx context.Context, input LoginInput, throttleKey string) (*LoginResult, error) {
	if err := validateLoginInput(input); err != nil {
		return nil, err
	}
	if err := m.checkThrottle("login:"+throttleKey, m.config.Throttle.LoginLimit, m.config.Throttle.LoginWindow); err != nil {
		return nil, err
	}

	user, err := m.authenticateUsing(ctx, m, input.Email, input.Password)
	if err != nil {
		m.observer.Record(ctx, "auth.login_failed", map[string]string{"email": normalizeIdentifier(input.Email)})
		return nil, err
	}

	state, err := m.twoFactor.FindByUserID(ctx, user.ID)
	hasTwoFactor := err == nil && state != nil && user.TwoFactorEnabled
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	now := m.clock.Now()
	session := &Session{
		ID:               m.ids.NewID(),
		UserID:           user.ID,
		PendingTwoFactor: hasTwoFactor,
		PendingRemember:  input.Remember,
		LastSeenAt:       now,
		CreatedAt:        now,
		ExpiresAt:        now.Add(m.config.SessionLifetime),
	}
	if !hasTwoFactor {
		session.AuthenticatedAt = &now
	}

	if err := m.sessions.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	result := &LoginResult{
		User:              user,
		Session:           session,
		RequiresTwoFactor: hasTwoFactor,
	}

	if !hasTwoFactor && input.Remember {
		token, err := m.issueRememberToken(ctx, user)
		if err != nil {
			return nil, err
		}
		result.RememberToken = token
	}

	m.observer.Record(ctx, "auth.login", map[string]string{"user_id": user.ID})
	return result, nil
}

// Logout invalidates the active session and remember-me token.
func (m *Manager) Logout(ctx context.Context, user *User, session *Session) error {
	if session != nil {
		if err := m.sessions.Delete(ctx, session.ID); err != nil {
			return fmt.Errorf("delete session: %w", err)
		}
	}

	if user != nil && user.RememberTokenHash != "" {
		user.RememberTokenHash = ""
		user.UpdatedAt = m.clock.Now()
		if err := m.users.Update(ctx, user); err != nil {
			return err
		}
	}

	m.observer.Record(ctx, "auth.logout", map[string]string{"user_id": safeUserID(user)})
	return nil
}

// SendPasswordResetLink issues and emails a password reset token.
func (m *Manager) SendPasswordResetLink(ctx context.Context, input ForgotPasswordInput, throttleKey string) error {
	if strings.TrimSpace(input.Email) == "" {
		return &ValidationError{Fields: map[string]string{"email": "email is required"}}
	}
	if err := m.checkThrottle("password-reset:"+throttleKey, m.config.Throttle.PasswordResetLimit, m.config.Throttle.PasswordResetWindow); err != nil {
		return err
	}

	user, err := m.users.FindByIdentifier(ctx, input.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil
		}
		return err
	}

	token, err := secure.RandomString(24)
	if err != nil {
		return fmt.Errorf("generate reset token: %w", err)
	}

	now := m.clock.Now()
	hashed := secure.HashString(token)
	if err := m.passwordResets.Save(ctx, &PasswordResetToken{
		UserID:    user.ID,
		TokenHash: hashed,
		CreatedAt: now,
		ExpiresAt: now.Add(m.config.PasswordResetTTL),
	}); err != nil {
		return fmt.Errorf("save reset token: %w", err)
	}

	body := fmt.Sprintf("Use this password reset token for %s: %s", user.Email, token)
	if m.config.BaseURL != "" {
		body += fmt.Sprintf("\nReset endpoint: %s/reset-password", strings.TrimRight(m.config.BaseURL, "/"))
	}

	if err := m.mailer.Send(ctx, MailMessage{
		To:      user.Email,
		Subject: "Reset your password",
		Body:    body,
		Metadata: map[string]string{
			"token": token,
		},
	}); err != nil {
		return fmt.Errorf("send reset email: %w", err)
	}

	m.observer.Record(ctx, "auth.password_reset_requested", map[string]string{"user_id": user.ID})
	return nil
}

// ResetPassword updates a user's password using a valid reset token.
func (m *Manager) ResetPassword(ctx context.Context, input ResetPasswordInput) (*User, error) {
	if err := validateResetPasswordInput(input); err != nil {
		return nil, err
	}

	tokenHash := secure.HashString(input.Token)
	resetToken, err := m.passwordResets.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if m.clock.Now().After(resetToken.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	user, err := m.users.FindByIdentifier(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if user.ID != resetToken.UserID {
		return nil, ErrInvalidToken
	}

	if err := m.resetUserPasswordsUsing(ctx, m, user, input.Password); err != nil {
		return nil, err
	}
	if err := m.passwordResets.DeleteByTokenHash(ctx, tokenHash); err != nil {
		return nil, fmt.Errorf("delete reset token: %w", err)
	}

	user.RememberTokenHash = ""
	user.UpdatedAt = m.clock.Now()
	if err := m.users.Update(ctx, user); err != nil {
		return nil, err
	}

	m.observer.Record(ctx, "auth.password_reset", map[string]string{"user_id": user.ID})
	return user, nil
}

// SendEmailVerification dispatches a signed email verification message.
func (m *Manager) SendEmailVerification(ctx context.Context, user *User) error {
	if user == nil {
		return ErrUserNotFound
	}
	if err := m.checkThrottle("verification:"+user.ID, m.config.Throttle.VerificationLimit, m.config.Throttle.VerificationWindow); err != nil {
		return err
	}

	expiresAt := m.clock.Now().Add(m.config.VerificationTTL)
	emailHash := secure.EmailHash(normalizeIdentifier(user.Email))
	signature, err := m.signer.Sign(ctx, "email-verification", []string{user.ID, emailHash}, expiresAt.Unix())
	if err != nil {
		return fmt.Errorf("sign verification link: %w", err)
	}

	base := strings.TrimRight(m.config.BaseURL, "/")
	link := fmt.Sprintf("/verify-email/%s/%s?expires=%d&signature=%s", user.ID, emailHash, expiresAt.Unix(), signature)
	if base != "" {
		link = base + link
	}

	body := fmt.Sprintf("Verify your email address by visiting %s", link)
	if err := m.mailer.Send(ctx, MailMessage{
		To:      user.Email,
		Subject: "Verify your email address",
		Body:    body,
		Metadata: map[string]string{
			"link":      link,
			"expiresAt": strconv.FormatInt(expiresAt.Unix(), 10),
		},
	}); err != nil {
		return fmt.Errorf("send verification email: %w", err)
	}

	return nil
}

// VerifyEmail validates a signed verification link and marks the user as verified.
func (m *Manager) VerifyEmail(ctx context.Context, userID string, emailHash string, expiresAt int64, signature string) (*User, error) {
	user, err := m.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	expectedHash := secure.EmailHash(normalizeIdentifier(user.Email))
	if expectedHash != emailHash {
		return nil, ErrEmailVerificationInvalid
	}

	if err := m.signer.Verify(ctx, "email-verification", []string{user.ID, emailHash}, expiresAt, signature); err != nil {
		return nil, err
	}

	now := m.clock.Now()
	user.EmailVerifiedAt = &now
	user.UpdatedAt = now
	if err := m.users.Update(ctx, user); err != nil {
		return nil, err
	}

	m.observer.Record(ctx, "auth.email_verified", map[string]string{"user_id": user.ID})
	return user, nil
}

// ConfirmPassword refreshes the password confirmation timestamp for the session.
func (m *Manager) ConfirmPassword(ctx context.Context, user *User, session *Session, input ConfirmPasswordInput) error {
	if user == nil || session == nil {
		return ErrUnauthorized
	}
	if strings.TrimSpace(input.Password) == "" {
		return &ValidationError{Fields: map[string]string{"password": "password is required"}}
	}
	if err := m.hasher.Compare(ctx, user.PasswordHash, input.Password); err != nil {
		return ErrInvalidCredentials
	}

	now := m.clock.Now()
	session.PasswordConfirmedAt = &now
	session.LastSeenAt = now
	if err := m.sessions.Update(ctx, session); err != nil {
		return fmt.Errorf("update session: %w", err)
	}

	m.observer.Record(ctx, "auth.password_confirmed", map[string]string{"user_id": user.ID})
	return nil
}

// EnableTwoFactor enables TOTP login after password confirmation.
func (m *Manager) EnableTwoFactor(ctx context.Context, user *User, session *Session) (*EnableTwoFactorResult, error) {
	if err := m.requirePasswordConfirmation(session); err != nil {
		return nil, err
	}

	secret, err := totp.GenerateSecret()
	if err != nil {
		return nil, fmt.Errorf("generate TOTP secret: %w", err)
	}

	recoveryCodes := make([]string, 0, 8)
	for range 8 {
		code, err := secure.RandomString(8)
		if err != nil {
			return nil, err
		}
		recoveryCodes = append(recoveryCodes, code[:10])
	}

	now := m.clock.Now()
	if err := m.twoFactor.Save(ctx, &TwoFactorState{
		UserID:        user.ID,
		Secret:        secret,
		RecoveryCodes: recoveryCodes,
		EnabledAt:     now,
		UpdatedAt:     now,
	}); err != nil {
		return nil, err
	}

	user.TwoFactorEnabled = true
	user.UpdatedAt = now
	if err := m.users.Update(ctx, user); err != nil {
		return nil, err
	}

	m.observer.Record(ctx, "auth.two_factor_enabled", map[string]string{"user_id": user.ID})
	return &EnableTwoFactorResult{
		Secret:        secret,
		RecoveryCodes: recoveryCodes,
		OTPAuthURL:    totp.OTPAuthURL(m.config.TwoFactorIssuer, user.Email, secret),
	}, nil
}

// DisableTwoFactor removes all two-factor state for a user.
func (m *Manager) DisableTwoFactor(ctx context.Context, user *User, session *Session) error {
	if err := m.requirePasswordConfirmation(session); err != nil {
		return err
	}

	if err := m.twoFactor.Delete(ctx, user.ID); err != nil {
		return err
	}
	user.TwoFactorEnabled = false
	user.UpdatedAt = m.clock.Now()
	if err := m.users.Update(ctx, user); err != nil {
		return err
	}

	m.observer.Record(ctx, "auth.two_factor_disabled", map[string]string{"user_id": user.ID})
	return nil
}

// RecoveryCodes returns the current recovery codes for a user.
func (m *Manager) RecoveryCodes(ctx context.Context, user *User, session *Session) ([]string, error) {
	if err := m.requirePasswordConfirmation(session); err != nil {
		return nil, err
	}

	state, err := m.twoFactor.FindByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return append([]string(nil), state.RecoveryCodes...), nil
}

// RegenerateRecoveryCodes replaces the current recovery codes for a user.
func (m *Manager) RegenerateRecoveryCodes(ctx context.Context, user *User, session *Session) ([]string, error) {
	if err := m.requirePasswordConfirmation(session); err != nil {
		return nil, err
	}

	state, err := m.twoFactor.FindByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	codes := make([]string, 0, 8)
	for range 8 {
		code, err := secure.RandomString(8)
		if err != nil {
			return nil, err
		}
		codes = append(codes, code[:10])
	}

	state.RecoveryCodes = codes
	state.UpdatedAt = m.clock.Now()
	if err := m.twoFactor.Save(ctx, state); err != nil {
		return nil, err
	}

	m.observer.Record(ctx, "auth.two_factor_recovery_regenerated", map[string]string{"user_id": user.ID})
	return append([]string(nil), codes...), nil
}

// ChallengeTwoFactor completes a pending two-factor session.
func (m *Manager) ChallengeTwoFactor(ctx context.Context, session *Session, input TwoFactorChallengeInput) (*LoginResult, error) {
	if session == nil || !session.PendingTwoFactor {
		return nil, ErrUnauthorized
	}
	if strings.TrimSpace(input.Code) == "" && strings.TrimSpace(input.RecoveryCode) == "" {
		return nil, &ValidationError{Fields: map[string]string{"code": "code or recovery_code is required"}}
	}
	if err := m.checkThrottle("two-factor:"+session.ID, m.config.Throttle.TwoFactorLimit, m.config.Throttle.TwoFactorWindow); err != nil {
		return nil, err
	}

	user, err := m.users.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}
	state, err := m.twoFactor.FindByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	valid := false
	if input.Code != "" {
		valid = totp.Validate(state.Secret, input.Code, m.clock.Now(), 1)
	}
	if !valid && input.RecoveryCode != "" {
		for index, code := range state.RecoveryCodes {
			if code == strings.TrimSpace(input.RecoveryCode) {
				state.RecoveryCodes = append(state.RecoveryCodes[:index], state.RecoveryCodes[index+1:]...)
				if err := m.twoFactor.Save(ctx, state); err != nil {
					return nil, err
				}
				valid = true
				break
			}
		}
	}
	if !valid {
		return nil, ErrTwoFactorInvalid
	}

	now := m.clock.Now()
	session.PendingTwoFactor = false
	session.AuthenticatedAt = &now
	session.LastSeenAt = now
	if err := m.sessions.Update(ctx, session); err != nil {
		return nil, err
	}

	result := &LoginResult{
		User:              user,
		Session:           session,
		PasswordConfirmed: session.PasswordConfirmedAt != nil,
	}
	if session.PendingRemember {
		token, err := m.issueRememberToken(ctx, user)
		if err != nil {
			return nil, err
		}
		result.RememberToken = token
		session.PendingRemember = false
		if err := m.sessions.Update(ctx, session); err != nil {
			return nil, err
		}
	}

	m.observer.Record(ctx, "auth.two_factor_challenged", map[string]string{"user_id": user.ID})
	return result, nil
}

// Authenticate loads and refreshes the current request session.
func (m *Manager) Authenticate(ctx context.Context, w http.ResponseWriter, r *http.Request) (*Session, *User, error) {
	return m.DefaultGuard().AuthenticateRequest(ctx, w, r)
}

func (m *Manager) checkThrottle(key string, limit int, window time.Duration) error {
	allowed, retryAfter := m.limiter.Allow(key, limit, window, m.clock.Now())
	if !allowed {
		return &ThrottleError{
			Scope:      key,
			RetryAfter: retryAfter,
		}
	}

	return nil
}

func (m *Manager) issueRememberToken(ctx context.Context, user *User) (string, error) {
	token, err := secure.RandomString(24)
	if err != nil {
		return "", fmt.Errorf("generate remember token: %w", err)
	}

	user.RememberTokenHash = secure.HashString(token)
	user.UpdatedAt = m.clock.Now()
	if err := m.users.Update(ctx, user); err != nil {
		return "", fmt.Errorf("save remember token: %w", err)
	}

	return token, nil
}

func (m *Manager) requirePasswordConfirmation(session *Session) error {
	if session == nil || session.PasswordConfirmedAt == nil {
		return ErrPasswordConfirmationRequired
	}
	if m.clock.Now().After(session.PasswordConfirmedAt.Add(m.config.PasswordConfirmationTimeout)) {
		return ErrPasswordConfirmationRequired
	}
	return nil
}

func (m *Manager) defaultAuthenticate(ctx context.Context, _ *Manager, identifier string, password string) (*User, error) {
	user, err := m.users.FindByIdentifier(ctx, identifier)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := m.hasher.Compare(ctx, user.PasswordHash, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (m *Manager) defaultCreateUser(ctx context.Context, _ *Manager, input RegisterInput) (*User, error) {
	if _, err := m.users.FindByIdentifier(ctx, input.Email); err == nil {
		return nil, ErrUserExists
	} else if !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	passwordHash, err := m.hasher.Hash(ctx, input.Password)
	if err != nil {
		return nil, err
	}

	now := m.clock.Now()
	user := &User{
		ID:           m.ids.NewID(),
		Email:        normalizeIdentifier(input.Email),
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := m.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (m *Manager) defaultUpdatePassword(ctx context.Context, _ *Manager, user *User, password string) error {
	passwordHash, err := m.hasher.Hash(ctx, password)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash
	user.UpdatedAt = m.clock.Now()
	return m.users.Update(ctx, user)
}

func validateRegisterInput(input RegisterInput) error {
	fields := map[string]string{}
	if !looksLikeEmail(input.Email) {
		fields["email"] = "email must be a valid email address"
	}
	if len(input.Password) < 8 {
		fields["password"] = "password must be at least 8 characters"
	}
	if input.PasswordConfirmation != input.Password {
		fields["password_confirmation"] = "password confirmation must match"
	}
	if len(fields) > 0 {
		return &ValidationError{Fields: fields}
	}
	return nil
}

func validateLoginInput(input LoginInput) error {
	fields := map[string]string{}
	if strings.TrimSpace(input.Email) == "" {
		fields["email"] = "email is required"
	}
	if strings.TrimSpace(input.Password) == "" {
		fields["password"] = "password is required"
	}
	if len(fields) > 0 {
		return &ValidationError{Fields: fields}
	}
	return nil
}

func validateResetPasswordInput(input ResetPasswordInput) error {
	fields := map[string]string{}
	if !looksLikeEmail(input.Email) {
		fields["email"] = "email must be a valid email address"
	}
	if strings.TrimSpace(input.Token) == "" {
		fields["token"] = "token is required"
	}
	if len(input.Password) < 8 {
		fields["password"] = "password must be at least 8 characters"
	}
	if input.PasswordConfirmation != input.Password {
		fields["password_confirmation"] = "password confirmation must match"
	}
	if len(fields) > 0 {
		return &ValidationError{Fields: fields}
	}
	return nil
}

func safeUserID(user *User) string {
	if user == nil {
		return ""
	}
	return user.ID
}

func looksLikeEmail(value string) bool {
	if strings.TrimSpace(value) == "" {
		return false
	}
	_, err := mail.ParseAddress(value)
	return err == nil
}

// DefaultPasswordHasher hashes passwords with PBKDF2-HMAC-SHA256.
type DefaultPasswordHasher struct{}

// Hash encodes a password.
func (DefaultPasswordHasher) Hash(_ context.Context, password string) (string, error) {
	salt, err := secure.RandomString(16)
	if err != nil {
		return "", err
	}

	iterations := 120_000
	derived := secure.PBKDF2([]byte(password), []byte(salt), iterations, 32, sha256.New)
	return fmt.Sprintf(
		"pbkdf2_sha256$%d$%s$%s",
		iterations,
		base64.RawURLEncoding.EncodeToString([]byte(salt)),
		base64.RawURLEncoding.EncodeToString(derived),
	), nil
}

// Compare validates a password against an encoded hash.
func (DefaultPasswordHasher) Compare(_ context.Context, encodedPassword string, password string) error {
	parts := strings.Split(encodedPassword, "$")
	if len(parts) != 4 || parts[0] != "pbkdf2_sha256" {
		return ErrInvalidCredentials
	}

	iterations, err := strconv.Atoi(parts[1])
	if err != nil {
		return ErrInvalidCredentials
	}
	salt, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return ErrInvalidCredentials
	}
	expected, err := base64.RawURLEncoding.DecodeString(parts[3])
	if err != nil {
		return ErrInvalidCredentials
	}

	actual := secure.PBKDF2([]byte(password), salt, iterations, len(expected), sha256.New)
	if !hmac.Equal(actual, expected) {
		return ErrInvalidCredentials
	}
	return nil
}

// HMACLinkSigner signs verification payloads with HMAC-SHA256.
type HMACLinkSigner struct {
	Key   []byte
	Clock Clock
}

// Sign signs a time-bound payload.
func (s HMACLinkSigner) Sign(_ context.Context, purpose string, values []string, expiresAt int64) (string, error) {
	payload := strings.Join(append([]string{purpose}, append(values, strconv.FormatInt(expiresAt, 10))...), "|")
	return secure.Sign(s.Key, payload), nil
}

// Verify validates a signed payload.
func (s HMACLinkSigner) Verify(_ context.Context, purpose string, values []string, expiresAt int64, signature string) error {
	now := time.Now().UTC().Unix()
	if s.Clock != nil {
		now = s.Clock.Now().Unix()
	}
	if now > expiresAt {
		return ErrTokenExpired
	}

	payload := strings.Join(append([]string{purpose}, append(values, strconv.FormatInt(expiresAt, 10))...), "|")
	if !secure.Verify(s.Key, payload, signature) {
		return ErrEmailVerificationInvalid
	}

	return nil
}
