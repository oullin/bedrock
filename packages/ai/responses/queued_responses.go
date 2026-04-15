package responses

// QueuedAgentResponse represents an agent prompt that has been dispatched to a queue.
// Mirrors Laravel\Ai\Responses\QueuedAgentResponse.
type QueuedAgentResponse struct {
	InvocationID string
	thenFns      []func()
	catchFns     []func(error)
}

// NewQueuedAgentResponse constructs a QueuedAgentResponse.
func NewQueuedAgentResponse(invocationID string) *QueuedAgentResponse {
	return &QueuedAgentResponse{InvocationID: invocationID}
}

// GetInvocationID returns the unique ID for this queued invocation.
func (r *QueuedAgentResponse) GetInvocationID() string { return r.InvocationID }

// Then registers a callback to be called on successful completion.
func (r *QueuedAgentResponse) Then(fn func()) *QueuedAgentResponse {
	r.thenFns = append(r.thenFns, fn)
	return r
}

// Catch registers a callback to be called on failure.
func (r *QueuedAgentResponse) Catch(fn func(error)) *QueuedAgentResponse {
	r.catchFns = append(r.catchFns, fn)
	return r
}

// QueuedImageResponse represents a queued image generation request.
type QueuedImageResponse struct{ InvocationID string }

// QueuedAudioResponse represents a queued audio generation request.
type QueuedAudioResponse struct{ InvocationID string }

// QueuedEmbeddingsResponse represents a queued embeddings generation request.
type QueuedEmbeddingsResponse struct{ InvocationID string }

// QueuedTranscriptionResponse represents a queued transcription request.
type QueuedTranscriptionResponse struct{ InvocationID string }
