package grammars

import (
	"fmt"
	"strings"

	dbcontract "github.com/bedrock/packages/contracts/database"
	"github.com/bedrock/packages/database"
	"github.com/bedrock/packages/database/query"
)

// ClickHouseGrammar compiles query builder instances into ClickHouse SQL.
// ClickHouse uses positional `?` placeholders (like MySQL) and backtick-wrapped
// identifiers. UPDATE and DELETE are emitted as ALTER TABLE mutations.
type ClickHouseGrammar struct {
	tablePrefix string
}

var _ query.Grammar = (*ClickHouseGrammar)(nil)

// NewClickHouseGrammar creates a new ClickHouse grammar.
func NewClickHouseGrammar() *ClickHouseGrammar {
	return &ClickHouseGrammar{}
}

func (g *ClickHouseGrammar) CompileSelect(b *query.Builder) string {
	if agg := b.GetAggregate(); agg != nil {
		return g.compileAggregate(b, agg)
	}

	sql := g.concatenate([]string{
		g.compileColumns(b),
		g.compileFrom(b),
		g.compileJoins(b),
		g.compileWheres(b),
		g.compileGroups(b),
		g.compileHavings(b),
		g.compileOrders(b),
		g.compileLimit(b),
		g.compileOffset(b),
	})

	if len(b.GetUnions()) > 0 {
		sql = "(" + sql + ") " + g.compileUnions(b)
	}

	return sql
}

func (g *ClickHouseGrammar) compileAggregate(b *query.Builder, agg *query.AggregateClause) string {
	cols := "*"

	if len(agg.Columns) > 0 && agg.Columns[0] != "*" {
		cols = g.Columnize(agg.Columns)
	}

	if b.IsDistinct() && cols != "*" {
		cols = "distinct " + cols
	}

	sql := "select " + agg.Function + "(" + cols + ") as aggregate"
	sql += " " + g.compileFrom(b)
	sql += " " + g.compileJoins(b)
	sql += " " + g.compileWheres(b)
	sql += " " + g.compileGroups(b)
	sql += " " + g.compileHavings(b)

	return strings.TrimSpace(sql)
}

func (g *ClickHouseGrammar) compileColumns(b *query.Builder) string {
	cols := b.GetColumns()

	if len(cols) == 0 {
		cols = []any{"*"}
	}

	sel := "select "

	if b.IsDistinct() {
		sel = "select distinct "
	}

	parts := make([]string, 0, len(cols))

	for _, col := range cols {
		parts = append(parts, g.wrapValue(col))
	}

	return sel + strings.Join(parts, ", ")
}

func (g *ClickHouseGrammar) compileFrom(b *query.Builder) string {
	from := b.GetFrom()

	if from == "" {
		return ""
	}

	return "from " + g.WrapTable(from)
}

func (g *ClickHouseGrammar) compileJoins(b *query.Builder) string {
	joins := b.GetJoins()

	if len(joins) == 0 {
		return ""
	}

	var parts []string

	for _, j := range joins {
		table := g.WrapTable(j.Table)
		joinSQL := string(j.Type) + " join " + table

		if len(j.Clauses) > 0 {
			joinSQL += " on " + g.compileJoinConditions(j)
		}

		parts = append(parts, joinSQL)
	}

	return strings.Join(parts, " ")
}

func (g *ClickHouseGrammar) compileJoinConditions(j *query.JoinClause) string {
	var parts []string

	for i, c := range j.Clauses {
		cond := ""

		if c.Where {
			cond = g.Wrap(c.First) + " " + c.Operator + " ?"
		} else {
			cond = g.Wrap(c.First) + " " + c.Operator + " " + g.wrapValue(c.Second)
		}

		if i > 0 {
			cond = c.Boolean + " " + cond
		}

		parts = append(parts, cond)
	}

	return strings.Join(parts, " ")
}

func (g *ClickHouseGrammar) compileWheres(b *query.Builder) string {
	wheres := b.GetWheres()

	if len(wheres) == 0 {
		return ""
	}

	var parts []string

	for i, w := range wheres {
		compiled := g.compileWhere(b, w)

		if i == 0 {
			parts = append(parts, compiled)
		} else {
			parts = append(parts, w.Boolean+" "+compiled)
		}
	}

	return "where " + strings.Join(parts, " ")
}

