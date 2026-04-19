package prompts

// ImagePrompt carries all parameters for an image generation request.
type ImagePrompt struct {
	Prompt      string
	Attachments []any
	Size        *string
	Quality     *string
	Provider    *string
	Model       *string
	Timeout     int // default 30
}
