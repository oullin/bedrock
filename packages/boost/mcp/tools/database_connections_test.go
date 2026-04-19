package tools_test

import (
	"testing"

	"github.com/bedrock/packages/boost/mcp/tools"
)

func TestDatabaseConnectionsEmpty(t *testing.T) {
	t.Parallel()

	tool := &tools.DatabaseConnections{}
	resp, err := tool.Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if resp.IsError {
		t.Errorf("Handle returned error response: %v", resp.Content)
	}

	data, ok := resp.Content[0].Data.(map[string]any)

	if !ok {
		t.Fatalf("content data type = %T, want map[string]any", resp.Content[0].Data)
	}

	if data["default"] != "" {
		t.Errorf("default = %v, want empty string", data["default"])
	}
}

func TestDatabaseConnectionsWithConnections(t *testing.T) {
	t.Parallel()

	tool := &tools.DatabaseConnections{
		Connections:    map[string]string{"sqlite": "file::memory:", "postgres": "postgres://localhost/app"},
		DefaultConnect: "sqlite",
	}

	resp, err := tool.Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if resp.IsError {
		t.Errorf("Handle returned error: %v", resp.Content)
	}

	data, ok := resp.Content[0].Data.(map[string]any)

	if !ok {
		t.Fatalf("content data type = %T", resp.Content[0].Data)
	}

	if data["default"] != "sqlite" {
		t.Errorf("default = %v, want \"sqlite\"", data["default"])
	}
}

func TestDatabaseConnectionsIsReadOnly(t *testing.T) {
	t.Parallel()

	if !(&tools.DatabaseConnections{}).IsReadOnly() {
		t.Error("DatabaseConnections should be read-only")
	}
}
