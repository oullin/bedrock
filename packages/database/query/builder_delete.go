package query

import "context"

// Delete deletes rows matching the query.
func (b *Builder) Delete(ctx context.Context, id ...any) (int64, error) {
	if len(id) > 0 {
		b.Where("id", "=", id[0])
	}

	compiled := b.grammar.CompileDelete(b)
	return b.connection.Delete(ctx, compiled, b.GetRawBindings()[BindingWhere]...)
}

// Truncate truncates the table.
func (b *Builder) Truncate(ctx context.Context) error {
	statements := b.grammar.CompileTruncate(b)
	for sql, bindings := range statements {
		_, err := b.connection.Statement(ctx, sql, bindings)
		if err != nil {
			return err
		}
	}
	return nil
}
