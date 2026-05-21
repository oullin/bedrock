package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	dbcontract "github.com/bedrock/packages/contracts/database"
	cevents "github.com/bedrock/packages/contracts/events"
	dbevents "github.com/bedrock/packages/database/events"
)

// Connection wraps a standard *sql.DB and adds query logging, event
// dispatching, transaction management, and prefix support. It is the Go port
// Ref: @bedrock/code-0202
type Connection struct {
	mu             sync.RWMutex
	db             *sql.DB
	readDB         *sql.DB
	name           string
	database       string
	tablePrefix    string
	driverName     string
	config         map[string]any
	events         cevents.Dispatcher
	queryLog       []QueryLog
	loggingQueries bool
	pretending     bool

	// Transaction state.
	tx      *sql.Tx
	txLevel int
	txMgr   *TransactionManager

	// Callbacks.
	beforeExecuting []func(string, []any)
}

// QueryLog records a single query execution.
type QueryLog struct {
	Query    string
	Bindings []any
	Duration time.Duration
}

// NewConnection creates a new Connection wrapping the given *sql.DB.

// SetDriverName sets the driver name for this connection.

// SetEventDispatcher sets the event dispatcher.

// GetEventDispatcher returns the event dispatcher.

// DB returns the underlying *sql.DB.

// SetReadDB sets a separate read connection.

// GetReadDB returns the read *sql.DB, falling back to the write connection.

// Table begins a fluent query against a database table. The returned value is
// a *query.Builder — the any return satisfies the contract without creating a
// circular import. Callers assert the concrete type.

// Implemented by the query package via SetQueryBuilder on connection init.
// This stub satisfies the interface; driver connections override it.

// Raw creates a raw Expression.

// SelectOne runs a SELECT query and returns the first row.

// Select runs a SELECT query and returns all rows as maps.

// Insert executes an INSERT statement.

// Update executes an UPDATE statement and returns affected rows.

// Delete executes a DELETE statement and returns affected rows.

// Statement executes a raw SQL statement.

// AffectingStatement executes a statement and returns affected rows.

// Unprepared runs a raw query without parameter binding.

// PrepareBindings prepares bindings for execution.

// GetTablePrefix returns the table prefix.

// SetTablePrefix sets the table prefix.

// GetDatabaseName returns the database name.

// SetDatabaseName sets the database name.

// GetDriverName returns the driver name.

// GetName returns the connection name.

// GetConfig returns a configuration value.

// EnableQueryLog enables query logging.

// DisableQueryLog disables query logging.

// IsLogging reports whether query logging is enabled.

// GetQueryLog returns the query log.

// FlushQueryLog clears the query log.

// Pretend runs queries in pretend mode (logging only, no execution).

// IsPretending returns true if the connection is in pretend mode.

// BeforeExecuting registers a callback to run before each query.

// Disconnect closes the underlying database connection.

// statement executes a statement and returns success.

// affectingStatement executes a statement and returns affected rows.

// run executes a query, fires events, and logs the query.

// logQuery logs a query and dispatches the QueryExecuted event.

type queryConnection interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type execConnection interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

var _ dbcontract.Connection = (*Connection)(nil)

func NewConnection(db *sql.DB, name, database, tablePrefix string, config map[string]any) *Connection {
	return &Connection{
		db:          db,
		name:        name,
		database:    database,
		tablePrefix: tablePrefix,
		driverName:  "",
		config:      config,
		txMgr:       NewTransactionManager(),
	}
}

func (c *Connection) SetDriverName(driver string) {
	c.driverName = driver
}

func (c *Connection) SetEventDispatcher(events cevents.Dispatcher) {
	c.events = events
}

func (c *Connection) GetEventDispatcher() cevents.Dispatcher {
	return c.events
}

func (c *Connection) DB() *sql.DB {
	return c.db
}

func (c *Connection) SetReadDB(db *sql.DB) {
	c.readDB = db
}

func (c *Connection) GetReadDB() *sql.DB {
	if c.readDB != nil {
		return c.readDB
	}

	return c.db
}

func (c *Connection) Table(_ context.Context, table string, as ...string) any {

	_ = table
	_ = as

	return nil
}

func (c *Connection) Raw(value string) dbcontract.Expression {
	return NewExpr(value)
}

func (c *Connection) SelectOne(ctx context.Context, query string, bindings ...any) (map[string]any, error) {
	rows, err := c.Select(ctx, query, bindings...)

	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows[0], nil
}

func (c *Connection) Select(ctx context.Context, query string, bindings ...any) ([]map[string]any, error) {
	return c.run(ctx, query, bindings, func(ctx context.Context, q string, b []any) ([]map[string]any, error) {
		if c.pretending {
			return nil, nil
		}

		stmt := c.getQueryConnection()

		rows, err := stmt.QueryContext(ctx, q, b...)

		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrQueryFailed, err)
		}

		defer rows.Close()

		return scanRows(rows)
	})
}

func (c *Connection) Insert(ctx context.Context, query string, bindings ...any) (bool, error) {
	return c.statement(ctx, query, bindings)
}

func (c *Connection) Update(ctx context.Context, query string, bindings ...any) (int64, error) {
	return c.affectingStatement(ctx, query, bindings)
}

