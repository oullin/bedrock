package drivers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bedrock/packages/queue"
)

// DBExecer is the minimal database interface for the database queue driver.
//
// Schema:
//
//	CREATE TABLE jobs (
//	  id           BIGINT PRIMARY KEY AUTO_INCREMENT,
//	  queue        TEXT NOT NULL,
//	  payload      LONGTEXT NOT NULL,
//	  attempts     INT NOT NULL DEFAULT 0,
//	  reserved_at  BIGINT,
//	  available_at BIGINT NOT NULL,
//	  created_at   BIGINT NOT NULL
//	);
//
//	CREATE TABLE failed_jobs (
//	  id           BIGINT PRIMARY KEY AUTO_INCREMENT,
//	  uuid         TEXT NOT NULL UNIQUE,
//	  connection   TEXT NOT NULL,
//	  queue        TEXT NOT NULL,
//	  payload      LONGTEXT NOT NULL,
//	  exception    LONGTEXT NOT NULL,
//	  failed_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
//	);
type DBExecer interface {
	QueryRow(ctx context.Context, query string, args ...any) DBRow
	Query(ctx context.Context, query string, args ...any) (DBRows, error)
	Exec(ctx context.Context, query string, args ...any) error
}

// DBRow is a single query result row.
type DBRow interface {
	Scan(dest ...any) error
}

// DBRows iterates the result of a multi-row query. It mirrors the
// shape of database/sql.Rows at the interface level: call Next to
// advance, Scan to read into destinations, Close when done. Err
// reports any error that terminated iteration. Callers must call
// Close on the returned DBRows.
type DBRows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

// InspectedJob is the decoded, queue-centric view of a row in the jobs
// table. It is the Go port of Laravel's Illuminate\Queue\Jobs\InspectedJob
// and is returned by PendingJobs, DelayedJobs, and ReservedJobs.
type InspectedJob struct {
	ID         int64
	Queue      string
	Name       string
	UUID       string
	Attempts   int
	CreatedAt  time.Time
	ReservedAt *time.Time
}

// DatabaseDriver stores jobs in a SQL table.
type DatabaseDriver struct {
	db         DBExecer
	table      string
	connection string
}

// NewDatabaseDriver creates a DatabaseDriver.

// Reserve the job.

type dbJob struct{ BaseJob }

func NewDatabaseDriver(db DBExecer, table, connection string) *DatabaseDriver {
	if table == "" {
		table = "jobs"
	}

	return &DatabaseDriver{db: db, table: table, connection: connection}
}

func (d *DatabaseDriver) Push(ctx context.Context, queueName string, payload []byte) (string, error) {
	now := time.Now().Unix()
	err := d.db.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (queue, payload, attempts, reserved_at, available_at, created_at) VALUES ($1,$2,0,NULL,$3,$4)", d.table),
		queueName, string(payload), now, now,
	)

	return "", err
}

func (d *DatabaseDriver) PushDelayed(ctx context.Context, queueName string, payload []byte, delay time.Duration) (string, error) {
	now := time.Now()
	availAt := now.Add(delay).Unix()
	err := d.db.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (queue, payload, attempts, reserved_at, available_at, created_at) VALUES ($1,$2,0,NULL,$3,$4)", d.table),
		queueName, string(payload), availAt, now.Unix(),
	)

	return "", err
}

