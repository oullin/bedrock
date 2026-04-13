package mail

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
)

// Message provides a handle to the email message being constructed.
// It exposes methods for setting addresses, content, attachments,
// and headers on the underlying message during send callbacks.
type Message struct {
	from        Address
	sender      Address
	returnPath  string
	to          []Address
	cc          []Address
	bcc         []Address
	replyTo     []Address
	subject     string
	htmlBody    string
	textBody    string
	priority    int
	headers     map[string][]string
	attachments []*Attachment
	embeds      []*Embed
}

// Embed represents an inline-embedded resource with a Content-ID.
type Embed struct {
	CID  string
	Data []byte
	Name string
	Mime string
}

// NewMessage creates a new empty Message.

// From sets the sender address.

// GetFrom returns the sender address.

// Sender sets the actual sender (Sender header), distinct from From.

// GetSender returns the Sender header address.

// ReturnPath sets the Return-Path address for bounce handling.

// GetReturnPath returns the Return-Path address.

// To adds primary recipients. If override is true, existing To
// addresses are replaced.

// ForgetTo removes all To recipients.

// GetTo returns the To recipients.

// CC adds carbon-copy recipients. If override is true, existing CC
// addresses are replaced.

// ForgetCC removes all CC recipients.

// GetCC returns the CC recipients.

// BCC adds blind-carbon-copy recipients. If override is true, existing
// BCC addresses are replaced.

// ForgetBCC removes all BCC recipients.

// GetBCC returns the BCC recipients.

// ReplyTo adds reply-to addresses.

// GetReplyTo returns the Reply-To addresses.

// Subject sets the email subject line.

// GetSubject returns the subject line.

// Priority sets the email priority (1 = highest, 5 = lowest).

// GetPriority returns the priority level.

// SetHTMLBody sets the HTML body content.

// GetHTMLBody returns the HTML body.

// SetTextBody sets the plain-text body content.

// GetTextBody returns the plain-text body.

// Attach adds a file attachment by path.

// AttachData adds an attachment from raw data.

// AttachWith adds a pre-built Attachment to the message.

// GetAttachments returns all file attachments.

// Embed embeds a file inline and returns the Content-ID reference.

// EmbedData embeds raw data inline and returns the Content-ID reference.

// GetEmbeds returns all inline embeds.

// SetHeader sets a custom header value.

// GetCustomHeaders returns all custom headers.

// AllRecipients returns a deduplicated list of all recipients
// (To + CC + BCC).

// AttachOption configures an attachment.
type AttachOption func(*Attachment)

func NewMessage() *Message {
	return &Message{
		headers: make(map[string][]string),
	}
}

func (m *Message) From(address string, name ...string) *Message {
	m.from = Address{Email: address}

	if len(name) > 0 {
		m.from.Name = name[0]
	}

	return m
}

func (m *Message) GetFrom() Address {
	return m.from
}

func (m *Message) Sender(address string, name ...string) *Message {
	m.sender = Address{Email: address}

	if len(name) > 0 {
		m.sender.Name = name[0]
	}

	return m
}

func (m *Message) GetSender() Address {
	return m.sender
}

func (m *Message) ReturnPath(address string) *Message {
	m.returnPath = address

	return m
}

func (m *Message) GetReturnPath() string {
	return m.returnPath
}

func (m *Message) To(addresses []Address, override ...bool) *Message {
	if len(override) > 0 && override[0] {
		m.to = nil
	}

	m.to = append(m.to, addresses...)

	return m
}

func (m *Message) ForgetTo() *Message {
	m.to = nil

	return m
}

func (m *Message) GetTo() []Address {
	return m.to
}

func (m *Message) CC(addresses []Address, override ...bool) *Message {
	if len(override) > 0 && override[0] {
		m.cc = nil
	}

	m.cc = append(m.cc, addresses...)

	return m
}

func (m *Message) ForgetCC() *Message {
	m.cc = nil

	return m
}

func (m *Message) GetCC() []Address {
	return m.cc
}

func (m *Message) BCC(addresses []Address, override ...bool) *Message {
	if len(override) > 0 && override[0] {
		m.bcc = nil
	}

	m.bcc = append(m.bcc, addresses...)

	return m
}

func (m *Message) ForgetBCC() *Message {
	m.bcc = nil

	return m
}

func (m *Message) GetBCC() []Address {
	return m.bcc
}

func (m *Message) ReplyTo(addresses []Address) *Message {
	m.replyTo = append(m.replyTo, addresses...)

	return m
}

func (m *Message) GetReplyTo() []Address {
	return m.replyTo
}

func (m *Message) Subject(subject string) *Message {
	m.subject = subject

	return m
}

func (m *Message) GetSubject() string {
	return m.subject
}

func (m *Message) Priority(level int) *Message {
	m.priority = level

	return m
}

func (m *Message) GetPriority() int {
	return m.priority
}

func (m *Message) SetHTMLBody(html string) *Message {
	m.htmlBody = html

	return m
}

func (m *Message) GetHTMLBody() string {
	return m.htmlBody
}

func (m *Message) SetTextBody(text string) *Message {
	m.textBody = text

	return m
}

func (m *Message) GetTextBody() string {
	return m.textBody
}

func (m *Message) Attach(path string, options ...AttachOption) *Message {
	a := &Attachment{Path: path}

	for _, opt := range options {
		opt(a)
	}

	m.attachments = append(m.attachments, a)

	return m
}

func (m *Message) AttachData(data []byte, name string, options ...AttachOption) *Message {
	a := FromData(func() (io.Reader, error) {
		return bytes.NewReader(data), nil
	}, name)

	for _, opt := range options {
		opt(a)
	}

	m.attachments = append(m.attachments, a)

	return m
}

func (m *Message) AttachWith(a *Attachment) *Message {
	m.attachments = append(m.attachments, a)

	return m
}

func (m *Message) GetAttachments() []*Attachment {
	return m.attachments
}

func (m *Message) Embed(path string, options ...AttachOption) string {
	cid := generateCID()
	a := &Attachment{Path: path, Inline: true}

	for _, opt := range options {
		opt(a)
	}

	m.embeds = append(m.embeds, &Embed{
		CID:  cid,
		Name: a.Name,
		Mime: a.Mime,
	})

	return "cid:" + cid
}

func (m *Message) EmbedData(data []byte, name string, mime ...string) string {
	cid := generateCID()
	mimeType := "application/octet-stream"

	if len(mime) > 0 {
		mimeType = mime[0]
	}

	m.embeds = append(m.embeds, &Embed{
		CID:  cid,
		Data: data,
		Name: name,
		Mime: mimeType,
	})

	return "cid:" + cid
}

func (m *Message) GetEmbeds() []*Embed {
	return m.embeds
}

func (m *Message) SetHeader(name string, values ...string) *Message {
	m.headers[name] = values

	return m
}

func (m *Message) GetCustomHeaders() map[string][]string {
	return m.headers
}

func (m *Message) AllRecipients() []Address {
	seen := make(map[string]struct{})

	var all []Address

	for _, lists := range [][]Address{m.to, m.cc, m.bcc} {
		for _, a := range lists {
			if _, ok := seen[a.Email]; !ok {
				seen[a.Email] = struct{}{}
				all = append(all, a)
			}
		}
	}

	return all
}

// WithName sets the attachment display filename.
func WithName(name string) AttachOption {
	return func(a *Attachment) {
		a.Name = name
	}
}

// WithMimeType sets the attachment MIME content type.
func WithMimeType(mime string) AttachOption {
	return func(a *Attachment) {
		a.Mime = mime
	}
}

func generateCID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)

	return fmt.Sprintf("%x", b)
}
