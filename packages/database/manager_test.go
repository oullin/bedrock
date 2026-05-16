package database_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/database"
)

func TestManagerGetDefaultConnection(t *testing.T) {
	t.Parallel()

	m := database.NewManager()
	m.SetDefaultConnection("sqlite")

	if m.GetDefaultConnection() != "sqlite" {
		t.Fatalf("expected sqlite, got %s", m.GetDefaultConnection())
	}
}

func TestManagerConnectionNotConfigured(t *testing.T) {
	t.Parallel()

	m := database.NewManager()
	m.SetDefaultConnection("missing")

	_, err := m.Connection(context.Background())

	if err == nil {
		t.Fatal("expected error for unconfigured connection")
	}
}

func TestManagerExtend(t *testing.T) {
	t.Parallel()

	m := database.NewManager()
	m.Extend("custom", func(config database.ConnectionConfig) (*database.Connection, error) {
		return database.NewConnection(nil, "custom", "test", "", nil), nil
	})

	m.AddConnection("custom", database.ConnectionConfig{Driver: "custom", Database: "test"})
	m.SetDefaultConnection("custom")

	conn, err := m.Connection(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if conn.GetDatabaseName() != "test" {
		t.Fatalf("expected database test, got %s", conn.GetDatabaseName())
	}
}

func TestManagerSupportedDrivers(t *testing.T) {
	t.Parallel()

	m := database.NewManager()
	drivers := m.SupportedDrivers()

	if len(drivers) != 5 {
		t.Fatalf("expected 5 drivers, got %d", len(drivers))
	}

	want := map[string]bool{"mysql": true, "mariadb": true, "pgsql": true, "sqlite": true, "clickhouse": true}

	for _, d := range drivers {
		if !want[d] {
			t.Fatalf("unexpected driver %q", d)
		}

		delete(want, d)
	}

	if len(want) > 0 {
		t.Fatalf("missing drivers: %v", want)
	}
}

func TestManagerPurge(t *testing.T) {
	t.Parallel()

	m := database.NewManager()
	m.Extend("test", func(config database.ConnectionConfig) (*database.Connection, error) {
		return database.NewConnection(nil, "test", "db", "", nil), nil
	})
	m.AddConnection("test", database.ConnectionConfig{Driver: "test"})

	_, err := m.Connection(context.Background(), "test")

	if err != nil {
		t.Fatal(err)
	}

	conns := m.GetConnections()

	if len(conns) != 1 {
		t.Fatalf("expected 1 connection, got %d", len(conns))
	}

	m.Purge("test")

	conns = m.GetConnections()

	if len(conns) != 0 {
		t.Fatalf("expected 0 connections after purge, got %d", len(conns))
	}
}
