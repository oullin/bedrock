package mcp_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/mcp"
)

// Port of Laravel\Mcp\Tests\PromptTest

func TestPromptsListReturnsAllPrompts(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddPrompt(writingPrompt()).AddPrompt(codePrompt())

	resp := sendRaw(t, srv, "prompts/list", nil)
	result := mustResult(t, resp)
	prompts, _ := result["prompts"].([]any)

	if len(prompts) != 2 {
		t.Fatalf("expected 2 prompts, got %d", len(prompts))
	}
}

func TestPromptsListIncludesArguments(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddPrompt(writingPrompt())

	resp := sendRaw(t, srv, "prompts/list", nil)
	result := mustResult(t, resp)
	prompts, _ := result["prompts"].([]any)
	p := prompts[0].(map[string]any)
	args, _ := p["arguments"].([]any)

	if len(args) != 2 {
		t.Fatalf("expected 2 arguments, got %d", len(args))
	}
}

func TestPromptsGetInvokesPromptWithArguments(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddPrompt(writingPrompt())

	result := srv.Test(t).GetPrompt("write-essay", map[string]any{
		"topic": "Go programming",
		"style": "formal",
	})
	result.AssertOK().AssertSee("Go programming")
}

func TestPromptsGetNotFoundReturnsError(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	result := srv.Test(t).GetPrompt("nonexistent", nil)
	result.AssertHasErrors()
}

func TestPromptsGetReturnsMessagesArray(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddPrompt(writingPrompt())

	resp := sendRaw(t, srv, "prompts/get", map[string]any{
		"name":      "write-essay",
		"arguments": map[string]any{"topic": "AI", "style": "casual"},
	})
	result := mustResult(t, resp)
	msgs, _ := result["messages"].([]any)

	if len(msgs) == 0 {
		t.Fatal("expected non-empty messages")
	}

	msg := msgs[0].(map[string]any)

	if msg["role"] != "user" {
		t.Fatalf("expected role=user, got %v", msg["role"])
	}

	if msg["content"] == nil {
		t.Fatal("expected content in message")
	}
}

func TestPromptsGetUserAndAssistantRoles(t *testing.T) {
	t.Parallel()

	multiTurnPrompt := mcp.NewPrompt(
		"multi-turn", "Multi-turn example",
		nil,
		func(_ context.Context, _ *mcp.Request) ([]*mcp.Message, error) {
			return []*mcp.Message{
				mcp.UserMessage(&mcp.TextContent{Text: "user turn"}),
				mcp.AssistantMessage(&mcp.TextContent{Text: "assistant turn"}),
			}, nil
		},
	)

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddPrompt(multiTurnPrompt)

	resp := sendRaw(t, srv, "prompts/get", map[string]any{"name": "multi-turn"})
	result := mustResult(t, resp)
	msgs, _ := result["messages"].([]any)

	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}

	if msgs[0].(map[string]any)["role"] != "user" {
		t.Fatal("expected first message role=user")
	}

	if msgs[1].(map[string]any)["role"] != "assistant" {
		t.Fatal("expected second message role=assistant")
	}
}

func TestPromptsGetIncludesDescription(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0")
	srv.AddPrompt(writingPrompt())

	resp := sendRaw(t, srv, "prompts/get", map[string]any{
		"name":      "write-essay",
		"arguments": map[string]any{"topic": "test"},
	})
	result := mustResult(t, resp)

	if result["description"] == nil || result["description"] == "" {
		t.Fatal("expected non-empty description in prompts/get result")
	}
}

func TestPromptsListPagination(t *testing.T) {
	t.Parallel()

	srv := mcp.NewServer("srv", "1.0.0", mcp.WithPagination(1, 5))

	for i := 0; i < 3; i++ {
		srv.AddPrompt(indexedPrompt(i))
	}

	resp := sendRaw(t, srv, "prompts/list", nil)
	result := mustResult(t, resp)
	prompts, _ := result["prompts"].([]any)

	if len(prompts) != 1 {
		t.Fatalf("expected 1 prompt per page, got %d", len(prompts))
	}

	if result["nextCursor"] == nil {
		t.Fatal("expected nextCursor")
	}
}

// --- helpers ---

func writingPrompt() mcp.Prompt {
	return mcp.NewPrompt(
		"write-essay", "Write an essay on a topic",
		[]*mcp.Argument{
			mcp.NewArgument("topic", "Essay topic", true),
			mcp.NewArgument("style", "Writing style"),
		},
		func(_ context.Context, req *mcp.Request) ([]*mcp.Message, error) {
			topic, _ := req.Get("topic").(string)
			style, _ := req.Get("style").(string)

			if style == "" {
				style = "neutral"
			}

			return []*mcp.Message{
				mcp.UserMessage(&mcp.TextContent{
					Text: "Write a " + style + " essay about: " + topic,
				}),
			}, nil
		},
	)
}

func codePrompt() mcp.Prompt {
	return mcp.NewPrompt("write-code", "Write code", nil,
		func(_ context.Context, _ *mcp.Request) ([]*mcp.Message, error) {
			return []*mcp.Message{mcp.UserMessage(&mcp.TextContent{Text: "Write the code"})}, nil
		})
}

func indexedPrompt(i int) mcp.Prompt {
	name := "prompt-" + string(rune('a'+i))

	return mcp.NewPrompt(name, "Prompt "+name, nil,
		func(_ context.Context, _ *mcp.Request) ([]*mcp.Message, error) {
			return []*mcp.Message{mcp.UserMessage(&mcp.TextContent{Text: name})}, nil
		})
}
