package gateway

import "context"

// ImageGenerateRequest carries parameters for image generation.
type ImageGenerateRequest struct {
	Model       string
	Prompt      string
	Attachments []any
	Size        *string
	Quality     *string
	Timeout     int
}

// GeneratedImageData holds the raw image data from a generation.
type GeneratedImageData struct {
	Image    string // base64-encoded or URL depending on provider
	MimeType string
}

// ImageGenerateResult is the normalised response from an image gateway.
type ImageGenerateResult struct {
	Images []GeneratedImageData
	Usage  TokenUsage
	Meta   ResponseMeta
}

// ImageGateway drives image generation for a single provider.
type ImageGateway interface {
	GenerateImage(ctx context.Context, req ImageGenerateRequest) (*ImageGenerateResult, error)
}
