package drivers

import (
	"context"
	"encoding/json"
	"fmt"
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
	Exec(ctx context.Context, query string, args ...any) error
}

// DBRow is a single query result row.
type DBRow interface {
	Scan(dest ...any) error
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

func (d *DatabaseDriver) count(ctx context.Context, query string, args ...any) (int64, error) {
	row := d.db.QueryRow(ctx, query, args...)

	var n int64

	if err := row.Scan(&n); err != nil {
		return 0, err
	}

	return n, nil
}
