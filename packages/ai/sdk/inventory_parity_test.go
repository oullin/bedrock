package ai_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	ai "github.com/bedrock/packages/ai/sdk"
	"github.com/bedrock/packages/ai/sdk/data"
	"github.com/bedrock/packages/ai/sdk/enums"
	"github.com/bedrock/packages/ai/sdk/fake"
	"github.com/bedrock/packages/ai/sdk/prompts"
	anthropic "github.com/bedrock/packages/ai/sdk/providers/anthropic"
	azure "github.com/bedrock/packages/ai/sdk/providers/azure"
	gemini "github.com/bedrock/packages/ai/sdk/providers/gemini"
	groq "github.com/bedrock/packages/ai/sdk/providers/groq"
	mistral "github.com/bedrock/packages/ai/sdk/providers/mistral"
	ollama "github.com/bedrock/packages/ai/sdk/providers/ollama"
	openai "github.com/bedrock/packages/ai/sdk/providers/openai"
	openrouter "github.com/bedrock/packages/ai/sdk/providers/openrouter"
	xai "github.com/bedrock/packages/ai/sdk/providers/xai"
	"github.com/bedrock/packages/ai/sdk/responses"
	"github.com/bedrock/packages/ai/sdk/stream"
	contractsai "github.com/bedrock/packages/contracts/ai"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"
)

type recordingT struct {
	failed bool
}

func (t *recordingT) Helper() {}

func (t *recordingT) Errorf(string, ...any) {
	t.failed = true
}

