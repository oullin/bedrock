package mcp

import (
	"github.com/bedrock/packages/ai/boost/mcp/tools"
)

// Server is the Boost MCP server. It wires together a ToolRegistry and a
// ToolExecutor and provides the high-level API that agents call to list and
// run tools.
// Mirrors upstream Boost\Mcp\Server.
type Server struct {
	registry *ToolRegistry
	executor *ToolExecutor
	readOnly bool
}

// NewServer returns a Server backed by the 9 default tools.
func NewServer() *Server {
	r := NewRegistry()
	e := NewExecutor(r)

	return &Server{registry: r, executor: e}
}

// NewServerWithRegistry returns a Server using a caller-supplied registry.
// Useful in tests and when custom tools must be registered.
func NewServerWithRegistry(r *ToolRegistry) *Server {
	return &Server{registry: r, executor: NewExecutor(r)}
}

// SetReadOnly restricts the server to read-only tools.
func (s *Server) SetReadOnly(v bool) { s.readOnly = v }

// Register adds a tool to the server's registry.
func (s *Server) Register(t tools.McpTool) { s.registry.Register(t) }

// AvailableTools returns the list of tools currently visible on this server.
func (s *Server) AvailableTools() []tools.McpTool { return s.registry.GetAvailableTools() }

// Call executes the named tool with the given arguments.
func (s *Server) Call(name string, args map[string]any) (tools.McpResponse, error) {
	return s.executor.ExecuteReadOnly(name, args, s.readOnly)
}

// Registry exposes the underlying ToolRegistry for advanced configuration.
func (s *Server) Registry() *ToolRegistry { return s.registry }

// Executor exposes the underlying ToolExecutor for timeout configuration.
func (s *Server) Executor() *ToolExecutor { return s.executor }
