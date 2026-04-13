package notifications

import (
	"bytes"
	"context"
	"fmt"
	"io"

	cn "github.com/bedrock/packages/contracts/notifications"

	"github.com/bedrock/packages/contracts/mail"
)

// MailChannel delivers notifications via email using the mail.Factory.
type MailChannel struct {
	mailerFactory mail.Factory
}

// NewMailChannel creates a MailChannel with the given mailer factory.

// compile-time interface check.

// Send delivers the notification as an email to the notifiable entity.

// Set notification ID metadata.

// notificationMailable adapts a MailMessage to the mail.Mailable interface.
type notificationMailable struct {
	message        *MailMessage
	notifiable     cn.Notifiable
	notificationID string
}

func NewMailChannel(factory mail.Factory) *MailChannel {
	return &MailChannel{mailerFactory: factory}
}

var _ cn.Channel = (*MailChannel)(nil)

func (c *MailChannel) Send(ctx context.Context, notifiable cn.Notifiable, notification any) error {
	mn, ok := notification.(MailNotification)

	if !ok {
		return fmt.Errorf("%w: %T", ErrMissingMailData, notification)
	}

	message, err := mn.ToMail(ctx, notifiable)

	if err != nil {
		return fmt.Errorf("notifications: build mail message: %w", err)
	}

	mailer, err := c.resolveMailer(message)

	if err != nil {
		return fmt.Errorf("notifications: resolve mailer: %w", err)
	}

	mailable := c.buildMailable(ctx, notifiable, notification, message)

	_, err = mailer.Send(ctx, mailable)

	return err
}

func (c *MailChannel) resolveMailer(message *MailMessage) (mail.Mailer, error) {
	if message.GetMailer() != "" {
		return c.mailerFactory.Mailer(message.GetMailer())
	}

	return c.mailerFactory.Mailer()
}

func (c *MailChannel) buildMailable(_ context.Context, notifiable cn.Notifiable, notification any, message *MailMessage) *notificationMailable {
	m := &notificationMailable{
		message:    message,
		notifiable: notifiable,
	}

	if n, ok := notification.(interface{ GetID() string }); ok {
		m.notificationID = n.GetID()
	}

	return m
}

func (m *notificationMailable) GetEnvelope() *mail.Envelope {
	env := &mail.Envelope{
		Subject:  m.message.GetSubject(),
		CC:       m.message.GetCC(),
		BCC:      m.message.GetBCC(),
		ReplyTo:  m.message.GetReplyTo(),
		Tags:     m.message.GetTags(),
		Metadata: m.message.GetMetadata(),
		Using:    m.message.GetCallbacks(),
	}

	if m.message.GetFrom() != nil {
		env.From = *m.message.GetFrom()
	}

	// Resolve recipient from notifiable routing.
	to := m.notifiable.RouteNotificationFor(context.Background(), "mail")

	switch v := to.(type) {
	case string:
		env.To = []mail.Address{{Email: v}}
	case mail.Address:
		env.To = []mail.Address{v}
	case []mail.Address:
		env.To = v
	case []string:
		addrs := make([]mail.Address, len(v))

		for i, email := range v {
			addrs[i] = mail.Address{Email: email}
		}

		env.To = addrs
	}

	return env
}

func (m *notificationMailable) GetContent() *mail.Content {
	content := &mail.Content{
		With: m.message.Data(),
	}

	if m.message.GetMarkdown() != "" {
		content.Markdown = m.message.GetMarkdown()
	} else if m.message.GetView() != "" {
		content.HTML = m.message.GetView()
	}

	return content
}

func (m *notificationMailable) GetAttachments() []*mail.Attachment {
	attachments := make([]*mail.Attachment, 0, len(m.message.GetAttachments())+len(m.message.GetRawAttachments()))
	attachments = append(attachments, m.message.GetAttachments()...)

	for _, raw := range m.message.GetRawAttachments() {
		data := raw.Data
		a := mail.FromData(func() (io.Reader, error) {
			return bytes.NewReader(data), nil
		}, raw.Name)

		if raw.Mime != "" {
			a.WithMime(raw.Mime)
		}

		attachments = append(attachments, a)
	}

	return attachments
}

func (m *notificationMailable) GetHeaders() *mail.Headers {
	return nil
}
