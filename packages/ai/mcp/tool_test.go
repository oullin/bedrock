package mcp_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bedrock/packages/ai/mcp"
)

// Port of \Mcp\Tests\ToolTest

func TestToolsListReturnsAllRegisteredTools(t *testing.T) {
	t.Parallel()

	// ListToolsTest::it_returns_a_valid_list_tools_response
	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddTool(echoTool()).AddTool(greetTool())

	resp := sendRaw(t, srv, "tools/list", nil)
	result := mustResult(t, resp)
	tools, _ := result["tools"].([]any)

	if len(tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(tools))
	}
}

func TestToolsListIncludesNameDescriptionSchema(t *testing.T) {
	t.Parallel()

	// ToolTest::it_includes_schema_properties_when_defined
	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddTool(echoTool())

	resp := sendRaw(t, srv, "tools/list", nil)
	result := mustResult(t, resp)
	tools, _ := result["tools"].([]any)
	tool := tools[0].(map[string]any)
	schema, _ := tool["inputSchema"].(map[string]any)

	if tool["name"] != "echo" {
		t.Fatalf("expected name=echo, got %v", tool["name"])
	}

	if tool["description"] == nil || tool["description"] == "" {
		t.Fatal("expected non-empty description")
	}

	if tool["inputSchema"] == nil {
		t.Fatal("expected inputSchema field")
	}

	if schema["type"] != "object" {
		t.Fatalf("expected schema type object, got %v", schema["type"])
	}

	if props, ok := schema["properties"].(map[string]any); !ok || len(props) != 1 {
		t.Fatalf("expected explicit schema properties to be preserved, got %#v", schema["properties"])
	}
}

func TestToolsListDefaultsNilSchemaToEmptyProperties(t *testing.T) {
	t.Parallel()

	// Unit/Tools/ToolTest::it_includes_an_empty_properties_object_when_the_schema_has_no_properties
	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddTool(greetTool())

	resp := sendRaw(t, srv, "tools/list", nil)
	result := mustResult(t, resp)
	tools, _ := result["tools"].([]any)
	tool := tools[0].(map[string]any)
	schema, _ := tool["inputSchema"].(map[string]any)

	if schema["type"] != "object" {
		t.Fatalf("expected default schema type object, got %v", schema["type"])
	}

	if props, ok := schema["properties"].(map[string]any); !ok || len(props) != 0 {
		t.Fatalf("expected empty properties object for nil schema, got %#v", schema["properties"])
	}
}

func TestToolsListPaginationFirstPage(t *testing.T) {
	t.Parallel()

	// ListToolsTest::it_handles_pagination_correctly
	// ListToolsTest::it_uses_default_per_page_when_not_provided
	srv := mcp.NewServer("srv", "1.0.0", mcp.WithPagination(2, 10))

	for i := 0; i < 5; i++ {
		srv.AddTool(indexedTool(i))
	}

	resp := sendRaw(t, srv, "tools/list", nil)
	result := mustResult(t, resp)
	tools, _ := result["tools"].([]any)

	if len(tools) != 2 {
		t.Fatalf("expected 2 tools on first page, got %d", len(tools))
	}

	if result["nextCursor"] == nil {
		t.Fatal("expected nextCursor on first page")
	}
}

func TestToolsListPaginationLastPage(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0", mcp.WithPagination(2, 10))

	for i := 0; i < 3; i++ {
		srv.AddTool(indexedTool(i))
	}

	// Get cursor from first page.
	r1 := sendRaw(t, srv, "tools/list", nil)
	res1 := mustResult(t, r1)
	cursor, _ := res1["nextCursor"].(string)

	r2 := sendRaw(t, srv, "tools/list", map[string]any{"cursor": cursor})
	res2 := mustResult(t, r2)
	tools, _ := res2["tools"].([]any)

	if len(tools) != 1 {
		t.Fatalf("expected 1 tool on last page, got %d", len(tools))
	}

	if res2["nextCursor"] != nil {
		t.Fatal("expected no nextCursor on last page")
	}
}

func TestToolsCallInvokesCorrectTool(t *testing.T) {
	t.Parallel()

	// CallToolTest::it_returns_a_valid_call_tool_response
	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddTool(echoTool())

	result := srv.Test(t).CallTool("echo", map[string]any{"text": "hello"})
	result.AssertOK().AssertSee("hello")
}

