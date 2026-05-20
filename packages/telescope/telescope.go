package telescope

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

// FilterFunc is a predicate that returns true if the entry should be recorded.
// When any filter returns false the entry is discarded.
type FilterFunc func(*IncomingEntry) bool

// FilterBatchFunc is a predicate applied to the full batch of queued entries
// before they are stored. Return false to prevent the entire batch from being
// persisted.
type FilterBatchFunc func([]*IncomingEntry) bool

// TagFunc returns additional tags to attach to an entry after it has passed
// all recording filters.
type TagFunc func(*IncomingEntry) []string

// AfterRecordingFunc is called synchronously after an entry is queued.
type AfterRecordingFunc func(*IncomingEntry)

// AfterStoringFunc is called after entries are successfully persisted. It
// receives the current batchID and the slice of stored entries.
type AfterStoringFunc func(batchID string, entries []*IncomingEntry)

// Option configures a Telescope instance.
type Option func(*Telescope)

// WithRepository sets the persistence backend.

// WithHiddenRequestHeaders registers header names whose values will be masked
// with "********" in recorded request entries.

// WithHiddenRequestParameters registers parameter names whose values will be
// masked in recorded request payload.

// WithHiddenResponseParameters registers parameter names whose values will be
// masked in recorded response bodies.

// Telescope is the central hub that collects, filters, tags, and stores
// telemetry entries. It mirrors the behaviour of the upstream Telescope class.
//
// A single Telescope instance should be created per application and shared
// with all Watchers.
type Telescope struct {
	mu sync.Mutex

	repository Repository

	entriesQueue []*IncomingEntry
	updatesQueue []*EntryUpdate

	filterCallbacks         []FilterFunc
	filterBatchCallbacks    []FilterBatchFunc
	tagCallbacks            []TagFunc
	afterRecordingCallbacks []AfterRecordingFunc
	afterStoringCallbacks   []AfterStoringFunc

	hiddenRequestHeaders     []string
	hiddenRequestParameters  []string
	hiddenResponseParameters []string

	recording bool
	paused    bool

	batchID string
}

func WithRepository(r Repository) Option {
	return func(t *Telescope) { t.repository = r }
}

func WithHiddenRequestHeaders(headers ...string) Option {
	return func(t *Telescope) {
		t.hiddenRequestHeaders = append(t.hiddenRequestHeaders, headers...)
	}
}

func WithHiddenRequestParameters(params ...string) Option {
	return func(t *Telescope) {
		t.hiddenRequestParameters = append(t.hiddenRequestParameters, params...)
	}
}

func WithHiddenResponseParameters(params ...string) Option {
	return func(t *Telescope) {
		t.hiddenResponseParameters = append(t.hiddenResponseParameters, params...)
	}
}

// New creates a Telescope instance with the given options. Recording is
// disabled by default; call StartRecording to begin capturing entries.
func New(opts ...Option) *Telescope {
	t := &Telescope{
		hiddenRequestHeaders: []string{
			"authorization",
			"php-auth-pw",
		},
		hiddenRequestParameters: []string{
			"password",
			"password_confirmation",
		},
	}

	for _, opt := range opts {
		opt(t)
	}

	return t
}

// ─── Recording control ──────────────────────────────────────────────────────

// StartRecording enables entry recording and assigns a new batch UUID.
func (t *Telescope) StartRecording() {
	t.mu.Lock()

	defer t.mu.Unlock()

	t.recording = true
	t.paused = false
	t.batchID = uuid.New().String()
}

// StopRecording disables entry recording.
func (t *Telescope) StopRecording() {
	t.mu.Lock()

	defer t.mu.Unlock()

	t.recording = false
}

// PauseRecording temporarily halts recording without resetting the batch.
func (t *Telescope) PauseRecording() {
	t.mu.Lock()

	defer t.mu.Unlock()

	t.paused = true
}

// ResumeRecording re-enables recording after a pause.
func (t *Telescope) ResumeRecording() {
	t.mu.Lock()

	defer t.mu.Unlock()

	t.paused = false
}

// IsRecording reports whether entries are currently being captured.
func (t *Telescope) IsRecording() bool {
	t.mu.Lock()

	defer t.mu.Unlock()

	return t.recording && !t.paused
}

