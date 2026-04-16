package query

// JoinClause represents a JOIN in the query.
type JoinClause struct {
	*Builder
	Type       JoinType
	Table      string
	Clauses    []JoinCondition
	ParentFrom string
}

// JoinCondition is a single ON condition in a JOIN.
type JoinCondition struct {
	First    string
	Operator string
	Second   any
	Boolean  string
	Where    bool // true if this is a WHERE condition (bound), not an ON (column)
	Nested   *JoinClause
}

// NewJoinClause creates a new JoinClause.
func NewJoinClause(parentBuilder *Builder, joinType JoinType, table string) *JoinClause {
	j := &JoinClause{
		Builder:    parentBuilder.NewQuery(),
		Type:       joinType,
		Table:      table,
		ParentFrom: parentBuilder.from,
	}
	return j
}

// On adds an ON condition to the join.
func (j *JoinClause) On(first, operator string, second ...any) *JoinClause {
	var secondVal any
	if len(second) > 0 {
		secondVal = second[0]
	}
	j.Clauses = append(j.Clauses, JoinCondition{
		First: first, Operator: operator, Second: secondVal, Boolean: "and",
	})
	return j
}

// OrOn adds an OR ON condition.
func (j *JoinClause) OrOn(first, operator string, second ...any) *JoinClause {
	var secondVal any
	if len(second) > 0 {
		secondVal = second[0]
	}
	j.Clauses = append(j.Clauses, JoinCondition{
		First: first, Operator: operator, Second: secondVal, Boolean: "or",
	})
	return j
}

// JoinWhere adds a WHERE-type condition (with binding) to the join.
func (j *JoinClause) JoinWhere(first, operator string, value any) *JoinClause {
	j.Clauses = append(j.Clauses, JoinCondition{
		First: first, Operator: operator, Second: value, Boolean: "and", Where: true,
	})
	j.AddBinding(BindingJoin, value)
	return j
}

// OrJoinWhere adds an OR WHERE-type condition to the join.
func (j *JoinClause) OrJoinWhere(first, operator string, value any) *JoinClause {
	j.Clauses = append(j.Clauses, JoinCondition{
		First: first, Operator: operator, Second: value, Boolean: "or", Where: true,
	})
	j.AddBinding(BindingJoin, value)
	return j
}

// Join adds an INNER JOIN to the query.
func (b *Builder) Join(table string, args ...any) *Builder {
	return b.join(JoinInner, table, args...)
}

// LeftJoin adds a LEFT JOIN to the query.
func (b *Builder) LeftJoin(table string, args ...any) *Builder {
	return b.join(JoinLeft, table, args...)
}

// RightJoin adds a RIGHT JOIN to the query.
func (b *Builder) RightJoin(table string, args ...any) *Builder {
	return b.join(JoinRight, table, args...)
}

// CrossJoin adds a CROSS JOIN to the query.
func (b *Builder) CrossJoin(table string, args ...any) *Builder {
	if len(args) == 0 {
		j := NewJoinClause(b, JoinCross, table)
		b.joins = append(b.joins, j)
		return b
	}
	return b.join(JoinCross, table, args...)
}

// JoinSub adds a JOIN to a subquery.
func (b *Builder) JoinSub(query any, as string, fn func(*JoinClause)) *Builder {
	return b.joinSub(JoinInner, query, as, fn)
}

// LeftJoinSub adds a LEFT JOIN to a subquery.
func (b *Builder) LeftJoinSub(query any, as string, fn func(*JoinClause)) *Builder {
	return b.joinSub(JoinLeft, query, as, fn)
}

// RightJoinSub adds a RIGHT JOIN to a subquery.
func (b *Builder) RightJoinSub(query any, as string, fn func(*JoinClause)) *Builder {
	return b.joinSub(JoinRight, query, as, fn)
}