func TestInventoryAgentFakeCoreParity(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("faked prompt records and returns response", func(t *testing.T) {
		// AgentFakeTest::test_agents_can_be_faked
		// AgentFakeTest::test_can_assert_agent_was_never_prompted
		// AgentFakeTest::test_timeout_defaults_to_sdk_default_when_not_provided
		m := ai.NewManager()
		rec := m.Fake("fake response")
		rec.AssertAgentNeverPrompted(t)

		resp, err := ai.NewAnonymousAgent(m, "Be precise.").Prompt(ctx, "hello")
		if err != nil {
			t.Fatalf("Prompt error: %v", err)
		}
		if resp.GetText() != "fake response" {
			t.Fatalf("response text = %q, want fake response", resp.GetText())
		}
		rec.AssertAgentWasPrompted(t, func(p *prompts.AgentPrompt) bool {
			return p.Text == "hello" && p.Timeout == fake.DefaultTextTimeout
		})
	})

	t.Run("empty and closure responses", func(t *testing.T) {
		// AgentFakeTest::test_agents_can_be_faked_with_no_predefined_responses
		// AgentFakeTest::test_agents_can_be_faked_with_a_single_closure_that_is_invoked_for_every_prompt
		m := ai.NewManager()
		m.Fake()

		empty, err := ai.NewAnonymousAgent(m, "Be brief.").Prompt(ctx, "empty")
		if err != nil {
			t.Fatalf("Prompt error: %v", err)
		}
		if empty.GetText() != "" {
			t.Fatalf("empty fake text = %q, want empty string", empty.GetText())
		}

		rec := m.Fake(func(p *prompts.AgentPrompt) (*responses.AgentResponse, error) {
			return responses.NewAgentResponse(p.InvocationID, "broadcastclient: "+p.Text, data.Usage{}, data.Meta{}), nil
		})
		got, err := ai.NewAnonymousAgent(m, "BroadcastClient.").Prompt(ctx, "closure")
		if err != nil {
			t.Fatalf("Prompt closure error: %v", err)
		}
		if got.GetText() != "broadcastclient: closure" {
			t.Fatalf("closure text = %q", got.GetText())
		}
		rec.AssertAgentWasPrompted(t, func(p *prompts.AgentPrompt) bool {
			return p.Text == "closure"
		})
	})

	t.Run("stray and failing closures", func(t *testing.T) {
		// AgentFakeTest::test_agents_can_prevent_stray_prompts
		// AgentFakeTest::test_fake_closures_can_throw_exceptions
		m := ai.NewManager()
		m.Fake()
		provider, err := m.TextProvider()
		if err != nil {
			t.Fatalf("TextProvider error: %v", err)
		}
		gw := fake.NewTextGateway(m.Recorder())
		gw.PreventStray()
		provider.UseTextGateway(gw)

		_, err = ai.NewAnonymousAgent(m, "No stray.").Prompt(ctx, "unexpected")
		if err == nil {
			t.Fatal("expected stray prompt error")
		}

		expected := errors.New("fake failed")
		m.Fake(func(*prompts.AgentPrompt) (*responses.AgentResponse, error) {
			return nil, expected
		})
		_, err = ai.NewAnonymousAgent(m, "Fail.").Prompt(ctx, "boom")
		if !errors.Is(err, expected) {
			t.Fatalf("Prompt error = %v, want %v", err, expected)
		}
	})

	t.Run("structured and streaming agents", func(t *testing.T) {
		// AgentFakeTest::test_agents_with_structured_output_can_be_faked
		// AgentFakeTest::test_agents_with_structured_output_can_be_faked_with_no_predefined_responses
		// AgentFakeTest::test_structured_agents_with_empty_schemas_fall_back_to_a_text_response
		// AgentFakeTest::test_agent_streams_can_be_faked
		m := ai.NewManager()
		m.Fake(`{"name":"Ada"}`)
		agent := ai.NewStructuredAnonymousAgent(m, "Return JSON.", contractsai.JsonSchema{"type": "object"})
		structured, err := agent.Prompt(ctx, "profile")
		if err != nil {
			t.Fatalf("Structured prompt error: %v", err)
		}
		if structured.ToMap()["name"] != "Ada" {
			t.Fatalf("structured name = %#v", structured.ToMap()["name"])
		}

		m.Fake()
		empty, err := ai.NewStructuredAnonymousAgent(m, "Return JSON.", nil).Prompt(ctx, "empty")
		if err != nil {
			t.Fatalf("Structured empty prompt error: %v", err)
		}
		if len(empty.ToMap()) != 0 || empty.GetText() != "" {
			t.Fatalf("empty structured response = %#v text %q", empty.ToMap(), empty.GetText())
		}

		m.Fake("stream response")
		streamed, err := ai.NewAnonymousAgent(m, "Stream.").Stream(ctx, "stream")
		if err != nil {
			t.Fatalf("Stream error: %v", err)
		}
		var b strings.Builder
		streamed.(*responses.StreamableAgentResponse).Each(func(e stream.Event) bool {
			if delta, ok := e.(stream.TextDelta); ok {
				b.WriteString(delta.Delta)
			}
			return true
		})
		if b.String() != "stream response" {
			t.Fatalf("stream text = %q", b.String())
		}
	})

	t.Run("queued and provider options", func(t *testing.T) {
		// AgentFakeTest::test_queued_agents_can_be_faked
		// AgentFakeTest::test_can_assert_agent_was_never_queued
		// AgentFakeTest::test_assert_queued_does_not_throw_undefined_key_when_agent_was_never_queued
		// AgentFakeTest::test_assert_not_queued_does_not_throw_undefined_key_when_agent_was_never_queued
		// AgentFakeTest::test_queued_agents_accept_ai_provider_enum
		// AgentFakeTest::test_prompt_accepts_ai_provider_enum
		// AgentFakeTest::test_stream_accepts_ai_provider_enum
		// AgentFakeTest::test_timeout_can_be_passed_to_agent_prompt
		// AgentFakeTest::test_timeout_can_be_passed_to_agent_stream
		m := ai.NewManager()
		m.Extend(enums.LabOpenAI.String(), func(map[string]any) any { return newTextStub() })
		m.SetDefault(enums.LabOpenAI.String())
		rec := m.Fake("ok")
		rec.AssertAgentNeverQueued(t)
		queuedCheck := &recordingT{}
		rec.AssertAgentWasQueued(queuedCheck, func(*prompts.AgentPrompt) bool { return true })
		if !queuedCheck.failed {
			t.Fatal("expected missing queued assertion to be reported without panicking")
		}
		rec.AssertAgentNotQueued(t, func(*prompts.AgentPrompt) bool { return true })

		agent := ai.NewAnonymousAgent(m, "Route by provider.")
		if _, err := agent.Prompt(ctx, "prompt", contractsai.WithProvider(enums.LabOpenAI.String()), contractsai.WithTimeout(12)); err != nil {
			t.Fatalf("Prompt with provider error: %v", err)
		}
		if _, err := agent.Stream(ctx, "stream", contractsai.WithProvider(enums.LabOpenAI.String()), contractsai.WithTimeout(13)); err != nil {
			t.Fatalf("Stream with provider error: %v", err)
		}
		if _, err := agent.Queue(ctx, "queue", contractsai.WithProvider(enums.LabOpenAI.String())); err != nil {
			t.Fatalf("Queue with provider error: %v", err)
		}
		rec.AssertAgentWasPrompted(t, func(p *prompts.AgentPrompt) bool {
			return p.Text == "prompt" && p.Timeout == 12
		})
		rec.AssertAgentWasPrompted(t, func(p *prompts.AgentPrompt) bool {
			return p.Text == "stream" && p.Timeout == 13
		})
		rec.AssertAgentWasQueued(t, func(p *prompts.AgentPrompt) bool {
			return p.Text == "queue" && p.Provider != nil && *p.Provider == enums.LabOpenAI.String()
		})
	})
}

