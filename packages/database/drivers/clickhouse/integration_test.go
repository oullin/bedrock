//go:build integration

// This file is compiled only with `-tags=integration`. It exercises the
// driver against a real ClickHouse instance pointed at by the CLICKHOUSE_DSN
// env var (e.g. `clickhouse://default:@localhost:9000/default`). The test
// skips when the env var is missing so contributors without a running
// ClickHouse instance can still run the integration suite.
//
// The clickhouse-go/v2 driver is required at runtime; it is registered via
// the blank import in clickhouse_driver_test.go.
package clickhouse_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/bedrock/packages/database"
	"github.com/bedrock/packages/database/drivers/clickhouse"
	"github.com/bedrock/packages/database/migrations"
	"github.com/bedrock/packages/database/schema"
	schemagrammars "github.com/bedrock/packages/database/schema/grammars"
)

func newIntegrationConn(t *testing.T) *database.Connection {
	t.Helper()

	dsn := os.Getenv("CLICKHOUSE_DSN")

	if dsn == "" {
		t.Skip("CLICKHOUSE_DSN not set; skipping integration test")
	}

	cfg, err := database.ParseDatabaseURL(dsn)

	if err != nil {
		t.Fatalf("invalid CLICKHOUSE_DSN: %v", err)
	}

	conn, err := (&clickhouse.Connector{}).Connect(*cfg)

	if err != nil {
		t.Fatalf("connect: %v", err)
	}

	t.Cleanup(func() {
		conn.DB().Close()
	})

	return conn
}

func TestClickHouse_PingAndQuery(t *testing.T) {
	conn := newIntegrationConn(t)

	if err := conn.DB().PingContext(context.Background()); err != nil {
		t.Fatalf("ping failed: %v", err)
	}

	rows, err := conn.Select(context.Background(), "SELECT 42 AS answer")

	if err != nil {
		t.Fatalf("select: %v", err)
	}

	if len(rows) != 1 || rows[0]["answer"] == nil {
		t.Fatalf("unexpected rows: %+v", rows)
	}
}

func TestClickHouse_CreateInsertSelectDeleteRoundtrip(t *testing.T) {
	conn := newIntegrationConn(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	table := "bedrock_clickhouse_it_" + time.Now().Format("150405")
	builder := schema.NewBuilder(conn, schemagrammars.NewClickHouseGrammar())

	defer func() {
		_, _ = conn.Statement(context.Background(), "DROP TABLE IF EXISTS "+table)
	}()

	err := builder.Create(ctx, table, func(bp *schema.Blueprint) {
		bp.UnsignedBigInteger("id")
		bp.String("payload")
		bp.DateTime("ts", 3)
		bp.SetOrderBy("id")
	})

	if err != nil {
		t.Fatalf("create: %v", err)
	}

	if _, err := conn.Insert(ctx,
		"INSERT INTO "+table+" (id, payload, ts) VALUES (?, ?, now64(3))",
		uint64(1), "hello world"); err != nil {
		t.Fatalf("insert: %v", err)
	}

	rows, err := conn.Select(ctx, "SELECT id, payload FROM "+table+" WHERE id = ?", uint64(1))

	if err != nil {
		t.Fatalf("select: %v", err)
	}

	if len(rows) != 1 || rows[0]["payload"] != "hello world" {
		t.Fatalf("unexpected row: %+v", rows)
	}

	// ALTER TABLE ... DELETE is asynchronous; wait briefly for it to settle.
	if _, err := conn.Statement(ctx, "ALTER TABLE "+table+" DELETE WHERE id = 1"); err != nil {
		t.Fatalf("alter delete: %v", err)
	}
}

func TestClickHouse_MigrationsRepository(t *testing.T) {
	conn := newIntegrationConn(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	table := "bedrock_clickhouse_migrations_it"

	defer func() {
		_, _ = conn.Statement(context.Background(), "DROP TABLE IF EXISTS "+table)
	}()

	repo := migrations.NewDatabaseRepository(conn, table)
	clickhouse.ConfigureMigrationsRepository(repo)

	if err := repo.CreateRepository(ctx); err != nil {
		t.Fatalf("create migrations table: %v", err)
	}

	exists, err := repo.RepositoryExists(ctx)

	if err != nil {
		t.Fatalf("repo exists: %v", err)
	}

	if !exists {
		t.Fatal("expected migrations repo to exist")
	}
}
