package query

import "context"

// Update updates rows matching the query.
func (b *Builder) Update(ctx context.Context, values map[string]any) (int64, error) {
	compiled := b.grammar.CompileUpdate(b, values)
	var bindings []any
	for _, v := range sortedValues(values) {
		bindings = append(bindings, v)
	}
	bindings = append(bindings, b.GetRawBindings()[BindingWhere]...)
	bindings = append(bindings, b.GetRawBindings()[BindingJoin]...)

	return b.connection.Update(ctx, compiled, bindings...)
}

// UpdateOrInsert updates a row if it exists, or inserts it.
func (b *Builder) UpdateOrInsert(ctx context.Context, attributes, values map[string]any) (bool, error) {
	exists, err := b.Clone().Where(attributes).Exists(ctx)
	if err != nil {
		return false, err
	}

	if exists {
		if len(values) == 0 {
			return true, nil
		}
		_, err := b.Clone().Where(attributes).Limit(1).Update(ctx, values)
		return err == nil, err
	}

	merged := make(map[string]any, len(attributes)+len(values))
	for k, v := range attributes {
		merged[k] = v
	}
	for k, v := range values {
		merged[k] = v
	}

	return b.Insert(ctx, merged)
}

// Increment increments a column value by the given amount.
func (b *Builder) Increment(ctx context.Context, column string, amount ...any) (int64, error) {
	amt := any(1)
	if len(amount) > 0 {
		amt = amount[0]
	}
	values := map[string]any{
		column: b.connection.Raw(b.grammar.Wrap(column) + " + " + parameterize(amt)),
	}
	return b.Update(ctx, values)
}

// Decrement decrements a column value by the given amount.
func (b *Builder) Decrement(ctx context.Context, column string, amount ...any) (int64, error) {
	amt := any(1)
	if len(amount) > 0 {
		amt = amount[0]
	}
	values := map[string]any{
		column: b.connection.Raw(b.grammar.Wrap(column) + " - " + parameterize(amt)),
	}
	return b.Update(ctx, values)
}

// IncrementEach increments multiple columns.
func (b *Builder) IncrementEach(ctx context.Context, columns map[string]any, extra ...map[string]any) (int64, error) {
	values := make(map[string]any, len(columns))
	for col, amt := range columns {
		values[col] = b.connection.Raw(b.grammar.Wrap(col) + " + " + parameterize(amt))
	}
	if len(extra) > 0 {
		for k, v := range extra[0] {
			values[k] = v
		}
	}
	return b.Update(ctx, values)
}

// DecrementEach decrements multiple columns.
func (b *Builder) DecrementEach(ctx context.Context, columns map[string]any, extra ...map[string]any) (int64, error) {
	values := make(map[string]any, len(columns))
	for col, amt := range columns {
		values[col] = b.connection.Raw(b.grammar.Wrap(col) + " - " + parameterize(amt))
	}
	if len(extra) > 0 {
		for k, v := range extra[0] {
			values[k] = v
		}
	}
	return b.Update(ctx, values)
}

func parameterize(v any) string {
	return "?"
}
