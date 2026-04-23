package providers_test

import (
	"context"
	"errors"
	"iter"
	"reflect"
	"testing"

	"github.com/bedrock/packages/ai/sdk/providers/azure"
	"github.com/bedrock/packages/ai/sdk/providers/cohere"
	"github.com/bedrock/packages/ai/sdk/providers/gemini"
	"github.com/bedrock/packages/ai/sdk/providers/openai"
	contractsgw "github.com/bedrock/packages/contracts/ai/gateway"
	contractsprovider "github.com/bedrock/packages/contracts/ai/provider"
)

type textGatewayStub struct {
	generateReq contractsgw.TextGenerateRequest
	streamReq   contractsgw.TextGenerateRequest
	streamID    string

	generateResult *contractsgw.TextGenerateResult
	generateErr    error
	streamEvents   iter.Seq[contractsgw.StreamEvent]
	streamErr      error
}

type imageGatewayStub struct {
	req    contractsgw.ImageGenerateRequest
	result *contractsgw.ImageGenerateResult
	err    error
}

type audioGatewayStub struct {
	req    contractsgw.AudioGenerateRequest
	result *contractsgw.AudioGenerateResult
	err    error
}

type embeddingGatewayStub struct {
	req    contractsgw.EmbeddingGenerateRequest
	result *contractsgw.EmbeddingGenerateResult
	err    error
}

type transcriptionGatewayStub struct {
	req    contractsgw.TranscriptionRequest
	result *contractsgw.TranscriptionResult
	err    error
}

type rerankingGatewayStub struct {
	req    contractsgw.RerankRequest
	result *contractsgw.RerankResult
	err    error
}

type fileGatewayStub struct {
	getID    string
	putFile  contractsgw.StorableFile
	deleteID string

	getResult *contractsgw.FileGetResult
	putResult *contractsgw.FilePutResult
	err       error
}

type storeGatewayStub struct {
	getID      string
	createReq  contractsgw.StoreCreateRequest
	storeID    string
	fileID     string
	file       contractsgw.StorableFile
	metadata   map[string]any
	deleteID   string
	getResult  *contractsgw.StoreData
	createData *contractsgw.StoreData
	addResult  *contractsgw.AddedDocumentData
	err        error
}

// RequestMappingTest::test_request_includes_model_and_input
// RequestMappingTest::test_system_instructions_are_sent_as_system_message
// RequestMappingTest::test_response_text_is_correctly_parsed
// RequestMappingTest::test_response_usage_is_correctly_parsed
// StreamingTest::test_streaming_emits_text_events
// StreamingTest::test_streaming_handles_tool_calls

// EmbeddingTest::test_embeddings_request_includes_model_input_and_dimensions
// EmbeddingTest::test_embeddings_response_is_correctly_parsed
// EmbeddingTest::test_multiple_inputs_return_multiple_embeddings
// RerankingTest::test_reranking_request_includes_model_query_and_documents
// RerankingTest::test_reranking_request_includes_top_n_when_limit_set
// RerankingTest::test_reranking_response_is_correctly_parsed_into_rankeddocuments

// FileProviderTest::test_get_put_delete_forwards_request_response_and_error
// StoreProviderTest::test_get_create_add_remove_delete_forwards_request_response_and_error

type streamEventStub struct {
	kind string
	id   string
	text string
}

func (g *textGatewayStub) GenerateText(_ context.Context, req contractsgw.TextGenerateRequest) (*contractsgw.TextGenerateResult, error) {
	g.generateReq = req

	return g.generateResult, g.generateErr
}

func (g *textGatewayStub) StreamText(_ context.Context, invocationID string, req contractsgw.TextGenerateRequest) (iter.Seq[contractsgw.StreamEvent], error) {
	g.streamID = invocationID
	g.streamReq = req

	return g.streamEvents, g.streamErr
}

func (g *textGatewayStub) OnToolInvocation(func(context.Context, string, string, map[string]any) (any, error)) {
}

func (g *imageGatewayStub) GenerateImage(_ context.Context, req contractsgw.ImageGenerateRequest) (*contractsgw.ImageGenerateResult, error) {
	g.req = req

	return g.result, g.err
}

func (g *audioGatewayStub) GenerateAudio(_ context.Context, req contractsgw.AudioGenerateRequest) (*contractsgw.AudioGenerateResult, error) {
	g.req = req

	return g.result, g.err
}

