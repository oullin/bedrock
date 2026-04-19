package prompts

// AudioPrompt carries all parameters for a text-to-speech request.
type AudioPrompt struct {
	Text         string
	Voice        string // default "default-female"
	Instructions *string
	Provider     *string
	Model        *string
	Timeout      int // default 30
}
