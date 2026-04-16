package database_test

import (
	"testing"

	"github.com/bedrock/packages/database"
)

func TestParseDatabaseURL_MySQL(t *testing.T) {
	t.Parallel()

	cfg, err := database.ParseDatabaseURL("mysql://user:pass@localhost:3306/mydb")
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Driver != "mysql" {
		t.Fatalf("expected driver mysql, got %s", cfg.Driver)
	}
	if cfg.Host != "localhost" {
		t.Fatalf("expected host localhost, got %s", cfg.Host)
	}
	if cfg.Port != 3306 {
		t.Fatalf("expected port 3306, got %d", cfg.Port)
	}
	if cfg.Database != "mydb" {
		t.Fatalf("expected database mydb, got %s", cfg.Database)
	}
	if cfg.Username != "user" {
		t.Fatalf("expected username user, got %s", cfg.Username)
	}
	if cfg.Password != "pass" {
		t.Fatalf("expected password pass, got %s", cfg.Password)
	}
}

func TestParseDatabaseURL_Postgres(t *testing.T) {
	t.Parallel()

	cfg, err := database.ParseDatabaseURL("postgres://admin:secret@db.example.com:5432/production?sslmode=require")
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Driver != "postgres" {
		t.Fatalf("expected driver postgres, got %s", cfg.Driver)
	}
	if cfg.Port != 5432 {
		t.Fatalf("expected port 5432, got %d", cfg.Port)
	}
	if cfg.SSLMode != "require" {
		t.Fatalf("expected sslmode require, got %s", cfg.SSLMode)
	}
}

func TestParseDatabaseURL_SQLite(t *testing.T) {
	t.Parallel()

	cfg, err := database.ParseDatabaseURL("sqlite:///path/to/database.db")
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Driver != "sqlite" {
		t.Fatalf("expected driver sqlite, got %s", cfg.Driver)
	}
	if cfg.Database != "path/to/database.db" {
		t.Fatalf("expected database path/to/database.db, got %s", cfg.Database)
	}
}
