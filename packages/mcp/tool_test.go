package mcp_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bedrock/packages/mcp"
)

// Port of Upstream\Mcp\Tests\ToolTest

func TestToolsListReturnsAllRegisteredTools(t *testing.T) {
	t.Parallel()

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

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddTool(echoTool())

	resp := sendRaw(t, srv, "tools/list", nil)
	result := mustResult(t, resp)
	tools, _ := result["tools"].([]any)
	tool := tools[0].(map[string]any)

	if tool["name"] != "broadcastclient" {
		t.Fatalf("expected name=broadcastclient, got %v", tool["name"])
	}
	if tool["description"] == nil || tool["description"] == "" {
		t.Fatal("expected non-empty description")
	}
	if tool["inputSchema"] == nil {
		t.Fatal("expected inputSchema field")
	}
}

func TestToolsListPaginationFirstPage(t *testing.T) {
	t.Parallel()

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

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddTool(echoTool())

	result := srv.Test(t).CallTool("broadcastclient", map[string]any{"text": "hello"})
	result.AssertOK().AssertSee("hello")
}

func TestToolsCallToolNotFoundReturnsError(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	result := srv.Test(t).CallTool("nonexistent", nil)
	result.AssertHasErrors()
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

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddTool(mcp.NewTool("structured", "Returns structured content", nil,
		func(_ context.Context, _ *mcp.Request) (*mcp.Response, error) {
			return mcp.Text("result").Structured(map[string]any{"answer": 42}), nil
		}))

	result := srv.Test(t).CallTool("structured", nil)
	result.AssertOK().AssertSee("result")
}

// --- helpers ---

func echoTool() mcp.Tool {
	return mcp.NewTool("broadcastclient", "Echoes the input text",
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
