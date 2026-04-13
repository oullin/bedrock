package mailx

import (
	"io"
	"testing"

	cmail "github.com/bedrock/packages/contracts/mail"
)

// Mailable is a base struct that implements contracts/mail.Mailable.
// Embed it and override GetEnvelope/GetContent/GetAttachments/GetHeaders
// or use the fluent setters to compose an email.
type Mailable struct {
	envelope    cmail.Envelope
	content     cmail.Content
	attachments []*cmail.Attachment
	headers     cmail.Headers
	mailerName  string
	localeName  string
	callbacks   []func(*cmail.Message)
}

var _ cmail.Mailable = (*Mailable)(nil)

// GetEnvelope returns the envelope.
func (m *Mailable) GetEnvelope() *cmail.Envelope {
	return &m.envelope
}

// GetContent returns the content.
func (m *Mailable) GetContent() *cmail.Content {
	return &m.content
}

// GetAttachments returns all attachments.
func (m *Mailable) GetAttachments() []*cmail.Attachment {
	return m.attachments
}

// GetHeaders returns custom headers.
func (m *Mailable) GetHeaders() *cmail.Headers {
	return &m.headers
}

// GetMailerName returns the mailer name this mailable should be sent with.
func (m *Mailable) GetMailerName() string {
	return m.mailerName
}

// GetLocale returns the locale for this mailable.
func (m *Mailable) GetLocale() string {
	return m.localeName
}

// GetCallbacks returns the message callbacks.
func (m *Mailable) GetCallbacks() []func(*cmail.Message) {
	return m.callbacks
}

// --- Fluent setters ---

// SetFrom sets the sender address.
func (m *Mailable) SetFrom(address string, name ...string) *Mailable {
	m.envelope.From = cmail.Address{Email: address}

	if len(name) > 0 {
		m.envelope.From.Name = name[0]
	}

	return m
}

// SetTo sets the primary recipients, replacing any existing ones.
func (m *Mailable) SetTo(addresses ...cmail.Address) *Mailable {
	m.envelope.To = addresses

	return m
}

// AddTo appends primary recipients.
func (m *Mailable) AddTo(addresses ...cmail.Address) *Mailable {
	m.envelope.To = append(m.envelope.To, addresses...)

	return m
}

// SetCC sets the CC recipients, replacing any existing ones.
func (m *Mailable) SetCC(addresses ...cmail.Address) *Mailable {
	m.envelope.CC = addresses

	return m
}

// AddCC appends CC recipients.
func (m *Mailable) AddCC(addresses ...cmail.Address) *Mailable {
	m.envelope.CC = append(m.envelope.CC, addresses...)

	return m
}

// SetBCC sets the BCC recipients, replacing any existing ones.
func (m *Mailable) SetBCC(addresses ...cmail.Address) *Mailable {
	m.envelope.BCC = addresses

	return m
}

// AddBCC appends BCC recipients.
func (m *Mailable) AddBCC(addresses ...cmail.Address) *Mailable {
	m.envelope.BCC = append(m.envelope.BCC, addresses...)

	return m
}

// SetReplyTo sets the reply-to addresses.
func (m *Mailable) SetReplyTo(addresses ...cmail.Address) *Mailable {
	m.envelope.ReplyTo = addresses

	return m
}

// AddReplyTo appends reply-to addresses.
func (m *Mailable) AddReplyTo(addresses ...cmail.Address) *Mailable {
	m.envelope.ReplyTo = append(m.envelope.ReplyTo, addresses...)

	return m
}

// SetSubject sets the email subject.
func (m *Mailable) SetSubject(subject string) *Mailable {
	m.envelope.Subject = subject

	return m
}

// SetPriority sets the email priority (1 = highest, 5 = lowest).
func (m *Mailable) SetPriority(level int) *Mailable {
	m.callbacks = append(m.callbacks, func(msg *cmail.Message) {
		msg.Priority(level)
	})

	return m
}

// Tag adds a tag to the mailable.
func (m *Mailable) Tag(value string) *Mailable {
	m.envelope.Tags = append(m.envelope.Tags, value)

	return m
}

// SetMetadata sets a metadata key-value pair.
func (m *Mailable) SetMetadata(key, value string) *Mailable {
	if m.envelope.Metadata == nil {
		m.envelope.Metadata = make(map[string]string)
	}

	m.envelope.Metadata[key] = value

	return m
}

// SetMailer specifies which mailer should send this mailable.
func (m *Mailable) SetMailer(name string) *Mailable {
	m.mailerName = name

	return m
}

// SetLocale sets the locale for template rendering.
func (m *Mailable) SetLocale(locale string) *Mailable {
	m.localeName = locale

	return m
}