// WithoutRecording executes fn while recording is paused, then resumes.
func (t *Telescope) WithoutRecording(fn func()) {
	t.PauseRecording()

	defer t.ResumeRecording()

	fn()
}

// NewBatch generates a fresh batch UUID, replaces the current one, and returns
// it. Call this at the start of each logical request/command.
func (t *Telescope) NewBatch() string {
	t.mu.Lock()

	defer t.mu.Unlock()

	t.batchID = uuid.New().String()

	return t.batchID
}

// CurrentBatchID returns the active batch UUID.
func (t *Telescope) CurrentBatchID() string {
	t.mu.Lock()

	defer t.mu.Unlock()

	return t.batchID
}

// ─── Filter / tag configuration ─────────────────────────────────────────────

// Filter registers a callback that controls whether an individual entry is
// recorded. If any registered filter returns false the entry is dropped.
func (t *Telescope) Filter(fn FilterFunc) {
	t.mu.Lock()

	defer t.mu.Unlock()

	t.filterCallbacks = append(t.filterCallbacks, fn)
}

// FilterBatch registers a callback that controls whether the current batch of
// queued entries is stored. Receives the full slice; return false to drop all.
func (t *Telescope) FilterBatch(fn FilterBatchFunc) {
	t.mu.Lock()

	defer t.mu.Unlock()

	t.filterBatchCallbacks = append(t.filterBatchCallbacks, fn)
}

// Tag registers a callback that returns additional tags for an entry after it
// passes all filters.
func (t *Telescope) Tag(fn TagFunc) {
	t.mu.Lock()

	defer t.mu.Unlock()

	t.tagCallbacks = append(t.tagCallbacks, fn)
}

// AfterRecording registers a callback invoked synchronously after an entry is
// appended to the queue.
func (t *Telescope) AfterRecording(fn AfterRecordingFunc) {
	t.mu.Lock()

	defer t.mu.Unlock()

	t.afterRecordingCallbacks = append(t.afterRecordingCallbacks, fn)
}

// AfterStoring registers a callback invoked after entries are successfully
// persisted. Receives the batch UUID and the stored entries.
func (t *Telescope) AfterStoring(fn AfterStoringFunc) {
	t.mu.Lock()

	defer t.mu.Unlock()

	t.afterStoringCallbacks = append(t.afterStoringCallbacks, fn)
}

// HiddenRequestHeaders returns the list of request header names that will be
// masked in recorded entries.
func (t *Telescope) HiddenRequestHeaders() []string {
	t.mu.Lock()

	defer t.mu.Unlock()

	out := make([]string, len(t.hiddenRequestHeaders))
	copy(out, t.hiddenRequestHeaders)

	return out
}

// HiddenRequestParameters returns the list of request parameter names that
// will be masked in recorded entries.
func (t *Telescope) HiddenRequestParameters() []string {
	t.mu.Lock()

	defer t.mu.Unlock()

	out := make([]string, len(t.hiddenRequestParameters))
	copy(out, t.hiddenRequestParameters)

	return out
}

// HiddenResponseParameters returns the list of response parameter names that
// will be masked in recorded entries.
func (t *Telescope) HiddenResponseParameters() []string {
	t.mu.Lock()

	defer t.mu.Unlock()

	out := make([]string, len(t.hiddenResponseParameters))
	copy(out, t.hiddenResponseParameters)

	return out
}

// ─── Recording ──────────────────────────────────────────────────────────────

// Record queues a single entry if recording is active and all filters pass.
// It stamps the entry with the current batch UUID, applies tag callbacks, and
// fires after-recording hooks. This is the central recording path that all
// type-specific Record* methods delegate to.
func (t *Telescope) Record(entry *IncomingEntry) {
	t.mu.Lock()

	if !t.recording || t.paused {
		t.mu.Unlock()

		return
	}

	entry.BatchID = t.batchID

	// Apply per-entry filters.
	for _, fn := range t.filterCallbacks {
		if !fn(entry) {
			t.mu.Unlock()

			return
		}
	}

	// Apply tag callbacks.
	for _, fn := range t.tagCallbacks {
		entry.AddTags(fn(entry)...)
	}

	t.entriesQueue = append(t.entriesQueue, entry)

	// Snapshot callbacks to call outside the lock.
	callbacks := make([]AfterRecordingFunc, len(t.afterRecordingCallbacks))
	copy(callbacks, t.afterRecordingCallbacks)

	t.mu.Unlock()

	for _, fn := range callbacks {
		fn(entry)
	}
}

