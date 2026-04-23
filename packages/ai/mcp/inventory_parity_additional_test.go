package mcp_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/ai/mcp"
)

func TestInventoryServerContextAndCapabilityParity(t *testing.T) {
	t.Parallel()

	// ServerContextTest::it_clamps_perpage_to_default_and_max_values
	ctx := mcp.ExportServerContext(mcp.NewServer("srv", "1.0.0", mcp.WithPagination(3, 5)))

	if got := ctx.PerPage(0); got != 3 {
		t.Fatalf("expected default per-page value 3, got %d", got)
	}

	if got := ctx.PerPage(2); got != 2 {
		t.Fatalf("expected requested per-page value 2, got %d", got)
	}

	if got := ctx.PerPage(9); got != 5 {
		t.Fatalf("expected max-clamped per-page value 5, got %d", got)
	}

	// ServerTest::it_can_add_a_capability
	// InitializeTest::it_returns_a_valid_initialize_response
	srv := mcp.NewServer("srv", "1.0.0")
	base := mcp.NewPrompt("suggest", "Suggest values", nil, func(_ context.Context, _ *mcp.Request) ([]*mcp.Message, error) {
		return nil, nil
	})
	srv.AddPrompt(&completablePrompt{Prompt: base})

	resp := sendRaw(t, srv, "initialize", map[string]any{"protocolVersion": "2025-11-25"})
	result := mustResult(t, resp)
	caps, _ := result["capabilities"].(map[string]any)

	if caps["completions"] == nil {
		t.Fatalf("expected completions capability when a completable prompt is registered, got %#v", caps)
	}
}

func TestInventoryReadResourceTemplateParity(t *testing.T) {
	t.Parallel()

	// ReadResourceTest::it_tries_static_resources_before_template_matching
	// ReadResourceTest::it_returns_the_actual_requested_uri_in_response_not_the_template_pattern
	// ReadResourceTest::it_extracts_variables_from_uri_template_and_passes_to_handler
	// ReadResourceTest::it_preserves_sessionid_and_meta_from_the_original_request_for_template_resources
	// ReadResourceTest::it_sets_uri_on_request_when_reading_resource_templates
	// ResourceTemplateTest::it_does_not_leak_variables_between_consecutive_template_resource_requests
	// ResourceTemplateTest::it_uri_is_correctly_set_and_isolated_for_consecutive_requests
	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddResource(mcp.NewResource("static", "Static resource", "file://users/42", "text/plain",
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			return mcp.Text("static"), nil
		}))

	srv.AddResource(mcp.NewResourceTemplate("user", "User template", "file://users/{id}", "text/plain",
		func(_ context.Context, req *mcp.Request) (*mcp.Response, error) {
			switch req.SessionID {
			case "session-1":
				if req.Meta["trace"] != "abc" {
					t.Fatalf("expected request meta to survive template dispatch, got %#v", req.Meta)
				}
			case "session-2":
				if req.Meta != nil {
					t.Fatalf("expected second request to carry no meta, got %#v", req.Meta)
				}
			default:
				t.Fatalf("expected session id to survive template dispatch, got %q", req.SessionID)
			}

			if req.URI == "" {
				t.Fatal("expected request URI to be set")
			}

			if req.URIVars["id"] == "" {
				t.Fatalf("expected extracted URI variable, got %#v", req.URIVars)
			}

			return mcp.Text(req.URI + ":" + req.URIVars["id"]), nil
		}))

	// The static resource should win over the template when both can match.
	srv.Test(t).ReadResource("file://users/42").AssertOK().AssertSee("static")

	rawRequest := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "resources/read",
		"params": map[string]any{
			"uri":   "file://users/7",
			"_meta": map[string]any{"trace": "abc"},
		},
	}
	payload, err := json.Marshal(rawRequest)

	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	raw, err := srv.Handle(context.Background(), string(payload), "session-1")

	if err != nil {
		t.Fatalf("unexpected handle error: %v", err)
	}

	first := decodeJSON(t, raw)
	result := mustResult(t, first)
	contents, _ := result["contents"].([]any)

	if len(contents) != 1 {
		t.Fatalf("expected one resource content item, got %#v", contents)
	}

	content := contents[0].(map[string]any)

	if content["uri"] != "file://users/7" {
		t.Fatalf("expected actual requested uri in response, got %v", content["uri"])
	}

	if !strings.Contains(content["text"].(string), "file://users/7") {
		t.Fatalf("expected requested uri in content text, got %#v", content["text"])
	}

	secondRequest := map[string]any{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "resources/read",
		"params": map[string]any{
			"uri": "file://users/8",
		},
	}
	secondPayload, err := json.Marshal(secondRequest)

	if err != nil {
		t.Fatalf("unexpected second marshal error: %v", err)
	}

	secondRaw, err := srv.Handle(context.Background(), string(secondPayload), "session-2")

	if err != nil {
		t.Fatalf("unexpected second handle error: %v", err)
	}

	second := decodeJSON(t, secondRaw)
	secondResult := mustResult(t, second)
	secondContents, _ := secondResult["contents"].([]any)

	if len(secondContents) != 1 {
		t.Fatalf("expected one resource content item on second read, got %#v", secondContents)
	}

	secondContent := secondContents[0].(map[string]any)

	if secondContent["uri"] != "file://users/8" {
		t.Fatalf("expected second requested uri in response, got %v", secondContent["uri"])
	}

	if !strings.Contains(secondContent["text"].(string), "file://users/8") {
		t.Fatalf("expected second requested uri in content text, got %#v", secondContent["text"])
	}
}

