package responses

import (
	"github.com/bedrock/packages/ai/data"
	"github.com/bedrock/packages/ai/stream"
)

// StreamableAgentResponse wraps a lazy stream of events.
// Mirrors Upstream\Ai\Responses\StreamableAgentResponse.
type StreamableAgentResponse struct {
	InvocationID     string
	events           func(yield func(stream.Event) bool)
	Usage            data.Usage
	Meta             data.Meta
	ConversationID   *string
	ConversationUser any
	thenFns          []func(*StreamedAgentResponse)
	useVercel        bool
}

// NewStreamableAgentResponse constructs a StreamableAgentResponse.
func NewStreamableAgentResponse(
	invocationID string,
	events func(yield func(stream.Event) bool),
	usage data.Usage,
	meta data.Meta,
) *StreamableAgentResponse {
	return &StreamableAgentResponse{
		InvocationID: invocationID,
		events:       events,
		Usage:        usage,
		Meta:         meta,
	}
}

// GetInvocationID returns the unique ID for this invocation.
func (r *StreamableAgentResponse) GetInvocationID() string { return r.InvocationID }

// Each iterates over streaming events, calling fn for each one.
// Iteration stops when fn returns false.
func (r *StreamableAgentResponse) Each(fn func(stream.Event) bool) {
	r.events(fn)
}

// Then registers a callback to be called with the final StreamedAgentResponse
// after the stream is fully consumed.
func (r *StreamableAgentResponse) Then(fn func(*StreamedAgentResponse)) *StreamableAgentResponse {
	r.thenFns = append(r.thenFns, fn)
	return r
}

// WithinConversation attaches conversation context.
func (r *StreamableAgentResponse) WithinConversation(id string, user any) *StreamableAgentResponse {
	r.ConversationID = &id
	r.ConversationUser = user
	return r
}

// UsingVercelDataProtocol marks this stream to be serialised using the
// Vercel AI SDK data protocol format.
func (r *StreamableAgentResponse) UsingVercelDataProtocol() *StreamableAgentResponse {
	r.useVercel = true
	return r
}

// Consume drains the stream, collects text, and calls Then callbacks.
func (r *StreamableAgentResponse) Consume() *StreamedAgentResponse {
	var text string
	r.events(func(e stream.Event) bool {
		if td, ok := e.(stream.TextDelta); ok {
			text += td.Delta
		}
		return true
	})
	streamed := &StreamedAgentResponse{
		AgentResponse: *NewAgentResponse(r.InvocationID, text, r.Usage, r.Meta),
	}
	if r.ConversationID != nil {
		streamed.ConversationID = r.ConversationID
		streamed.ConversationUser = r.ConversationUser
	}
	for _, fn := range r.thenFns {
		fn(streamed)
	}
	return streamed
}

// StreamedAgentResponse is the fully-consumed result of a streamed generation.
// Mirrors Upstream\Ai\Responses\StreamedAgentResponse.
type StreamedAgentResponse struct {
	AgentResponse
}
