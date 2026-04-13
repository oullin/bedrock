package notifications_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/notifications"
)

func TestMailChannelSend(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mailer := &mockMailer{}
	factory := &mockMailerFactory{mailer: mailer}
	channel := notifications.NewMailChannel(factory)
	notifiable := newMockNotifiable("user-1")
	notifiable.routes["mail"] = "user@example.com"
	notification := newTestMailNotification()

	err := channel.Send(ctx, notifiable, notification)

	if err != nil {
		t.Fatalf("Send error: %v", err)
	}

	if mailer.CallCount() != 1 {
		t.Fatalf("expected 1 mailer.Send call, got %d", mailer.CallCount())
	}
}

func TestMailChannelSendSetsRecipient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mailer := &mockMailer{}
	factory := &mockMailerFactory{mailer: mailer}
	channel := notifications.NewMailChannel(factory)
	notifiable := newMockNotifiable("user-2")
	notifiable.routes["mail"] = "test@example.com"
	notification := newTestMailNotification()

	_ = channel.Send(ctx, notifiable, notification)

	mailable := mailer.calls[0].Mailable
	envelope := mailable.GetEnvelope()

	if len(envelope.To) != 1 {
		t.Fatalf("To count = %d, want 1", len(envelope.To))
	}

	if envelope.To[0].Email != "test@example.com" {
		t.Fatalf("To email = %q, want %q", envelope.To[0].Email, "test@example.com")
	}
}

func TestMailChannelSendSetsSubject(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mailer := &mockMailer{}
	factory := &mockMailerFactory{mailer: mailer}
	channel := notifications.NewMailChannel(factory)
	notifiable := newMockNotifiable("user-3")
	notifiable.routes["mail"] = "test@example.com"
	notification := newTestMailNotification()

	_ = channel.Send(ctx, notifiable, notification)

	mailable := mailer.calls[0].Mailable
	envelope := mailable.GetEnvelope()

	if envelope.Subject != "Test Subject" {
		t.Fatalf("subject = %q, want %q", envelope.Subject, "Test Subject")
	}
}

func TestMailChannelSendMissingMailNotificationError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mailer := &mockMailer{}
	factory := &mockMailerFactory{mailer: mailer}
	channel := notifications.NewMailChannel(factory)
	notifiable := newMockNotifiable("user-4")
	notification := newTestNotification("mail")

	err := channel.Send(ctx, notifiable, notification)

	if err == nil {
		t.Fatal("expected error for missing MailNotification")
	}
}

func TestMailChannelSendWithStringSliceRecipients(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mailer := &mockMailer{}
	factory := &mockMailerFactory{mailer: mailer}
	channel := notifications.NewMailChannel(factory)
	notifiable := newMockNotifiable("user-5")
	notifiable.routes["mail"] = []string{"a@example.com", "b@example.com"}
	notification := newTestMailNotification()

	_ = channel.Send(ctx, notifiable, notification)

	mailable := mailer.calls[0].Mailable
	envelope := mailable.GetEnvelope()

	if len(envelope.To) != 2 {
		t.Fatalf("To count = %d, want 2", len(envelope.To))
	}
}
