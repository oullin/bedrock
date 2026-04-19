package tools_test

import (
	"testing"

	"github.com/bedrock/packages/boost/mcp/tools"
)

func TestDatabaseSchemaNoConnection(t *testing.T) {
	t.Parallel()

	tool := &tools.DatabaseSchema{}
	resp, err := tool.Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if !resp.IsError {
		t.Error("expected error response when no DB connection")
	}
}

func TestDatabaseSchemaSchema(t *testing.T) {
	t.Parallel()

	tool := &tools.DatabaseSchema{}
	schema := tool.Schema()

	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want \"object\"", schema["type"])
	}

	props, ok := schema["properties"].(map[string]any)

	if !ok {
		t.Fatal("schema missing properties map")
	}

	for _, k := range []string{"filter", "summary", "include_views"} {
		if _, exists := props[k]; !exists {
			t.Errorf("schema missing property %q", k)
		}
	}
}

func TestDatabaseSchemaIsReadOnly(t *testing.T) {
	t.Parallel()

	if !(&tools.DatabaseSchema{}).IsReadOnly() {
		t.Error("DatabaseSchema should be read-only")
	}
}
