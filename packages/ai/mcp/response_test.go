package mcp_test

import (
	"testing"

	"github.com/bedrock/packages/ai/mcp"
)

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

	// ResponseFactoryTest::it_supports_fluent_withmeta_for_result_level_metadata
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

func TestResponseWithMetaMergesMultipleCalls(t *testing.T) {
	t.Parallel()

	// ResponseFactoryTest::it_supports_withmeta_with_key_value_signature
	// ResponseFactoryTest::it_merges_multiple_withmeta_calls
	resp := mcp.Text("data").
		WithMeta("source", "db").
		WithMeta("trace", "abc").
		WithMeta("source", "api")

	result := mcp.ExportToolResult(resp)
	meta, ok := result["_meta"].(map[string]any)

	if !ok {
		t.Fatal("expected _meta in tool result")
	}

	if meta["source"] != "api" {
		t.Fatalf("expected later WithMeta call to overwrite source, got %v", meta["source"])
	}

	if meta["trace"] != "abc" {
		t.Fatalf("expected trace=abc in _meta, got %v", meta["trace"])
	}
}

func TestResponseStructuredAttachesData(t *testing.T) {
	t.Parallel()

	// ResponseFactoryTest::it_creates_a_structured_content_response_with_response_structured
	// ResponseFactoryTest::it_creates_a_structured_content_response_with_meta_using_response_structured
	// ResponseFactoryTest::it_adds_structured_content_to_existing_responsefactory_with_withstructuredcontent
	// ResponseFactoryTest::it_adds_structured_content_with_meta_to_responsefactory
	resp := mcp.Text("ok").
		Structured(map[string]any{"count": 3}).
		WithMeta("trace", "abc")
	result := mcp.ExportToolResult(resp)
	structured, ok := result["structuredContent"].(map[string]any)

	if !ok {
		t.Fatal("expected structuredContent in tool result")
	}

	if structured["count"] != 3 {
		t.Fatalf("expected count=3, got %v", structured["count"])
	}

	if meta := result["_meta"].(map[string]any); meta["trace"] != "abc" {
		t.Fatalf("expected trace=abc in _meta, got %#v", meta)
	}
}

func TestResponseContentsReturnsItems(t *testing.T) {
	t.Parallel()

	resp := mcp.Text("a")

	if len(resp.Contents()) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(resp.Contents()))
	}
}
