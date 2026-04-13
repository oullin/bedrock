package mail

import "strings"

// TextHeader is a custom MIME header key-value pair.
type TextHeader struct {
	Name  string
	Value string
}

// Headers holds custom header entries for an email message.
type Headers struct {
	MessageID  string
	References []string
	Text       []TextHeader
}

// ReferencesString returns the References header value formatted as a
// space-separated list of message-IDs.
func (h *Headers) ReferencesString() string {
	return strings.Join(h.References, " ")
}
