package query

import (
	"context"
	"strings"
)

// OrderBy adds an ORDER BY clause.
func (b *Builder) OrderBy(column string, direction ...string) *Builder {
	dir := "asc"

	if len(direction) > 0 {
		dir = strings.ToLower(direction[0])
	}

	if dir != "asc" && dir != "desc" {
		dir = "asc"
	}

	b.orders = append(b.orders, OrderClause{Column: column, Direction: dir})

	return b
}

// OrderByDesc adds a descending ORDER BY clause.
func (b *Builder) OrderByDesc(column string) *Builder {
	return b.OrderBy(column, "desc")
}

// OrderByRaw adds a raw ORDER BY expression.
func (b *Builder) OrderByRaw(sql string, bindings ...any) *Builder {
	b.orders = append(b.orders, OrderClause{SQL: sql})
	b.AddBinding(BindingOrder, bindings...)

	return b
}

// Latest orders by the given column descending (default: "created_at").
func (b *Builder) Latest(column ...string) *Builder {
	col := "created_at"

	if len(column) > 0 {
		col = column[0]
	}

	return b.OrderByDesc(col)
}

// Oldest orders by the given column ascending (default: "created_at").
func (b *Builder) Oldest(column ...string) *Builder {
	col := "created_at"

	if len(column) > 0 {
		col = column[0]
	}

	return b.OrderBy(col, "asc")
}

// InRandomOrder orders by random.
func (b *Builder) InRandomOrder(seed ...string) *Builder {
	s := ""

	if len(seed) > 0 {
		s = seed[0]
	}

	b.orders = append(b.orders, OrderClause{SQL: b.grammar.CompileRandom(s)})

	return b
}

// Limit sets the maximum number of rows to return.
func (b *Builder) Limit(n int) *Builder {
	if n >= 0 {
		b.limit_ = &n
	} else {
		b.limit_ = nil
	}

	return b
}

// Take is an alias for Limit.
func (b *Builder) Take(n int) *Builder {
	return b.Limit(n)
}

// Offset sets the number of rows to skip.
func (b *Builder) Offset(n int) *Builder {
	if n >= 0 {
		b.offset_ = &n
	} else {
		b.offset_ = nil
	}

	return b
}

// Skip is an alias for Offset.
func (b *Builder) Skip(n int) *Builder {
	return b.Offset(n)
}

// ForPage sets offset/limit for a given page number and page size.
func (b *Builder) ForPage(page, perPage int) *Builder {
	return b.Offset((page - 1) * perPage).Limit(perPage)
}

// ForPageBeforeId constrains results to rows with an id less than the given value.
func (b *Builder) ForPageBeforeId(perPage int, lastId any, column ...string) *Builder {
	col := "id"

	if len(column) > 0 {
		col = column[0]
	}

	b.orders = nil
	b.bindings[BindingOrder] = nil

	if lastId != nil {
		return b.Where(col, "<", lastId).OrderByDesc(col).Limit(perPage)
	}

	return b.OrderByDesc(col).Limit(perPage)
}

// ForPageAfterId constrains results to rows with an id greater than the given value.
func (b *Builder) ForPageAfterId(perPage int, lastId any, column ...string) *Builder {
	col := "id"

	if len(column) > 0 {
		col = column[0]
	}

	b.orders = nil
	b.bindings[BindingOrder] = nil

	if lastId != nil {
		return b.Where(col, ">", lastId).OrderBy(col).Limit(perPage)
	}

	return b.OrderBy(col).Limit(perPage)
}

// InOrderOf orders by matching against a specific list of values.
func (b *Builder) InOrderOf(column string, values []any) *Builder {
	if len(values) == 0 {
		return b
	}
	// Generates: ORDER BY FIELD(column, val1, val2, ...) or CASE WHEN expression.
	raw := "field(" + b.grammar.Wrap(column) + ", " + b.grammar.Parameterize(values) + ")"
	b.orders = append(b.orders, OrderClause{SQL: raw})
	b.AddBinding(BindingOrder, values...)

	return b
}

// GroupLimit sets a limit per group (for use with GROUP BY).
func (b *Builder) GroupLimit(n int) *Builder {
	return b.Limit(n)
}

// Cursor returns results one at a time for memory efficiency. It calls the
// callback for each row. Returns on first error or when callback returns false.
func (b *Builder) Cursor(ctx context.Context, fn func(map[string]any) bool) error {
	return b.Chunk(ctx, 100, func(rows []map[string]any, _ int) bool {
		for _, row := range rows {
			if !fn(row) {
				return false
			}
		}

		return true
	})
}

// DD dumps the SQL and bindings and returns them. Go equivalent of PHP's dd().
func (b *Builder) DD() string {
	return b.Dump()
}
