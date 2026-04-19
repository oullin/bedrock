package ai

import (
	"fmt"
	"sync"

	contractsai "github.com/bedrock/packages/contracts/ai"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"

	"github.com/bedrock/packages/ai/sdk/enums"
	"github.com/bedrock/packages/ai/sdk/fake"
	"github.com/bedrock/packages/ai/sdk/prompts"
)

// DriverFactory is a function that constructs a concrete provider from a config map.
// Each concrete provider package registers one of these with the Manager.
type DriverFactory func(config map[string]any) any

// Manager resolves AI providers by capability and lab name.
// It mirrors Upstream\Ai\AiManager.
type Manager struct {
	mu         sync.RWMutex
	configs    map[string]map[string]any // lab → config
	factories  map[string]DriverFactory  // lab → constructor
	instances  map[string]any            // lab → cached concrete provider
	defaultLab string

	// Shared fake state — nil when not in fake mode.
	recorder *fake.Recorder
}

// NewManager constructs an unconfigured Manager.
func NewManager() *Manager {
	return &Manager{
		configs:   make(map[string]map[string]any),
		factories: make(map[string]DriverFactory),
		instances: make(map[string]any),
	}
}

// SetDefault sets the default provider lab name.
func (m *Manager) SetDefault(lab string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.defaultLab = lab
}

// Default returns the default provider lab name.
func (m *Manager) Default() string {
	m.mu.RLock()

	defer m.mu.RUnlock()

	return m.defaultLab
}

// Extend registers a driver factory for the given lab name.
func (m *Manager) Extend(lab string, factory DriverFactory) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.factories[lab] = factory
}

// Configure stores provider configuration for the given lab.
func (m *Manager) Configure(lab string, config map[string]any) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.configs[lab] = config
}

// RecordQueuedAgent records an agent prompt as queued in the shared recorder.
// Called by AnonymousAgent.Queue when in fake mode; no-op otherwise.
func (m *Manager) RecordQueuedAgent(p *prompts.AgentPrompt) {
	m.mu.Lock()
	rec := m.recorder
	m.mu.Unlock()

	if rec != nil {
		rec.RecordQueuedPrompt(p)
	}
}

// Reset clears all cached instances and fake state.
func (m *Manager) Reset() {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.instances = make(map[string]any)
	m.recorder = nil
}

// Recorder returns the shared fake recorder, creating one if needed.
// This is used by static accessors to run assertions.
func (m *Manager) Recorder() *fake.Recorder {
	m.mu.Lock()

	defer m.mu.Unlock()

	return m.ensureRecorderLocked()
}

// --- internal resolution ---

// resolve returns the cached or newly-constructed provider for lab.
// Callers must NOT hold m.mu.
func (m *Manager) resolve(lab string) (any, error) {
	m.mu.Lock()

	defer m.mu.Unlock()

	return m.resolveLocked(lab)
}

// resolveLocked creates or returns a cached provider. Must be called with m.mu held.
func (m *Manager) resolveLocked(lab string) (any, error) {
	if inst, ok := m.instances[lab]; ok {
		return inst, nil
	}

	factory, ok := m.factories[lab]

	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedProvider, lab)
	}

	cfg := m.configs[lab]

	if cfg == nil {
		cfg = make(map[string]any)
	}

	inst := factory(cfg)
	m.instances[lab] = inst

	return inst, nil
}

func (m *Manager) ensureRecorderLocked() *fake.Recorder {
	if m.recorder == nil {
		m.recorder = fake.NewRecorder()
	}

	return m.recorder
}

// labFor returns the lab name to use, falling back to default.
func (m *Manager) labFor(labs []enums.Lab) string {
	if len(labs) > 0 {
		return string(labs[0])
	}

	return m.Default()
}

// --- Text capability ---

// TextProvider resolves a TextProvider for the given lab (defaults to default lab).
// Returns ErrProviderCapability if the provider does not support text generation.
func (m *Manager) TextProvider(labs ...enums.Lab) (contractsprovider.TextProvider, error) {
	lab := m.labFor(labs)
	inst, err := m.resolve(lab)

	if err != nil {
		return nil, err
	}

	tp, ok := inst.(contractsprovider.TextProvider)

	if !ok {
		return nil, fmt.Errorf("%w: %q does not support text generation", ErrProviderCapability, lab)
	}

	return tp, nil
}

