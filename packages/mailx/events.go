package mailx

import cmail "github.com/bedrock/packages/contracts/mail"

// MessageSending is dispatched before a message is sent. Listeners
// may inspect or modify the message before delivery.
type MessageSending struct {
	Mailer  string
	Message *cmail.Message
	Data    map[string]any
}

// MessageSent is dispatched after a message has been successfully sent.
type MessageSent struct {
	Mailer      string
	SentMessage *cmail.SentMessage
	Data        map[string]any
}
