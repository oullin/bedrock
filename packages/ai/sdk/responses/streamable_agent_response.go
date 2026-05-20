package responses

import (
	"github.com/bedrock/packages/ai/sdk/data"
	"github.com/bedrock/packages/ai/sdk/stream"
)

// StreamableAgentResponse wraps a lazy stream of events.
// Mirrors upstream Ai\Responses\StreamableAgentResponse.
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

// GetInvocationID returns the unique ID for this invocation.

// Each iterates over streaming events, calling fn for each one.
// Iteration stops when fn returns false.

// Then registers a callback to be called with the final StreamedAgentResponse
// after the stream is fully consumed.

// WithinConversation attaches conversation context.

// UsingVercelDataProtocol marks this stream to be serialised using the
// Vercel AI SDK data protocol format.

// Consume drains the stream, collects text, and calls Then callbacks.

// StreamedAgentResponse is the fully-consumed result of a streamed generation.
// Mirrors upstream Ai\Responses\StreamedAgentResponse.
type StreamedAgentResponse struct {
	AgentResponse
}

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

func (r *StreamableAgentResponse) GetInvocationID() string { return r.InvocationID }

func (r *StreamableAgentResponse) Each(fn func(stream.Event) bool) {
	r.events(fn)
}

func (r *StreamableAgentResponse) Then(fn func(*StreamedAgentResponse)) *StreamableAgentResponse {
	r.thenFns = append(r.thenFns, fn)

	return r
}

func (r *StreamableAgentResponse) WithinConversation(id string, user any) *StreamableAgentResponse {
	r.ConversationID = &id
	r.ConversationUser = user

	return r
}

func (r *StreamableAgentResponse) UsingVercelDataProtocol() *StreamableAgentResponse {
	r.useVercel = true

	return r
}

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
