package bus

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// DBExecutor abstracts database/sql for testability.
type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// DBTransactor extends DBExecutor with transaction support.
type DBTransactor interface {
	DBExecutor
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// DatabaseBatchRepository persists batch state in a SQL database.
//
// Expected table schema:
//
//	CREATE TABLE job_batches (
//	    id VARCHAR(255) PRIMARY KEY,
//	    name VARCHAR(255) NOT NULL,
//	    total_jobs INT NOT NULL DEFAULT 0,
//	    pending_jobs INT NOT NULL DEFAULT 0,
//	    failed_jobs INT NOT NULL DEFAULT 0,
//	    failed_job_ids TEXT NOT NULL DEFAULT '[]',
//	    options TEXT NOT NULL DEFAULT '{}',
//	    created_at TIMESTAMP NOT NULL,
//	    cancelled_at TIMESTAMP NULL,
//	    finished_at TIMESTAMP NULL
//	);
type DatabaseBatchRepository struct {
	db    DBExecutor
	table string
}

// NewDatabaseBatchRepository creates a DatabaseBatchRepository.
func NewDatabaseBatchRepository(db DBExecutor, table string) *DatabaseBatchRepository {
	if table == "" {
		table = "job_batches"
	}

	return &DatabaseBatchRepository{db: db, table: table}
}

// Get retrieves a batch by ID.
func (r *DatabaseBatchRepository) Get(ctx context.Context, id string) (*Batch, error) {
	query := fmt.Sprintf(
		"SELECT id, name, total_jobs, pending_jobs, failed_jobs, failed_job_ids, options, created_at, cancelled_at, finished_at FROM %s WHERE id = ?",
		r.table,
	)

	row := r.db.QueryRowContext(ctx, query, id)

	return r.scanBatch(row)
}

// GetList retrieves batches with pagination.
func (r *DatabaseBatchRepository) GetList(ctx context.Context, limit int, before string) ([]*Batch, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if before != "" {
		query := fmt.Sprintf(
			"SELECT id, name, total_jobs, pending_jobs, failed_jobs, failed_job_ids, options, created_at, cancelled_at, finished_at FROM %s WHERE id < ? ORDER BY id DESC LIMIT ?",
			r.table,
		)

		rows, err = r.db.QueryContext(ctx, query, before, limit)
	} else {
		query := fmt.Sprintf(
			"SELECT id, name, total_jobs, pending_jobs, failed_jobs, failed_job_ids, options, created_at, cancelled_at, finished_at FROM %s ORDER BY id DESC LIMIT ?",
			r.table,
		)

		rows, err = r.db.QueryContext(ctx, query, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("bus: get batches: %w", err)
	}

	defer rows.Close()

	var batches []*Batch

	for rows.Next() {
		b, err := r.scanBatchFromRows(rows)

		if err != nil {
			return batches, err
		}

		batches = append(batches, b)
	}

	return batches, rows.Err()
}

// Store persists a new batch.
func (r *DatabaseBatchRepository) Store(ctx context.Context, batch *Batch) error {
	failedIDs, err := json.Marshal(batch.FailedJobIDs)

	if err != nil {
		return fmt.Errorf("bus: marshal failed_job_ids: %w", err)
	}

	opts, err := json.Marshal(batch.Options)

	if err != nil {
		return fmt.Errorf("bus: marshal options: %w", err)
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (id, name, total_jobs, pending_jobs, failed_jobs, failed_job_ids, options, created_at, cancelled_at, finished_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		r.table,
	)

	_, err = r.db.ExecContext(ctx, query,
		batch.ID, batch.Name, batch.TotalJobs, batch.PendingJobs, batch.FailedJobs,
		string(failedIDs), string(opts), batch.CreatedAt, batch.CancelledAt, batch.FinishedAt,
	)

	if err != nil {
		return fmt.Errorf("bus: store batch: %w", err)
	}

	return nil
}

// IncrementTotalJobs increments the total job count.
func (r *DatabaseBatchRepository) IncrementTotalJobs(ctx context.Context, id string, amount int) error {
	query := fmt.Sprintf(
		"UPDATE %s SET total_jobs = total_jobs + ?, pending_jobs = pending_jobs + ? WHERE id = ?",
		r.table,
	)

	_, err := r.db.ExecContext(ctx, query, amount, amount, id)

	if err != nil {
		return fmt.Errorf("bus: increment total jobs: %w", err)
	}

	return nil
}

// DecrementPendingJobs decrements the pending job count and returns updated counts.
func (r *DatabaseBatchRepository) DecrementPendingJobs(ctx context.Context, id string) (*UpdatedBatchJobCounts, error) {
	query := fmt.Sprintf(
		"UPDATE %s SET pending_jobs = pending_jobs - 1 WHERE id = ?",
		r.table,
	)

	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return nil, fmt.Errorf("bus: decrement pending jobs: %w", err)
	}

	return r.fetchCounts(ctx, id)
}

// IncrementFailedJobs increments the failed job count and appends the job ID.
func (r *DatabaseBatchRepository) IncrementFailedJobs(ctx context.Context, id string, failedJobID string) (*UpdatedBatchJobCounts, error) {
	// Fetch current failed_job_ids, append, and update.
	var failedIDsJSON string

	query := fmt.Sprintf("SELECT failed_job_ids FROM %s WHERE id = ?", r.table)

	row := r.db.QueryRowContext(ctx, query, id)

	if err := row.Scan(&failedIDsJSON); err != nil {
		return nil, fmt.Errorf("bus: read failed_job_ids: %w", err)
	}

	var failedIDs []string

	if err := json.Unmarshal([]byte(failedIDsJSON), &failedIDs); err != nil {
		failedIDs = []string{}
	}

	failedIDs = append(failedIDs, failedJobID)

	updatedJSON, err := json.Marshal(failedIDs)

	if err != nil {
		return nil, fmt.Errorf("bus: marshal failed_job_ids: %w", err)
	}

	updateQuery := fmt.Sprintf(
		"UPDATE %s SET failed_jobs = failed_jobs + 1, pending_jobs = pending_jobs - 1, failed_job_ids = ? WHERE id = ?",
		r.table,
	)

	if _, err = r.db.ExecContext(ctx, updateQuery, string(updatedJSON), id); err != nil {
		return nil, fmt.Errorf("bus: increment failed jobs: %w", err)
	}

	return r.fetchCounts(ctx, id)
}

// MarkAsFinished sets the finished_at timestamp.
func (r *DatabaseBatchRepository) MarkAsFinished(ctx context.Context, id string) error {
	query := fmt.Sprintf("UPDATE %s SET finished_at = ? WHERE id = ?", r.table)

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)

	if err != nil {
		return fmt.Errorf("bus: mark as finished: %w", err)
	}

	return nil
}

// Cancel sets the cancelled_at timestamp.
func (r *DatabaseBatchRepository) Cancel(ctx context.Context, id string) error {
	query := fmt.Sprintf("UPDATE %s SET cancelled_at = ? WHERE id = ?", r.table)

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)

	if err != nil {
		return fmt.Errorf("bus: cancel batch: %w", err)
	}

	return nil
}

