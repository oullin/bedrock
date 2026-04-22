package providers_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	anthropic "github.com/bedrock/packages/ai/sdk/providers/anthropic"
	azure "github.com/bedrock/packages/ai/sdk/providers/azure"
	cohere "github.com/bedrock/packages/ai/sdk/providers/cohere"
	deepseek "github.com/bedrock/packages/ai/sdk/providers/deepseek"
	gemini "github.com/bedrock/packages/ai/sdk/providers/gemini"
	groq "github.com/bedrock/packages/ai/sdk/providers/groq"
	jina "github.com/bedrock/packages/ai/sdk/providers/jina"
	mistral "github.com/bedrock/packages/ai/sdk/providers/mistral"
	ollama "github.com/bedrock/packages/ai/sdk/providers/ollama"
	openai "github.com/bedrock/packages/ai/sdk/providers/openai"
	openrouter "github.com/bedrock/packages/ai/sdk/providers/openrouter"
	voyageai "github.com/bedrock/packages/ai/sdk/providers/voyageai"
	xai "github.com/bedrock/packages/ai/sdk/providers/xai"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"
)

type textWrapper interface {
	Name() string
	DefaultTextModel() string
	CheapestTextModel() string
	SmartestTextModel() string
	Prompt(context.Context, contractsprovider.TextPromptRequest) (*contractsprovider.TextPromptResult, error)
	Stream(context.Context, contractsprovider.TextPromptRequest) (contractsprovider.StreamTextResult, error)
	UseTextGateway(contractsgw.TextGateway) contractsprovider.TextProvider
}

type transcriptionWrapper interface {
	DefaultTranscriptionModel() string
	Transcribe(context.Context, contractsprovider.TranscriptionRequest) (*contractsgw.TranscriptionResult, error)
	UseTranscriptionGateway(contractsgw.TranscriptionGateway) contractsprovider.TranscriptionProvider
}

type embeddingWrapper interface {
	DefaultEmbeddingsModel() string
	DefaultEmbeddingsDimensions() int
	Embeddings(context.Context, contractsprovider.EmbeddingRequest) (*contractsgw.EmbeddingGenerateResult, error)
	UseEmbeddingGateway(contractsgw.EmbeddingGateway) contractsprovider.EmbeddingProvider
}

type rerankingWrapper interface {
	DefaultRerankingModel() string
	Rerank(context.Context, contractsprovider.RerankingRequest) (*contractsgw.RerankResult, error)
	UseRerankingGateway(contractsgw.RerankingGateway) contractsprovider.RerankingProvider
}

type fileWrapper interface {
	GetFile(context.Context, string) (*contractsgw.FileGetResult, error)
	PutFile(context.Context, contractsgw.StorableFile) (*contractsgw.FilePutResult, error)
	DeleteFile(context.Context, string) error
	UseFileGateway(contractsgw.FileGateway) contractsprovider.FileProvider
}

