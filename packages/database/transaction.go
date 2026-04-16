package database

import (
	"context"
	"database/sql"
	"fmt"

	dbcontract "github.com/bedrock/packages/contracts/database"
	dbevents "github.com/bedrock/packages/database/events"
)

// Transaction runs a closure inside a database transaction. If the closure
// returns an error the transaction is rolled back; otherwise it is committed.
// The attempts parameter controls automatic retry on deadlock (default 1).
func (c *Connection) Transaction(ctx context.Context, fn func(dbcontract.Connection) error, attempts ...int) error {
	maxAttempts := 1
	if len(attempts) > 0 && attempts[0] > 0 {
		maxAttempts = attempts[0]
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := c.BeginTransaction(ctx); err != nil {
			return err
		}

		err := fn(c)
		if err != nil {
			rbErr := c.Rollback(ctx)
			if rbErr != nil {
				return fmt.Errorf("%w (rollback: %v)", err, rbErr)
			}
			if isDeadlock(err) && attempt < maxAttempts {
				continue
			}
			return err
		}

		if err := c.Commit(ctx); err != nil {
			return err
		}
		return nil
	}

	return ErrTransactionFailed
}

// BeginTransaction starts a new database transaction or creates a savepoint
// for nested transactions.
func (c *Connection) BeginTransaction(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.txLevel == 0 {
		tx, err := c.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
		}
		c.tx = tx
	} else {
		savepointName := c.savepointName(c.txLevel + 1)
		_, err := c.tx.ExecContext(ctx, "SAVEPOINT "+savepointName)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
		}
	}

	c.txLevel++

	c.fireTransactionEvent(ctx, dbevents.TransactionBeginning{
		ConnectionName: c.name,
	})

	return nil
}

// Commit commits the active database transaction.
func (c *Connection) Commit(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.txLevel <= 0 {
		return nil
	}

	c.fireTransactionEvent(ctx, dbevents.TransactionCommitting{
		ConnectionName: c.name,
	})

	if c.txLevel == 1 {
		if err := c.tx.Commit(); err != nil {
			return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
		}
		c.tx = nil
		c.txMgr.Commit(c.txLevel)
	} else {
		savepointName := c.savepointName(c.txLevel)
		_, _ = c.tx.ExecContext(ctx, "RELEASE SAVEPOINT "+savepointName)
	}

	c.txLevel--

	c.fireTransactionEvent(ctx, dbevents.TransactionCommitted{
		ConnectionName: c.name,
	})

	return nil
}

// Rollback rolls back the active database transaction or to a savepoint.
func (c *Connection) Rollback(ctx context.Context, toLevel ...int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.txLevel <= 0 {
		return nil
	}

	level := 0
	if len(toLevel) > 0 {
		level = toLevel[0]
	}

	if c.txLevel == 1 || level == 0 {
		if c.tx != nil {
			if err := c.tx.Rollback(); err != nil && err != sql.ErrTxDone {
				return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
			}
			c.tx = nil
		}
		c.txMgr.Rollback(c.txLevel)
		c.txLevel = 0
	} else {
		savepointName := c.savepointName(c.txLevel)
		if c.tx != nil {
			_, _ = c.tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT "+savepointName)
		}
		c.txMgr.Rollback(c.txLevel)
		c.txLevel--
	}

	c.fireTransactionEvent(ctx, dbevents.TransactionRolledBack{
		ConnectionName: c.name,
	})

	return nil
}

// TransactionLevel returns the current transaction nesting depth.
func (c *Connection) TransactionLevel() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.txLevel
}

// AfterCommit registers a callback to run after the outermost transaction commits.
func (c *Connection) AfterCommit(fn func()) {
	c.txMgr.AfterCommit(c.txLevel, fn)
}

func (c *Connection) savepointName(level int) string {
	return fmt.Sprintf("savepoint_%d", level)
}

func (c *Connection) fireTransactionEvent(ctx context.Context, event any) {
	if c.events != nil {
		c.events.Dispatch(ctx, event)
	}
}

// isDeadlock checks if the error is a deadlock.
func isDeadlock(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return contains(msg, "deadlock") || contains(msg, "lock wait timeout")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
