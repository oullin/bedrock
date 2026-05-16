package grammars_test

import (
	"strings"
	"testing"

	"github.com/bedrock/packages/database/schema"
	"github.com/bedrock/packages/database/schema/grammars"
)

func TestClickHouseCreate_DefaultsToMergeTreeWithFirstColumnOrderBy(t *testing.T) {
	t.Parallel()

	bp := schema.NewBlueprint("events")
	bp.UnsignedBigInteger("id")
	bp.String("name", 64)

	g := grammars.NewClickHouseGrammar()
	stmts := g.CompileCreate(bp)

	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}

	sql := stmts[0]

	if !strings.Contains(sql, "engine = MergeTree()") {
		t.Fatalf("expected default MergeTree engine, got %q", sql)
	}

	if !strings.Contains(sql, "order by (`id`)") {
		t.Fatalf("expected default ORDER BY (id), got %q", sql)
	}
}

func TestClickHouseCreate_HonorsExplicitEngineAndOrderBy(t *testing.T) {
	t.Parallel()

	bp := schema.NewBlueprint("events")
	bp.UnsignedBigInteger("id")
	bp.DateTime("ts", 3)
	bp.SetEngine("ReplacingMergeTree")
	bp.SetOrderBy("ts", "id")
	bp.SetPartitionBy("toYYYYMM(ts)")

	g := grammars.NewClickHouseGrammar()
	stmts := g.CompileCreate(bp)
	sql := stmts[0]

	if !strings.Contains(sql, "engine = ReplacingMergeTree()") {
		t.Fatalf("expected ReplacingMergeTree, got %q", sql)
	}

	if !strings.Contains(sql, "partition by toYYYYMM(ts)") {
		t.Fatalf("expected partition by clause, got %q", sql)
	}

	if !strings.Contains(sql, "order by (`ts`, `id`)") {
		t.Fatalf("expected order by (ts, id), got %q", sql)
	}
}

func TestClickHouseNullableWrapsType(t *testing.T) {
	t.Parallel()

	bp := schema.NewBlueprint("events")
	bp.String("note").Nullable()

	g := grammars.NewClickHouseGrammar()
	sql := g.CompileCreate(bp)[0]

	if !strings.Contains(sql, "`note` Nullable(") {
		t.Fatalf("expected Nullable wrapper, got %q", sql)
	}
}

func TestClickHouseTypeMapping(t *testing.T) {
	t.Parallel()

	bp := schema.NewBlueprint("typed")
	bp.BigInteger("i64")
	bp.UnsignedBigInteger("u64")
	bp.Integer("i32")
	bp.UnsignedSmallInteger("u16")
	bp.Float("f32")
	bp.Double("f64")
	bp.Decimal("amount", 18, 4)
	bp.Boolean("flag")
	bp.UUID("uid")
	bp.DateTime("ts", 6)
	bp.Date("d")
	bp.Text("body")
	bp.JSON("payload")

	g := grammars.NewClickHouseGrammar()
	sql := g.CompileCreate(bp)[0]

	checks := map[string]string{
		"Int64":          "`i64` Int64",
		"UInt64":         "`u64` UInt64",
		"Int32":          "`i32` Int32",
		"UInt16":         "`u16` UInt16",
		"Float32":        "`f32` Float32",
		"Float64":        "`f64` Float64",
		"Decimal":        "`amount` Decimal(18, 4)",
		"Boolean->UInt8": "`flag` UInt8",
		"UUID":           "`uid` UUID",
		"DateTime64":     "`ts` DateTime64(6)",
		"Date":           "`d` Date",
		"Text->String":   "`body` String",
		"JSON->String":   "`payload` String",
	}

	for label, fragment := range checks {
		if !strings.Contains(sql, fragment) {
			t.Errorf("%s: expected fragment %q in %q", label, fragment, sql)
		}
	}
}

func TestClickHouseAlterStatements(t *testing.T) {
	t.Parallel()

	g := grammars.NewClickHouseGrammar()
	bp := schema.NewBlueprint("events")
	bp.String("note", 32)

	if got := g.CompileAdd(bp); len(got) != 1 || !strings.HasPrefix(got[0], "alter table `events` add column ") {
		t.Fatalf("expected ALTER TABLE ADD COLUMN, got %v", got)
	}

	if got := g.CompileDropColumn(bp, []string{"note"}); got != "alter table `events` drop column `note`" {
		t.Fatalf("expected drop column statement, got %q", got)
	}

	if got := g.CompileRenameColumn(bp, "note", "comment"); got != "alter table `events` rename column `note` to `comment`" {
		t.Fatalf("expected rename column statement, got %q", got)
	}
}

func TestClickHouseForeignKeysAreNoOp(t *testing.T) {
	t.Parallel()

	g := grammars.NewClickHouseGrammar()
	bp := schema.NewBlueprint("events")

	if got := g.CompileCreateForeignKey(bp, &schema.ForeignKeyDefinition{}); got != "" {
		t.Fatalf("expected empty SQL for foreign key create, got %q", got)
	}

	if got := g.CompileDropForeignKey(bp, "fk"); got != "" {
		t.Fatalf("expected empty SQL for foreign key drop, got %q", got)
	}

	if got := g.CompileEnableForeignKeyConstraints(); got != "" {
		t.Fatalf("expected empty SQL for enable FK constraints, got %q", got)
	}
}

func TestClickHouseTableExistsUsesSystemTables(t *testing.T) {
	t.Parallel()

	g := grammars.NewClickHouseGrammar()

	if !strings.Contains(g.CompileTableExists(), "system.tables") {
		t.Fatalf("expected lookup against system.tables, got %q", g.CompileTableExists())
	}

	if !strings.Contains(g.CompileColumnListing("events"), "system.columns") {
		t.Fatalf("expected lookup against system.columns, got %q", g.CompileColumnListing("events"))
	}
}
