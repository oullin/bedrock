package ai

import "context"

// ToolRequest represents an invocation request from the LLM to a tool.
type ToolRequest struct {
	ID        string
	Name      string
	Arguments map[string]any
}

// Tool is the contract every callable AI tool must satisfy.
// Mirrors Laravel\Ai\Contracts\Tool.
type Tool interface {
	// Description returns a human-readable description for the LLM.
	Description() string
	// Handle executes the tool and returns its result.
	Handle(ctx context.Context, req *ToolRequest) (any, error)
	// Schema returns the JSON schema describing the tool's input.
	Schema(schema JsonSchema) map[string]any
}