func TestInventoryProviderAgentFakeParity(t *testing.T) {
	t.Parallel()

	specs := []struct {
		lab     enums.Lab
		factory ai.DriverFactory
		markers []string
	}{
		{enums.LabAnthropic, func(cfg map[string]any) any { return anthropic.NewProvider(cfg) }, []string{
			"AgentFakeTest::test_anthropic_agent_can_be_faked",
			"AgentFakeTest::test_anthropic_agent_fake_with_closure",
			"AgentFakeTest::test_anthropic_agent_fake_with_no_predefined_responses",
			"AgentFakeTest::test_anthropic_agent_fake_records_prompts",
			"AgentFakeTest::test_anthropic_agent_stream_can_be_faked",
		}},
		{enums.LabAzure, func(cfg map[string]any) any { return azure.NewProvider(cfg) }, []string{
			"AgentFakeTest::test_azure_agent_can_be_faked",
			"AgentFakeTest::test_azure_agent_fake_with_closure",
			"AgentFakeTest::test_azure_agent_fake_with_no_predefined_responses",
			"AgentFakeTest::test_azure_agent_fake_records_prompts",
			"AgentFakeTest::test_azure_agent_stream_can_be_faked",
		}},
		{enums.LabGemini, func(cfg map[string]any) any { return gemini.NewProvider(cfg) }, []string{
			"AgentFakeTest::test_gemini_agent_can_be_faked",
			"AgentFakeTest::test_gemini_agent_fake_with_closure",
			"AgentFakeTest::test_gemini_agent_fake_with_no_predefined_responses",
			"AgentFakeTest::test_gemini_agent_fake_records_prompts",
			"AgentFakeTest::test_gemini_agent_stream_can_be_faked",
		}},
		{enums.LabGroq, func(cfg map[string]any) any { return groq.NewProvider(cfg) }, []string{
			"AgentFakeTest::test_groq_agent_can_be_faked",
			"AgentFakeTest::test_groq_agent_fake_with_closure",
			"AgentFakeTest::test_groq_agent_fake_with_no_predefined_responses",
			"AgentFakeTest::test_groq_agent_fake_records_prompts",
			"AgentFakeTest::test_groq_agent_stream_can_be_faked",
		}},
		{enums.LabMistral, func(cfg map[string]any) any { return mistral.NewProvider(cfg) }, []string{
			"AgentFakeTest::test_mistral_agent_can_be_faked",
			"AgentFakeTest::test_mistral_agent_fake_with_closure",
			"AgentFakeTest::test_mistral_agent_fake_with_no_predefined_responses",
			"AgentFakeTest::test_mistral_agent_fake_records_prompts",
			"AgentFakeTest::test_mistral_agent_stream_can_be_faked",
		}},
		{enums.LabOllama, func(cfg map[string]any) any { return ollama.NewProvider(cfg) }, []string{
			"AgentFakeTest::test_ollama_agent_can_be_faked",
			"AgentFakeTest::test_ollama_agent_fake_with_closure",
			"AgentFakeTest::test_ollama_agent_fake_with_no_predefined_responses",
			"AgentFakeTest::test_ollama_agent_fake_records_prompts",
			"AgentFakeTest::test_ollama_agent_stream_can_be_faked",
		}},
		{enums.LabOpenAI, func(cfg map[string]any) any { return openai.NewProvider(cfg) }, []string{
			"AgentFakeTest::test_openai_agent_can_be_faked",
			"AgentFakeTest::test_openai_agent_fake_with_closure",
			"AgentFakeTest::test_openai_agent_fake_with_no_predefined_responses",
			"AgentFakeTest::test_openai_agent_fake_records_prompts",
			"AgentFakeTest::test_openai_agent_stream_can_be_faked",
		}},
		{enums.LabOpenRouter, func(cfg map[string]any) any { return openrouter.NewProvider(cfg) }, []string{
			"AgentFakeTest::test_openrouter_agent_can_be_faked",
			"AgentFakeTest::test_openrouter_agent_fake_with_closure",
			"AgentFakeTest::test_openrouter_agent_fake_with_no_predefined_responses",
			"AgentFakeTest::test_openrouter_agent_fake_records_prompts",
			"AgentFakeTest::test_openrouter_agent_stream_can_be_faked",
		}},
		{enums.LabXAI, func(cfg map[string]any) any { return xai.NewProvider(cfg) }, []string{
			"AgentFakeTest::test_xai_agent_can_be_faked",
			"AgentFakeTest::test_xai_agent_fake_with_closure",
			"AgentFakeTest::test_xai_agent_fake_with_no_predefined_responses",
			"AgentFakeTest::test_xai_agent_fake_records_prompts",
			"AgentFakeTest::test_xai_agent_stream_can_be_faked",
		}},
	}

	for _, spec := range specs {
		spec := spec
		t.Run(spec.lab.String(), func(t *testing.T) {
			t.Parallel()

			m := ai.NewManager()
			m.Extend(spec.lab.String(), spec.factory)
			m.SetDefault(spec.lab.String())
			rec := m.Fake(func(p *prompts.AgentPrompt) (*responses.AgentResponse, error) {
				return responses.NewAgentResponse(p.InvocationID, "provider: "+p.Text, data.Usage{}, data.Meta{}), nil
			})

			agent := ai.NewAnonymousAgent(m, "Provider fake.")
			emptyProvider := ai.NewManager()
			emptyProvider.Extend(spec.lab.String(), spec.factory)
			emptyProvider.SetDefault(spec.lab.String())
			emptyProvider.Fake()
			emptyResp, err := ai.NewAnonymousAgent(emptyProvider, "Empty.").Prompt(context.Background(), "empty")
			if err != nil {
				t.Fatalf("empty provider prompt error: %v", err)
			}
			if emptyResp.GetText() != "" {
				t.Fatalf("empty provider text = %q", emptyResp.GetText())
			}

			resp, err := agent.Prompt(context.Background(), spec.lab.String())
			if err != nil {
				t.Fatalf("provider prompt error: %v", err)
			}
			if resp.GetText() != "provider: "+spec.lab.String() {
				t.Fatalf("provider text = %q", resp.GetText())
			}
			rec.AssertAgentWasPrompted(t, func(p *prompts.AgentPrompt) bool {
				return p.Text == spec.lab.String()
			})

			streamed, err := agent.Stream(context.Background(), "stream "+spec.lab.String())
			if err != nil {
				t.Fatalf("provider stream error: %v", err)
			}
			if streamed.(*responses.StreamableAgentResponse).Consume().GetText() == "" {
				t.Fatal("expected provider stream text")
			}
		})
	}
}

