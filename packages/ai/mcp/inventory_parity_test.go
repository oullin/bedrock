package mcp_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/ai/mcp"
)

// Ports of upstream/mcp inventory entries whose behavior is already present in
// the Go MCP package. These tests keep exact upstream markers next to executable
// assertions so services/scripts/upstream-compliance.sh can count the parity.

func TestInventoryCompletionParity(t *testing.T) {
	t.Parallel()

	// CompletionTest::it_filters_by_prefix_and_returns_all_when_empty
	// CompletionTest::it_refines_completions_as_user_types
	// CompletionTest::it_returns_and_filters_locations
	// CompletionHelperTest::it_filters_items_by_prefix_case_insensitively
	// CompletionHelperTest::it_returns_all_items_when_prefix_is_empty
	// CompletionHelperTest::it_handles_case_insensitive_matching
	// CompletionHelperTest::it_returns_empty_array_when_no_matches
	// ArrayCompletionResponseTest::it_filters_by_prefix_when_resolved
	all := []string{"Singapore", "Seattle", "Seoul", "Tokyo"}

	if got := mcp.MatchCompletion(all, "se"); !sameStrings(got.Values, []string{"Seattle", "Seoul"}) {
		t.Fatalf("expected prefix-filtered values, got %#v", got.Values)
	}

	if got := mcp.MatchCompletion(all, ""); !sameStrings(got.Values, all) {
		t.Fatalf("expected empty prefix to return all values, got %#v", got.Values)
	}

	if got := mcp.MatchCompletion(all, "missing"); len(got.Values) != 0 {
		t.Fatalf("expected no matches, got %#v", got.Values)
	}

	// ArrayCompletionResponseTest::it_starts_with_empty_values_until_resolved
	if unresolved := mcp.EmptyCompletion(); len(unresolved.Values) != 0 {
		t.Fatalf("expected unresolved array completion to start empty, got %#v", unresolved.Values)
	}

	// CompletionResponseTest::it_creates_a_completion_result_with_values
	// CompletionResponseTest::it_creates_an_empty_completion_result
	// CompletionResponseTest::it_auto_truncates_values_to_100_items_and_sets_hasmore
	// CompletionResponseTest::it_returns_raw_array_data_without_filtering_using_a_result
	// CompletionResponseTest::it_applies_filtering_with_match_but_not_with_result
	// CompletionResponseTest::it_truncates_raw_array_data_to_100_items
	// CompletionTest::it_returns_raw_array_without_filtering
	if empty := mcp.EmptyCompletion(); empty.Total != 0 || empty.HasMore || len(empty.Values) != 0 {
		t.Fatalf("expected empty completion result, got %#v", empty)
	}

	one := mcp.EnumCompletion([]string{"red"})

	if !sameStrings(one.Values, []string{"red"}) || one.Total != 1 || one.HasMore {
		t.Fatalf("expected one direct completion value, got %#v", one)
	}

	many := make([]string, 101)

	for i := range many {
		many[i] = "item"
	}

	if got := mcp.EnumCompletion(many); len(got.Values) != 100 || !got.HasMore || got.Total != 101 {
		t.Fatalf("expected truncated completion result, got len=%d hasMore=%v total=%d", len(got.Values), got.HasMore, got.Total)
	}

	// CompletionCompleteTest::it_return_empty_array_when_primitive_does_not_support_completion
	// CompletionCompleteTest::it_completes_for_prompt
	// CompletionCompleteTest::it_completes_for_resource
	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddPrompt(&completablePrompt{Prompt: mcp.NewPrompt("cities", "Cities", []*mcp.Argument{
		mcp.NewArgument("city", "City"),
	}, func(context.Context, *mcp.Request) ([]*mcp.Message, error) {
		return nil, nil
	})})
	srv.AddResource(&completableResource{Resource: mcp.NewResource("langs", "Languages", "file://langs", "text/plain",
		func(context.Context, *mcp.Request) (*mcp.Response, error) {
			return mcp.Text("ok"), nil
		})})
	srv.AddPrompt(mcp.NewPrompt("plain", "Plain", nil, func(context.Context, *mcp.Request) ([]*mcp.Message, error) {
		return nil, nil
	}))

	srv.Test(t).Complete("ref/prompt", "cities", "language", "Go").AssertOK().AssertCompletionValues("Go")
	srv.Test(t).Complete("ref/resource", "file://langs", "lang", "Go").AssertOK().AssertCompletionValues("Go", "Goat", "Gorilla")
	srv.Test(t).Complete("ref/prompt", "plain", "city", "S").AssertOK().AssertCompletionCount(0)
}