// RecordRequest records an HTTP request entry.
func (t *Telescope) RecordRequest(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeRequest))
}

// RecordQuery records a database query entry.
func (t *Telescope) RecordQuery(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeQuery))
}

// RecordException records an exception entry.
func (t *Telescope) RecordException(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeException))
}

// RecordLog records a log message entry.
func (t *Telescope) RecordLog(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeLog))
}

// RecordEvent records an application event entry.
func (t *Telescope) RecordEvent(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeEvent))
}

// RecordJob records a queued job entry.
func (t *Telescope) RecordJob(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeJob))
}

// RecordCache records a cache operation entry.
func (t *Telescope) RecordCache(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeCache))
}

// RecordMail records a mail message entry.
func (t *Telescope) RecordMail(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeMail))
}

// RecordNotification records a notification entry.
func (t *Telescope) RecordNotification(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeNotification))
}

// RecordModel records a model lifecycle entry.
func (t *Telescope) RecordModel(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeModel))
}

// RecordView records a view render entry.
func (t *Telescope) RecordView(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeView))
}

// RecordCommand records an Artisan/CLI command entry.
func (t *Telescope) RecordCommand(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeCommand))
}

// RecordScheduledTask records a scheduled task execution entry.
func (t *Telescope) RecordScheduledTask(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeScheduledTask))
}

// RecordRedis records a Redis command entry.
func (t *Telescope) RecordRedis(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeRedis))
}

// RecordGate records an authorization gate check entry.
func (t *Telescope) RecordGate(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeGate))
}

// RecordClientRequest records an outbound HTTP client request entry.
func (t *Telescope) RecordClientRequest(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeClientRequest))
}

// RecordDump records a debug dump entry.
func (t *Telescope) RecordDump(entry *IncomingEntry) {
	t.Record(entry.WithType(EntryTypeDump))
}

// ─── Storage ────────────────────────────────────────────────────────────────

// Store persists the current entries queue to the repository and fires
// after-storing callbacks. The queue is cleared regardless of whether storage
// succeeds. Mirrors the upstream Telescope::store().
func (t *Telescope) Store(ctx context.Context) error {
	t.mu.Lock()

	entries := t.entriesQueue
	updates := t.updatesQueue
	batchID := t.batchID
	t.entriesQueue = nil
	t.updatesQueue = nil

	// Apply batch-level filters.
	for _, fn := range t.filterBatchCallbacks {
		if !fn(entries) {
			t.mu.Unlock()

			return nil
		}
	}

	afterStoring := make([]AfterStoringFunc, len(t.afterStoringCallbacks))
	copy(afterStoring, t.afterStoringCallbacks)

	t.mu.Unlock()

	if t.repository == nil {
		return nil
	}

	if len(entries) > 0 {
		if err := t.repository.Store(entries); err != nil {
			return err
		}
	}

	if len(updates) > 0 {
		if err := t.repository.Update(updates); err != nil {
			return err
		}
	}

	for _, fn := range afterStoring {
		fn(batchID, entries)
	}

	return nil
}

// Update queues a set of entry mutations to be flushed on the next Store call.
func (t *Telescope) Update(updates []*EntryUpdate) {
	t.mu.Lock()

	defer t.mu.Unlock()

	t.updatesQueue = append(t.updatesQueue, updates...)
}

// Flush discards all queued entries and updates without persisting them.
func (t *Telescope) Flush() {
	t.mu.Lock()

	defer t.mu.Unlock()

	t.entriesQueue = nil
	t.updatesQueue = nil
}

// QueuedEntries returns a snapshot of the currently queued entries (useful in
// tests to inspect state before Store is called).
func (t *Telescope) QueuedEntries() []*IncomingEntry {
	t.mu.Lock()

	defer t.mu.Unlock()

	out := make([]*IncomingEntry, len(t.entriesQueue))
	copy(out, t.entriesQueue)

	return out
}

// Repository returns the configured persistence backend.
func (t *Telescope) Repository() Repository {
	t.mu.Lock()

	defer t.mu.Unlock()

	return t.repository
}