// SetHTML sets the HTML content as a literal string.
func (m *Mailable) SetHTML(html string) *Mailable {
	m.content.HTML = html
	m.content.HTMLString = true

	return m
}

// SetHTMLView sets the HTML template name.
func (m *Mailable) SetHTMLView(view string) *Mailable {
	m.content.HTML = view
	m.content.HTMLString = false

	return m
}

// SetText sets the plain-text template name.
func (m *Mailable) SetText(text string) *Mailable {
	m.content.Text = text

	return m
}

// SetMarkdown sets the Markdown content.
func (m *Mailable) SetMarkdown(markdown string) *Mailable {
	m.content.Markdown = markdown

	return m
}

// With adds data for template rendering.
func (m *Mailable) With(key string, value any) *Mailable {
	if m.content.With == nil {
		m.content.With = make(map[string]any)
	}

	m.content.With[key] = value

	return m
}

// WithData sets all template data at once.
func (m *Mailable) WithData(data map[string]any) *Mailable {
	m.content.With = data

	return m
}

// Attach adds a file attachment.
func (m *Mailable) Attach(attachment *cmail.Attachment) *Mailable {
	for _, existing := range m.attachments {
		if existing.IsEquivalent(attachment) {
			return m
		}
	}

	m.attachments = append(m.attachments, attachment)

	return m
}

// AttachMany adds multiple attachments.
func (m *Mailable) AttachMany(attachments []*cmail.Attachment) *Mailable {
	for _, a := range attachments {
		m.Attach(a)
	}

	return m
}

// AttachFromPath attaches a file by path.
func (m *Mailable) AttachFromPath(path string, options ...cmail.AttachOption) *Mailable {
	a := &cmail.Attachment{Path: path}

	for _, opt := range options {
		opt(a)
	}

	return m.Attach(a)
}

// AttachData attaches raw data with a filename.
func (m *Mailable) AttachData(data func() (io.Reader, error), name string, options ...cmail.AttachOption) *Mailable {
	a := cmail.FromData(data, name)

	for _, opt := range options {
		opt(a)
	}

	return m.Attach(a)
}

// AttachFromStorage attaches a file from the default storage disk.
func (m *Mailable) AttachFromStorage(path string, options ...cmail.AttachOption) *Mailable {
	a := cmail.FromStorage(path)

	for _, opt := range options {
		opt(a)
	}

	return m.Attach(a)
}

// AttachFromStorageDisk attaches a file from the given storage disk.
func (m *Mailable) AttachFromStorageDisk(disk, path string, options ...cmail.AttachOption) *Mailable {
	a := cmail.FromStorageDisk(disk, path)

	for _, opt := range options {
		opt(a)
	}

	return m.Attach(a)
}

// WithCallback registers a callback to customize the Message before sending.
func (m *Mailable) WithCallback(callback func(*cmail.Message)) *Mailable {
	m.callbacks = append(m.callbacks, callback)

	return m
}

// SetMessageID sets the Message-ID header.
func (m *Mailable) SetMessageID(id string) *Mailable {
	m.headers.MessageID = id

	return m
}

// SetReferences sets the References header.
func (m *Mailable) SetReferences(refs []string) *Mailable {
	m.headers.References = refs

	return m
}

// AddHeader adds a custom text header.
func (m *Mailable) AddHeader(name, value string) *Mailable {
	m.headers.Text = append(m.headers.Text, cmail.TextHeader{
		Name:  name,
		Value: value,
	})

	return m
}

// Tap calls the given function with the mailable for inline configuration.
func (m *Mailable) Tap(fn func(*Mailable)) *Mailable {
	fn(m)

	return m
}

// --- Has methods ---

// HasFrom reports whether the mailable is from the given address.
func (m *Mailable) HasFrom(address string, name ...string) bool {
	if m.envelope.From.Email != address {
		return false
	}

	if len(name) > 0 && m.envelope.From.Name != name[0] {
		return false
	}

	return true
}

// HasTo reports whether the given address is in the To recipients.
func (m *Mailable) HasTo(address string, name ...string) bool {
	return hasAddress(m.envelope.To, address, name)
}

// HasCC reports whether the given address is in the CC recipients.
func (m *Mailable) HasCC(address string, name ...string) bool {
	return hasAddress(m.envelope.CC, address, name)
}

// HasBCC reports whether the given address is in the BCC recipients.
func (m *Mailable) HasBCC(address string, name ...string) bool {
	return hasAddress(m.envelope.BCC, address, name)
}

// HasReplyTo reports whether the given address is in the ReplyTo list.
func (m *Mailable) HasReplyTo(address string, name ...string) bool {
	return hasAddress(m.envelope.ReplyTo, address, name)
}

// HasSubject reports whether the subject matches.
func (m *Mailable) HasSubject(subject string) bool {
	return m.envelope.Subject == subject
}

