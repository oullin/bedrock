package mcp_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/bedrock/packages/ai/mcp"
)

// Port of \Mcp\Tests\ServerTest

func TestServerInitializeNegotiatesPreferredVersion(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("test-server", "1.0.0")
	resp := sendRaw(t, srv, "initialize", map[string]any{
		"protocolVersion": "2025-11-25",
		"clientInfo":      map[string]any{"name": "test-client", "version": "0.1"},
	})

	result := mustResult(t, resp)

	if result["protocolVersion"] != "2025-11-25" {
		t.Fatalf("expected protocolVersion=2025-11-25, got %v", result["protocolVersion"])
	}
}

func TestServerInitializeFallsBackForUnsupportedVersion(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("test-server", "1.0.0")
	resp := sendRaw(t, srv, "initialize", map[string]any{
		"protocolVersion": "1999-01-01",
	})

	result := mustResult(t, resp)
	// Should fall back to most recent supported version.
	if result["protocolVersion"] != "2025-11-25" {
		t.Fatalf("expected fallback to 2025-11-25, got %v", result["protocolVersion"])
	}
}

func TestServerInitializeIncludesInstructionsForNewProtocol(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0", mcp.WithInstructions("Be helpful."))
	resp := sendRaw(t, srv, "initialize", map[string]any{
		"protocolVersion": "2025-11-25",
	})

	result := mustResult(t, resp)

	if result["instructions"] != "Be helpful." {
		t.Fatalf("expected instructions in response, got %v", result["instructions"])
	}
}

func TestServerInitializeOmitsInstructionsForLegacyProtocol(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0", mcp.WithInstructions("Be helpful."))
	resp := sendRaw(t, srv, "initialize", map[string]any{
		"protocolVersion": "2024-11-05",
	})

	result := mustResult(t, resp)

	if _, ok := result["instructions"]; ok {
		t.Fatal("expected no instructions for 2024-11-05 protocol")
	}
}

func TestServerInitializeCapabilitiesMatchRegisteredPrimitives(t *testing.T) {
	t.Parallel()

	tool := mcp.NewTool("greet", "Greet", nil, func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
		return mcp.Text("hi"), nil
	})
	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddTool(tool)

	resp := sendRaw(t, srv, "initialize", map[string]any{"protocolVersion": "2025-11-25"})
	result := mustResult(t, resp)
	caps, _ := result["capabilities"].(map[string]any)

	if caps["tools"] == nil {
		t.Fatal("expected tools capability when tools are registered")
	}

	if caps["resources"] != nil {
		t.Fatal("expected no resources capability when no resources are registered")
	}
}

func TestServerPingReturnsEmptyObject(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	resp := sendRaw(t, srv, "ping", nil)
	result := mustResult(t, resp)

	if len(result) != 0 {
		t.Fatalf("expected empty result for ping, got %v", result)
	}
}

func TestServerUnknownMethodReturnsMethodNotFound(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	resp := sendRaw(t, srv, "nonexistent/method", nil)

	errObj := mustError(t, resp)
	code := int(errObj["code"].(float64))

	if code != mcp.CodeMethodNotFound {
		t.Fatalf("expected code %d, got %d", mcp.CodeMethodNotFound, code)
	}
}

func TestServerMalformedJSONReturnsParseError(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	raw, _ := srv.Handle(context.Background(), "{not json}", "")

	var m map[string]any

	json.Unmarshal([]byte(raw), &m) //nolint:errcheck
	errObj := mustError(t, m)
	code := int(errObj["code"].(float64))

	if code != mcp.CodeParseError {
		t.Fatalf("expected parse error code %d, got %d", mcp.CodeParseError, code)
	}
}

func TestServerInitializeServerInfoIncludesNameAndVersion(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("my-app", "2.3.4")
	resp := sendRaw(t, srv, "initialize", map[string]any{"protocolVersion": "2025-11-25"})
	result := mustResult(t, resp)
	info, _ := result["serverInfo"].(map[string]any)

	if info["name"] != "my-app" {
		t.Fatalf("expected name=my-app, got %v", info["name"])
	}

	if info["version"] != "2.3.4" {
		t.Fatalf("expected version=2.3.4, got %v", info["version"])
	}
}

// sendRaw sends a raw JSON-RPC message to the server and returns the parsed
// response map.
func sendRaw(t testing.TB, srv *mcp.Server, method string, params map[string]any) map[string]any {
	t.Helper()

	if params == nil {
		params = map[string]any{}
	}

	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	}
	b, _ := json.Marshal(req)
	resp, _ := srv.Handle(context.Background(), string(b), "")

	var m map[string]any

	json.Unmarshal([]byte(resp), &m) //nolint:errcheck

	return m
}

// mustResult extracts the "result" field or fails the test.
func mustResult(t testing.TB, m map[string]any) map[string]any {
	t.Helper()

	if e, ok := m["error"]; ok {
		t.Fatalf("unexpected error in response: %v", e)
	}

	result, ok := m["result"].(map[string]any)

	if !ok {
		t.Fatalf("expected result map, got %T: %v", m["result"], m["result"])
	}

	return result
}

// mustError extracts the "error" object or fails the test.
func mustError(t testing.TB, m map[string]any) map[string]any {
	t.Helper()
	errObj, ok := m["error"].(map[string]any)

	if !ok {
		t.Fatalf("expected error in response, got result: %v", m)
	}

	return errObj
}
