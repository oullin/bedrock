package mail

// SentMessage wraps the result of a successfully sent email.
type SentMessage struct {
	// Original is the Message that was sent.
	Original *Message
	// MessageID is the Message-ID assigned by the transport.
	MessageID string
}