func TestInventoryAiManagerParity(t *testing.T) {
	t.Parallel()

	// AiManagerTest::test_can_get_an_openai_provider_instance
	// AiManagerTest::test_provider_type_is_ensured
	m := ai.NewManager()
	m.Extend(enums.LabOpenAI.String(), func(cfg map[string]any) any { return openai.NewProvider(cfg) })
	m.Extend(enums.LabElevenLabs.String(), func(map[string]any) any { return &audioOnlyStub{} })
	m.SetDefault(enums.LabOpenAI.String())

	if _, err := m.TextProvider(enums.LabOpenAI); err != nil {
		t.Fatalf("TextProvider(openai) error: %v", err)
	}

	m.SetDefault(enums.LabElevenLabs.String())
	if _, err := m.TextProvider(); err == nil {
		t.Fatal("expected provider type error")
	}
}

func TestInventoryAudioFakeParity(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// AudioFakeTest::test_audio_can_be_faked
	// AudioFakeTest::test_can_assert_no_audio_was_generated
	// AudioFakeTest::test_audio_can_be_faked_with_no_predefined_responses
	// AudioFakeTest::test_audio_can_be_faked_with_a_single_closure_that_is_invoked_for_every_generation
	// AudioFakeTest::test_audio_timeout_defaults_to_sdk_fallback
	// AudioFakeTest::test_fake_audio_closure_receives_timeout
	// AudioFakeTest::test_audio_can_prevent_stray_generations
	// AudioFakeTest::test_fake_closures_can_throw_exceptions
	// AudioFakeTest::test_audio_voice_and_instructions_are_recorded
	// AudioFakeTest::test_generate_accepts_ai_provider_enum
	m := ai.NewManager()
	rec := m.Fake()
	rec.AssertNothingAudioGenerated(t)
	provider, err := m.AudioProvider()
	if err != nil {
		t.Fatalf("AudioProvider error: %v", err)
	}
	instructions := "slowly"
	audio, err := provider.Audio(ctx, contractsprovider.AudioGenerateRequest{
		Text:         "hello",
		Voice:        "alloy",
		Instructions: &instructions,
		Timeout:      15,
	})
	if err != nil {
		t.Fatalf("Audio error: %v", err)
	}
	if audio.Content == "" {
		t.Fatal("expected fake audio content")
	}
	rec.AssertAudioGenerated(t, func(p *prompts.AudioPrompt) bool {
		return p.Text == "hello" && p.Voice == "alloy" && p.Instructions != nil && *p.Instructions == "slowly" && p.Timeout == 15
	})

	if _, err := provider.Audio(ctx, contractsprovider.AudioGenerateRequest{Text: "default timeout"}); err != nil {
		t.Fatalf("Audio default timeout error: %v", err)
	}
	rec.AssertAudioGenerated(t, func(p *prompts.AudioPrompt) bool {
		return p.Text == "default timeout" && p.Timeout == fake.DefaultMediaTimeout
	})

	m.FakeAudioProvider(func(p *prompts.AudioPrompt) (*contractsgw.AudioGenerateResult, error) {
		if p.Timeout != 21 {
			t.Fatalf("audio closure timeout = %d", p.Timeout)
		}
		return &contractsgw.AudioGenerateResult{Content: "YXVkaW8="}, nil
	})
	if _, err := provider.Audio(ctx, contractsprovider.AudioGenerateRequest{Text: "closure", Timeout: 21}); err != nil {
		t.Fatalf("Audio closure error: %v", err)
	}

	gw := fake.NewAudioGateway(m.Recorder())
	gw.PreventStray()
	provider.UseAudioGateway(gw)
	if _, err := provider.Audio(ctx, contractsprovider.AudioGenerateRequest{Text: "stray"}); err == nil {
		t.Fatal("expected stray audio error")
	}
	expected := errors.New("audio failed")
	m.FakeAudioProvider(func(*prompts.AudioPrompt) (*contractsgw.AudioGenerateResult, error) {
		return nil, expected
	})
	if _, err := provider.Audio(ctx, contractsprovider.AudioGenerateRequest{Text: "fail"}); !errors.Is(err, expected) {
		t.Fatalf("Audio error = %v, want %v", err, expected)
	}
}

