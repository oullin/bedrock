package user_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	authaccess "github.com/gollin/packages/auth/access"
	configpkg "github.com/gollin/packages/config"
	databasepkg "github.com/gollin/packages/database"
	"github.com/gollin/packages/user"
)

func TestConfigNewMemoryRepositoryAndSQLRepository(t *testing.T) {
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

	cfg, err := user.ConfigFromRepository(repo)

	if err != nil {
		t.Fatalf("ConfigFromRepository: %v", err)
	}

	record := user.New(cfg, "user-1", "", " USER@example.com ", time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC))

	if record.Name != "User" {
		t.Fatalf("unexpected default name: %q", record.Name)
	}

	if record.Email != "user@example.com" {
		t.Fatalf("unexpected normalized email: %q", record.Email)
	}

	hasher, err := auth.NewDefaultPasswordHasher()

	if err != nil {
		t.Fatalf("NewDefaultPasswordHasher: %v", err)
	}

	record.PasswordHash, err = hasher.Hash(context.Background(), "secret")

	if err != nil {
		t.Fatalf("Hash: %v", err)
	}

	users, err := user.NewMemoryRepository(hasher)

	if err != nil {
		t.Fatalf("NewMemoryRepository: %v", err)
	}

	if err := users.Create(context.Background(), record); err != nil {
		t.Fatalf("Create: %v", err)
	}

	found, err := users.FindByEmail(context.Background(), record.Email)

	if err != nil {
		t.Fatalf("FindByEmail: %v", err)
	}

	if found.GetAuthIdentifier() != record.GetAuthIdentifier() {
		t.Fatalf("unexpected auth identifier: %q", found.GetAuthIdentifier())
	}

	valid, err := users.ValidateCredentials(context.Background(), found, map[string]string{"password": "secret"})

	if err != nil {
		t.Fatalf("ValidateCredentials: %v", err)
	}

	if !valid {
		t.Fatal("expected credentials to be valid")
	}

	gate := authaccess.NewGate()
	gate.Define("view-dashboard", func(context.Context, auth.Authenticatable, ...any) authaccess.Response {
		return authaccess.Allow()
	})

	if !record.Can(context.Background(), gate, "view-dashboard") {
		t.Fatal("expected user authorization helper to allow ability")
	}

	dbPath := filepath.Join(t.TempDir(), "users.sqlite")
	db, err := databasepkg.Open(context.Background(), databasepkg.Config{
		Driver:          "sqlite",
		DSN:             dbPath,
		MigrationsTable: "schema_migrations",
	})

	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	defer db.Close()

	if err := databasepkg.Migrate(context.Background(), db, databasepkg.Config{MigrationsTable: "schema_migrations"}); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	sqlUsers, err := user.NewSQLRepository(db, hasher, "")

	if err != nil {
		t.Fatalf("NewSQLRepository: %v", err)
	}

	if err := sqlUsers.Create(context.Background(), record); err != nil {
		t.Fatalf("SQL Create: %v", err)
	}

	sqlFound, err := sqlUsers.FindByEmail(context.Background(), record.Email)

	if err != nil {
		t.Fatalf("SQL FindByEmail: %v", err)
	}

	if sqlFound.GetAuthIdentifier() != record.GetAuthIdentifier() {
		t.Fatalf("unexpected SQL auth identifier: %q", sqlFound.GetAuthIdentifier())
	}
}
