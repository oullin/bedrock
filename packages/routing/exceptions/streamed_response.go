package exceptions

// StreamedResponseException wraps an error encountered while writing a
// streamed response back to the client.
//
// Mirrors @bedrock\Routing\Exceptions\StreamedResponseException.
type StreamedResponseException struct{ Inner error }

func (e *StreamedResponseException) Error() string {
	if e.Inner == nil {
		return "streamed response failed"
	}

	return "streamed response failed: " + e.Inner.Error()
}

func (e *StreamedResponseException) Unwrap() error { return e.Inner }
