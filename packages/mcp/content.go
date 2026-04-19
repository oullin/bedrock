package mcp

import "encoding/base64"

// Content is a piece of content returned by a tool, resource, or prompt
// handler. The three serialisation methods produce the format required by
// each MCP endpoint.
type Content interface {
	// ToTool serialises the content for a tools/call response.
	ToTool() map[string]any
	// ToPrompt serialises the content for a prompts/get response.
	ToPrompt() map[string]any
	// ToResource serialises the content for a resources/read response.
	// uri is the resource URI included in each content item.
	ToResource(uri string) map[string]any
}

// TextContent holds plain-text content.
type TextContent struct {
	Text string
	meta map[string]any
}

// newTextContent creates a TextContent, optionally merging extra metadata.

// ToTool implements Content.

// ToPrompt implements Content.

// ToResource implements Content.

// ImageContent holds a base64-encoded image.
type ImageContent struct {
	Data     string
	MIMEType string
	meta     map[string]any
}

// ToTool implements Content.

// ToPrompt implements Content.

// ToResource implements Content.

// AudioContent holds base64-encoded audio.
type AudioContent struct {
	Data     string
	MIMEType string
	meta     map[string]any
}

// ToTool implements Content.

// ToPrompt implements Content.

// ToResource implements Content.

// BlobContent holds raw binary data for resource responses.
type BlobContent struct {
	Blob     []byte
	MIMEType string
}

func newTextContent(text string, meta map[string]any) *TextContent {
	return &TextContent{Text: text, meta: meta}
}

func (c *TextContent) ToTool() map[string]any {
	m := map[string]any{"type": "text", "text": c.Text}

	for k, v := range c.meta {
		m[k] = v
	}

	return m
}

func (c *TextContent) ToPrompt() map[string]any {
	return c.ToTool()
}

func (c *TextContent) ToResource(uri string) map[string]any {
	m := map[string]any{"uri": uri, "text": c.Text}

	if c.meta != nil {
		if mt, ok := c.meta["mimeType"]; ok {
			m["mimeType"] = mt
		}
	}

	return m
}

func (c *ImageContent) ToTool() map[string]any {
	m := map[string]any{"type": "image", "data": c.Data, "mimeType": c.MIMEType}

	for k, v := range c.meta {
		m[k] = v
	}

	return m
}

func (c *ImageContent) ToPrompt() map[string]any {
	return c.ToTool()
}

func (c *ImageContent) ToResource(uri string) map[string]any {
	return map[string]any{"uri": uri, "blob": c.Data, "mimeType": c.MIMEType}
}

func (c *AudioContent) ToTool() map[string]any {
	m := map[string]any{"type": "audio", "data": c.Data, "mimeType": c.MIMEType}

	for k, v := range c.meta {
		m[k] = v
	}

	return m
}

func (c *AudioContent) ToPrompt() map[string]any {
	return c.ToTool()
}

func (c *AudioContent) ToResource(uri string) map[string]any {
	return map[string]any{"uri": uri, "blob": c.Data, "mimeType": c.MIMEType}
}

// ToTool implements Content — binary is base64-encoded as image for tools.
func (c *BlobContent) ToTool() map[string]any {
	return map[string]any{
		"type":     "image",
		"data":     base64.StdEncoding.EncodeToString(c.Blob),
		"mimeType": c.MIMEType,
	}
}

// ToPrompt implements Content.
func (c *BlobContent) ToPrompt() map[string]any {
	return c.ToTool()
}

// ToResource implements Content.
func (c *BlobContent) ToResource(uri string) map[string]any {
	return map[string]any{
		"uri":      uri,
		"blob":     base64.StdEncoding.EncodeToString(c.Blob),
		"mimeType": c.MIMEType,
	}
}
