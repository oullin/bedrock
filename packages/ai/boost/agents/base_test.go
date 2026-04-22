package agents_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/ai/boost"
	"github.com/bedrock/packages/ai/boost/agents"
)

func TestBaseAgentInstallMcpMissingConfigPathReturnsSharedSentinel(t *testing.T) {
	t.Parallel()

	a := agents.NewBaseAgent(agents.AgentOptions{})

	_, err := a.InstallMcp("boost", "go", []string{"run", "."}, nil)

	if !errors.Is(err, agents.ErrNoMcpConfigPath) {
		t.Fatalf("InstallMcp() error = %v, want %v", err, agents.ErrNoMcpConfigPath)
	}

	if !errors.Is(err, boost.ErrNoMcpConfigPath) {
		t.Fatalf("InstallMcp() error = %v, want %v", err, boost.ErrNoMcpConfigPath)
	}
}
