package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/bedrock/packages/boost/mcp/tools"
)

const defaultTimeout = 180 * time.Second

// ToolExecutor dispatches MCP tool requests to the appropriate handler.
// Mirrors Upstream\Boost\Mcp\ToolExecutor.
type ToolExecutor struct {
	registry *ToolRegistry
	timeout  time.Duration
}

// NewExecutor returns a ToolExecutor backed by registry with the default 180 s
// timeout.
func NewExecutor(registry *ToolRegistry) *ToolExecutor {
	return &ToolExecutor{registry: registry, timeout: defaultTimeout}
}

// WithTimeout returns a new executor with the given timeout applied to every
// tool invocation.
func (e *ToolExecutor) WithTimeout(d time.Duration) *ToolExecutor {
	return &ToolExecutor{registry: e.registry, timeout: d}
}

// Execute runs the named tool with args, enforcing the configured timeout.
// Returns an error McpResponse (not a Go error) for well-known tool errors so
// that MCP clients receive structured responses.
func (e *ToolExecutor) Execute(name string, args map[string]any) (tools.McpResponse, error) {
	tool := e.registry.Find(name)
	if tool == nil {
		return tools.ErrorResponse(fmt.Sprintf("tool not found: %s", name)), nil
	}

	req := tools.McpRequest{Args: args}

	if e.timeout <= 0 {
		return tool.Handle(req)
	}

	type result struct {
		resp tools.McpResponse
		err  error
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	ch := make(chan result, 1)

	go func() {
		r, err := tool.Handle(req)
		ch <- result{r, err}
	}()

	select {
	case res := <-ch:
		return res.resp, res.err
	case <-ctx.Done():
		return tools.ErrorResponse(fmt.Sprintf("tool %s: execution timed out after %s", name, e.timeout)), nil
	}
}

// ExecuteReadOnly is like Execute but returns an error for non-read-only tools
// when readOnlyMode is true.
func (e *ToolExecutor) ExecuteReadOnly(name string, args map[string]any, readOnlyMode bool) (tools.McpResponse, error) {
	if readOnlyMode {
		tool := e.registry.Find(name)
		if tool != nil && !tool.IsReadOnly() {
			return tools.ErrorResponse(fmt.Sprintf("tool %s is not allowed in read-only mode", name)), nil
		}
	}

	return e.Execute(name, args)
}
