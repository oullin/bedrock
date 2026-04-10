package consumer

import (
	"errors"
	"net/http"

	"github.com/bedrock/packages/auth"
)

// EmailVerificationRequest handles email verification logic for the current user.
type EmailVerificationRequest struct {
	user auth.Authenticatable
}

// NewEmailVerificationRequest creates a request for the given user.
func NewEmailVerificationRequest(user auth.Authenticatable) *EmailVerificationRequest {
	return &EmailVerificationRequest{user: user}
}

// Fulfill marks the user's email as verified if it has not been verified yet.
func (e *EmailVerificationRequest) Fulfill(r *http.Request) error {
	mv, ok := e.user.(auth.MustVerifyEmail)

	if !ok {
		return errors.New("user does not implement MustVerifyEmail")
	}

	if mv.HasVerifiedEmail() {
		return nil
	}

	if err := mv.MarkEmailAsVerified(); err != nil {
		return err
	}

	return nil
}

// HasVerifiedEmail reports whether the user's email is already verified.
func (e *EmailVerificationRequest) HasVerifiedEmail() bool {
	mv, ok := e.user.(auth.MustVerifyEmail)

	if !ok {
		return false
	}

	return mv.HasVerifiedEmail()
}