func (c *Connection) Delete(ctx context.Context, query string, bindings ...any) (int64, error) {
	return c.affectingStatement(ctx, query, bindings)
}

func (c *Connection) Statement(ctx context.Context, query string, bindings ...any) (bool, error) {
	return c.statement(ctx, query, bindings)
}

func (c *Connection) AffectingStatement(ctx context.Context, query string, bindings ...any) (int64, error) {
	return c.affectingStatement(ctx, query, bindings)
}

func (c *Connection) Unprepared(ctx context.Context, query string) (bool, error) {
	return c.statement(ctx, query, nil)
}

func (c *Connection) PrepareBindings(bindings []any) []any {
	prepared := make([]any, len(bindings))

	for i, b := range bindings {
		switch v := b.(type) {
		case bool:
			if v {
				prepared[i] = 1
			} else {
				prepared[i] = 0
			}
		default:
			prepared[i] = v
		}
	}

	return prepared
}

func (c *Connection) GetTablePrefix() string {
	return c.tablePrefix
}

func (c *Connection) SetTablePrefix(prefix string) {
	c.tablePrefix = prefix
}

func (c *Connection) GetDatabaseName() string {
	return c.database
}

func (c *Connection) SetDatabaseName(name string) {
	c.database = name
}

func (c *Connection) GetDriverName() string {
	return c.driverName
}

func (c *Connection) GetName() string {
	return c.name
}

func (c *Connection) GetConfig(key string) any {
	if c.config == nil {
		return nil
	}

	return c.config[key]
}

func (c *Connection) EnableQueryLog() {
	c.loggingQueries = true
}

func (c *Connection) DisableQueryLog() {
	c.loggingQueries = false
}

func (c *Connection) IsLogging() bool {
	return c.loggingQueries
}

func (c *Connection) GetQueryLog() []QueryLog {
	c.mu.RLock()

	defer c.mu.RUnlock()

	log := make([]QueryLog, len(c.queryLog))
	copy(log, c.queryLog)

	return log
}

func (c *Connection) FlushQueryLog() {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.queryLog = c.queryLog[:0]
}

func (c *Connection) Pretend(fn func()) []QueryLog {
	c.pretending = true
	c.EnableQueryLog()

	initialLen := len(c.queryLog)

	fn()

	c.pretending = false

	return c.queryLog[initialLen:]
}

func (c *Connection) IsPretending() bool {
	return c.pretending
}

func (c *Connection) BeforeExecuting(fn func(string, []any)) {
	c.beforeExecuting = append(c.beforeExecuting, fn)
}

func (c *Connection) Disconnect() error {
	if c.db != nil {
		return c.db.Close()
	}

	return nil
}

func (c *Connection) statement(ctx context.Context, query string, bindings []any) (bool, error) {
	_, err := c.run(ctx, query, bindings, func(ctx context.Context, q string, b []any) ([]map[string]any, error) {
		if c.pretending {
			return nil, nil
		}

		conn := c.getExecConnection()

		_, err := conn.ExecContext(ctx, q, b...)

		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrQueryFailed, err)
		}

		return nil, nil
	})

	return err == nil, err
}

func (c *Connection) affectingStatement(ctx context.Context, query string, bindings []any) (int64, error) {
	var affected int64

	_, err := c.run(ctx, query, bindings, func(ctx context.Context, q string, b []any) ([]map[string]any, error) {
		if c.pretending {
			return nil, nil
		}

		conn := c.getExecConnection()

		result, err := conn.ExecContext(ctx, q, b...)

		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrQueryFailed, err)
		}

		affected, _ = result.RowsAffected()

		return nil, nil
	})

	return affected, err
}

func (c *Connection) run(ctx context.Context, query string, bindings []any, callback func(context.Context, string, []any) ([]map[string]any, error)) ([]map[string]any, error) {
	for _, fn := range c.beforeExecuting {
		fn(query, bindings)
	}

	prepared := c.PrepareBindings(bindings)

	start := time.Now()
	result, err := callback(ctx, query, prepared)
	duration := time.Since(start)

	c.logQuery(query, prepared, duration)

	return result, err
}

func (c *Connection) logQuery(query string, bindings []any, duration time.Duration) {
	if c.loggingQueries {
		c.mu.Lock()
		c.queryLog = append(c.queryLog, QueryLog{
			Query:    query,
			Bindings: bindings,
			Duration: duration,
		})
		c.mu.Unlock()
	}

	if c.events != nil {
		c.events.Dispatch(context.Background(), dbevents.QueryExecuted{
			SQL:            query,
			Bindings:       bindings,
			Duration:       duration,
			ConnectionName: c.name,
		})
	}
}

// getQueryConnection returns the active query executor (tx or read db).
func (c *Connection) getQueryConnection() queryConnection {
	if c.tx != nil {
		return c.tx
	}

	return c.GetReadDB()
}

// getExecConnection returns the active exec executor (tx or write db).
func (c *Connection) getExecConnection() execConnection {
	if c.tx != nil {
		return c.tx
	}

	return c.db
}

// scanRows converts *sql.Rows into a slice of maps.
func scanRows(rows *sql.Rows) ([]map[string]any, error) {
	columns, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	var results []map[string]any

	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))

		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		row := make(map[string]any, len(columns))

		for i, col := range columns {
			val := values[i]

			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}

		results = append(results, row)
	}

	return results, rows.Err()
}
