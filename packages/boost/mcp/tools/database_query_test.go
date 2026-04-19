package tools_test

import (
	"testing"

	"github.com/bedrock/packages/boost/mcp/tools"
)

func TestDatabaseQueryNoConnection(t *testing.T) {
	t.Parallel()

	tool := &tools.DatabaseQuery{}
	resp, err := tool.Handle(tools.McpRequest{Args: map[string]any{"query": "SELECT 1"}})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if !resp.IsError {
		t.Error("expected error response when no DB connection")
	}
}

func TestDatabaseQueryMissingQueryArg(t *testing.T) {
	t.Parallel()

	tool := &tools.DatabaseQuery{}
	resp, err := tool.Handle(tools.McpRequest{Args: map[string]any{}})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if !resp.IsError {
		t.Error("expected error response when query arg missing")
	}
}

func TestDatabaseQuerySchema(t *testing.T) {
	t.Parallel()

	tool := &tools.DatabaseQuery{}
	schema := tool.Schema()

	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want \"object\"", schema["type"])
	}

	props, ok := schema["properties"].(map[string]any)

	if !ok {
		t.Fatal("schema missing properties map")
	}

	if _, hasQuery := props["query"]; !hasQuery {
		t.Error("schema missing 'query' property")
	}
}

func TestDatabaseQueryIsReadOnly(t *testing.T) {
	t.Parallel()

	if !(&tools.DatabaseQuery{}).IsReadOnly() {
		t.Error("DatabaseQuery should be read-only")
	}
}