// Delete removes a batch from the database.
func (r *DatabaseBatchRepository) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", r.table)

	_, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("bus: delete batch: %w", err)
	}

	return nil
}

// Prune removes finished batches older than the given time, in chunks of 1000.
func (r *DatabaseBatchRepository) Prune(ctx context.Context, before time.Time) (int, error) {
	return r.pruneWhere(ctx, "finished_at IS NOT NULL AND finished_at < ?", before)
}

// PruneCancelled removes cancelled batches older than the given time.
func (r *DatabaseBatchRepository) PruneCancelled(ctx context.Context, before time.Time) (int, error) {
	return r.pruneWhere(ctx, "cancelled_at IS NOT NULL AND cancelled_at < ?", before)
}

// PruneUnfinished removes unfinished batches older than the given time.
func (r *DatabaseBatchRepository) PruneUnfinished(ctx context.Context, before time.Time) (int, error) {
	return r.pruneWhere(ctx, "finished_at IS NULL AND created_at < ?", before)
}

// Transaction executes fn within a database transaction.
func (r *DatabaseBatchRepository) Transaction(ctx context.Context, fn func(BatchRepository) error) error {
	txDB, ok := r.db.(DBTransactor)

	if !ok {
		return fn(r)
	}

	tx, err := txDB.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("bus: begin transaction: %w", err)
	}

	txRepo := &DatabaseBatchRepository{db: tx, table: r.table}

	if err = fn(txRepo); err != nil {
		_ = tx.Rollback()

		return err
	}

	return tx.Commit()
}

