package messages

import (
	"github.com/bedrock/packages/ai/sdk/data"
	"github.com/bedrock/packages/ai/sdk/enums"
)

// AssistantMessage is a message from the AI assistant.
// Mirrors Laravel\Ai\Messages\AssistantMessage.
type AssistantMessage struct {
	Message
	ToolCalls []data.ToolCall
}

// NewAssistantMessage constructs an AssistantMessage.
func NewAssistantMessage(content string, toolCalls ...data.ToolCall) *AssistantMessage {
	return &AssistantMessage{
		Message:   Message{Role: enums.RoleAssistant, Content: &content},
		ToolCalls: toolCalls,
	}
}
