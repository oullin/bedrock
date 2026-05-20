// Port of Laravel\Ai\Tests\Feature\SubAgentTest
package ai_test

import (
	"context"
	"strings"
	"testing"

	ai "github.com/bedrock/packages/ai/sdk"
	"github.com/bedrock/packages/ai/sdk/fake"
	"github.com/bedrock/packages/ai/sdk/prompts"
	contractsai "github.com/bedrock/packages/contracts/ai"
)

// refundsAgent is a minimal Promptable that also implements CanActAsTool, used
// to verify the explicit name/description path.
type refundsAgent struct {
	*ai.AnonymousAgent
}

func (refundsAgent) Name() string        { return "refunds_agent" }
func (refundsAgent) Description() string { return "Handles refund requests." }

func TestSubAgentExposesExplicitNameAndDescription(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	m.Fake() // ensures a default provider is wired so the agent can be built

	inner := ai.NewAnonymousAgent(m, "refunds")
	tool := ai.AsTool(refundsAgent{AnonymousAgent: inner})

	if got := tool.Name(); got != "refunds_agent" {
		t.Errorf("Name() = %q, want %q", got, "refunds_agent")
	}

	if got := tool.Description(); got != "Handles refund requests." {
		t.Errorf("Description() = %q, want %q", got, "Handles refund requests.")
	}
}

func TestSubAgentDerivesNameFromTypeWhenNotCanActAsTool(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	m.Fake()

	inner := ai.NewAnonymousAgent(m, "anything")
	tool := ai.AsTool(inner)

	// AnonymousAgent → snake_case → "anonymous_agent".
	if got := tool.Name(); got != "anonymous_agent" {
		t.Errorf("Name() = %q, want %q", got, "anonymous_agent")
	}

	if !strings.Contains(tool.Description(), "anonymous_agent") {
		t.Errorf("Description() = %q, expected to contain %q", tool.Description(), "anonymous_agent")
	}
}

func TestSubAgentSchemaShape(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	m.Fake()

	tool := ai.AsTool(ai.NewAnonymousAgent(m, ""))
	schema := tool.Schema(contractsai.JsonSchema{})

	if schema["type"] != "object" {
		t.Fatalf("schema type = %v, want object", schema["type"])
	}

	props, ok := schema["properties"].(map[string]any)

	if !ok {
		t.Fatalf("schema.properties has wrong type: %T", schema["properties"])
	}

	task, ok := props["task"].(map[string]any)

	if !ok {
		t.Fatalf("schema.properties.task missing or wrong type: %T", props["task"])
	}

	if task["type"] != "string" {
		t.Errorf("task.type = %v, want string", task["type"])
	}

	required, ok := schema["required"].([]string)

	if !ok || len(required) != 1 || required[0] != "task" {
		t.Errorf("required = %v, want [task]", schema["required"])
	}
}

func TestSubAgentHandleInvokesWrappedAgentInIsolation(t *testing.T) {
	t.Parallel()

	// The wrapped agent is configured with three prior messages; SubAgent.Handle
	// must NOT forward those — it only passes the "task" text through.
	m := ai.NewManager()
	rec := m.Fake("Refund issued.")

	inner := ai.NewAnonymousAgent(m, "refunds expert").
		WithMessages([]any{"prior parent turn 1", "prior parent turn 2"})

	tool := ai.AsTool(inner)

	out, err := tool.Handle(context.Background(), &contractsai.ToolRequest{
		ID:   "call_1",
		Name: tool.Name(),
		Arguments: map[string]any{
			"task": "refund order 42",
		},
	})

	if err != nil {
		t.Fatalf("Handle error: %v", err)
	}

	if out != "Refund issued." {
		t.Errorf("Handle() = %v, want %q", out, "Refund issued.")
	}

	// Verify the recorded prompt that reached the gateway carried only the task.
	rec.AssertAgentWasPrompted(t, func(p *prompts.AgentPrompt) bool {
		return p.Text == "refund order 42"
	})
}

func TestSubAgentHandleErrorsWhenTaskMissing(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	m.Fake()

	tool := ai.AsTool(ai.NewAnonymousAgent(m, ""))

	_, err := tool.Handle(context.Background(), &contractsai.ToolRequest{
		ID:        "x",
		Name:      tool.Name(),
		Arguments: map[string]any{},
	})

	if err == nil {
		t.Fatal("expected error for missing task argument, got nil")
	}
}

