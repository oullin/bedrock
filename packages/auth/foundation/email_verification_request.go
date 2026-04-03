package foundation

import (
	"context"
	"fmt"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/events"
	securitycrypto "github.com/gollin/packages/security/crypto"
)

// VerifiedDispatcher dispatches verification events.
type VerifiedDispatcher interface {
	DispatchVerified(ctx context.Context, event events.Verified) error
}

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

// Fulfill marks the user's email as verified and dispatches a Verified event once.
func (r EmailVerificationRequest) Fulfill(ctx context.Context, dispatcher VerifiedDispatcher) error {
	return r.FulfillAt(ctx, time.Now().UTC(), dispatcher)
}

// FulfillAt fulfils the request at a specific time.
func (r EmailVerificationRequest) FulfillAt(ctx context.Context, now time.Time, dispatcher VerifiedDispatcher) error {
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

	if dispatcher != nil {
		return dispatcher.DispatchVerified(ctx, events.Verified{
			User: r.User,
			At:   now,
		})
	}

	return nil
}
