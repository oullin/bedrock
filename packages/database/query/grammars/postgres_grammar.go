package grammars

import (
	"fmt"
	"strings"

	dbcontract "github.com/bedrock/packages/contracts/database"
	"github.com/bedrock/packages/database/query"
)

// PostgresGrammar compiles query builder instances into PostgreSQL SQL.
type PostgresGrammar struct {
	tablePrefix  string
	paramCounter int
}

var _ query.Grammar = (*PostgresGrammar)(nil)

// NewPostgresGrammar creates a new PostgreSQL grammar.
func NewPostgresGrammar() *PostgresGrammar {
	return &PostgresGrammar{}
}

func (g *PostgresGrammar) resetParams() { g.paramCounter = 0 }
func (g *PostgresGrammar) nextParam() string {
	g.paramCounter++

	return fmt.Sprintf("$%d", g.paramCounter)
}

func (g *PostgresGrammar) CompileSelect(b *query.Builder) string {
	g.resetParams()

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
		g.compileLock(b),
	})

	if len(b.GetUnions()) > 0 {
		sql = "(" + sql + ") " + g.compileUnions(b)
	}

	return sql
}

func (g *PostgresGrammar) compileAggregate(b *query.Builder, agg *AggregateClause) string {
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

func (g *PostgresGrammar) compileColumns(b *query.Builder) string {
	cols := b.GetColumns()

	if len(cols) == 0 {
		cols = []any{"*"}
	}

	sel := "select "

	if b.IsDistinct() {
		sel = "select distinct "
	}

	var parts []string

	for _, col := range cols {
		parts = append(parts, g.wrapValue(col))
	}

	return sel + strings.Join(parts, ", ")
}

func (g *PostgresGrammar) compileFrom(b *query.Builder) string {
	from := b.GetFrom()

	if from == "" {
		return ""
	}

	return "from " + g.WrapTable(from)
}

func (g *PostgresGrammar) compileJoins(b *query.Builder) string {
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

func (g *PostgresGrammar) compileJoinConditions(j *query.JoinClause) string {
	var parts []string

	for i, c := range j.Clauses {
		cond := ""

		if c.Where {
			cond = g.Wrap(c.First) + " " + c.Operator + " " + g.nextParam()
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

func (g *PostgresGrammar) compileWheres(b *query.Builder) string {
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

func (g *PostgresGrammar) compileWhere(b *query.Builder, w query.WhereClause) string {
	switch w.Type {
	case query.WhereBasic:
		return g.Wrap(w.Column) + " " + w.Operator + " " + g.nextParam()
	case query.WhereColumn:
		second, _ := w.Value.(string)

		return g.Wrap(w.Column) + " " + w.Operator + " " + g.Wrap(second)
	case query.WhereIn:
		params := g.nParams(len(w.Values))

		return g.Wrap(w.Column) + " in (" + params + ")"
	case query.WhereNotIn:
		params := g.nParams(len(w.Values))

		return g.Wrap(w.Column) + " not in (" + params + ")"
	case query.WhereNull:
		return g.Wrap(w.Column) + " is null"
	case query.WhereNotNull:
		return g.Wrap(w.Column) + " is not null"
	case query.WhereBetween:
		return g.Wrap(w.Column) + " between " + g.nextParam() + " and " + g.nextParam()
	case query.WhereNotBetween:
		return g.Wrap(w.Column) + " not between " + g.nextParam() + " and " + g.nextParam()
	case query.WhereBetweenColumns:
		if len(w.Columns) >= 2 {
			return g.Wrap(w.Column) + " between " + g.Wrap(w.Columns[0]) + " and " + g.Wrap(w.Columns[1])
		}

		return ""
	case query.WhereDate:
		return g.Wrap(w.Column) + "::date " + w.Operator + " " + g.nextParam()
	case query.WhereTime:
		return g.Wrap(w.Column) + "::time " + w.Operator + " " + g.nextParam()
	case query.WhereDay:
		return "extract(day from " + g.Wrap(w.Column) + ") " + w.Operator + " " + g.nextParam()
	case query.WhereMonth:
		return "extract(month from " + g.Wrap(w.Column) + ") " + w.Operator + " " + g.nextParam()
	case query.WhereYear:
		return "extract(year from " + g.Wrap(w.Column) + ") " + w.Operator + " " + g.nextParam()
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
		return g.Wrap(w.Column) + " like " + g.nextParam()
	case query.WhereNotLike:
		return g.Wrap(w.Column) + " not like " + g.nextParam()
	case query.WhereJsonContains:
		return g.Wrap(w.Column) + " @> " + g.nextParam()
	case query.WhereJsonLength:
		return "jsonb_array_length(" + g.Wrap(w.Column) + ") " + w.Operator + " " + g.nextParam()
	case query.WhereFullText:
		cols := g.Columnize(w.Columns)

		return "to_tsvector(" + cols + ") @@ plainto_tsquery(" + g.nextParam() + ")"
	case query.WhereRowValues:
		cols := g.Columnize(w.Columns)
		params := g.nParams(len(w.Values))

		return "(" + cols + ") " + w.Operator + " (" + params + ")"
	default:
		return ""
	}
}

func (g *PostgresGrammar) compileGroups(b *query.Builder) string {
	groups := b.GetGroups()

	if len(groups) == 0 {
		return ""
	}

	return "group by " + g.Columnize(groups)
}

func (g *PostgresGrammar) compileHavings(b *query.Builder) string {
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

func (g *PostgresGrammar) compileHaving(h query.HavingClause) string {
	switch h.Type {
	case "Basic":
		return g.Wrap(h.Column) + " " + h.Operator + " " + g.nextParam()
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

		return g.Wrap(h.Column) + " " + not + "between " + g.nextParam() + " and " + g.nextParam()
	default:
		return ""
	}
}

func (g *PostgresGrammar) compileOrders(b *query.Builder) string {
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

func (g *PostgresGrammar) compileLimit(b *query.Builder) string {
	if b.GetLimit() < 0 {
		return ""
	}

	return fmt.Sprintf("limit %d", b.GetLimit())
}

func (g *PostgresGrammar) compileOffset(b *query.Builder) string {
	if b.GetOffset() < 0 {
		return ""
	}

	return fmt.Sprintf("offset %d", b.GetOffset())
}

func (g *PostgresGrammar) compileUnions(b *query.Builder) string {
	var parts []string

	for _, u := range b.GetUnions() {
		keyword := "union"

		if u.All {
			keyword = "union all"
		}

		parts = append(parts, keyword+" ("+g.CompileSelect(u.Query)+")")
	}

	return strings.Join(parts, " ")
}

func (g *PostgresGrammar) compileLock(b *query.Builder) string {
	lock := b.GetLock()

	if lock == nil {
		return ""
	}

	switch v := lock.(type) {
	case string:
		return v
	case bool:
		if v {
			return "for update"
		}

		return "for share"
	default:
		return ""
	}
}

func (g *PostgresGrammar) CompileExists(b *query.Builder) string {
	return "select exists(" + g.CompileSelect(b) + ") as \"exists\""
}

func (g *PostgresGrammar) CompileInsert(b *query.Builder, values []map[string]any) string {
	g.resetParams()
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

func (g *PostgresGrammar) CompileInsertOrIgnore(b *query.Builder, values []map[string]any) string {
	return g.CompileInsert(b, values) + " on conflict do nothing"
}

func (g *PostgresGrammar) CompileInsertGetId(b *query.Builder, values map[string]any, sequence string) string {
	return g.CompileInsert(b, []map[string]any{values}) + " returning " + g.Wrap(sequence)
}

func (g *PostgresGrammar) CompileInsertUsing(b *query.Builder, columns []string, sql string) string {
	table := g.WrapTable(b.GetFrom())

	return "insert into " + table + " (" + g.Columnize(columns) + ") " + sql
}

func (g *PostgresGrammar) CompileUpdate(b *query.Builder, values map[string]any) string {
	g.resetParams()
	table := g.WrapTable(b.GetFrom())
	keys := sortedKeys(values)

	var sets []string

	for _, k := range keys {
		val := values[k]

		if isExpression(val) {
			sets = append(sets, g.Wrap(k)+" = "+getExprValue(val))
		} else {
			sets = append(sets, g.Wrap(k)+" = "+g.nextParam())
		}
	}

	sql := "update " + table + " set " + strings.Join(sets, ", ")
	wheres := g.compileWheres(b)

	if wheres != "" {
		sql += " " + wheres
	}

	return sql
}

func (g *PostgresGrammar) CompileUpsert(b *query.Builder, values []map[string]any, uniqueBy []string, update []string) string {
	sql := g.CompileInsert(b, values)
	sql += " on conflict (" + g.Columnize(uniqueBy) + ") do update set "

	var sets []string

	for _, col := range update {
		sets = append(sets, g.Wrap(col)+" = "+g.Wrap("excluded."+col))
	}

	return sql + strings.Join(sets, ", ")
}

func (g *PostgresGrammar) CompileDelete(b *query.Builder) string {
	g.resetParams()
	table := g.WrapTable(b.GetFrom())
	sql := "delete from " + table
	wheres := g.compileWheres(b)

	if wheres != "" {
		sql += " " + wheres
	}

	return sql
}

func (g *PostgresGrammar) CompileTruncate(b *query.Builder) map[string]string {
	return map[string]string{
		"truncate " + g.WrapTable(b.GetFrom()) + " restart identity cascade": "",
	}
}

func (g *PostgresGrammar) CompileRandom(seed string) string {
	return "RANDOM()"
}

func (g *PostgresGrammar) Wrap(value string) string {
	if value == "*" {
		return value
	}

	if strings.Contains(value, " as ") {
		parts := strings.SplitN(value, " as ", 2)

		return g.wrapSegments(parts[0]) + " as " + g.wrapSingle(strings.TrimSpace(parts[1]))
	}

	return g.wrapSegments(value)
}

func (g *PostgresGrammar) wrapSegments(value string) string {
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

func (g *PostgresGrammar) wrapSingle(value string) string {
	if value == "*" {
		return value
	}

	return "\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\""
}

func (g *PostgresGrammar) WrapTable(table string) string {
	if strings.Contains(table, " as ") {
		parts := strings.SplitN(table, " as ", 2)

		return g.wrapSingle(g.tablePrefix+strings.TrimSpace(parts[0])) + " as " + g.wrapSingle(strings.TrimSpace(parts[1]))
	}

	if strings.Contains(table, "(") {
		return table
	}

	return g.wrapSingle(g.tablePrefix + table)
}

func (g *PostgresGrammar) Columnize(columns []string) string {
	wrapped := make([]string, len(columns))

	for i, col := range columns {
		wrapped[i] = g.Wrap(col)
	}

	return strings.Join(wrapped, ", ")
}

func (g *PostgresGrammar) Parameterize(values []any) string {
	params := make([]string, len(values))

	for i, v := range values {
		params[i] = g.Parameter(v)
	}

	return strings.Join(params, ", ")
}

func (g *PostgresGrammar) Parameter(value any) string {
	if isExpression(value) {
		return getExprValue(value)
	}

	return g.nextParam()
}

func (g *PostgresGrammar) GetTablePrefix() string         { return g.tablePrefix }
func (g *PostgresGrammar) SetTablePrefix(p string)        { g.tablePrefix = p }
func (g *PostgresGrammar) IsExpression(v any) bool        { return isExpression(v) }
func (g *PostgresGrammar) GetValue(expression any) string { return getExprValue(expression) }

func (g *PostgresGrammar) wrapValue(v any) string {
	switch val := v.(type) {
	case string:
		return g.Wrap(val)
	case dbcontract.Expression:
		return val.GetValue()
	default:
		return fmt.Sprintf("%v", val)
	}
}

func (g *PostgresGrammar) concatenate(parts []string) string {
	var nonEmpty []string

	for _, p := range parts {
		if p != "" {
			nonEmpty = append(nonEmpty, p)
		}
	}

	return strings.Join(nonEmpty, " ")
}

func (g *PostgresGrammar) nParams(n int) string {
	if n <= 0 {
		return ""
	}

	params := make([]string, n)

	for i := range params {
		params[i] = g.nextParam()
	}

	return strings.Join(params, ", ")
}
