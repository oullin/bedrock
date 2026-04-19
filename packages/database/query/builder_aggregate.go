package query

import (
	"context"
	"fmt"
)

// Count returns the count of rows matching the query.
func (b *Builder) Count(ctx context.Context, columns ...string) (int64, error) {
	if len(columns) == 0 {
		columns = []string{"*"}
	}

	return b.aggregateInt(ctx, "count", columns)
}

// Min returns the minimum value of a column.
func (b *Builder) Min(ctx context.Context, column string) (any, error) {
	return b.aggregateValue(ctx, "min", []string{column})
}

// Max returns the maximum value of a column.
func (b *Builder) Max(ctx context.Context, column string) (any, error) {
	return b.aggregateValue(ctx, "max", []string{column})
}

// Sum returns the sum of a column.
func (b *Builder) Sum(ctx context.Context, column string) (float64, error) {
	val, err := b.aggregateValue(ctx, "sum", []string{column})

	if err != nil {
		return 0, err
	}

	return toFloat64(val), nil
}

// Avg returns the average of a column.
func (b *Builder) Avg(ctx context.Context, column string) (float64, error) {
	val, err := b.aggregateValue(ctx, "avg", []string{column})

	if err != nil {
		return 0, err
	}

	return toFloat64(val), nil
}

// Average is an alias for Avg.
func (b *Builder) Average(ctx context.Context, column string) (float64, error) {
	return b.Avg(ctx, column)
}

// Exists checks if any rows match the query.
func (b *Builder) Exists(ctx context.Context) (bool, error) {
	compiled := b.grammar.CompileExists(b)
	results, err := b.connection.Select(ctx, compiled, b.GetBindings()...)

	if err != nil {
		return false, err
	}

	if len(results) == 0 {
		return false, nil
	}

	val, ok := results[0]["exists"]

	if !ok {
		// Some drivers may use different column names.
		for _, v := range results[0] {
			val = v

			break
		}
	}

	return toBool(val), nil
}

// DoesntExist checks if no rows match the query.
func (b *Builder) DoesntExist(ctx context.Context) (bool, error) {
	exists, err := b.Exists(ctx)

	return !exists, err
}

func (b *Builder) aggregateInt(ctx context.Context, fn string, columns []string) (int64, error) {
	val, err := b.aggregateValue(ctx, fn, columns)

	if err != nil {
		return 0, err
	}

	return toInt64(val), nil
}

func (b *Builder) aggregateValue(ctx context.Context, fn string, columns []string) (any, error) {
	clone := b.CloneWithout("columns")
	clone.SetAggregate(fn, columns)
	clone.columns = nil

	results, err := clone.Get(ctx)

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	return results[0]["aggregate"], nil
}

func toInt64(v any) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case float64:
		return int64(n)
	case string:
		var i int64

		fmt.Sscanf(n, "%d", &i)

		return i
	default:
		return 0
	}
}

func toFloat64(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int64:
		return float64(n)
	case int:
		return float64(n)
	case string:
		var f float64

		fmt.Sscanf(n, "%f", &f)

		return f
	default:
		return 0
	}
}

func toBool(v any) bool {
	switch n := v.(type) {
	case bool:
		return n
	case int64:
		return n != 0
	case int:
		return n != 0
	case float64:
		return n != 0
	case string:
		return n == "1" || n == "true"
	default:
		return false
	}
}
