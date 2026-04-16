package grammars_test

import (
	"testing"

	"github.com/bedrock/packages/database/query"
	"github.com/bedrock/packages/database/query/grammars"
)

type mockConnection struct{}

func (c *mockConnection) Select(_ interface{}, _ string, _ ...any) ([]map[string]any, error) {
	return nil, nil
}
func (c *mockConnection) Insert(_ interface{}, _ string, _ ...any) (bool, error)   { return true, nil }
func (c *mockConnection) Update(_ interface{}, _ string, _ ...any) (int64, error)  { return 0, nil }
func (c *mockConnection) Delete(_ interface{}, _ string, _ ...any) (int64, error)  { return 0, nil }
func (c *mockConnection) Statement(_ interface{}, _ string, _ ...any) (bool, error) {
	return true, nil
}
func (c *mockConnection) AffectingStatement(_ interface{}, _ string, _ ...any) (int64, error) {
	return 0, nil
}
func (c *mockConnection) Raw(value string) interface{ GetValue() string } {
	return &rawExpr{value}
}
func (c *mockConnection) GetTablePrefix() string { return "" }

type rawExpr struct{ v string }

func (e *rawExpr) GetValue() string { return e.v }

func newMySQLBuilder() *query.Builder {
	g := grammars.NewMySQLGrammar()
	return query.NewBuilder(nil, g, nil).From("users")
}

func TestMySQLBasicSelect(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder()
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `users`"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLSelectColumns(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder().Select("name", "email")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select `name`, `email` from `users`"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLDistinct(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder().Distinct().Select("name")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select distinct `name` from `users`"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLWhereBasic(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder().Where("id", "=", 1)
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `users` where `id` = ?"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLWhereIn(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder().WhereIn("id", []any{1, 2, 3})
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `users` where `id` in (?, ?, ?)"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLWhereNull(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder().WhereNull("deleted_at")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `users` where `deleted_at` is null"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLWhereNotNull(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder().WhereNotNull("email")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `users` where `email` is not null"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLWhereBetween(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder().WhereBetween("age", [2]any{18, 65})
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `users` where `age` between ? and ?"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLOrderBy(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder().OrderBy("name").OrderByDesc("id")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `users` order by `name` asc, `id` desc"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLLimitOffset(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder().Limit(10).Offset(20)
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `users` limit 10 offset 20"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLGroupByHaving(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder().GroupBy("status").Having("count", ">", 5)
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `users` group by `status` having `count` > ?"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLJoin(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder().Join("orders", "users.id", "=", "orders.user_id")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `users` inner join `orders` on `users`.`id` = `orders`.`user_id`"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLLeftJoin(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder().LeftJoin("orders", "users.id", "=", "orders.user_id")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `users` left join `orders` on `users`.`id` = `orders`.`user_id`"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLCompileInsert(t *testing.T) {
	t.Parallel()
	g := grammars.NewMySQLGrammar()
	b := query.NewBuilder(nil, g, nil).From("users")

	sql := g.CompileInsert(b, []map[string]any{{"email": "test@test.com", "name": "Test"}})
	expected := "insert into `users` (`email`, `name`) values (?, ?)"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLCompileUpdate(t *testing.T) {
	t.Parallel()
	g := grammars.NewMySQLGrammar()
	b := query.NewBuilder(nil, g, nil).From("users").Where("id", "=", 1)

	sql := g.CompileUpdate(b, map[string]any{"name": "New Name"})
	expected := "update `users` set `name` = ? where `id` = ?"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLCompileDelete(t *testing.T) {
	t.Parallel()
	g := grammars.NewMySQLGrammar()
	b := query.NewBuilder(nil, g, nil).From("users").Where("id", "=", 1)

	sql := g.CompileDelete(b)
	expected := "delete from `users` where `id` = ?"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLCompileExists(t *testing.T) {
	t.Parallel()
	g := grammars.NewMySQLGrammar()
	b := query.NewBuilder(nil, g, nil).From("users").Where("id", "=", 1)

	sql := g.CompileExists(b)
	expected := "select exists(select * from `users` where `id` = ?) as `exists`"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLLockForUpdate(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder().LockForUpdate()
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `users` for update"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}

func TestMySQLOrWhere(t *testing.T) {
	t.Parallel()
	b := newMySQLBuilder().Where("name", "=", "John").OrWhere("name", "=", "Jane")
	sql := b.GetGrammar().CompileSelect(b)
	expected := "select * from `users` where `name` = ? or `name` = ?"
	if sql != expected {
		t.Fatalf("expected %q, got %q", expected, sql)
	}
}
