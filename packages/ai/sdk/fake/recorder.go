// Package fake provides in-memory gateway implementations for testing.
// Every fake gateway records interactions and exposes assert helpers that
// mirror the upstream AI test assertion API.
package fake

import (
	"fmt"
	"sync"

	"github.com/bedrock/packages/ai/sdk/prompts"
)

// TestingT is satisfied by *testing.T without requiring an import of the
// testing package in production code.
type TestingT interface {
	Helper()
	Errorf(format string, args ...any)
}

// agentRecord stores a single recorded agent prompt interaction.
type agentRecord struct {
	prompt *prompts.AgentPrompt
	queued bool
}

// imageRecord stores a single recorded image generation interaction.
type imageRecord struct {
	prompt *prompts.ImagePrompt
	queued bool
}

// audioRecord stores a single recorded audio generation interaction.
type audioRecord struct {
	prompt *prompts.AudioPrompt
	queued bool
}

// embeddingsRecord stores a single recorded embeddings generation interaction.
type embeddingsRecord struct {
	prompt *prompts.EmbeddingsPrompt
	queued bool
}

// transcriptionRecord stores a single recorded transcription interaction.
type transcriptionRecord struct {
	prompt *prompts.TranscriptionPrompt
	queued bool
}

// rerankingRecord stores a single recorded reranking interaction.
type rerankingRecord struct {
	prompt *prompts.RerankingPrompt
}

// fileRecord stores a single recorded file operation.
type fileRecord struct {
	op       string // "get", "put", "delete"
	id       string
	filename string
}

// storeRecord stores a single recorded store operation.
type storeRecord struct {
	op      string // "get", "create", "add_file", "remove_file", "delete"
	storeID string
	fileID  string
}

// Recorder is the shared interaction log used by all fake gateways.
// A single Recorder instance is shared across all fakes within a Manager
// so assertions can be made against any capability uniformly.
type Recorder struct {
	mu             sync.Mutex
	agents         []agentRecord
	images         []imageRecord
	audio          []audioRecord
	embeddings     []embeddingsRecord
	transcriptions []transcriptionRecord
	rerankings     []rerankingRecord
	files          []fileRecord
	stores         []storeRecord
}

// NewRecorder creates a fresh Recorder.
func NewRecorder() *Recorder { return &Recorder{} }

// Reset clears all recorded interactions.
func (r *Recorder) Reset() {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.agents = nil
	r.images = nil
	r.audio = nil
	r.embeddings = nil
	r.transcriptions = nil
	r.rerankings = nil
	r.files = nil
	r.stores = nil
}

// --- Agent recording ---

func (r *Recorder) recordAgent(p *prompts.AgentPrompt, queued bool) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.agents = append(r.agents, agentRecord{prompt: p, queued: queued})
}

// AssertAgentWasPrompted fails t if no recorded prompt matches fn.
func (r *Recorder) AssertAgentWasPrompted(t TestingT, fn func(*prompts.AgentPrompt) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.agents {
		if !rec.queued && fn(rec.prompt) {
			return
		}
	}

	t.Errorf("ai: expected agent to have been prompted but no matching prompt was found")
}

// AssertAgentNotPrompted fails t if any recorded prompt matches fn.
func (r *Recorder) AssertAgentNotPrompted(t TestingT, fn func(*prompts.AgentPrompt) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.agents {
		if !rec.queued && fn(rec.prompt) {
			t.Errorf("ai: expected agent not to have been prompted but a matching prompt was found")

			return
		}
	}
}

// AssertAgentNeverPrompted fails t if any prompt was recorded.
func (r *Recorder) AssertAgentNeverPrompted(t TestingT) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.agents {
		if !rec.queued {
			t.Errorf("ai: expected agent to never have been prompted but %d prompt(s) were recorded", len(r.agents))

			return
		}
	}
}

// AssertAgentWasQueued fails t if no queued prompt matches fn.
func (r *Recorder) AssertAgentWasQueued(t TestingT, fn func(*prompts.AgentPrompt) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.agents {
		if rec.queued && fn(rec.prompt) {
			return
		}
	}

	t.Errorf("ai: expected agent to have been queued but no matching queued prompt was found")
}

