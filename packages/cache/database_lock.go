package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/bedrock/packages/contracts"
)

// DatabaseLock is a distributed lock backed by a SQL database. It stores lock
// state in the cache table using the same DBConnection as DatabaseStore.
type DatabaseLock struct {
	conn    DBConnection
	table   string
	name    string
	owner   string
	ttl     time.Duration
	clock   contracts.Clock
	sleepMs int
}

var _ Lock = (*DatabaseLock)(nil)

// NewDatabaseLock creates a database-backed lock.
func NewDatabaseLock(conn DBConnection, table, name, owner string, ttl time.Duration, clock contracts.Clock) *DatabaseLock {
	if table == "" {
		table = "cache_locks"
	}

	return &DatabaseLock{conn: conn, table: table, name: name, owner: owner, ttl: ttl, clock: clock, sleepMs: 50}
}

func (l *DatabaseLock) now() time.Time {
	if l.clock != nil {
		return l.clock.Now()
	}

	return time.Now()
}

func (l *DatabaseLock) Acquire(ctx context.Context) (bool, error) {
	now := l.now()
	exp := now.Add(l.ttl).Unix()

	// Try to insert or update if expired.
	err := l.conn.Exec(ctx,
		fmt.Sprintf(`INSERT INTO %s (key, owner, expiration) VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE SET owner = $2, expiration = $3
		WHERE %s.expiration < $4 OR %s.owner = $2`, l.table, l.table, l.table),
		l.name, l.owner, exp, now.Unix(),
	)

	if err != nil {
		return false, err
	}

	// Verify we own the lock.
	var currentOwner string

	row := l.conn.QueryRow(ctx,
		fmt.Sprintf("SELECT owner FROM %s WHERE key = $1", l.table),
		l.name,
	)

	if err := row.Scan(&currentOwner); err != nil {
		return false, err
	}

	return currentOwner == l.owner, nil
}

func (l *DatabaseLock) Release(ctx context.Context) (bool, error) {
	var currentOwner string

	row := l.conn.QueryRow(ctx,
		fmt.Sprintf("SELECT owner FROM %s WHERE key = $1", l.table),
		l.name,
	)

	if err := row.Scan(&currentOwner); err != nil {
		return false, nil
	}

	if currentOwner != l.owner {
		return false, nil
	}

	err := l.conn.Exec(ctx,
		fmt.Sprintf("DELETE FROM %s WHERE key = $1 AND owner = $2", l.table),
		l.name, l.owner,
	)

	return err == nil, err
}

func (l *DatabaseLock) ForceRelease(ctx context.Context) error {
	return l.conn.Exec(ctx,
		fmt.Sprintf("DELETE FROM %s WHERE key = $1", l.table),
		l.name,
	)
}

func (l *DatabaseLock) Get(ctx context.Context, fn func() error) error {
	ok, err := l.Acquire(ctx)

	if err != nil {
		return err
	}

	if !ok {
		return ErrLockTimeout
	}

	defer l.Release(ctx) //nolint:errcheck

	return fn()
}

func (l *DatabaseLock) Block(ctx context.Context, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for {
		ok, err := l.Acquire(ctx)

		if err != nil {
			return err
		}

		if ok {
			return nil
		}

		if time.Now().After(deadline) {
			return ErrLockTimeout
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(l.sleepMs) * time.Millisecond):
		}
	}
}

func (l *DatabaseLock) Owner() string { return l.owner }

func (l *DatabaseLock) IsOwnedByCurrentProcess(ctx context.Context) (bool, error) {
	return l.IsOwnedBy(ctx, l.owner)
}

func (l *DatabaseLock) IsOwnedBy(ctx context.Context, owner string) (bool, error) {
	var currentOwner string

	row := l.conn.QueryRow(ctx,
		fmt.Sprintf("SELECT owner FROM %s WHERE key = $1 AND expiration > $2", l.table),
		l.name, l.now().Unix(),
	)

	if err := row.Scan(&currentOwner); err != nil {
		return false, nil
	}

	return currentOwner == owner, nil
}

func (l *DatabaseLock) BetweenBlockedAttemptsSleepFor(ms int) Lock {
	l.sleepMs = ms

	return l
}

func (l *DatabaseLock) Blocked(ctx context.Context) (bool, error) {
	var currentOwner string

	row := l.conn.QueryRow(ctx,
		fmt.Sprintf("SELECT owner FROM %s WHERE key = $1 AND expiration > $2", l.table),
		l.name, l.now().Unix(),
	)

	if err := row.Scan(&currentOwner); err != nil {
		return false, nil
	}

	return currentOwner != l.owner, nil
}