func (g *embeddingGatewayStub) GenerateEmbeddings(_ context.Context, req contractsgw.EmbeddingGenerateRequest) (*contractsgw.EmbeddingGenerateResult, error) {
	g.req = req

	return g.result, g.err
}

func (g *transcriptionGatewayStub) GenerateTranscription(_ context.Context, req contractsgw.TranscriptionRequest) (*contractsgw.TranscriptionResult, error) {
	g.req = req

	return g.result, g.err
}

func (g *rerankingGatewayStub) Rerank(_ context.Context, req contractsgw.RerankRequest) (*contractsgw.RerankResult, error) {
	g.req = req

	return g.result, g.err
}

func (g *fileGatewayStub) GetFile(_ context.Context, id string) (*contractsgw.FileGetResult, error) {
	g.getID = id

	return g.getResult, g.err
}

func (g *fileGatewayStub) PutFile(_ context.Context, file contractsgw.StorableFile) (*contractsgw.FilePutResult, error) {
	g.putFile = file

	return g.putResult, g.err
}

func (g *fileGatewayStub) DeleteFile(_ context.Context, id string) error {
	g.deleteID = id

	return g.err
}

func (g *storeGatewayStub) GetStore(_ context.Context, id string) (*contractsgw.StoreData, error) {
	g.getID = id

	return g.getResult, g.err
}

func (g *storeGatewayStub) CreateStore(_ context.Context, req contractsgw.StoreCreateRequest) (*contractsgw.StoreData, error) {
	g.createReq = req

	return g.createData, g.err
}

func (g *storeGatewayStub) AddFile(_ context.Context, storeID string, file contractsgw.StorableFile, metadata map[string]any) (*contractsgw.AddedDocumentData, error) {
	g.storeID = storeID
	g.file = file
	g.metadata = metadata

	return g.addResult, g.err
}

func (g *storeGatewayStub) RemoveFile(_ context.Context, storeID, fileID string) error {
	g.storeID = storeID
	g.fileID = fileID

	return g.err
}

func (g *storeGatewayStub) DeleteStore(_ context.Context, id string) error {
	g.deleteID = id

	return g.err
}

