package listeners

import (
	"context"

	"github.com/bedrock/packages/auth/events"
	cauth "github.com/bedrock/packages/contracts/auth"
)

// EmailVerificationSender is kept for backward compatibility. New code should
// use contracts/auth.EmailVerificationNotificationSender.
type EmailVerificationSender = cauth.EmailVerificationNotificationSender

// SendEmailVerificationNotification sends an email verification notification
// when a new user registers, if the user implements EmailVerificationSender and has
// not yet verified their email.
type SendEmailVerificationNotification struct{}

// Handle processes a Registered event.
func (l *SendEmailVerificationNotification) Handle(_ context.Context, event events.Registered) {
	mv, ok := event.User.(cauth.EmailVerificationNotificationSender)

	if !ok {
		return
	}

	if mv.HasVerifiedEmail() {
		return
	}

	mv.SendEmailVerificationNotification()
}
