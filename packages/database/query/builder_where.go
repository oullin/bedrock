package query

import (
	"strings"
)

// Where adds a WHERE clause. It supports multiple calling patterns:
//   - Where("column", value)           → column = value
//   - Where("column", "op", value)     → column op value
//   - Where(func(b *Builder))          → nested WHERE group
//   - Where(map[string]any{...})       → multiple column = value
func (b *Builder) Where(args ...any) *Builder {
	return b.where("and", args...)
}

// OrWhere adds an OR WHERE clause.
func (b *Builder) OrWhere(args ...any) *Builder {
	return b.where("or", args...)
}

// WhereNot adds a WHERE NOT group.
func (b *Builder) WhereNot(args ...any) *Builder {
	return b.Where(func(q *Builder) {
		q.where("and", args...)
	})
}

// OrWhereNot adds an OR WHERE NOT group.
func (b *Builder) OrWhereNot(args ...any) *Builder {
	return b.OrWhere(func(q *Builder) {
		q.where("and", args...)
	})
}

func (b *Builder) where(boolean string, args ...any) *Builder {
	if len(args) == 0 {
		return b
	}

	switch v := args[0].(type) {
	case func(*Builder):
		return b.whereNested(boolean, v)
	case map[string]any:
		return b.whereMap(boolean, v)
	}

	switch len(args) {
	case 2:
		column, _ := args[0].(string)
		b.wheres = append(b.wheres, WhereClause{
			Type: WhereBasic, Column: column, Operator: "=", Value: args[1], Boolean: boolean,
		})
		b.AddBinding(BindingWhere, args[1])
	case 3:
		column, _ := args[0].(string)
		operator, _ := args[1].(string)
		b.wheres = append(b.wheres, WhereClause{
			Type: WhereBasic, Column: column, Operator: operator, Value: args[2], Boolean: boolean,
		})
		b.AddBinding(BindingWhere, args[2])
	}

	return b
}

func (b *Builder) whereNested(boolean string, fn func(*Builder)) *Builder {
	query := b.NewQuery().From(b.from)
	fn(query)

	if len(query.wheres) > 0 {
		b.wheres = append(b.wheres, WhereClause{
			Type: WhereNested, Query: query, Boolean: boolean,
		})
		b.AddBinding(BindingWhere, query.GetRawBindings()[BindingWhere]...)
	}

	return b
}

func (b *Builder) whereMap(boolean string, m map[string]any) *Builder {
	for col, val := range m {
		b.wheres = append(b.wheres, WhereClause{
			Type: WhereBasic, Column: col, Operator: "=", Value: val, Boolean: boolean,
		})
		b.AddBinding(BindingWhere, val)
	}

	return b
}

// WhereColumn adds a column-to-column comparison.
func (b *Builder) WhereColumn(first string, args ...any) *Builder {
	return b.whereColumn("and", first, args...)
}

// OrWhereColumn adds an OR column-to-column comparison.
func (b *Builder) OrWhereColumn(first string, args ...any) *Builder {
	return b.whereColumn("or", first, args...)
}

func (b *Builder) whereColumn(boolean, first string, args ...any) *Builder {
	operator := "="

	var second string

	if len(args) == 1 {
		second, _ = args[0].(string)
	} else if len(args) >= 2 {
		operator, _ = args[0].(string)
		second, _ = args[1].(string)
	}

	b.wheres = append(b.wheres, WhereClause{
		Type: WhereColumn, Column: first, Operator: operator, Value: second, Boolean: boolean,
	})

	return b
}

// WhereIn adds a WHERE IN clause.
func (b *Builder) WhereIn(column string, values []any) *Builder {
	return b.whereIn("and", column, values, false)
}

// OrWhereIn adds an OR WHERE IN clause.
func (b *Builder) OrWhereIn(column string, values []any) *Builder {
	return b.whereIn("or", column, values, false)
}