// AssertAgentNotQueued fails t if any queued prompt matches fn.
func (r *Recorder) AssertAgentNotQueued(t TestingT, fn func(*prompts.AgentPrompt) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.agents {
		if rec.queued && fn(rec.prompt) {
			t.Errorf("ai: expected agent not to have been queued but a matching queued prompt was found")

			return
		}
	}
}

// AssertAgentNeverQueued fails t if any prompt was queued.
func (r *Recorder) AssertAgentNeverQueued(t TestingT) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.agents {
		if rec.queued {
			t.Errorf("ai: expected agent to never have been queued but queued prompts were recorded")

			return
		}
	}
}

// --- Image recording ---

func (r *Recorder) recordImage(p *prompts.ImagePrompt, queued bool) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.images = append(r.images, imageRecord{prompt: p, queued: queued})
}

// AssertImageGenerated fails t if no image prompt matches fn.
func (r *Recorder) AssertImageGenerated(t TestingT, fn func(*prompts.ImagePrompt) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.images {
		if !rec.queued && fn(rec.prompt) {
			return
		}
	}

	t.Errorf("ai: expected image to have been generated but no matching prompt was found")
}

// AssertImageNotGenerated fails t if any image prompt matches fn.
func (r *Recorder) AssertImageNotGenerated(t TestingT, fn func(*prompts.ImagePrompt) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.images {
		if !rec.queued && fn(rec.prompt) {
			t.Errorf("ai: expected image not to have been generated but a matching prompt was found")

			return
		}
	}
}

// AssertNothingImageGenerated fails t if any image was generated.
func (r *Recorder) AssertNothingImageGenerated(t TestingT) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.images {
		if !rec.queued {
			t.Errorf("ai: expected no images to have been generated but %d generation(s) were recorded", len(r.images))

			return
		}
	}
}

// AssertImageQueued fails t if no queued image prompt matches fn.
func (r *Recorder) AssertImageQueued(t TestingT, fn func(*prompts.ImagePrompt) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.images {
		if rec.queued && fn(rec.prompt) {
			return
		}
	}

	t.Errorf("ai: expected image to have been queued but no matching queued prompt was found")
}

// AssertNothingImageQueued fails t if any image was queued.
func (r *Recorder) AssertNothingImageQueued(t TestingT) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.images {
		if rec.queued {
			t.Errorf("ai: expected no images to have been queued but queued generations were recorded")

			return
		}
	}
}

// --- Audio recording ---

func (r *Recorder) recordAudio(p *prompts.AudioPrompt, queued bool) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.audio = append(r.audio, audioRecord{prompt: p, queued: queued})
}

// AssertAudioGenerated fails t if no audio prompt matches fn.
func (r *Recorder) AssertAudioGenerated(t TestingT, fn func(*prompts.AudioPrompt) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.audio {
		if !rec.queued && fn(rec.prompt) {
			return
		}
	}

	t.Errorf("ai: expected audio to have been generated but no matching prompt was found")
}

// AssertNothingAudioGenerated fails t if any audio was generated.
func (r *Recorder) AssertNothingAudioGenerated(t TestingT) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.audio {
		if !rec.queued {
			t.Errorf("ai: expected no audio to have been generated but generations were recorded")

			return
		}
	}
}

// --- Embeddings recording ---

func (r *Recorder) recordEmbeddings(p *prompts.EmbeddingsPrompt, queued bool) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.embeddings = append(r.embeddings, embeddingsRecord{prompt: p, queued: queued})
}

// AssertEmbeddingsGenerated fails t if no embeddings prompt matches fn.
func (r *Recorder) AssertEmbeddingsGenerated(t TestingT, fn func(*prompts.EmbeddingsPrompt) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.embeddings {
		if !rec.queued && fn(rec.prompt) {
			return
		}
	}

	t.Errorf("ai: expected embeddings to have been generated but no matching prompt was found")
}

// AssertNothingEmbeddingsGenerated fails t if any embeddings were generated.
func (r *Recorder) AssertNothingEmbeddingsGenerated(t TestingT) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.embeddings {
		if !rec.queued {
			t.Errorf("ai: expected no embeddings to have been generated but generations were recorded")

			return
		}
	}
}

// --- Transcription recording ---

func (r *Recorder) recordTranscription(p *prompts.TranscriptionPrompt, queued bool) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.transcriptions = append(r.transcriptions, transcriptionRecord{prompt: p, queued: queued})
}

