package mcp_test

import (
	"testing"
	"time"

	"github.com/bedrock/packages/boost/mcp"
	"github.com/bedrock/packages/boost/mcp/tools"
)

// TestExecutorCallsHandler verifies Execute routes to the correct tool.

// Register a known-good stub that echoes the "value" arg.

// TestExecutorUnknownToolReturnsErrorResponse ensures unknown tools are handled gracefully.

// TestExecutorTimeout ensures the executor respects the configured timeout.

// TestExecutorReadOnlyMode blocks non-read-only tools.

// ---- Test tools -----------------------------------------------------------

type echoTool struct{}

type slowTool struct{ delay time.Duration }

type writableTool struct{}

func TestExecutorCallsHandler(t *testing.T) {
	t.Parallel()

	r := mcp.NewRegistry()

	r.Register(&echoTool{})

	e := mcp.NewExecutor(r)
	resp, err := e.Execute("echo_tool", map[string]any{"value": "hello"})

	if err != nil {
		t.Fatalf("Execute: unexpected error: %v", err)
	}

	if resp.IsError {
		t.Errorf("Execute returned error response: %v", resp.Content)
	}
}

func TestExecutorUnknownToolReturnsErrorResponse(t *testing.T) {
	t.Parallel()

	r := mcp.NewRegistry()
	e := mcp.NewExecutor(r)
	resp, err := e.Execute("no_such_tool", nil)

	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}

	if !resp.IsError {
		t.Error("expected error response for unknown tool")
	}
}

func TestExecutorTimeout(t *testing.T) {
	t.Parallel()

	r := mcp.NewRegistry()
	r.Register(&slowTool{delay: 200 * time.Millisecond})

	e := mcp.NewExecutor(r).WithTimeout(10 * time.Millisecond)
	resp, err := e.Execute("slow_tool", nil)

	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}

	if !resp.IsError {
		t.Error("expected timeout error response")
	}
}

func TestExecutorReadOnlyMode(t *testing.T) {
	t.Parallel()

	r := mcp.NewRegistry()
	r.Register(&writableTool{})

	e := mcp.NewExecutor(r)
	resp, err := e.ExecuteReadOnly("writable_tool", nil, true)

	if err != nil {
		t.Fatalf("unexpected Go error: %v", err)
	}

	if !resp.IsError {
		t.Error("expected error response in read-only mode for non-read-only tool")
	}
}

func (e *echoTool) Name() string           { return "echo_tool" }
func (e *echoTool) Description() string    { return "broadcastclient" }
func (e *echoTool) Schema() map[string]any { return nil }
func (e *echoTool) IsReadOnly() bool       { return true }
func (e *echoTool) Handle(req tools.McpRequest) (tools.McpResponse, error) {
	return tools.OkResponse(req.Args), nil
}

func (s *slowTool) Name() string           { return "slow_tool" }
func (s *slowTool) Description() string    { return "slow" }
func (s *slowTool) Schema() map[string]any { return nil }
func (s *slowTool) IsReadOnly() bool       { return true }
func (s *slowTool) Handle(_ tools.McpRequest) (tools.McpResponse, error) {
	time.Sleep(s.delay)

	return tools.TextResponse("done"), nil
}

func (w *writableTool) Name() string           { return "writable_tool" }
func (w *writableTool) Description() string    { return "writable" }
func (w *writableTool) Schema() map[string]any { return nil }
func (w *writableTool) IsReadOnly() bool       { return false }
func (w *writableTool) Handle(_ tools.McpRequest) (tools.McpResponse, error) {
	return tools.TextResponse("wrote something"), nil
}
