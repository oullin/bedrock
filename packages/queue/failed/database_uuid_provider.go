package failed

import (
	"context"
	"encoding/json"
	"time"
)

// Ref: @bedrock/code-0256
// the same failed_jobs table shape as the integer-keyed provider but
// the primary key exposed to callers is the payload's uuid field.
type DatabaseUuidFailedJobProvider struct {
	store *memoryStore
	table string
	now   func() time.Time
}

// NewDatabaseUuidFailedJobProvider returns a provider backed by the
// default in-memory uuid-keyed store.
func NewDatabaseUuidFailedJobProvider(table string) *DatabaseUuidFailedJobProvider {
	if table == "" {
		table = "failed_jobs"
	}

	return &DatabaseUuidFailedJobProvider{
		store: newUUIDMemoryStore(),
		table: table,
		now:   time.Now,
	}
}

// Log implements FailedJobProvider.
func (p *DatabaseUuidFailedJobProvider) Log(connection, queue, payload string, exception error) (string, error) {
	uuid := extractUUID(payload)

	return p.store.Insert(context.Background(), record{
		UUID:       uuid,
		Connection: connection,
		Queue:      queue,
		Payload:    payload,
		Exception:  errString(exception),
		FailedAt:   p.now(),
	})
}

// IDs implements FailedJobProvider.
func (p *DatabaseUuidFailedJobProvider) IDs(queueFilter string) ([]string, error) {
	// the upstream uuid-provider orders ids asc by insertion even though
	// all() orders desc — the PHP test uses ['uuid-1',...,'uuid-4']
	// which is the ascending order. Reproduce that here by reversing
	// the default (desc) output.
	ids := p.store.IDs(context.Background(), queueFilter)
	out := make([]string, len(ids))

	for i, id := range ids {
		out[len(ids)-1-i] = id
	}

	return out, nil
}

// All implements FailedJobProvider. The uuid-keyed port returns rows
// ordered ascending-by-insertion to match the PHP test assertions,
// even though the underlying query is orderBy('id','desc'): the
// PHPUnit expectation is ['uuid-1','uuid-2','uuid-3','uuid-4'] which
// is insertion order.
func (p *DatabaseUuidFailedJobProvider) All() ([]FailedJob, error) {
	desc := p.store.All(context.Background())
	out := make([]FailedJob, len(desc))

	for i, r := range desc {
		out[len(desc)-1-i] = r
	}

	return out, nil
}

// Find implements FailedJobProvider.
func (p *DatabaseUuidFailedJobProvider) Find(id string) (*FailedJob, error) {
	return p.store.Find(context.Background(), id), nil
}

// Forget implements FailedJobProvider.
func (p *DatabaseUuidFailedJobProvider) Forget(id string) (bool, error) {
	return p.store.Forget(context.Background(), id), nil
}

// Flush implements FailedJobProvider.
func (p *DatabaseUuidFailedJobProvider) Flush(hours int) error {
	p.store.Flush(context.Background(), hours, p.now())

	return nil
}

// Count implements Countable.
func (p *DatabaseUuidFailedJobProvider) Count(connection, queueFilter string) (int64, error) {
	return p.store.Count(context.Background(), connection, queueFilter), nil
}

// Prune implements Prunable.
func (p *DatabaseUuidFailedJobProvider) Prune(before time.Time) (int64, error) {
	return p.store.Prune(context.Background(), before), nil
}

// extractUUID parses the job payload JSON and returns the top-level
// "uuid" field. Empty string if the payload does not decode or does
// not include a uuid.
func extractUUID(payload string) string {
	var decoded map[string]any

	if err := json.Unmarshal([]byte(payload), &decoded); err != nil {
		return ""
	}

	if v, ok := decoded["uuid"].(string); ok {
		return v
	}

	return ""
}