// RollBack is a no-op for DatabaseBatchRepository as rollback
// is handled internally by the Transaction method.
func (r *DatabaseBatchRepository) RollBack(_ context.Context) error {
	return nil
}

func (r *DatabaseBatchRepository) fetchCounts(ctx context.Context, id string) (*UpdatedBatchJobCounts, error) {
	query := fmt.Sprintf("SELECT pending_jobs, failed_jobs FROM %s WHERE id = ?", r.table)

	var counts UpdatedBatchJobCounts

	if err := r.db.QueryRowContext(ctx, query, id).Scan(&counts.PendingJobs, &counts.FailedJobs); err != nil {
		return nil, fmt.Errorf("bus: fetch counts: %w", err)
	}

	return &counts, nil
}

func (r *DatabaseBatchRepository) pruneWhere(ctx context.Context, where string, before time.Time) (int, error) {
	total := 0

	for {
		query := fmt.Sprintf("DELETE FROM %s WHERE %s LIMIT 1000", r.table, where)

		result, err := r.db.ExecContext(ctx, query, before)

		if err != nil {
			return total, fmt.Errorf("bus: prune: %w", err)
		}

		affected, err := result.RowsAffected()

		if err != nil {
			return total, fmt.Errorf("bus: prune rows affected: %w", err)
		}

		total += int(affected)

		if affected == 0 {
			break
		}
	}

	return total, nil
}

func (r *DatabaseBatchRepository) scanBatch(row *sql.Row) (*Batch, error) {
	var (
		b             Batch
		failedIDsJSON string
		optionsJSON   string
		cancelledAt   sql.NullTime
		finishedAt    sql.NullTime
	)

	err := row.Scan(
		&b.ID, &b.Name, &b.TotalJobs, &b.PendingJobs, &b.FailedJobs,
		&failedIDsJSON, &optionsJSON, &b.CreatedAt, &cancelledAt, &finishedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("bus: scan batch: %w", err)
	}

	if err = json.Unmarshal([]byte(failedIDsJSON), &b.FailedJobIDs); err != nil {
		b.FailedJobIDs = []string{}
	}

	if err = json.Unmarshal([]byte(optionsJSON), &b.Options); err != nil {
		b.Options = make(map[string]any)
	}

	if cancelledAt.Valid {
		b.CancelledAt = &cancelledAt.Time
	}

	if finishedAt.Valid {
		b.FinishedAt = &finishedAt.Time
	}

	return &b, nil
}

func (r *DatabaseBatchRepository) scanBatchFromRows(rows *sql.Rows) (*Batch, error) {
	var (
		b             Batch
		failedIDsJSON string
		optionsJSON   string
		cancelledAt   sql.NullTime
		finishedAt    sql.NullTime
	)

	err := rows.Scan(
		&b.ID, &b.Name, &b.TotalJobs, &b.PendingJobs, &b.FailedJobs,
		&failedIDsJSON, &optionsJSON, &b.CreatedAt, &cancelledAt, &finishedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("bus: scan batch row: %w", err)
	}

	if err = json.Unmarshal([]byte(failedIDsJSON), &b.FailedJobIDs); err != nil {
		b.FailedJobIDs = []string{}
	}

	if err = json.Unmarshal([]byte(optionsJSON), &b.Options); err != nil {
		b.Options = make(map[string]any)
	}

	if cancelledAt.Valid {
		b.CancelledAt = &cancelledAt.Time
	}

	if finishedAt.Valid {
		b.FinishedAt = &finishedAt.Time
	}

	return &b, nil
}
