// Package mcp provides the Boost MCP server, tool registry, and executor.
package mcp

import (
	"sync"

	"github.com/bedrock/packages/boost/mcp/tools"
)

// ToolRegistry holds the set of available MCP tools and manages per-tool
// allow-list enforcement. Mirrors Upstream\Boost\Mcp\ToolRegistry.
type ToolRegistry struct {
	mu         sync.RWMutex
	registered []tools.McpTool
	allowed    []string // nil = allow all
}

// NewRegistry returns a ToolRegistry pre-loaded with the 9 default tools.
// The caller can optionally pass a *config.Repository; if omitted the tools use
// sensible defaults (e.g. default log path, no DB connection).
func NewRegistry() *ToolRegistry {
	r := &ToolRegistry{}
	r.loadDefaults()

	return r
}

// loadDefaults registers the 9 built-in tools.
func (r *ToolRegistry) loadDefaults() {
	r.registered = []tools.McpTool{
		&tools.ApplicationInfo{},
		&tools.BrowserLogs{},
		&tools.DatabaseConnections{},
		&tools.DatabaseQuery{},
		&tools.DatabaseSchema{},
		&tools.GetAbsoluteUrl{},
		&tools.LastError{},
		&tools.ReadLogEntries{},
		&tools.SearchDocs{},
	}
}

// Register adds a tool to the registry. If a tool with the same name already
// exists it is replaced.
func (r *ToolRegistry) Register(t tools.McpTool) {
	r.mu.Lock()

	defer r.mu.Unlock()

	for i, existing := range r.registered {
		if existing.Name() == t.Name() {
			r.registered[i] = t

			return
		}
	}

	r.registered = append(r.registered, t)
}

// SetAllowed restricts which tools are visible to callers. Passing nil or an
// empty slice re-enables all tools.
func (r *ToolRegistry) SetAllowed(names []string) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.allowed = names
}

// IsToolAllowed reports whether the named tool is in the allow-list (or that no
// allow-list is set, in which case all tools are allowed).
func (r *ToolRegistry) IsToolAllowed(name string) bool {
	r.mu.RLock()

	defer r.mu.RUnlock()

	if len(r.allowed) == 0 {
		return true
	}

	for _, a := range r.allowed {
		if a == name {
			return true
		}
	}

	return false
}

// GetAvailableTools returns the subset of registered tools that pass the
// allow-list filter.
func (r *ToolRegistry) GetAvailableTools() []tools.McpTool {
	r.mu.RLock()

	defer r.mu.RUnlock()

	if len(r.allowed) == 0 {
		out := make([]tools.McpTool, len(r.registered))
		copy(out, r.registered)

		return out
	}

	var out []tools.McpTool

	for _, t := range r.registered {
		if r.isAllowed(t.Name()) {
			out = append(out, t)
		}
	}

	return out
}

// GetToolNames returns the names of all available tools.
func (r *ToolRegistry) GetToolNames() []string {
	available := r.GetAvailableTools()
	names := make([]string, len(available))

	for i, t := range available {
		names[i] = t.Name()
	}

	return names
}

// Find returns the tool with the given name, or nil if not found or not allowed.
func (r *ToolRegistry) Find(name string) tools.McpTool {
	r.mu.RLock()

	defer r.mu.RUnlock()

	if !r.isAllowed(name) {
		return nil
	}

	for _, t := range r.registered {
		if t.Name() == name {
			return t
		}
	}

	return nil
}

// ClearCache removes the allow-list, making all registered tools visible again.
func (r *ToolRegistry) ClearCache() {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.allowed = nil
}

// isAllowed is the internal (non-locking) helper.
func (r *ToolRegistry) isAllowed(name string) bool {
	if len(r.allowed) == 0 {
		return true
	}

	for _, a := range r.allowed {
		if a == name {
			return true
		}
	}

	return false
}