// WhereNotIn adds a WHERE NOT IN clause.
func (b *Builder) WhereNotIn(column string, values []any) *Builder {
	return b.whereIn("and", column, values, true)
}

// OrWhereNotIn adds an OR WHERE NOT IN clause.
func (b *Builder) OrWhereNotIn(column string, values []any) *Builder {
	return b.whereIn("or", column, values, true)
}

func (b *Builder) whereIn(boolean, column string, values []any, not bool) *Builder {
	t := WhereIn

	if not {
		t = WhereNotIn
	}

	b.wheres = append(b.wheres, WhereClause{
		Type: t, Column: column, Values: values, Boolean: boolean,
	})
	b.AddBinding(BindingWhere, values...)

	return b
}

// WhereIntegerInRaw adds a WHERE IN with raw integer values (no bindings).
func (b *Builder) WhereIntegerInRaw(column string, values []int64) *Builder {
	anyValues := make([]any, len(values))

	for i, v := range values {
		anyValues[i] = v
	}

	b.wheres = append(b.wheres, WhereClause{
		Type: WhereIn, Column: column, Values: anyValues, Boolean: "and",
	})

	return b
}

// WhereNull adds a WHERE IS NULL clause.
func (b *Builder) WhereNull(columns ...string) *Builder {
	for _, col := range columns {
		b.wheres = append(b.wheres, WhereClause{
			Type: WhereNull, Column: col, Boolean: "and",
		})
	}

	return b
}

// OrWhereNull adds an OR WHERE IS NULL clause.
func (b *Builder) OrWhereNull(column string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereNull, Column: column, Boolean: "or",
	})

	return b
}

// WhereNotNull adds a WHERE IS NOT NULL clause.
func (b *Builder) WhereNotNull(columns ...string) *Builder {
	for _, col := range columns {
		b.wheres = append(b.wheres, WhereClause{
			Type: WhereNotNull, Column: col, Boolean: "and",
		})
	}

	return b
}

// OrWhereNotNull adds an OR WHERE IS NOT NULL clause.
func (b *Builder) OrWhereNotNull(column string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereNotNull, Column: column, Boolean: "or",
	})

	return b
}

// WhereBetween adds a WHERE BETWEEN clause.
func (b *Builder) WhereBetween(column string, values [2]any) *Builder {
	return b.whereBetween("and", column, values, false)
}

// OrWhereBetween adds an OR WHERE BETWEEN clause.
func (b *Builder) OrWhereBetween(column string, values [2]any) *Builder {
	return b.whereBetween("or", column, values, false)
}

// WhereNotBetween adds a WHERE NOT BETWEEN clause.
func (b *Builder) WhereNotBetween(column string, values [2]any) *Builder {
	return b.whereBetween("and", column, values, true)
}

// OrWhereNotBetween adds an OR WHERE NOT BETWEEN clause.
func (b *Builder) OrWhereNotBetween(column string, values [2]any) *Builder {
	return b.whereBetween("or", column, values, true)
}

func (b *Builder) whereBetween(boolean, column string, values [2]any, not bool) *Builder {
	t := WhereBetween

	if not {
		t = WhereNotBetween
	}

	b.wheres = append(b.wheres, WhereClause{
		Type: t, Column: column, Values: []any{values[0], values[1]}, Boolean: boolean, Not: not,
	})
	b.AddBinding(BindingWhere, values[0], values[1])

	return b
}

// WhereBetweenColumns adds a WHERE column BETWEEN col1 AND col2.
func (b *Builder) WhereBetweenColumns(column string, columns [2]string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereBetweenColumns, Column: column, Columns: columns[:], Boolean: "and",
	})

	return b
}

// OrWhereBetweenColumns adds an OR WHERE column BETWEEN col1 AND col2.
func (b *Builder) OrWhereBetweenColumns(column string, columns [2]string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereBetweenColumns, Column: column, Columns: columns[:], Boolean: "or",
	})

	return b
}

