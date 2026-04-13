package mailx

import cmail "github.com/bedrock/packages/contracts/mail"

// TextMessage wraps a Message for plain-text email rendering. Embed
// and EmbedData are stubs that return empty strings because inline
// images are meaningless in plain text.
type TextMessage struct {
	*cmail.Message
}

// NewTextMessage creates a TextMessage wrapping the given Message.
func NewTextMessage(m *cmail.Message) *TextMessage {
	return &TextMessage{Message: m}
}

// Embed is a no-op for text messages and returns an empty string.
func (t *TextMessage) Embed(_ string, _ ...cmail.AttachOption) string {
	return ""
}

// EmbedData is a no-op for text messages and returns an empty string.
func (t *TextMessage) EmbedData(_ []byte, _ string, _ ...string) string {
	return ""
}