func TestInventoryProviderWrapperTextDefaultsAndDelegation(t *testing.T) {
	t.Parallel()

	// RequestMappingTest::test_request_includes_model_and_messages
	// RequestMappingTest::test_response_text_is_correctly_parsed
	// RequestMappingTest::test_response_usage_is_correctly_parsed
	// StreamingTest::test_streaming_emits_text_events
	// RequestMappingTest::test_request_includes_model_and_input
	// RequestMappingTest::test_response_text_is_correctly_parsed
	// RequestMappingTest::test_response_usage_is_correctly_parsed
	// StreamingTest::test_streaming_emits_text_events
	// RequestMappingTest::test_request_includes_model_and_messages
	// RequestMappingTest::test_response_text_is_correctly_parsed
	// RequestMappingTest::test_response_usage_is_correctly_parsed
	// StreamingTest::test_streaming_emits_text_events
	// RequestMappingTest::test_request_includes_model_and_messages
	// RequestMappingTest::test_response_text_is_correctly_parsed
	// RequestMappingTest::test_response_usage_is_correctly_parsed
	// StreamingTest::test_streaming_emits_text_events
	// RequestMappingTest::test_request_includes_model_in_url_and_contents
	// RequestMappingTest::test_response_text_is_correctly_parsed
	// RequestMappingTest::test_response_usage_is_correctly_parsed
	// StreamingTest::test_streaming_emits_text_events
	// RequestMappingTest::test_request_includes_model_and_messages
	// RequestMappingTest::test_response_text_is_correctly_parsed
	// RequestMappingTest::test_response_usage_is_correctly_parsed
	// StreamingTest::test_streaming_emits_text_events
	// RequestMappingTest::test_request_includes_model_and_messages
	// RequestMappingTest::test_response_text_is_correctly_parsed
	// RequestMappingTest::test_response_usage_is_correctly_parsed
	// StreamingTest::test_streaming_emits_text_events
	// RequestMappingTest::test_request_includes_model_and_messages
	// RequestMappingTest::test_response_text_is_correctly_parsed
	// RequestMappingTest::test_response_usage_is_correctly_parsed
	// StreamingTest::test_streaming_emits_text_events
	// RequestMappingTest::test_request_includes_model_and_messages
	// RequestMappingTest::test_response_text_is_correctly_parsed
	// RequestMappingTest::test_response_usage_is_correctly_parsed
	// StreamingTest::test_streaming_emits_text_events
	// RequestMappingTest::test_request_includes_model_and_messages
	// RequestMappingTest::test_response_text_is_correctly_parsed
	// RequestMappingTest::test_response_usage_is_correctly_parsed
	// StreamingTest::test_streaming_emits_text_events
	// RequestMappingTest::test_request_includes_model_and_input
	// RequestMappingTest::test_response_text_is_correctly_parsed
	// RequestMappingTest::test_response_usage_is_correctly_parsed
	// StreamingTest::test_streaming_emits_text_events
	ctx := context.Background()

	specs := []struct {
		name     string
		factory  func(map[string]any) textWrapper
		text     string
		cheapest string
		smartest string
	}{
		{
			name:     "anthropic",
			factory:  func(cfg map[string]any) textWrapper { return anthropic.NewProvider(cfg) },
			text:     "claude-sonnet-4-5",
			cheapest: "claude-haiku-4-5",
			smartest: "claude-opus-4-5",
		},
		{
			name:     "azure",
			factory:  func(cfg map[string]any) textWrapper { return azure.NewProvider(cfg) },
			text:     "gpt-4o",
			cheapest: "gpt-4o-mini",
			smartest: "gpt-4",
		},
		{
			name:     "cohere",
			factory:  func(cfg map[string]any) textWrapper { return cohere.NewProvider(cfg) },
			text:     "command-r-plus",
			cheapest: "command-r",
			smartest: "command-r-plus",
		},
		{
			name:     "deepseek",
			factory:  func(cfg map[string]any) textWrapper { return deepseek.NewProvider(cfg) },
			text:     "deepseek-chat",
			cheapest: "deepseek-chat",
			smartest: "deepseek-reasoner",
		},
		{
			name:     "gemini",
			factory:  func(cfg map[string]any) textWrapper { return gemini.NewProvider(cfg) },
			text:     "gemini-2.0-flash",
			cheapest: "gemini-2.0-flash-lite",
			smartest: "gemini-2.5-pro",
		},
		{
			name:     "groq",
			factory:  func(cfg map[string]any) textWrapper { return groq.NewProvider(cfg) },
			text:     "llama-3.3-70b-versatile",
			cheapest: "llama-3.1-8b-instant",
			smartest: "llama-3.3-70b-versatile",
		},
		{
			name:     "mistral",
			factory:  func(cfg map[string]any) textWrapper { return mistral.NewProvider(cfg) },
			text:     "mistral-small-latest",
			cheapest: "open-mistral-nemo",
			smartest: "mistral-large-latest",
		},
		{
			name:     "ollama",
			factory:  func(cfg map[string]any) textWrapper { return ollama.NewProvider(cfg) },
			text:     "llama3.2",
			cheapest: "llama3.2:1b",
			smartest: "llama3.3:70b",
		},
		{
			name:     "openai",
			factory:  func(cfg map[string]any) textWrapper { return openai.NewProvider(cfg) },
			text:     "gpt-4o",
			cheapest: "gpt-4o-mini",
			smartest: "o3",
		},
		{
			name:     "openrouter",
			factory:  func(cfg map[string]any) textWrapper { return openrouter.NewProvider(cfg) },
			text:     "openai/gpt-4o",
			cheapest: "openai/gpt-4o-mini",
			smartest: "anthropic/claude-opus-4-5",
		},
		{
			name:     "xai",
			factory:  func(cfg map[string]any) textWrapper { return xai.NewProvider(cfg) },
			text:     "grok-3",
			cheapest: "grok-3-mini",
			smartest: "grok-3",
		},
	}

	for _, spec := range specs {
		spec := spec

		t.Run(spec.name, func(t *testing.T) {
			t.Parallel()

			provider := spec.factory(nil)
			if provider.Name() != spec.name {
				t.Fatalf("Name() = %q, want %q", provider.Name(), spec.name)
			}
			if provider.DefaultTextModel() != spec.text {
				t.Fatalf("DefaultTextModel() = %q, want %q", provider.DefaultTextModel(), spec.text)
			}
			if provider.CheapestTextModel() != spec.cheapest {
				t.Fatalf("CheapestTextModel() = %q, want %q", provider.CheapestTextModel(), spec.cheapest)
			}
			if provider.SmartestTextModel() != spec.smartest {
				t.Fatalf("SmartestTextModel() = %q, want %q", provider.SmartestTextModel(), spec.smartest)
			}

			textGW := &textGatewayStub{
				generateResult: &contractsgw.TextGenerateResult{
					Text: "reply",
					Usage: contractsgw.TokenUsage{
						PromptTokens:     3,
						CompletionTokens: 5,
					},
					Meta: contractsgw.ResponseMeta{Provider: stringPtr(spec.name), Model: stringPtr(spec.text)},
				},
				streamEvents: streamSeq(
					contractsgw.StreamEvent(streamEvent("stream.start", "stream", nil)),
					contractsgw.StreamEvent(streamEvent("text.delta", "hello", nil)),
					contractsgw.StreamEvent(streamEvent("stream.end", "stream", nil)),
				),
			}
			provider.UseTextGateway(textGW)

			model := provider.DefaultTextModel()
			instructions := "Be precise."
			got, err := provider.Prompt(ctx, contractsprovider.TextPromptRequest{
				InvocationID: "inv-1",
				Instructions: &instructions,
				Text:         "hello",
				Model:        &model,
				Messages:     []any{"user message"},
				Timeout:      21,
			})
			if err != nil {
				t.Fatalf("Prompt() returned error: %v", err)
			}
			if got.InvocationID != "inv-1" || got.Text != "reply" {
				t.Fatalf("Prompt() result = %#v", got)
			}
			if got.Usage.PromptTokens != 3 || got.Usage.CompletionTokens != 5 {
				t.Fatalf("Prompt() usage = %#v", got.Usage)
			}

			wantReq := contractsgw.TextGenerateRequest{
				Model:        spec.text,
				Instructions: &instructions,
				Text:         "hello",
				Messages:     []any{"user message"},
				Timeout:      21,
			}
			if !reflect.DeepEqual(textGW.generateReq, wantReq) {
				t.Fatalf("GenerateText request mismatch:\n got: %#v\nwant: %#v", textGW.generateReq, wantReq)
			}

			streamed, err := provider.Stream(ctx, contractsprovider.TextPromptRequest{
				InvocationID: "inv-stream",
				Text:         "stream me",
				Model:        &model,
				Timeout:      13,
			})
			if err != nil {
				t.Fatalf("Stream() returned error: %v", err)
			}
			if streamed.InvocationID != "inv-stream" {
				t.Fatalf("Stream() invocation ID = %q, want %q", streamed.InvocationID, "inv-stream")
			}

			var gotEvents []contractsgw.StreamEvent
			streamed.Events(func(event contractsgw.StreamEvent) bool {
				gotEvents = append(gotEvents, event)

				return true
			})
			wantEvents := []contractsgw.StreamEvent{
				contractsgw.StreamEvent(streamEvent("stream.start", "stream", nil)),
				contractsgw.StreamEvent(streamEvent("text.delta", "hello", nil)),
				contractsgw.StreamEvent(streamEvent("stream.end", "stream", nil)),
			}
			if !reflect.DeepEqual(gotEvents, wantEvents) {
				t.Fatalf("Stream() events mismatch:\n got: %#v\nwant: %#v", gotEvents, wantEvents)
			}

			wantErr := errors.New("text gateway failed")
			provider.UseTextGateway(&textGatewayStub{generateErr: wantErr, streamErr: wantErr})
			if _, err := provider.Prompt(ctx, contractsprovider.TextPromptRequest{}); !errors.Is(err, wantErr) {
				t.Fatalf("Prompt() error = %v, want %v", err, wantErr)
			}
			if _, err := provider.Stream(ctx, contractsprovider.TextPromptRequest{}); !errors.Is(err, wantErr) {
				t.Fatalf("Stream() error = %v, want %v", err, wantErr)
			}
		})
	}
}