func TestToolsCallToolNotFoundReturnsError(t *testing.T) {
	t.Parallel()

	// CallToolTest::it_throws_an_exception_when_the_tool_is_not_found
	srv := mcp.NewServer("srv", "1.0.0")
	result := srv.Test(t).CallTool("nonexistent", nil)
	result.AssertHasErrors()
}

func TestToolsCallDoesNotSetURIReturnsResult(t *testing.T) {
	t.Parallel()

	// CallToolTest::it_does_not_set_uri_on_request_when_calling_tools
	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddTool(mcp.NewTool("uri-check", "Checks that tools do not receive a URI", nil,
		func(_ context.Context, req *mcp.Request) (*mcp.Response, error) {
			if req.URI != "" {
				t.Fatalf("expected tool request URI to be empty, got %q", req.URI)
			}

			return mcp.Text("ok"), nil
		}))

	result := srv.Test(t).CallTool("uri-check", nil)
	result.AssertOK().AssertSee("ok")
}

func TestToolsCallHandlerErrorProducesIsErrorResult(t *testing.T) {
	t.Parallel()

	failTool := mcp.NewTool("fail", "Always fails", nil,
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			return nil, errors.New("something broke")
		})

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddTool(failTool)

	result := srv.Test(t).CallTool("fail", nil)
	result.AssertHasErrors()
}

func TestToolsCallExplicitErrorResponse(t *testing.T) {
	t.Parallel()

	errTool := mcp.NewTool("bad", "Returns error response", nil,
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			return mcp.Error("custom error message"), nil
		})

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddTool(errTool)

	result := srv.Test(t).CallTool("bad", nil)
	result.AssertHasErrors().AssertSee("custom error message")
}

func TestToolsCallTextResponseContent(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddTool(mcp.NewTool("hi", "Say hi", nil,
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			return mcp.Text("Hello, World!"), nil
		}))

	result := srv.Test(t).CallTool("hi", nil)
	result.AssertOK().AssertSee("Hello, World!")
}

func TestToolsCallStructuredContent(t *testing.T) {
	t.Parallel()

	// CallToolTest::it_returns_structured_content_in_tool_response
	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddTool(mcp.NewTool("structured", "Returns structured content", nil,
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			return mcp.Text("result").Structured(map[string]any{"answer": 42}), nil
		}))

	result := srv.Test(t).CallTool("structured", nil)
	result.AssertOK().AssertSee("result")
}

func TestToolsCallStructuredContentWithMeta(t *testing.T) {
	t.Parallel()

	// CallToolTest::it_includes_result_meta_when_responses_provide_it
	// CallToolTest::it_returns_a_result_with_result_level_meta_when_using_responsefactory
	// CallToolTest::it_returns_structured_content_with_meta_in_tool_response
	// CallToolTest::it_returns_responsefactory_with_structured_content_added_via_withstructuredcontent
	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddTool(mcp.NewTool("structured-meta", "Returns structured content with meta", nil,
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			return mcp.Text("result").
				WithMeta("trace", "abc").
				Structured(map[string]any{"answer": 42}), nil
		}))

	resp := sendRaw(t, srv, "tools/call", map[string]any{
		"name":      "structured-meta",
		"arguments": map[string]any{},
	})
	result := mustResult(t, resp)

	if result["_meta"].(map[string]any)["trace"] != "abc" {
		t.Fatalf("expected result meta to survive tools/call, got %#v", result["_meta"])
	}

	structured, ok := result["structuredContent"].(map[string]any)

	if !ok || structured["answer"] != float64(42) {
		t.Fatalf("expected structured content in tool result, got %#v", result["structuredContent"])
	}
}

// --- helpers ---

func echoTool() mcp.Tool {
	return mcp.NewTool("echo", "Echoes the input text",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"text": map[string]any{"type": "string"},
			},
			"required": []string{"text"},
		},
		func(_ context.Context, req *mcp.Request) (*mcp.Response, error) {
			text, _ := req.Get("text").(string)

			return mcp.Text(text), nil
		})
}

func greetTool() mcp.Tool {
	return mcp.NewTool("greet", "Greet a user", nil,
		func(_ context.Context, req *mcp.Request) (*mcp.Response, error) {
			name, _ := req.Get("name").(string)

			return mcp.Text("Hello, " + name + "!"), nil
		})
}

func indexedTool(i int) mcp.Tool {
	name := "tool-" + string(rune('a'+i))

	return mcp.NewTool(name, "Tool "+name, nil,
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			return mcp.Text(name), nil
		})
}