func TestInventoryProviderTextMapping(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	model := "gpt-4o"
	instructions := "Be precise."
	provider := openai.NewProvider(nil)

	textGW := &textGatewayStub{
		generateResult: &contractsgw.TextGenerateResult{
			Text: "reply",
			Usage: contractsgw.TokenUsage{
				PromptTokens:          7,
				CompletionTokens:      11,
				CacheWriteInputTokens: 13,
				CacheReadInputTokens:  17,
				ReasoningTokens:       19,
			},
			Meta: contractsgw.ResponseMeta{
				Provider:  stringPtr("openai"),
				Model:     stringPtr("gpt-4o"),
				Citations: []any{"citation"},
			},
			Messages: []any{"assistant message"},
			Steps:    []any{"step-one"},
		},
		streamEvents: streamSeq(
			contractsgw.StreamEvent(streamEvent("stream.start", "stream", nil)),
			contractsgw.StreamEvent(streamEvent("text.start", "", nil)),
			contractsgw.StreamEvent(streamEvent("text.delta", "hello ", nil)),
			contractsgw.StreamEvent(streamEvent("text.delta", "world", nil)),
			contractsgw.StreamEvent(streamEvent("tool_call", "tool-1", nil)),
			contractsgw.StreamEvent(streamEvent("tool_result", "tool-1", nil)),
			contractsgw.StreamEvent(streamEvent("text.end", "hello world", nil)),
			contractsgw.StreamEvent(streamEvent("stream.end", "stream", nil)),
		),
	}
	provider.UseTextGateway(textGW)

	got, err := provider.Prompt(ctx, contractsprovider.TextPromptRequest{
		InvocationID: "inv-1",
		Instructions: &instructions,
		Text:         "hello",
		Messages:     []any{"user message"},
		Tools:        []any{"tool"},
		Schema:       map[string]any{"type": "object"},
		Options:      map[string]any{"temperature": 0.2},
		Model:        &model,
		Timeout:      21,
	})

	if err != nil {
		t.Fatalf("Prompt() returned error: %v", err)
	}

	wantReq := contractsgw.TextGenerateRequest{
		Model:        model,
		Instructions: &instructions,
		Text:         "hello",
		Messages:     []any{"user message"},
		Tools:        []any{"tool"},
		Schema:       map[string]any{"type": "object"},
		Options:      map[string]any{"temperature": 0.2},
		Timeout:      21,
	}

	if !reflect.DeepEqual(textGW.generateReq, wantReq) {
		t.Fatalf("GenerateText request mismatch:\n got: %#v\nwant: %#v", textGW.generateReq, wantReq)
	}

	wantResp := &contractsprovider.TextPromptResult{
		InvocationID: "inv-1",
		Text:         "reply",
		Usage: contractsgw.TokenUsage{
			PromptTokens:          7,
			CompletionTokens:      11,
			CacheWriteInputTokens: 13,
			CacheReadInputTokens:  17,
			ReasoningTokens:       19,
		},
		Meta: contractsgw.ResponseMeta{
			Provider:  stringPtr("openai"),
			Model:     stringPtr("gpt-4o"),
			Citations: []any{"citation"},
		},
		Messages: []any{"assistant message"},
		Steps:    []any{"step-one"},
	}

	if !reflect.DeepEqual(got, wantResp) {
		t.Fatalf("Prompt() response mismatch:\n got: %#v\nwant: %#v", got, wantResp)
	}

	streamed, err := provider.Stream(ctx, contractsprovider.TextPromptRequest{
		InvocationID: "inv-stream",
		Instructions: &instructions,
		Text:         "hello stream",
		Messages:     []any{"user stream"},
		Tools:        []any{"tool"},
		Schema:       map[string]any{"type": "object"},
		Options:      map[string]any{"temperature": 0.5},
		Model:        &model,
		Timeout:      13,
	})

	if err != nil {
		t.Fatalf("Stream() returned error: %v", err)
	}

	wantStreamReq := contractsgw.TextGenerateRequest{
		Model:        model,
		Instructions: &instructions,
		Text:         "hello stream",
		Messages:     []any{"user stream"},
		Tools:        []any{"tool"},
		Schema:       map[string]any{"type": "object"},
		Options:      map[string]any{"temperature": 0.5},
		Timeout:      13,
	}

	if textGW.streamID != "inv-stream" {
		t.Fatalf("StreamText invocation ID = %q, want %q", textGW.streamID, "inv-stream")
	}

	if !reflect.DeepEqual(textGW.streamReq, wantStreamReq) {
		t.Fatalf("StreamText request mismatch:\n got: %#v\nwant: %#v", textGW.streamReq, wantStreamReq)
	}

	if streamed.InvocationID != "inv-stream" {
		t.Fatalf("StreamText result invocation ID = %q, want %q", streamed.InvocationID, "inv-stream")
	}

	var gotEvents []contractsgw.StreamEvent

	streamed.Events(func(event contractsgw.StreamEvent) bool {
		gotEvents = append(gotEvents, event)

		return true
	})

	wantEvents := []contractsgw.StreamEvent{
		streamEvent("stream.start", "stream", nil),
		streamEvent("text.start", "", nil),
		streamEvent("text.delta", "hello ", nil),
		streamEvent("text.delta", "world", nil),
		streamEvent("tool_call", "tool-1", nil),
		streamEvent("tool_result", "tool-1", nil),
		streamEvent("text.end", "hello world", nil),
		streamEvent("stream.end", "stream", nil),
	}

	if !reflect.DeepEqual(gotEvents, wantEvents) {
		t.Fatalf("StreamText events mismatch:\n got: %#v\nwant: %#v", gotEvents, wantEvents)
	}

	wantErr := errors.New("text gateway failed")
	provider.UseTextGateway(&textGatewayStub{generateErr: wantErr, streamErr: wantErr})

	if _, err := provider.Prompt(ctx, contractsprovider.TextPromptRequest{}); !errors.Is(err, wantErr) {
		t.Fatalf("Prompt() error = %v, want %v", err, wantErr)
	}

	if _, err := provider.Stream(ctx, contractsprovider.TextPromptRequest{}); !errors.Is(err, wantErr) {
		t.Fatalf("Stream() error = %v, want %v", err, wantErr)
	}
}

