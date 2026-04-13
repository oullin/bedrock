package mail

// Attachable is implemented by types that can be converted to a mail
// attachment (e.g. an uploaded file or a storage object).
type Attachable interface {
	ToMailAttachment() (*Attachment, error)
}
