package messages

import (
	"github.com/bedrock/packages/ai/sdk/data"
	"github.com/bedrock/packages/ai/sdk/enums"
)

// ToolResultMessage carries the results of tool invocations back to the LLM.
// Mirrors Upstream\Ai\Messages\ToolResultMessage.
type ToolResultMessage struct {
	Message
	ToolResults []data.ToolResult
}

// NewToolResultMessage constructs a ToolResultMessage.
func NewToolResultMessage(results ...data.ToolResult) *ToolResultMessage {
	return &ToolResultMessage{
		Message:     Message{Role: enums.RoleToolResult},
		ToolResults: results,
	}
}
