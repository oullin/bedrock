package mail

import (
	"context"
	"time"
)

// Mailer sends email messages.
type Mailer interface {
	// Raw sends a plain-text email, optionally configuring the message via
	// callback functions.
	Raw(ctx context.Context, text string, callbacks ...func(*Message)) (*SentMessage, error)
	// Send sends the given mailable.
	Send(ctx context.Context, mailable Mailable) (*SentMessage, error)
	// SendNow sends the mailable immediately, bypassing any queue.
	SendNow(ctx context.Context, mailable Mailable) (*SentMessage, error)
}

// Factory resolves named mailer instances.
type Factory interface {
	Mailer(name ...string) (Mailer, error)
}

// MailQueue adds queue support to a mailer.
type MailQueue interface {
	Queue(ctx context.Context, mailable Mailable) error
	Later(ctx context.Context, delay time.Duration, mailable Mailable) error
}
