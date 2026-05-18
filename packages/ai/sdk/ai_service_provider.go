package ai

import (
	"github.com/bedrock/packages/container"

	"github.com/bedrock/packages/ai/sdk/enums"
	"github.com/bedrock/packages/ai/sdk/providers/anthropic"
	"github.com/bedrock/packages/ai/sdk/providers/gemini"
	"github.com/bedrock/packages/ai/sdk/providers/openai"
)

// AiServiceProvider registers the AI Manager as a singleton in the container.
// It mirrors Upstream\Ai\AiServiceProvider.
type AiServiceProvider struct {
	app             *container.Container
	defaultProvider string
	configs         map[string]map[string]any
}

// NewAiServiceProvider constructs the provider.
// defaultProvider is the lab name used when no explicit provider is specified.
// configs maps lab names to their configuration (api_key, etc.).
func NewAiServiceProvider(app *container.Container, defaultProvider string, configs map[string]map[string]any) *AiServiceProvider {
	if configs == nil {
		configs = make(map[string]map[string]any)
	}

	return &AiServiceProvider{
		app:             app,
		defaultProvider: defaultProvider,
		configs:         configs,
	}
}

// Register binds the AI Manager as a singleton under "ai".
func (p *AiServiceProvider) Register() {
	p.app.Singleton("ai", func(_ *container.Container) (any, error) {
		m := NewManager()
		m.SetDefault(p.defaultProvider)

		m.Extend(string(enums.LabOpenAI), func(cfg map[string]any) any {
			return openai.NewProvider(cfg)
		})
		m.Extend(string(enums.LabAnthropic), func(cfg map[string]any) any {
			return anthropic.NewProvider(cfg)
		})
		m.Extend(string(enums.LabGemini), func(cfg map[string]any) any {
			return gemini.NewProvider(cfg)
		})

		// Apply per-provider configurations.
		for lab, cfg := range p.configs {
			m.Configure(lab, cfg)
		}

		// Wire the Manager into the global accessor so package-level functions work.
		SetManager(m)

		return m, nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *AiServiceProvider) Provides() []string { return []string{"ai"} }
