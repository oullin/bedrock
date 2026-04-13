package mailx

import (
	"context"
	"time"

	cmail "github.com/bedrock/packages/contracts/mail"
)

// PendingMail is a fluent builder created by Mailer.To, Mailer.CC, or
// Mailer.BCC. It collects recipients and locale before delegating to the
// underlying Mailer for sending.
type PendingMail struct {
	mailer *Mailer
	to     []cmail.Address
	cc     []cmail.Address
	bcc    []cmail.Address
	locale string
}

// To sets the primary recipients.
func (p *PendingMail) To(addresses ...cmail.Address) *PendingMail {
	p.to = append(p.to, addresses...)

	return p
}

// CC sets the CC recipients.
func (p *PendingMail) CC(addresses ...cmail.Address) *PendingMail {
	p.cc = append(p.cc, addresses...)

	return p
}

// BCC sets the BCC recipients.
func (p *PendingMail) BCC(addresses ...cmail.Address) *PendingMail {
	p.bcc = append(p.bcc, addresses...)

	return p
}

// Locale sets the locale for template rendering.
func (p *PendingMail) Locale(locale string) *PendingMail {
	p.locale = locale

	return p
}

// Send sends the mailable after merging the pending recipients.
func (p *PendingMail) Send(ctx context.Context, mailable cmail.Mailable) (*cmail.SentMessage, error) {
	return p.mailer.Send(ctx, p.fill(mailable))
}

// SendNow sends the mailable immediately, bypassing any queue.
func (p *PendingMail) SendNow(ctx context.Context, mailable cmail.Mailable) (*cmail.SentMessage, error) {
	return p.mailer.SendNow(ctx, p.fill(mailable))
}

// Queue queues the mailable for later delivery.
func (p *PendingMail) Queue(ctx context.Context, mailable cmail.Mailable) error {
	return p.mailer.Queue(ctx, p.fill(mailable))
}

// Later queues the mailable with a delay.
func (p *PendingMail) Later(ctx context.Context, delay time.Duration, mailable cmail.Mailable) error {
	return p.mailer.Later(ctx, delay, p.fill(mailable))
}

func (p *PendingMail) fill(mailable cmail.Mailable) cmail.Mailable {
	env := mailable.GetEnvelope()

	if len(p.to) > 0 {
		env.To = append(env.To, p.to...)
	}

	if len(p.cc) > 0 {
		env.CC = append(env.CC, p.cc...)
	}

	if len(p.bcc) > 0 {
		env.BCC = append(env.BCC, p.bcc...)
	}

	return mailable
}
