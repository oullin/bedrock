package ai

import (
	"github.com/bedrock/packages/container"

	"github.com/bedrock/packages/ai/enums"
	"github.com/bedrock/packages/ai/providers/anthropic"
	"github.com/bedrock/packages/ai/providers/azure"
	"github.com/bedrock/packages/ai/providers/cohere"
	"github.com/bedrock/packages/ai/providers/deepseek"
	"github.com/bedrock/packages/ai/providers/elevenlabs"
	"github.com/bedrock/packages/ai/providers/gemini"
	"github.com/bedrock/packages/ai/providers/groq"
	"github.com/bedrock/packages/ai/providers/jina"
	"github.com/bedrock/packages/ai/providers/mistral"
	"github.com/bedrock/packages/ai/providers/ollama"
	"github.com/bedrock/packages/ai/providers/openai"
	"github.com/bedrock/packages/ai/providers/openrouter"
	"github.com/bedrock/packages/ai/providers/voyageai"
	"github.com/bedrock/packages/ai/providers/xai"
)

// AiServiceProvider registers the AI Manager as a singleton in the container.
// It mirrors Laravel\Ai\AiServiceProvider.
type AiServiceProvider struct {
	app            *container.Container
	defaultProvider string
	configs        map[string]map[string]any
}

// NewAiServiceProvider constructs the provider.
// defaultProvider is the lab name used when no explicit provider is specified.
// configs maps lab names to their configuration (api_key, etc.).
func NewAiServiceProvider(app *container.Container, defaultProvider string, configs map[string]map[string]any) *AiServiceProvider {
	if configs == nil {
		configs = make(map[string]map[string]any)
	}
	return &AiServiceProvider{
		app:            app,
		defaultProvider: defaultProvider,
		configs:         configs,
	}
}

// Register binds the AI Manager as a singleton under "ai".
func (p *AiServiceProvider) Register() {
	p.app.Singleton("ai", func(_ *container.Container) (any, error) {
		m := NewManager()
		m.SetDefault(p.defaultProvider)

		// Register all 14 providers.
		m.Extend(string(enums.LabOpenAI), func(cfg map[string]any) any {
			return openai.NewProvider(cfg)
		})
		m.Extend(string(enums.LabAnthropic), func(cfg map[string]any) any {
			return anthropic.NewProvider(cfg)
		})
		m.Extend(string(enums.LabAzure), func(cfg map[string]any) any {
			return azure.NewProvider(cfg)
		})
		m.Extend(string(enums.LabGemini), func(cfg map[string]any) any {
			return gemini.NewProvider(cfg)
		})
		m.Extend(string(enums.LabGroq), func(cfg map[string]any) any {
			return groq.NewProvider(cfg)
		})
		m.Extend(string(enums.LabMistral), func(cfg map[string]any) any {
			return mistral.NewProvider(cfg)
		})
		m.Extend(string(enums.LabDeepSeek), func(cfg map[string]any) any {
			return deepseek.NewProvider(cfg)
		})
		m.Extend(string(enums.LabXAI), func(cfg map[string]any) any {
			return xai.NewProvider(cfg)
		})
		m.Extend(string(enums.LabCohere), func(cfg map[string]any) any {
			return cohere.NewProvider(cfg)
		})
		m.Extend(string(enums.LabOllama), func(cfg map[string]any) any {
			return ollama.NewProvider(cfg)
		})
		m.Extend(string(enums.LabElevenLabs), func(cfg map[string]any) any {
			return elevenlabs.NewProvider(cfg)
		})
		m.Extend(string(enums.LabVoyageAI), func(cfg map[string]any) any {
			return voyageai.NewProvider(cfg)
		})
		m.Extend(string(enums.LabJina), func(cfg map[string]any) any {
			return jina.NewProvider(cfg)
		})
		m.Extend(string(enums.LabOpenRouter), func(cfg map[string]any) any {
			return openrouter.NewProvider(cfg)
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