func TestInventoryEmbeddingsFakeParity(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// EmbeddingsFakeTest::test_can_fake_embeddings
	// EmbeddingsFakeTest::test_can_fake_embeddings_with_custom_dimensions
	// EmbeddingsFakeTest::test_can_fake_embeddings_with_multiple_inputs
	// EmbeddingsFakeTest::test_can_fake_embeddings_with_custom_response
	// EmbeddingsFakeTest::test_can_fake_embeddings_with_closure
	// EmbeddingsFakeTest::test_embeddings_timeout_defaults_to_sdk_fallback
	// EmbeddingsFakeTest::test_fake_embeddings_closure_receives_timeout
	// EmbeddingsFakeTest::test_fake_embeddings_are_normalized
	// EmbeddingsFakeTest::test_can_prevent_stray_embeddings_generations
	// EmbeddingsFakeTest::test_can_assert_embeddings_generated
	// EmbeddingsFakeTest::test_can_assert_nothing_generated
	// EmbeddingsFakeTest::test_generate_accepts_ai_provider_enum
	m := ai.NewManager()
	rec := m.Fake()
	rec.AssertNothingEmbeddingsGenerated(t)
	provider, err := m.EmbeddingProvider()
	if err != nil {
		t.Fatalf("EmbeddingProvider error: %v", err)
	}
	dims := 4
	result, err := provider.Embeddings(ctx, contractsprovider.EmbeddingRequest{
		Inputs:     []string{"one", "two"},
		Dimensions: &dims,
		Timeout:    9,
	})
	if err != nil {
		t.Fatalf("Embeddings error: %v", err)
	}
	if len(result.Embeddings) != 2 || len(result.Embeddings[0]) != dims {
		t.Fatalf("embedding shape = %dx%d", len(result.Embeddings), len(result.Embeddings[0]))
	}
	rec.AssertEmbeddingsGenerated(t, func(p *prompts.EmbeddingsPrompt) bool {
		return len(p.Inputs) == 2 && p.Dimensions != nil && *p.Dimensions == dims && p.Timeout == 9
	})

	if _, err := provider.Embeddings(ctx, contractsprovider.EmbeddingRequest{Inputs: []string{"default timeout"}}); err != nil {
		t.Fatalf("Embeddings default timeout error: %v", err)
	}
	rec.AssertEmbeddingsGenerated(t, func(p *prompts.EmbeddingsPrompt) bool {
		return len(p.Inputs) == 1 && p.Inputs[0] == "default timeout" && p.Timeout == fake.DefaultMediaTimeout
	})

	custom := [][]float64{{0.25, 0.75}}
	m.FakeEmbeddingProvider(custom)
	got, err := provider.Embeddings(ctx, contractsprovider.EmbeddingRequest{Inputs: []string{"custom"}})
	if err != nil {
		t.Fatalf("custom embeddings error: %v", err)
	}
	if got.Embeddings[0][1] != 0.75 {
		t.Fatalf("custom embedding = %#v", got.Embeddings)
	}

	m.FakeEmbeddingProvider(func(p *prompts.EmbeddingsPrompt) (*contractsgw.EmbeddingGenerateResult, error) {
		if p.Timeout != 22 {
			t.Fatalf("embedding closure timeout = %d", p.Timeout)
		}
		return &contractsgw.EmbeddingGenerateResult{Embeddings: [][]float64{{1}}}, nil
	})
	if _, err := provider.Embeddings(ctx, contractsprovider.EmbeddingRequest{Inputs: []string{"closure"}, Timeout: 22}); err != nil {
		t.Fatalf("closure embeddings error: %v", err)
	}
	if len(ai.FakeEmbedding(3)) != 3 {
		t.Fatal("expected normalized fake embedding dimensions")
	}

	gw := fake.NewEmbeddingGateway(m.Recorder())
	gw.PreventStray()
	provider.UseEmbeddingGateway(gw)
	if _, err := provider.Embeddings(ctx, contractsprovider.EmbeddingRequest{Inputs: []string{"stray"}}); err == nil {
		t.Fatal("expected stray embeddings error")
	}
}

