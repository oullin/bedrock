package listeners

import (
	"context"

	"github.com/bedrock/packages/auth/events"
	contracts "github.com/bedrock/packages/contracts/auth"
)

// SendEmailVerificationNotification sends an email verification notification
// when a new user registers, if the user implements MustVerifyEmail and has
// not yet verified their email.
type SendEmailVerificationNotification struct{}

// Handle processes a Registered event.
func (l *SendEmailVerificationNotification) Handle(_ context.Context, event events.Registered) {
	mv, ok := event.User.(contracts.MustVerifyEmail)

	if !ok {
		return
	}

	if mv.HasVerifiedEmail() {
		return
	}

	mv.SendEmailVerificationNotification()
}
