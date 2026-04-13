package mailx_test

import (
	"bytes"
	"io"
	"testing"

	cmail "github.com/bedrock/packages/contracts/mail"
	"github.com/bedrock/packages/mailx"
)

func TestMailableSetsRecipientsCorrectly(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetTo(
		cmail.Address{Email: "foo@bar.baz", Name: "Foo"},
		cmail.Address{Email: "bar@bar.baz"},
	)

	if !m.HasTo("foo@bar.baz", "Foo") {
		t.Fatal("expected to have foo@bar.baz with name Foo")
	}

	if !m.HasTo("bar@bar.baz") {
		t.Fatal("expected to have bar@bar.baz")
	}

	if m.HasTo("baz@bar.baz") {
		t.Fatal("should not have baz@bar.baz")
	}
}

func TestMailableSetsCCRecipientsCorrectly(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetCC(
		cmail.Address{Email: "cc@bar.baz", Name: "CC"},
	)

	if !m.HasCC("cc@bar.baz", "CC") {
		t.Fatal("expected to have cc@bar.baz")
	}

	if m.HasCC("nope@bar.baz") {
		t.Fatal("should not have nope@bar.baz")
	}
}

func TestMailableSetsBCCRecipientsCorrectly(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetBCC(
		cmail.Address{Email: "bcc@bar.baz"},
	)

	if !m.HasBCC("bcc@bar.baz") {
		t.Fatal("expected to have bcc@bar.baz")
	}
}

func TestMailableSetsReplyToCorrectly(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetReplyTo(
		cmail.Address{Email: "reply@bar.baz", Name: "Reply"},
	)

	if !m.HasReplyTo("reply@bar.baz", "Reply") {
		t.Fatal("expected to have reply@bar.baz")
	}
}

func TestMailableSetsFromCorrectly(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetFrom("from@bar.baz", "From")

	if !m.HasFrom("from@bar.baz", "From") {
		t.Fatal("expected from from@bar.baz")
	}

	if m.HasFrom("other@bar.baz") {
		t.Fatal("should not have other@bar.baz as from")
	}
}

func TestMailableSetsSubjectCorrectly(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetSubject("Test Subject")

	if !m.HasSubject("Test Subject") {
		t.Fatal("expected subject Test Subject")
	}

	if m.HasSubject("Wrong Subject") {
		t.Fatal("should not have Wrong Subject")
	}
}

func TestMailableIgnoresDuplicateRawAttachments(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	a := &cmail.Attachment{Name: "file.txt", Mime: "text/plain"}

	m.Attach(a)
	m.Attach(a)

	attachments := m.GetAttachments()

	if len(attachments) != 1 {
		t.Fatalf("expected 1 attachment after dedup, got %d", len(attachments))
	}
}

func TestMailableIgnoresDuplicateStorageAttachments(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}

	m.AttachFromStorageDisk("local", "/path/file.txt")
	m.AttachFromStorageDisk("local", "/path/file.txt")

	attachments := m.GetAttachments()

	if len(attachments) != 1 {
		t.Fatalf("expected 1 attachment after dedup, got %d", len(attachments))
	}
}

func TestMailableBuildsViewData(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.With("key", "value")

	content := m.GetContent()

	if content.With["key"] != "value" {
		t.Fatalf("expected with key=value, got %v", content.With["key"])
	}
}

func TestMailableMailerMayBeSet(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetMailer("postmark")

	if m.GetMailerName() != "postmark" {
		t.Fatalf("expected mailer postmark, got %s", m.GetMailerName())
	}
}

func TestMailablePriority(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetPriority(2)

	callbacks := m.GetCallbacks()

	if len(callbacks) != 1 {
		t.Fatalf("expected 1 callback for priority, got %d", len(callbacks))
	}

	msg := cmail.NewMessage()
	callbacks[0](msg)

	if msg.GetPriority() != 2 {
		t.Fatalf("expected priority 2, got %d", msg.GetPriority())
	}
}

func TestMailableMetadata(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetMetadata("color", "blue")

	if !m.HasMetadata("color", "blue") {
		t.Fatal("expected metadata color=blue")
	}

	if m.HasMetadata("color", "red") {
		t.Fatal("should not have metadata color=red")
	}
}

func TestMailableMergeMetadata(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetMetadata("a", "1")
	m.SetMetadata("b", "2")

	if !m.HasMetadata("a", "1") {
		t.Fatal("expected metadata a=1")
	}

	if !m.HasMetadata("b", "2") {
		t.Fatal("expected metadata b=2")
	}
}

