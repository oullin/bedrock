package query

// GroupBy adds GROUP BY columns.
func (b *Builder) GroupBy(groups ...string) *Builder {
	b.groups = append(b.groups, groups...)

	return b
}

// GroupByRaw adds a raw GROUP BY expression.
func (b *Builder) GroupByRaw(sql string, bindings ...any) *Builder {
	b.groups = append(b.groups, sql)
	b.AddBinding(BindingGroupBy, bindings...)

	return b
}

// Having adds a HAVING clause.
func (b *Builder) Having(column string, args ...any) *Builder {
	return b.having("and", column, args...)
}

// OrHaving adds an OR HAVING clause.
func (b *Builder) OrHaving(column string, args ...any) *Builder {
	return b.having("or", column, args...)
}

func (b *Builder) having(boolean, column string, args ...any) *Builder {
	operator := "="

	var value any

	if len(args) == 1 {
		value = args[0]
	} else if len(args) >= 2 {
		operator, _ = args[0].(string)
		value = args[1]
	}

	b.havings = append(b.havings, HavingClause{
		Type: "Basic", Column: column, Operator: operator, Value: value, Boolean: boolean,
	})
	b.AddBinding(BindingHaving, value)

	return b
}

// HavingRaw adds a raw HAVING clause.
func (b *Builder) HavingRaw(sql string, bindings ...any) *Builder {
	b.havings = append(b.havings, HavingClause{
		Type: "Raw", SQL: sql, Boolean: "and",
	})
	b.AddBinding(BindingHaving, bindings...)

	return b
}

// OrHavingRaw adds a raw OR HAVING clause.
func (b *Builder) OrHavingRaw(sql string, bindings ...any) *Builder {
	b.havings = append(b.havings, HavingClause{
		Type: "Raw", SQL: sql, Boolean: "or",
	})
	b.AddBinding(BindingHaving, bindings...)

	return b
}

// HavingNull adds a HAVING IS NULL clause.
func (b *Builder) HavingNull(column string) *Builder {
	b.havings = append(b.havings, HavingClause{
		Type: "Null", Column: column, Boolean: "and",
	})

	return b
}

// OrHavingNull adds an OR HAVING IS NULL clause.
func (b *Builder) OrHavingNull(column string) *Builder {
	b.havings = append(b.havings, HavingClause{
		Type: "Null", Column: column, Boolean: "or",
	})

	return b
}

// HavingNotNull adds a HAVING IS NOT NULL clause.
func (b *Builder) HavingNotNull(column string) *Builder {
	b.havings = append(b.havings, HavingClause{
		Type: "NotNull", Column: column, Boolean: "and",
	})

	return b
}

// OrHavingNotNull adds an OR HAVING IS NOT NULL clause.
func (b *Builder) OrHavingNotNull(column string) *Builder {
	b.havings = append(b.havings, HavingClause{
		Type: "NotNull", Column: column, Boolean: "or",
	})

	return b
}

// HavingBetween adds a HAVING BETWEEN clause.
func (b *Builder) HavingBetween(column string, values [2]any) *Builder {
	b.havings = append(b.havings, HavingClause{
		Type: "Between", Column: column, Values: []any{values[0], values[1]}, Boolean: "and",
	})
	b.AddBinding(BindingHaving, values[0], values[1])

	return b
}

// HavingNotBetween adds a HAVING NOT BETWEEN clause.
func (b *Builder) HavingNotBetween(column string, values [2]any) *Builder {
	b.havings = append(b.havings, HavingClause{
		Type: "Between", Column: column, Values: []any{values[0], values[1]}, Boolean: "and", Not: true,
	})
	b.AddBinding(BindingHaving, values[0], values[1])

	return b
}

// OrHavingBetween adds an OR HAVING BETWEEN clause.
func (b *Builder) OrHavingBetween(column string, values [2]any) *Builder {
	b.havings = append(b.havings, HavingClause{
		Type: "Between", Column: column, Values: []any{values[0], values[1]}, Boolean: "or",
	})
	b.AddBinding(BindingHaving, values[0], values[1])

	return b
}

// OrHavingNotBetween adds an OR HAVING NOT BETWEEN clause.
func (b *Builder) OrHavingNotBetween(column string, values [2]any) *Builder {
	b.havings = append(b.havings, HavingClause{
		Type: "Between", Column: column, Values: []any{values[0], values[1]}, Boolean: "or", Not: true,
	})
	b.AddBinding(BindingHaving, values[0], values[1])

	return b
}

// HavingNested adds a nested HAVING group.
func (b *Builder) HavingNested(fn func(*Builder)) *Builder {
	q := b.NewQuery().From(b.from)
	fn(q)

	if len(q.havings) > 0 {
		b.havings = append(b.havings, HavingClause{
			Type: "Nested", Boolean: "and",
		})
		b.AddBinding(BindingHaving, q.GetRawBindings()[BindingHaving]...)
	}

	return b
}
