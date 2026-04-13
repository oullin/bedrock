package notifications

import (
	"github.com/bedrock/packages/contracts/mail"
)

// RawAttachment holds raw byte data for an email attachment.
type RawAttachment struct {
	Data []byte
	Name string
	Mime string
}

// MailMessage is a fluent builder for email notification messages. It embeds
// SimpleMessage for body/level/action and adds mail-specific fields like
// view, from, recipients, attachments, tags, and metadata.
type MailMessage struct {
	*SimpleMessage
	view           string
	viewData       map[string]any
	markdown       string
	theme          string
	from           *mail.Address
	replyTo        []mail.Address
	cc             []mail.Address
	bcc            []mail.Address
	attachments    []*mail.Attachment
	rawAttachments []RawAttachment
	tags           []string
	metadata       map[string]string
	priority       int
	callbacks      []func(*mail.Message)
}

// NewMailMessage creates a MailMessage with default info level.
func NewMailMessage() *MailMessage {
	return &MailMessage{
		SimpleMessage: NewSimpleMessage(),
		viewData:      make(map[string]any),
		metadata:      make(map[string]string),
	}
}

// View sets the view template name.
func (m *MailMessage) View(view string, data ...map[string]any) *MailMessage {
	m.view = view

	if len(data) > 0 {
		for k, v := range data[0] {
			m.viewData[k] = v
		}
	}

	return m
}

// GetView returns the view template name.
func (m *MailMessage) GetView() string { return m.view }

// GetViewData returns the view template data.
func (m *MailMessage) GetViewData() map[string]any { return m.viewData }

// Text sets a plain-text view for the message.
func (m *MailMessage) Text(view string, data ...map[string]any) *MailMessage {
	return m.View(view, data...)
}

// Markdown sets the markdown template name.
func (m *MailMessage) Markdown(view string, data ...map[string]any) *MailMessage {
	m.markdown = view

	if len(data) > 0 {
		for k, v := range data[0] {
			m.viewData[k] = v
		}
	}

	return m
}

// GetMarkdown returns the markdown template name.
func (m *MailMessage) GetMarkdown() string { return m.markdown }

// Template sets the Template template used for rendering.
func (m *MailMessage) Template(template string) *MailMessage {
	m.view = template

	return m
}

// Theme sets the markdown theme.
func (m *MailMessage) Theme(theme string) *MailMessage {
	m.theme = theme

	return m
}

// GetTheme returns the markdown theme.
func (m *MailMessage) GetTheme() string { return m.theme }

// From sets the sender address.
func (m *MailMessage) From(address string, name ...string) *MailMessage {
	a := mail.Address{Email: address}

	if len(name) > 0 {
		a.Name = name[0]
	}

	m.from = &a

	return m
}

// GetFrom returns the sender address.
func (m *MailMessage) GetFrom() *mail.Address { return m.from }

// ReplyTo adds a reply-to address.
func (m *MailMessage) ReplyTo(address string, name ...string) *MailMessage {
	a := mail.Address{Email: address}

	if len(name) > 0 {
		a.Name = name[0]
	}

	m.replyTo = append(m.replyTo, a)

	return m
}

// GetReplyTo returns the reply-to addresses.
func (m *MailMessage) GetReplyTo() []mail.Address { return m.replyTo }

// CC adds a carbon-copy recipient.
func (m *MailMessage) CC(address string, name ...string) *MailMessage {
	a := mail.Address{Email: address}

	if len(name) > 0 {
		a.Name = name[0]
	}

	m.cc = append(m.cc, a)

	return m
}

// GetCC returns the CC recipients.
func (m *MailMessage) GetCC() []mail.Address { return m.cc }

// BCC adds a blind-carbon-copy recipient.
func (m *MailMessage) BCC(address string, name ...string) *MailMessage {
	a := mail.Address{Email: address}

	if len(name) > 0 {
		a.Name = name[0]
	}

	m.bcc = append(m.bcc, a)

	return m
}

// GetBCC returns the BCC recipients.
func (m *MailMessage) GetBCC() []mail.Address { return m.bcc }

// Attach adds a file attachment by path.
func (m *MailMessage) Attach(path string, opts ...mail.AttachOption) *MailMessage {
	a := &mail.Attachment{Path: path}

	for _, opt := range opts {
		opt(a)
	}

	m.attachments = append(m.attachments, a)

	return m
}