func TestMailableTag(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.Tag("promo")

	if !m.HasTag("promo") {
		t.Fatal("expected tag promo")
	}

	if m.HasTag("other") {
		t.Fatal("should not have tag other")
	}
}

func TestMailableAttachMultipleFiles(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.AttachMany([]*cmail.Attachment{
		{Path: "/a.txt", Name: "a.txt"},
		{Path: "/b.txt", Name: "b.txt"},
	})

	if len(m.GetAttachments()) != 2 {
		t.Fatalf("expected 2 attachments, got %d", len(m.GetAttachments()))
	}
}

func TestMailableAttachViaAttachableFromPath(t *testing.T) {
	t.Parallel()

	attachable := &testPathAttachable{path: "/photo.jpg", name: "photo.jpg", mime: "image/jpeg"}
	a, _ := attachable.ToMailAttachment()

	m := &mailx.Mailable{}
	m.Attach(a)

	if !m.HasAttachmentFromPath("/photo.jpg", cmail.WithName("photo.jpg")) {
		t.Fatal("expected attachment from path /photo.jpg")
	}
}

func TestMailableAttachViaAttachableFromData(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.AttachData(func() (io.Reader, error) {
		return bytes.NewReader([]byte("data")), nil
	}, "doc.txt", cmail.WithMimeType("text/plain"))

	if len(m.GetAttachments()) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(m.GetAttachments()))
	}

	if m.GetAttachments()[0].Name != "doc.txt" {
		t.Fatalf("expected name doc.txt, got %s", m.GetAttachments()[0].Name)
	}
}

func TestMailableCanJitNameAttachments(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.AttachFromPath("/report.pdf", cmail.WithName("Q4 Report.pdf"))

	if !m.HasAttachmentFromPath("/report.pdf", cmail.WithName("Q4 Report.pdf")) {
		t.Fatal("expected attachment with jit name")
	}
}

func TestMailableHasAttachmentWithJitNamedAttachment(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.AttachFromPath("/invoice.pdf", cmail.WithName("Invoice #123.pdf"))

	if !m.HasAttachmentFromPath("/invoice.pdf", cmail.WithName("Invoice #123.pdf")) {
		t.Fatal("expected to find jit-named attachment")
	}
}

func TestMailableCheckForPathBasedAttachments(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.AttachFromPath("/file.txt", cmail.WithName("file.txt"), cmail.WithMimeType("text/plain"))

	if !m.HasAttachmentFromPath("/file.txt", cmail.WithName("file.txt"), cmail.WithMimeType("text/plain")) {
		t.Fatal("expected path-based attachment")
	}

	if m.HasAttachmentFromPath("/other.txt") {
		t.Fatal("should not find non-existent attachment")
	}
}

func TestMailableCheckForAttachmentBasedAttachments(t *testing.T) {
	t.Parallel()

	a := &cmail.Attachment{Path: "/doc.pdf", Name: "doc.pdf", Mime: "application/pdf"}
	m := &mailx.Mailable{}
	m.Attach(a)

	if !m.HasAttachment(&cmail.Attachment{Path: "/doc.pdf", Name: "doc.pdf", Mime: "application/pdf"}) {
		t.Fatal("expected to find attachment")
	}
}

func TestMailableCheckForDataBasedAttachments(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.AttachData(func() (io.Reader, error) {
		return bytes.NewReader([]byte("data")), nil
	}, "data.csv", cmail.WithMimeType("text/csv"))

	if !m.HasAttachment(&cmail.Attachment{Name: "data.csv", Mime: "text/csv"}) {
		t.Fatal("expected data-based attachment")
	}
}

func TestMailableCheckForStorageBasedAttachments(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.AttachFromStorageDisk("s3", "/uploads/file.zip")

	if !m.HasAttachmentFromStorageDisk("s3", "/uploads/file.zip") {
		t.Fatal("expected storage-based attachment")
	}
}

func TestMailableAssertHasNoAttachments(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.AssertHasNoAttachments(t)
}

func TestMailableAssertHasAttachment(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	a := &cmail.Attachment{Path: "/file.pdf", Name: "file.pdf"}
	m.Attach(a)

	m.AssertHasAttachment(t, a)
}

func TestMailableAssertHasAttachedData(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.AttachData(func() (io.Reader, error) {
		return bytes.NewReader([]byte("contents")), nil
	}, "report.csv")

	m.AssertHasAttachment(t, &cmail.Attachment{Name: "report.csv"})
}