// TextProviderFor resolves a TextProvider honouring the agent's provider hint.
// If the agent implements HasProviderOptions and provides a "provider" key,
// that lab is used; otherwise the Manager default is used.
func (m *Manager) TextProviderFor(agent contractsai.Agent) (contractsprovider.TextProvider, error) {
	lab := m.Default()

	if opts, ok := agent.(contractsai.HasProviderOptions); ok {
		if p, ok := opts.ProviderOptions()["provider"].(string); ok && p != "" {
			lab = p
		}
	}

	return m.TextProvider(enums.Lab(lab))
}

// FakeTextProvider injects a fake text gateway into the default provider.
// The gateway is pre-loaded with responses (string | *responses.AgentResponse | func).
// Returns the provider as TextProvider so callers can chain usage.
func (m *Manager) FakeTextProvider(responses ...any) contractsprovider.TextProvider {
	m.mu.Lock()
	rec := m.ensureRecorderLocked()
	inst, err := m.resolveLocked(m.defaultLab)
	m.mu.Unlock()

	if err != nil {
		panic(fmt.Sprintf("ai: FakeTextProvider: cannot resolve default provider %q: %v", m.defaultLab, err))
	}

	gw := fake.NewTextGateway(rec)

	if len(responses) > 0 {
		gw.SetResponses(responses)
	}

	tp, ok := inst.(contractsprovider.TextProvider)

	if !ok {
		panic(fmt.Sprintf("ai: FakeTextProvider: default provider %q does not support text", m.defaultLab))
	}

	return tp.UseTextGateway(gw)
}

// --- Image capability ---

// ImageProvider resolves an ImageProvider for the given lab.
func (m *Manager) ImageProvider(labs ...enums.Lab) (contractsprovider.ImageProvider, error) {
	lab := m.labFor(labs)
	inst, err := m.resolve(lab)

	if err != nil {
		return nil, err
	}

	ip, ok := inst.(contractsprovider.ImageProvider)

	if !ok {
		return nil, fmt.Errorf("%w: %q does not support image generation", ErrProviderCapability, lab)
	}

	return ip, nil
}

// FakeImageProvider injects a fake image gateway into the default provider.
func (m *Manager) FakeImageProvider(responses ...any) contractsprovider.ImageProvider {
	m.mu.Lock()
	rec := m.ensureRecorderLocked()
	inst, err := m.resolveLocked(m.defaultLab)
	m.mu.Unlock()

	if err != nil {
		panic(fmt.Sprintf("ai: FakeImageProvider: cannot resolve default provider: %v", err))
	}

	gw := fake.NewImageGateway(rec)

	if len(responses) > 0 {
		gw.SetResponses(responses)
	}

	ip, ok := inst.(contractsprovider.ImageProvider)

	if !ok {
		panic(fmt.Sprintf("ai: FakeImageProvider: default provider %q does not support images", m.defaultLab))
	}

	return ip.UseImageGateway(gw)
}

// --- Audio capability ---

// AudioProvider resolves an AudioProvider for the given lab.
func (m *Manager) AudioProvider(labs ...enums.Lab) (contractsprovider.AudioProvider, error) {
	lab := m.labFor(labs)
	inst, err := m.resolve(lab)

	if err != nil {
		return nil, err
	}

	ap, ok := inst.(contractsprovider.AudioProvider)

	if !ok {
		return nil, fmt.Errorf("%w: %q does not support audio generation", ErrProviderCapability, lab)
	}

	return ap, nil
}

// FakeAudioProvider injects a fake audio gateway into the default provider.
func (m *Manager) FakeAudioProvider(responses ...any) contractsprovider.AudioProvider {
	m.mu.Lock()
	rec := m.ensureRecorderLocked()
	inst, err := m.resolveLocked(m.defaultLab)
	m.mu.Unlock()

	if err != nil {
		panic(fmt.Sprintf("ai: FakeAudioProvider: cannot resolve default provider: %v", err))
	}

	gw := fake.NewAudioGateway(rec)

	if len(responses) > 0 {
		gw.SetResponses(responses)
	}

	ap, ok := inst.(contractsprovider.AudioProvider)

	if !ok {
		panic(fmt.Sprintf("ai: FakeAudioProvider: default provider %q does not support audio", m.defaultLab))
	}

	return ap.UseAudioGateway(gw)
}

// --- Embedding capability ---

// EmbeddingProvider resolves an EmbeddingProvider for the given lab.
func (m *Manager) EmbeddingProvider(labs ...enums.Lab) (contractsprovider.EmbeddingProvider, error) {
	lab := m.labFor(labs)
	inst, err := m.resolve(lab)

	if err != nil {
		return nil, err
	}

	ep, ok := inst.(contractsprovider.EmbeddingProvider)

	if !ok {
		return nil, fmt.Errorf("%w: %q does not support embeddings", ErrProviderCapability, lab)
	}

	return ep, nil
}

