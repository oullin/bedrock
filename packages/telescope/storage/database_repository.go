package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bedrock/packages/telescope"
)

// DatabaseRepository persists Telescope entries to a relational database. It
// mirrors the upstream DatabaseEntriesRepository, supporting the same three-table
// schema: telescope_entries, telescope_entries_tags, and telescope_monitoring.
//
// The connection must be a *database/sql.DB configured with the appropriate
// driver (SQLite, MySQL, PostgreSQL).
type DatabaseRepository struct {
	mu            sync.RWMutex
	db            *sql.DB
	chunkSize     int
	monitoredTags map[string]struct{}
}

// NewDatabaseRepository creates a DatabaseRepository backed by db.
// chunkSize controls how many entries are inserted per transaction (0 uses
// the default of 1000).

// compile-time interface satisfaction check.

// Migrate creates the Telescope database tables if they do not already exist.
// Suitable for SQLite and most SQL databases. For production use, prefer
// applying migrations through your migration tool.

// Find retrieves a single entry by UUID.

// Get retrieves entries of the given type filtered by EntryQueryOptions.

// Load tags for all entries.

// Store persists a batch of incoming entries in chunks.

// storeChunk persists a slice of entries in a single transaction.

//nolint:errcheck

// Ignore duplicate tag constraint violations.

// Update applies field mutations to stored entries.

//nolint:errcheck

// Load current content.

// Merge changes.

// Add tags.

// Remove tags.

// LoadMonitoredTags loads monitored tags from the database into memory.

// IsMonitoring reports whether any of the given tags are monitored.

// Monitoring returns the list of currently monitored tags.

// Monitor activates monitoring for the given tags.

//nolint:errcheck

// StopMonitoring deactivates monitoring for the given tags.

//nolint:errcheck

// Clear removes all stored entries, tags, and monitored tags.

// Prune removes entries older than before. When keepExceptions is true,
// exception entries are retained regardless of age.

// ─── Helpers ─────────────────────────────────────────────────────────────────

type scanner interface {
	Scan(dest ...any) error
}

const defaultChunkSize = 1000

func NewDatabaseRepository(db *sql.DB, chunkSize int) *DatabaseRepository {
	if chunkSize <= 0 {
		chunkSize = defaultChunkSize
	}

	return &DatabaseRepository{
		db:            db,
		chunkSize:     chunkSize,
		monitoredTags: make(map[string]struct{}),
	}
}

var _ telescope.Repository = (*DatabaseRepository)(nil)

