package query

// Select sets the columns to be selected.
func (b *Builder) Select(columns ...any) *Builder {
	b.columns = nil
	b.bindings[BindingSelect] = nil

	for _, col := range columns {
		switch v := col.(type) {
		case string:
			b.columns = append(b.columns, v)
		case *Builder:
			b.columns = append(b.columns, v)
		default:
			b.columns = append(b.columns, v)
		}
	}

	return b
}

// SelectRaw adds a raw expression to the select clause.
func (b *Builder) SelectRaw(expression string, bindings ...any) *Builder {
	b.AddSelect(b.connection.Raw(expression))
	b.AddBinding(BindingSelect, bindings...)

	return b
}

// SelectSub adds a subquery as a column.
func (b *Builder) SelectSub(query any, as string) *Builder {
	switch q := query.(type) {
	case *Builder:
		compiled := b.grammar.CompileSelect(q)
		b.columns = append(b.columns, b.connection.Raw("("+compiled+") as "+b.grammar.Wrap(as)))
		b.AddBinding(BindingSelect, q.GetBindings()...)
	case string:
		b.columns = append(b.columns, b.connection.Raw("("+q+") as "+b.grammar.Wrap(as)))
	}

	return b
}

// AddSelect adds columns to the existing select clause.
func (b *Builder) AddSelect(columns ...any) *Builder {
	for _, col := range columns {
		b.columns = append(b.columns, col)
	}

	return b
}

// Distinct marks the query as a DISTINCT query.
func (b *Builder) Distinct(columns ...string) *Builder {
	b.distinct = true

	if len(columns) > 0 {
		b.distinctColumns = columns
	}

	return b
}