func TestInventoryGroqTranscriptionWrapper(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	provider := groq.NewProvider(nil)

	if provider.DefaultTranscriptionModel() != "whisper-large-v3-turbo" {
		t.Fatalf("DefaultTranscriptionModel() = %q, want whisper-large-v3-turbo", provider.DefaultTranscriptionModel())
	}

	model := provider.DefaultTranscriptionModel()
	language := "en"
	audio := contractsgw.TranscribableAudio{Content: "YWJj", MimeType: "audio/mpeg"}
	gw := &transcriptionGatewayStub{
		result: &contractsgw.TranscriptionResult{
			Text:     "transcribed text",
			Segments: []contractsgw.TranscriptionSegmentData{{Start: 0, End: 1.2, Text: "transcribed text"}},
			Usage:    contractsgw.TokenUsage{PromptTokens: 5, CompletionTokens: 6},
			Meta:     contractsgw.ResponseMeta{Provider: stringPtr("groq"), Model: stringPtr(model)},
		},
	}
	provider.UseTranscriptionGateway(gw)

	got, err := provider.Transcribe(ctx, contractsprovider.TranscriptionRequest{
		Audio:    audio,
		Language: &language,
		Diarize:  true,
		Model:    &model,
		Timeout:  14,
	})
	if err != nil {
		t.Fatalf("Transcribe() returned error: %v", err)
	}
	if !reflect.DeepEqual(got, gw.result) {
		t.Fatalf("Transcribe() response mismatch:\n got: %#v\nwant: %#v", got, gw.result)
	}

	wantReq := contractsgw.TranscriptionRequest{
		Model:    model,
		Audio:    audio,
		Language: &language,
		Diarize:  true,
		Timeout:  14,
	}
	if !reflect.DeepEqual(gw.req, wantReq) {
		t.Fatalf("GenerateTranscription request mismatch:\n got: %#v\nwant: %#v", gw.req, wantReq)
	}

	wantErr := errors.New("transcription gateway failed")
	provider.UseTranscriptionGateway(&transcriptionGatewayStub{err: wantErr})
	if _, err := provider.Transcribe(ctx, contractsprovider.TranscriptionRequest{}); !errors.Is(err, wantErr) {
		t.Fatalf("Transcribe() error = %v, want %v", err, wantErr)
	}
}