func TestParentAgentDelegatesToSubAgentEndToEnd(t *testing.T) {
	t.Parallel()

	// Two managers because each fake-gateway scope is per-Manager; the sub-agent
	// runs through its own faked text stream while the parent's queue contains
	// (ToolCall, "Done.").
	parentMgr := ai.NewManager()
	parentRec := parentMgr.Fake(
		fake.ToolCall{ID: "call_1", Name: "refunds_agent", Args: map[string]any{"task": "refund order 42"}},
		"Refund processed successfully.",
	)

	subMgr := ai.NewManager()
	subRec := subMgr.Fake("Refund issued for order 42.")

	subAgent := ai.NewAnonymousAgent(subMgr, "refunds expert")
	tool := ai.AsTool(refundsAgent{AnonymousAgent: subAgent})

	parent := ai.NewAnonymousAgent(parentMgr, "router").
		WithTools([]contractsai.Tool{tool})

	resp, err := parent.Prompt(context.Background(), "I want a refund on order 42")

	if err != nil {
		t.Fatalf("parent.Prompt error: %v", err)
	}

	if got := resp.GetText(); got != "Refund processed successfully." {
		t.Errorf("parent text = %q, want %q", got, "Refund processed successfully.")
	}

	parentRec.AssertAgentWasPrompted(t, func(p *prompts.AgentPrompt) bool {
		return p.Text == "I want a refund on order 42"
	})

	// Sub-agent saw ONLY the task, not the parent's user turn.
	subRec.AssertAgentWasPrompted(t, func(p *prompts.AgentPrompt) bool {
		return p.Text == "refund order 42"
	})

	subRec.AssertAgentNotPrompted(t, func(p *prompts.AgentPrompt) bool {
		return strings.Contains(p.Text, "I want a refund")
	})
}

func TestSubAgentUsesItsOwnProviderHint(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	m.Fake("ok")

	inner := ai.NewAnonymousAgent(m, "specialist").WithProvider("stub")
	tool := ai.AsTool(inner)

	// Sanity: the wrapped agent surfaces its own provider hint via ProviderOptions,
	// independent of any parent.
	opts := inner.ProviderOptions()

	if opts["provider"] != "stub" {
		t.Errorf("ProviderOptions[provider] = %v, want %q", opts["provider"], "stub")
	}

	if tool.Agent() != inner {
		t.Error("Agent() should return the wrapped Promptable")
	}
}

func TestAsToolPanicsOnNilAgent(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic from AsTool(nil), got none")
		}
	}()

	_ = ai.AsTool(nil)
}

func TestDuplicateToolNamesSurfaceError(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	m.Fake("ok")

	dup1 := ai.AsTool(ai.NewAnonymousAgent(ai.NewManager(), "")).WithName("twin")
	dup2 := ai.AsTool(ai.NewAnonymousAgent(ai.NewManager(), "")).WithName("twin")

	parent := ai.NewAnonymousAgent(m, "router").
		WithTools([]contractsai.Tool{dup1, dup2})

	_, err := parent.Prompt(context.Background(), "go")

	if err == nil {
		t.Fatal("expected error for duplicate tool names, got nil")
	}

	if !strings.Contains(err.Error(), "duplicate tool name") {
		t.Errorf("error = %v, want substring %q", err, "duplicate tool name")
	}
}

func TestUnknownToolNameSurfacesError(t *testing.T) {
	t.Parallel()

	parentMgr := ai.NewManager()
	parentMgr.Fake(
		fake.ToolCall{ID: "x", Name: "no_such_tool", Args: map[string]any{"task": "do thing"}},
	)

	parent := ai.NewAnonymousAgent(parentMgr, "router").
		WithTools([]contractsai.Tool{
			ai.AsTool(ai.NewAnonymousAgent(ai.NewManager(), "")).WithName("known_tool"),
		})

	_, err := parent.Prompt(context.Background(), "go")

	if err == nil {
		t.Fatal("expected error for unknown tool, got nil")
	}

	if !strings.Contains(err.Error(), "no tool registered") {
		t.Errorf("error = %v, want substring %q", err, "no tool registered")
	}
}