// AssertTranscriptionGenerated fails t if no transcription prompt matches fn.
func (r *Recorder) AssertTranscriptionGenerated(t TestingT, fn func(*prompts.TranscriptionPrompt) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.transcriptions {
		if !rec.queued && fn(rec.prompt) {
			return
		}
	}

	t.Errorf("ai: expected transcription to have been generated but no matching prompt was found")
}

// AssertNothingTranscriptionGenerated fails t if any transcription was generated.
func (r *Recorder) AssertNothingTranscriptionGenerated(t TestingT) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.transcriptions {
		if !rec.queued {
			t.Errorf("ai: expected no transcriptions to have been generated but generations were recorded")

			return
		}
	}
}

// --- Reranking recording ---

func (r *Recorder) recordReranking(p *prompts.RerankingPrompt) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.rerankings = append(r.rerankings, rerankingRecord{prompt: p})
}

// AssertReranked fails t if no reranking prompt matches fn.
func (r *Recorder) AssertReranked(t TestingT, fn func(*prompts.RerankingPrompt) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.rerankings {
		if fn(rec.prompt) {
			return
		}
	}

	t.Errorf("ai: expected reranking to have been performed but no matching prompt was found")
}

// AssertNothingReranked fails t if any reranking was performed.
func (r *Recorder) AssertNothingReranked(t TestingT) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	if len(r.rerankings) > 0 {
		t.Errorf("ai: expected no rerankings to have been performed but %d were recorded", len(r.rerankings))
	}
}

// --- Store recording ---

func (r *Recorder) recordStore(op, storeID, fileID string) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.stores = append(r.stores, storeRecord{op: op, storeID: storeID, fileID: fileID})
}

// AssertStoreFileAdded fails t if no add_file operation matches fn (receives fileID).
func (r *Recorder) AssertStoreFileAdded(t TestingT, fn func(fileID string) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.stores {
		if rec.op == "add_file" && fn(rec.fileID) {
			return
		}
	}

	t.Errorf("ai: expected a file to have been added to a store but no matching operation was found")
}

// AssertStoreFileRemoved fails t if no remove_file operation matches fn (receives fileID).
func (r *Recorder) AssertStoreFileRemoved(t TestingT, fn func(fileID string) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.stores {
		if rec.op == "remove_file" && fn(rec.fileID) {
			return
		}
	}

	t.Errorf("ai: expected a file to have been removed from a store but no matching operation was found")
}

// --- File recording ---

func (r *Recorder) recordFile(op, id, filename string) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.files = append(r.files, fileRecord{op: op, id: id, filename: filename})
}

// AssertFileStored fails t if no stored file matches fn (receives filename).
func (r *Recorder) AssertFileStored(t TestingT, fn func(filename string) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.files {
		if rec.op == "put" && fn(rec.filename) {
			return
		}
	}

	t.Errorf("ai: expected a file to have been stored but no matching operation was found")
}

// AssertNothingFileStored fails t if any file was stored.
func (r *Recorder) AssertNothingFileStored(t TestingT) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.files {
		if rec.op == "put" {
			t.Errorf("ai: expected no files to have been stored but storage operations were recorded")

			return
		}
	}
}

// AssertFileDeleted fails t if no deleted file matches fn (receives id).
func (r *Recorder) AssertFileDeleted(t TestingT, fn func(id string) bool) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.files {
		if rec.op == "delete" && fn(rec.id) {
			return
		}
	}

	t.Errorf("ai: expected a file to have been deleted but no matching operation was found")
}

// AssertNothingFileDeleted fails t if any file was deleted.
func (r *Recorder) AssertNothingFileDeleted(t TestingT) {
	t.Helper()
	r.mu.Lock()

	defer r.mu.Unlock()

	for _, rec := range r.files {
		if rec.op == "delete" {
			t.Errorf("ai: expected no files to have been deleted but delete operations were recorded")

			return
		}
	}
}

// RecordQueuedPrompt records an agent prompt that was dispatched to a queue.
// This is called by the AnonymousAgent.Queue method to track queue calls in fake mode.
func (r *Recorder) RecordQueuedPrompt(p *prompts.AgentPrompt) {
	r.recordAgent(p, true)
}

// FakeFileID generates a deterministic fake file ID for testing.
func FakeFileID() string {
	return fmt.Sprintf("file_%d", fakeIDCounter.Add(1))
}
