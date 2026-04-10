package auth

// Authenticatable is any entity that can be authenticated.
type Authenticatable interface {
	GetAuthIdentifierName() string
	GetAuthIdentifier() any
	GetAuthPassword() string
	GetRememberToken() string
	SetRememberToken(token string)
	GetRememberTokenName() string
}

// MustVerifyEmail is implemented by users that require email verification.
type MustVerifyEmail interface {
	HasVerifiedEmail() bool
	MarkEmailAsVerified() error
	SendEmailVerificationNotification()
	GetEmailForVerification() string
}

// CanResetPassword is implemented by users that support password resets.
type CanResetPassword interface {
	GetEmailForPasswordReset() string
	SendPasswordResetNotification(token string)
}