// FakeEmbeddingProvider injects a fake embedding gateway into the default provider.
func (m *Manager) FakeEmbeddingProvider(responses ...any) contractsprovider.EmbeddingProvider {
	m.mu.Lock()
	rec := m.ensureRecorderLocked()
	inst, err := m.resolveLocked(m.defaultLab)
	m.mu.Unlock()

	if err != nil {
		panic(fmt.Sprintf("ai: FakeEmbeddingProvider: cannot resolve default provider: %v", err))
	}

	gw := fake.NewEmbeddingGateway(rec)

	if len(responses) > 0 {
		gw.SetResponses(responses)
	}

	ep, ok := inst.(contractsprovider.EmbeddingProvider)

	if !ok {
		panic(fmt.Sprintf("ai: FakeEmbeddingProvider: default provider %q does not support embeddings", m.defaultLab))
	}

	return ep.UseEmbeddingGateway(gw)
}

// --- Transcription capability ---

// TranscriptionProvider resolves a TranscriptionProvider for the given lab.
func (m *Manager) TranscriptionProvider(labs ...enums.Lab) (contractsprovider.TranscriptionProvider, error) {
	lab := m.labFor(labs)
	inst, err := m.resolve(lab)

	if err != nil {
		return nil, err
	}

	tp, ok := inst.(contractsprovider.TranscriptionProvider)

	if !ok {
		return nil, fmt.Errorf("%w: %q does not support transcription", ErrProviderCapability, lab)
	}

	return tp, nil
}

// FakeTranscriptionProvider injects a fake transcription gateway into the default provider.
func (m *Manager) FakeTranscriptionProvider(responses ...any) contractsprovider.TranscriptionProvider {
	m.mu.Lock()
	rec := m.ensureRecorderLocked()
	inst, err := m.resolveLocked(m.defaultLab)
	m.mu.Unlock()

	if err != nil {
		panic(fmt.Sprintf("ai: FakeTranscriptionProvider: cannot resolve default provider: %v", err))
	}

	gw := fake.NewTranscriptionGateway(rec)

	if len(responses) > 0 {
		gw.SetResponses(responses)
	}

	tp, ok := inst.(contractsprovider.TranscriptionProvider)

	if !ok {
		panic(fmt.Sprintf("ai: FakeTranscriptionProvider: default provider %q does not support transcription", m.defaultLab))
	}

	return tp.UseTranscriptionGateway(gw)
}

// --- Reranking capability ---

// RerankingProvider resolves a RerankingProvider for the given lab.
func (m *Manager) RerankingProvider(labs ...enums.Lab) (contractsprovider.RerankingProvider, error) {
	lab := m.labFor(labs)
	inst, err := m.resolve(lab)

	if err != nil {
		return nil, err
	}

	rp, ok := inst.(contractsprovider.RerankingProvider)

	if !ok {
		return nil, fmt.Errorf("%w: %q does not support reranking", ErrProviderCapability, lab)
	}

	return rp, nil
}

// FakeRerankingProvider injects a fake reranking gateway into the default provider.
func (m *Manager) FakeRerankingProvider(responses ...any) contractsprovider.RerankingProvider {
	m.mu.Lock()
	rec := m.ensureRecorderLocked()
	inst, err := m.resolveLocked(m.defaultLab)
	m.mu.Unlock()

	if err != nil {
		panic(fmt.Sprintf("ai: FakeRerankingProvider: cannot resolve default provider: %v", err))
	}

	gw := fake.NewRerankingGateway(rec)

	if len(responses) > 0 {
		gw.SetResponses(responses)
	}

	rp, ok := inst.(contractsprovider.RerankingProvider)

	if !ok {
		panic(fmt.Sprintf("ai: FakeRerankingProvider: default provider %q does not support reranking", m.defaultLab))
	}

	return rp.UseRerankingGateway(gw)
}

// --- File capability ---

// FileProvider resolves a FileProvider for the given lab.
func (m *Manager) FileProvider(labs ...enums.Lab) (contractsprovider.FileProvider, error) {
	lab := m.labFor(labs)
	inst, err := m.resolve(lab)

	if err != nil {
		return nil, err
	}

	fp, ok := inst.(contractsprovider.FileProvider)

	if !ok {
		return nil, fmt.Errorf("%w: %q does not support file management", ErrProviderCapability, lab)
	}

	return fp, nil
}