func (g *ClickHouseGrammar) compileWhere(b *query.Builder, w query.WhereClause) string {
	switch w.Type {
	case query.WhereBasic:
		return g.Wrap(w.Column) + " " + w.Operator + " ?"
	case query.WhereColumn:
		second, _ := w.Value.(string)

		return g.Wrap(w.Column) + " " + w.Operator + " " + g.Wrap(second)
	case query.WhereIn:
		return g.Wrap(w.Column) + " in (" + g.Parameterize(w.Values) + ")"
	case query.WhereNotIn:
		return g.Wrap(w.Column) + " not in (" + g.Parameterize(w.Values) + ")"
	case query.WhereNull:
		return g.Wrap(w.Column) + " is null"
	case query.WhereNotNull:
		return g.Wrap(w.Column) + " is not null"
	case query.WhereBetween:
		return g.Wrap(w.Column) + " between ? and ?"
	case query.WhereNotBetween:
		return g.Wrap(w.Column) + " not between ? and ?"
	case query.WhereBetweenColumns:
		if len(w.Columns) >= 2 {
			return g.Wrap(w.Column) + " between " + g.Wrap(w.Columns[0]) + " and " + g.Wrap(w.Columns[1])
		}

		return ""
	case query.WhereDate:
		return "toDate(" + g.Wrap(w.Column) + ") " + w.Operator + " ?"
	case query.WhereTime:
		return "toTime(" + g.Wrap(w.Column) + ") " + w.Operator + " ?"
	case query.WhereDay:
		return "toDayOfMonth(" + g.Wrap(w.Column) + ") " + w.Operator + " ?"
	case query.WhereMonth:
		return "toMonth(" + g.Wrap(w.Column) + ") " + w.Operator + " ?"
	case query.WhereYear:
		return "toYear(" + g.Wrap(w.Column) + ") " + w.Operator + " ?"
	case query.WhereRaw:
		return w.SQL
	case query.WhereExists:
		return "exists (" + g.CompileSelect(w.Query) + ")"
	case query.WhereNotExists:
		return "not exists (" + g.CompileSelect(w.Query) + ")"
	case query.WhereNested:
		nested := g.compileWheres(w.Query)
		nested = strings.TrimPrefix(nested, "where ")

		return "(" + nested + ")"
	case query.WhereSub:
		return g.Wrap(w.Column) + " " + w.Operator + " (" + g.CompileSelect(w.Query) + ")"
	case query.WhereLike:
		return g.Wrap(w.Column) + " like ?"
	case query.WhereNotLike:
		return g.Wrap(w.Column) + " not like ?"
	case query.WhereJsonContains:
		return "position(" + g.Wrap(w.Column) + ", ?) > 0"
	case query.WhereJsonLength:
		return "length(" + g.Wrap(w.Column) + ") " + w.Operator + " ?"
	case query.WhereFullText:
		cols := g.Columnize(w.Columns)

		return "match(" + cols + ", ?)"
	case query.WhereRowValues:
		cols := g.Columnize(w.Columns)
		params := g.Parameterize(w.Values)

		return "(" + cols + ") " + w.Operator + " (" + params + ")"
	default:
		return ""
	}
}

func (g *ClickHouseGrammar) compileGroups(b *query.Builder) string {
	groups := b.GetGroups()

	if len(groups) == 0 {
		return ""
	}

	return "group by " + g.Columnize(groups)
}

func (g *ClickHouseGrammar) compileHavings(b *query.Builder) string {
	havings := b.GetHavings()

	if len(havings) == 0 {
		return ""
	}

	var parts []string

	for i, h := range havings {
		compiled := g.compileHaving(h)

		if i == 0 {
			parts = append(parts, compiled)
		} else {
			parts = append(parts, h.Boolean+" "+compiled)
		}
	}

	return "having " + strings.Join(parts, " ")
}

