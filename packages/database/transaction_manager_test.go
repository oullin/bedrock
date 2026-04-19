package database_test

import (
	"testing"

	"github.com/bedrock/packages/database"
)

func TestTransactionManagerAfterCommit(t *testing.T) {
	t.Parallel()

	mgr := database.NewTransactionManager()

	called := false
	mgr.AfterCommit(1, func() { called = true })
	mgr.Commit(1)

	if !called {
		t.Fatal("expected afterCommit callback to be called")
	}
}

func TestTransactionManagerRollbackDiscards(t *testing.T) {
	t.Parallel()

	mgr := database.NewTransactionManager()

	called := false
	mgr.AfterCommit(1, func() { called = true })
	mgr.Rollback(1)

	if called {
		t.Fatal("expected afterCommit callback NOT to be called after rollback")
	}
}
