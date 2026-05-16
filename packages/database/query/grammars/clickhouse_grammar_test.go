package grammars_test

import (
	"strings"
	"testing"

	"github.com/bedrock/packages/database/query"
	"github.com/bedrock/packages/database/query/grammars"
)

func newClickHouseBuilder() *query.Builder {
	g := grammars.NewClickHouseGrammar()

	return query.NewBuilder(nil, g, nil).From("events")
}

func TestClickHouseBasicSelect(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder()
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `events`"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseWhereWithBacktickQuoting(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().Where("user_id", "=", 42)
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `events` where `user_id` = ?"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseLimitOffset(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().Limit(50).Offset(100)
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `events` limit 50 offset 100"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseInsert(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	b := query.NewBuilder(nil, g, nil).From("events")
	sql := g.CompileInsert(b, []map[string]any{
		{"id": 1, "name": "hello"},
	})
	expected := "insert into `events` (`id`, `name`) values (?, ?)"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

// CompileUpdate must emit ALTER TABLE ... UPDATE; per-row UPDATE is not valid.
func TestClickHouseUpdateIsAlterMutation(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	b := query.NewBuilder(nil, g, nil).From("events").Where("id", "=", 1)
	sql := g.CompileUpdate(b, map[string]any{"name": "renamed"})
	expected := "alter table `events` update `name` = ? where `id` = ?"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

// CompileDelete must emit ALTER TABLE ... DELETE.
func TestClickHouseDeleteIsAlterMutation(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	b := query.NewBuilder(nil, g, nil).From("events").Where("id", "=", 1)
	sql := g.CompileDelete(b)
	expected := "alter table `events` delete where `id` = ?"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

// ClickHouse rejects mutations without a WHERE clause; the grammar synthesises
// `where 1=1` so callers cannot accidentally emit a server-side error.
func TestClickHouseDeleteWithoutWhereGetsTautology(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	b := query.NewBuilder(nil, g, nil).From("events")
	sql := g.CompileDelete(b)
	expected := "alter table `events` delete where 1=1"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

// Upsert must surface a server-side throwIf so the caller fails loudly.
func TestClickHouseUpsertEmitsThrow(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	b := query.NewBuilder(nil, g, nil).From("events")
	sql := g.CompileUpsert(b, []map[string]any{{"id": 1}}, []string{"id"}, []string{"id"})

	if !strings.Contains(sql, "throwIf") || !strings.Contains(sql, "upsert is not supported") {
		t.Fatalf("expected throwIf with upsert message, got %q", sql)
	}
}

func TestClickHouseParameterUsesQuestionMark(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()

	if got := g.Parameter("anything"); got != "?" {
		t.Fatalf("expected ?, got %q", got)
	}
}
