package query

// Union adds a UNION query.
func (b *Builder) Union(query *Builder) *Builder {
	b.unions = append(b.unions, UnionClause{Query: query, All: false})
	b.AddBinding(BindingUnion, query.GetBindings()...)
	return b
}

// UnionAll adds a UNION ALL query.
func (b *Builder) UnionAll(query *Builder) *Builder {
	b.unions = append(b.unions, UnionClause{Query: query, All: true})
	b.AddBinding(BindingUnion, query.GetBindings()...)
	return b
}

// UnionOrderBy adds an ORDER BY clause to the union query.
func (b *Builder) UnionOrderBy(column, direction string) *Builder {
	b.unionOrders = append(b.unionOrders, OrderClause{Column: column, Direction: direction})
	return b
}

// UnionLimit sets the limit on the union query.
func (b *Builder) UnionLimit(n int) *Builder {
	b.unionLimit = &n
	return b
}

// UnionOffset sets the offset on the union query.
func (b *Builder) UnionOffset(n int) *Builder {
	b.unionOffset = &n
	return b
}
