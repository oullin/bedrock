package mailx_test

import (
	"strings"
	"testing"

	cmail "github.com/bedrock/packages/contracts/mail"
)

// Name should remain empty; the transport layer determines the name.

// --- Test helpers ---

type testPathAttachable struct {
	path string
	name string
	mime string
}

type testDataAttachable struct {
	data []byte
	name string
	mime string
}

func TestMessageFromMethod(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.From("foo@bar.baz", "Foo")

	from := m.GetFrom()

	if from.Email != "foo@bar.baz" {
		t.Fatalf("expected email foo@bar.baz, got %s", from.Email)
	}

	if from.Name != "Foo" {
		t.Fatalf("expected name Foo, got %s", from.Name)
	}
}

func TestMessageSenderMethod(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.Sender("foo@bar.baz", "Foo")

	sender := m.GetSender()

	if sender.Email != "foo@bar.baz" {
		t.Fatalf("expected sender foo@bar.baz, got %s", sender.Email)
	}

	if sender.Name != "Foo" {
		t.Fatalf("expected name Foo, got %s", sender.Name)
	}
}

func TestMessageReturnPathMethod(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.ReturnPath("foo@bar.baz")

	if m.GetReturnPath() != "foo@bar.baz" {
		t.Fatalf("expected return path foo@bar.baz, got %s", m.GetReturnPath())
	}
}

func TestMessageToMethod(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.To([]cmail.Address{
		{Email: "foo@bar.baz", Name: "Foo"},
		{Email: "bar@bar.baz"},
	})

	to := m.GetTo()

	if len(to) != 2 {
		t.Fatalf("expected 2 To recipients, got %d", len(to))
	}

	if to[0].Email != "foo@bar.baz" {
		t.Fatalf("expected foo@bar.baz, got %s", to[0].Email)
	}

	if to[0].Name != "Foo" {
		t.Fatalf("expected Foo, got %s", to[0].Name)
	}
}

func TestMessageToMethodWithOverride(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.To([]cmail.Address{{Email: "foo@bar.baz"}})
	m.To([]cmail.Address{{Email: "bar@bar.baz"}}, true)

	to := m.GetTo()

	if len(to) != 1 {
		t.Fatalf("expected 1 To recipient after override, got %d", len(to))
	}

	if to[0].Email != "bar@bar.baz" {
		t.Fatalf("expected bar@bar.baz, got %s", to[0].Email)
	}
}

func TestMessageCCMethod(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.CC([]cmail.Address{
		{Email: "foo@bar.baz", Name: "Foo"},
	})

	cc := m.GetCC()

	if len(cc) != 1 {
		t.Fatalf("expected 1 CC, got %d", len(cc))
	}

	if cc[0].Email != "foo@bar.baz" {
		t.Fatalf("expected foo@bar.baz, got %s", cc[0].Email)
	}
}

func TestMessageBCCMethod(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.BCC([]cmail.Address{
		{Email: "foo@bar.baz", Name: "Foo"},
	})

	bcc := m.GetBCC()

	if len(bcc) != 1 {
		t.Fatalf("expected 1 BCC, got %d", len(bcc))
	}

	if bcc[0].Email != "foo@bar.baz" {
		t.Fatalf("expected foo@bar.baz, got %s", bcc[0].Email)
	}
}

func TestMessageReplyToMethod(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.ReplyTo([]cmail.Address{
		{Email: "foo@bar.baz", Name: "Foo"},
	})

	replyTo := m.GetReplyTo()

	if len(replyTo) != 1 {
		t.Fatalf("expected 1 ReplyTo, got %d", len(replyTo))
	}

	if replyTo[0].Email != "foo@bar.baz" {
		t.Fatalf("expected foo@bar.baz, got %s", replyTo[0].Email)
	}
}

func TestMessageSubjectMethod(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.Subject("Test Subject")

	if m.GetSubject() != "Test Subject" {
		t.Fatalf("expected Test Subject, got %s", m.GetSubject())
	}
}

func TestMessagePriorityMethod(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.Priority(1)

	if m.GetPriority() != 1 {
		t.Fatalf("expected priority 1, got %d", m.GetPriority())
	}
}

func TestMessageBasicAttachment(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.Attach("/path/to/file.pdf", cmail.WithName("report.pdf"), cmail.WithMimeType("application/pdf"))

	attachments := m.GetAttachments()

	if len(attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(attachments))
	}

	if attachments[0].Path != "/path/to/file.pdf" {
		t.Fatalf("expected path /path/to/file.pdf, got %s", attachments[0].Path)
	}

	if attachments[0].Name != "report.pdf" {
		t.Fatalf("expected name report.pdf, got %s", attachments[0].Name)
	}

	if attachments[0].Mime != "application/pdf" {
		t.Fatalf("expected mime application/pdf, got %s", attachments[0].Mime)
	}
}

func TestMessageDataAttachment(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.AttachData([]byte("file contents"), "data.txt", cmail.WithMimeType("text/plain"))

	attachments := m.GetAttachments()

	if len(attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(attachments))
	}

	if attachments[0].Name != "data.txt" {
		t.Fatalf("expected name data.txt, got %s", attachments[0].Name)
	}

	if attachments[0].Mime != "text/plain" {
		t.Fatalf("expected mime text/plain, got %s", attachments[0].Mime)
	}
}

