package mcp

// Response is returned by tool, resource, and prompt handlers. It wraps one
// or more Content items along with optional metadata.
type Response struct {
	contents       []Content
	role           string
	meta           map[string]any
	isError        bool
	isNotification bool
	notifMethod    string
	notifParams    map[string]any
	structured     map[string]any
}

// Text returns a Response containing a single TextContent.
func Text(text string) *Response {
	return &Response{contents: []Content{newTextContent(text, nil)}}
}

// Image returns a Response containing a single ImageContent.
// data must be a base64-encoded string and mimeType the image MIME type
// (e.g. "image/png").
func Image(data, mimeType string) *Response {
	return &Response{contents: []Content{&ImageContent{Data: data, MIMEType: mimeType}}}
}

// Audio returns a Response containing a single AudioContent.
// data must be a base64-encoded string and mimeType the audio MIME type
// (e.g. "audio/wav").
func Audio(data, mimeType string) *Response {
	return &Response{contents: []Content{&AudioContent{Data: data, MIMEType: mimeType}}}
}

// Blob returns a Response containing raw binary data for resource responses.
func Blob(data []byte, mimeType string) *Response {
	return &Response{contents: []Content{&BlobContent{Blob: data, MIMEType: mimeType}}}
}

// Error returns a Response that represents an error result for tools/call.
// The optional errs map is merged into the structured content.
func Error(message string, errs ...map[string]any) *Response {
	r := &Response{
		contents: []Content{newTextContent(message, nil)},
		isError:  true,
	}

	if len(errs) > 0 {
		r.structured = errs[0]
	}

	return r
}

// Notification returns a server-side notification (not a tool result). method
// is the JSON-RPC notification method name and params is optional.
func Notification(method string, params ...map[string]any) *Response {
	r := &Response{isNotification: true, notifMethod: method}

	if len(params) > 0 {
		r.notifParams = params[0]
	}

	return r
}

// WithMeta attaches an arbitrary metadata key-value pair to the response.
func (r *Response) WithMeta(key string, value any) *Response {
	if r.meta == nil {
		r.meta = make(map[string]any)
	}

	r.meta[key] = value

	return r
}

// Structured attaches structured output data to the response (tools only).
func (r *Response) Structured(data map[string]any) *Response {
	r.structured = data

	return r
}

// AsAssistant marks the response role as "assistant" (prompts only).
func (r *Response) AsAssistant() *Response {
	r.role = "assistant"

	return r
}

// IsError reports whether the response is an error result.
func (r *Response) IsError() bool { return r.isError }

// IsNotification reports whether the response is a server notification.
func (r *Response) IsNotification() bool { return r.isNotification }

// Contents returns the list of content items.
func (r *Response) Contents() []Content { return r.contents }

// Role returns the message role ("user" by default, "assistant" if
// AsAssistant() was called).
func (r *Response) Role() string {
	if r.role != "" {
		return r.role
	}

	return "user"
}

// toToolResult serialises the response for a tools/call result.
func (r *Response) toToolResult() map[string]any {
	content := make([]map[string]any, 0, len(r.contents))

	for _, c := range r.contents {
		content = append(content, c.ToTool())
	}

	result := map[string]any{"content": content, "isError": r.isError}

	if r.structured != nil {
		result["structuredContent"] = r.structured
	}

	if r.meta != nil {
		result["_meta"] = r.meta
	}

	return result
}

// toResourceResult serialises the response for a resources/read result.
func (r *Response) toResourceResult(uri string) []map[string]any {
	out := make([]map[string]any, 0, len(r.contents))

	for _, c := range r.contents {
		out = append(out, c.ToResource(uri))
	}

	return out
}

// toPromptMessages serialises the response as a prompt message.
func (r *Response) toPromptMessages() []map[string]any {
	msgs := make([]map[string]any, 0, len(r.contents))

	for _, c := range r.contents {
		msgs = append(msgs, map[string]any{
			"role":    r.Role(),
			"content": c.ToPrompt(),
		})
	}

	return msgs
}
