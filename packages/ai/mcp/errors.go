package mcp

import "errors"

// Sentinel errors returned by the MCP server and its primitives.
var (
	// ErrToolNotFound is returned when tools/call names a tool that is not
	// registered on the server.
	ErrToolNotFound = errors.New("mcp: tool not found")

	// ErrResourceNotFound is returned when resources/read supplies a URI
	// that matches no registered resource or resource template.
	ErrResourceNotFound = errors.New("mcp: resource not found")

	// ErrPromptNotFound is returned when prompts/get names a prompt that is
	// not registered on the server.
	ErrPromptNotFound = errors.New("mcp: prompt not found")

	// ErrMethodNotFound is returned when the incoming JSON-RPC request
	// specifies a method that the server does not handle.
	ErrMethodNotFound = errors.New("mcp: method not found")

	// ErrInvalidRequest is returned when the incoming JSON-RPC message is
	// structurally invalid (e.g. missing the "method" field).
	ErrInvalidRequest = errors.New("mcp: invalid request")

	// ErrParseError is returned when the raw message is not valid JSON.
	ErrParseError = errors.New("mcp: parse error")
)