func TestMailableAssertHasAttachmentFromStorage(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.AttachFromStorage("/docs/manual.pdf")

	if !m.HasAttachmentFromStorage("/docs/manual.pdf") {
		t.Fatal("expected storage attachment")
	}
}

func TestMailableAssertHasSubject(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetSubject("Hello World")

	m.AssertHasSubject(t, "Hello World")
}

func TestMailableHeaders(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetMessageID("<abc@example.com>")
	m.SetReferences([]string{"<ref1@example.com>", "<ref2@example.com>"})
	m.AddHeader("X-Custom", "value")

	headers := m.GetHeaders()

	if headers.MessageID != "<abc@example.com>" {
		t.Fatalf("expected message id <abc@example.com>, got %s", headers.MessageID)
	}

	if len(headers.References) != 2 {
		t.Fatalf("expected 2 references, got %d", len(headers.References))
	}

	if len(headers.Text) != 1 {
		t.Fatalf("expected 1 text header, got %d", len(headers.Text))
	}

	if headers.Text[0].Value != "value" {
		t.Fatalf("expected header value 'value', got %s", headers.Text[0].Value)
	}
}

func TestMailableAttributesInBuild(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetFrom("build@example.com", "Builder").
		SetSubject("Built").
		SetTo(cmail.Address{Email: "to@example.com"}).
		SetHTML("<h1>Hello</h1>")

	m.AssertFrom(t, "build@example.com", "Builder")
	m.AssertHasSubject(t, "Built")
	m.AssertTo(t, "to@example.com")
}

func TestMailableCanBeTapped(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}

	m.Tap(func(mb *mailx.Mailable) {
		mb.SetSubject("Tapped Subject")
	})

	m.AssertHasSubject(t, "Tapped Subject")
}

func TestMailableAddTo(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.AddTo(cmail.Address{Email: "a@test.com"})
	m.AddTo(cmail.Address{Email: "b@test.com"})

	if len(m.GetEnvelope().To) != 2 {
		t.Fatalf("expected 2 To, got %d", len(m.GetEnvelope().To))
	}
}

func TestMailableAddCC(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.AddCC(cmail.Address{Email: "a@test.com"})
	m.AddCC(cmail.Address{Email: "b@test.com"})

	if len(m.GetEnvelope().CC) != 2 {
		t.Fatalf("expected 2 CC, got %d", len(m.GetEnvelope().CC))
	}
}

func TestMailableAddBCC(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.AddBCC(cmail.Address{Email: "a@test.com"})
	m.AddBCC(cmail.Address{Email: "b@test.com"})

	if len(m.GetEnvelope().BCC) != 2 {
		t.Fatalf("expected 2 BCC, got %d", len(m.GetEnvelope().BCC))
	}
}

func TestMailableAddReplyTo(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.AddReplyTo(cmail.Address{Email: "a@test.com"})
	m.AddReplyTo(cmail.Address{Email: "b@test.com"})

	if len(m.GetEnvelope().ReplyTo) != 2 {
		t.Fatalf("expected 2 ReplyTo, got %d", len(m.GetEnvelope().ReplyTo))
	}
}

func TestMailableLocale(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetLocale("es")

	if m.GetLocale() != "es" {
		t.Fatalf("expected locale es, got %s", m.GetLocale())
	}
}

func TestMailableWithData(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.WithData(map[string]any{"a": 1, "b": 2})

	content := m.GetContent()

	if content.With["a"] != 1 {
		t.Fatalf("expected a=1, got %v", content.With["a"])
	}
}

func TestMailableSetHTMLView(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetHTMLView("emails.welcome")

	content := m.GetContent()

	if content.HTML != "emails.welcome" {
		t.Fatalf("expected view emails.welcome, got %s", content.HTML)
	}

	if content.HTMLString {
		t.Fatal("expected HTMLString to be false for views")
	}
}

func TestMailableSetText(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetText("emails.welcome_text")

	content := m.GetContent()

	if content.Text != "emails.welcome_text" {
		t.Fatalf("expected text emails.welcome_text, got %s", content.Text)
	}
}

func TestMailableSetMarkdown(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetMarkdown("# Hello")

	content := m.GetContent()

	if content.Markdown != "# Hello" {
		t.Fatalf("expected markdown # Hello, got %s", content.Markdown)
	}
}

func TestMailableAssertHasTag(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.Tag("welcome")

	m.AssertHasTag(t, "welcome")
}

func TestMailableAssertHasMetadata(t *testing.T) {
	t.Parallel()

	m := &mailx.Mailable{}
	m.SetMetadata("campaign", "q4")

	m.AssertHasMetadata(t, "campaign", "q4")
}
