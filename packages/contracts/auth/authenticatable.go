package auth

import "time"

// Authenticatable is any entity that can be authenticated.
type Authenticatable interface {
	GetAuthIdentifierName() string
	GetAuthIdentifier() string
	GetAuthPasswordName() string
	GetAuthPassword() string
	SetAuthPassword(password string)
	GetRememberToken() string
	SetRememberToken(token string)
	GetRememberTokenName() string
}

// BroadcastingAuthenticatable is implemented by users that expose a stable
// authentication identifier for private broadcast channel authorization.
type BroadcastingAuthenticatable interface {
	Authenticatable
	GetAuthIdentifierForBroadcasting() string
}

// MustVerifyEmail is implemented by users that require email verification.
type MustVerifyEmail interface {
	HasVerifiedEmail() bool
	MarkEmailAsVerified(at time.Time)
	MarkEmailAsUnverified()
	GetEmailForVerification() string
}

// EmailVerificationNotificationSender is implemented by users that can send
// their own email verification notification.
type EmailVerificationNotificationSender interface {
	MustVerifyEmail
	SendEmailVerificationNotification()
}

// CanResetPassword is implemented by users that support password resets.
type CanResetPassword interface {
	GetEmailForPasswordReset() string
}

// PasswordResetNotificationSender is implemented by users that can send their
// own password reset notification.
type PasswordResetNotificationSender interface {
	CanResetPassword
	SendPasswordResetNotification(token string)
}

// TwoFactorAuthenticatable is implemented by users that support 2FA.
type TwoFactorAuthenticatable interface {
	IsTwoFactorEnabled() bool
	SetTwoFactorEnabled(enabled bool)
	GetTwoFactorSecret() string
	SetTwoFactorSecret(secret string)
	GetTwoFactorRecoveryCodes() []string
	SetTwoFactorRecoveryCodes(codes []string)
	GetTwoFactorConfirmedAt() *time.Time
	SetTwoFactorConfirmedAt(at *time.Time)
}