func (r *DatabaseRepository) Migrate() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS telescope_entries (
			sequence    INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid        TEXT    NOT NULL UNIQUE,
			batch_id    TEXT    NOT NULL,
			family_hash TEXT,
			type        TEXT    NOT NULL,
			content     TEXT    NOT NULL,
			created_at  TEXT    NOT NULL
		);

		CREATE TABLE IF NOT EXISTS telescope_entries_tags (
			entry_uuid TEXT NOT NULL REFERENCES telescope_entries(uuid) ON DELETE CASCADE,
			tag        TEXT NOT NULL,
			PRIMARY KEY (entry_uuid, tag)
		);

		CREATE INDEX IF NOT EXISTS telescope_entries_tags_tag_idx
			ON telescope_entries_tags (tag);

		CREATE TABLE IF NOT EXISTS telescope_monitoring (
			tag TEXT PRIMARY KEY
		);
	`)

	return err
}

func (r *DatabaseRepository) Find(id string) (*telescope.EntryResult, error) {
	row := r.db.QueryRow(
		`SELECT sequence, uuid, batch_id, family_hash, type, content, created_at
		   FROM telescope_entries WHERE uuid = ?`, id,
	)

	entry, err := r.scanEntry(row)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	tags, err := r.loadTags(id)

	if err != nil {
		return nil, err
	}

	entry.Tags = tags

	return entry, nil
}

func (r *DatabaseRepository) Get(entryType string, opts telescope.EntryQueryOptions) ([]*telescope.EntryResult, error) {
	limit := opts.Limit

	if limit <= 0 {
		limit = 50
	}

	var (
		conditions []string
		args       []any
	)

	if entryType != "" {
		conditions = append(conditions, "e.type = ?")
		args = append(args, entryType)
	}

	if opts.BatchID != "" {
		conditions = append(conditions, "e.batch_id = ?")
		args = append(args, opts.BatchID)
	}

	if opts.FamilyHash != "" {
		conditions = append(conditions, "e.family_hash = ?")
		args = append(args, opts.FamilyHash)
	}

	if opts.BeforeSequence > 0 {
		conditions = append(conditions, "e.sequence < ?")
		args = append(args, opts.BeforeSequence)
	}

	if opts.Tag != "" {
		conditions = append(conditions, "EXISTS (SELECT 1 FROM telescope_entries_tags t WHERE t.entry_uuid = e.uuid AND t.tag = ?)")
		args = append(args, opts.Tag)
	}

	where := ""

	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	args = append(args, limit)

	rows, err := r.db.Query(
		fmt.Sprintf(`SELECT e.sequence, e.uuid, e.batch_id, e.family_hash, e.type, e.content, e.created_at
			FROM telescope_entries e %s ORDER BY e.sequence DESC LIMIT ?`, where),
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []*telescope.EntryResult

	for rows.Next() {
		entry, err := r.scanEntry(rows)

		if err != nil {
			return nil, err
		}

		results = append(results, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, e := range results {
		tags, err := r.loadTags(e.ID)

		if err != nil {
			return nil, err
		}

		e.Tags = tags
	}

	return results, nil
}

func (r *DatabaseRepository) Store(entries []*telescope.IncomingEntry) error {
	for i := 0; i < len(entries); i += r.chunkSize {
		end := i + r.chunkSize

		if end > len(entries) {
			end = len(entries)
		}

		if err := r.storeChunk(entries[i:end]); err != nil {
			return err
		}
	}

	return nil
}

func (r *DatabaseRepository) storeChunk(entries []*telescope.IncomingEntry) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	entryStmt, err := tx.Prepare(
		`INSERT OR IGNORE INTO telescope_entries (uuid, batch_id, family_hash, type, content, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
	)

	if err != nil {
		return err
	}

	defer entryStmt.Close()

	tagStmt, err := tx.Prepare(
		`INSERT OR IGNORE INTO telescope_entries_tags (entry_uuid, tag) VALUES (?, ?)`,
	)

	if err != nil {
		return err
	}

	defer tagStmt.Close()

	for _, e := range entries {
		content, err := json.Marshal(e.Content)

		if err != nil {
			return fmt.Errorf("telescope: marshal content for %s: %w", e.UUID, err)
		}

		if _, err = entryStmt.Exec(
			e.UUID,
			e.BatchID,
			e.FamilyHash,
			e.Type,
			string(content),
			e.RecordedAt.UTC().Format(time.RFC3339Nano),
		); err != nil {
			return err
		}

		for _, tag := range e.Tags {
			if _, err = tagStmt.Exec(e.UUID, tag); err != nil {

				continue
			}
		}
	}

	return tx.Commit()
}

func (r *DatabaseRepository) Update(updates []*telescope.EntryUpdate) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, u := range updates {

		var rawContent string
		row := tx.QueryRow(`SELECT content FROM telescope_entries WHERE uuid = ?`, u.UUID)

		if err = row.Scan(&rawContent); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}

			return err
		}

		var content map[string]any

		if err = json.Unmarshal([]byte(rawContent), &content); err != nil {
			return err
		}

		for k, v := range u.Changes {
			content[k] = v
		}

		merged, err := json.Marshal(content)

		if err != nil {
			return err
		}

		if _, err = tx.Exec(`UPDATE telescope_entries SET content = ? WHERE uuid = ?`, string(merged), u.UUID); err != nil {
			return err
		}

		for _, tag := range u.Tags.Add {
			if _, err = tx.Exec(`INSERT OR IGNORE INTO telescope_entries_tags (entry_uuid, tag) VALUES (?, ?)`, u.UUID, tag); err != nil {
				continue
			}
		}

		for _, tag := range u.Tags.Remove {
			if _, err = tx.Exec(`DELETE FROM telescope_entries_tags WHERE entry_uuid = ? AND tag = ?`, u.UUID, tag); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *DatabaseRepository) LoadMonitoredTags() error {
	rows, err := r.db.Query(`SELECT tag FROM telescope_monitoring`)

	if err != nil {
		return err
	}

	defer rows.Close()

	r.mu.Lock()

	defer r.mu.Unlock()

	r.monitoredTags = make(map[string]struct{})

	for rows.Next() {
		var tag string

		if err := rows.Scan(&tag); err != nil {
			return err
		}

		r.monitoredTags[tag] = struct{}{}
	}

	return rows.Err()
}