// CrossJoinSub adds a CROSS JOIN to a subquery.
func (b *Builder) CrossJoinSub(query any, as string) *Builder {
	j := NewJoinClause(b, JoinCross, as)
	// Store subquery info on the join clause.
	switch q := query.(type) {
	case *Builder:
		j.Table = "(" + b.grammar.CompileSelect(q) + ") as " + b.grammar.Wrap(as)
		b.AddBinding(BindingJoin, q.GetBindings()...)
	case string:
		j.Table = "(" + q + ") as " + b.grammar.Wrap(as)
	}
	b.joins = append(b.joins, j)
	return b
}

// JoinLateral adds a LATERAL JOIN.
func (b *Builder) JoinLateral(query any, as string) *Builder {
	j := NewJoinClause(b, JoinLateral, as)
	switch q := query.(type) {
	case *Builder:
		j.Table = "lateral (" + b.grammar.CompileSelect(q) + ") as " + b.grammar.Wrap(as)
		b.AddBinding(BindingJoin, q.GetBindings()...)
	case string:
		j.Table = "lateral (" + q + ") as " + b.grammar.Wrap(as)
	}
	b.joins = append(b.joins, j)
	return b
}

// LeftJoinLateral adds a LEFT LATERAL JOIN.
func (b *Builder) LeftJoinLateral(query any, as string) *Builder {
	j := NewJoinClause(b, JoinLeft, as)
	switch q := query.(type) {
	case *Builder:
		j.Table = "lateral (" + b.grammar.CompileSelect(q) + ") as " + b.grammar.Wrap(as)
		b.AddBinding(BindingJoin, q.GetBindings()...)
	case string:
		j.Table = "lateral (" + q + ") as " + b.grammar.Wrap(as)
	}
	b.joins = append(b.joins, j)
	return b
}

// JoinWhere adds a JOIN with WHERE-type conditions.
func (b *Builder) JoinWhere(table, first, operator string, value any) *Builder {
	j := NewJoinClause(b, JoinInner, table)
	j.JoinWhere(first, operator, value)
	b.joins = append(b.joins, j)
	b.AddBinding(BindingJoin, j.GetRawBindings()[BindingJoin]...)
	return b
}

// LeftJoinWhere adds a LEFT JOIN with WHERE-type conditions.
func (b *Builder) LeftJoinWhere(table, first, operator string, value any) *Builder {
	j := NewJoinClause(b, JoinLeft, table)
	j.JoinWhere(first, operator, value)
	b.joins = append(b.joins, j)
	b.AddBinding(BindingJoin, j.GetRawBindings()[BindingJoin]...)
	return b
}

// RightJoinWhere adds a RIGHT JOIN with WHERE-type conditions.
func (b *Builder) RightJoinWhere(table, first, operator string, value any) *Builder {
	j := NewJoinClause(b, JoinRight, table)
	j.JoinWhere(first, operator, value)
	b.joins = append(b.joins, j)
	b.AddBinding(BindingJoin, j.GetRawBindings()[BindingJoin]...)
	return b
}

func (b *Builder) join(joinType JoinType, table string, args ...any) *Builder {
	j := NewJoinClause(b, joinType, table)

	if len(args) == 1 {
		if fn, ok := args[0].(func(*JoinClause)); ok {
			fn(j)
		}
	} else if len(args) >= 2 {
		first, _ := args[0].(string)
		if len(args) == 2 {
			second, _ := args[1].(string)
			j.On(first, "=", second)
		} else {
			operator, _ := args[1].(string)
			j.On(first, operator, args[2])
		}
	}

	b.joins = append(b.joins, j)
	b.AddBinding(BindingJoin, j.GetRawBindings()[BindingJoin]...)

	return b
}

func (b *Builder) joinSub(joinType JoinType, query any, as string, fn func(*JoinClause)) *Builder {
	j := NewJoinClause(b, joinType, as)

	switch q := query.(type) {
	case *Builder:
		j.Table = "(" + b.grammar.CompileSelect(q) + ") as " + b.grammar.Wrap(as)
		b.AddBinding(BindingJoin, q.GetBindings()...)
	case string:
		j.Table = "(" + q + ") as " + b.grammar.Wrap(as)
	}

	fn(j)
	b.joins = append(b.joins, j)
	b.AddBinding(BindingJoin, j.GetRawBindings()[BindingJoin]...)

	return b
}