// AttachMany attaches multiple files.
func (m *MailMessage) AttachMany(files []*mail.Attachment) *MailMessage {
	m.attachments = append(m.attachments, files...)

	return m
}

// AttachData attaches raw data with a filename.
func (m *MailMessage) AttachData(data []byte, name string, mime ...string) *MailMessage {
	r := RawAttachment{Data: data, Name: name}

	if len(mime) > 0 {
		r.Mime = mime[0]
	}

	m.rawAttachments = append(m.rawAttachments, r)

	return m
}

// GetAttachments returns the file attachments.
func (m *MailMessage) GetAttachments() []*mail.Attachment { return m.attachments }

// GetRawAttachments returns the raw data attachments.
func (m *MailMessage) GetRawAttachments() []RawAttachment { return m.rawAttachments }

// Tag adds a tag for mail provider categorisation.
func (m *MailMessage) Tag(value string) *MailMessage {
	m.tags = append(m.tags, value)

	return m
}

// GetTags returns the tags.
func (m *MailMessage) GetTags() []string { return m.tags }

// Metadata sets a metadata key-value pair.
func (m *MailMessage) Metadata(key, value string) *MailMessage {
	m.metadata[key] = value

	return m
}

// GetMetadata returns the metadata map.
func (m *MailMessage) GetMetadata() map[string]string { return m.metadata }

// Priority sets the email priority (1 = highest, 5 = lowest).
func (m *MailMessage) Priority(level int) *MailMessage {
	m.priority = level

	return m
}

// GetPriority returns the priority level.
func (m *MailMessage) GetPriority() int { return m.priority }

// WithSymfonyMessage registers a callback to configure the underlying mail
// message directly.
func (m *MailMessage) WithSymfonyMessage(callback func(*mail.Message)) *MailMessage {
	m.callbacks = append(m.callbacks, callback)

	return m
}

// GetCallbacks returns the message callbacks.
func (m *MailMessage) GetCallbacks() []func(*mail.Message) { return m.callbacks }

// Data returns the combined data for template rendering.
func (m *MailMessage) Data() map[string]any {
	data := m.ToMap()

	for k, v := range m.viewData {
		data[k] = v
	}

	return data
}

// The following methods shadow SimpleMessage methods so that MailMessage's
// fluent API returns *MailMessage instead of *SimpleMessage.

// Success sets the message level to success.
func (m *MailMessage) Success() *MailMessage {
	m.SimpleMessage.Success()

	return m
}

// Error sets the message level to error.
func (m *MailMessage) Error() *MailMessage {
	m.SimpleMessage.Error()

	return m
}

// Level sets the message level.
func (m *MailMessage) Level(level string) *MailMessage {
	m.SimpleMessage.Level(level)

	return m
}

// Subject sets the message subject.
func (m *MailMessage) Subject(subject string) *MailMessage {
	m.SimpleMessage.Subject(subject)

	return m
}

// Greeting sets the greeting line.
func (m *MailMessage) Greeting(greeting string) *MailMessage {
	m.SimpleMessage.Greeting(greeting)

	return m
}

// Salutation sets the closing salutation.
func (m *MailMessage) Salutation(salutation string) *MailMessage {
	m.SimpleMessage.Salutation(salutation)

	return m
}

// Line appends a line to the message body.
func (m *MailMessage) Line(line string) *MailMessage {
	m.SimpleMessage.Line(line)

	return m
}

// LineIf appends a line only if the condition is true.
func (m *MailMessage) LineIf(condition bool, line string) *MailMessage {
	m.SimpleMessage.LineIf(condition, line)

	return m
}

// Lines appends multiple lines to the message body.
func (m *MailMessage) Lines(lines []string) *MailMessage {
	m.SimpleMessage.Lines(lines)

	return m
}

// LinesIf appends multiple lines only if the condition is true.
func (m *MailMessage) LinesIf(condition bool, lines []string) *MailMessage {
	m.SimpleMessage.LinesIf(condition, lines)

	return m
}

// With is an alias for Line.
func (m *MailMessage) With(line string) *MailMessage {
	m.SimpleMessage.With(line)

	return m
}

// Action sets the call-to-action button text and URL.
func (m *MailMessage) Action(text, url string) *MailMessage {
	m.SimpleMessage.Action(text, url)

	return m
}

// Mailer sets the mailer name to use for delivery.
func (m *MailMessage) Mailer(name string) *MailMessage {
	m.SimpleMessage.Mailer(name)

	return m
}