func (r *DatabaseRepository) IsMonitoring(tags []string) bool {
	r.mu.RLock()

	defer r.mu.RUnlock()

	for _, t := range tags {
		if _, ok := r.monitoredTags[t]; ok {
			return true
		}
	}

	return false
}

func (r *DatabaseRepository) Monitoring() []string {
	r.mu.RLock()

	defer r.mu.RUnlock()

	out := make([]string, 0, len(r.monitoredTags))

	for t := range r.monitoredTags {
		out = append(out, t)
	}

	sort.Strings(out)

	return out
}

func (r *DatabaseRepository) Monitor(tags []string) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, tag := range tags {
		if _, err = tx.Exec(`INSERT OR IGNORE INTO telescope_monitoring (tag) VALUES (?)`, tag); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	r.mu.Lock()

	defer r.mu.Unlock()

	for _, t := range tags {
		r.monitoredTags[t] = struct{}{}
	}

	return nil
}

func (r *DatabaseRepository) StopMonitoring(tags []string) error {
	tx, err := r.db.Begin()

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, tag := range tags {
		if _, err = tx.Exec(`DELETE FROM telescope_monitoring WHERE tag = ?`, tag); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	r.mu.Lock()

	defer r.mu.Unlock()

	for _, t := range tags {
		delete(r.monitoredTags, t)
	}

	return nil
}

func (r *DatabaseRepository) Clear() error {
	_, err := r.db.Exec(`DELETE FROM telescope_entries`)

	return err
}

func (r *DatabaseRepository) Prune(before time.Time, keepExceptions bool) (int64, error) {
	var (
		query string
		args  []any
	)

	ts := before.UTC().Format(time.RFC3339Nano)

	if keepExceptions {
		query = `DELETE FROM telescope_entries WHERE created_at < ? AND type != ?`
		args = []any{ts, telescope.EntryTypeException}
	} else {
		query = `DELETE FROM telescope_entries WHERE created_at < ?`
		args = []any{ts}
	}

	result, err := r.db.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *DatabaseRepository) scanEntry(s scanner) (*telescope.EntryResult, error) {
	var (
		sequence   int64
		id         string
		batchID    string
		familyHash sql.NullString
		entryType  string
		rawContent string
		createdAt  string
	)

	if err := s.Scan(&sequence, &id, &batchID, &familyHash, &entryType, &rawContent, &createdAt); err != nil {
		return nil, err
	}

	var content map[string]any

	if err := json.Unmarshal([]byte(rawContent), &content); err != nil {
		return nil, fmt.Errorf("telescope: unmarshal content for %s: %w", id, err)
	}

	ts, _ := time.Parse(time.RFC3339Nano, createdAt)

	return &telescope.EntryResult{
		ID:         id,
		Sequence:   sequence,
		BatchID:    batchID,
		FamilyHash: familyHash.String,
		Type:       entryType,
		Content:    content,
		CreatedAt:  ts,
	}, nil
}

func (r *DatabaseRepository) loadTags(uuid string) ([]string, error) {
	rows, err := r.db.Query(
		`SELECT tag FROM telescope_entries_tags WHERE entry_uuid = ?`, uuid,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []string

	for rows.Next() {
		var tag string

		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, rows.Err()
}
