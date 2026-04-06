package consumer

import (
	"context"
	"fmt"
	"time"

	auth "github.com/bedrock/packages/auth"
	securitycrypto "github.com/bedrock/packages/support/crypto"
)

// EmailVerificationRequest ports Foundation/Auth email verification semantics.
type EmailVerificationRequest struct {
	User      auth.Authenticatable
	RouteID   string
	RouteHash string
}

// Authorize validates the route id and email hash against the current user.
func (r EmailVerificationRequest) Authorize() bool {
	if r.User == nil {
		return false
	}

	verifiable, ok := r.User.(auth.MustVerifyEmail)
	if !ok {
		return false
	}

	return r.RouteID == r.User.GetAuthIdentifier() && r.RouteHash == securitycrypto.EmailHash(verifiable.GetEmailForVerification())
}

// Fulfill marks the user's email as verified.
func (r EmailVerificationRequest) Fulfill(ctx context.Context) error {
	return r.FulfillAt(ctx, time.Now().UTC())
}

// FulfillAt fulfils the request at a specific time.
func (r EmailVerificationRequest) FulfillAt(_ context.Context, now time.Time) error {
	if !r.Authorize() {
		return auth.ErrUnauthorized
	}

	verifiable, ok := r.User.(auth.MustVerifyEmail)
	if !ok {
		return fmt.Errorf("auth: user does not support email verification")
	}

	if verifiable.HasVerifiedEmail() {
		return nil
	}

	verifiable.MarkEmailAsVerified(now)

	return nil
}
