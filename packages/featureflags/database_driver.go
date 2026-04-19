package featureflags

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

// DBExecutor abstracts *sql.DB for testability. Any type that implements these
// three methods (including *sql.DB itself) satisfies the interface.
type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// DatabaseDriver persists feature-flag state in a SQL table.
//
// Expected schema (PostgreSQL):
//
//	CREATE TABLE features (
//	    id         BIGSERIAL PRIMARY KEY,
//	    name       VARCHAR(255) NOT NULL,
//	    scope      VARCHAR(255) NOT NULL,
//	    value      TEXT NOT NULL,       -- JSON-encoded
//	    created_at TIMESTAMP NOT NULL,
//	    updated_at TIMESTAMP NOT NULL,
//	    UNIQUE (name, scope)
//	);
//
// JSON note: json.Unmarshal into any decodes numbers as float64. This is a
// known limitation when storing integer feature values — callers should cast
// appropriately after retrieval.
type DatabaseDriver struct {
	mu         sync.Mutex
	db         DBExecutor
	table      string
	resolvers  map[string]func(ctx context.Context, scope any) (any, error)
	dispatcher EventDispatcher
	maxRetries int
}

var _ Driver = (*DatabaseDriver)(nil)
var _ StoredFeaturesLister = (*DatabaseDriver)(nil)
var _ BulkFeatureSetter = (*DatabaseDriver)(nil)

// NewDatabaseDriver creates a DatabaseDriver backed by db, storing rows in
// table. Event dispatch is disabled.
func NewDatabaseDriver(db DBExecutor, table string) *DatabaseDriver {
	return &DatabaseDriver{
		db:         db,
		table:      table,
		resolvers:  make(map[string]func(ctx context.Context, scope any) (any, error)),
		maxRetries: 3,
	}
}

// NewDatabaseDriverWithDispatcher creates a DatabaseDriver that dispatches
// featureflags events via d.
func NewDatabaseDriverWithDispatcher(db DBExecutor, table string, d EventDispatcher) *DatabaseDriver {
	drv := NewDatabaseDriver(db, table)
	drv.dispatcher = d

	return drv
}

// Define registers a resolver for a named feature. A second call for the same
// name replaces the previous resolver.
func (d *DatabaseDriver) Define(name string, resolver func(ctx context.Context, scope any) (any, error)) {
	d.mu.Lock()

	defer d.mu.Unlock()

	d.resolvers[name] = resolver
}

// Defined returns the names of all features with registered resolvers.
func (d *DatabaseDriver) Defined() []string {
	d.mu.Lock()

	defer d.mu.Unlock()

	names := make([]string, 0, len(d.resolvers))

	for name := range d.resolvers {
		names = append(names, name)
	}

	return names
}

// Get resolves a single feature for a single scope.
//
//  1. SerializeScope(scope) produces the storage key.
//  2. A SELECT is issued; if a row is found the JSON value is unmarshalled and
//     returned.
//  3. On sql.ErrNoRows the registered resolver is invoked and the result is
//     persisted via Set.
//  4. When no resolver exists, UnknownFeatureResolved is dispatched and
//     ErrFeatureNotDefined is returned.
func (d *DatabaseDriver) Get(ctx context.Context, feature string, scope any) (any, error) {
	key, err := SerializeScope(scope)

	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(
		"SELECT value FROM %s WHERE name=$1 AND scope=$2",
		d.table,
	)

	row := d.db.QueryRowContext(ctx, query, feature, key)

	var raw string

	err = row.Scan(&raw)

	if err == nil {
		var value any

		if unmarshalErr := json.Unmarshal([]byte(raw), &value); unmarshalErr != nil {
			return nil, unmarshalErr
		}

		return value, nil
	}

	if err != sql.ErrNoRows {
		return nil, err
	}

	// No stored value — look for a resolver.
	d.mu.Lock()
	resolver := d.resolvers[feature]
	d.mu.Unlock()

	if resolver == nil {
		d.dispatch(ctx, UnknownFeatureResolved{Feature: feature, Scope: scope})

		return nil, fmt.Errorf("%w: %q", ErrFeatureNotDefined, feature)
	}

	value, err := resolver(ctx, scope)

	if err != nil {
		return nil, err
	}

	if setErr := d.Set(ctx, feature, scope, value); setErr != nil {
		return nil, setErr
	}

	return value, nil
}

