package data

import (
	"fmt"
	"math/rand/v2"
)

// GeneratedImage holds a single generated image.
// Mirrors Laravel\Ai\Responses\Data\GeneratedImage.
type GeneratedImage struct {
	Image    string `json:"image"`    // base64-encoded or URL
	MimeType string `json:"mime_type"` // default: image/png
}

// Content returns the raw image data.
func (g GeneratedImage) Content() string { return g.Image }

// RandomStorageName returns a random filename with the correct extension.
func (g GeneratedImage) RandomStorageName() string {
	ext := "png"
	switch g.MimeType {
	case "image/jpeg":
		ext = "jpg"
	case "image/webp":
		ext = "webp"
	case "image/gif":
		ext = "gif"
	}
	return fmt.Sprintf("%d.%s", rand.Int64(), ext)
}

// String returns the image data.
func (g GeneratedImage) String() string { return g.Image }

// ToMap returns a map representation.
func (g GeneratedImage) ToMap() map[string]any {
	return map[string]any{"image": g.Image, "mime_type": g.MimeType}
}
