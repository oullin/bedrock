package handlers

import (
	"context"
	"time"
)

// DBConn is the minimal database interface required by DatabaseHandler.
type DBConn interface {
	// QueryRow executes a query that returns at most one row.
	QueryRow(ctx context.Context, query string, args ...any) DBRow
	// Exec executes a query that doesn't return rows.
	Exec(ctx context.Context, query string, args ...any) error
}

// DBRow represents a single result row.
type DBRow interface {
	Scan(dest ...any) error
}

// DatabaseHandler stores sessions in a SQL table with the schema:
//
//	CREATE TABLE sessions (
//	  id         TEXT PRIMARY KEY,
//	  user_id    TEXT,
//	  ip_address TEXT,
//	  user_agent TEXT,
//	  payload    TEXT NOT NULL,
//	  last_activity INT NOT NULL
//	);
type DatabaseHandler struct {
	db        DBConn
	table     string
	exists    bool
	userID    string
	ipAddress string
	userAgent string
}

// NewDatabaseHandler creates a DatabaseHandler using the given connection.
// table is the sessions table name (e.g. "sessions").
func NewDatabaseHandler(db DBConn, table string) *DatabaseHandler {
	return &DatabaseHandler{db: db, table: table}
}

// SetExists implements session.ExistenceAware.
func (h *DatabaseHandler) SetExists(exists bool) { h.exists = exists }

// WithRequest enriches the handler with per-request metadata written on each save.
func (h *DatabaseHandler) WithRequest(userID, ipAddress, userAgent string) {
	h.userID = userID
	h.ipAddress = ipAddress
	h.userAgent = userAgent
}

func (h *DatabaseHandler) Open(_ context.Context, _, _ string) error { return nil }

func (h *DatabaseHandler) Close(_ context.Context) error { return nil }

func (h *DatabaseHandler) Read(ctx context.Context, id string) (string, error) {
	query := "SELECT payload FROM " + h.table + " WHERE id = $1"
	row := h.db.QueryRow(ctx, query, id)

	var payload string

	err := row.Scan(&payload)

	if err != nil {
		// No row found — return empty string (not an error for sessions).
		return "", nil
	}

	h.exists = true

	return payload, nil
}

func (h *DatabaseHandler) Write(ctx context.Context, id, data string) error {
	now := time.Now().Unix()

	if h.exists {
		return h.db.Exec(ctx,
			"UPDATE "+h.table+" SET user_id=$1, ip_address=$2, user_agent=$3, payload=$4, last_activity=$5 WHERE id=$6",
			h.userID, h.ipAddress, h.userAgent, data, now, id,
		)
	}

	return h.db.Exec(ctx,
		"INSERT INTO "+h.table+" (id, user_id, ip_address, user_agent, payload, last_activity) VALUES ($1,$2,$3,$4,$5,$6)",
		id, h.userID, h.ipAddress, h.userAgent, data, now,
	)
}

func (h *DatabaseHandler) Destroy(ctx context.Context, id string) error {
	return h.db.Exec(ctx, "DELETE FROM "+h.table+" WHERE id=$1", id)
}

func (h *DatabaseHandler) GC(ctx context.Context, maxLifetime int) error {
	cutoff := time.Now().Unix() - int64(maxLifetime)

	return h.db.Exec(ctx, "DELETE FROM "+h.table+" WHERE last_activity < $1", cutoff)
}