func TestInventoryResourceTemplateOrderingParity(t *testing.T) {
	t.Parallel()

	// ResourceTemplateTest::it_returns_the_first_matching_template_when_multiple_templates_exist
	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddResource(mcp.NewResourceTemplate("first", "First template", "file://dup/{id}", "text/plain",
		func(_ context.Context, req *mcp.Request) (*mcp.Response, error) {
			return mcp.Text("first:" + req.URIVars["id"]), nil
		}))
	srv.AddResource(mcp.NewResourceTemplate("second", "Second template", "file://dup/{id}", "text/plain",
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			t.Fatal("expected first matching template to be selected before later templates")

			return nil, nil
		}))

	srv.Test(t).ReadResource("file://dup/123").AssertOK().AssertSee("first:123")
}

func TestInventoryJSONRPCValidationParity(t *testing.T) {
	t.Parallel()

	// JsonRpcRequestTest::it_throws_exception_for_missing_jsonrpc_version
	// JsonRpcNotificationTest::it_throws_exception_for_missing_jsonrpc_version_in_notification
	if _, err := mcp.ParseJsonRpcRequest([]byte(`{"id":1,"method":"ping","params":{}}`), ""); err == nil {
		t.Fatal("expected missing JSON-RPC version to fail")
	}

	// JsonRpcRequestTest::it_throws_exception_for_incorrect_jsonrpc_version
	// JsonRpcNotificationTest::it_throws_exception_for_incorrect_jsonrpc_version_in_notification
	if _, err := mcp.ParseJsonRpcRequest([]byte(`{"jsonrpc":"1.0","id":1,"method":"ping","params":{}}`), ""); err == nil {
		t.Fatal("expected incorrect JSON-RPC version to fail")
	}

	// JsonRpcRequestTest::it_throws_exception_for_invalid_id_type
	if _, err := mcp.ParseJsonRpcRequest([]byte(`{"jsonrpc":"2.0","id":{"bad":true},"method":"ping","params":{}}`), ""); err == nil {
		t.Fatal("expected object id to fail")
	}

	// JsonRpcRequestTest::it_throws_exception_for_non_string_method
	// JsonRpcNotificationTest::it_throws_exception_for_non_string_method_in_notification
	if _, err := mcp.ParseJsonRpcRequest([]byte(`{"jsonrpc":"2.0","id":1,"method":42,"params":{}}`), ""); err == nil {
		t.Fatal("expected non-string method to fail")
	}

	// JsonRpcNotificationTest::it_throws_exception_for_missing_method_in_notification
	if _, err := mcp.ParseJsonRpcRequest([]byte(`{"jsonrpc":"2.0","params":{}}`), ""); err == nil {
		t.Fatal("expected missing method to fail")
	}

	// JsonRpcResponseTest::it_converts_empty_array_params_in_notification_to_object
	notification, err := mcp.NotificationResponse("notifications/progress", nil).ToJSON()

	if err != nil {
		t.Fatalf("unexpected notification encode error: %v", err)
	}

	var decoded map[string]any

	if err := json.Unmarshal(notification, &decoded); err != nil {
		t.Fatalf("unexpected notification decode error: %v", err)
	}

	params := decoded["result"].(map[string]any)["params"].(map[string]any)

	if len(params) != 0 {
		t.Fatalf("expected empty notification params object, got %#v", params)
	}

	// JsonRpcResponseTest::it_includes_meta_in_result_when_provided_in_result_array
	// JsonRpcResponseTest::it_does_not_include_meta_when_not_in_result
	withMeta := mustResult(t, decodeJSONBytes(t, mustJSON(t, mcp.ResultResponse(1, map[string]any{"_meta": map[string]any{"trace": "abc"}}))))

	if withMeta["_meta"].(map[string]any)["trace"] != "abc" {
		t.Fatalf("expected result meta to be preserved, got %#v", withMeta)
	}

	withoutMeta := mustResult(t, decodeJSONBytes(t, mustJSON(t, mcp.ResultResponse(1, map[string]any{"ok": true}))))

	if _, ok := withoutMeta["_meta"]; ok {
		t.Fatalf("expected no meta when result has none, got %#v", withoutMeta)
	}

	// JsonRpcExceptionTest::it_converts_to_jsonrpc_error_with_id_and_without_data
	// JsonRpcExceptionTest::it_converts_to_jsonrpc_error_without_id_and_with_data
	withID := mustError(t, decodeJSONBytes(t, mustJSON(t, mcp.ErrorResponse(9, mcp.CodeInvalidParams, "invalid"))))

	if withID["message"] != "invalid" || withID["data"] != nil {
		t.Fatalf("expected error with id and no data, got %#v", withID)
	}

	withData := mustError(t, decodeJSONBytes(t, mustJSON(t, mcp.ErrorResponse(nil, mcp.CodeInvalidParams, "invalid", map[string]any{"field": "name"}))))

	if withData["data"].(map[string]any)["field"] != "name" {
		t.Fatalf("expected error data to be preserved, got %#v", withData)
	}
}