func TestInventoryProviderMediaAndIndexMapping(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("image", func(t *testing.T) {
		t.Parallel()

		provider := azure.NewProvider(nil)
		size := "1024x1024"
		quality := "high"
		model := "dall-e-3"
		gw := &imageGatewayStub{
			result: &contractsgw.ImageGenerateResult{
				Images: []contractsgw.GeneratedImageData{{Image: "img", MimeType: "image/png"}},
				Usage:  contractsgw.TokenUsage{PromptTokens: 1, CompletionTokens: 2},
				Meta:   contractsgw.ResponseMeta{Provider: stringPtr("azure"), Model: stringPtr("dall-e-3")},
			},
		}
		provider.UseImageGateway(gw)

		got, err := provider.Image(ctx, contractsprovider.ImageGenerateRequest{
			Prompt:      "a diagram",
			Attachments: []any{"attachment"},
			Size:        &size,
			Quality:     &quality,
			Model:       &model,
			Timeout:     9,
		})

		if err != nil {
			t.Fatalf("Image() returned error: %v", err)
		}

		wantReq := contractsgw.ImageGenerateRequest{
			Model:       model,
			Prompt:      "a diagram",
			Attachments: []any{"attachment"},
			Size:        &size,
			Quality:     &quality,
			Timeout:     9,
		}

		if !reflect.DeepEqual(gw.req, wantReq) {
			t.Fatalf("GenerateImage request mismatch:\n got: %#v\nwant: %#v", gw.req, wantReq)
		}

		if !reflect.DeepEqual(got, gw.result) {
			t.Fatalf("Image() response mismatch:\n got: %#v\nwant: %#v", got, gw.result)
		}

		wantErr := errors.New("image gateway failed")
		provider.UseImageGateway(&imageGatewayStub{err: wantErr})

		if _, err := provider.Image(ctx, contractsprovider.ImageGenerateRequest{}); !errors.Is(err, wantErr) {
			t.Fatalf("Image() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("audio", func(t *testing.T) {
		t.Parallel()

		provider := openai.NewProvider(nil)
		instructions := "Whisper it."
		model := "tts-1"
		gw := &audioGatewayStub{
			result: &contractsgw.AudioGenerateResult{
				Content: "YXVkaW8=",
				Usage:   contractsgw.TokenUsage{PromptTokens: 3, CompletionTokens: 4},
				Meta:    contractsgw.ResponseMeta{Provider: stringPtr("openai"), Model: stringPtr("tts-1")},
			},
		}
		provider.UseAudioGateway(gw)

		got, err := provider.Audio(ctx, contractsprovider.AudioGenerateRequest{
			Text:         "hello",
			Voice:        "alloy",
			Instructions: &instructions,
			Model:        &model,
			Timeout:      8,
		})

		if err != nil {
			t.Fatalf("Audio() returned error: %v", err)
		}

		wantReq := contractsgw.AudioGenerateRequest{
			Model:        model,
			Text:         "hello",
			Voice:        "alloy",
			Instructions: &instructions,
			Timeout:      8,
		}

		if !reflect.DeepEqual(gw.req, wantReq) {
			t.Fatalf("GenerateAudio request mismatch:\n got: %#v\nwant: %#v", gw.req, wantReq)
		}

		if !reflect.DeepEqual(got, gw.result) {
			t.Fatalf("Audio() response mismatch:\n got: %#v\nwant: %#v", got, gw.result)
		}

		wantErr := errors.New("audio gateway failed")
		provider.UseAudioGateway(&audioGatewayStub{err: wantErr})

		if _, err := provider.Audio(ctx, contractsprovider.AudioGenerateRequest{}); !errors.Is(err, wantErr) {
			t.Fatalf("Audio() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("embeddings", func(t *testing.T) {
		t.Parallel()

		provider := azure.NewProvider(nil)
		model := "text-embedding-3-small"
		dimensions := 1536
		gw := &embeddingGatewayStub{
			result: &contractsgw.EmbeddingGenerateResult{
				Embeddings: [][]float64{{0.25, 0.75}},
				Tokens:     2,
				Meta:       contractsgw.ResponseMeta{Provider: stringPtr("azure"), Model: stringPtr(model)},
			},
		}
		provider.UseEmbeddingGateway(gw)

		got, err := provider.Embeddings(ctx, contractsprovider.EmbeddingRequest{
			Inputs:     []string{"one", "two"},
			Dimensions: &dimensions,
			Model:      &model,
			Timeout:    12,
		})

		if err != nil {
			t.Fatalf("Embeddings() returned error: %v", err)
		}

		wantReq := contractsgw.EmbeddingGenerateRequest{
			Model:      model,
			Inputs:     []string{"one", "two"},
			Dimensions: dimensions,
			Timeout:    12,
		}

		if !reflect.DeepEqual(gw.req, wantReq) {
			t.Fatalf("GenerateEmbeddings request mismatch:\n got: %#v\nwant: %#v", gw.req, wantReq)
		}

		if !reflect.DeepEqual(got, gw.result) {
			t.Fatalf("Embeddings() response mismatch:\n got: %#v\nwant: %#v", got, gw.result)
		}

		wantErr := errors.New("embeddings gateway failed")
		provider.UseEmbeddingGateway(&embeddingGatewayStub{err: wantErr})

		if _, err := provider.Embeddings(ctx, contractsprovider.EmbeddingRequest{}); !errors.Is(err, wantErr) {
			t.Fatalf("Embeddings() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("transcription", func(t *testing.T) {
		t.Parallel()

		provider := gemini.NewProvider(nil)
		language := "en"
		model := "gemini-2.0-flash"
		audio := contractsgw.TranscribableAudio{Content: "YWJj", MimeType: "audio/mpeg"}
		gw := &transcriptionGatewayStub{
			result: &contractsgw.TranscriptionResult{
				Text:     "transcribed text",
				Segments: []contractsgw.TranscriptionSegmentData{{Start: 0, End: 1.2, Text: "transcribed text"}},
				Usage:    contractsgw.TokenUsage{PromptTokens: 5, CompletionTokens: 6},
				Meta:     contractsgw.ResponseMeta{Provider: stringPtr("gemini"), Model: stringPtr("gemini-2.0-flash")},
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

		if !reflect.DeepEqual(got, gw.result) {
			t.Fatalf("Transcribe() response mismatch:\n got: %#v\nwant: %#v", got, gw.result)
		}

		wantErr := errors.New("transcription gateway failed")
		provider.UseTranscriptionGateway(&transcriptionGatewayStub{err: wantErr})

		if _, err := provider.Transcribe(ctx, contractsprovider.TranscriptionRequest{}); !errors.Is(err, wantErr) {
			t.Fatalf("Transcribe() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("reranking", func(t *testing.T) {
		t.Parallel()

		provider := cohere.NewProvider(nil)
		limit := 2
		model := "rerank-v3.5"
		gw := &rerankingGatewayStub{
			result: &contractsgw.RerankResult{
				Results: []contractsgw.RankedDocumentData{
					{Index: 1, Document: "second", Score: 0.9},
					{Index: 0, Document: "first", Score: 0.8},
				},
				Usage: contractsgw.TokenUsage{PromptTokens: 9, CompletionTokens: 10},
				Meta:  contractsgw.ResponseMeta{Provider: stringPtr("cohere"), Model: stringPtr("rerank-v3.5")},
			},
		}
		provider.UseRerankingGateway(gw)

		got, err := provider.Rerank(ctx, contractsprovider.RerankingRequest{
			Documents: []string{"first", "second"},
			Query:     "match",
			Limit:     &limit,
			Model:     &model,
		})

		if err != nil {
			t.Fatalf("Rerank() returned error: %v", err)
		}

		wantReq := contractsgw.RerankRequest{
			Model:     model,
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
		provider.UseRerankingGateway(&rerankingGatewayStub{err: wantErr})

		if _, err := provider.Rerank(ctx, contractsprovider.RerankingRequest{}); !errors.Is(err, wantErr) {
			t.Fatalf("Rerank() error = %v, want %v", err, wantErr)
		}
	})
}

func TestInventoryProviderFileAndStoreMapping(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	provider := openai.NewProvider(nil)

	fileGW := &fileGatewayStub{
		getResult: &contractsgw.FileGetResult{ID: "file-1", Filename: "notes.txt", Size: 12, Usage: contractsgw.TokenUsage{PromptTokens: 1}, Meta: contractsgw.ResponseMeta{Provider: stringPtr("openai")}},
		putResult: &contractsgw.FilePutResult{ID: "file-2", Filename: "notes.txt", Usage: contractsgw.TokenUsage{PromptTokens: 2}, Meta: contractsgw.ResponseMeta{Provider: stringPtr("openai")}},
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

	if got, err := provider.PutFile(ctx, contractsgw.StorableFile{Content: []byte("abc"), Filename: "notes.txt", MimeType: "text/plain"}); err != nil {
		t.Fatalf("PutFile() returned error: %v", err)
	} else if !reflect.DeepEqual(got, fileGW.putResult) {
		t.Fatalf("PutFile() response mismatch:\n got: %#v\nwant: %#v", got, fileGW.putResult)
	}

	if !reflect.DeepEqual(fileGW.putFile, contractsgw.StorableFile{Content: []byte("abc"), Filename: "notes.txt", MimeType: "text/plain"}) {
		t.Fatalf("PutFile() request mismatch:\n got: %#v\nwant: %#v", fileGW.putFile, contractsgw.StorableFile{Content: []byte("abc"), Filename: "notes.txt", MimeType: "text/plain"})
	}

	wantErr := errors.New("file gateway failed")
	provider.UseFileGateway(&fileGatewayStub{err: wantErr})

	if err := provider.DeleteFile(ctx, "file-9"); !errors.Is(err, wantErr) {
		t.Fatalf("DeleteFile() error = %v, want %v", err, wantErr)
	}

	storeGW := &storeGatewayStub{
		getResult:  &contractsgw.StoreData{ID: "store-1", Ready: true},
		createData: &contractsgw.StoreData{ID: "store-2", Ready: false},
		addResult:  &contractsgw.AddedDocumentData{ID: "doc-1", Filename: "notes.txt", Usage: contractsgw.TokenUsage{PromptTokens: 3}, Meta: contractsgw.ResponseMeta{Provider: stringPtr("openai")}},
	}
	provider.UseStoreGateway(storeGW)

	if got, err := provider.GetStore(ctx, "store-1"); err != nil {
		t.Fatalf("GetStore() returned error: %v", err)
	} else if !reflect.DeepEqual(got, storeGW.getResult) {
		t.Fatalf("GetStore() response mismatch:\n got: %#v\nwant: %#v", got, storeGW.getResult)
	}

	if storeGW.getID != "store-1" {
		t.Fatalf("GetStore() ID = %q, want %q", storeGW.getID, "store-1")
	}

	name := "My Store"
	description := "Store description"
	expiresIn := 3600

	if got, err := provider.CreateStore(ctx, contractsgw.StoreCreateRequest{
		Name:        &name,
		Description: &description,
		Files:       []contractsgw.StorableFile{{Filename: "notes.txt", MimeType: "text/plain"}},
		ExpiresIn:   &expiresIn,
	}); err != nil {
		t.Fatalf("CreateStore() returned error: %v", err)
	} else if !reflect.DeepEqual(got, storeGW.createData) {
		t.Fatalf("CreateStore() response mismatch:\n got: %#v\nwant: %#v", got, storeGW.createData)
	}

	if !reflect.DeepEqual(storeGW.createReq, contractsgw.StoreCreateRequest{
		Name:        &name,
		Description: &description,
		Files:       []contractsgw.StorableFile{{Filename: "notes.txt", MimeType: "text/plain"}},
		ExpiresIn:   &expiresIn,
	}) {
		t.Fatalf("CreateStore() request mismatch:\n got: %#v", storeGW.createReq)
	}

	if got, err := provider.AddFileToStore(ctx, "store-1", contractsgw.StorableFile{Content: []byte("abc"), Filename: "notes.txt", MimeType: "text/plain"}, map[string]any{"scope": "test"}); err != nil {
		t.Fatalf("AddFileToStore() returned error: %v", err)
	} else if !reflect.DeepEqual(got, storeGW.addResult) {
		t.Fatalf("AddFileToStore() response mismatch:\n got: %#v\nwant: %#v", got, storeGW.addResult)
	}

	if storeGW.storeID != "store-1" {
		t.Fatalf("AddFileToStore() store ID = %q, want %q", storeGW.storeID, "store-1")
	}

	wantErr = errors.New("store gateway failed")
	provider.UseStoreGateway(&storeGatewayStub{err: wantErr})

	if err := provider.RemoveFileFromStore(ctx, "store-1", "file-9"); !errors.Is(err, wantErr) {
		t.Fatalf("RemoveFileFromStore() error = %v, want %v", err, wantErr)
	}

	if err := provider.DeleteStore(ctx, "store-9"); !errors.Is(err, wantErr) {
		t.Fatalf("DeleteStore() error = %v, want %v", err, wantErr)
	}
}

func streamSeq(events ...contractsgw.StreamEvent) iter.Seq[contractsgw.StreamEvent] {
	return func(yield func(contractsgw.StreamEvent) bool) {
		for _, event := range events {
			if !yield(event) {
				return
			}
		}
	}
}

func streamEvent(kind, id string, _ any) streamEventStub {
	return streamEventStub{kind: kind, id: id}
}

func (e streamEventStub) EventType() string { return e.kind }

func stringPtr(v string) *string { return &v }