func TestInventoryProviderEmbeddingWrappers(t *testing.T) {
	t.Parallel()

	// EmbeddingTest::test_embeddings_request_includes_model_input_and_dimensions
	// EmbeddingTest::test_embeddings_response_is_correctly_parsed
	// EmbeddingTest::test_multiple_inputs_return_multiple_embeddings
	// EmbeddingTest::test_embeddings_request_includes_model_texts_input_type_and_embedding_types
	// EmbeddingTest::test_embeddings_response_is_correctly_parsed
	// EmbeddingTest::test_multiple_inputs_return_multiple_embeddings
	// EmbeddingTest::test_embeddings_request_includes_model_input_and_dimensions
	// EmbeddingTest::test_embeddings_response_is_correctly_parsed
	// EmbeddingTest::test_multiple_inputs_return_multiple_embeddings
	// EmbeddingTest::test_embeddings_request_includes_model_input_and_dimensions
	// EmbeddingTest::test_embeddings_response_is_correctly_parsed
	// EmbeddingTest::test_multiple_inputs_return_multiple_embeddings
	// EmbeddingTest::test_embeddings_request_includes_model_and_input
	// EmbeddingTest::test_embeddings_response_is_correctly_parsed
	// EmbeddingTest::test_multiple_inputs_return_multiple_embeddings
	// EmbeddingTest::test_embeddings_request_includes_model_and_input
	// EmbeddingTest::test_embeddings_response_is_correctly_parsed
	// EmbeddingTest::test_multiple_inputs_return_multiple_embeddings
	// EmbeddingTest::test_embeddings_request_includes_model_input_and_dimensions
	// EmbeddingTest::test_embeddings_response_is_correctly_parsed
	// EmbeddingTest::test_multiple_inputs_return_multiple_embeddings
	// EmbeddingTest::test_embeddings_request_includes_model_input_and_output_dimension
	// EmbeddingTest::test_embeddings_response_is_correctly_parsed
	// EmbeddingTest::test_multiple_inputs_return_multiple_embeddings
	ctx := context.Background()

	embeddingCases := []struct {
		name        string
		provider    embeddingWrapper
		model       string
		dimensions  int
		providerTag string
	}{
		{
			name:        "azure",
			provider:    azure.NewProvider(nil),
			model:       "text-embedding-3-small",
			dimensions:  1536,
			providerTag: "azure",
		},
		{
			name:        "cohere",
			provider:    cohere.NewProvider(nil),
			model:       "embed-v4.0",
			dimensions:  1024,
			providerTag: "cohere",
		},
		{
			name:        "gemini",
			provider:    gemini.NewProvider(nil),
			model:       "text-embedding-004",
			dimensions:  768,
			providerTag: "gemini",
		},
		{
			name:        "jina",
			provider:    jina.NewProvider(nil),
			model:       "jina-embeddings-v3",
			dimensions:  1024,
			providerTag: "jina",
		},
		{
			name:        "mistral",
			provider:    mistral.NewProvider(nil),
			model:       "mistral-embed",
			dimensions:  1024,
			providerTag: "mistral",
		},
		{
			name:        "ollama",
			provider:    ollama.NewProvider(nil),
			model:       "nomic-embed-text",
			dimensions:  768,
			providerTag: "ollama",
		},
		{
			name:        "openai",
			provider:    openai.NewProvider(nil),
			model:       "text-embedding-3-small",
			dimensions:  1536,
			providerTag: "openai",
		},
		{
			name:        "voyageai",
			provider:    voyageai.NewProvider(nil),
			model:       "voyage-3-large",
			dimensions:  1024,
			providerTag: "voyageai",
		},
		{
			name:        "xai",
			provider:    xai.NewProvider(nil),
			model:       "v1",
			dimensions:  1536,
			providerTag: "xai",
		},
	}

	for _, tc := range embeddingCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.provider.DefaultEmbeddingsModel() != tc.model {
				t.Fatalf("DefaultEmbeddingsModel() = %q, want %q", tc.provider.DefaultEmbeddingsModel(), tc.model)
			}
			if tc.provider.DefaultEmbeddingsDimensions() != tc.dimensions {
				t.Fatalf("DefaultEmbeddingsDimensions() = %d, want %d", tc.provider.DefaultEmbeddingsDimensions(), tc.dimensions)
			}

			gw := &embeddingGatewayStub{
				result: &contractsgw.EmbeddingGenerateResult{
					Embeddings: [][]float64{{0.25, 0.75}, {0.5, 0.5}},
					Tokens:     2,
					Meta:       contractsgw.ResponseMeta{Provider: stringPtr(tc.providerTag), Model: stringPtr(tc.model)},
				},
			}
			tc.provider.UseEmbeddingGateway(gw)

			got, err := tc.provider.Embeddings(ctx, contractsprovider.EmbeddingRequest{
				Inputs:     []string{"one", "two"},
				Dimensions: &tc.dimensions,
				Model:      &tc.model,
				Timeout:    12,
			})
			if err != nil {
				t.Fatalf("Embeddings() returned error: %v", err)
			}
			if !reflect.DeepEqual(got, gw.result) {
				t.Fatalf("Embeddings() response mismatch:\n got: %#v\nwant: %#v", got, gw.result)
			}
			if len(got.Embeddings) != 2 {
				t.Fatalf("Embeddings() count = %d, want 2", len(got.Embeddings))
			}

			wantReq := contractsgw.EmbeddingGenerateRequest{
				Model:      tc.model,
				Inputs:     []string{"one", "two"},
				Dimensions: tc.dimensions,
				Timeout:    12,
			}
			if !reflect.DeepEqual(gw.req, wantReq) {
				t.Fatalf("GenerateEmbeddings request mismatch:\n got: %#v\nwant: %#v", gw.req, wantReq)
			}

			wantErr := errors.New("embedding gateway failed")
			tc.provider.UseEmbeddingGateway(&embeddingGatewayStub{err: wantErr})
			if _, err := tc.provider.Embeddings(ctx, contractsprovider.EmbeddingRequest{}); !errors.Is(err, wantErr) {
				t.Fatalf("Embeddings() error = %v, want %v", err, wantErr)
			}
		})
	}

	t.Run("mistral file gateway", func(t *testing.T) {
		t.Parallel()

		// FileGatewayTest::test_get_file_sends_correct_request
		// FileGatewayTest::test_put_file_sends_multipart_upload
		// FileGatewayTest::test_delete_file_sends_correct_request
		provider := gemini.NewProvider(nil)
		fileGW := &fileGatewayStub{
			getResult: &contractsgw.FileGetResult{ID: "file-1", Filename: "notes.txt", Size: 12, Usage: contractsgw.TokenUsage{PromptTokens: 1}, Meta: contractsgw.ResponseMeta{Provider: stringPtr("gemini")}},
			putResult: &contractsgw.FilePutResult{ID: "file-2", Filename: "notes.txt", Usage: contractsgw.TokenUsage{PromptTokens: 2}, Meta: contractsgw.ResponseMeta{Provider: stringPtr("gemini")}},
		}
		provider.UseFileGateway(fileGW)

		if got, err := provider.GetFile(ctx, "file-1"); err != nil {
			t.Fatalf("GetFile() returned error: %v", err)
		} else if !reflect.DeepEqual(got, fileGW.getResult) {
			t.Fatalf("GetFile() response mismatch:\n got: %#v\nwant: %#v", got, fileGW.getResult)
		}
		if fileGW.getID != "file-1" {
			t.Fatalf("GetFile() ID = %q, want %q", fileGW.getID, "file-1")
		}

		put := contractsgw.StorableFile{Content: []byte("abc"), Filename: "notes.txt", MimeType: "text/plain"}
		if got, err := provider.PutFile(ctx, put); err != nil {
			t.Fatalf("PutFile() returned error: %v", err)
		} else if !reflect.DeepEqual(got, fileGW.putResult) {
			t.Fatalf("PutFile() response mismatch:\n got: %#v\nwant: %#v", got, fileGW.putResult)
		}
		if !reflect.DeepEqual(fileGW.putFile, put) {
			t.Fatalf("PutFile() request mismatch:\n got: %#v\nwant: %#v", fileGW.putFile, put)
		}

		wantErr := errors.New("file gateway failed")
		provider.UseFileGateway(&fileGatewayStub{err: wantErr})
		if err := provider.DeleteFile(ctx, "file-9"); !errors.Is(err, wantErr) {
			t.Fatalf("DeleteFile() error = %v, want %v", err, wantErr)
		}
	})
}

