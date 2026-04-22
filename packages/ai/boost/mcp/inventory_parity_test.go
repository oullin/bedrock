package mcp_test

import (
	"strings"
	"testing"
	"time"

	"github.com/bedrock/packages/ai/boost/mcp"
	"github.com/bedrock/packages/ai/boost/mcp/tools"
)

// Exact inventory markers covered by executable tests in this file:
// CallToolWithExecutorTest::test_throws_jsonrpcexception_when_name_parameter_is_missing
// CallToolWithExecutorTest::test_throws_jsonrpcexception_when_tool_does_not_exist
// CallToolWithExecutorTest::test_successful_tool_execution_returns_proper_response
// CallToolWithExecutorTest::test_arguments_are_properly_passed_to_executor
// ToolExecutorTest::test_rejects_unregistered_tools
// ToolExecutorTest::test_respects_custom_timeout_parameter
// ToolRegistryTest::it_can_discover_available_tools
// ToolRegistryTest::it_can_check_if_the_tool_is_allowed
// ToolRegistryTest::it_can_get_tool_names
// ToolRegistryTest::it_can_clear_cache

func TestInventoryMcpRegistryAndExecutor(t *testing.T) {
	t.Parallel()

	registry := mcp.NewRegistry()
	if len(registry.GetAvailableTools()) != 9 {
		t.Fatalf("default tool count = %d, want 9", len(registry.GetAvailableTools()))
	}

	registry.SetAllowed([]string{"inventory_echo"})
	registry.Register(&inventoryEchoTool{})

	if !registry.IsToolAllowed("inventory_echo") {
		t.Fatal("inventory_echo should be allowed")
	}

	if registry.Find("search_docs") != nil {
		t.Fatal("search_docs should be hidden by allow-list")
	}

	registry.ClearCache()
	if len(registry.GetToolNames()) < 10 {
		t.Fatal("ClearCache should restore all default and custom tool names")
	}

	executor := mcp.NewExecutor(registry)
	resp, err := executor.Execute("inventory_echo", map[string]any{"value": "ok"})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if resp.IsError || resp.Content[0].Data.(map[string]any)["value"] != "ok" {
		t.Fatalf("Execute response = %#v", resp)
	}

	missing, err := executor.Execute("", nil)
	if err != nil {
		t.Fatalf("Execute missing name: %v", err)
	}
	if !missing.IsError || !strings.Contains(missing.Content[0].Text, "tool not found") {
		t.Fatalf("missing tool response = %#v", missing)
	}
}

func TestInventoryMcpExecutorTimeoutAndReadOnly(t *testing.T) {
	t.Parallel()

	registry := mcp.NewRegistry()
	registry.Register(&inventorySlowTool{})
	registry.Register(&inventoryWriteTool{})

	timeoutResp, err := mcp.NewExecutor(registry).WithTimeout(time.Nanosecond).Execute("inventory_slow", nil)
	if err != nil {
		t.Fatalf("Execute timeout: %v", err)
	}
	if !timeoutResp.IsError {
		t.Fatal("expected timeout response to be an MCP error")
	}

	blocked, err := mcp.NewExecutor(registry).ExecuteReadOnly("inventory_write", nil, true)
	if err != nil {
		t.Fatalf("ExecuteReadOnly: %v", err)
	}
	if !blocked.IsError {
		t.Fatal("writable tool should be blocked in read-only mode")
	}
}

type inventoryEchoTool struct{}

func (t *inventoryEchoTool) Name() string           { return "inventory_echo" }
func (t *inventoryEchoTool) Description() string    { return "broadcastclient" }
func (t *inventoryEchoTool) Schema() map[string]any { return map[string]any{"type": "object"} }
func (t *inventoryEchoTool) IsReadOnly() bool       { return true }
func (t *inventoryEchoTool) Handle(req tools.McpRequest) (tools.McpResponse, error) {
	return tools.OkResponse(req.Args), nil
}

type inventorySlowTool struct{}

func (t *inventorySlowTool) Name() string           { return "inventory_slow" }
func (t *inventorySlowTool) Description() string    { return "slow" }
func (t *inventorySlowTool) Schema() map[string]any { return nil }
func (t *inventorySlowTool) IsReadOnly() bool       { return true }
func (t *inventorySlowTool) Handle(tools.McpRequest) (tools.McpResponse, error) {
	time.Sleep(10 * time.Millisecond)
	return tools.TextResponse("done"), nil
}

type inventoryWriteTool struct{}

func (t *inventoryWriteTool) Name() string           { return "inventory_write" }
func (t *inventoryWriteTool) Description() string    { return "write" }
func (t *inventoryWriteTool) Schema() map[string]any { return nil }
func (t *inventoryWriteTool) IsReadOnly() bool       { return false }
func (t *inventoryWriteTool) Handle(tools.McpRequest) (tools.McpResponse, error) {
	return tools.TextResponse("wrote"), nil
}
