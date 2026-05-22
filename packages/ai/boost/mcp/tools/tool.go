// Package tools provides the nine built-in MCP server tool implementations.
package tools

// McpRequest carries the arguments passed to an MCP tool invocation.
// Named McpRequest to avoid collision with http.Request and other request types.
type McpRequest struct {
	Args map[string]any
}

// Content is one item within an MCP tool response.
type Content struct {
	// Type is "text" or "json".
	Type string
	// Text holds string payloads when Type is "text".
	Text string
	// Data holds structured payloads when Type is "json".
	Data any
}

// McpResponse is the structured result of an MCP tool invocation.
type McpResponse struct {
	Content []Content
	IsError bool
}

// TextContent constructs a text Content item.

// JSONContent constructs a JSON Content item.

// McpTool is the interface every MCP server tool must satisfy.
// Named McpTool to avoid collision with contracts/ai.Tool (LLM function-calling tools).
type McpTool interface {
	// Name returns the tool's canonical identifier.
	Name() string
	// Description returns a human-readable description for AI agents.
	Description() string
	// Schema returns a JSON-schema-like map describing the tool's input parameters.
	Schema() map[string]any
	// Handle executes the tool and returns the response.
	Handle(req McpRequest) (McpResponse, error)
	// IsReadOnly reports whether this tool has side effects.
	IsReadOnly() bool
}

func TextContent(text string) Content { return Content{Type: "text", Text: text} }

func JSONContent(data any) Content { return Content{Type: "json", Data: data} }

// OkResponse constructs a successful McpResponse with a single JSON content item.
func OkResponse(data any) McpResponse {
	return McpResponse{Content: []Content{JSONContent(data)}}
}

// TextResponse constructs a successful McpResponse with a single text content item.
func TextResponse(text string) McpResponse {
	return McpResponse{Content: []Content{TextContent(text)}}
}

// ErrorResponse constructs an error McpResponse.
func ErrorResponse(text string) McpResponse {
	return McpResponse{Content: []Content{TextContent(text)}, IsError: true}
}