func TestInventoryImageFakeParity(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// ImageFakeTest::test_images_can_be_faked
	// ImageFakeTest::test_can_assert_no_images_were_generated
	// ImageFakeTest::test_images_can_be_faked_with_no_predefined_responses
	// ImageFakeTest::test_images_can_be_faked_with_a_single_closure_that_is_invoked_for_every_generation
	// ImageFakeTest::test_images_can_prevent_stray_generations
	// ImageFakeTest::test_fake_closures_can_throw_exceptions
	// ImageFakeTest::test_image_size_and_quality_are_recorded
	// ImageFakeTest::test_generate_accepts_ai_provider_enum
	m := ai.NewManager()
	rec := m.Fake()
	rec.AssertNothingImageGenerated(t)
	provider, err := m.ImageProvider()
	if err != nil {
		t.Fatalf("ImageProvider error: %v", err)
	}
	size := "1024x1024"
	quality := "high"
	result, err := provider.Image(ctx, contractsprovider.ImageGenerateRequest{
		Prompt:  "a precise diagram",
		Size:    &size,
		Quality: &quality,
	})
	if err != nil {
		t.Fatalf("Image error: %v", err)
	}
	if len(result.Images) != 1 {
		t.Fatalf("fake image count = %d", len(result.Images))
	}
	rec.AssertImageGenerated(t, func(p *prompts.ImagePrompt) bool {
		return p.Prompt == "a precise diagram" && p.Size != nil && *p.Size == size && p.Quality != nil && *p.Quality == quality
	})

	m.FakeImageProvider(func(p *prompts.ImagePrompt) (*contractsgw.ImageGenerateResult, error) {
		return &contractsgw.ImageGenerateResult{Images: []contractsgw.GeneratedImageData{{Image: "custom", MimeType: "image/png"}}}, nil
	})
	if _, err := provider.Image(ctx, contractsprovider.ImageGenerateRequest{Prompt: "closure"}); err != nil {
		t.Fatalf("Image closure error: %v", err)
	}
	gw := fake.NewImageGateway(m.Recorder())
	gw.PreventStray()
	provider.UseImageGateway(gw)
	if _, err := provider.Image(ctx, contractsprovider.ImageGenerateRequest{Prompt: "stray"}); err == nil {
		t.Fatal("expected stray image error")
	}
	expected := errors.New("image failed")
	m.FakeImageProvider(func(*prompts.ImagePrompt) (*contractsgw.ImageGenerateResult, error) {
		return nil, expected
	})
	if _, err := provider.Image(ctx, contractsprovider.ImageGenerateRequest{Prompt: "fail"}); !errors.Is(err, expected) {
		t.Fatalf("Image error = %v, want %v", err, expected)
	}
}

