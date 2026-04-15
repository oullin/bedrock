package tools_test

import (
	"runtime"
	"testing"

	"github.com/bedrock/packages/boost/mcp/tools"
)

// TestApplicationInfoHandle verifies the response contains required runtime fields.
func TestApplicationInfoHandle(t *testing.T) {
	t.Parallel()

	tool := &tools.ApplicationInfo{}
	resp, err := tool.Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("Handle: unexpected error: %v", err)
	}

	if resp.IsError {
		t.Fatalf("Handle returned error response: %v", resp.Content)
	}

	if len(resp.Content) == 0 {
		t.Fatal("response has no content")
	}

	data, ok := resp.Content[0].Data.(map[string]any)
	if !ok {
		t.Fatalf("content data type = %T, want map[string]any", resp.Content[0].Data)
	}

	wantKeys := []string{"go_version", "os", "arch"}
	for _, k := range wantKeys {
		if _, exists := data[k]; !exists {
			t.Errorf("response missing key %q", k)
		}
	}

	if data["os"] != runtime.GOOS {
		t.Errorf("os = %v, want %v", data["os"], runtime.GOOS)
	}

	if data["arch"] != runtime.GOARCH {
		t.Errorf("arch = %v, want %v", data["arch"], runtime.GOARCH)
	}
}

func TestApplicationInfoSchema(t *testing.T) {
	t.Parallel()

	tool := &tools.ApplicationInfo{}
	schema := tool.Schema()

	if schema["type"] != "object" {
		t.Errorf("schema type = %v, want \"object\"", schema["type"])
	}
}

func TestApplicationInfoIsReadOnly(t *testing.T) {
	t.Parallel()

	if !(&tools.ApplicationInfo{}).IsReadOnly() {
		t.Error("ApplicationInfo should be read-only")
	}
}
