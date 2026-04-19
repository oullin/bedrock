package mcp_test

import (
	"testing"

	"github.com/bedrock/packages/mcp"
)

// Port of Laravel\Mcp\Tests\ResponseTest

func TestTextResponseIsNotError(t *testing.T) {
	t.Parallel()

	resp := mcp.Text("hello")

	if resp.IsError() {
		t.Fatal("expected non-error response for Text()")
	}
}

func TestErrorResponseIsError(t *testing.T) {
	t.Parallel()

	resp := mcp.Error("something went wrong")

	if !resp.IsError() {
		t.Fatal("expected isError=true for Error()")
	}
}

func TestNotificationResponseIsNotification(t *testing.T) {
	t.Parallel()

	resp := mcp.Notification("notifications/progress")

	if !resp.IsNotification() {
		t.Fatal("expected IsNotification=true")
	}
}

func TestResponseDefaultRoleIsUser(t *testing.T) {
	t.Parallel()

	resp := mcp.Text("hello")

	if resp.Role() != "user" {
		t.Fatalf("expected role=user, got %q", resp.Role())
	}
}

func TestResponseAsAssistantSetsRole(t *testing.T) {
	t.Parallel()

	resp := mcp.Text("hello").AsAssistant()

	if resp.Role() != "assistant" {
		t.Fatalf("expected role=assistant, got %q", resp.Role())
	}
}

func TestResponseWithMetaAttachesMeta(t *testing.T) {
	t.Parallel()

	resp := mcp.Text("data").WithMeta("source", "db")
	// Meta is internal; verify it's included in the tool result.
	result := mcp.ExportToolResult(resp)
	meta, ok := result["_meta"].(map[string]any)

	if !ok {
		t.Fatal("expected _meta in tool result")
	}

	if meta["source"] != "db" {
		t.Fatalf("expected source=db in _meta, got %v", meta["source"])
	}
}

func TestResponseStructuredAttachesData(t *testing.T) {
	t.Parallel()

	resp := mcp.Text("ok").Structured(map[string]any{"count": 3})
	result := mcp.ExportToolResult(resp)
	structured, ok := result["structuredContent"].(map[string]any)

	if !ok {
		t.Fatal("expected structuredContent in tool result")
	}

	if structured["count"] != 3 {
		t.Fatalf("expected count=3, got %v", structured["count"])
	}
}

func TestResponseContentsReturnsItems(t *testing.T) {
	t.Parallel()

	resp := mcp.Text("a")

	if len(resp.Contents()) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(resp.Contents()))
	}
}
