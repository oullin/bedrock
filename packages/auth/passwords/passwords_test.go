package passwords_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	auth "github.com/gollin/packages/auth"
	"github.com/gollin/packages/auth/passwords"
	configpkg "github.com/gollin/packages/config"
	databasepkg "github.com/gollin/packages/database"
	"github.com/gollin/packages/user"
)

type resetter struct {
	hasher auth.PasswordHasher
}

func (r resetter) Reset(ctx context.Context, user auth.Authenticatable, password string) error {
	hash, err := r.hasher.Hash(ctx, password)
	if err != nil {
		return err
	}

	user.SetAuthPassword(hash)

	return nil
}

func TestBrokerCreateAndValidateToken(t *testing.T) {
	t.Parallel()

	hasher, err := auth.NewDefaultPasswordHasher()
	if err != nil {
		t.Fatalf("NewDefaultPasswordHasher: %v", err)
	}

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	userConfigDir := filepath.Join(filepath.Dir(currentFile), "..", "..", "user", "config")
	repo, err := configpkg.NewBuilder(userConfigDir).Build(context.Background())
	if err != nil {
		t.Fatalf("build user config: %v", err)
	}

	userCfg, err := user.ConfigFromRepository(repo)
	if err != nil {
		t.Fatalf("ConfigFromRepository: %v", err)
	}

	record := user.New(userCfg, "user-1", "User", "user@example.com", time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC))
	record.PasswordHash, err = hasher.Hash(context.Background(), "secret")
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}

	users, err := user.NewMemoryRepository(hasher)
	if err != nil {
		t.Fatalf("NewMemoryRepository: %v", err)
	}

	if err := users.Create(context.Background(), record); err != nil {
		t.Fatalf("Create user: %v", err)
	}

	broker := &passwords.Broker{
		Config: passwords.Config{
			Name:     "users",
			Expire:   time.Hour,
			Throttle: time.Minute,
		},
		Users:  users,
		Tokens: passwords.NewMemoryTokenRepository(),
		Clock:  auth.SystemClock{},
	}

	token, err := broker.CreateToken(context.Background(), record)
	if err != nil {
		t.Fatalf("CreateToken: %v", err)
	}

	valid, err := broker.TokenExists(context.Background(), record, token)
	if err != nil {
		t.Fatalf("TokenExists: %v", err)
	}

	if !valid {
		t.Fatal("expected token to be valid")
	}

	dbPath := filepath.Join(t.TempDir(), "passwords.sqlite")
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

	sqlTokens := passwords.NewSQLTokenRepository(db, "")
	if err := sqlTokens.Save(context.Background(), &passwords.Token{
		UserID:    record.ID,
		TokenHash: "hash",
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}); err != nil {
		t.Fatalf("Save SQL token: %v", err)
	}
}
