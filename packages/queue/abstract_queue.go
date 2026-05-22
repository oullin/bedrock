package queue

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"reflect"
	"time"
)

// Ref: @bedrock/code-0268
// class — the shared base every concrete driver extends in PHP. In Go we
// expose it as a set of stateless helpers that drivers call directly,
// rather than as embedded state, so existing drivers can opt in without
// inheritance.
//
// The helpers cover:
//
//   - DisplayName   → Queue::getDisplayName
//   - NewUUIDv4     → Str::uuid (upstream leans on ramsey/uuid)
//   - CreatePayloadFor → Queue::createPayload (+ createPayloadUsing hooks)
//   - ShouldDispatchAfterCommit → Queue::shouldDispatchAfterCommit (lightweight
//     version; the full transaction-aware wiring lands in Step 14).

// Namer is implemented by job values that want to override the default
// reflect-based display name.
type Namer interface {
	QueueDisplayName() string
}

// AfterCommitMarker is implemented by jobs that should be dispatched
// only after the surrounding database transaction has committed. The
// concrete transaction plumbing arrives in Step 14; for now the marker
// exists so drivers can detect the intent without a compile-time dep
// on a higher-level transaction package.
//
// Ref: @bedrock/code-0197
type AfterCommitMarker interface {
	QueueAfterCommit() bool
}

// BeforeCommitMarker is the inverse of AfterCommitMarker: a job that
// Ref: @bedrock/code-0196
// override. Step 14 wires the full decision tree; today the type exists
// so drivers and tests can reference it.
type BeforeCommitMarker interface {
	QueueBeforeCommit() bool
}

// DisplayName returns the display name for a job value.
//
//   - A nil job yields "".
//   - A string is returned verbatim.
//   - A value implementing Namer uses the method's return.
//   - Everything else falls back to the reflect type path, e.g.
//     "github.com/acme/jobs.SendEmail".
func DisplayName(job any) string {
	if job == nil {
		return ""
	}

	if s, ok := job.(string); ok {
		return s
	}

	if n, ok := job.(Namer); ok {
		return n.QueueDisplayName()
	}

	t := reflect.TypeOf(job)

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if path := t.PkgPath(); path != "" {
		return path + "." + t.Name()
	}

	return t.String()
}

// NewUUIDv4 returns an RFC 4122 v4 UUID string. Used by CreatePayloadFor
// to populate Payload.UUID. Kept package-private-ish (exported so tests
// and callers that need a matching UUID generator can use it) but not
// intended as a general UUID library — add a proper dependency later if
// this grows legs.
func NewUUIDv4() string {
	var b [16]byte

	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant RFC 4122

	hexed := hex.EncodeToString(b[:])

	return fmt.Sprintf("%s-%s-%s-%s-%s", hexed[0:8], hexed[8:12], hexed[12:16], hexed[16:20], hexed[20:32])
}

// CreatePayloadFor builds, hook-applies, and serialises a payload ready
// Ref: @bedrock/code-0268
// Drivers should call this instead of constructing Payload structs by
// hand so that CreatePayloadUsing hooks fire uniformly, the UUID format
// matches upstream's, and timestamp/field shapes stay consistent.
//
// The returned *Payload is the hook-applied in-memory form; the []byte
// is the marshalled JSON. Both are returned because some drivers want
// the decoded form for bookkeeping (e.g. database `uuid` column) and
// others only need the raw bytes.
func CreatePayloadFor(connection, queueName string, job any, data map[string]any, opts JobOptions) (*Payload, []byte, error) {
	now := time.Now().UTC()

	if opts.BatchID != "" {
		if data == nil {
			data = make(map[string]any, 1)
		}

		data["batchId"] = opts.BatchID
	}

	p := &Payload{
		UUID:          NewUUIDv4(),
		DisplayName:   DisplayName(job),
		Job:           DisplayName(job),
		Data:          data,
		Tries:         0,
		MaxTries:      opts.MaxTries,
		Timeout:       int(opts.Timeout / time.Second),
		Backoff:       backoffSeconds(opts.Backoff),
		MaxExceptions: opts.MaxExceptions,
		CreatedAt:     &now,
	}

	if !opts.RetryUntil.IsZero() {
		retry := opts.RetryUntil
		p.RetryUntil = &retry
	}

	ApplyPayloadHooks(connection, queueName, p)

	raw, err := p.Marshal()

	if err != nil {
		return nil, nil, NewInvalidPayloadError(fmt.Sprintf("queue: marshal payload: %v", err), p)
	}

	return p, raw, nil
}

// ShouldDispatchAfterCommit reports whether a dispatch of job on
// connection connName should wait for the enclosing transaction to
// commit. Precedence (matching upstream):
//
//  1. A BeforeCommitMarker that returns true forces before-commit.
//  2. An AfterCommitMarker that returns true forces after-commit.
//  3. The connection config's `after_commit` key, if present, wins next.
//  4. Default: false.
//
// The actual transaction integration (registering a callback on a
// TxManager) is wired up in Step 14.
func ShouldDispatchAfterCommit(job any, config map[string]any) bool {
	if b, ok := job.(BeforeCommitMarker); ok && b.QueueBeforeCommit() {
		return false
	}

	if a, ok := job.(AfterCommitMarker); ok && a.QueueAfterCommit() {
		return true
	}

	if config != nil {
		if v, ok := config["after_commit"]; ok {
			if b, ok := v.(bool); ok {
				return b
			}
		}
	}

	return false
}

// backoffSeconds converts a slice of durations into the per-second
// integer slice that upstream payloads carry.
func backoffSeconds(d []time.Duration) []int {
	if len(d) == 0 {
		return nil
	}

	out := make([]int, len(d))

	for i, v := range d {
		out[i] = int(v / time.Second)
	}

	return out
}
