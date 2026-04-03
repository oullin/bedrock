package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

// Migrate applies the embedded workspace migrations.
func Migrate(ctx context.Context, db *sql.DB, cfg Config) error {
	return MigrateFS(ctx, db, embeddedMigrations, cfg.MigrationsTable)
}

// MigrateFS applies SQL files from an fs in filename order.
func MigrateFS(ctx context.Context, db *sql.DB, fsys fs.FS, table string) error {
	if db == nil {
		return fmt.Errorf("database: db is required")
	}

	table = strings.TrimSpace(table)
	if table == "" {
		table = "schema_migrations"
	}

	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS `+table+` (name TEXT PRIMARY KEY)`); err != nil {
		return fmt.Errorf("database: create migrations table: %w", err)
	}

	entries, err := fs.ReadDir(fsys, "migrations")
	if err != nil {
		return fmt.Errorf("database: read migrations: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if filepath.Ext(name) == ".sql" {
			names = append(names, name)
		}
	}

	slices.Sort(names)

	for _, name := range names {
		var applied string
		err := db.QueryRowContext(ctx, `SELECT name FROM `+table+` WHERE name = ?`, name).Scan(&applied)
		if err == nil && applied == name {
			continue
		}

		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("database: check migration %q: %w", name, err)
		}

		payload, err := fs.ReadFile(fsys, "migrations/"+name)
		if err != nil {
			return fmt.Errorf("database: read migration %q: %w", name, err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("database: begin migration %q: %w", name, err)
		}

		if _, err := tx.ExecContext(ctx, string(payload)); err != nil {
			_ = tx.Rollback()

			return fmt.Errorf("database: run migration %q: %w", name, err)
		}

		if _, err := tx.ExecContext(ctx, `INSERT INTO `+table+` (name) VALUES (?)`, name); err != nil {
			_ = tx.Rollback()

			return fmt.Errorf("database: record migration %q: %w", name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("database: commit migration %q: %w", name, err)
		}
	}

	return nil
}
