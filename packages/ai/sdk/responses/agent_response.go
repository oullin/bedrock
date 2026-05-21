package responses

import (
	"github.com/bedrock/packages/ai/sdk/data"
)

// AgentResponse extends TextResponse with conversation tracking metadata.
type AgentResponse struct {
	TextResponse
	InvocationID     string
	ConversationID   *string
	ConversationUser any
}

// NewAgentResponse constructs an AgentResponse.
func NewAgentResponse(invocationID string, text string, usage data.Usage, meta data.Meta) *AgentResponse {
	return &AgentResponse{
		TextResponse: TextResponse{
			Text:  text,
			Usage: usage,
			Meta:  meta,
		},
		InvocationID: invocationID,
	}
}

// GetText returns the response text (satisfies AgentResponder contract).
func (r *AgentResponse) GetText() string { return r.Text }

// GetInvocationID returns the unique ID for this invocation.
func (r *AgentResponse) GetInvocationID() string { return r.InvocationID }

// WithinConversation attaches conversation context to the response.
func (r *AgentResponse) WithinConversation(id string, user any) *AgentResponse {
	r.ConversationID = &id
	r.ConversationUser = user

	return r
}

// Then executes a callback with this response and returns it (fluent chaining).
func (r *AgentResponse) Then(fn func(*AgentResponse)) *AgentResponse {
	fn(r)

	return r
}
