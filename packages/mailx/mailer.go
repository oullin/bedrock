package mailx

import (
	"context"
	"fmt"
	"sync"
	"time"

	cevents "github.com/bedrock/packages/contracts/events"
	cmail "github.com/bedrock/packages/contracts/mail"
)

// Mailer is the primary email sending engine. It builds messages from
// Mailable instances, applies global overrides, fires events, and
// delegates to a Transport for delivery.
type Mailer struct {
	name      string
	transport Transport

	mu               sync.RWMutex
	dispatcher       cevents.Dispatcher
	alwaysFrom       *cmail.Address
	alwaysReplyTo    *cmail.Address
	alwaysReturnPath string
	alwaysTo         []cmail.Address
}

// MailerOption configures a Mailer.
type MailerOption func(*Mailer)

// WithDispatcher sets the event dispatcher.
func WithDispatcher(d cevents.Dispatcher) MailerOption {
	return func(m *Mailer) {
		m.dispatcher = d
	}
}

// NewMailer creates a Mailer with the given transport.
func NewMailer(name string, transport Transport, opts ...MailerOption) *Mailer {
	m := &Mailer{
		name:      name,
		transport: transport,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// AlwaysFrom sets a global sender address applied to all messages.
func (m *Mailer) AlwaysFrom(address string, name ...string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	a := cmail.Address{Email: address}

	if len(name) > 0 {
		a.Name = name[0]
	}

	m.alwaysFrom = &a
}

// AlwaysReplyTo sets a global reply-to address applied to all messages.
func (m *Mailer) AlwaysReplyTo(address string, name ...string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	a := cmail.Address{Email: address}

	if len(name) > 0 {
		a.Name = name[0]
	}

	m.alwaysReplyTo = &a
}

// AlwaysReturnPath sets a global return-path applied to all messages.
func (m *Mailer) AlwaysReturnPath(address string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.alwaysReturnPath = address
}

// AlwaysTo sets global recipients that override per-message To addresses.
func (m *Mailer) AlwaysTo(addresses ...cmail.Address) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.alwaysTo = addresses
}

// To begins a fluent mail builder addressed to the given recipients.
func (m *Mailer) To(addresses ...cmail.Address) *PendingMail {
	return &PendingMail{mailer: m, to: addresses}
}

// CC begins a fluent mail builder with CC recipients.
func (m *Mailer) CC(addresses ...cmail.Address) *PendingMail {
	return &PendingMail{mailer: m, cc: addresses}
}

// BCC begins a fluent mail builder with BCC recipients.
func (m *Mailer) BCC(addresses ...cmail.Address) *PendingMail {
	return &PendingMail{mailer: m, bcc: addresses}
}

// Raw sends a plain-text email, optionally configuring the message via
// callbacks.
func (m *Mailer) Raw(ctx context.Context, text string, callbacks ...func(*cmail.Message)) (*cmail.SentMessage, error) {
	msg := cmail.NewMessage()
	msg.SetTextBody(text)

	for _, cb := range callbacks {
		cb(msg)
	}

	return m.sendMessage(ctx, msg)
}

// HTML sends a raw HTML email.
func (m *Mailer) HTML(ctx context.Context, html string, callbacks ...func(*cmail.Message)) (*cmail.SentMessage, error) {
	msg := cmail.NewMessage()
	msg.SetHTMLBody(html)

	for _, cb := range callbacks {
		cb(msg)
	}

	return m.sendMessage(ctx, msg)
}

// Plain sends a plain-text email with the given content.
func (m *Mailer) Plain(ctx context.Context, text string, callbacks ...func(*cmail.Message)) (*cmail.SentMessage, error) {
	return m.Raw(ctx, text, callbacks...)
}

// Send builds a Message from the Mailable and sends it.
func (m *Mailer) Send(ctx context.Context, mailable cmail.Mailable) (*cmail.SentMessage, error) {
	return m.sendMailable(ctx, mailable)
}

// SendNow sends the mailable immediately, bypassing any queue.
func (m *Mailer) SendNow(ctx context.Context, mailable cmail.Mailable) (*cmail.SentMessage, error) {
	return m.sendMailable(ctx, mailable)
}

// Render builds the message from a Mailable without sending and returns
// the HTML body.
func (m *Mailer) Render(_ context.Context, mailable cmail.Mailable) (string, error) {
	msg := m.buildMessage(mailable)

	return msg.GetHTMLBody(), nil
}

// Queue queues the mailable for later delivery.
func (m *Mailer) Queue(_ context.Context, _ cmail.Mailable) error {
	return ErrNoQueue
}

// Later queues the mailable with a delay.
func (m *Mailer) Later(_ context.Context, _ time.Duration, _ cmail.Mailable) error {
	return ErrNoQueue
}

// GetTransport returns the underlying transport.
func (m *Mailer) GetTransport() Transport {
	return m.transport
}

// SetTransport replaces the underlying transport.
func (m *Mailer) SetTransport(t Transport) {
	m.transport = t
}

// Name returns the mailer name.
func (m *Mailer) Name() string {
	return m.name
}

func (m *Mailer) sendMailable(ctx context.Context, mailable cmail.Mailable) (*cmail.SentMessage, error) {
	msg := m.buildMessage(mailable)

	// Apply mailable callbacks.
	if mb, ok := mailable.(*Mailable); ok {
		for _, cb := range mb.GetCallbacks() {
			cb(msg)
		}
	}

	return m.sendMessage(ctx, msg)
}

func (m *Mailer) buildMessage(mailable cmail.Mailable) *cmail.Message {
	env := mailable.GetEnvelope()
	content := mailable.GetContent()
	headers := mailable.GetHeaders()
	attachments := mailable.GetAttachments()

	msg := cmail.NewMessage()

	// Envelope.
	if env.From.Email != "" {
		msg.From(env.From.Email, env.From.Name)
	}

	if len(env.To) > 0 {
		msg.To(env.To)
	}

	if len(env.CC) > 0 {
		msg.CC(env.CC)
	}

	if len(env.BCC) > 0 {
		msg.BCC(env.BCC)
	}

	if len(env.ReplyTo) > 0 {
		msg.ReplyTo(env.ReplyTo)
	}

	if env.Subject != "" {
		msg.Subject(env.Subject)
	}

	// Tags and metadata as custom headers.
	for _, tag := range env.Tags {
		msg.SetHeader("X-Tag", tag)
	}

	for k, v := range env.Metadata {
		msg.SetHeader("X-Metadata-"+k, v)
	}

	// Content.
	if content != nil {
		if content.HTMLString {
			msg.SetHTMLBody(content.HTML)
		} else if content.HTML != "" {
			msg.SetHTMLBody(m.renderTemplate(content.HTML, content.With))
		}

		if content.Text != "" {
			msg.SetTextBody(m.renderTemplate(content.Text, content.With))
		}

		if content.Markdown != "" {
			msg.SetHTMLBody(content.Markdown)
		}
	}

	// Headers.
	if headers != nil {
		if headers.MessageID != "" {
			msg.SetHeader("Message-ID", headers.MessageID)
		}

		if len(headers.References) > 0 {
			msg.SetHeader("References", headers.ReferencesString())
		}

		for _, h := range headers.Text {
			msg.SetHeader(h.Name, h.Value)
		}
	}

	// Attachments.
	for _, a := range attachments {
		msg.AttachWith(a)
	}

	// Envelope.Using callbacks.
	for _, cb := range env.Using {
		cb(msg)
	}

	return msg
}

func (m *Mailer) sendMessage(ctx context.Context, msg *cmail.Message) (*cmail.SentMessage, error) {
	m.applyGlobalOverrides(msg)

	recipients := msg.AllRecipients()

	if len(recipients) == 0 {
		return nil, ErrNoRecipients
	}

	// Dispatch MessageSending event.
	m.mu.RLock()
	dispatcher := m.dispatcher
	m.mu.RUnlock()

	if dispatcher != nil {
		event := MessageSending{
			Mailer:  m.name,
			Message: msg,
		}

		_, err := dispatcher.Until(ctx, event)

		if err != nil {
			return nil, fmt.Errorf("mail: message sending cancelled: %w", err)
		}
	}

	// Send via transport.
	sent, err := m.transport.Send(ctx, msg)

	if err != nil {
		return nil, err
	}

	// Dispatch MessageSent event.
	if dispatcher != nil {
		event := MessageSent{
			Mailer:      m.name,
			SentMessage: sent,
		}

		dispatcher.Dispatch(ctx, event)
	}

	return sent, nil
}

func (m *Mailer) applyGlobalOverrides(msg *cmail.Message) {
	m.mu.RLock()

	defer m.mu.RUnlock()

	if m.alwaysFrom != nil && msg.GetFrom().Email == "" {
		msg.From(m.alwaysFrom.Email, m.alwaysFrom.Name)
	}

	if m.alwaysReplyTo != nil {
		msg.ReplyTo([]cmail.Address{*m.alwaysReplyTo})
	}

	if m.alwaysReturnPath != "" {
		msg.ReturnPath(m.alwaysReturnPath)
	}

	if len(m.alwaysTo) > 0 {
		original := msg.GetTo()

		msg.To(m.alwaysTo, true)

		if len(original) > 0 {
			for _, a := range original {
				msg.SetHeader("X-Original-To", a.String())
			}
		}
	}
}

func (m *Mailer) renderTemplate(tmpl string, data map[string]any) string {
	// For now, return the template string as-is. A full template
	// engine integration (html/template or text/template) can be
	// added when view rendering support is needed.
	return tmpl
}