func TestInventoryListToolsRequestedPerPageParity(t *testing.T) {
	t.Parallel()

	// ListToolsTest::it_uses_requested_per_page_when_valid
	// ListToolsTest::it_respects_per_page_when_bigger_than_default
	// ListToolsTest::it_caps_per_page_at_max_pagination_length
	srv := mcp.NewServer("inventory", "1.0.0", mcp.WithPagination(2, 4))

	for i := 0; i < 6; i++ {
		srv.AddTool(indexedTool(i))
	}

	requested := mustResult(t, sendRaw(t, srv, "tools/list", map[string]any{"perPage": 3}))

	if tools := requested["tools"].([]any); len(tools) != 3 {
		t.Fatalf("expected requested perPage to return 3 tools, got %#v", tools)
	}

	capped := mustResult(t, sendRaw(t, srv, "tools/list", map[string]any{"perPage": 99}))

	if tools := capped["tools"].([]any); len(tools) != 4 {
		t.Fatalf("expected perPage to be capped at max 4, got %#v", tools)
	}
}

func TestInventoryHTTPResourceAndMethodParity(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("inventory", "1.0.0")
	srv.AddResource(staticFileResource())
	srv.AddTool(echoTool())

	// StartCommandTest::it_can_list_resources_over_http
	resourcesResp := serveJSONRPC(t, srv, http.MethodPost, map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "resources/list",
		"params":  map[string]any{},
	})

	if resourcesResp.Code != http.StatusOK || !strings.Contains(resourcesResp.Body.String(), `"resources"`) {
		t.Fatalf("unexpected HTTP resources/list response: code=%d body=%s", resourcesResp.Code, resourcesResp.Body.String())
	}

	// StartCommandTest::it_can_read_a_resource_over_http
	readResp := serveJSONRPC(t, srv, http.MethodPost, map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "resources/read",
		"params":  map[string]any{"uri": "file://static/config.json"},
	})

	if readResp.Code != http.StatusOK || !strings.Contains(readResp.Body.String(), "config data") {
		t.Fatalf("unexpected HTTP resources/read response: code=%d body=%s", readResp.Code, readResp.Body.String())
	}

	// StartCommandTest::it_can_list_dynamically_added_tools
	before := serveJSONRPC(t, srv, http.MethodPost, map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/list",
		"params":  map[string]any{},
	})

	if !strings.Contains(before.Body.String(), `"broadcastclient"`) {
		t.Fatalf("expected initial tool list to include broadcastclient, got %s", before.Body.String())
	}

	srv.AddTool(greetTool())
	after := serveJSONRPC(t, srv, http.MethodPost, map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/list",
		"params":  map[string]any{},
	})

	if !strings.Contains(after.Body.String(), `"greet"`) {
		t.Fatalf("expected dynamically added tool to be listed, got %s", after.Body.String())
	}

	// StartCommandTest::it_returns_405_for_get_requests_to_mcp_web_routes
	getResp := serveJSONRPC(t, srv, http.MethodGet, nil)

	if getResp.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected GET to return 405, got code=%d body=%s", getResp.Code, getResp.Body.String())
	}

	// StartCommandTest::it_returns_405_for_delete_requests_to_mcp_web_routes
	deleteResp := serveJSONRPC(t, srv, http.MethodDelete, nil)

	if deleteResp.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected DELETE to return 405, got code=%d body=%s", deleteResp.Code, deleteResp.Body.String())
	}
}