// GetAll resolves multiple features for multiple scopes. The returned map is
// parallel-indexed: result[feature][i] corresponds to features[feature][i].
func (d *DatabaseDriver) GetAll(ctx context.Context, features map[string][]any) (map[string][]any, error) {
	result := make(map[string][]any, len(features))

	for feature, scopes := range features {
		values := make([]any, len(scopes))

		for i, scope := range scopes {
			val, err := d.Get(ctx, feature, scope)

			if err != nil {
				return nil, err
			}

			values[i] = val
		}

		result[feature] = values
	}

	return result, nil
}

// Set stores a resolved value for the given feature and scope, bypassing the
// resolver. It uses an upsert statement and retries up to maxRetries times on
// unique-constraint violations.
func (d *DatabaseDriver) Set(ctx context.Context, feature string, scope any, value any) error {
	key, err := SerializeScope(scope)

	if err != nil {
		return err
	}

	encoded, err := json.Marshal(value)

	if err != nil {
		return err
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (name, scope, value, created_at, updated_at) VALUES ($1,$2,$3,$4,$4) "+
			"ON CONFLICT (name, scope) DO UPDATE SET value=$3, updated_at=$4",
		d.table,
	)

	now := time.Now().UTC()

	var lastErr error

	for attempt := 0; attempt <= d.maxRetries; attempt++ {
		_, lastErr = d.db.ExecContext(ctx, query, feature, key, string(encoded), now)

		if lastErr == nil {
			return nil
		}

		msg := lastErr.Error()
		isConflict := strings.Contains(msg, "UNIQUE") ||
			strings.Contains(msg, "unique") ||
			strings.Contains(msg, "duplicate")

		if !isConflict {
			return lastErr
		}
	}

	return fmt.Errorf("%w: %v", ErrStorageConflict, lastErr)
}

// SetAll stores multiple (feature, scope, value) entries by calling Set for
// each entry.
func (d *DatabaseDriver) SetAll(ctx context.Context, entries []FeatureEntry) error {
	for _, e := range entries {
		if err := d.Set(ctx, e.Feature, e.Scope, e.Value); err != nil {
			return err
		}
	}

	return nil
}

// SetForAllScopes updates the resolved value for every scope that already has
// stored state for the feature.
func (d *DatabaseDriver) SetForAllScopes(ctx context.Context, feature string, value any) error {
	encoded, err := json.Marshal(value)

	if err != nil {
		return err
	}

	query := fmt.Sprintf(
		"UPDATE %s SET value=$1, updated_at=$2 WHERE name=$3",
		d.table,
	)

	now := time.Now().UTC()

	_, err = d.db.ExecContext(ctx, query, string(encoded), now, feature)

	return err
}

// Delete removes the stored resolved value for the given feature and scope.
// The next Get call will re-invoke the resolver.
func (d *DatabaseDriver) Delete(ctx context.Context, feature string, scope any) error {
	key, err := SerializeScope(scope)

	if err != nil {
		return err
	}

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE name=$1 AND scope=$2",
		d.table,
	)

	_, err = d.db.ExecContext(ctx, query, feature, key)

	return err
}

// Purge removes stored state.
//   - nil: deletes every row in the table.
//   - empty non-nil slice: no-op.
//   - non-empty slice: deletes rows whose name is in the given list.
func (d *DatabaseDriver) Purge(ctx context.Context, features []string) error {
	if features == nil {
		query := fmt.Sprintf("DELETE FROM %s", d.table)

		_, err := d.db.ExecContext(ctx, query)

		return err
	}

	if len(features) == 0 {
		return nil
	}

	// Build dynamic placeholder list: $1, $2, ...
	placeholders := make([]string, len(features))
	args := make([]any, len(features))

	for i, name := range features {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = name
	}

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE name IN (%s)",
		d.table,
		strings.Join(placeholders, ","),
	)

	_, err := d.db.ExecContext(ctx, query, args...)

	return err
}

// Stored returns the distinct feature names that have at least one row stored
// in the database.
func (d *DatabaseDriver) Stored(ctx context.Context) ([]string, error) {
	query := fmt.Sprintf("SELECT DISTINCT name FROM %s", d.table)

	rows, err := d.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close() //nolint:errcheck

	var names []string

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			return nil, err
		}

		names = append(names, name)
	}

	return names, rows.Err()
}

func (d *DatabaseDriver) dispatch(ctx context.Context, event Event) {
	if d.dispatcher != nil {
		d.dispatcher.Dispatch(ctx, event)
	}
}