func TestMessageAttachViaAttachableFromPath(t *testing.T) {
	t.Parallel()

	attachable := &testPathAttachable{path: "/path/to/file.pdf", name: "report.pdf", mime: "application/pdf"}

	a, err := attachable.ToMailAttachment()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := cmail.NewMessage()
	m.AttachWith(a)

	attachments := m.GetAttachments()

	if len(attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(attachments))
	}

	if attachments[0].Path != "/path/to/file.pdf" {
		t.Fatalf("expected path /path/to/file.pdf, got %s", attachments[0].Path)
	}
}

func TestMessageAttachViaAttachableFromData(t *testing.T) {
	t.Parallel()

	attachable := &testDataAttachable{data: []byte("contents"), name: "doc.txt", mime: "text/plain"}

	a, err := attachable.ToMailAttachment()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := cmail.NewMessage()
	m.AttachWith(a)

	attachments := m.GetAttachments()

	if len(attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(attachments))
	}

	if attachments[0].Name != "doc.txt" {
		t.Fatalf("expected name doc.txt, got %s", attachments[0].Name)
	}
}

func TestMessageEmbedPath(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	cid := m.Embed("/path/to/image.png", cmail.WithName("logo.png"), cmail.WithMimeType("image/png"))

	if !strings.HasPrefix(cid, "cid:") {
		t.Fatalf("expected cid: prefix, got %s", cid)
	}

	embeds := m.GetEmbeds()

	if len(embeds) != 1 {
		t.Fatalf("expected 1 embed, got %d", len(embeds))
	}

	if embeds[0].Name != "logo.png" {
		t.Fatalf("expected name logo.png, got %s", embeds[0].Name)
	}
}

func TestMessageEmbedData(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	cid := m.EmbedData([]byte("image data"), "photo.jpg", "image/jpeg")

	if !strings.HasPrefix(cid, "cid:") {
		t.Fatalf("expected cid: prefix, got %s", cid)
	}

	embeds := m.GetEmbeds()

	if len(embeds) != 1 {
		t.Fatalf("expected 1 embed, got %d", len(embeds))
	}

	if embeds[0].Mime != "image/jpeg" {
		t.Fatalf("expected mime image/jpeg, got %s", embeds[0].Mime)
	}
}

func TestMessageEmbedViaAttachableFromPath(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	cid := m.Embed("/path/to/logo.png", cmail.WithName("logo.png"), cmail.WithMimeType("image/png"))

	if !strings.HasPrefix(cid, "cid:") {
		t.Fatalf("expected cid: prefix, got %s", cid)
	}

	embeds := m.GetEmbeds()

	if len(embeds) != 1 {
		t.Fatalf("expected 1 embed, got %d", len(embeds))
	}
}

func TestMessageGeneratesRandomNameWhenAttachableHasNone(t *testing.T) {
	t.Parallel()

	attachable := &testPathAttachable{path: "/path/to/file.pdf", name: "", mime: "application/pdf"}

	a, err := attachable.ToMailAttachment()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := cmail.NewMessage()
	m.AttachWith(a)

	attachments := m.GetAttachments()

	if len(attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(attachments))
	}

	if attachments[0].Name != "" {
		t.Fatalf("expected empty name for unnamed attachable, got %s", attachments[0].Name)
	}
}

func TestMessageForgetTo(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.To([]cmail.Address{{Email: "foo@bar.baz"}})
	m.ForgetTo()

	if len(m.GetTo()) != 0 {
		t.Fatalf("expected 0 To recipients after ForgetTo, got %d", len(m.GetTo()))
	}
}

func TestMessageForgetCC(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.CC([]cmail.Address{{Email: "foo@bar.baz"}})
	m.ForgetCC()

	if len(m.GetCC()) != 0 {
		t.Fatalf("expected 0 CC recipients after ForgetCC, got %d", len(m.GetCC()))
	}
}

func TestMessageForgetBCC(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.BCC([]cmail.Address{{Email: "foo@bar.baz"}})
	m.ForgetBCC()

	if len(m.GetBCC()) != 0 {
		t.Fatalf("expected 0 BCC recipients after ForgetBCC, got %d", len(m.GetBCC()))
	}
}

func TestMessageAllRecipients(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.To([]cmail.Address{{Email: "a@test.com"}})
	m.CC([]cmail.Address{{Email: "b@test.com"}})
	m.BCC([]cmail.Address{{Email: "c@test.com"}, {Email: "a@test.com"}})

	all := m.AllRecipients()

	if len(all) != 3 {
		t.Fatalf("expected 3 unique recipients, got %d", len(all))
	}
}

func TestMessageSetHeader(t *testing.T) {
	t.Parallel()

	m := cmail.NewMessage()
	m.SetHeader("X-Custom", "value1", "value2")

	headers := m.GetCustomHeaders()

	if len(headers["X-Custom"]) != 2 {
		t.Fatalf("expected 2 header values, got %d", len(headers["X-Custom"]))
	}
}

func (a *testPathAttachable) ToMailAttachment() (*cmail.Attachment, error) {
	att := cmail.FromPath(a.path)

	if a.name != "" {
		att.As(a.name)
	}

	if a.mime != "" {
		att.WithMime(a.mime)
	}

	return att, nil
}

func (a *testDataAttachable) ToMailAttachment() (*cmail.Attachment, error) {
	att := &cmail.Attachment{
		Name: a.name,
		Mime: a.mime,
	}

	return att, nil
}
