package tools_test

import (
	"testing"

	"github.com/bedrock/packages/ai/boost/mcp/tools"
)

func TestGetAbsoluteUrlPath(t *testing.T) {
	t.Parallel()

	tool := &tools.GetAbsoluteUrl{BaseURL: "http://localhost:8080"}
	resp, err := tool.Handle(tools.McpRequest{Args: map[string]any{"path": "/users/1"}})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if resp.IsError {
		t.Errorf("Handle returned error: %v", resp.Content)
	}

	data, ok := resp.Content[0].Data.(map[string]any)

	if !ok {
		t.Fatalf("data type = %T", resp.Content[0].Data)
	}

	if data["url"] != "http://localhost:8080/users/1" {
		t.Errorf("url = %v, want \"http://localhost:8080/users/1\"", data["url"])
	}
}

func TestGetAbsoluteUrlNamedRoute(t *testing.T) {
	t.Parallel()

	tool := &tools.GetAbsoluteUrl{
		BaseURL: "http://app.test",
		Routes:  map[string]string{"home": "/", "users.index": "/users"},
	}

	resp, err := tool.Handle(tools.McpRequest{Args: map[string]any{"route": "users.index"}})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	data, ok := resp.Content[0].Data.(map[string]any)

	if !ok {
		t.Fatalf("data type = %T", resp.Content[0].Data)
	}

	if data["url"] != "http://app.test/users" {
		t.Errorf("url = %v, want \"http://app.test/users\"", data["url"])
	}
}

func TestGetAbsoluteUrlDefaultBase(t *testing.T) {
	t.Parallel()

	tool := &tools.GetAbsoluteUrl{} // no BaseURL
	resp, err := tool.Handle(tools.McpRequest{Args: map[string]any{"path": "/api"}})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	data, ok := resp.Content[0].Data.(map[string]any)

	if !ok {
		t.Fatalf("data type = %T", resp.Content[0].Data)
	}

	u, _ := data["url"].(string)

	if u == "" {
		t.Error("url should not be empty")
	}
}

func TestGetAbsoluteUrlIsReadOnly(t *testing.T) {
	t.Parallel()

	if !(&tools.GetAbsoluteUrl{}).IsReadOnly() {
		t.Error("GetAbsoluteUrl should be read-only")
	}
}
