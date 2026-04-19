// Package ai provides the global singleton accessor pattern for the AI Manager.
// These package-level functions mirror the Laravel AI facade's static API.
package ai

import (
	"sync"

	contractsai "github.com/bedrock/packages/contracts/ai"

	"github.com/bedrock/packages/ai/sdk/fake"
	"github.com/bedrock/packages/ai/sdk/prompts"
)

var (
	globalMu  sync.Mutex
	globalMgr *Manager
)

// globalManager returns the package-level singleton Manager, creating it if needed.
func globalManager() *Manager {
	globalMu.Lock()

	defer globalMu.Unlock()

	if globalMgr == nil {
		globalMgr = NewManager()
	}

	return globalMgr
}

// SetManager replaces the global Manager.
// Called by AiServiceProvider.Register so the whole application shares one Manager.
func SetManager(m *Manager) {
	globalMu.Lock()

	defer globalMu.Unlock()

	globalMgr = m
}

// GetManager returns the current global Manager.
func GetManager() *Manager {
	return globalManager()
}

// Fake sets up fake gateways for all capabilities and returns a Recorder for assertions.
// In-test usage: defer ai.Reset() to clean up.
func Fake(textResponses ...any) *fake.Recorder {
	return globalManager().Fake(textResponses...)
}

// FakeText queues text responses and injects a fake text gateway.
func FakeText(responses ...any) {
	globalManager().FakeTextProvider(responses...)
}

// Reset clears all cached provider instances and fake state.
func Reset() {
	globalManager().Reset()
}

// NewAgent creates a new AnonymousAgent using the global Manager.
func NewAgent(instructions string) *AnonymousAgent {
	return NewAnonymousAgent(globalManager(), instructions)
}

// NewStructuredAgent creates a new StructuredAnonymousAgent using the global Manager.
func NewStructuredAgent(instructions string, schema contractsai.JsonSchema) *StructuredAnonymousAgent {
	return NewStructuredAnonymousAgent(globalManager(), instructions, schema)
}

// --- Agent assertion helpers ---

// AssertAgentWasPrompted fails t if no recorded prompt matches fn.
func AssertAgentWasPrompted(t fake.TestingT, fn func(*prompts.AgentPrompt) bool) {
	globalManager().Recorder().AssertAgentWasPrompted(t, fn)
}

// AssertAgentNotPrompted fails t if any recorded prompt matches fn.
func AssertAgentNotPrompted(t fake.TestingT, fn func(*prompts.AgentPrompt) bool) {
	globalManager().Recorder().AssertAgentNotPrompted(t, fn)
}

// AssertAgentNeverPrompted fails t if any prompt was ever recorded.
func AssertAgentNeverPrompted(t fake.TestingT) {
	globalManager().Recorder().AssertAgentNeverPrompted(t)
}

// AssertAgentWasQueued fails t if no queued prompt matches fn.
func AssertAgentWasQueued(t fake.TestingT, fn func(*prompts.AgentPrompt) bool) {
	globalManager().Recorder().AssertAgentWasQueued(t, fn)
}

// AssertAgentNotQueued fails t if any queued prompt matches fn.
func AssertAgentNotQueued(t fake.TestingT, fn func(*prompts.AgentPrompt) bool) {
	globalManager().Recorder().AssertAgentNotQueued(t, fn)
}

// AssertAgentNeverQueued fails t if any prompt was ever queued.
func AssertAgentNeverQueued(t fake.TestingT) {
	globalManager().Recorder().AssertAgentNeverQueued(t)
}