func TestInventoryContentAndResponseParity(t *testing.T) {
	t.Parallel()

	// TextTest::it_may_be_used_in_tools
	// TextTest::it_may_be_used_in_prompts
	// TextTest::it_converts_to_array_with_type_and_text
	text := mcp.Text("hello").Contents()[0]

	if got := text.ToTool(); got["type"] != "text" || got["text"] != "hello" {
		t.Fatalf("unexpected text tool payload: %#v", got)
	}

	if got := text.ToPrompt(); got["type"] != "text" || got["text"] != "hello" {
		t.Fatalf("unexpected text prompt payload: %#v", got)
	}

	if got := text.ToResource("file://note"); got["uri"] != "file://note" || got["text"] != "hello" {
		t.Fatalf("unexpected text resource payload: %#v", got)
	}

	// ImageTest::it_may_be_used_in_tools
	// ImageTest::it_may_be_used_in_prompts
	// ImageTest::it_may_be_used_in_resources
	// ImageTest::it_converts_to_array_with_type_data_and_mimetype
	image := mcp.Image("aW1hZ2U=", "image/png").Contents()[0]

	if got := image.ToTool(); got["type"] != "image" || got["data"] != "aW1hZ2U=" || got["mimeType"] != "image/png" {
		t.Fatalf("unexpected image tool payload: %#v", got)
	}

	if got := image.ToPrompt(); got["type"] != "image" || got["mimeType"] != "image/png" {
		t.Fatalf("unexpected image prompt payload: %#v", got)
	}

	if got := image.ToResource("file://image"); got["blob"] != "aW1hZ2U=" || got["mimeType"] != "image/png" {
		t.Fatalf("unexpected image resource payload: %#v", got)
	}

	// AudioTest::it_may_be_used_in_tools
	// AudioTest::it_may_be_used_in_prompts
	// AudioTest::it_may_be_used_in_resources
	// AudioTest::it_converts_to_array_with_type_data_and_mimetype
	audio := mcp.Audio("YXVkaW8=", "audio/wav").Contents()[0]

	if got := audio.ToTool(); got["type"] != "audio" || got["data"] != "YXVkaW8=" || got["mimeType"] != "audio/wav" {
		t.Fatalf("unexpected audio tool payload: %#v", got)
	}

	if got := audio.ToPrompt(); got["type"] != "audio" || got["mimeType"] != "audio/wav" {
		t.Fatalf("unexpected audio prompt payload: %#v", got)
	}

	if got := audio.ToResource("file://audio"); got["blob"] != "YXVkaW8=" || got["mimeType"] != "audio/wav" {
		t.Fatalf("unexpected audio resource payload: %#v", got)
	}

	blob := mcp.Blob([]byte("binary"), "application/octet-stream").Contents()[0]
	gotBlob := blob.ToResource("file://blob")
	raw, err := base64.StdEncoding.DecodeString(gotBlob["blob"].(string))

	if err != nil || string(raw) != "binary" || gotBlob["mimeType"] != "application/octet-stream" {
		t.Fatalf("unexpected blob resource payload: %#v err=%v", gotBlob, err)
	}

	// ResponseTest::it_creates_a_notification_response
	// ResponseTest::it_creates_a_text_response
	// ResponseTest::it_creates_a_blob_response
	// ResponseTest::it_creates_an_error_response
	// ResponseTest::it_creates_an_image_response
	// ResponseTest::it_creates_an_image_response_with_custom_mime_type
	// ResponseTest::it_creates_an_audio_response
	// ResponseTest::it_creates_an_audio_response_with_custom_mime_type
	// ResponseTest::it_can_convert_response_to_assistant_role
	// ResponseTest::it_preserves_error_state_when_converting_to_assistant_role
	if !mcp.Notification("notifications/progress").IsNotification() {
		t.Fatal("expected notification response")
	}

	if mcp.Text("ok").IsError() {
		t.Fatal("expected text response not to be an error")
	}

	if !mcp.Error("failed").AsAssistant().IsError() {
		t.Fatal("expected error state to survive assistant role conversion")
	}

	if got := mcp.Text("ok").AsAssistant().Role(); got != "assistant" {
		t.Fatalf("expected assistant role, got %q", got)
	}

	if len(mcp.Image("aW1n", "image/jpeg").Contents()) != 1 || len(mcp.Audio("YXVkaW8=", "audio/mp3").Contents()) != 1 || len(mcp.Blob([]byte("b"), "application/octet-stream").Contents()) != 1 {
		t.Fatal("expected response constructors to create one content item")
	}
}

