package query

import "context"

// Insert inserts one or more rows.
func (b *Builder) Insert(ctx context.Context, values ...map[string]any) (bool, error) {
	if len(values) == 0 {
		return true, nil
	}

	compiled := b.grammar.CompileInsert(b, values)
	var bindings []any
	for _, row := range values {
		for _, v := range sortedValues(row) {
			bindings = append(bindings, v)
		}
	}

	return b.connection.Insert(ctx, compiled, bindings...)
}

// InsertOrIgnore inserts rows, ignoring duplicate key errors.
func (b *Builder) InsertOrIgnore(ctx context.Context, values ...map[string]any) (int64, error) {
	if len(values) == 0 {
		return 0, nil
	}

	compiled := b.grammar.CompileInsertOrIgnore(b, values)
	var bindings []any
	for _, row := range values {
		for _, v := range sortedValues(row) {
			bindings = append(bindings, v)
		}
	}

	return b.connection.AffectingStatement(ctx, compiled, bindings...)
}

// InsertGetId inserts a row and returns the auto-incrementing ID.
func (b *Builder) InsertGetId(ctx context.Context, values map[string]any, sequence ...string) (int64, error) {
	seq := "id"
	if len(sequence) > 0 {
		seq = sequence[0]
	}

	compiled := b.grammar.CompileInsertGetId(b, values, seq)
	var bindings []any
	for _, v := range sortedValues(values) {
		bindings = append(bindings, v)
	}

	return b.processor.ProcessInsertGetId(b, compiled, bindings, seq)
}

// InsertUsing inserts using a subquery.
func (b *Builder) InsertUsing(ctx context.Context, columns []string, query any) (int64, error) {
	var sql string
	var queryBindings []any

	switch q := query.(type) {
	case *Builder:
		sql = b.grammar.CompileSelect(q)
		queryBindings = q.GetBindings()
	case string:
		sql = q
	}

	compiled := b.grammar.CompileInsertUsing(b, columns, sql)

	return b.connection.AffectingStatement(ctx, compiled, queryBindings...)
}

// Upsert inserts or updates rows based on unique columns.
func (b *Builder) Upsert(ctx context.Context, values []map[string]any, uniqueBy []string, update []string) (int64, error) {
	if len(values) == 0 {
		return 0, nil
	}

	compiled := b.grammar.CompileUpsert(b, values, uniqueBy, update)
	var bindings []any
	for _, row := range values {
		for _, v := range sortedValues(row) {
			bindings = append(bindings, v)
		}
	}
	// Add update values bindings.
	for _, row := range values {
		for _, col := range update {
			if v, ok := row[col]; ok {
				bindings = append(bindings, v)
			}
		}
	}

	return b.connection.AffectingStatement(ctx, compiled, bindings...)
}

// sortedValues returns map values in a deterministic key order.
func sortedValues(m map[string]any) []any {
	keys := sortedKeys(m)
	vals := make([]any, 0, len(keys))
	for _, k := range keys {
		vals = append(vals, m[k])
	}
	return vals
}

// sortedKeys returns map keys in sorted order for deterministic SQL generation.
func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	// Simple insertion sort — maps are typically small.
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}