func TestInventoryProviderRerankingWrappers(t *testing.T) {
	t.Parallel()

	// RerankingTest::test_reranking_request_includes_model_query_and_documents
	// RerankingTest::test_reranking_request_includes_top_n_when_limit_set
	// RerankingTest::test_reranking_response_is_correctly_parsed_into_rankeddocuments
	// RerankingTest::test_reranking_request_includes_model_query_and_documents
	// RerankingTest::test_reranking_request_includes_top_k_when_limit_set
	// RerankingTest::test_reranking_response_is_correctly_parsed_into_rankeddocuments
	ctx := context.Background()

	cases := []struct {
		name        string
		provider    rerankingWrapper
		model       string
		providerTag string
	}{
		{
			name:        "cohere",
			provider:    cohere.NewProvider(nil),
			model:       "rerank-v3.5",
			providerTag: "cohere",
		},
		{
			name:        "jina",
			provider:    jina.NewProvider(nil),
			model:       "jina-reranker-v2-base-multilingual",
			providerTag: "jina",
		},
		{
			name:        "voyageai",
			provider:    voyageai.NewProvider(nil),
			model:       "rerank-2",
			providerTag: "voyageai",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.provider.DefaultRerankingModel() != tc.model {
				t.Fatalf("DefaultRerankingModel() = %q, want %q", tc.provider.DefaultRerankingModel(), tc.model)
			}

			limit := 2
			gw := &rerankingGatewayStub{
				result: &contractsgw.RerankResult{
					Results: []contractsgw.RankedDocumentData{
						{Index: 1, Document: "second", Score: 0.9},
						{Index: 0, Document: "first", Score: 0.8},
					},
					Usage: contractsgw.TokenUsage{PromptTokens: 9, CompletionTokens: 10},
					Meta:  contractsgw.ResponseMeta{Provider: stringPtr(tc.providerTag), Model: stringPtr(tc.model)},
				},
			}
			tc.provider.UseRerankingGateway(gw)

			got, err := tc.provider.Rerank(ctx, contractsprovider.RerankingRequest{
				Documents: []string{"first", "second"},
				Query:     "match",
				Limit:     &limit,
				Model:     &tc.model,
			})
			if err != nil {
				t.Fatalf("Rerank() returned error: %v", err)
			}

			wantReq := contractsgw.RerankRequest{
				Model:     tc.model,
				Documents: []string{"first", "second"},
				Query:     "match",
				Limit:     &limit,
			}
			if !reflect.DeepEqual(gw.req, wantReq) {
				t.Fatalf("Rerank request mismatch:\n got: %#v\nwant: %#v", gw.req, wantReq)
			}
			if !reflect.DeepEqual(got, gw.result) {
				t.Fatalf("Rerank() response mismatch:\n got: %#v\nwant: %#v", got, gw.result)
			}

			wantErr := errors.New("reranking gateway failed")
			tc.provider.UseRerankingGateway(&rerankingGatewayStub{err: wantErr})
			if _, err := tc.provider.Rerank(ctx, contractsprovider.RerankingRequest{}); !errors.Is(err, wantErr) {
				t.Fatalf("Rerank() error = %v, want %v", err, wantErr)
			}
		})
	}
}
