package mailx_test

import (
	"context"
	"sync/atomic"
	"testing"

	cevents "github.com/bedrock/packages/contracts/events"
	cmail "github.com/bedrock/packages/contracts/mail"
	"github.com/bedrock/packages/mailx"
)

// For HTMLString, the HTML content is literal.

// Original should be preserved in X-Original-To header.

// --- Mock event dispatcher ---

type mockDispatcher struct {
	count atomic.Int64
}

func newTestMailer(t *testing.T, opts ...mailx.MailerOption) (*mailx.Mailer, *mailx.ArrayTransport) {
	t.Helper()

	transport := mailx.NewArrayTransport()
	mailer := mailx.NewMailer("test", transport, opts...)

	return mailer, transport
}

func TestMailerSendWithViewContent(t *testing.T) {
	t.Parallel()

	mailer, transport := newTestMailer(t)

	m := &mailx.Mailable{}
	m.SetFrom("sender@example.com").
		SetTo(cmail.Address{Email: "recipient@example.com"}).
		SetSubject("Hello").
		SetHTML("<h1>Welcome</h1>")

	_, err := mailer.Send(context.Background(), m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	messages := transport.Messages()

	if len(messages) != 1 {
		t.Fatalf("expected 1 sent message, got %d", len(messages))
	}

	if messages[0].Original.GetHTMLBody() != "<h1>Welcome</h1>" {
		t.Fatalf("expected HTML body, got %s", messages[0].Original.GetHTMLBody())
	}
}

func TestMailerSendWithCCAndBCC(t *testing.T) {
	t.Parallel()

	mailer, transport := newTestMailer(t)

	m := &mailx.Mailable{}
	m.SetFrom("sender@example.com").
		SetTo(cmail.Address{Email: "to@example.com"}).
		SetCC(cmail.Address{Email: "cc@example.com"}).
		SetBCC(cmail.Address{Email: "bcc@example.com"}).
		SetSubject("Test").
		SetHTML("<p>Body</p>")

	_, err := mailer.Send(context.Background(), m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := transport.Messages()[0].Original

	if len(msg.GetCC()) != 1 || msg.GetCC()[0].Email != "cc@example.com" {
		t.Fatal("expected CC recipient cc@example.com")
	}

	if len(msg.GetBCC()) != 1 || msg.GetBCC()[0].Email != "bcc@example.com" {
		t.Fatal("expected BCC recipient bcc@example.com")
	}
}

func TestMailerSendWithHTMLString(t *testing.T) {
	t.Parallel()

	mailer, transport := newTestMailer(t)

	m := &mailx.Mailable{}
	m.SetFrom("sender@example.com").
		SetTo(cmail.Address{Email: "to@example.com"}).
		SetSubject("HTML String").
		SetHTML("<strong>Bold</strong>")

	_, err := mailer.Send(context.Background(), m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := transport.Messages()[0].Original.GetHTMLBody()

	if body != "<strong>Bold</strong>" {
		t.Fatalf("expected literal HTML, got %s", body)
	}
}

func TestMailerSendWithCallableContent(t *testing.T) {
	t.Parallel()

	mailer, transport := newTestMailer(t)

	m := &mailx.Mailable{}
	m.SetFrom("sender@example.com").
		SetTo(cmail.Address{Email: "to@example.com"}).
		SetSubject("Callback").
		WithCallback(func(msg *cmail.Message) {
			msg.SetHTMLBody("<p>From callback</p>")
		})

	_, err := mailer.Send(context.Background(), m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := transport.Messages()[0].Original.GetHTMLBody()

	if body != "<p>From callback</p>" {
		t.Fatalf("expected callback HTML, got %s", body)
	}
}

func TestMailerSendWithHTMLMethod(t *testing.T) {
	t.Parallel()

	mailer, transport := newTestMailer(t)

	_, err := mailer.HTML(context.Background(), "<p>Direct HTML</p>", func(msg *cmail.Message) {
		msg.From("sender@example.com")
		msg.To([]cmail.Address{{Email: "to@example.com"}})
		msg.Subject("Direct")
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body := transport.Messages()[0].Original.GetHTMLBody()

	if body != "<p>Direct HTML</p>" {
		t.Fatalf("expected direct HTML, got %s", body)
	}
}

func TestMailerSendPlainText(t *testing.T) {
	t.Parallel()

	mailer, transport := newTestMailer(t)

	_, err := mailer.Raw(context.Background(), "Hello plain text", func(msg *cmail.Message) {
		msg.From("sender@example.com")
		msg.To([]cmail.Address{{Email: "to@example.com"}})
		msg.Subject("Plain")
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	text := transport.Messages()[0].Original.GetTextBody()

	if text != "Hello plain text" {
		t.Fatalf("expected plain text, got %s", text)
	}
}

func TestMailerSendExplicitViews(t *testing.T) {
	t.Parallel()

	mailer, transport := newTestMailer(t)

	m := &mailx.Mailable{}
	m.SetFrom("sender@example.com").
		SetTo(cmail.Address{Email: "to@example.com"}).
		SetSubject("Dual").
		SetHTML("<p>HTML part</p>").
		SetText("Text part")

	m.GetContent().HTMLString = true

	_, err := mailer.Send(context.Background(), m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := transport.Messages()[0].Original

	if msg.GetHTMLBody() != "<p>HTML part</p>" {
		t.Fatalf("expected HTML body, got %s", msg.GetHTMLBody())
	}

	if msg.GetTextBody() != "Text part" {
		t.Fatalf("expected text body, got %s", msg.GetTextBody())
	}
}

func TestMailerToWithEmailAndName(t *testing.T) {
	t.Parallel()

	mailer, transport := newTestMailer(t)

	m := &mailx.Mailable{}
	m.SetFrom("sender@example.com").
		SetTo(cmail.Address{Email: "john@example.com", Name: "John Doe"}).
		SetSubject("Named").
		SetHTML("<p>Hi</p>")

	_, err := mailer.Send(context.Background(), m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	to := transport.Messages()[0].Original.GetTo()

	if len(to) != 1 || to[0].Name != "John Doe" {
		t.Fatalf("expected recipient John Doe, got %v", to)
	}
}

func TestMailerGlobalFrom(t *testing.T) {
	t.Parallel()

	mailer, transport := newTestMailer(t)
	mailer.AlwaysFrom("global@example.com", "Global Sender")

	m := &mailx.Mailable{}
	m.SetTo(cmail.Address{Email: "to@example.com"}).
		SetSubject("Global From").
		SetHTML("<p>Body</p>")

	_, err := mailer.Send(context.Background(), m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	from := transport.Messages()[0].Original.GetFrom()

	if from.Email != "global@example.com" {
		t.Fatalf("expected global from, got %s", from.Email)
	}

	if from.Name != "Global Sender" {
		t.Fatalf("expected name Global Sender, got %s", from.Name)
	}
}

func TestMailerGlobalReplyTo(t *testing.T) {
	t.Parallel()

	mailer, transport := newTestMailer(t)
	mailer.AlwaysReplyTo("reply@example.com", "Support")

	m := &mailx.Mailable{}
	m.SetFrom("sender@example.com").
		SetTo(cmail.Address{Email: "to@example.com"}).
		SetSubject("Reply To").
		SetHTML("<p>Body</p>")

	_, err := mailer.Send(context.Background(), m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	replyTo := transport.Messages()[0].Original.GetReplyTo()

	if len(replyTo) == 0 || replyTo[0].Email != "reply@example.com" {
		t.Fatal("expected global reply-to")
	}
}

func TestMailerGlobalTo(t *testing.T) {
	t.Parallel()

	mailer, transport := newTestMailer(t)
	mailer.AlwaysTo(cmail.Address{Email: "always@example.com"})

	m := &mailx.Mailable{}
	m.SetFrom("sender@example.com").
		SetTo(cmail.Address{Email: "original@example.com"}).
		SetSubject("Always To").
		SetHTML("<p>Body</p>")

	_, err := mailer.Send(context.Background(), m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := transport.Messages()[0].Original
	to := msg.GetTo()

	if len(to) != 1 || to[0].Email != "always@example.com" {
		t.Fatal("expected always-to override")
	}

	headers := msg.GetCustomHeaders()

	if len(headers["X-Original-To"]) == 0 {
		t.Fatal("expected X-Original-To header with original recipient")
	}
}

func TestMailerGlobalReturnPath(t *testing.T) {
	t.Parallel()

	mailer, transport := newTestMailer(t)
	mailer.AlwaysReturnPath("bounce@example.com")

	m := &mailx.Mailable{}
	m.SetFrom("sender@example.com").
		SetTo(cmail.Address{Email: "to@example.com"}).
		SetSubject("Return Path").
		SetHTML("<p>Body</p>")

	_, err := mailer.Send(context.Background(), m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	rp := transport.Messages()[0].Original.GetReturnPath()

	if rp != "bounce@example.com" {
		t.Fatalf("expected return path bounce@example.com, got %s", rp)
	}
}

func TestMailerEventsDispatched(t *testing.T) {
	t.Parallel()

	dispatcher := &mockDispatcher{}
	mailer, _ := newTestMailer(t, mailx.WithDispatcher(dispatcher))

	m := &mailx.Mailable{}
	m.SetFrom("sender@example.com").
		SetTo(cmail.Address{Email: "to@example.com"}).
		SetSubject("Events").
		SetHTML("<p>Body</p>")

	_, err := mailer.Send(context.Background(), m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sendingCount := dispatcher.count.Load()

	if sendingCount != 2 {
		t.Fatalf("expected 2 event dispatches (Until + Dispatch), got %d", sendingCount)
	}
}

func TestMailerSendReturnsErrorWhenNoRecipients(t *testing.T) {
	t.Parallel()

	mailer, _ := newTestMailer(t)

	m := &mailx.Mailable{}
	m.SetFrom("sender@example.com").
		SetSubject("No Recipients").
		SetHTML("<p>Body</p>")

	_, err := mailer.Send(context.Background(), m)

	if err == nil {
		t.Fatal("expected error for no recipients")
	}
}

func TestMailerPendingMailSend(t *testing.T) {
	t.Parallel()

	mailer, transport := newTestMailer(t)

	m := &mailx.Mailable{}
	m.SetFrom("sender@example.com").
		SetSubject("Pending").
		SetHTML("<p>Body</p>")

	_, err := mailer.To(cmail.Address{Email: "pending@example.com"}).Send(context.Background(), m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	to := transport.Messages()[0].Original.GetTo()

	if len(to) != 1 || to[0].Email != "pending@example.com" {
		t.Fatal("expected pending mail to set To recipient")
	}
}

func TestMailerRender(t *testing.T) {
	t.Parallel()

	mailer, _ := newTestMailer(t)

	m := &mailx.Mailable{}
	m.SetHTML("<h1>Rendered</h1>")

	html, err := mailer.Render(context.Background(), m)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if html != "<h1>Rendered</h1>" {
		t.Fatalf("expected rendered HTML, got %s", html)
	}
}

func TestMailerFlushArrayTransport(t *testing.T) {
	t.Parallel()

	mailer, transport := newTestMailer(t)

	m := &mailx.Mailable{}
	m.SetFrom("sender@example.com").
		SetTo(cmail.Address{Email: "to@example.com"}).
		SetSubject("Flush").
		SetHTML("<p>Body</p>")

	mailer.Send(context.Background(), m)
	transport.Flush()

	if len(transport.Messages()) != 0 {
		t.Fatalf("expected 0 messages after flush, got %d", len(transport.Messages()))
	}
}

func (d *mockDispatcher) Listen(_ any, _ ...cevents.Listener) {}

func (d *mockDispatcher) HasListeners(_ any) bool { return false }

func (d *mockDispatcher) HasWildcardListeners(_ any) bool { return false }

func (d *mockDispatcher) Subscribe(_ cevents.Subscriber) {}

func (d *mockDispatcher) Until(_ context.Context, _ any) (any, error) {
	d.count.Add(1)

	return nil, nil
}

func (d *mockDispatcher) Dispatch(_ context.Context, _ any) ([]any, error) {
	d.count.Add(1)

	return nil, nil
}

func (d *mockDispatcher) Push(_ context.Context, _ any) {}

func (d *mockDispatcher) Flush(_ context.Context, _ string) error { return nil }

func (d *mockDispatcher) Forget(_ any) {}

func (d *mockDispatcher) ForgetPushed() {}

func (d *mockDispatcher) GetListeners(_ any) []cevents.Listener { return nil }
