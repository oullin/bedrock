package database_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"

	configpkg "github.com/gollin/packages/config"
	databasepkg "github.com/gollin/packages/database"
)

func TestConfigFromRepositoryAndMigrate(t *testing.T) {
	t.Parallel()

	_, currentFile, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	configDir := filepath.Join(filepath.Dir(currentFile), "config")
	repo, err := configpkg.NewBuilder(configDir).Build(context.Background())

	if err != nil {
		t.Fatalf("build config: %v", err)
	}

	dbPath := filepath.Join(t.TempDir(), "database.sqlite")
	repo.Set("database.connections.sqlite.dsn", dbPath)

	cfg, err := databasepkg.ConfigFromRepository(repo)

	if err != nil {
		t.Fatalf("ConfigFromRepository: %v", err)
	}

	db, err := databasepkg.Open(context.Background(), cfg)

	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	defer db.Close()

	if err := databasepkg.Migrate(context.Background(), db, cfg); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	var tableCount int

	if err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name IN ('users', 'auth_sessions', 'password_reset_tokens')`).Scan(&tableCount); err != nil {
		t.Fatalf("count migrated tables: %v", err)
	}

	if tableCount != 3 {
		t.Fatalf("expected 3 migrated tables, got %d", tableCount)
	}
}
