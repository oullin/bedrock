package database

import (
	"fmt"
	"strings"

	configpkg "github.com/gollin/packages/config"
)

// Config controls shared database behavior.
type Config struct {
	DefaultConnection string
	Driver            string
	DSN               string
	MigrationsTable   string
}

// ConfigFromRepository loads the active database connection config.
func ConfigFromRepository(repo *configpkg.Repository) (Config, error) {
	if repo == nil {
		return Config{}, fmt.Errorf("database: config repository is required")
	}

	connection, err := repo.String("database.default")
	if err != nil {
		return Config{}, err
	}

	connection = strings.TrimSpace(connection)
	if connection == "" {
		connection = "sqlite"
	}

	driver, err := repo.String("database.connections." + connection + ".driver")
	if err != nil {
		return Config{}, err
	}

	dsn, err := repo.String("database.connections." + connection + ".dsn")
	if err != nil {
		return Config{}, err
	}

	table, err := repo.String("database.migrations.table")
	if err != nil {
		return Config{}, err
	}

	table = strings.TrimSpace(table)
	if table == "" {
		table = "schema_migrations"
	}

	return Config{
		DefaultConnection: connection,
		Driver:            strings.TrimSpace(driver),
		DSN:               strings.TrimSpace(dsn),
		MigrationsTable:   table,
	}, nil
}
