package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

// Open opens a shared SQL connection.
func Open(ctx context.Context, cfg Config) (*sql.DB, error) {
	driver := strings.TrimSpace(strings.ToLower(cfg.Driver))

	if driver == "" {
		driver = "sqlite"
	}

	switch driver {
	case "sqlite":
	default:
		return nil, fmt.Errorf("database: unsupported driver %q", cfg.Driver)
	}

	db, err := sql.Open(driver, cfg.DSN)

	if err != nil {
		return nil, fmt.Errorf("database: open connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()

		return nil, fmt.Errorf("database: ping connection: %w", err)
	}

	return db, nil
}
