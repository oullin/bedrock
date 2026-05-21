package responses

import (
	"strings"

	"github.com/bedrock/packages/ai/sdk/data"
)

// ImageResponse holds the result of an image generation request.
type ImageResponse struct {
	Images []data.GeneratedImage
	Usage  data.Usage
	Meta   data.Meta
}

// NewImageResponse constructs an ImageResponse.
func NewImageResponse(images []data.GeneratedImage, usage data.Usage, meta data.Meta) *ImageResponse {
	return &ImageResponse{Images: images, Usage: usage, Meta: meta}
}

// FirstImage returns the first generated image, or an empty one if none.
func (r *ImageResponse) FirstImage() data.GeneratedImage {
	if len(r.Images) == 0 {
		return data.GeneratedImage{}
	}

	return r.Images[0]
}

// Count returns the number of generated images.
func (r *ImageResponse) Count() int { return len(r.Images) }

// ToHTML returns an HTML string of <img> tags for all images.
func (r *ImageResponse) ToHTML() string {
	var sb strings.Builder

	for _, img := range r.Images {
		sb.WriteString(`<img src="data:`)
		sb.WriteString(img.MimeType)
		sb.WriteString(`;base64,`)
		sb.WriteString(img.Image)
		sb.WriteString(`" />`)
	}

	return sb.String()
}

// String returns the ToHTML representation.
func (r *ImageResponse) String() string { return r.ToHTML() }