func (d *DatabaseDriver) PushMultiple(ctx context.Context, queueName string, payloads [][]byte) ([]string, error) {
	ids := make([]string, 0, len(payloads))

	for _, p := range payloads {
		id, err := d.Push(ctx, queueName, p)

		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (d *DatabaseDriver) Pop(ctx context.Context, queueName string) (queue.Job, error) {
	now := time.Now().Unix()
	row := d.db.QueryRow(ctx,
		fmt.Sprintf("SELECT id, payload, attempts FROM %s WHERE queue=$1 AND reserved_at IS NULL AND available_at<=$2 ORDER BY id ASC LIMIT 1", d.table),
		queueName, now,
	)

	var id int64

	var payload string

	var attempts int

	if err := row.Scan(&id, &payload, &attempts); err != nil {
		return nil, queue.ErrNoJob
	}

	_ = d.db.Exec(ctx,
		fmt.Sprintf("UPDATE %s SET reserved_at=$1, attempts=attempts+1 WHERE id=$2", d.table),
		now, id,
	)

	job := &dbJob{
		BaseJob: BaseJob{
			id:       fmt.Sprintf("%d", id),
			payload:  []byte(payload),
			queue:    queueName,
			attempts: attempts + 1,
		},
	}
	job.releaseFunc = func(delay time.Duration) error {
		availAt := time.Now().Add(delay).Unix()

		return d.db.Exec(ctx,
			fmt.Sprintf("UPDATE %s SET reserved_at=NULL, available_at=$1 WHERE id=$2", d.table),
			availAt, id,
		)
	}

	job.deleteFunc = func() error {
		return d.db.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE id=$1", d.table), id)
	}

	job.failFunc = func(err error) error {
		var errMsg string

		if err != nil {
			errMsg = err.Error()
		}

		errBytes, _ := json.Marshal(map[string]string{"exception": errMsg})
		_ = d.db.Exec(ctx,
			"INSERT INTO failed_jobs (uuid, connection, queue, payload, exception) VALUES ($1,$2,$3,$4,$5)",
			job.uuid, d.connection, queueName, payload, string(errBytes),
		)

		return job.deleteFunc()
	}

	return job, nil
}

func (d *DatabaseDriver) Size(ctx context.Context, queueName string) (int64, error) {
	return d.count(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE queue=$1 AND reserved_at IS NULL AND available_at<=$2", d.table), queueName, time.Now().Unix())
}

func (d *DatabaseDriver) PendingSize(ctx context.Context, queueName string) (int64, error) {
	return d.count(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE queue=$1 AND reserved_at IS NULL", d.table), queueName)
}

func (d *DatabaseDriver) DelayedSize(ctx context.Context, queueName string) (int64, error) {
	return d.count(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE queue=$1 AND reserved_at IS NULL AND available_at>$2", d.table), queueName, time.Now().Unix())
}

func (d *DatabaseDriver) ReservedSize(ctx context.Context, queueName string) (int64, error) {
	return d.count(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE queue=$1 AND reserved_at IS NOT NULL", d.table), queueName)
}

func (d *DatabaseDriver) ConnectionName() string { return d.connection }

// ClearQueue deletes every row for queueName from the jobs table. Go
// port of Laravel's DatabaseQueue::clear — a bulk DELETE WHERE queue=?
// matching the observable side effect of calling $queue->clear($name).
// Laravel returns the deleted row count; the Go driver surfaces only
// the error (row count would require a wider DBExecer interface).
func (d *DatabaseDriver) ClearQueue(ctx context.Context, queueName string) error {
	return d.db.Exec(ctx,
		fmt.Sprintf("DELETE FROM %s WHERE queue=$1", d.table),
		queueName,
	)
}

func (d *DatabaseDriver) count(ctx context.Context, query string, args ...any) (int64, error) {
	row := d.db.QueryRow(ctx, query, args...)

	var n int64

	if err := row.Scan(&n); err != nil {
		return 0, err
	}

	return n, nil
}

// Bulk inserts every payload in a single multi-row INSERT statement.
// It is the Go port of Laravel's DatabaseQueue::bulk and matches the
// observable side effects of a $db->insert([$record1, $record2, ...])
// call: one Exec, one INSERT, one round-trip. Returns the number of
// rows attempted.
//
// Prefer Bulk over PushMultiple when enqueuing a large batch: Bulk
// executes a single statement, whereas PushMultiple loops Push and
// incurs one round-trip per payload.
func (d *DatabaseDriver) Bulk(ctx context.Context, queueName string, payloads [][]byte) error {
	if len(payloads) == 0 {
		return nil
	}

	now := time.Now().Unix()

	var sb strings.Builder

	fmt.Fprintf(&sb, "INSERT INTO %s (queue, payload, attempts, reserved_at, available_at, created_at) VALUES ", d.table)

	args := make([]any, 0, 4*len(payloads))

	for i, p := range payloads {
		if i > 0 {
			sb.WriteString(",")
		}

		base := i * 4

		fmt.Fprintf(&sb, "($%d,$%d,0,NULL,$%d,$%d)", base+1, base+2, base+3, base+4)

		args = append(args, queueName, string(p), now, now)
	}

	return d.db.Exec(ctx, sb.String(), args...)
}

// PendingJobs returns the pending (unreserved, ready-to-run) rows for
// queueName. Go port of Laravel's DatabaseQueue::pendingJobs — selects
// rows where reserved_at IS NULL AND available_at <= now.
func (d *DatabaseDriver) PendingJobs(ctx context.Context, queueName string) ([]InspectedJob, error) {
	return d.fetchInspected(ctx,
		fmt.Sprintf("SELECT id, queue, payload, attempts, reserved_at FROM %s WHERE queue=$1 AND reserved_at IS NULL AND available_at<=$2", d.table),
		queueName, time.Now().Unix(),
	)
}

// DelayedJobs returns the delayed (unreserved, not-yet-available) rows
// for queueName. Go port of DatabaseQueue::delayedJobs.
func (d *DatabaseDriver) DelayedJobs(ctx context.Context, queueName string) ([]InspectedJob, error) {
	return d.fetchInspected(ctx,
		fmt.Sprintf("SELECT id, queue, payload, attempts, reserved_at FROM %s WHERE queue=$1 AND reserved_at IS NULL AND available_at>$2", d.table),
		queueName, time.Now().Unix(),
	)
}

// ReservedJobs returns the currently-reserved (in-flight) rows for
// queueName. Go port of DatabaseQueue::reservedJobs.
func (d *DatabaseDriver) ReservedJobs(ctx context.Context, queueName string) ([]InspectedJob, error) {
	return d.fetchInspected(ctx,
		fmt.Sprintf("SELECT id, queue, payload, attempts, reserved_at FROM %s WHERE queue=$1 AND reserved_at IS NOT NULL", d.table),
		queueName,
	)
}

// fetchInspected runs the given query and decodes each row into an
// InspectedJob. The payload JSON is parsed to extract displayName,
// uuid, and createdAt; any decoding error is recorded on the job's
// Name/UUID with empty strings so the caller can still page through.
func (d *DatabaseDriver) fetchInspected(ctx context.Context, query string, args ...any) ([]InspectedJob, error) {
	rows, err := d.db.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var out []InspectedJob

	for rows.Next() {
		var (
			id         int64
			queueName  string
			payload    string
			attempts   int
			reservedAt *int64
		)

		if err := rows.Scan(&id, &queueName, &payload, &attempts, &reservedAt); err != nil {
			return nil, err
		}

		job := InspectedJob{
			ID:       id,
			Queue:    queueName,
			Attempts: attempts,
		}

		if reservedAt != nil {
			t := time.Unix(*reservedAt, 0)
			job.ReservedAt = &t
		}

		var decoded map[string]any

		if err := json.Unmarshal([]byte(payload), &decoded); err == nil {
			if v, ok := decoded["displayName"].(string); ok {
				job.Name = v
			}

			if v, ok := decoded["uuid"].(string); ok {
				job.UUID = v
			}

			if v, ok := decoded["createdAt"].(float64); ok {
				job.CreatedAt = time.Unix(int64(v), 0)
			}
		}

		out = append(out, job)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