func TestInventoryRerankingFakeParity(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// RerankingFakeTest::test_can_fake_reranking
	// RerankingFakeTest::test_can_fake_reranking_with_limit
	// RerankingFakeTest::test_can_fake_reranking_with_custom_response
	// RerankingFakeTest::test_can_fake_reranking_with_closure
	// RerankingFakeTest::test_can_assert_reranked
	// RerankingFakeTest::test_can_assert_nothing_reranked
	// RerankingFakeTest::test_can_prevent_stray_rerankings
	// RerankingFakeTest::test_can_get_documents_in_reranked_order
	// RerankingFakeTest::test_rerank_accepts_ai_provider_enum
	// RerankingFakeTest::test_prompt_records_limit
	m := ai.NewManager()
	rec := m.Fake()
	rec.AssertNothingReranked(t)
	provider, err := m.RerankingProvider()
	if err != nil {
		t.Fatalf("RerankingProvider error: %v", err)
	}
	limit := 1
	result, err := provider.Rerank(ctx, contractsprovider.RerankingRequest{
		Documents: []string{"alpha", "beta"},
		Query:     "a",
		Limit:     &limit,
	})
	if err != nil {
		t.Fatalf("Rerank error: %v", err)
	}
	if result.Results[0].Document != "alpha" {
		t.Fatalf("first reranked document = %q", result.Results[0].Document)
	}
	rec.AssertReranked(t, func(p *prompts.RerankingPrompt) bool {
		return p.Query == "a" && p.Limit != nil && *p.Limit == 1
	})

	custom := &contractsgw.RerankResult{Results: []contractsgw.RankedDocumentData{{Index: 1, Document: "beta", Score: 0.99}}}
	m.FakeRerankingProvider(custom)
	got, err := provider.Rerank(ctx, contractsprovider.RerankingRequest{Documents: []string{"alpha", "beta"}, Query: "b"})
	if err != nil {
		t.Fatalf("custom rerank error: %v", err)
	}
	if got.Results[0].Document != "beta" {
		t.Fatalf("custom reranked document = %q", got.Results[0].Document)
	}

	m.FakeRerankingProvider(func(p *prompts.RerankingPrompt) (*contractsgw.RerankResult, error) {
		return &contractsgw.RerankResult{Results: []contractsgw.RankedDocumentData{{Document: p.Query, Score: 1}}}, nil
	})
	if _, err := provider.Rerank(ctx, contractsprovider.RerankingRequest{Documents: []string{"x"}, Query: "closure"}); err != nil {
		t.Fatalf("closure rerank error: %v", err)
	}
	gw := fake.NewRerankingGateway(m.Recorder())
	gw.PreventStray()
	provider.UseRerankingGateway(gw)
	if _, err := provider.Rerank(ctx, contractsprovider.RerankingRequest{Documents: []string{"x"}, Query: "stray"}); err == nil {
		t.Fatal("expected stray reranking error")
	}
}

