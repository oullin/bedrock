// Port of Laravel\Ai\Tests\Feature\AgentFakeTest
package ai_test

import (
	"context"
	"strings"
	"testing"

	"github.com/bedrock/packages/ai"
	"github.com/bedrock/packages/ai/fake"
	"github.com/bedrock/packages/ai/prompts"
	"github.com/bedrock/packages/ai/responses"
	"github.com/bedrock/packages/ai/stream"
)

// newTestManager returns a Manager pre-configured with a stub provider for testing.

// Use the stub provider (auto-registered when Fake() is called).

// TestAgentCanBeFakedWithAString mirrors test_agent_can_be_faked_with_a_string.

// TestAgentCanBeFakedWithMultipleStrings mirrors test_agent_can_be_faked_with_multiple_strings.

// TestAgentLastResponseRepeats verifies FIFO queue with last-response repeat.

// TestAgentCanBeFakedWithAnAgentResponse mirrors faking with a struct response.

// TestAgentCanBeFakedWithAClosure mirrors faking with a function.

// TestAgentAssertNeverPrompted mirrors the "never prompted" assertion.

// No Prompt call was made.

// TestAgentAssertNotPrompted mirrors the "not prompted with specific match" assertion.

//nolint:errcheck

// TestAgentCanBeQueued verifies queue recording.

// TestAgentAssertNeverQueued mirrors the "never queued" assertion.

// No Queue call.

// TestAgentCanBeStreamed verifies the streaming path.

// TestManagerProviderTypeCheck verifies ErrProviderCapability.

// Register ElevenLabs (audio only) under its lab name.

// Return a minimal audio-only stub (won't be called, just registered).

// audioOnlyStub satisfies AudioProvider but NOT TextProvider.
type audioOnlyStub struct{}

func newTestManager(t *testing.T) *ai.Manager {
	t.Helper()
	m := ai.NewManager()

	return m
}

func TestAgentCanBeFakedWithAString(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake("Hello, Agent!")

	agent := ai.NewAnonymousAgent(m, "You are helpful.")
	resp, err := agent.Prompt(context.Background(), "hi")

	if err != nil {
		t.Fatalf("Prompt error: %v", err)
	}

	if resp.GetText() != "Hello, Agent!" {
		t.Errorf("expected %q got %q", "Hello, Agent!", resp.GetText())
	}

	rec.AssertAgentWasPrompted(t, func(p *prompts.AgentPrompt) bool {
		return p.Text == "hi"
	})
}

func TestAgentCanBeFakedWithMultipleStrings(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake("First response", "Second response")

	agent := ai.NewAnonymousAgent(m, "Be helpful.")

	resp1, _ := agent.Prompt(context.Background(), "call 1")
	resp2, _ := agent.Prompt(context.Background(), "call 2")

	if resp1.GetText() != "First response" {
		t.Errorf("call 1: expected %q got %q", "First response", resp1.GetText())
	}

	if resp2.GetText() != "Second response" {
		t.Errorf("call 2: expected %q got %q", "Second response", resp2.GetText())
	}

	rec.AssertAgentWasPrompted(t, func(p *prompts.AgentPrompt) bool { return true })
}

func TestAgentLastResponseRepeats(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake("Only response")

	agent := ai.NewAnonymousAgent(m, "Be helpful.")

	resp1, _ := agent.Prompt(context.Background(), "call 1")
	resp2, _ := agent.Prompt(context.Background(), "call 2")

	if resp1.GetText() != "Only response" {
		t.Errorf("call 1: got %q", resp1.GetText())
	}

	if resp2.GetText() != "Only response" {
		t.Errorf("call 2: should repeat last response, got %q", resp2.GetText())
	}

	_ = rec
}

func TestAgentCanBeFakedWithAnAgentResponse(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	fakeResp := responses.NewAgentResponse("inv-1", "Structured response", fake.DataUsage(), fake.DataMeta())
	rec := m.Fake(fakeResp)

	agent := ai.NewAnonymousAgent(m, "Be helpful.")
	resp, err := agent.Prompt(context.Background(), "hi")

	if err != nil {
		t.Fatalf("Prompt error: %v", err)
	}

	if resp.GetText() != "Structured response" {
		t.Errorf("expected %q got %q", "Structured response", resp.GetText())
	}

	rec.AssertAgentWasPrompted(t, func(p *prompts.AgentPrompt) bool { return true })
}

func TestAgentCanBeFakedWithAClosure(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake(func(p *prompts.AgentPrompt) (*responses.AgentResponse, error) {
		return responses.NewAgentResponse("", "echo: "+p.Text, fake.DataUsage(), fake.DataMeta()), nil
	})

	agent := ai.NewAnonymousAgent(m, "Be helpful.")
	resp, err := agent.Prompt(context.Background(), "hello world")

	if err != nil {
		t.Fatalf("Prompt error: %v", err)
	}

	if resp.GetText() != "echo: hello world" {
		t.Errorf("expected %q got %q", "echo: hello world", resp.GetText())
	}

	rec.AssertAgentWasPrompted(t, func(p *prompts.AgentPrompt) bool {
		return p.Text == "hello world"
	})
}

func TestAgentAssertNeverPrompted(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake("hello")

	rec.AssertAgentNeverPrompted(t)
}

func TestAgentAssertNotPrompted(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake("hello")

	agent := ai.NewAnonymousAgent(m, "Be helpful.")
	agent.Prompt(context.Background(), "something else")

	rec.AssertAgentNotPrompted(t, func(p *prompts.AgentPrompt) bool {
		return p.Text == "specific text we never sent"
	})
}

func TestAgentCanBeQueued(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	agent := ai.NewAnonymousAgent(m, "Be helpful.")
	queued, err := agent.Queue(context.Background(), "background task")

	if err != nil {
		t.Fatalf("Queue error: %v", err)
	}

	if queued.GetInvocationID() == "" {
		t.Error("expected non-empty InvocationID")
	}

	rec.AssertAgentWasQueued(t, func(p *prompts.AgentPrompt) bool {
		return p.Text == "background task"
	})
	rec.AssertAgentNeverPrompted(t)
}

func TestAgentAssertNeverQueued(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake()

	rec.AssertAgentNeverQueued(t)
}

func TestAgentCanBeStreamed(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake("Hello stream world")

	agent := ai.NewAnonymousAgent(m, "Be helpful.")
	streamResp, err := agent.Stream(context.Background(), "stream me")

	if err != nil {
		t.Fatalf("Stream error: %v", err)
	}

	var collected []string

	streamResp.(*responses.StreamableAgentResponse).Each(func(e stream.Event) bool {
		if td, ok := e.(stream.TextDelta); ok {
			collected = append(collected, td.Delta)
		}

		return true
	})

	got := strings.Join(collected, "")

	if got == "" {
		t.Error("expected non-empty streamed text")
	}

	rec.AssertAgentWasPrompted(t, func(p *prompts.AgentPrompt) bool {
		return p.Text == "stream me"
	})
}

func TestManagerProviderTypeCheck(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()

	m.Extend("elevenlabs", func(cfg map[string]any) any {

		return &audioOnlyStub{}
	})
	m.SetDefault("elevenlabs")

	_, err := m.TextProvider()

	if err == nil {
		t.Fatal("expected ErrProviderCapability but got nil")
	}
}

func (s *audioOnlyStub) Name() string { return "elevenlabs-stub" }
