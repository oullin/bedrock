package notifications_test

import (
	"testing"

	"github.com/bedrock/packages/contracts/mail"
	"github.com/bedrock/packages/notifications"
)

func TestMailMessageFluentChaining(t *testing.T) {
	t.Parallel()

	m := notifications.NewMailMessage().
		Subject("Welcome").
		Greeting("Hello!").
		Line("Welcome aboard.").
		Action("Get Started", "https://example.com/start").
		Line("Thanks for joining.").
		Salutation("Best")

	if m.GetSubject() != "Welcome" {
		t.Fatalf("subject = %q, want %q", m.GetSubject(), "Welcome")
	}

	if m.GetGreeting() != "Hello!" {
		t.Fatalf("greeting = %q, want %q", m.GetGreeting(), "Hello!")
	}

	if len(m.GetIntroLines()) != 1 {
		t.Fatalf("intro lines = %d, want 1", len(m.GetIntroLines()))
	}

	if len(m.GetOutroLines()) != 1 {
		t.Fatalf("outro lines = %d, want 1", len(m.GetOutroLines()))
	}

	if m.GetSalutation() != "Best" {
		t.Fatalf("salutation = %q, want %q", m.GetSalutation(), "Best")
	}
}

func TestMailMessageView(t *testing.T) {
	t.Parallel()

	m := notifications.NewMailMessage().
		View("emails.welcome", map[string]any{"name": "John"})

	if m.GetView() != "emails.welcome" {
		t.Fatalf("view = %q, want %q", m.GetView(), "emails.welcome")
	}

	if m.GetViewData()["name"] != "John" {
		t.Fatalf("view data name = %v, want %q", m.GetViewData()["name"], "John")
	}
}

func TestMailMessageMarkdown(t *testing.T) {
	t.Parallel()

	m := notifications.NewMailMessage().Markdown("markdown.welcome")

	if m.GetMarkdown() != "markdown.welcome" {
		t.Fatalf("markdown = %q, want %q", m.GetMarkdown(), "markdown.welcome")
	}
}

func TestMailMessageTheme(t *testing.T) {
	t.Parallel()

	m := notifications.NewMailMessage().Theme("dark")

	if m.GetTheme() != "dark" {
		t.Fatalf("theme = %q, want %q", m.GetTheme(), "dark")
	}
}

func TestMailMessageFrom(t *testing.T) {
	t.Parallel()

	m := notifications.NewMailMessage().From("noreply@example.com", "App")

	from := m.GetFrom()

	if from == nil {
		t.Fatal("expected non-nil from address")
	}

	if from.Email != "noreply@example.com" {
		t.Fatalf("from email = %q, want %q", from.Email, "noreply@example.com")
	}

	if from.Name != "App" {
		t.Fatalf("from name = %q, want %q", from.Name, "App")
	}
}

func TestMailMessageReplyTo(t *testing.T) {
	t.Parallel()

	m := notifications.NewMailMessage().
		ReplyTo("reply@example.com").
		ReplyTo("other@example.com", "Other")

	if len(m.GetReplyTo()) != 2 {
		t.Fatalf("replyTo count = %d, want 2", len(m.GetReplyTo()))
	}
}

func TestMailMessageCCAndBCC(t *testing.T) {
	t.Parallel()

	m := notifications.NewMailMessage().
		CC("cc@example.com").
		BCC("bcc@example.com", "Secret")

	if len(m.GetCC()) != 1 {
		t.Fatalf("CC count = %d, want 1", len(m.GetCC()))
	}

	if len(m.GetBCC()) != 1 {
		t.Fatalf("BCC count = %d, want 1", len(m.GetBCC()))
	}

	if m.GetBCC()[0].Name != "Secret" {
		t.Fatalf("BCC name = %q, want %q", m.GetBCC()[0].Name, "Secret")
	}
}

func TestMailMessageAttach(t *testing.T) {
	t.Parallel()

	m := notifications.NewMailMessage().
		Attach("/path/to/file.pdf", mail.WithName("document.pdf"))

	attachments := m.GetAttachments()

	if len(attachments) != 1 {
		t.Fatalf("attachments count = %d, want 1", len(attachments))
	}

	if attachments[0].Path != "/path/to/file.pdf" {
		t.Fatalf("path = %q, want %q", attachments[0].Path, "/path/to/file.pdf")
	}

	if attachments[0].Name != "document.pdf" {
		t.Fatalf("name = %q, want %q", attachments[0].Name, "document.pdf")
	}
}

func TestMailMessageAttachData(t *testing.T) {
	t.Parallel()

	m := notifications.NewMailMessage().
		AttachData([]byte("hello"), "greeting.txt", "text/plain")

	raw := m.GetRawAttachments()

	if len(raw) != 1 {
		t.Fatalf("raw attachments count = %d, want 1", len(raw))
	}

	if raw[0].Name != "greeting.txt" {
		t.Fatalf("name = %q, want %q", raw[0].Name, "greeting.txt")
	}

	if raw[0].Mime != "text/plain" {
		t.Fatalf("mime = %q, want %q", raw[0].Mime, "text/plain")
	}
}

func TestMailMessageAttachMany(t *testing.T) {
	t.Parallel()

	files := []*mail.Attachment{
		{Path: "/a.txt"},
		{Path: "/b.txt"},
	}

	m := notifications.NewMailMessage().AttachMany(files)

	if len(m.GetAttachments()) != 2 {
		t.Fatalf("attachments count = %d, want 2", len(m.GetAttachments()))
	}
}

func TestMailMessageTagsAndMetadata(t *testing.T) {
	t.Parallel()

	m := notifications.NewMailMessage().
		Tag("transactional").
		Tag("welcome").
		Metadata("campaign", "onboarding").
		Metadata("version", "2")

	if len(m.GetTags()) != 2 {
		t.Fatalf("tags count = %d, want 2", len(m.GetTags()))
	}

	if m.GetMetadata()["campaign"] != "onboarding" {
		t.Fatalf("metadata campaign = %q, want %q", m.GetMetadata()["campaign"], "onboarding")
	}
}

func TestMailMessagePriority(t *testing.T) {
	t.Parallel()

	m := notifications.NewMailMessage().Priority(1)

	if m.GetPriority() != 1 {
		t.Fatalf("priority = %d, want 1", m.GetPriority())
	}
}

func TestMailMessageCallback(t *testing.T) {
	t.Parallel()

	called := false
	m := notifications.NewMailMessage().WithSymfonyMessage(func(_ *mail.Message) {
		called = true
	})

	if len(m.GetCallbacks()) != 1 {
		t.Fatalf("callbacks count = %d, want 1", len(m.GetCallbacks()))
	}

	m.GetCallbacks()[0](mail.NewMessage())

	if !called {
		t.Fatal("callback was not invoked")
	}
}

func TestMailMessageData(t *testing.T) {
	t.Parallel()

	m := notifications.NewMailMessage().
		Subject("Test").
		View("view", map[string]any{"extra": "value"})

	data := m.Data()

	if data["subject"] != "Test" {
		t.Fatalf("subject = %v, want %q", data["subject"], "Test")
	}

	if data["extra"] != "value" {
		t.Fatalf("extra = %v, want %q", data["extra"], "value")
	}
}

func TestMailMessageLevelChaining(t *testing.T) {
	t.Parallel()

	m := notifications.NewMailMessage().
		Success().
		Subject("OK")

	if m.GetLevel() != "success" {
		t.Fatalf("level = %q, want %q", m.GetLevel(), "success")
	}

	if m.GetSubject() != "OK" {
		t.Fatalf("subject = %q, want %q", m.GetSubject(), "OK")
	}
}