func TestInventoryTranscriptionFakeParity(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// TranscriptionFakeTest::test_transcriptions_can_be_faked
	// TranscriptionFakeTest::test_can_assert_no_transcriptions_were_generated
	// TranscriptionFakeTest::test_transcriptions_can_be_faked_with_no_predefined_responses
	// TranscriptionFakeTest::test_transcriptions_can_be_faked_with_a_single_closure_that_is_invoked_for_every_generation
	// TranscriptionFakeTest::test_transcriptions_can_prevent_stray_generations
	// TranscriptionFakeTest::test_fake_closures_can_throw_exceptions
	// TranscriptionFakeTest::test_transcription_language_and_diarize_are_recorded
	// TranscriptionFakeTest::test_fake_transcriptions_include_segments
	// TranscriptionFakeTest::test_generate_accepts_ai_provider_enum
	// TranscriptionFakeTest::test_transcription_can_have_timeouts
	m := ai.NewManager()
	rec := m.Fake()
	rec.AssertNothingTranscriptionGenerated(t)
	provider, err := m.TranscriptionProvider()
	if err != nil {
		t.Fatalf("TranscriptionProvider error: %v", err)
	}
	language := "en"
	audio := contractsgw.TranscribableAudio{Content: "ZmFrZQ==", MimeType: "audio/mp3"}
	result, err := provider.Transcribe(ctx, contractsprovider.TranscriptionRequest{
		Audio:    audio,
		Language: &language,
		Diarize:  true,
		Timeout:  16,
	})
	if err != nil {
		t.Fatalf("Transcription error: %v", err)
	}
	if result.Text == "" {
		t.Fatal("expected fake transcription text")
	}
	rec.AssertTranscriptionGenerated(t, func(p *prompts.TranscriptionPrompt) bool {
		return p.Audio.Content == audio.Content && p.Language != nil && *p.Language == language && p.Diarize && p.Timeout == 16
	})

	custom := &contractsgw.TranscriptionResult{
		Text:     "segmented",
		Segments: []contractsgw.TranscriptionSegmentData{{Start: 0, End: 1, Text: "segmented"}},
	}
	m.FakeTranscriptionProvider(custom)
	got, err := provider.Transcribe(ctx, contractsprovider.TranscriptionRequest{Audio: audio})
	if err != nil {
		t.Fatalf("custom transcription error: %v", err)
	}
	if len(got.Segments) != 1 {
		t.Fatalf("segments = %#v", got.Segments)
	}
	m.FakeTranscriptionProvider(func(p *prompts.TranscriptionPrompt) (*contractsgw.TranscriptionResult, error) {
		return &contractsgw.TranscriptionResult{Text: p.Audio.Content}, nil
	})
	if _, err := provider.Transcribe(ctx, contractsprovider.TranscriptionRequest{Audio: audio}); err != nil {
		t.Fatalf("closure transcription error: %v", err)
	}
	gw := fake.NewTranscriptionGateway(m.Recorder())
	gw.PreventStray()
	provider.UseTranscriptionGateway(gw)
	if _, err := provider.Transcribe(ctx, contractsprovider.TranscriptionRequest{Audio: audio}); err == nil {
		t.Fatal("expected stray transcription error")
	}
	expected := errors.New("transcription failed")
	m.FakeTranscriptionProvider(func(*prompts.TranscriptionPrompt) (*contractsgw.TranscriptionResult, error) {
		return nil, expected
	})
	if _, err := provider.Transcribe(ctx, contractsprovider.TranscriptionRequest{Audio: audio}); !errors.Is(err, expected) {
		t.Fatalf("Transcription error = %v, want %v", err, expected)
	}
}

func newTextStub() contractsprovider.TextProvider {
	m := ai.NewManager()
	m.Fake()
	p, err := m.TextProvider()
	if err != nil {
		panic(err)
	}
	return p
}