func TestInventoryResourceTemplateUnitParity(t *testing.T) {
	t.Parallel()

	// ResourceTemplateTest::it_matches_uris_against_a_template_pattern
	// ResourceTemplateTest::it_extracts_variables_from_matching_uri
	// ResourceTemplateTest::it_handles_template_with_single_variable
	single, err := mcp.NewUriTemplate("file://users/{id}")

	if err != nil {
		t.Fatal(err)
	}

	singleVars, ok := single.Match("file://users/42")

	if !ok || singleVars["id"] != "42" {
		t.Fatalf("expected single-variable template match, got vars=%#v ok=%v", singleVars, ok)
	}

	// ResourceTemplateTest::it_handles_complex_uri_templates_with_multiple_path_segments
	complex, err := mcp.NewUriTemplate("file://orgs/{org}/users/{user}/files/{file}")

	if err != nil {
		t.Fatal(err)
	}

	complexVars, ok := complex.Match("file://orgs/bedrock/users/ada/files/report")

	if !ok || complexVars["org"] != "bedrock" || complexVars["user"] != "ada" || complexVars["file"] != "report" {
		t.Fatalf("expected complex template variables, got vars=%#v ok=%v", complexVars, ok)
	}

	// ResourceTemplateTest::it_does_not_match_uris_with_different_path_structure
	if _, ok := complex.Match("file://orgs/bedrock/users/ada/files/report/extra"); ok {
		t.Fatal("expected template not to match a different path structure")
	}

	// ResourceTemplateTest::it_handles_template_resource_with_extracted_variables
	// ResourceTemplateTest::it_end_to_end_template_reads_uri_extracts_variables_and_returns_response
	// ReadResourceTest::it_template_handler_receives_variables_via_request_get_method
	srv := mcp.NewServer("inventory", "1.0.0")
	srv.AddResource(mcp.NewResourceTemplate("file", "File", "file://orgs/{org}/files/{file}", "text/plain",
		func(_ context.Context, req *mcp.Request) (*mcp.Response, error) {
			return mcp.Text(req.Get("org").(string) + ":" + req.Get("file").(string)), nil
		}))
	srv.Test(t).ReadResource("file://orgs/bedrock/files/report").AssertOK().AssertSee("bedrock:report")
}

func mustJSON(t testing.TB, resp *mcp.JsonRpcResponse) []byte {
	t.Helper()
	b, err := resp.ToJSON()

	if err != nil {
		t.Fatalf("unexpected JSON encode error: %v", err)
	}

	return b
}

func decodeJSONBytes(t testing.TB, raw []byte) map[string]any {
	t.Helper()

	var out map[string]any

	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("invalid JSON %q: %v", raw, err)
	}

	return out
}

func serveJSONRPC(t testing.TB, srv *mcp.Server, method string, payload map[string]any) *httptest.ResponseRecorder {
	t.Helper()

	var body *strings.Reader

	if payload == nil {
		body = strings.NewReader("")
	} else {
		b, err := json.Marshal(payload)

		if err != nil {
			t.Fatalf("unexpected request marshal error: %v", err)
		}

		body = strings.NewReader(string(b))
	}

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(method, "/mcp", body)
	srv.ServeHTTP(resp, req)

	return resp
}
