package mcp_test

import (
	"testing"

	"github.com/bedrock/packages/boost/mcp"
	"github.com/bedrock/packages/boost/mcp/tools"
)

// TestRegistryDefaultTools verifies 9 tools are registered by default.

// TestRegistryGetToolNames verifies the tool names list.

// TestRegistryIsToolAllowedDefault verifies all tools are allowed by default.

// TestRegistryAllowList verifies the allow-list restricts visible tools.

// TestRegistryClearCache re-enables all tools after an allow-list is set.

// TestRegistryRegisterCustomTool verifies adding a custom tool.

// stubTool satisfies tools.McpTool for registry tests.
type stubTool struct{ name string }

func TestRegistryDefaultTools(t *testing.T) {
	t.Parallel()

	r := mcp.NewRegistry()
	got := r.GetAvailableTools()

	if len(got) != 9 {
		t.Errorf("GetAvailableTools() returned %d tools, want 9", len(got))
	}
}

func TestRegistryGetToolNames(t *testing.T) {
	t.Parallel()

	r := mcp.NewRegistry()
	names := r.GetToolNames()

	wantNames := []string{
		"application_info",
		"browser_logs",
		"database_connections",
		"database_query",
		"database_schema",
		"get_absolute_url",
		"last_error",
		"read_log_entries",
		"search_docs",
	}

	if len(names) != len(wantNames) {
		t.Fatalf("GetToolNames() len = %d, want %d", len(names), len(wantNames))
	}

	nameSet := make(map[string]bool, len(names))

	for _, n := range names {
		nameSet[n] = true
	}

	for _, want := range wantNames {
		if !nameSet[want] {
			t.Errorf("missing tool %q in GetToolNames()", want)
		}
	}
}

func TestRegistryIsToolAllowedDefault(t *testing.T) {
	t.Parallel()

	r := mcp.NewRegistry()

	if !r.IsToolAllowed("search_docs") {
		t.Error("IsToolAllowed(\"search_docs\") should be true by default")
	}
}

func TestRegistryAllowList(t *testing.T) {
	t.Parallel()

	r := mcp.NewRegistry()
	r.SetAllowed([]string{"search_docs", "last_error"})

	if !r.IsToolAllowed("search_docs") {
		t.Error("search_docs should be allowed")
	}

	if r.IsToolAllowed("database_query") {
		t.Error("database_query should not be allowed")
	}

	available := r.GetAvailableTools()

	if len(available) != 2 {
		t.Errorf("GetAvailableTools() with allow-list returned %d, want 2", len(available))
	}
}

func TestRegistryClearCache(t *testing.T) {
	t.Parallel()

	r := mcp.NewRegistry()
	r.SetAllowed([]string{"search_docs"})
	r.ClearCache()

	if got := r.GetAvailableTools(); len(got) != 9 {
		t.Errorf("after ClearCache, GetAvailableTools() = %d, want 9", len(got))
	}
}

func TestRegistryRegisterCustomTool(t *testing.T) {
	t.Parallel()

	r := mcp.NewRegistry()
	r.Register(&stubTool{name: "my_custom_tool"})

	if r.Find("my_custom_tool") == nil {
		t.Error("custom tool should be discoverable via Find")
	}
}

func (s *stubTool) Name() string           { return s.name }
func (s *stubTool) Description() string    { return "stub" }
func (s *stubTool) Schema() map[string]any { return nil }
func (s *stubTool) Handle(_ tools.McpRequest) (tools.McpResponse, error) {
	return tools.TextResponse("ok"), nil
}
func (s *stubTool) IsReadOnly() bool { return true }
