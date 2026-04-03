package auth

import "time"

// User is the shared auth user model stored by pluggable user providers.
type User struct {
	ID                string     `json:"id"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"`
	EmailVerifiedAt   *time.Time `json:"emailVerifiedAt,omitempty"`
	RememberTokenHash string     `json:"-"`
	TwoFactorEnabled  bool       `json:"twoFactorEnabled"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

// Session represents an authenticated or pending two-factor session.
type Session struct {
	ID                  string     `json:"id"`
	UserID              string     `json:"userId"`
	PendingTwoFactor    bool       `json:"pendingTwoFactor"`
	PendingRemember     bool       `json:"pendingRemember"`
	PasswordConfirmedAt *time.Time `json:"passwordConfirmedAt,omitempty"`
	AuthenticatedAt     *time.Time `json:"authenticatedAt,omitempty"`
	LastSeenAt          time.Time  `json:"lastSeenAt"`
	CreatedAt           time.Time  `json:"createdAt"`
	ExpiresAt           time.Time  `json:"expiresAt"`
}

// PasswordResetToken stores a reset token hash and expiry.
type PasswordResetToken struct {
	UserID    string
	TokenHash string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// TwoFactorState stores a TOTP secret and recovery codes for a user.
type TwoFactorState struct {
	UserID        string
	Secret        string
	RecoveryCodes []string
	EnabledAt     time.Time
	UpdatedAt     time.Time
}

// MailMessage is dispatched by the configured mailer.
type MailMessage struct {
	To       string
	Subject  string
	Body     string
	Metadata map[string]string
}

// RegisterInput is the AuthFlows-style registration payload.
type RegisterInput struct {
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

// LoginInput is the AuthFlows-style login payload.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

// ForgotPasswordInput requests a password reset link.
type ForgotPasswordInput struct {
	Email string `json:"email"`
}

// ResetPasswordInput resets a password with a broker token.
type ResetPasswordInput struct {
	Email                string `json:"email"`
	Token                string `json:"token"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

// ConfirmPasswordInput confirms the current password for sensitive actions.
type ConfirmPasswordInput struct {
	Password string `json:"password"`
}

// TwoFactorChallengeInput submits either a TOTP code or a recovery code.
type TwoFactorChallengeInput struct {
	Code         string `json:"code"`
	RecoveryCode string `json:"recovery_code"`
}

// LoginResult reports the outcome of an attempted login or two-factor challenge.
type LoginResult struct {
	User               *User
	Session            *Session
	RememberToken      string
	RequiresTwoFactor  bool
	PasswordConfirmed  bool
	AuthenticatedState string
}

// EnableTwoFactorResult contains the initial two-factor secret and recovery codes.
type EnableTwoFactorResult struct {
	Secret        string   `json:"secret"`
	RecoveryCodes []string `json:"recoveryCodes"`
	OTPAuthURL    string   `json:"otpAuthUrl"`
}
