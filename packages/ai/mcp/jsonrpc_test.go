package mcp_test

import (
	"encoding/json"
	"testing"

	"github.com/bedrock/packages/ai/mcp"
)

// Port of Upstream\Mcp\Tests\JsonRpcTest

func TestParseJsonRpcRequestValid(t *testing.T) {
	t.Parallel()

	raw := `{"jsonrpc":"2.0","id":1,"method":"ping","params":{}}`
	req, err := mcp.ParseJsonRpcRequest([]byte(raw), "sid-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Method != "ping" {
		t.Fatalf("expected method=ping, got %q", req.Method)
	}

	if req.SessionID != "sid-1" {
		t.Fatalf("expected sessionID=sid-1, got %q", req.SessionID)
	}
}

func TestParseJsonRpcRequestMalformedJSON(t *testing.T) {
	t.Parallel()

	_, err := mcp.ParseJsonRpcRequest([]byte("{not json"), "")

	if err == nil {
		t.Fatal("expected parse error for malformed JSON")
	}
}

func TestParseJsonRpcRequestMissingMethod(t *testing.T) {
	t.Parallel()

	raw := `{"jsonrpc":"2.0","id":1}`
	_, err := mcp.ParseJsonRpcRequest([]byte(raw), "")

	if err == nil {
		t.Fatal("expected error for missing method")
	}
}

func TestParseJsonRpcRequestNilParamsBecomeEmptyMap(t *testing.T) {
	t.Parallel()

	raw := `{"jsonrpc":"2.0","id":1,"method":"ping"}`
	req, err := mcp.ParseJsonRpcRequest([]byte(raw), "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Params == nil {
		t.Fatal("expected non-nil Params map")
	}
}

func TestJsonRpcRequestGetReturnsValue(t *testing.T) {
	t.Parallel()

	raw := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"greet"}}`
	req, _ := mcp.ParseJsonRpcRequest([]byte(raw), "")

	if req.Get("name") != "greet" {
		t.Fatalf("expected name=greet, got %v", req.Get("name"))
	}
}

func TestJsonRpcRequestGetReturnsFallback(t *testing.T) {
	t.Parallel()

	raw := `{"jsonrpc":"2.0","id":1,"method":"ping","params":{}}`
	req, _ := mcp.ParseJsonRpcRequest([]byte(raw), "")

	if req.Get("missing", "default") != "default" {
		t.Fatal("expected fallback value")
	}
}

func TestJsonRpcRequestCursorExtracted(t *testing.T) {
	t.Parallel()

	raw := `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{"cursor":"abc"}}`
	req, _ := mcp.ParseJsonRpcRequest([]byte(raw), "")

	if req.Cursor() != "abc" {
		t.Fatalf("expected cursor=abc, got %q", req.Cursor())
	}
}

func TestResultResponseIsValidJSON(t *testing.T) {
	t.Parallel()

	resp := mcp.ResultResponse(1, map[string]any{"ok": true})
	b, err := resp.ToJSON()

	if err != nil {
		t.Fatalf("serialisation error: %v", err)
	}

	var m map[string]any

	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if m["jsonrpc"] != "2.0" {
		t.Fatalf("expected jsonrpc=2.0, got %v", m["jsonrpc"])
	}
}

func TestErrorResponseHasCorrectCode(t *testing.T) {
	t.Parallel()

	resp := mcp.ErrorResponse(1, mcp.CodeMethodNotFound, "method not found")
	b, _ := resp.ToJSON()

	var m map[string]any

	json.Unmarshal(b, &m) //nolint:errcheck
	errObj, _ := m["error"].(map[string]any)

	if errObj["code"] != float64(mcp.CodeMethodNotFound) {
		t.Fatalf("expected code %d, got %v", mcp.CodeMethodNotFound, errObj["code"])
	}
}

func TestNotificationResponseHasNoID(t *testing.T) {
	t.Parallel()

	resp := mcp.NotificationResponse("notifications/progress", nil)
	b, _ := resp.ToJSON()

	var m map[string]any

	json.Unmarshal(b, &m) //nolint:errcheck

	if _, ok := m["id"]; ok {
		t.Fatal("expected no 'id' field in notification response")
	}
}

func TestToRequestExtractsArguments(t *testing.T) {
	t.Parallel()

	raw := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"greet","arguments":{"who":"World"}}}`
	req, _ := mcp.ParseJsonRpcRequest([]byte(raw), "sid")
	mcpReq := req.ToRequest()

	if mcpReq.Get("who") != "World" {
		t.Fatalf("expected who=World from arguments sub-object, got %v", mcpReq.Get("who"))
	}

	if mcpReq.SessionID != "sid" {
		t.Fatalf("expected sessionID=sid, got %q", mcpReq.SessionID)
	}
}
