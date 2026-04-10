package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// DBConnection is the database connection interface required by DatabaseStore.
type DBConnection interface {
	// QueryRow executes a query and returns a single row result.
	QueryRow(ctx context.Context, query string, args ...any) DBRow
	// Exec executes a statement.
	Exec(ctx context.Context, query string, args ...any) error
}

// DBRow is a single database row.
type DBRow interface {
	Scan(dest ...any) error
}

// DatabaseStore stores cached values in a SQL table named "cache".
// The table must have columns: key TEXT PRIMARY KEY, value TEXT, expiration BIGINT.
type DatabaseStore struct {
	conn   DBConnection
	table  string
	prefix string
	clock  Clock
}

var _ Store = (*DatabaseStore)(nil)
var _ Locker = (*DatabaseStore)(nil)
var _ TaggableStore = (*DatabaseStore)(nil)

// NewDatabaseStore creates a DatabaseStore using the given connection.
// table is the cache table name (default "cache").
func NewDatabaseStore(conn DBConnection, table, prefix string) *DatabaseStore {
	if table == "" {
		table = "cache"
	}

	return &DatabaseStore{conn: conn, table: table, prefix: prefix}
}

func (s *DatabaseStore) now() time.Time {
	if s.clock != nil {
		return s.clock.Now()
	}

	return time.Now()
}

func (s *DatabaseStore) GetPrefix() string { return s.prefix }

// Tags returns a tag-scoped view of the store.
func (s *DatabaseStore) Tags(tags ...string) TaggedCache {
	return NewTaggedCache(s, NewTagSet(s, tags))
}

func (s *DatabaseStore) prefixed(key string) string {
	if s.prefix == "" {
		return key
	}

	return s.prefix + key
}

func (s *DatabaseStore) Get(ctx context.Context, key string) (any, error) {
	var encoded string

	var expiration int64

	row := s.conn.QueryRow(ctx,
		fmt.Sprintf("SELECT value, expiration FROM %s WHERE key = $1", s.table),
		s.prefixed(key),
	)

	if err := row.Scan(&encoded, &expiration); err != nil {
		return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
	}

	if expiration > 0 && s.now().Unix() > expiration {
		_ = s.Forget(ctx, key)

		return nil, fmt.Errorf("%w: %q", ErrNotFound, key)
	}

	var value any

	if err := json.Unmarshal([]byte(encoded), &value); err != nil {
		return nil, err
	}

	return value, nil
}

func (s *DatabaseStore) GetMany(ctx context.Context, keys []string) (map[string]any, error) {
	out := make(map[string]any, len(keys))

	for _, key := range keys {
		if v, err := s.Get(ctx, key); err == nil {
			out[key] = v
		}
	}

	return out, nil
}

func (s *DatabaseStore) Put(ctx context.Context, key string, value any, ttl time.Duration) error {
	encoded, err := json.Marshal(value)

	if err != nil {
		return err
	}

	var exp int64

	if ttl > 0 {
		exp = s.now().Add(ttl).Unix()
	}

	return s.conn.Exec(ctx,
		fmt.Sprintf(`INSERT INTO %s (key, value, expiration) VALUES ($1, $2, $3)
		             ON CONFLICT (key) DO UPDATE SET value = $2, expiration = $3`, s.table),
		s.prefixed(key), string(encoded), exp,
	)
}

func (s *DatabaseStore) PutMany(ctx context.Context, values map[string]any, ttl time.Duration) error {
	for k, v := range values {
		if err := s.Put(ctx, k, v, ttl); err != nil {
			return err
		}
	}

	return nil
}

func (s *DatabaseStore) Add(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	if _, err := s.Get(ctx, key); err == nil {
		return false, nil
	}

	return true, s.Put(ctx, key, value, ttl)
}

func (s *DatabaseStore) Forever(ctx context.Context, key string, value any) error {
	return s.Put(ctx, key, value, 0)
}

func (s *DatabaseStore) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	current, err := s.Get(ctx, key)

	var val int64

	if err == nil {
		val, err = toInt64(current)

		if err != nil {
			return 0, fmt.Errorf("%w: key %q", ErrInvalidValue, key)
		}
	}

	result := val + delta

	return result, s.Put(ctx, key, result, 0)
}

func (s *DatabaseStore) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return s.Increment(ctx, key, -delta)
}

func (s *DatabaseStore) Touch(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	v, err := s.Get(ctx, key)

	if err != nil {
		return false, nil
	}

	return true, s.Put(ctx, key, v, ttl)
}

func (s *DatabaseStore) Forget(ctx context.Context, key string) error {
	return s.conn.Exec(ctx,
		fmt.Sprintf("DELETE FROM %s WHERE key = $1", s.table),
		s.prefixed(key),
	)
}

// Lock returns a database-backed lock for the named resource.
func (s *DatabaseStore) Lock(name, owner string, ttl time.Duration) Lock {
	return NewDatabaseLock(s.conn, s.table+"_locks", name, owner, ttl, s.clock)
}

func (s *DatabaseStore) Flush(ctx context.Context) error {
	return s.conn.Exec(ctx, fmt.Sprintf("DELETE FROM %s", s.table))
}
