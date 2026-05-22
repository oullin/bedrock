package failed

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Ref: @bedrock/code-0259
// job as a record inside a single JSON file on disk (newest first),
// capped at `limit` entries.
//
// the upstream provider supports an optional lockProviderResolver closure;
// the Go port keeps a process-local sync.
// the observable serialisation behaviour without dragging a lock
// provider interface into this package.
type FileFailedJobProvider struct {
	path  string
	limit int
	now   func() time.Time
	mu    sync.Mutex
}

// fileRecord matches the on-disk JSON shape the upstream file provider
// writes: id, connection, queue, payload, exception, failed_at
// (string "Y-m-d H:i:s"), failed_at_timestamp (Unix seconds).
type fileRecord struct {
	ID                string `json:"id"`
	Connection        string `json:"connection"`
	Queue             string `json:"queue"`
	Payload           string `json:"payload"`
	Exception         string `json:"exception"`
	FailedAt          string `json:"failed_at"`
	FailedAtTimestamp int64  `json:"failed_at_timestamp"`
}

// NewFileFailedJobProvider returns a provider that persists to the
// given path. A zero limit defaults to 100, matching upstream.
func NewFileFailedJobProvider(path string, limit int) *FileFailedJobProvider {
	if limit <= 0 {
		limit = 100
	}

	return &FileFailedJobProvider{path: path, limit: limit, now: time.Now}
}

// Log implements FailedJobProvider.
func (p *FileFailedJobProvider) Log(connection, queue, payload string, exception error) (string, error) {
	p.mu.Lock()

	defer p.mu.Unlock()

	id := extractUUID(payload)
	jobs, err := p.read()

	if err != nil {
		return "", err
	}

	now := p.now()
	entry := fileRecord{
		ID:                id,
		Connection:        connection,
		Queue:             queue,
		Payload:           payload,
		Exception:         errString(exception),
		FailedAt:          now.Format("2006-01-02 15:04:05"),
		FailedAtTimestamp: now.Unix(),
	}
	jobs = append([]fileRecord{entry}, jobs...)

	if len(jobs) > p.limit {
		jobs = jobs[:p.limit]
	}

	if err := p.write(jobs); err != nil {
		return "", err
	}

	return id, nil
}

// IDs implements FailedJobProvider.
func (p *FileFailedJobProvider) IDs(queueFilter string) ([]string, error) {
	p.mu.Lock()

	defer p.mu.Unlock()

	jobs, err := p.read()

	if err != nil {
		return nil, err
	}

	out := []string{}

	for _, j := range jobs {
		if queueFilter != "" && j.Queue != queueFilter {
			continue
		}

		out = append(out, j.ID)
	}

	return out, nil
}

// All implements FailedJobProvider.
func (p *FileFailedJobProvider) All() ([]FailedJob, error) {
	p.mu.Lock()

	defer p.mu.Unlock()

	jobs, err := p.read()

	if err != nil {
		return nil, err
	}

	out := make([]FailedJob, 0, len(jobs))

	for _, j := range jobs {
		out = append(out, j.toFailedJob())
	}

	return out, nil
}

// Find implements FailedJobProvider.
func (p *FileFailedJobProvider) Find(id string) (*FailedJob, error) {
	p.mu.Lock()

	defer p.mu.Unlock()

	jobs, err := p.read()

	if err != nil {
		return nil, err
	}

	for _, j := range jobs {
		if j.ID == id {
			fj := j.toFailedJob()

			return &fj, nil
		}
	}

	return nil, nil
}

// Forget implements FailedJobProvider.
func (p *FileFailedJobProvider) Forget(id string) (bool, error) {
	p.mu.Lock()

	defer p.mu.Unlock()

	jobs, err := p.read()

	if err != nil {
		return false, err
	}

	kept := make([]fileRecord, 0, len(jobs))
	removed := false

	for _, j := range jobs {
		if j.ID == id {
			removed = true

			continue
		}

		kept = append(kept, j)
	}

	if err := p.write(kept); err != nil {
		return false, err
	}

	return removed, nil
}

// Flush implements FailedJobProvider. It delegates to Prune using a
func (p *FileFailedJobProvider) Flush(hours int) error {
	cutoff := p.now().Add(-time.Duration(hours) * time.Hour)
	_, err := p.Prune(cutoff)

	return err
}

// Prune implements Prunable. It removes every entry whose
// failed_at_timestamp is <= before.Unix().
func (p *FileFailedJobProvider) Prune(before time.Time) (int64, error) {
	p.mu.Lock()

	defer p.mu.Unlock()

	jobs, err := p.read()

	if err != nil {
		return 0, err
	}

	cutoff := before.Unix()
	kept := make([]fileRecord, 0, len(jobs))

	for _, j := range jobs {
		if j.FailedAtTimestamp <= cutoff {
			continue
		}

		kept = append(kept, j)
	}

	if err := p.write(kept); err != nil {
		return 0, err
	}

	return int64(len(jobs) - len(kept)), nil
}

// Count implements Countable.
func (p *FileFailedJobProvider) Count(connection, queueFilter string) (int64, error) {
	p.mu.Lock()

	defer p.mu.Unlock()

	jobs, err := p.read()

	if err != nil {
		return 0, err
	}

	if connection == "" && queueFilter == "" {
		return int64(len(jobs)), nil
	}

	var n int64

	for _, j := range jobs {
		if connection != "" && j.Connection != connection {
			continue
		}

		if queueFilter != "" && j.Queue != queueFilter {
			continue
		}

		n++
	}

	return n, nil
}

func (p *FileFailedJobProvider) read() ([]fileRecord, error) {
	data, err := os.ReadFile(p.path)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	var jobs []fileRecord

	if err := json.Unmarshal(data, &jobs); err != nil {
		// Upstream silently treats malformed JSON as empty.
		return nil, nil
	}

	return jobs, nil
}

func (p *FileFailedJobProvider) write(jobs []fileRecord) error {
	if jobs == nil {
		jobs = []fileRecord{}
	}

	data, err := json.MarshalIndent(jobs, "", "    ")

	if err != nil {
		return err
	}

	return os.WriteFile(p.path, data, 0o644)
}

func (r fileRecord) toFailedJob() FailedJob {
	t := time.Unix(r.FailedAtTimestamp, 0)

	return FailedJob{
		ID:         r.ID,
		UUID:       r.ID,
		Connection: r.Connection,
		Queue:      r.Queue,
		Payload:    r.Payload,
		Exception:  r.Exception,
		FailedAt:   t,
	}
}
