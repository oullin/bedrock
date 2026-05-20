// Port of \Ai\Tests\Feature\AiManagerTest
package ai_test

import (
	"context"
	"testing"

	ai "github.com/bedrock/packages/ai/sdk"
	"github.com/bedrock/packages/ai/sdk/enums"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"
)

// textOnlyStub is a minimal TextProvider stub for Manager tests.
type textOnlyStub struct{ name string }

func (s *textOnlyStub) Prompt(_ context.Context, _ contractsprovider.TextPromptRequest) (*contractsprovider.TextPromptResult, error) {
	return &contractsprovider.TextPromptResult{}, nil
}
func (s *textOnlyStub) Stream(_ context.Context, _ contractsprovider.TextPromptRequest) (contractsprovider.StreamTextResult, error) {
	return contractsprovider.StreamTextResult{}, nil
}
func (s *textOnlyStub) TextGateway() contractsgw.TextGateway { return nil }
func (s *textOnlyStub) UseTextGateway(_ contractsgw.TextGateway) contractsprovider.TextProvider {
	return s
}
func (s *textOnlyStub) DefaultTextModel() string  { return "stub" }
func (s *textOnlyStub) CheapestTextModel() string { return "stub-cheap" }
func (s *textOnlyStub) SmartestTextModel() string { return "stub-smart" }

// TestManagerSetAndGetDefault verifies SetDefault / Default round-trips.
func TestManagerSetAndGetDefault(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	m.SetDefault("openai")

	if m.Default() != "openai" {
		t.Errorf("expected %q got %q", "openai", m.Default())
	}
}

// TestManagerExtendAndResolve verifies that a factory registered via Extend is used.
func TestManagerExtendAndResolve(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	called := false
	m.Extend("test-provider", func(cfg map[string]any) any {
		called = true

		return &textOnlyStub{}
	})
	m.SetDefault("test-provider")

	_, err := m.TextProvider()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !called {
		t.Error("expected factory to be called")
	}
}

// TestManagerProviderCached verifies the same instance is returned on repeated calls.
func TestManagerProviderCached(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	m.Extend("test-provider", func(cfg map[string]any) any {
		return &textOnlyStub{}
	})
	m.SetDefault("test-provider")

	p1, _ := m.TextProvider()
	p2, _ := m.TextProvider()

	if p1 != p2 {
		t.Error("expected cached provider instance to be reused")
	}
}

// TestManagerTextProviderCapabilityError verifies ErrProviderCapability for wrong type.
func TestManagerTextProviderCapabilityError(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	m.Extend("audio-only", func(cfg map[string]any) any {
		return &audioOnlyStub{}
	})
	m.SetDefault("audio-only")

	_, err := m.TextProvider()

	if err == nil {
		t.Fatal("expected ErrProviderCapability but got nil")
	}
}

// TestManagerFakeRecorderReturned verifies Fake() returns a non-nil Recorder.
func TestManagerFakeRecorderReturned(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	rec := m.Fake("hello")

	if rec == nil {
		t.Fatal("expected non-nil recorder from Fake()")
	}
}

// TestManagerResetClearsInstances verifies Reset() clears cached providers.
func TestManagerResetClearsInstances(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	count := 0
	m.Extend("cnt", func(cfg map[string]any) any {
		count++

		return &textOnlyStub{}
	})
	m.SetDefault("cnt")

	m.TextProvider() //nolint:errcheck
	m.Reset()
	m.TextProvider() //nolint:errcheck

	if count != 2 {
		t.Errorf("expected factory called twice after reset, got %d", count)
	}
}

// TestManagerNamedProvider verifies that a named lab can be specified.
func TestManagerNamedProvider(t *testing.T) {
	t.Parallel()

	m := ai.NewManager()
	m.Extend("alpha", func(cfg map[string]any) any { return &textOnlyStub{name: "alpha"} })
	m.Extend("beta", func(cfg map[string]any) any { return &textOnlyStub{name: "beta"} })
	m.SetDefault("alpha")

	p, err := m.TextProvider(enums.Lab("beta"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := p.(*textOnlyStub); !ok {
		t.Errorf("expected textOnlyStub, got %T", p)
	}
}
