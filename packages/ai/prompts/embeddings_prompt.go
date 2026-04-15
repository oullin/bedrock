package prompts

// EmbeddingsPrompt carries all parameters for an embedding generation request.
type EmbeddingsPrompt struct {
	Inputs     []string
	Dimensions *int
	Provider   *string
	Model      *string
	Timeout    int // default 30
}