// WhereDate adds a WHERE date(column) clause.
func (b *Builder) WhereDate(column, operator, value string) *Builder {
	return b.whereDateTime("and", WhereDate, column, operator, value)
}

// OrWhereDate adds an OR WHERE date(column) clause.
func (b *Builder) OrWhereDate(column, operator, value string) *Builder {
	return b.whereDateTime("or", WhereDate, column, operator, value)
}

// WhereTime adds a WHERE time(column) clause.
func (b *Builder) WhereTime(column, operator, value string) *Builder {
	return b.whereDateTime("and", WhereTime, column, operator, value)
}

// OrWhereTime adds an OR WHERE time(column) clause.
func (b *Builder) OrWhereTime(column, operator, value string) *Builder {
	return b.whereDateTime("or", WhereTime, column, operator, value)
}

// WhereDay adds a WHERE day(column) clause.
func (b *Builder) WhereDay(column, operator, value string) *Builder {
	return b.whereDateTime("and", WhereDay, column, operator, value)
}

// OrWhereDay adds an OR WHERE day(column) clause.
func (b *Builder) OrWhereDay(column, operator, value string) *Builder {
	return b.whereDateTime("or", WhereDay, column, operator, value)
}

// WhereMonth adds a WHERE month(column) clause.
func (b *Builder) WhereMonth(column, operator, value string) *Builder {
	return b.whereDateTime("and", WhereMonth, column, operator, value)
}

// OrWhereMonth adds an OR WHERE month(column) clause.
func (b *Builder) OrWhereMonth(column, operator, value string) *Builder {
	return b.whereDateTime("or", WhereMonth, column, operator, value)
}

// WhereYear adds a WHERE year(column) clause.
func (b *Builder) WhereYear(column, operator, value string) *Builder {
	return b.whereDateTime("and", WhereYear, column, operator, value)
}

// OrWhereYear adds an OR WHERE year(column) clause.
func (b *Builder) OrWhereYear(column, operator, value string) *Builder {
	return b.whereDateTime("or", WhereYear, column, operator, value)
}

func (b *Builder) whereDateTime(boolean string, whereType WhereType, column, operator, value string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: whereType, Column: column, Operator: operator, Value: value, Boolean: boolean,
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// WhereRaw adds a raw WHERE clause.
func (b *Builder) WhereRaw(sql string, bindings ...any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereRaw, SQL: sql, Boolean: "and",
	})
	b.AddBinding(BindingWhere, bindings...)

	return b
}

// OrWhereRaw adds a raw OR WHERE clause.
func (b *Builder) OrWhereRaw(sql string, bindings ...any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereRaw, SQL: sql, Boolean: "or",
	})
	b.AddBinding(BindingWhere, bindings...)

	return b
}

// WhereExists adds a WHERE EXISTS clause.
func (b *Builder) WhereExists(query any) *Builder {
	return b.whereExists("and", query, false)
}

// OrWhereExists adds an OR WHERE EXISTS clause.
func (b *Builder) OrWhereExists(query any) *Builder {
	return b.whereExists("or", query, false)
}

// WhereNotExists adds a WHERE NOT EXISTS clause.
func (b *Builder) WhereNotExists(query any) *Builder {
	return b.whereExists("and", query, true)
}

// OrWhereNotExists adds an OR WHERE NOT EXISTS clause.
func (b *Builder) OrWhereNotExists(query any) *Builder {
	return b.whereExists("or", query, true)
}

func (b *Builder) whereExists(boolean string, query any, not bool) *Builder {
	t := WhereExists

	if not {
		t = WhereNotExists
	}

	switch q := query.(type) {
	case func(*Builder):
		sub := b.NewQuery()
		q(sub)
		b.wheres = append(b.wheres, WhereClause{
			Type: t, Query: sub, Boolean: boolean,
		})
		b.AddBinding(BindingWhere, sub.GetBindings()...)
	case *Builder:
		b.wheres = append(b.wheres, WhereClause{
			Type: t, Query: q, Boolean: boolean,
		})
		b.AddBinding(BindingWhere, q.GetBindings()...)
	}

	return b
}

