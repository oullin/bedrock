package fake

import (
	"context"
	"errors"
	"strings"
	"testing"

	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
)

// TestGenerateTextDispatchesToolCallThenReturnsFinalText verifies that a queued
// ToolCall triggers the registered OnToolInvocation handler and the next queued
// response becomes the assistant turn.
func TestGenerateTextDispatchesToolCallThenReturnsFinalText(t *testing.T) {
	t.Parallel()

	gw := NewTextGateway(NewRecorder())
	gw.SetResponses([]TextResponse{
		ToolCall{ID: "call_1", Name: "refunds_agent", Args: map[string]any{"task": "refund 42"}},
		"Refund processed.",
	})

	var seenName, seenTask string

	gw.OnToolInvocation(func(_ context.Context, _, name string, args map[string]any) (any, error) {
		seenName = name

		if v, ok := args["task"].(string); ok {
			seenTask = v
		}

		return "ok", nil
	})

	res, err := gw.GenerateText(context.Background(), contractsgw.TextGenerateRequest{Text: "hi"})

	if err != nil {
		t.Fatalf("GenerateText error: %v", err)
	}

	if res.Text != "Refund processed." {
		t.Errorf("Text = %q, want %q", res.Text, "Refund processed.")
	}

	if seenName != "refunds_agent" || seenTask != "refund 42" {
		t.Errorf("handler saw name=%q task=%q, want refunds_agent / refund 42", seenName, seenTask)
	}
}

// behavior for the streaming path.
func TestStreamTextDispatchesToolCallBeforeStreaming(t *testing.T) {
	t.Parallel()

	gw := NewTextGateway(NewRecorder())
	gw.SetResponses([]TextResponse{
		ToolCall{ID: "c1", Name: "t1", Args: nil},
		"hello world",
	})

	invoked := false

	gw.OnToolInvocation(func(_ context.Context, _, _ string, _ map[string]any) (any, error) {
		invoked = true

		return nil, nil
	})

	seq, err := gw.StreamText(context.Background(), "inv", contractsgw.TextGenerateRequest{Text: "hi"})

	if err != nil {
		t.Fatalf("StreamText error: %v", err)
	}

	count := 0

	for ev := range seq {
		_ = ev

		count++
	}

	if !invoked {
		t.Error("expected tool handler to be invoked")
	}

	if count == 0 {
		t.Error("expected stream events to be yielded")
	}
}

// TestInvokeToolCallErrorsWhenNoHandler verifies that a queued ToolCall with
// no OnToolInvocation handler produces a clear error.
func TestInvokeToolCallErrorsWhenNoHandler(t *testing.T) {
	t.Parallel()

	gw := NewTextGateway(NewRecorder())
	gw.SetResponses([]TextResponse{
		ToolCall{ID: "x", Name: "orphan", Args: nil},
		"never reached",
	})

	_, err := gw.GenerateText(context.Background(), contractsgw.TextGenerateRequest{Text: "go"})

	if err == nil {
		t.Fatal("expected error when no handler is registered, got nil")
	}

	if !strings.Contains(err.Error(), "no OnToolInvocation handler") {
		t.Errorf("error = %v, want substring %q", err, "no OnToolInvocation handler")
	}
}

// TestInvokeToolCallPropagatesHandlerError verifies that handler errors abort
// the GenerateText call.
func TestInvokeToolCallPropagatesHandlerError(t *testing.T) {
	t.Parallel()

	gw := NewTextGateway(NewRecorder())
	gw.SetResponses([]TextResponse{
		ToolCall{ID: "x", Name: "bad", Args: nil},
	})

	want := errors.New("handler failed")

	gw.OnToolInvocation(func(context.Context, string, string, map[string]any) (any, error) {
		return nil, want
	})

	_, err := gw.GenerateText(context.Background(), contractsgw.TextGenerateRequest{Text: "go"})

	if !errors.Is(err, want) {
		t.Errorf("err = %v, want %v", err, want)
	}
}

// TestDequeueRawStrayCallWhenEmpty verifies the prevent-stray error path.
func TestDequeueRawStrayCallWhenEmpty(t *testing.T) {
	t.Parallel()

	gw := NewTextGateway(NewRecorder())
	gw.PreventStray()

	_, err := gw.GenerateText(context.Background(), contractsgw.TextGenerateRequest{Text: "go"})

	if !errors.Is(err, ErrStrayCall) {
		t.Errorf("err = %v, want ErrStrayCall", err)
	}
}

// TestDequeueRawDefaultsToEmptyWhenNotPreventing verifies that an empty queue
// yields a zero-value AgentResponse rather than an error.
func TestDequeueRawDefaultsToEmptyWhenNotPreventing(t *testing.T) {
	t.Parallel()

	gw := NewTextGateway(NewRecorder())

	res, err := gw.GenerateText(context.Background(), contractsgw.TextGenerateRequest{Text: "go"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Text != "" {
		t.Errorf("Text = %q, want empty", res.Text)
	}
}
