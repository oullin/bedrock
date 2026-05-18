package grammars_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/bedrock/packages/database"
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

// Upsert must surface the typed sentinel so callers fail loudly before any
// SQL is sent.
func TestClickHouseUpsertReturnsSentinel(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	b := query.NewBuilder(nil, g, nil).From("events")
	sql, err := g.CompileUpsert(b, []map[string]any{{"id": 1}}, []string{"id"}, []string{"id"})

	if sql != "" {
		t.Fatalf("expected empty SQL, got %q", sql)
	}

	if !errors.Is(err, database.ErrUpsertNotSupported) {
		t.Fatalf("expected ErrUpsertNotSupported, got %v", err)
	}
}

func TestClickHouseParameterUsesQuestionMark(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()

	if got := g.Parameter("anything"); got != "?" {
		t.Fatalf("expected ?, got %q", got)
	}
}

func TestClickHouseSelectColumns(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().Select("id", "name")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select `id`, `name` from `events`"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseSelectDistinct(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().Distinct().Select("country")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select distinct `country` from `events`"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseWhereIn(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereIn("status", []any{"active", "pending"})
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `events` where `status` in (?, ?)"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseWhereNotIn(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereNotIn("status", []any{"deleted"})
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `events` where `status` not in (?)"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseWhereNull(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereNull("deleted_at")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `events` where `deleted_at` is null"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseWhereNotNull(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereNotNull("verified_at")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `events` where `verified_at` is not null"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseWhereBetween(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereBetween("age", [2]any{18, 65})
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `events` where `age` between ? and ?"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseOrWhere(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().Where("type", "=", "click").OrWhere("type", "=", "view")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `events` where `type` = ? or `type` = ?"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseGroupBy(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().GroupBy("country", "city")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `events` group by `country`, `city`"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseHaving(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().GroupBy("country").Having("count", ">", 100)
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `events` group by `country` having `count` > ?"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseOrderBy(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().OrderBy("created_at", "desc")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `events` order by `created_at` desc"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseJoin(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().Join("users", "events.user_id", "=", "users.id")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `events` inner join `users` on `events`.`user_id` = `users`.`id`"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseAggregateCount(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	b := query.NewBuilder(nil, g, nil).From("events")

	b.SetAggregate("count", []string{"*"})

	sql := g.CompileSelect(b)
	expected := "select count(*) as aggregate from `events`"

	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestClickHouseTruncateReturnsMap(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	b := query.NewBuilder(nil, g, nil).From("events")

	got := g.CompileTruncate(b)

	if len(got) == 0 {
		t.Fatal("expected non-empty truncate map")
	}

	found := false

	for sql := range got {
		if strings.Contains(sql, "truncate") && strings.Contains(sql, "`events`") {
			found = true

			break
		}
	}

	if !found {
		t.Fatalf("expected truncate sql referencing events; got %v", got)
	}
}

func TestClickHouseCompileExists(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	b := query.NewBuilder(nil, g, nil).From("events").Where("id", "=", 1)
	sql := g.CompileExists(b)

	if !strings.Contains(sql, "exists") {
		t.Fatalf("expected 'exists' keyword in %q", sql)
	}
}

func TestClickHouseRandomReturnsRand(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()

	if got := g.CompileRandom(""); got == "" {
		t.Fatal("expected non-empty rand expression")
	}
}

func TestClickHouseInsertGetIdReturnsSentinel(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	b := query.NewBuilder(nil, g, nil).From("events")
	sql, err := g.CompileInsertGetId(b, map[string]any{"name": "click"}, "id")

	if sql != "" {
		t.Fatalf("expected empty SQL, got %q", sql)
	}

	if !errors.Is(err, database.ErrInsertGetIdNotSupported) {
		t.Fatalf("expected ErrInsertGetIdNotSupported, got %v", err)
	}
}

func TestClickHouseInsertUsing(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	b := query.NewBuilder(nil, g, nil).From("events_copy")
	sql := g.CompileInsertUsing(b, []string{"id", "name"}, "select id, name from events")

	if !strings.Contains(sql, "`events_copy`") || !strings.Contains(sql, "select id, name from events") {
		t.Fatalf("expected insert ... select sql; got %q", sql)
	}
}

func TestClickHouseInsertOrIgnoreReturnsSentinel(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	sql, err := g.CompileInsertOrIgnore(nil, nil)

	if sql != "" {
		t.Fatalf("expected empty SQL, got %q", sql)
	}

	if !errors.Is(err, database.ErrInsertOrIgnoreNotSupported) {
		t.Fatalf("expected ErrInsertOrIgnoreNotSupported, got %v", err)
	}
}

func TestClickHouseWrap(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()

	cases := []struct {
		in, want string
	}{
		{"id", "`id`"},
		{"table.col", "`table`.`col`"},
		{"*", "*"},
		{"table.*", "`table`.*"},
	}

	for _, c := range cases {
		if got := g.Wrap(c.in); got != c.want {
			t.Errorf("Wrap(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestClickHouseWrapTablePrefixes(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	g.SetTablePrefix("acme_")

	if g.GetTablePrefix() != "acme_" {
		t.Errorf("GetTablePrefix = %q, want acme_", g.GetTablePrefix())
	}

	if got := g.WrapTable("events"); got != "`acme_events`" {
		t.Errorf("WrapTable = %q, want `acme_events`", got)
	}
}

func TestClickHouseColumnize(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	got := g.Columnize([]string{"id", "name"})

	if got != "`id`, `name`" {
		t.Errorf("Columnize = %q", got)
	}
}

func TestClickHouseParameterize(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	got := g.Parameterize([]any{1, "x", nil})

	if got != "?, ?, ?" {
		t.Errorf("Parameterize = %q", got)
	}
}

func TestClickHouseWhereDate(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereDate("created_at", "=", "2026-01-01")
	sql := b.GetGrammar().CompileSelect(b)

	if !strings.Contains(sql, "toDate(`created_at`)") {
		t.Errorf("expected toDate(); got %q", sql)
	}
}

func TestClickHouseWhereTime(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereTime("created_at", "<", "12:00:00")
	sql := b.GetGrammar().CompileSelect(b)

	if !strings.Contains(sql, "toTime(`created_at`)") {
		t.Errorf("expected toTime(); got %q", sql)
	}
}

func TestClickHouseWhereDayMonthYear(t *testing.T) {
	t.Parallel()

	for _, c := range []struct {
		method string
		fn     func(*query.Builder, string) *query.Builder
		want   string
	}{
		{"WhereDay", func(b *query.Builder, c string) *query.Builder { return b.WhereDay(c, "=", "15") }, "toDayOfMonth"},
		{"WhereMonth", func(b *query.Builder, c string) *query.Builder { return b.WhereMonth(c, "=", "3") }, "toMonth"},
		{"WhereYear", func(b *query.Builder, c string) *query.Builder { return b.WhereYear(c, "=", "2026") }, "toYear"},
	} {
		b := c.fn(newClickHouseBuilder(), "created_at")
		sql := b.GetGrammar().CompileSelect(b)

		if !strings.Contains(sql, c.want+"(`created_at`)") {
			t.Errorf("%s: expected %s(); got %q", c.method, c.want, sql)
		}
	}
}

func TestClickHouseWhereLike(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereLike("name", "%foo%")
	sql := b.GetGrammar().CompileSelect(b)

	if !strings.Contains(sql, "`name` like ?") {
		t.Errorf("expected 'like ?'; got %q", sql)
	}
}

func TestClickHouseWhereNotLike(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereNotLike("name", "%foo%")
	sql := b.GetGrammar().CompileSelect(b)

	if !strings.Contains(sql, "`name` not like ?") {
		t.Errorf("expected 'not like ?'; got %q", sql)
	}
}

func TestClickHouseWhereRaw(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereRaw("length(name) > ?", 10)
	sql := b.GetGrammar().CompileSelect(b)

	if !strings.Contains(sql, "length(name) > ?") {
		t.Errorf("expected raw fragment; got %q", sql)
	}
}

func TestClickHouseWhereJsonContains(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereJsonContains("data", "needle")
	sql := b.GetGrammar().CompileSelect(b)

	if !strings.Contains(sql, "position(`data`, ?) > 0") {
		t.Errorf("expected position()/>0; got %q", sql)
	}
}

func TestClickHouseWhereJsonLength(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereJsonLength("tags", ">", 0)
	sql := b.GetGrammar().CompileSelect(b)

	if !strings.Contains(sql, "length(`tags`) > ?") {
		t.Errorf("expected length()/?; got %q", sql)
	}
}

func TestClickHouseWhereFullText(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereFullText([]string{"title", "body"}, "needle")
	sql := b.GetGrammar().CompileSelect(b)

	if !strings.Contains(sql, "match(`title`, `body`, ?)") {
		t.Errorf("expected match() with columns; got %q", sql)
	}
}

func TestClickHouseWhereBetweenColumns(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereBetweenColumns("age", [2]string{"min_age", "max_age"})
	sql := b.GetGrammar().CompileSelect(b)

	if !strings.Contains(sql, "`age` between `min_age` and `max_age`") {
		t.Errorf("expected column-between form; got %q", sql)
	}
}

func TestClickHouseWhereColumn(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().WhereColumn("a", "=", "b")
	sql := b.GetGrammar().CompileSelect(b)

	if !strings.Contains(sql, "`a` = `b`") {
		t.Errorf("expected column-vs-column compare; got %q", sql)
	}
}

func TestClickHouseUnion(t *testing.T) {
	t.Parallel()
	left := newClickHouseBuilder().Where("type", "=", "a")
	right := newClickHouseBuilder().Where("type", "=", "b")
	left.Union(right)
	sql := left.GetGrammar().CompileSelect(left)

	if !strings.Contains(sql, "union") {
		t.Errorf("expected 'union' keyword; got %q", sql)
	}
}

func TestClickHouseInsertOrIgnoreSurfacesUnsupported(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	sql, err := g.CompileInsertOrIgnore(nil, nil)

	if sql != "" {
		t.Errorf("expected empty SQL, got %q", sql)
	}

	if !errors.Is(err, database.ErrInsertOrIgnoreNotSupported) {
		t.Errorf("expected ErrInsertOrIgnoreNotSupported, got %v", err)
	}
}

func TestClickHouseHavingNull(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().GroupBy("country").HavingNull("count")
	sql := b.GetGrammar().CompileSelect(b)

	if !strings.Contains(sql, "having `count` is null") {
		t.Errorf("expected 'having `count` is null'; got %q", sql)
	}
}

func TestClickHouseHavingBetween(t *testing.T) {
	t.Parallel()
	b := newClickHouseBuilder().GroupBy("country").HavingBetween("count", [2]any{10, 100})
	sql := b.GetGrammar().CompileSelect(b)

	if !strings.Contains(sql, "having `count` between ? and ?") {
		t.Errorf("expected having-between; got %q", sql)
	}
}

func TestClickHouseIsExpressionAndGetValue(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	expr := &rawExpr{v: "now()"}

	if !g.IsExpression(expr) {
		t.Errorf("expected IsExpression(expr) = true")
	}

	if g.IsExpression("plain") {
		t.Errorf("string should not be an Expression")
	}

	if got := g.GetValue(expr); got != "now()" {
		t.Errorf("GetValue = %q, want now()", got)
	}
}

func TestClickHouseWrapTableWithAlias(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()

	if got := g.WrapTable("events as e"); got != "`events` as `e`" {
		t.Errorf("WrapTable alias = %q, want backtick-aliased form", got)
	}
}

func TestClickHouseWrapValueExpression(t *testing.T) {
	t.Parallel()
	g := grammars.NewClickHouseGrammar()
	// Raw expressions should pass through unwrapped — reuse the package's
	// existing rawExpr helper from mysql_grammar_test.go.
	b := query.NewBuilder(nil, g, nil).From("events").Select(&rawExpr{v: "count(*)"})
	sql := g.CompileSelect(b)

	if !strings.Contains(sql, "count(*)") {
		t.Errorf("expected raw expression preserved; got %q", sql)
	}
}