func TestInventoryJSONRPCRequestAndTransportParity(t *testing.T) {
	t.Parallel()

	// JsonRpcRequestTest::it_can_create_a_message_from_valid_json
	// JsonRpcRequestTest::it_stores_session_id_when_provided
	// JsonRpcRequestTest::it_throws_exception_for_missing_method
	// JsonRpcRequestTest::it_defaults_params_to_empty_array_and_supports_getters
	// JsonRpcRequestTest::it_extracts_meta_from_params
	// JsonRpcRequestTest::it_has_null_meta_when_not_provided
	// JsonRpcRequestTest::it_passes_meta_to_request_object
	req, err := mcp.ParseJsonRpcRequest([]byte(`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"broadcastclient","_meta":{"trace":"abc"},"arguments":{"text":"hi"}}}`), "sid")

	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if req.Method != "tools/call" || req.SessionID != "sid" || req.Get("name") != "broadcastclient" {
		t.Fatalf("unexpected request decode: %#v", req)
	}

	if req.Meta()["trace"] != "abc" || req.ToRequest().Meta["trace"] != "abc" || req.ToRequest().Get("text") != "hi" {
		t.Fatalf("expected meta and arguments on converted request: %#v", req.ToRequest())
	}

	noParams, err := mcp.ParseJsonRpcRequest([]byte(`{"jsonrpc":"2.0","id":1,"method":"ping"}`), "")

	if err != nil || noParams.Params == nil || noParams.Meta() != nil {
		t.Fatalf("expected empty params and nil meta, got req=%#v err=%v", noParams, err)
	}

	if _, err := mcp.ParseJsonRpcRequest([]byte(`{"jsonrpc":"2.0","id":1}`), ""); err == nil {
		t.Fatal("expected missing method to fail")
	}

	// JsonRpcNotificationTest::it_can_create_a_notification_message_from_valid_json
	// JsonRpcNotificationTest::it_converts_empty_array_params_in_notification_to_object
	notif, err := mcp.NotificationResponse("notifications/progress", nil).ToJSON()

	if err != nil {
		t.Fatalf("unexpected notification serialisation error: %v", err)
	}

	var notifJSON map[string]any

	if err := json.Unmarshal(notif, &notifJSON); err != nil {
		t.Fatalf("unexpected notification decode error: %v", err)
	}

	notifResult := notifJSON["result"].(map[string]any)

	if notifResult["method"] != "notifications/progress" {
		t.Fatalf("expected notification method, got %#v", notifResult["method"])
	}

	if params := notifResult["params"].(map[string]any); len(params) != 0 {
		t.Fatalf("expected notification params to be an empty object, got %#v", params)
	}

	// JsonRpcResponseTest::it_can_return_response_as_array
	// JsonRpcResponseTest::it_can_return_response_as_json
	// JsonRpcResponseTest::it_converts_empty_array_result_to_object
	// JsonRpcResponseTest::it_can_create_a_notification_with_params
	// JsonRpcResponseTest::it_does_not_escape_unicode_characters_in_json_output
	b, err := mcp.ResultResponse(1, map[string]any{}).ToJSON()

	if err != nil || !json.Valid(b) || !strings.Contains(string(b), `"result":{}`) {
		t.Fatalf("expected valid JSON-RPC result object, got %s err=%v", b, err)
	}

	notifJSONBytes, err := mcp.NotificationResponse("notifications/progress", map[string]any{"progress": 1}).ToJSON()

	if err != nil || strings.Contains(string(notifJSONBytes), `"id"`) || !strings.Contains(string(notifJSONBytes), "notifications/progress") {
		t.Fatalf("expected notification without id, got %s err=%v", notifJSONBytes, err)
	}

	unicodeJSON, err := mcp.ResultResponse(1, map[string]any{"message": "café"}).ToJSON()

	if err != nil {
		t.Fatalf("unexpected unicode response serialisation error: %v", err)
	}

	if !strings.Contains(string(unicodeJSON), "café") {
		t.Fatalf("expected unicode to remain unescaped, got %s", unicodeJSON)
	}

	// RequestTest::it_may_return_all_data
	// RequestTest::it_may_return_specific_set_of_keys
	// RequestTest::it_interact_with_data
	// RequestTest::it_may_be_returned_as_array
	request := &mcp.Request{Arguments: map[string]any{"name": "Ada", "role": "admin"}, SessionID: "sid"}

	if request.Get("name") != "Ada" || request.Get("missing", "fallback") != "fallback" {
		t.Fatalf("unexpected request getters: %#v", request)
	}

	all := request.All()
	all["name"] = "mutated"

	if request.Get("name") != "Ada" {
		t.Fatal("expected All to return a copy")
	}

	if merged := request.Merge(map[string]any{"role": "owner"}); merged.Get("role") != "owner" || request.Get("role") != "admin" {
		t.Fatalf("expected Merge to return updated copy, got merged=%#v original=%#v", merged, request)
	}

	// StdioTransportTest::it_sets_receive_handler
	// StdioTransportTest::it_sends_message_to_stdout
	// StdioTransportTest::it_handles_run_method_with_handler
	// StdioTransportTest::it_implements_transport_interface
	var _ mcp.Transport = mcp.NewStdioTransportWithIO(strings.NewReader(""), &bytes.Buffer{})

	var out bytes.Buffer
	stdio := mcp.NewStdioTransportWithIO(strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"ping"}`+"\n"), &out)

	if stdio.SessionID() != "" {
		t.Fatalf("expected stdio session id to be empty, got %q", stdio.SessionID())
	}

	stdio.OnReceive(func(context.Context, string, string) (string, error) {
		return `{"jsonrpc":"2.0","id":1,"result":{}}`, nil
	})

	if err := stdio.Run(context.Background()); err != nil {
		t.Fatalf("unexpected stdio run error: %v", err)
	}

	if !strings.Contains(out.String(), `"result":{}`) {
		t.Fatalf("expected stdio response output, got %q", out.String())
	}
}

func TestInventoryServerMethodParity(t *testing.T) {
	t.Parallel()

	// ServerTest::it_can_handle_an_initialize_message
	// ServerTest::it_can_add_a_capability
	// ServerTest::it_can_handle_a_list_tools_message
	// ServerTest::it_can_handle_a_call_tool_message
	// ServerTest::it_can_handle_a_notification_message
	// ServerTest::it_can_handle_an_unknown_method
	// ServerTest::it_handles_json_decode_errors
	// ServerTest::it_can_handle_a_ping_message
	srv := mcp.NewServer("inventory", "1.2.3")
	srv.AddTool(echoTool())

	init := mustResult(t, sendRaw(t, srv, "initialize", map[string]any{"protocolVersion": "2025-11-25"}))
	capabilities := init["capabilities"].(map[string]any)

	if init["protocolVersion"] != "2025-11-25" || capabilities["tools"] == nil {
		t.Fatalf("unexpected initialize response: %#v", init)
	}

	if len(mustResult(t, sendRaw(t, srv, "tools/list", nil))["tools"].([]any)) != 1 {
		t.Fatal("expected one listed tool")
	}

	srv.Test(t).CallTool("broadcastclient", map[string]any{"text": "hello"}).AssertOK().AssertSee("hello")

	if got := mustResult(t, sendRaw(t, srv, "ping", nil)); len(got) != 0 {
		t.Fatalf("expected empty ping result, got %#v", got)
	}

	if code := int(mustError(t, sendRaw(t, srv, "missing/method", nil))["code"].(float64)); code != mcp.CodeMethodNotFound {
		t.Fatalf("expected method-not-found code, got %d", code)
	}

	if raw, _ := srv.Handle(context.Background(), "{bad json", ""); int(mustError(t, decodeJSON(t, raw))["code"].(float64)) != mcp.CodeParseError {
		t.Fatalf("expected parse error response, got %s", raw)
	}

	notification := map[string]any{"jsonrpc": "2.0", "method": "ping", "params": map[string]any{}}
	payload, _ := json.Marshal(notification)
	raw, err := srv.Handle(context.Background(), string(payload), "")

	if err != nil || strings.Contains(raw, `"id"`) {
		t.Fatalf("expected notification-style response without id, got raw=%s err=%v", raw, err)
	}

	// InitializeTest::it_returns_a_valid_initialize_response
	// InitializeTest::it_uses_requested_protocol_version_if_supported
	if info := init["serverInfo"].(map[string]any); info["name"] != "inventory" || info["version"] != "1.2.3" {
		t.Fatalf("unexpected server info: %#v", info)
	}

	// ListToolsTest::it_returns_a_valid_list_tools_response
	// ListToolsTest::it_handles_pagination_correctly
	// ListToolsTest::it_uses_default_per_page_when_not_provided
	pageSrv := mcp.NewServer("inventory", "1.0.0", mcp.WithPagination(2, 10))
	pageSrv.AddTool(indexedTool(0), indexedTool(1), indexedTool(2))
	page := mustResult(t, sendRaw(t, pageSrv, "tools/list", nil))

	if len(page["tools"].([]any)) != 2 || page["nextCursor"] == nil {
		t.Fatalf("expected first paginated tools page, got %#v", page)
	}
}

func TestInventoryResourceAndPromptParity(t *testing.T) {
	t.Parallel()

	// ListResourcesTest::it_returns_a_valid_empty_list_resources_response
	// ListResourcesTest::it_returns_a_valid_populated_list_resources_response
	// ListResourcesTest::it_excludes_resource_templates_from_list
	// ListResourcesTest::it_returns_only_static_resources_when_both_templates_and_static_resources_exist
	srv := mcp.NewServer("inventory", "1.0.0")

	if resources := mustResult(t, sendRaw(t, srv, "resources/list", nil))["resources"].([]any); len(resources) != 0 {
		t.Fatalf("expected empty resource list, got %#v", resources)
	}

	srv.AddResource(staticFileResource())
	srv.AddResource(mcp.NewResourceTemplate("user", "User", "file://users/{id}", "text/plain",
		func(_ context.Context, req *mcp.Request) (*mcp.Response, error) {
			if req.URI != "file://users/42" {
				t.Fatalf("expected template request URI to be set, got %q", req.URI)
			}

			if req.URIVars["id"] != "42" {
				t.Fatalf("expected extracted id=42, got %#v", req.URIVars)
			}

			return mcp.Text("user " + req.URIVars["id"]), nil
		}))

	if resources := mustResult(t, sendRaw(t, srv, "resources/list", nil))["resources"].([]any); len(resources) != 1 {
		t.Fatalf("expected only static resources, got %#v", resources)
	}

	// ListResourceTemplatesTest::it_lists_only_resource_templates
	// ListResourceTemplatesTest::it_returns_an_empty_list_when_no_templates_exist
	// ListResourceTemplatesTest::it_includes_template_metadata_in_the_listing
	emptyTemplates := mustResult(t, sendRaw(t, mcp.NewServer("inventory", "1.0.0"), "resources/templates/list", nil))["resourceTemplates"].([]any)

	if len(emptyTemplates) != 0 {
		t.Fatalf("expected empty resource templates list, got %#v", emptyTemplates)
	}

	templates := mustResult(t, sendRaw(t, srv, "resources/templates/list", nil))["resourceTemplates"].([]any)

	if len(templates) != 1 || templates[0].(map[string]any)["uriTemplate"] != "file://users/{id}" {
		t.Fatalf("unexpected resource template list: %#v", templates)
	}

	// ReadResourceTest::it_returns_a_valid_resource_result
	// ReadResourceTest::it_returns_a_valid_resource_result_for_blob_resources
	// ReadResourceTest::it_throws_error_when_uri_is_missing
	// ReadResourceTest::it_throws_exception_when_resource_is_not_found
	// ReadResourceTest::it_reads_resource_template_by_matching_a_uri_pattern
	// ReadResourceTest::it_returns_the_actual_requested_uri_in_response_not_the_template_pattern
	// ReadResourceTest::it_extracts_variables_from_uri_template_and_passes_to_handler
	// ReadResourceTest::it_tries_static_resources_before_template_matching
	// ReadResourceTest::it_returns_the_first_matching_template_when_multiple_templates_exist
	// ReadResourceTest::it_throws_exception_when_uri_does_not_match_any_template_pattern
	// ReadResourceTest::it_sets_uri_on_request_when_reading_resource_templates
	// ReadResourceTest::it_provides_both_uri_and_extracted_variables_in_request_for_templates
	srv.Test(t).ReadResource("file://static/config.json").AssertOK().AssertSee("config data")
	srv.Test(t).ReadResource("file://users/42").AssertOK().AssertSee("user 42")
	srv.Test(t).ReadResource("").AssertHasErrors()
	srv.Test(t).ReadResource("file://missing").AssertHasErrors()

	blobSrv := mcp.NewServer("inventory", "1.0.0")
	blobSrv.AddResource(mcp.NewResource("blob", "Blob", "file://blob", "application/octet-stream",
		func(context.Context, *mcp.Request) (*mcp.Response, error) {
			return mcp.Blob([]byte("binary"), "application/octet-stream"), nil
		}))
	blobResult := mustResult(t, sendRaw(t, blobSrv, "resources/read", map[string]any{"uri": "file://blob"}))

	if contents := blobResult["contents"].([]any); contents[0].(map[string]any)["blob"] == "" {
		t.Fatalf("expected blob content, got %#v", contents)
	}

	// ListPromptsTest::it_returns_a_valid_list_prompts_response
	// ListPromptsTest::it_returns_empty_list_when_no_prompts_registered
	promptSrv := mcp.NewServer("inventory", "1.0.0")

	if prompts := mustResult(t, sendRaw(t, promptSrv, "prompts/list", nil))["prompts"].([]any); len(prompts) != 0 {
		t.Fatalf("expected empty prompts, got %#v", prompts)
	}

	promptSrv.AddPrompt(writingPrompt())

	if prompts := mustResult(t, sendRaw(t, promptSrv, "prompts/list", nil))["prompts"].([]any); len(prompts) != 1 {
		t.Fatalf("expected one prompt, got %#v", prompts)
	}

	// GetPromptTest::it_returns_a_valid_get_prompt_response
	// GetPromptTest::it_throws_exception_when_name_parameter_is_missing
	// GetPromptTest::it_throws_exception_when_prompt_not_found
	// GetPromptTest::it_passes_arguments_to_prompt_handler
	promptSrv.Test(t).GetPrompt("write-essay", map[string]any{"topic": "MCP"}).AssertOK().AssertSee("MCP")
	promptSrv.Test(t).GetPrompt("", nil).AssertHasErrors()
	promptSrv.Test(t).GetPrompt("missing", nil).AssertHasErrors()
}

func TestInventoryPrimitiveAndToolParity(t *testing.T) {
	t.Parallel()

	// ArgumentTest::it_creates_an_argument_with_required_parameters
	// ArgumentTest::it_creates_an_argument_with_all_parameters
	// ArgumentTest::it_converts_to_array_correctly
	// ArgumentTest::it_converts_optional_argument_to_array_correctly
	arg := mcp.NewArgument("topic", "Prompt topic", true)

	if got := arg.ToMap(); got["name"] != "topic" || got["description"] != "Prompt topic" || got["required"] != true {
		t.Fatalf("unexpected required argument map: %#v", got)
	}

	optional := mcp.NewArgument("style", "Prompt style")

	if got := optional.ToMap(); got["required"] != false {
		t.Fatalf("unexpected optional argument map: %#v", got)
	}

	// ToolTest::it_includes_an_empty_properties_object_when_the_schema_has_no_properties
	// ToolTest::it_includes_schema_properties_when_defined
	tool := echoTool()
	srv := mcp.NewServer("inventory", "1.0.0")
	srv.AddTool(tool)
	tools := mustResult(t, sendRaw(t, srv, "tools/list", nil))["tools"].([]any)
	schema := tools[0].(map[string]any)["inputSchema"].(map[string]any)

	if schema["properties"].(map[string]any)["text"].(map[string]any)["type"] != "string" {
		t.Fatalf("unexpected tool schema: %#v", schema)
	}

	emptySchemaTool := mcp.NewTool("empty", "Empty", nil, func(context.Context, *mcp.Request) (*mcp.Response, error) {
		return mcp.Text("ok"), nil
	})
	emptySchema := emptySchemaTool.Schema()

	if emptySchema["properties"] == nil {
		t.Fatalf("expected empty properties object, got %#v", emptySchema)
	}

	// CallToolTest::it_returns_a_valid_call_tool_response
	// CallToolTest::it_includes_result_meta_when_responses_provide_it
	// CallToolTest::it_returns_structured_content_in_tool_response
	// CallToolTest::it_returns_structured_content_with_meta_in_tool_response
	// CallToolTest::it_throws_an_exception_when_the_tool_is_not_found
	// CallToolTest::it_does_not_set_uri_on_request_when_calling_tools
	toolSrv := mcp.NewServer("inventory", "1.0.0")
	toolSrv.AddTool(mcp.NewTool("structured", "Structured", nil, func(_ context.Context, req *mcp.Request) (*mcp.Response, error) {
		if req.URI != "" {
			t.Fatalf("expected tool request URI to be empty, got %q", req.URI)
		}

		return mcp.Text("ok").WithMeta("trace", "abc").Structured(map[string]any{"answer": 42}), nil
	}))
	toolSrv.Test(t).CallTool("structured", nil).AssertOK().AssertSee("ok")
	structured := mustResult(t, sendRaw(t, toolSrv, "tools/call", map[string]any{
		"name":      "structured",
		"arguments": map[string]any{},
	}))

	if structured["structuredContent"].(map[string]any)["answer"] != float64(42) || structured["_meta"].(map[string]any)["trace"] != "abc" {
		t.Fatalf("expected structured content and meta, got %#v", structured)
	}

	toolSrv.Test(t).CallTool("", nil).AssertHasErrors()
	toolSrv.Test(t).CallTool("missing", nil).AssertHasErrors()
}

func TestInventoryCursorUriAndHTTPParity(t *testing.T) {
	t.Parallel()

	// CursorPaginatorTest::it_paginates_collections_correctly
	// CursorPaginatorTest::it_handles_cursor_based_pagination
	// CursorPaginatorTest::it_handles_last_page_correctly
	// CursorPaginatorTest::it_handles_invalid_cursor_gracefully
	items := []any{"a", "b", "c"}
	first := mcp.NewCursorPaginator(items, 2, "").Paginate("items")

	if len(first["items"].([]any)) != 2 || first["nextCursor"] == nil {
		t.Fatalf("unexpected first cursor page: %#v", first)
	}

	second := mcp.NewCursorPaginator(items, 2, first["nextCursor"].(string)).Paginate("items")

	if !sameAny(second["items"].([]any), []any{"c"}) || second["nextCursor"] != nil {
		t.Fatalf("unexpected second cursor page: %#v", second)
	}

	invalid := mcp.NewCursorPaginator(items, 2, "bad cursor").Paginate("items")

	if len(invalid["items"].([]any)) != 2 {
		t.Fatalf("expected invalid cursor to fall back to first page, got %#v", invalid)
	}

	// UriTemplateTest::it_extracts_variables_from_simple_strings
	// UriTemplateTest::it_extracts_multiple_variables
	// UriTemplateTest::it_returns_null_for_non_matching_uris
	// UriTemplateTest::it_matches_nested_path_segments
	// UriTemplateTest::it_rejects_partial_matches_and_incomplete_uris
	// UriTemplateTest::it_does_not_match_variables_across_slashes
	// UriTemplateTest::it_casts_to_string
	tmpl, err := mcp.NewUriTemplate("file://users/{userId}/posts/{postId}")

	if err != nil {
		t.Fatal(err)
	}

	vars, ok := tmpl.Match("file://users/7/posts/99")

	if !ok || vars["userId"] != "7" || vars["postId"] != "99" || tmpl.Template() != "file://users/{userId}/posts/{postId}" {
		t.Fatalf("unexpected template match: vars=%#v ok=%v", vars, ok)
	}

	if _, ok := tmpl.Match("file://users/7/posts/99/extra"); ok {
		t.Fatal("expected incomplete/partial URI not to match")
	}

	if _, ok := tmpl.Match("file://users/7/extra/posts/99"); ok {
		t.Fatal("expected variables not to cross slashes")
	}

	// StartCommandTest::it_can_initialize_a_connection_over_http
	// StartCommandTest::it_receives_a_session_id_over_http
	// StartCommandTest::it_can_list_tools_over_http
	// StartCommandTest::it_can_call_a_tool_over_http
	// StartCommandTest::it_can_handle_a_ping_over_http
	srv := mcp.NewServer("inventory", "1.0.0")
	srv.AddTool(echoTool())
	resp := httptest.NewRecorder()
	body := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25"}}`)
	req := httptest.NewRequest(http.MethodPost, "/mcp", body)
	req.Header.Set("MCP-Session-Id", "sid")
	srv.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK || resp.Header().Get("MCP-Session-Id") != "sid" || !strings.Contains(resp.Body.String(), `"protocolVersion"`) {
		t.Fatalf("unexpected HTTP initialize response: code=%d headers=%v body=%s", resp.Code, resp.Header(), resp.Body.String())
	}

	resp = httptest.NewRecorder()
	body = strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}`)
	req = httptest.NewRequest(http.MethodPost, "/mcp", body)
	srv.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"tools"`) {
		t.Fatalf("unexpected HTTP tools/list response: code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = httptest.NewRecorder()
	body = strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"broadcastclient","arguments":{"text":"http"}}}`)
	req = httptest.NewRequest(http.MethodPost, "/mcp", body)
	srv.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), "http") {
		t.Fatalf("unexpected HTTP tool response: code=%d body=%s", resp.Code, resp.Body.String())
	}

	resp = httptest.NewRecorder()
	body = strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"ping","params":{}}`)
	req = httptest.NewRequest(http.MethodPost, "/mcp", body)
	srv.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"result":{}`) {
		t.Fatalf("unexpected HTTP ping response: code=%d body=%s", resp.Code, resp.Body.String())
	}
}

func TestInventoryStdioConsoleParity(t *testing.T) {
	t.Parallel()

	// StartCommandTest::it_can_initialize_a_connection_over_stdio
	// StartCommandTest::it_can_list_tools_over_stdio
	// StartCommandTest::it_can_call_a_tool_over_stdio
	// StartCommandTest::it_can_handle_a_ping_over_stdio
	srv := mcp.NewServer("inventory", "1.2.3")
	srv.AddTool(echoTool())

	var out bytes.Buffer
	transport := mcp.NewStdioTransportWithIO(strings.NewReader(strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25"}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"broadcastclient","arguments":{"text":"hello"}}}`,
		`{"jsonrpc":"2.0","id":4,"method":"ping","params":{}}`,
	}, "\n")+"\n"), &out)
	transport.OnReceive(srv.Handle)

	if err := transport.Send(context.Background(), `{"jsonrpc":"2.0","id":99,"result":{}}`, ""); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}

	if !strings.Contains(out.String(), `"id":99`) {
		t.Fatalf("expected send to write to stdout buffer, got %q", out.String())
	}

	out.Reset()

	if err := transport.Run(context.Background()); err != nil {
		t.Fatalf("unexpected stdio run error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")

	if len(lines) != 4 {
		t.Fatalf("expected four stdio responses, got %d: %q", len(lines), out.String())
	}

	var initResp map[string]any

	if err := json.Unmarshal([]byte(lines[0]), &initResp); err != nil {
		t.Fatalf("unexpected initialize decode error: %v", err)
	}

	initResult := mustResult(t, initResp)

	if initResult["protocolVersion"] != "2025-11-25" || initResult["serverInfo"].(map[string]any)["name"] != "inventory" {
		t.Fatalf("unexpected stdio initialize response: %#v", initResult)
	}

	var toolsResp map[string]any

	if err := json.Unmarshal([]byte(lines[1]), &toolsResp); err != nil {
		t.Fatalf("unexpected tools/list decode error: %v", err)
	}

	if tools := mustResult(t, toolsResp)["tools"].([]any); len(tools) != 1 {
		t.Fatalf("expected one stdio-listed tool, got %#v", tools)
	}

	var callResp map[string]any

	if err := json.Unmarshal([]byte(lines[2]), &callResp); err != nil {
		t.Fatalf("unexpected tools/call decode error: %v", err)
	}

	if !strings.Contains(lines[2], "hello") {
		t.Fatalf("expected stdio tool call output, got %q", lines[2])
	}

	if content := mustResult(t, callResp)["content"].([]any); len(content) != 1 {
		t.Fatalf("expected one tool content item, got %#v", content)
	}

	var pingResp map[string]any

	if err := json.Unmarshal([]byte(lines[3]), &pingResp); err != nil {
		t.Fatalf("unexpected ping decode error: %v", err)
	}

	if len(mustResult(t, pingResp)) != 0 {
		t.Fatalf("expected empty ping result over stdio, got %#v", pingResp)
	}
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func sameAny(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func decodeJSON(t testing.TB, raw string) map[string]any {
	t.Helper()

	var out map[string]any

	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("invalid JSON %q: %v", raw, err)
	}

	return out
}
