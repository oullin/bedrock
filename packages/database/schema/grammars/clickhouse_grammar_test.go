package grammars_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/bedrock/packages/database"
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

func TestClickHouseForeignKeysReturnSentinel(t *testing.T) {
	t.Parallel()

	g := grammars.NewClickHouseGrammar()
	bp := schema.NewBlueprint("events")

	sql, err := g.CompileCreateForeignKey(bp, &schema.ForeignKeyDefinition{})

	if sql != "" {
		t.Fatalf("expected empty SQL for foreign key create, got %q", sql)
	}

	if !errors.Is(err, database.ErrForeignKeysNotSupported) {
		t.Fatalf("expected ErrForeignKeysNotSupported, got %v", err)
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

func TestClickHouseDropStatements(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()

	if got := g.CompileDrop("events"); got != "drop table `events`" {
		t.Errorf("CompileDrop = %q", got)
	}

	if got := g.CompileDropIfExists("events"); got != "drop table if exists `events`" {
		t.Errorf("CompileDropIfExists = %q", got)
	}

	if got := g.CompileRename("events", "events_v2"); got != "rename table `events` to `events_v2`" {
		t.Errorf("CompileRename = %q", got)
	}
}

func TestClickHouseDisableForeignKeyConstraintsIsNoOp(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()

	if got := g.CompileDisableForeignKeyConstraints(); got != "" {
		t.Errorf("expected empty SQL; got %q", got)
	}
}

func TestClickHouseCompileChangeUsesModifyColumn(t *testing.T) {
	t.Parallel()

	bp := schema.NewBlueprint("events")
	bp.String("note", 64).Change()

	g := grammars.NewClickHouseGrammar()
	stmts := g.CompileChange(bp)

	if len(stmts) == 0 || !strings.Contains(stmts[0], "modify column `note`") {
		t.Fatalf("expected modify column; got %v", stmts)
	}
}

func TestClickHouseCompileCreateIndex(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	bp := schema.NewBlueprint("events")

	cases := []struct {
		name, cmdName, indexName string
		wantContains             []string
	}{
		{"unique", "unique", "events_user_id_unique", []string{"add index", "minmax granularity 1"}},
		{"index", "index", "events_user_id_index", []string{"add index", "minmax granularity 1"}},
		{"fulltext", "fulltext", "events_body_fulltext", []string{"add index", "tokenbf_v1"}},
	}

	for _, c := range cases {
		got := g.CompileCreateIndex(bp, schema.BlueprintCommand{
			Name: c.cmdName, Index: c.indexName, Columns: []string{"user_id"},
		})

		for _, frag := range c.wantContains {
			if !strings.Contains(got, frag) {
				t.Errorf("%s: missing %q in %q", c.name, frag, got)
			}
		}
	}
}

func TestClickHouseCompileCreateIndexSkipsPrimary(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	bp := schema.NewBlueprint("events")

	got := g.CompileCreateIndex(bp, schema.BlueprintCommand{
		Name: "primary", Index: "pk", Columns: []string{"id"},
	})

	if got != "" {
		t.Errorf("expected empty SQL for primary command; got %q", got)
	}
}

func TestClickHouseCompileDropIndex(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	bp := schema.NewBlueprint("events")

	if got := g.CompileDropIndex(bp, "events_user_id_index"); got != "alter table `events` drop index `events_user_id_index`" {
		t.Errorf("CompileDropIndex = %q", got)
	}
}

func TestClickHouseColumnDefaultValues(t *testing.T) {
	t.Parallel()
	bp := schema.NewBlueprint("events")
	bp.String("status", 8).Default("active")
	bp.Boolean("active").Default(true)
	bp.Boolean("disabled").Default(false)
	bp.Integer("count").Default(42)

	g := grammars.NewClickHouseGrammar()
	sql := g.CompileCreate(bp)[0]

	if !strings.Contains(sql, "`status` FixedString(8) default 'active'") {
		t.Errorf("expected string default; got %q", sql)
	}

	if !strings.Contains(sql, "`active` UInt8 default 1") {
		t.Errorf("expected boolean=true default 1; got %q", sql)
	}

	if !strings.Contains(sql, "`disabled` UInt8 default 0") {
		t.Errorf("expected boolean=false default 0; got %q", sql)
	}

	if !strings.Contains(sql, "`count` Int32 default '42'") {
		t.Errorf("expected int default fallback; got %q", sql)
	}
}

func TestClickHouseColumnComment(t *testing.T) {
	t.Parallel()
	bp := schema.NewBlueprint("events")
	bp.String("note", 16).Comment("a note")

	g := grammars.NewClickHouseGrammar()
	sql := g.CompileCreate(bp)[0]

	if !strings.Contains(sql, "comment 'a note'") {
		t.Errorf("expected comment clause; got %q", sql)
	}
}

func TestClickHouseTinyAndMediumIntegers(t *testing.T) {
	t.Parallel()
	bp := schema.NewBlueprint("ints")
	bp.TinyInteger("i8")
	bp.UnsignedTinyInteger("u8")
	bp.MediumInteger("i24")
	bp.UnsignedMediumInteger("u24")

	g := grammars.NewClickHouseGrammar()
	sql := g.CompileCreate(bp)[0]

	for _, frag := range []string{"`i8` Int8", "`u8` UInt8", "`i24` Int32", "`u24` UInt32"} {
		if !strings.Contains(sql, frag) {
			t.Errorf("missing %q in %q", frag, sql)
		}
	}
}

func TestClickHouseDefaultOrderColumnPicksAutoIncrement(t *testing.T) {
	t.Parallel()
	bp := schema.NewBlueprint("events")
	bp.String("name", 16)
	bp.BigIncrements("id") // auto-increment via helper

	g := grammars.NewClickHouseGrammar()
	sql := g.CompileCreate(bp)[0]

	if !strings.Contains(sql, "order by (`id`)") {
		t.Errorf("expected auto-increment column as ORDER BY; got %q", sql)
	}
}

func TestClickHouseEmptyBlueprintFallsBackToTuple(t *testing.T) {
	t.Parallel()
	bp := schema.NewBlueprint("empty")

	g := grammars.NewClickHouseGrammar()
	stmts := g.CompileCreate(bp)

	if len(stmts) == 0 || !strings.Contains(stmts[0], "order by (`tuple()`)") {
		// tuple() may be unwrapped or wrapped; accept either form provided
		// "tuple" appears in the order-by.
		if len(stmts) == 0 || !strings.Contains(stmts[0], "tuple") {
			t.Errorf("expected tuple() fallback ORDER BY; got %v", stmts)
		}
	}
}
