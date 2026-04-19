package messages

import "github.com/bedrock/packages/ai/sdk/enums"

// Attachment represents a file or data attachment added to a user message.
type Attachment struct {
	Content  []byte
	MimeType string
	Name     string
}

// UserMessage is a message from the human participant.
// Mirrors Laravel\Ai\Messages\UserMessage.
type UserMessage struct {
	Message
	Attachments []Attachment
}

// NewUserMessage constructs a UserMessage.
func NewUserMessage(content string, attachments ...Attachment) *UserMessage {
	return &UserMessage{
		Message:     Message{Role: enums.RoleUser, Content: &content},
		Attachments: attachments,
	}
}