func (g *ClickHouseGrammar) compileHaving(h query.HavingClause) string {
	switch h.Type {
	case "Basic":
		return g.Wrap(h.Column) + " " + h.Operator + " ?"
	case "Raw":
		return h.SQL
	case "Null":
		return g.Wrap(h.Column) + " is null"
	case "NotNull":
		return g.Wrap(h.Column) + " is not null"
	case "Between":
		not := ""

		if h.Not {
			not = "not "
		}

		return g.Wrap(h.Column) + " " + not + "between ? and ?"
	default:
		return ""
	}
}

func (g *ClickHouseGrammar) compileOrders(b *query.Builder) string {
	orders := b.GetOrders()

	if len(orders) == 0 {
		return ""
	}

	var parts []string

	for _, o := range orders {
		if o.SQL != "" {
			parts = append(parts, o.SQL)
		} else {
			parts = append(parts, g.Wrap(o.Column)+" "+o.Direction)
		}
	}

	return "order by " + strings.Join(parts, ", ")
}

func (g *ClickHouseGrammar) compileLimit(b *query.Builder) string {
	if b.GetLimit() < 0 {
		return ""
	}

	return fmt.Sprintf("limit %d", b.GetLimit())
}

func (g *ClickHouseGrammar) compileOffset(b *query.Builder) string {
	if b.GetOffset() < 0 {
		return ""
	}

	return fmt.Sprintf("offset %d", b.GetOffset())
}

func (g *ClickHouseGrammar) compileUnions(b *query.Builder) string {
	var parts []string

	for _, u := range b.GetUnions() {
		keyword := "union distinct"

		if u.All {
			keyword = "union all"
		}

		parts = append(parts, keyword+" ("+g.CompileSelect(u.Query)+")")
	}

	return strings.Join(parts, " ")
}

func (g *ClickHouseGrammar) CompileExists(b *query.Builder) string {
	return "select count() > 0 as `exists` from (" + g.CompileSelect(b) + ")"
}

func (g *ClickHouseGrammar) CompileInsert(b *query.Builder, values []map[string]any) string {
	table := g.WrapTable(b.GetFrom())

	if len(values) == 0 {
		return "insert into " + table + " default values"
	}

	columns := sortedKeys(values[0])
	cols := g.Columnize(columns)

	var paramRows []string

	for range values {
		paramRows = append(paramRows, "("+g.nParams(len(columns))+")")
	}

	return "insert into " + table + " (" + cols + ") values " + strings.Join(paramRows, ", ")
}

func (g *ClickHouseGrammar) CompileInsertOrIgnore(_ *query.Builder, _ []map[string]any) string {
	return throwSQL("InsertOrIgnore is not supported; use ReplacingMergeTree at the table level")
}

// CompileInsertGetId returns database.ErrInsertGetIdNotSupported. ClickHouse
// has no auto-increment / last-insert-id semantics — callers must assign IDs
// explicitly and use Insert instead.
func (g *ClickHouseGrammar) CompileInsertGetId(_ *query.Builder, _ map[string]any, _ string) (string, error) {
	return "", database.ErrInsertGetIdNotSupported
}

func (g *ClickHouseGrammar) CompileInsertUsing(b *query.Builder, columns []string, sql string) string {
	table := g.WrapTable(b.GetFrom())

	return "insert into " + table + " (" + g.Columnize(columns) + ") " + sql
}

func (g *ClickHouseGrammar) CompileUpdate(b *query.Builder, values map[string]any) string {
	table := g.WrapTable(b.GetFrom())
	keys := sortedKeys(values)

	var sets []string

	for _, k := range keys {
		val := values[k]

		if isExpression(val) {
			sets = append(sets, g.Wrap(k)+" = "+getExprValue(val))
		} else {
			sets = append(sets, g.Wrap(k)+" = ?")
		}
	}

	sql := "alter table " + table + " update " + strings.Join(sets, ", ")
	wheres := g.compileWheres(b)

	if wheres == "" {
		sql += " where 1=1"
	} else {
		sql += " " + wheres
	}

	return sql
}

// CompileUpsert returns database.ErrUpsertNotSupported. ClickHouse exposes
// deduplication through ReplacingMergeTree at the table level rather than an
// INSERT ... ON CONFLICT path.
func (g *ClickHouseGrammar) CompileUpsert(_ *query.Builder, _ []map[string]any, _ []string, _ []string) (string, error) {
	return "", database.ErrUpsertNotSupported
}

