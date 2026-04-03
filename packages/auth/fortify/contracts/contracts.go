package contracts

import (
	"context"
	"net/http"
	"time"

	auth "github.com/gollin/packages/auth"
)

// RegisterInput is the create-user action input.
type RegisterInput struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

// ForgotPasswordInput is the password reset link request input.
type ForgotPasswordInput struct {
	Email string `json:"email"`
}

// ResetPasswordInput is the new password submission input.
type ResetPasswordInput struct {
	Email                string `json:"email"`
	Token                string `json:"token"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

// ConfirmPasswordInput is the password confirmation input.
type ConfirmPasswordInput struct {
	Password string `json:"password"`
}

// TwoFactorChallengeInput is the two-factor login input.
type TwoFactorChallengeInput struct {
	Code         string `json:"code"`
	RecoveryCode string `json:"recovery_code"`
}

// UpdatePasswordInput is the user password update input.
type UpdatePasswordInput struct {
	CurrentPassword      string `json:"current_password"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

// UpdateProfileInformationInput is the profile update input.
type UpdateProfileInformationInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreatesNewUsers mirrors Fortify's create-user action contract.
type CreatesNewUsers interface {
	Create(ctx context.Context, input RegisterInput) (auth.Authenticatable, error)
}

// ResetsUserPasswords mirrors Fortify's reset-password action contract.
type ResetsUserPasswords interface {
	Reset(ctx context.Context, user auth.Authenticatable, password string) error
}

// UpdatesUserPasswords mirrors Fortify's password update action contract.
type UpdatesUserPasswords interface {
	Update(ctx context.Context, user auth.Authenticatable, input UpdatePasswordInput) error
}

// UpdatesUserProfileInformation mirrors Fortify's profile update action contract.
type UpdatesUserProfileInformation interface {
	Update(ctx context.Context, user auth.Authenticatable, input UpdateProfileInformationInput) error
}

// TwoFactorAuthenticationProvider mirrors Fortify's two-factor provider contract.
type TwoFactorAuthenticationProvider interface {
	GenerateSecret() (string, error)
	GenerateRecoveryCodes() ([]string, error)
	Validate(secret string, code string, now time.Time, allowedSkew int) bool
	OTPAuthURL(issuer string, account string, secret string) string
}

// RedirectsIfTwoFactorAuthenticatable mirrors Fortify's redirect contract.
type RedirectsIfTwoFactorAuthenticatable interface {
	Redirect(user auth.Authenticatable) string
}

// LoginViewResponse renders the login page response.
type LoginViewResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// RegisterViewResponse renders the registration page response.
type RegisterViewResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// RequestPasswordResetLinkViewResponse renders the forgot-password view response.
type RequestPasswordResetLinkViewResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// ResetPasswordViewResponse renders the reset-password view response.
type ResetPasswordViewResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// VerifyEmailViewResponse renders the verify-email view response.
type VerifyEmailViewResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// ConfirmPasswordViewResponse renders the confirm-password view response.
type ConfirmPasswordViewResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// TwoFactorChallengeViewResponse renders the two-factor challenge view response.
type TwoFactorChallengeViewResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// RegisterResponse renders a successful registration response.
type RegisterResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// LoginResponse renders a successful login response.
type LoginResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// LogoutResponse renders a successful logout response.
type LogoutResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// PasswordResetResponse renders a successful password reset response.
type PasswordResetResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// SuccessfulPasswordResetLinkRequestResponse renders a successful reset-link response.
type SuccessfulPasswordResetLinkRequestResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// FailedPasswordResetLinkRequestResponse renders a failed reset-link response.
type FailedPasswordResetLinkRequestResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// FailedPasswordResetResponse renders a failed password reset response.
type FailedPasswordResetResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// FailedPasswordConfirmationResponse renders a failed password confirmation response.
type FailedPasswordConfirmationResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// PasswordConfirmedResponse renders a successful password confirmation response.
type PasswordConfirmedResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// VerifyEmailResponse renders a successful verification response.
type VerifyEmailResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// EmailVerificationNotificationSentResponse renders a successful verification-notification response.
type EmailVerificationNotificationSentResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// TwoFactorEnabledResponse renders a successful enable-two-factor response.
type TwoFactorEnabledResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// TwoFactorConfirmedResponse renders a successful confirm-two-factor response.
type TwoFactorConfirmedResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// TwoFactorDisabledResponse renders a successful disable-two-factor response.
type TwoFactorDisabledResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// RecoveryCodesGeneratedResponse renders the recovery-codes response.
type RecoveryCodesGeneratedResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// TwoFactorLoginResponse renders a successful or pending two-factor login response.
type TwoFactorLoginResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// FailedTwoFactorLoginResponse renders a failed two-factor login response.
type FailedTwoFactorLoginResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// ProfileInformationUpdatedResponse renders a profile-update response.
type ProfileInformationUpdatedResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// PasswordUpdateResponse renders a password-update response.
type PasswordUpdateResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}

// LockoutResponse renders a throttle response.
type LockoutResponse interface {
	ToResponse(w http.ResponseWriter, r *http.Request, payload any)
}
