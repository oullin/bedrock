package mail

// Mailable describes an email message. It is a pure data carrier: the
// envelope defines recipients and metadata, content defines the body,
// attachments lists files, and headers provides custom header entries.
type Mailable interface {
	GetEnvelope() *Envelope
	GetContent() *Content
	GetAttachments() []*Attachment
	GetHeaders() *Headers
}