func (g *ClickHouseGrammar) CompileDelete(b *query.Builder) string {
	table := g.WrapTable(b.GetFrom())
	sql := "alter table " + table + " delete"
	wheres := g.compileWheres(b)

	if wheres == "" {
		sql += " where 1=1"
	} else {
		sql += " " + wheres
	}

	return sql
}

func (g *ClickHouseGrammar) CompileTruncate(b *query.Builder) map[string]string {
	return map[string]string{
		"truncate table " + g.WrapTable(b.GetFrom()): "",
	}
}

func (g *ClickHouseGrammar) CompileRandom(_ string) string {
	return "rand()"
}

func (g *ClickHouseGrammar) Wrap(value string) string {
	if value == "*" {
		return value
	}

	if strings.Contains(value, " as ") {
		parts := strings.SplitN(value, " as ", 2)

		return g.wrapSegments(parts[0]) + " as " + g.wrapSingle(strings.TrimSpace(parts[1]))
	}

	return g.wrapSegments(value)
}

func (g *ClickHouseGrammar) wrapSegments(value string) string {
	segments := strings.Split(value, ".")
	wrapped := make([]string, len(segments))

	for i, seg := range segments {
		if i == 0 && len(segments) > 1 {
			wrapped[i] = g.WrapTable(seg)
		} else {
			wrapped[i] = g.wrapSingle(seg)
		}
	}

	return strings.Join(wrapped, ".")
}

func (g *ClickHouseGrammar) wrapSingle(value string) string {
	if value == "*" {
		return value
	}

	return "`" + strings.ReplaceAll(value, "`", "``") + "`"
}

func (g *ClickHouseGrammar) WrapTable(table string) string {
	if strings.Contains(table, " as ") {
		parts := strings.SplitN(table, " as ", 2)

		return g.wrapSingle(g.tablePrefix+strings.TrimSpace(parts[0])) + " as " + g.wrapSingle(strings.TrimSpace(parts[1]))
	}

	if strings.Contains(table, "(") {
		return table
	}

	return g.wrapSingle(g.tablePrefix + table)
}

func (g *ClickHouseGrammar) Columnize(columns []string) string {
	wrapped := make([]string, len(columns))

	for i, col := range columns {
		wrapped[i] = g.Wrap(col)
	}

	return strings.Join(wrapped, ", ")
}

func (g *ClickHouseGrammar) Parameterize(values []any) string {
	params := make([]string, len(values))

	for i, v := range values {
		params[i] = g.Parameter(v)
	}

	return strings.Join(params, ", ")
}

func (g *ClickHouseGrammar) Parameter(value any) string {
	if isExpression(value) {
		return getExprValue(value)
	}

	return "?"
}

func (g *ClickHouseGrammar) GetTablePrefix() string         { return g.tablePrefix }
func (g *ClickHouseGrammar) SetTablePrefix(p string)        { g.tablePrefix = p }
func (g *ClickHouseGrammar) IsExpression(v any) bool        { return isExpression(v) }
func (g *ClickHouseGrammar) GetValue(expression any) string { return getExprValue(expression) }

func (g *ClickHouseGrammar) wrapValue(v any) string {
	switch val := v.(type) {
	case string:
		return g.Wrap(val)
	case dbcontract.Expression:
		return val.GetValue()
	default:
		return fmt.Sprintf("%v", val)
	}
}

func (g *ClickHouseGrammar) concatenate(parts []string) string {
	var nonEmpty []string

	for _, p := range parts {
		if p != "" {
			nonEmpty = append(nonEmpty, p)
		}
	}

	return strings.Join(nonEmpty, " ")
}

func (g *ClickHouseGrammar) nParams(n int) string {
	if n <= 0 {
		return ""
	}

	params := make([]string, n)

	for i := range params {
		params[i] = "?"
	}

	return strings.Join(params, ", ")
}

// throwSQL emits a statement that raises a server-side error with the given
// message. Used for operations ClickHouse cannot support, so callers fail
// loudly rather than silently no-op.
func throwSQL(reason string) string {
	escaped := strings.ReplaceAll(reason, "'", "\\'")

	return fmt.Sprintf("select throwIf(1, 'clickhouse: %s')", escaped)
}