// FakeFileProvider injects a fake file gateway into the default provider.
func (m *Manager) FakeFileProvider() contractsprovider.FileProvider {
	m.mu.Lock()
	rec := m.ensureRecorderLocked()
	inst, err := m.resolveLocked(m.defaultLab)
	m.mu.Unlock()

	if err != nil {
		panic(fmt.Sprintf("ai: FakeFileProvider: cannot resolve default provider: %v", err))
	}

	gw := fake.NewFileGateway(rec)

	fp, ok := inst.(contractsprovider.FileProvider)

	if !ok {
		panic(fmt.Sprintf("ai: FakeFileProvider: default provider %q does not support files", m.defaultLab))
	}

	return fp.UseFileGateway(gw)
}

// --- Store capability ---

// StoreProvider resolves a StoreProvider for the given lab.
func (m *Manager) StoreProvider(labs ...enums.Lab) (contractsprovider.StoreProvider, error) {
	lab := m.labFor(labs)
	inst, err := m.resolve(lab)

	if err != nil {
		return nil, err
	}

	sp, ok := inst.(contractsprovider.StoreProvider)

	if !ok {
		return nil, fmt.Errorf("%w: %q does not support vector stores", ErrProviderCapability, lab)
	}

	return sp, nil
}

// FakeStoreProvider injects a fake store gateway into the default provider.
func (m *Manager) FakeStoreProvider() contractsprovider.StoreProvider {
	m.mu.Lock()
	rec := m.ensureRecorderLocked()
	inst, err := m.resolveLocked(m.defaultLab)
	m.mu.Unlock()

	if err != nil {
		panic(fmt.Sprintf("ai: FakeStoreProvider: cannot resolve default provider: %v", err))
	}

	gw := fake.NewStoreGateway(rec)

	sp, ok := inst.(contractsprovider.StoreProvider)

	if !ok {
		panic(fmt.Sprintf("ai: FakeStoreProvider: default provider %q does not support stores", m.defaultLab))
	}

	return sp.UseStoreGateway(gw)
}

// --- Fake all ---

// Fake configures fake gateways for ALL capabilities on the default provider
// (whichever the provider supports). Returns the shared Recorder for assertions.
// textResponses are pre-queued for the text gateway.
// If no provider has been registered, a stub provider is automatically registered.
func (m *Manager) Fake(textResponses ...any) *fake.Recorder {
	m.mu.Lock()
	// Auto-register a stub when no provider is configured (common in unit tests).
	if m.defaultLab == "" || len(m.factories) == 0 {
		if _, ok := m.factories["stub"]; !ok {
			m.factories["stub"] = newStubProvider
		}

		m.defaultLab = "stub"
	}

	rec := m.ensureRecorderLocked()
	inst, err := m.resolveLocked(m.defaultLab)
	m.mu.Unlock()

	if err != nil {
		panic(fmt.Sprintf("ai: Fake: cannot resolve default provider %q: %v", m.defaultLab, err))
	}

	// Inject fake gateways for each supported capability.
	if tp, ok := inst.(contractsprovider.TextProvider); ok {
		gw := fake.NewTextGateway(rec)

		if len(textResponses) > 0 {
			gw.SetResponses(textResponses)
		}

		tp.UseTextGateway(gw)
	}

	if ip, ok := inst.(contractsprovider.ImageProvider); ok {
		ip.UseImageGateway(fake.NewImageGateway(rec))
	}

	if ap, ok := inst.(contractsprovider.AudioProvider); ok {
		ap.UseAudioGateway(fake.NewAudioGateway(rec))
	}

	if ep, ok := inst.(contractsprovider.EmbeddingProvider); ok {
		ep.UseEmbeddingGateway(fake.NewEmbeddingGateway(rec))
	}

	if tp, ok := inst.(contractsprovider.TranscriptionProvider); ok {
		tp.UseTranscriptionGateway(fake.NewTranscriptionGateway(rec))
	}

	if rp, ok := inst.(contractsprovider.RerankingProvider); ok {
		rp.UseRerankingGateway(fake.NewRerankingGateway(rec))
	}

	if fp, ok := inst.(contractsprovider.FileProvider); ok {
		fp.UseFileGateway(fake.NewFileGateway(rec))
	}

	if sp, ok := inst.(contractsprovider.StoreProvider); ok {
		sp.UseStoreGateway(fake.NewStoreGateway(rec))
	}

	return rec
}