// HasTag reports whether the mailable has the given tag.
func (m *Mailable) HasTag(tag string) bool {
	for _, t := range m.envelope.Tags {
		if t == tag {
			return true
		}
	}

	return false
}

// HasMetadata reports whether the mailable has the given metadata entry.
func (m *Mailable) HasMetadata(key, value string) bool {
	v, ok := m.envelope.Metadata[key]

	return ok && v == value
}

// HasAttachment reports whether the mailable has an equivalent attachment.
func (m *Mailable) HasAttachment(other *cmail.Attachment) bool {
	for _, a := range m.attachments {
		if a.IsEquivalent(other) {
			return true
		}
	}

	return false
}

// HasAttachmentFromPath reports whether a path-based attachment exists.
func (m *Mailable) HasAttachmentFromPath(path string, options ...cmail.AttachOption) bool {
	check := &cmail.Attachment{Path: path}

	for _, opt := range options {
		opt(check)
	}

	return m.HasAttachment(check)
}

// HasAttachmentFromStorage reports whether a storage-based attachment exists.
func (m *Mailable) HasAttachmentFromStorage(path string, options ...cmail.AttachOption) bool {
	check := cmail.FromStorage(path)

	for _, opt := range options {
		opt(check)
	}

	return m.HasAttachment(check)
}

// HasAttachmentFromStorageDisk reports whether a disk-based attachment exists.
func (m *Mailable) HasAttachmentFromStorageDisk(disk, path string, options ...cmail.AttachOption) bool {
	check := cmail.FromStorageDisk(disk, path)

	for _, opt := range options {
		opt(check)
	}

	return m.HasAttachment(check)
}

// --- Assert methods ---

// AssertFrom asserts the sender matches.
func (m *Mailable) AssertFrom(t testing.TB, address string, name ...string) {
	t.Helper()

	if !m.HasFrom(address, name...) {
		t.Errorf("expected from %q, got %q", address, m.envelope.From.Email)
	}
}

// AssertTo asserts the given address is in To.
func (m *Mailable) AssertTo(t testing.TB, address string, name ...string) {
	t.Helper()

	if !m.HasTo(address, name...) {
		t.Errorf("expected to contain %q in To recipients", address)
	}
}

// AssertHasCC asserts the given address is in CC.
func (m *Mailable) AssertHasCC(t testing.TB, address string, name ...string) {
	t.Helper()

	if !m.HasCC(address, name...) {
		t.Errorf("expected to contain %q in CC recipients", address)
	}
}

// AssertHasBCC asserts the given address is in BCC.
func (m *Mailable) AssertHasBCC(t testing.TB, address string, name ...string) {
	t.Helper()

	if !m.HasBCC(address, name...) {
		t.Errorf("expected to contain %q in BCC recipients", address)
	}
}

// AssertHasReplyTo asserts the given address is in ReplyTo.
func (m *Mailable) AssertHasReplyTo(t testing.TB, address string, name ...string) {
	t.Helper()

	if !m.HasReplyTo(address, name...) {
		t.Errorf("expected to contain %q in ReplyTo", address)
	}
}

// AssertHasSubject asserts the subject matches.
func (m *Mailable) AssertHasSubject(t testing.TB, subject string) {
	t.Helper()

	if !m.HasSubject(subject) {
		t.Errorf("expected subject %q, got %q", subject, m.envelope.Subject)
	}
}

// AssertHasTag asserts the tag exists.
func (m *Mailable) AssertHasTag(t testing.TB, tag string) {
	t.Helper()

	if !m.HasTag(tag) {
		t.Errorf("expected tag %q", tag)
	}
}

// AssertHasMetadata asserts the metadata entry exists.
func (m *Mailable) AssertHasMetadata(t testing.TB, key, value string) {
	t.Helper()

	if !m.HasMetadata(key, value) {
		t.Errorf("expected metadata %q=%q", key, value)
	}
}

// AssertHasAttachment asserts an equivalent attachment exists.
func (m *Mailable) AssertHasAttachment(t testing.TB, attachment *cmail.Attachment) {
	t.Helper()

	if !m.HasAttachment(attachment) {
		t.Errorf("expected attachment %q", attachment.Name)
	}
}

// AssertHasNoAttachments asserts there are no attachments.
func (m *Mailable) AssertHasNoAttachments(t testing.TB) {
	t.Helper()

	if len(m.attachments) > 0 {
		t.Errorf("expected no attachments, got %d", len(m.attachments))
	}
}

func hasAddress(list []cmail.Address, email string, name []string) bool {
	for _, a := range list {
		if a.Email != email {
			continue
		}

		if len(name) > 0 && a.Name != name[0] {
			continue
		}

		return true
	}

	return false
}
