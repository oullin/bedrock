package consumer

import (
	"errors"
	"net/http"
	"time"

	cauth "github.com/bedrock/packages/contracts/auth"
)

// EmailVerificationRequest handles email verification logic for the current user.
type EmailVerificationRequest struct {
	user cauth.Authenticatable
}

// NewEmailVerificationRequest creates a request for the given user.
func NewEmailVerificationRequest(user cauth.Authenticatable) *EmailVerificationRequest {
	return &EmailVerificationRequest{user: user}
}

// Fulfill marks the user's email as verified if it has not been verified yet.
func (e *EmailVerificationRequest) Fulfill(_ *http.Request) error {
	mv, ok := e.user.(cauth.MustVerifyEmail)

	if !ok {
		return errors.New("user does not implement MustVerifyEmail")
	}

	if mv.HasVerifiedEmail() {
		return nil
	}

	mv.MarkEmailAsVerified(time.Now())

	return nil
}

// HasVerifiedEmail reports whether the user's email is already verified.
func (e *EmailVerificationRequest) HasVerifiedEmail() bool {
	mv, ok := e.user.(cauth.MustVerifyEmail)

	if !ok {
		return false
	}

	return mv.HasVerifiedEmail()
}