// WhereSub adds a WHERE with a subquery as the value.
func (b *Builder) WhereSub(column, operator string, query any) *Builder {
	switch q := query.(type) {
	case func(*Builder):
		sub := b.NewQuery()
		q(sub)
		b.wheres = append(b.wheres, WhereClause{
			Type: WhereSub, Column: column, Operator: operator, Query: sub, Boolean: "and",
		})
		b.AddBinding(BindingWhere, sub.GetBindings()...)
	case *Builder:
		b.wheres = append(b.wheres, WhereClause{
			Type: WhereSub, Column: column, Operator: operator, Query: q, Boolean: "and",
		})
		b.AddBinding(BindingWhere, q.GetBindings()...)
	}

	return b
}

// WhereLike adds a WHERE LIKE clause.
func (b *Builder) WhereLike(column string, value any, caseSensitive ...bool) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereLike, Column: column, Value: value, Boolean: "and",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// OrWhereLike adds an OR WHERE LIKE clause.
func (b *Builder) OrWhereLike(column string, value any, caseSensitive ...bool) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereLike, Column: column, Value: value, Boolean: "or",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// WhereNotLike adds a WHERE NOT LIKE clause.
func (b *Builder) WhereNotLike(column string, value any, caseSensitive ...bool) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereNotLike, Column: column, Value: value, Boolean: "and",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// OrWhereNotLike adds an OR WHERE NOT LIKE clause.
func (b *Builder) OrWhereNotLike(column string, value any, caseSensitive ...bool) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereNotLike, Column: column, Value: value, Boolean: "or",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// WhereJsonContains adds a WHERE JSON_CONTAINS clause.
func (b *Builder) WhereJsonContains(column string, value any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereJsonContains, Column: column, Value: value, Boolean: "and",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// OrWhereJsonContains adds an OR WHERE JSON_CONTAINS clause.
func (b *Builder) OrWhereJsonContains(column string, value any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereJsonContains, Column: column, Value: value, Boolean: "or",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// WhereJsonLength adds a WHERE JSON_LENGTH clause.
func (b *Builder) WhereJsonLength(column, operator string, value any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereJsonLength, Column: column, Operator: operator, Value: value, Boolean: "and",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// OrWhereJsonLength adds an OR WHERE JSON_LENGTH clause.
func (b *Builder) OrWhereJsonLength(column, operator string, value any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereJsonLength, Column: column, Operator: operator, Value: value, Boolean: "or",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// WhereFullText adds a WHERE MATCH ... AGAINST clause.
func (b *Builder) WhereFullText(columns []string, value string, options ...map[string]any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereFullText, Columns: columns, Value: value, Boolean: "and",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// OrWhereFullText adds an OR WHERE MATCH ... AGAINST clause.
func (b *Builder) OrWhereFullText(columns []string, value string, options ...map[string]any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereFullText, Columns: columns, Value: value, Boolean: "or",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// WhereAll adds WHERE conditions for ALL given columns matching a value.
func (b *Builder) WhereAll(columns []string, operator string, value any) *Builder {
	return b.Where(func(q *Builder) {
		for _, col := range columns {
			q.Where(col, operator, value)
		}
	})
}

// OrWhereAll adds OR WHERE conditions for ALL given columns.
func (b *Builder) OrWhereAll(columns []string, operator string, value any) *Builder {
	return b.OrWhere(func(q *Builder) {
		for _, col := range columns {
			q.Where(col, operator, value)
		}
	})
}

// WhereAny adds WHERE conditions for ANY of the given columns.
func (b *Builder) WhereAny(columns []string, operator string, value any) *Builder {
	return b.Where(func(q *Builder) {
		for _, col := range columns {
			q.OrWhere(col, operator, value)
		}
	})
}

// OrWhereAny adds OR WHERE conditions for ANY of the given columns.
func (b *Builder) OrWhereAny(columns []string, operator string, value any) *Builder {
	return b.OrWhere(func(q *Builder) {
		for _, col := range columns {
			q.OrWhere(col, operator, value)
		}
	})
}

// WhereNone adds WHERE conditions such that NONE of the given columns match.
func (b *Builder) WhereNone(columns []string, operator string, value any) *Builder {
	return b.Where(func(q *Builder) {
		for _, col := range columns {
			q.Where(col, negateOperator(operator), value)
		}
	})
}

// WhereRowValues adds a WHERE with row value comparisons.
func (b *Builder) WhereRowValues(columns []string, operator string, values []any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereRowValues, Columns: columns, Operator: operator, Values: values, Boolean: "and",
	})
	b.AddBinding(BindingWhere, values...)

	return b
}

// OrWhereRowValues adds an OR WHERE with row value comparisons.
func (b *Builder) OrWhereRowValues(columns []string, operator string, values []any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereRowValues, Columns: columns, Operator: operator, Values: values, Boolean: "or",
	})
	b.AddBinding(BindingWhere, values...)

	return b
}

// OrWhereNone adds OR WHERE conditions such that NONE of the given columns match.
func (b *Builder) OrWhereNone(columns []string, operator string, value any) *Builder {
	return b.OrWhere(func(q *Builder) {
		for _, col := range columns {
			q.Where(col, negateOperator(operator), value)
		}
	})
}

// WhereNotBetweenColumns adds a WHERE column NOT BETWEEN col1 AND col2.
func (b *Builder) WhereNotBetweenColumns(column string, columns [2]string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereBetweenColumns, Column: column, Columns: columns[:], Boolean: "and", Not: true,
	})

	return b
}

