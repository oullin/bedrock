package provider

import (
	"context"

	"github.com/bedrock/packages/contracts/ai/gateway"
)

// ImageGenerateRequest carries parameters for image generation at the provider level.
type ImageGenerateRequest struct {
	Prompt      string
	Attachments []any
	Size        *string
	Quality     *string
	Model       *string
	Timeout     int
}

// ImageProvider is the provider-level contract for image generation.
// Mirrors Upstream\Ai\Contracts\Providers\ImageProvider.
type ImageProvider interface {
	Image(ctx context.Context, req ImageGenerateRequest) (*gateway.ImageGenerateResult, error)
	ImageGateway() gateway.ImageGateway
	UseImageGateway(gw gateway.ImageGateway) ImageProvider
	DefaultImageModel() string
}
