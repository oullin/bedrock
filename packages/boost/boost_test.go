package boost_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/boost"
)

// TestManagerDefaultAgents mirrors BoostManagerTest::test_has_default_agents.
func TestManagerDefaultAgents(t *testing.T) {
	t.Parallel()

	m := boost.New()
	got := m.GetAgents()

	wantKeys := []string{
		"amp", "junie", "cursor", "claude_code",
		"codex", "copilot", "kiro", "opencode", "gemini",
	}

	if len(got) != len(wantKeys) {
		t.Fatalf("want %d default agents, got %d", len(wantKeys), len(got))
	}

	for _, k := range wantKeys {
		if _, ok := got[k]; !ok {
			t.Errorf("missing default agent key %q", k)
		}
	}
}

// TestManagerRegisterAgent mirrors BoostManagerTest::test_can_register_agent.
func TestManagerRegisterAgent(t *testing.T) {
	t.Parallel()

	m := boost.New()

	stub := &stubAgent{name: "my_agent"}
	if err := m.RegisterAgent("my_agent", stub); err != nil {
		t.Fatalf("RegisterAgent: unexpected error: %v", err)
	}

	agents := m.GetAgents()
	if _, ok := agents["my_agent"]; !ok {
		t.Error("registered agent not found via GetAgents")
	}
}

// TestManagerRegisterAgentDuplicate mirrors BoostManagerTest::test_duplicate_key_returns_error.
func TestManagerRegisterAgentDuplicate(t *testing.T) {
	t.Parallel()

	m := boost.New()

	// "claude_code" is a default key; attempting to register it again must fail.
	err := m.RegisterAgent("claude_code", &stubAgent{name: "claude_code"})
	if err == nil {
		t.Fatal("expected error for duplicate key, got nil")
	}

	if !errors.Is(err, boost.ErrAgentAlreadyRegistered) {
		t.Errorf("want ErrAgentAlreadyRegistered, got %v", err)
	}
}

// TestManagerGetAgentsSnapshot mirrors BoostManagerTest::test_get_agents_returns_snapshot.
func TestManagerGetAgentsSnapshot(t *testing.T) {
	t.Parallel()

	m := boost.New()
	snap1 := m.GetAgents()

	if err := m.RegisterAgent("extra", &stubAgent{name: "extra"}); err != nil {
		t.Fatalf("RegisterAgent: %v", err)
	}

	snap2 := m.GetAgents()

	// snap1 must not include "extra" (it's a snapshot taken before registration).
	if _, ok := snap1["extra"]; ok {
		t.Error("snapshot 1 should not contain 'extra'")
	}

	if _, ok := snap2["extra"]; !ok {
		t.Error("snapshot 2 should contain 'extra'")
	}
}

// ---------------------------------------------------------------------------
// stubAgent satisfies boost.CodingAgent for test purposes only.
// ---------------------------------------------------------------------------

type stubAgent struct {
	name string
}

func (s *stubAgent) Name() string                         { return s.name }
func (s *stubAgent) DisplayName() string                  { return s.name }
func (s *stubAgent) McpConfigPath() string                { return "" }
func (s *stubAgent) McpConfigKey() string                 { return "mcpServers" }
func (s *stubAgent) ShellMcpCommand() string              { return "" }
func (s *stubAgent) DefaultMcpConfig() map[string]any    { return nil }
func (s *stubAgent) Frontmatter() bool                    { return false }
func (s *stubAgent) McpInstallationStrategy() boost.McpInstallationStrategy {
	return boost.McpStrategyNone
}
func (s *stubAgent) DetectOnSystem(_ boost.Platform) bool    { return false }
func (s *stubAgent) DetectInProject(_ string) bool           { return false }
func (s *stubAgent) InstallMcp(_, _ string, _ []string, _ map[string]string) (bool, error) {
	return false, nil
}
func (s *stubAgent) InstallHttpMcp(_, _ string) (bool, error) { return false, nil }
func (s *stubAgent) HttpMcpServerConfig(_ string) map[string]any {
	return map[string]any{"type": "http"}
}
func (s *stubAgent) McpServerConfig(_ string, _ []string, _ map[string]string) map[string]any {
	return nil
}
func (s *stubAgent) UseAbsolutePathForMcp() bool              { return false }
func (s *stubAgent) GoBinaryPath(_ bool) string               { return "go" }
func (s *stubAgent) EntryPointPath(_ bool) string             { return "main.go" }
func (s *stubAgent) TransformGuidelines(md string) string     { return md }