// OrWhereNotBetweenColumns adds an OR WHERE column NOT BETWEEN col1 AND col2.
func (b *Builder) OrWhereNotBetweenColumns(column string, columns [2]string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereBetweenColumns, Column: column, Columns: columns[:], Boolean: "or", Not: true,
	})

	return b
}

// WhereValueBetween adds a WHERE ? BETWEEN col1 AND col2.
func (b *Builder) WhereValueBetween(value any, columns [2]string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereBetweenColumns, Value: value, Columns: columns[:], Boolean: "and",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// OrWhereValueBetween adds an OR WHERE ? BETWEEN col1 AND col2.
func (b *Builder) OrWhereValueBetween(value any, columns [2]string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereBetweenColumns, Value: value, Columns: columns[:], Boolean: "or",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// WhereValueNotBetween adds a WHERE ? NOT BETWEEN col1 AND col2.
func (b *Builder) WhereValueNotBetween(value any, columns [2]string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereBetweenColumns, Value: value, Columns: columns[:], Boolean: "and", Not: true,
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// OrWhereValueNotBetween adds an OR WHERE ? NOT BETWEEN col1 AND col2.
func (b *Builder) OrWhereValueNotBetween(value any, columns [2]string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereBetweenColumns, Value: value, Columns: columns[:], Boolean: "or", Not: true,
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// OrWhereIntegerInRaw adds an OR WHERE IN with raw integer values.
func (b *Builder) OrWhereIntegerInRaw(column string, values []int64) *Builder {
	anyValues := make([]any, len(values))

	for i, v := range values {
		anyValues[i] = v
	}

	b.wheres = append(b.wheres, WhereClause{
		Type: WhereIn, Column: column, Values: anyValues, Boolean: "or",
	})

	return b
}

// WhereIntegerNotInRaw adds a WHERE NOT IN with raw integer values.
func (b *Builder) WhereIntegerNotInRaw(column string, values []int64) *Builder {
	anyValues := make([]any, len(values))

	for i, v := range values {
		anyValues[i] = v
	}

	b.wheres = append(b.wheres, WhereClause{
		Type: WhereNotIn, Column: column, Values: anyValues, Boolean: "and",
	})

	return b
}

// OrWhereIntegerNotInRaw adds an OR WHERE NOT IN with raw integer values.
func (b *Builder) OrWhereIntegerNotInRaw(column string, values []int64) *Builder {
	anyValues := make([]any, len(values))

	for i, v := range values {
		anyValues[i] = v
	}

	b.wheres = append(b.wheres, WhereClause{
		Type: WhereNotIn, Column: column, Values: anyValues, Boolean: "or",
	})

	return b
}

// WhereJsonDoesntContain adds a WHERE NOT JSON_CONTAINS clause.
func (b *Builder) WhereJsonDoesntContain(column string, value any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereJsonContains, Column: column, Value: value, Boolean: "and", Not: true,
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// OrWhereJsonDoesntContain adds an OR WHERE NOT JSON_CONTAINS clause.
func (b *Builder) OrWhereJsonDoesntContain(column string, value any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereJsonContains, Column: column, Value: value, Boolean: "or", Not: true,
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// WhereJsonOverlaps adds a WHERE JSON_OVERLAPS clause.
func (b *Builder) WhereJsonOverlaps(column string, value any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereJsonContains, Column: column, Value: value, Boolean: "and",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// OrWhereJsonOverlaps adds an OR WHERE JSON_OVERLAPS clause.
func (b *Builder) OrWhereJsonOverlaps(column string, value any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereJsonContains, Column: column, Value: value, Boolean: "or",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// WhereJsonDoesntOverlap adds a WHERE NOT JSON_OVERLAPS clause.
func (b *Builder) WhereJsonDoesntOverlap(column string, value any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereJsonContains, Column: column, Value: value, Boolean: "and", Not: true,
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// OrWhereJsonDoesntOverlap adds an OR WHERE NOT JSON_OVERLAPS clause.
func (b *Builder) OrWhereJsonDoesntOverlap(column string, value any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereJsonContains, Column: column, Value: value, Boolean: "or", Not: true,
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// WhereJsonContainsKey adds a WHERE JSON_CONTAINS_KEY clause.
func (b *Builder) WhereJsonContainsKey(column string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereJsonContains, Column: column, Boolean: "and",
	})

	return b
}

// OrWhereJsonContainsKey adds an OR WHERE JSON_CONTAINS_KEY clause.
func (b *Builder) OrWhereJsonContainsKey(column string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereJsonContains, Column: column, Boolean: "or",
	})

	return b
}

// WhereJsonDoesntContainKey adds a WHERE NOT JSON_CONTAINS_KEY clause.
func (b *Builder) WhereJsonDoesntContainKey(column string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereJsonContains, Column: column, Boolean: "and", Not: true,
	})

	return b
}

// OrWhereJsonDoesntContainKey adds an OR WHERE NOT JSON_CONTAINS_KEY clause.
func (b *Builder) OrWhereJsonDoesntContainKey(column string) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereJsonContains, Column: column, Boolean: "or", Not: true,
	})

	return b
}

// WhereNullSafeEquals adds a WHERE column <=> value clause (MySQL null-safe).
func (b *Builder) WhereNullSafeEquals(column string, value any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereBasic, Column: column, Operator: "<=>", Value: value, Boolean: "and",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

// OrWhereNullSafeEquals adds an OR WHERE column <=> value clause.
func (b *Builder) OrWhereNullSafeEquals(column string, value any) *Builder {
	b.wheres = append(b.wheres, WhereClause{
		Type: WhereBasic, Column: column, Operator: "<=>", Value: value, Boolean: "or",
	})
	b.AddBinding(BindingWhere, value)

	return b
}

func negateOperator(op string) string {
	switch strings.ToLower(op) {
	case "=":
		return "!="
	case "!=", "<>":
		return "="
	case "<":
		return ">="
	case ">":
		return "<="
	case "<=":
		return ">"
	case ">=":
		return "<"
	case "like":
		return "not like"
	case "not like":
		return "like"
	default:
		return op
	}
}
