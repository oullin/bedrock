package mcp

import "context"

// Tool is an MCP tool that AI agents can invoke. Implementations must provide
// a name, a human-readable description, a JSON Schema for the input
// parameters, and a Handle function.
type Tool interface {
	Name() string
	Description() string
	// Schema returns a JSON Schema object (type "object") describing the
	// accepted input parameters.
	Schema() map[string]any
	Handle(ctx context.Context, req *Request) (*Response, error)
}

// funcTool is an adapter that implements Tool from plain functions.
type funcTool struct {
	name        string
	description string
	schema      map[string]any
	handler     func(ctx context.Context, req *Request) (*Response, error)
}

func (t *funcTool) Name() string           { return t.name }
func (t *funcTool) Description() string    { return t.description }
func (t *funcTool) Schema() map[string]any { return t.schema }
func (t *funcTool) Handle(ctx context.Context, req *Request) (*Response, error) {
	return t.handler(ctx, req)
}

// NewTool creates a Tool from plain functions. schema is a JSON Schema object
// (may be nil for tools with no input parameters).
func NewTool(
	name, description string,
	schema map[string]any,
	handler func(ctx context.Context, req *Request) (*Response, error),
) Tool {
	if schema == nil {
		schema = map[string]any{"type": "object", "properties": map[string]any{}}
	}

	return &funcTool{name: name, description: description, schema: schema, handler: handler}
}

// toolToMap serialises a Tool for inclusion in a tools/list response.
func toolToMap(t Tool) map[string]any {
	return map[string]any{
		"name":        t.Name(),
		"description": t.Description(),
		"inputSchema": t.Schema(),
	}
}
